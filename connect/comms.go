///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

// Handles the basic top-level comms object used across all packages

package connect

import (
	"crypto/tls"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/xx_network/comms/connect/token"
	"gitlab.com/xx_network/crypto/signature/rsa"
	"gitlab.com/xx_network/primitives/id"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
	"math"
	"net"
	"strings"
	"time"
)

// TODO: Set these via config

// KaOpts are Keepalive options for servers
var KaOpts = keepalive.ServerParameters{
	// Idle for at most 60s
	MaxConnectionIdle: 60 * time.Second,
	// Reset after an hour
	MaxConnectionAge: 1 * time.Hour,
	// w/ 1m grace shutdown
	MaxConnectionAgeGrace: 1 * time.Minute,
	// Never ping to keepalive
	Time: infinityTime,
	// Close connection 60 seconds after ping
	Timeout: 60 * time.Second,
}

// KaEnforcement are keepalive enforcement options for servers
var KaEnforcement = keepalive.EnforcementPolicy{
	// Client should never send keep alive ping
	MinTime: infinityTime,
	// Doing KA on non-streams is OK
	PermitWithoutStream: true,
}

// MaxConcurrentStreams is the number of server-side streams to allow open
var MaxConcurrentStreams = uint32(250000)

// Proto object containing a gRPC server
type ProtoComms struct {
	// Inherit the Manager object
	*Manager

	// The network ID of this comms server
	Id *id.ID

	// Private key of the local comms instance
	privateKey *rsa.PrivateKey

	// Disables the checking of authentication signatures for testing setups
	disableAuth bool

	// SERVER-ONLY FIELDS ------------------------------------------------------

	// A map of reverse-authentication tokens
	tokens *token.Map

	// Local network server
	LocalServer *grpc.Server

	// Listening address of the local server
	ListeningAddr string

	// CLIENT-ONLY FIELDS ------------------------------------------------------

	// Used to store the public key used for generating Client Id
	pubKeyPem []byte

	// Used to store the salt used for generating Client Id
	salt []byte

	// -------------------------------------------------------------------------
}

// Creates a ProtoComms client-type object to be used in various initializers
func CreateCommClient(id *id.ID, pubKeyPem, privKeyPem,
	salt []byte) (*ProtoComms, error) {
	// Build the ProtoComms object
	pc := &ProtoComms{
		Id:        id,
		pubKeyPem: pubKeyPem,
		salt:      salt,
		tokens:    token.NewMap(),
		Manager:   newManager(),
	}

	// Set the private key if specified
	if privKeyPem != nil {
		err := pc.setPrivateKey(privKeyPem)
		if err != nil {
			return nil, errors.Errorf("Could not set private key: %+v", err)
		}
	}
	return pc, nil
}

// Creates a ProtoComms server-type object to be used in various initializers
func StartCommServer(id *id.ID, localServer string, certPEMblock,
	keyPEMblock []byte) (*ProtoComms, net.Listener, error) {

	// Build the ProtoComms object
	pc := &ProtoComms{
		Id:            id,
		ListeningAddr: localServer,
		tokens:        token.NewMap(),
		Manager:       newManager(),
	}

listen:
	// Listen on the given address
	lis, err := net.Listen("tcp", localServer)
	if err != nil {
		if strings.Contains(err.Error(), "bind: address already in use") {
			jww.WARN.Printf("Could not listen on %s, is port in use? waiting 30s: %s", localServer, err.Error())
			time.Sleep(30 * time.Second)
			goto listen
		}
		return nil, nil, errors.New(err.Error())
	}

	// If TLS was specified
	if certPEMblock != nil && keyPEMblock != nil {

		// Create the TLS certificate
		x509cert, err := tls.X509KeyPair(certPEMblock, keyPEMblock)
		if err != nil {
			return nil, nil, errors.Errorf("Could not load TLS keys: %+v", err)
		}

		// Set the private key
		err = pc.setPrivateKey(keyPEMblock)
		if err != nil {
			return nil, nil, errors.Errorf("Could not set private key: %+v", err)
		}

		// Create the gRPC server with TLS
		jww.INFO.Printf("Starting server with TLS...")
		creds := credentials.NewServerTLSFromCert(&x509cert)
		pc.LocalServer = grpc.NewServer(grpc.Creds(creds),
			grpc.MaxConcurrentStreams(MaxConcurrentStreams),
			grpc.MaxRecvMsgSize(math.MaxInt32),
			grpc.KeepaliveParams(KaOpts),
			grpc.KeepaliveEnforcementPolicy(KaEnforcement))
	} else {
		// Create the gRPC server without TLS
		jww.WARN.Printf("Starting server with TLS disabled...")
		pc.LocalServer = grpc.NewServer(
			grpc.MaxConcurrentStreams(MaxConcurrentStreams),
			grpc.MaxRecvMsgSize(math.MaxInt32),
			grpc.KeepaliveParams(KaOpts),
			grpc.KeepaliveEnforcementPolicy(KaEnforcement))
	}

	return pc, lis, nil
}

// Performs a graceful shutdown of the local server
func (c *ProtoComms) Shutdown() {
	c.DisconnectAll()
	c.LocalServer.GracefulStop()
	time.Sleep(time.Millisecond * 500)
}

// Stringer method
func (c *ProtoComms) String() string {
	return c.ListeningAddr
}

// Setter for local server's private key
func (c *ProtoComms) setPrivateKey(data []byte) error {
	key, err := rsa.LoadPrivateKeyFromPem(data)
	if err != nil {
		return errors.Errorf("Failed to form private key file from data at %s: %+v", data, err)
	}

	c.privateKey = key
	return nil
}

// Getter for local server's private key
func (c *ProtoComms) GetPrivateKey() *rsa.PrivateKey {
	return c.privateKey
}

const (
	//numerical designators for 3 operations of send, used to ensure the same
	//operation isn't repeated
	con  = 1
	auth = 2
	send = 3
)

// Sets up or recovers the Host's connection
// Then runs the given Send function
func (c *ProtoComms) Send(host *Host, f func(conn *grpc.ClientConn) (*any.Any,
	error)) (result *any.Any, err error) {

	jww.TRACE.Printf("Attempting to send to host: %s", host)
	fSh := func(conn *grpc.ClientConn) (interface{}, error) {
		return f(conn)
	}

	anyFace, err := c.transmit(host, fSh)
	if err != nil {
		return nil, err
	}

	return anyFace.(*any.Any), err
}

// Sets up or recovers the Host's connection
// Then runs the given Stream function
func (c *ProtoComms) Stream(host *Host, f func(conn *grpc.ClientConn) (
	interface{}, error)) (client interface{}, err error) {

	// Ensure the connection is running
	jww.TRACE.Printf("Attempting to stream to host: %s", host)
	return c.transmit(host, f)
}

// returns true if the connection error is one of the connection errors which
// should be retried
func isConnError(err error) bool {
	return strings.Contains(err.Error(), "context deadline exceeded") ||
		strings.Contains(err.Error(), "connection refused") ||
		strings.Contains(err.Error(), "host disconnected")
}
