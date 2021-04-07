///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

// Contains functionality for describing and creating connections

package connect

import (
	"fmt"
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/xx_network/comms/connect/token"
	"gitlab.com/xx_network/crypto/signature/rsa"
	tlsCreds "gitlab.com/xx_network/crypto/tls"
	"gitlab.com/xx_network/primitives/id"
	"gitlab.com/xx_network/primitives/rateLimiting"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
	"math"
	"net"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

const infinityTime = time.Duration(math.MaxInt64)

// KaClientOpts are the keepalive options for clients
// TODO: Set via configuration
var KaClientOpts = keepalive.ClientParameters{
	// Never ping to keepalive
	Time: infinityTime,
	// 60s after ping before closing
	Timeout: 60 * time.Second,
	// For all connections, streaming and nonstreaming
	PermitWithoutStream: true,
}

// Information used to describe a connection to a host
type Host struct {
	// System-wide ID of the Host
	id *id.ID

	// address:Port being connected to
	addressAtomic atomic.Value

	// PEM-format TLS Certificate
	certificate []byte

	/* Tokens shared with this Host establishing reverse authentication */

	//  Live used for receiving from this host
	receptionToken *token.Live

	// Live used for sending to this host
	transmissionToken *token.Live

	// GRPC connection object
	connection      *grpc.ClientConn
	connectionCount uint64

	// TLS credentials object used to establish the connection
	credentials credentials.TransportCredentials

	// RSA Public Key corresponding to the TLS Certificate
	rsaPublicKey *rsa.PublicKey

	// Indicates whether dynamic authentication was used for this Host
	// This is useful for determining whether a Host's key was hardcoded
	dynamicHost bool

	// State tracking for host metric
	metrics *Metric

	// Send lock
	sendMux sync.RWMutex

	coolOffBucket *rateLimiting.Bucket
	inCoolOff     bool

	// Stored default values (should be non-mutated)
	params HostParams
}

// Creates a new Host object
func NewHost(id *id.ID, address string, cert []byte, params HostParams) (host *Host, err error) {

	// Initialize the Host object
	host = &Host{
		id:                id,
		certificate:       cert,
		transmissionToken: token.NewLive(),
		receptionToken:    token.NewLive(),
		metrics:           newMetric(),
		params:            params,
	}

	if params.EnableCoolOff {
		host.coolOffBucket = rateLimiting.CreateBucket(
			params.NumSendsBeforeCoolOff+1, params.NumSendsBeforeCoolOff+1,
			params.CoolOffTimeout, nil)
	}

	if host.params.MaxRetries == 0 {
		host.params.MaxRetries = math.MaxUint32
	}

	host.UpdateAddress(address)

	// Configure the host credentials
	err = host.setCredentials()

	//print logs
	jww.INFO.Printf("New Host Created: %s", host)
	jww.TRACE.Printf("New Host Certificate for %v: %s...", id, cert)
	return
}

// Creates a new dynamic-authenticated Host object
func newDynamicHost(id *id.ID, publicKey []byte) (host *Host, err error) {

	// Initialize the Host object
	// IMPORTANT: This flag must be set to true for all dynamic Hosts
	//            because the security properties for these Hosts differ
	host = &Host{
		id:                id,
		dynamicHost:       true,
		transmissionToken: token.NewLive(),
		receptionToken:    token.NewLive(),
	}

	// Create the RSA Public Key object
	host.rsaPublicKey, err = rsa.LoadPublicKeyFromPem(publicKey)
	if err != nil {
		err = errors.Errorf("Error extracting PublicKey: %+v", err)
	}
	return
}

// Simple getter for the dynamicHost value
func (h *Host) IsDynamicHost() bool {
	return h.dynamicHost
}

// Simple getter for the public key
func (h *Host) GetPubKey() *rsa.PublicKey {
	return h.rsaPublicKey
}

// Connected checks if the given Host's connection is alive
// the uint is the connection count, it increments every time a reconnect occurs
func (h *Host) Connected() (bool, uint64) {
	h.sendMux.RLock()
	defer h.sendMux.RUnlock()

	return h.isAlive() && !h.authenticationRequired(), h.connectionCount
}

// GetId returns the id of the host
func (h *Host) GetId() *id.ID {
	if h == nil {
		return &id.ID{}
	}
	return h.id
}

// GetAddress returns the address of the host.
func (h *Host) GetAddress() string {
	a := h.addressAtomic.Load()
	if a == nil {
		return ""
	}
	return a.(string)
}

// UpdateAddress updates the address of the host
func (h *Host) UpdateAddress(address string) {
	h.addressAtomic.Store(address)
}

// GetMetrics returns a deep copy of Host's Metric
// This resets the state of metrics
func (h *Host) GetMetrics() *Metric {
	return h.metrics.get()
}

// isExcludedMetricError determines if err is within the list
// of excludeMetricErrors.  Returns true if it's an excluded error,
// false if it is not
func (h *Host) isExcludedMetricError(err string) bool {
	for _, excludedErr := range h.params.ExcludeMetricErrors {
		if strings.Contains(excludedErr, err) {
			return true
		}
	}
	return false
}

// Sets the host metrics to an arbitrary value. Used for testing
// purposes only
func (h *Host) SetMetricsTesting(m *Metric, face interface{}) {
	// Ensure that this function is only run in testing environments
	switch face.(type) {
	case *testing.T, *testing.M, *testing.B:
		break
	default:
		panic("SetMetricsTesting() can only be used for testing.")
	}

	h.metrics = m

}

// Disconnect closes a the Host connection under the write lock
func (h *Host) Disconnect() {
	h.sendMux.Lock()
	defer h.sendMux.Unlock()

	h.disconnect()
}

// ConditionalDisconnect closes a the Host connection under the write lock only
// if the connection count has not increased
func (h *Host) conditionalDisconnect(count uint64) {
	h.sendMux.Lock()
	defer h.sendMux.Unlock()

	if count == h.connectionCount {
		h.disconnect()
	}
}

// Returns whether or not the Host is able to be contacted
// by attempting to dial a tcp connection
func (h *Host) IsOnline() bool {
	addr := h.GetAddress()
	timeout := 5 * time.Second
	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		// If we cannot connect, mark the node as failed
		jww.DEBUG.Printf("Failed to verify connectivity for address %s", addr)
		return false
	}
	// Attempt to close the connection
	if conn != nil {
		errClose := conn.Close()
		if errClose != nil {
			jww.DEBUG.Printf("Failed to close connection for address %s",
				addr)
		}
	}
	return true
}

// send checks that the host has a connection and sends if it does.
// Operates under the host's read lock.
func (h *Host) transmit(f func(conn *grpc.ClientConn) (interface{},
	error)) (interface{}, error) {

	h.sendMux.RLock()
	defer h.sendMux.RUnlock()

	// Check if connection is down
	if h.connection == nil {
		return nil, errors.New("Failed to transmit: host disconnected")
	}

	a, err := f(h.connection)

	if h.params.EnableMetrics && err != nil {
		// Checks if the received error is a among excluded errors
		// If it is not an excluded error, update host's metrics
		if !h.isExcludedMetricError(err.Error()) {
			h.metrics.incrementErrors()
		}
	}

	return a, err
}

// connect attempts to connect to the host if it does not have a valid connection
func (h *Host) connect() error {

	//connect to remote
	if err := h.connectHelper(); err != nil {
		return err
	}

	h.connectionCount++

	return nil
}

// authenticationRequired Checks if new authentication is required with
// the remote.  This is used exclusively under the lock in protocoms.transmit so
// no lock is needed
func (h *Host) authenticationRequired() bool {
	return h.params.AuthEnabled && !h.transmissionToken.Has()
}

// isAlive returns true if the connection is non-nil and alive
func (h *Host) isAlive() bool {
	if h.connection == nil {
		return false
	}
	state := h.connection.GetState()
	return state == connectivity.Idle || state == connectivity.Connecting ||
		state == connectivity.Ready
}

// disconnect closes a the Host connection while not under a write lock.
// undefined behavior if the caller has not taken the write lock
func (h *Host) disconnect() {
	// its possible to close a host which never sent so it never made a
	// connection. In that case, we should not close a connection which does not
	// exist
	if h.connection != nil {
		err := h.connection.Close()
		if err != nil {
			jww.ERROR.Printf("Unable to close connection to %s: %+v",
				h.GetAddress(), errors.New(err.Error()))
		} else {
			h.connection = nil
		}
	}
	h.transmissionToken.Clear()
}

// connectHelper creates a connection while not under a write lock.
// undefined behavior if the caller has not taken the write lock
func (h *Host) connectHelper() (err error) {

	// Configure TLS options
	var securityDial grpc.DialOption
	if h.credentials != nil {
		// Create the gRPC client with TLS
		securityDial = grpc.WithTransportCredentials(h.credentials)
	} else {
		// Create the gRPC client without TLS
		jww.WARN.Printf("Connecting to %v without TLS!", h.GetAddress())
		securityDial = grpc.WithInsecure()
	}

	jww.DEBUG.Printf("Attempting to establish connection to %s using"+
		" credentials: %+v", h.GetAddress(), securityDial)

	// Attempt to establish a new connection
	var numRetries uint32
	//todo-remove this retry block when grpc is updated
	for numRetries = 0; numRetries < h.params.MaxRetries && !h.isAlive(); numRetries++ {
		h.disconnect()

		jww.DEBUG.Printf("Connecting to %+v Attempt number %+v of %+v",
			h.GetAddress(), numRetries, h.params.MaxRetries)

		// If timeout is enabled, the max wait time becomes
		// ~14 seconds (with maxRetries=100)
		backoffTime := 2000 * (numRetries/16 + 1)
		if backoffTime > 15000 {
			backoffTime = 15000
		}
		ctx, cancel := ConnectionContext(time.Duration(backoffTime) * time.Millisecond)

		// Create the connection
		h.connection, err = grpc.DialContext(ctx, h.GetAddress(),
			securityDial,
			grpc.WithBlock(),
			grpc.WithKeepaliveParams(KaClientOpts))
		if err != nil {
			jww.DEBUG.Printf("Attempt number %+v to connect to %s failed\n",
				numRetries, h.GetAddress())
		}
		cancel()
	}

	// Verify that the connection was established successfully
	if !h.isAlive() {
		h.disconnect()
		return errors.New(fmt.Sprintf(
			"Last try to connect to %s failed. Giving up",
			h.GetAddress()))
	}

	// Add the successful connection to the Manager
	jww.INFO.Printf("Successfully connected to %v", h.GetAddress())
	return
}

// setCredentials sets TransportCredentials and RSA PublicKey objects
// using a PEM-encoded TLS Certificate
func (h *Host) setCredentials() error {

	// If no TLS Certificate specified, print a warning and do nothing
	if h.certificate == nil || len(h.certificate) == 0 {
		jww.WARN.Printf("No TLS Certificate specified!")
		return nil
	}

	// Obtain the DNS name included with the certificate
	dnsName := ""
	cert, err := tlsCreds.LoadCertificate(string(h.certificate))
	if err != nil {
		return errors.Errorf("Error forming transportCredentials: %+v", err)
	}
	if len(cert.DNSNames) > 0 {
		dnsName = cert.DNSNames[0]
	}

	// Create the TLS Credentials object
	h.credentials, err = tlsCreds.NewCredentialsFromPEM(string(h.certificate),
		dnsName)
	if err != nil {
		return errors.Errorf("Error forming transportCredentials: %+v", err)
	}

	// Create the RSA Public Key object
	h.rsaPublicKey, err = tlsCreds.NewPublicKeyFromPEM(h.certificate)
	if err != nil {
		err = errors.Errorf("Error extracting PublicKey: %+v", err)
	}

	return err
}

// Stringer interface for connection
func (h *Host) String() string {
	h.sendMux.RLock()
	defer h.sendMux.RUnlock()
	addr := h.GetAddress()

	return fmt.Sprintf(
		"ID: %v\tAddr: %v",
		h.id, addr)
}

// Stringer interface for connection
func (h *Host) StringVerbose() string {
	return fmt.Sprintf("%s\t CERTIFICATE: %s", h, h.certificate)
}

func (h *Host) SetTestPublicKey(key *rsa.PublicKey, t interface{}) {
	switch t.(type) {
	case *testing.T:
		break
	case *testing.M:
		break
	default:
		jww.FATAL.Panicf("SetTestPublicKey is restricted to testing only. Got %T", t)
	}
	h.rsaPublicKey = key
}

// Set host to dynamic (for testing use)
func (h *Host) SetTestDynamic(t interface{}) {
	switch t.(type) {
	case *testing.T:
		break
	case *testing.M:
		break
	default:
		jww.FATAL.Panicf("SetTestDynamic is restricted to testing only. Got %T", t)
	}
	h.dynamicHost = true
}
