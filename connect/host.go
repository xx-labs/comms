////////////////////////////////////////////////////////////////////////////////
// Copyright © 2024 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

// Contains functionality for describing and creating connections

package connect

import (
	"context"
	"crypto/x509"
	"fmt"
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/xx_network/comms/connect/token"
	"gitlab.com/xx_network/crypto/signature/rsa"
	tlsCreds "gitlab.com/xx_network/crypto/tls"
	"gitlab.com/xx_network/primitives/exponential"
	"gitlab.com/xx_network/primitives/id"
	"gitlab.com/xx_network/primitives/rateLimiting"
	"google.golang.org/grpc/credentials"
	"math"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// Host information used to describe a remote connection
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
	connection      Connection
	connectionCount uint64
	// lock which ensures only a single thread is connecting at a time and
	// that connections do not interrupt sends
	connectionMux sync.RWMutex

	// TLS credentials object used to establish the connection
	credentials credentials.TransportCredentials

	// RSA Public Key corresponding to the TLS Certificate
	rsaPublicKey *rsa.PublicKey

	// State tracking for host metric
	metrics *Metric

	// Tracks the exponential moving average of proxy error messages so that if
	// too many connection error occur, the layer above can be informed
	proxyErrorMetric *exponential.MovingAvg

	coolOffBucket *rateLimiting.Bucket
	inCoolOff     bool

	// Stored default values (should be non-mutated)
	params HostParams

	// the amount of data, when streaming, that a sender can send before receiving an ACK
	// keep at zero to use the default GRPC algorithm to determine
	windowSize *int32
}

// NewHost creates a new host object which will use GRPC.
func NewHost(id *id.ID, address string, cert []byte,
	params HostParams) (host *Host, err error) {
	return newHost(id, address, cert, params)
}

// newHost is a helper which creates a new Host object
func newHost(id *id.ID, address string, cert []byte, params HostParams) (host *Host, err error) {

	windowSize := int32(0)

	// Initialize the Host object
	host = &Host{
		id:                id,
		certificate:       cert,
		transmissionToken: token.NewLive(),
		receptionToken:    token.NewLive(),
		metrics:           newMetric(),
		proxyErrorMetric:  exponential.NewMovingAvg(params.ProxyErrorMetricParams),
		params:            params,
		windowSize:        &windowSize,
	}

	host.connection = newConnection(params.ConnectionType, host)

	if params.EnableCoolOff {
		host.coolOffBucket = rateLimiting.CreateBucket(
			params.NumSendsBeforeCoolOff+1, params.NumSendsBeforeCoolOff+1,
			params.CoolOffTimeout, nil)
	}

	if host.params.MaxRetries == 0 {
		host.params.MaxRetries = math.MaxUint32
	}

	host.UpdateAddress(address)

	// Configure the transport credentials for GRPC hosts
	//if !host.IsWeb() {
	err = host.setCredentials()
	if err != nil {
		return
	}
	//}

	// Connect immediately if configured to do so
	if params.DisableLazyConnection {
		// No mutex required
		err = host.connect()
	}
	return
}

// SetWindowSize sets the amount of data, when streaming, that a sender can send before receiving an ACK
// keep at zero to use the default GRPC algorithm to determine
func (h *Host) SetWindowSize(size int32) {
	atomic.StoreInt32(h.windowSize, size)
}

// GetPubKey simple getter for the public key
func (h *Host) GetPubKey() *rsa.PublicKey {
	return h.rsaPublicKey
}

// Connected checks if the given Host's connection is alive
// the uint is the connection count, it increments every time a reconnect occurs
func (h *Host) Connected() (bool, uint64) {
	h.connectionMux.RLock()
	defer h.connectionMux.RUnlock()

	return h.connectedUnsafe()
}

// connectedUnsafe checks if the given Host's connection is alive without taking
// a connection lock. Only use if already under a connection lock. The uint is
// the connection count, it increments every time a reconnect occurs
func (h *Host) connectedUnsafe() (bool, uint64) {
	return h.connection != nil && h.connection.isAlive() && !h.authenticationRequired(), h.connectionCount
}

// GetMessagingContext returns a context object for message sending configured according to HostParams
func (h *Host) GetMessagingContext() (context.Context, context.CancelFunc) {
	return h.GetMessagingContextWithTimeout(h.params.SendTimeout)
}

// GetMessagingContextWithTimeout returns a context object for message sending configured according to HostParams
func (h *Host) GetMessagingContextWithTimeout(
	timeout time.Duration) (context.Context, context.CancelFunc) {
	return newContext(timeout)
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

// IsWeb returns the connection type of the host
func (h *Host) IsWeb() bool {
	return h.connection.IsWeb()
}

// GetRemoteCertificate returns the tls certificate from the server for web hosts
// Note that this will return an error when used on grpc hosts, and will not
// have a certificate ready until something has been sent over the connection.
func (h *Host) GetRemoteCertificate() (*x509.Certificate, error) {
	return h.connection.GetRemoteCertificate()
}

// SetMetricsTesting sets the host metrics to an arbitrary value. Used for testing
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

// Disconnect closes the Host connection under the write lock
// Due to asynchronous connection handling, this may result in
// killing a good connection and could result in an immediate
// reconnection by a separate thread
func (h *Host) Disconnect() {
	h.connectionMux.Lock()
	defer h.connectionMux.Unlock()
	h.disconnect()
}

// ConditionalDisconnect closes the Host connection under the write lock only
// if the connection count has not increased
func (h *Host) conditionalDisconnect(count uint64) {
	if count == h.connectionCount {
		h.disconnect()
	}
}

// IsOnline returns whether the Host is able to be contacted
// before the timeout by attempting to dial a tcp connection
// Returns how long the ping took, and whether it was successful
func (h *Host) IsOnline() (time.Duration, bool) {
	return h.connection.IsOnline()
}

// send checks that the host has a connection and sends if it does.
// must be called under host's connection read lock.
func (h *Host) transmit(f func(conn Connection) (interface{},
	error)) (interface{}, error) {

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

	if err != nil {
		// Check if the received error is a connection timeout and add it to the
		// moving average. If the cutoff is reached for too many timeouts,
		// return TooManyProxyError instead so that the host can be removed from
		// the host pool on the layer above.
		err2 := h.proxyErrorMetric.Intake(
			exponential.BoolToFloat(strings.Contains(err.Error(), ProxyError)))
		if err2 != nil {
			err = errors.Errorf("%s: %+v", TooManyProxyError, err2)
		}
	}

	return a, err
}

// Connect allows manual connection to the host if it does not have a valid connection
func (h *Host) Connect() error {
	h.connectionMux.Lock()
	defer h.connectionMux.Unlock()

	return h.connect()
}

// connect attempts to connect to the host if it does not have a valid connection
func (h *Host) connect() error {
	//connect to remote
	if err := h.connection.Connect(); err != nil {
		return err
	}

	h.connectionCount++

	return nil
}

// authenticationRequired Checks if new authentication is required with
// the remote.  This is used exclusively under the lock in protocomm.transmit so
// no lock is needed
func (h *Host) authenticationRequired() bool {
	return h.params.AuthEnabled && !h.transmissionToken.Has()
}

// isAlive returns true if the connection is non-nil and alive
// must already be under the connectionMux
func (h *Host) isAlive() bool {
	if h.connection == nil {
		return false
	}
	return h.connection.isAlive()
}

// disconnect closes the Host connection while not under a write lock.
// undefined behavior if the caller has not taken the write lock
func (h *Host) disconnect() {
	h.connection.disconnect()
	h.transmissionToken.Clear()
}

// setCredentials sets GRPC TransportCredentials and RSA PublicKey objects
// using a PEM-encoded TLS Certificate
func (h *Host) setCredentials() error {

	// If no TLS Certificate specified, print a warning and do nothing
	if h.certificate == nil || len(h.certificate) == 0 {
		if TestingOnlyDisableTLS {
			jww.WARN.Printf("No TLS Certificate specified!")
			return nil
		} else {
			jww.FATAL.Panicf("TLS cannot be disabled in production, only for testing suites!")
		}
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
	h.connectionMux.RLock()
	defer h.connectionMux.RUnlock()
	addr := h.GetAddress()

	return fmt.Sprintf(
		"ID: %v\tAddr: %v",
		h.id, addr)
}

// StringVerbose stringer interface for connection
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
