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
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/pkg/errors"
	"github.com/soheilhy/cmux"
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/xx_network/comms/connect/token"
	"gitlab.com/xx_network/crypto/signature/rsa"
	"gitlab.com/xx_network/primitives/id"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
	"math"
	"net"
	"net/http"
	"strings"
	"time"
)

// MaxWindowSize 4 MB
const MaxWindowSize = math.MaxInt32

// TestingOnlyDisableTLS is the variable set for testing
// which allows for the disabled TLS code-path. Production
// code-path will only function with TLS enabled.
var TestingOnlyDisableTLS = false

// KaOpts are Keepalive options for servers
// TODO: Set these via config
var KaOpts = keepalive.ServerParameters{
	// Idle for at most 60s
	MaxConnectionIdle: 60 * time.Second,
	// Reset after an hour
	MaxConnectionAge: 1 * time.Hour,
	// w/ 1m grace shutdown
	MaxConnectionAgeGrace: 1 * time.Minute,
	// Send keepAlive every Time interval
	Time: 5 * time.Second,
	// Timeout after last successful keepAlive to close connection
	Timeout: 60 * time.Second,
}

// KaEnforcement are keepalive enforcement options for servers
var KaEnforcement = keepalive.EnforcementPolicy{
	// Send keepAlive every Time interval
	MinTime: 3 * time.Second,
	// Doing KA on non-streams is OK
	PermitWithoutStream: true,
}

// MaxConcurrentStreams is the number of server-side streams to allow open
var MaxConcurrentStreams = uint32(250000)

// ProtoComms is a proto object containing a gRPC server logic.
type ProtoComms struct {
	// Inherit the Manager object
	*Manager

	// The network ID of this comms server
	receptionId *id.ID

	// Private key of the local comms instance
	privateKey *rsa.PrivateKey

	// Disables the checking of authentication signatures for testing setups
	disableAuth bool

	// SERVER-ONLY FIELDS ------------------------------------------------------

	// A map of reverse-authentication tokens
	tokens *token.Map

	// Low-level net.Listener object that listens at listeningAddress
	netListener net.Listener

	// Local network server
	grpcServer *grpc.Server

	// Listening address of the local server
	listeningAddress string

	// CLIENT-ONLY FIELDS ------------------------------------------------------

	// Used to store the public key used for generating Client Id
	pubKeyPem []byte

	// Used to store the salt used for generating Client Id
	salt []byte

	// -------------------------------------------------------------------------
}

func (c *ProtoComms) GetId() *id.ID {
	return c.receptionId.DeepCopy()
}

func (c *ProtoComms) GetListener() net.Listener {
	return c.netListener
}

func (c *ProtoComms) GetServer() *grpc.Server {
	return c.grpcServer
}

func (c *ProtoComms) GetListeningAddress() string {
	return c.listeningAddress
}

// CreateCommClient creates a ProtoComms client-type object to be
// used in various initializers.
func CreateCommClient(id *id.ID, pubKeyPem, privKeyPem,
	salt []byte) (*ProtoComms, error) {
	// Build the ProtoComms object
	pc := &ProtoComms{
		receptionId: id,
		pubKeyPem:   pubKeyPem,
		salt:        salt,
		tokens:      token.NewMap(),
		Manager:     newManager(),
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

// StartCommServer creates a ProtoComms server-type object to be used in various initializers.
func StartCommServer(id *id.ID, localServer string,
	certPEMblock, keyPEMblock []byte, preloadedHosts []*Host) (*ProtoComms, error) {

listen:
	// Listen on the given address
	lis, err := net.Listen("tcp", localServer)
	if err != nil {
		if strings.Contains(err.Error(), "bind: address already in use") {
			jww.WARN.Printf("Could not listen on %s, is port in use? waiting 30s: %s", localServer, err.Error())
			time.Sleep(30 * time.Second)
			goto listen
		}
		return nil, errors.New(err.Error())
	}

	// Build the comms object
	pc := &ProtoComms{
		receptionId:      id,
		listeningAddress: localServer,
		netListener:      lis,
		tokens:           token.NewMap(),
		Manager:          newManager(),
	}

	for _, h := range preloadedHosts {
		pc.Manager.addHost(h)
	}

	// If TLS was specified
	if certPEMblock != nil && keyPEMblock != nil {

		// Create the TLS certificate
		x509cert, err := tls.X509KeyPair(certPEMblock, keyPEMblock)
		if err != nil {
			return nil, errors.Errorf("Could not load TLS keys: %+v", err)
		}

		// Set the private key
		err = pc.setPrivateKey(keyPEMblock)
		if err != nil {
			return nil, errors.Errorf("Could not set private key: %+v", err)
		}

		// Create the gRPC server with TLS
		jww.INFO.Printf("Starting server with TLS...")
		creds := credentials.NewServerTLSFromCert(&x509cert)
		pc.grpcServer = grpc.NewServer(grpc.Creds(creds),
			grpc.MaxConcurrentStreams(MaxConcurrentStreams),
			grpc.MaxRecvMsgSize(math.MaxInt32),
			grpc.KeepaliveParams(KaOpts),
			grpc.KeepaliveEnforcementPolicy(KaEnforcement))
	} else if TestingOnlyDisableTLS {
		// Create the gRPC server without TLS
		jww.WARN.Printf("Starting server with TLS disabled...")
		pc.grpcServer = grpc.NewServer(
			grpc.MaxConcurrentStreams(MaxConcurrentStreams),
			grpc.MaxRecvMsgSize(math.MaxInt32),
			grpc.KeepaliveParams(KaOpts),
			grpc.KeepaliveEnforcementPolicy(KaEnforcement))
	} else {
		jww.FATAL.Panicf("TLS cannot be disabled in production, only for testing suites!")
	}

	return pc, nil
}

// Serve is a non-blocking call that begins serving content
// for GRPC. GRPC endpoints must be registered before making this call.
func (c *ProtoComms) Serve() {
	listenGRPC := func(l net.Listener) {
		// This blocks for the lifetime of the listener.
		if err := c.GetServer().Serve(l); err != nil {
			jww.FATAL.Panicf("Failed to serve GRPC: %+v", err)
		}
		jww.INFO.Printf("Shutting down GRPC server listener")
	}
	go listenGRPC(c.GetListener())
}

// ServeWithWeb is a non-blocking call that begins serving content
// for grpcWeb (over HTTP) and GRPC on the same port.
// GRPC endpoints must be registered before making this call.
func (c *ProtoComms) ServeWithWeb() {
	netListener := c.GetListener()
	grpcServer := c.GetServer()

	// Split netListener into two distinct listeners for GRPC and HTTP
	mux := cmux.New(netListener)
	grpcL := mux.MatchWithWriters(
		cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc"))
	httpL := mux.Match(cmux.HTTP1Fast())

	jww.INFO.Printf("Started HTTP listener on GRPC endpoints: %+v",
		grpcweb.ListGRPCResources(grpcServer))

	listenHTTP := func(l net.Listener) {
		httpServer := grpcweb.WrapServer(grpcServer,
			grpcweb.WithOriginFunc(func(origin string) bool { return true }))
		// This blocks for the lifetime of the listener.
		if err := http.Serve(l, httpServer); err != nil {
			jww.FATAL.Panicf("Failed to serve HTTP: %v", err)
		}
		jww.INFO.Printf("Shutting down HTTP server listener")
	}
	listenGRPC := func(l net.Listener) {
		// This blocks for the lifetime of the listener.
		if err := grpcServer.Serve(l); err != nil {
			jww.FATAL.Panicf("Failed to serve GRPC: %+v", err)
		}
		jww.INFO.Printf("Shutting down GRPC server listener")
	}
	listenPort := func() {
		if err := mux.Serve(); err != nil {
			jww.FATAL.Panicf("Failed to serve port: %+v", err)
		}
		jww.INFO.Printf("Shutting down port server listener")
	}
	go listenGRPC(grpcL)
	go listenHTTP(httpL)
	go listenPort()
}

// Shutdown performs a graceful shutdown of the local server.
func (c *ProtoComms) Shutdown() {
	c.grpcServer.GracefulStop()
	c.DisconnectAll()
	if err := c.netListener.Close(); err != nil {
		jww.ERROR.Printf("Unable to close net listener: %+v", err)
		return
	}
	time.Sleep(100 * time.Millisecond)
	jww.INFO.Printf("Comms server successfully shut down")
}

// Stringer method
func (c *ProtoComms) String() string {
	return c.listeningAddress
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

// GetPrivateKey is the getter for local server's private key.
func (c *ProtoComms) GetPrivateKey() *rsa.PrivateKey {
	return c.privateKey
}

// Send sets up or recovers the Host's connection,
// then runs the given transmit function.
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

// Stream sets up or recovers the Host's connection,
// then runs the given Stream function.
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
