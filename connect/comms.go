////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

// Handles the basic top-level comms object used across all packages

package connect

import (
	"crypto/tls"
	"crypto/x509"
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
	networkId *id.ID

	// Private key of the local comms instance
	privateKey *rsa.PrivateKey

	// Disables the checking of authentication signatures for testing setups
	disableAuth bool

	listeningAddress string

	// SERVER-ONLY FIELDS ------------------------------------------------------

	// A map of reverse-authentication tokens
	tokens *token.Map

	// Low-level net.Listener object that listens at listeningAddress
	netListener net.Listener

	// Local network server
	grpcServer *grpc.Server
	grpcCreds  credentials.TransportCredentials

	// CLIENT-ONLY FIELDS ------------------------------------------------------

	// Used to store the public key used for generating Client Id
	pubKeyPem []byte

	// Used to store the salt used for generating Client Id
	salt []byte

	// cmux interface for starting http listener when serving with web
	mux cmux.CMux

	// https certificate infra for replacing tls cert with no downtime
	httpsCertificate *tls.Certificate

	// -------------------------------------------------------------------------
}

// GetId returns a copy of the ProtoComms networkId
func (c *ProtoComms) GetId() *id.ID {
	return c.networkId.DeepCopy()
}

// GetServer returns the ProtoComms grpc.Server object
func (c *ProtoComms) GetServer() *grpc.Server {
	return c.grpcServer
}

// CreateCommClient creates a ProtoComms client-type object to be
// used in various initializers.
func CreateCommClient(id *id.ID, pubKeyPem, privKeyPem,
	salt []byte) (*ProtoComms, error) {
	// Build the ProtoComms object
	pc := &ProtoComms{
		networkId: id,
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

// StartCommServer creates a ProtoComms server-type object to be used in various initializers.
// Opens a net.Listener the local address specified by listeningAddr.
func StartCommServer(id *id.ID, listeningAddr string,
	certPEMblock, keyPEMblock []byte, preloadedHosts []*Host) (*ProtoComms, error) {

listen:
	// Listen on the given address
	lis, err := net.Listen("tcp", listeningAddr)
	if err != nil {
		if strings.Contains(err.Error(), "bind: address already in use") {
			jww.WARN.Printf("Could not listen on %s, is port in use? waiting 30s: %s", listeningAddr, err.Error())
			time.Sleep(30 * time.Second)
			goto listen
		}
		return nil, errors.New(err.Error())
	}

	// Build the comms object
	pc := &ProtoComms{
		networkId:        id,
		netListener:      lis,
		tokens:           token.NewMap(),
		Manager:          newManager(),
		listeningAddress: listeningAddr,
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

		// Set the public key data
		pc.pubKeyPem = certPEMblock

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

func (c *ProtoComms) Restart() error {
	if TestingOnlyDisableTLS {
		c.grpcServer = grpc.NewServer(
			grpc.MaxConcurrentStreams(MaxConcurrentStreams),
			grpc.MaxRecvMsgSize(math.MaxInt32),
			grpc.KeepaliveParams(KaOpts),
			grpc.KeepaliveEnforcementPolicy(KaEnforcement))
	} else {
		c.grpcServer = grpc.NewServer(grpc.Creds(c.grpcCreds),
			grpc.MaxConcurrentStreams(MaxConcurrentStreams),
			grpc.MaxRecvMsgSize(math.MaxInt32),
			grpc.KeepaliveParams(KaOpts),
			grpc.KeepaliveEnforcementPolicy(KaEnforcement))
	}

	if c.netListener != nil {
		return errors.New("ProtoComms is already listening")
	}
listen:
	// Listen on the given address
	lis, err := net.Listen("tcp", c.listeningAddress)
	if err != nil {
		if strings.Contains(err.Error(), "bind: address already in use") {
			jww.WARN.Printf("Could not listen on %s, is port in use? waiting 30s: %s", c.listeningAddress, err.Error())
			time.Sleep(30 * time.Second)
			goto listen
		}
		return errors.New(err.Error())
	}

	c.netListener = lis
	return nil
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
	go listenGRPC(c.netListener)
}

// ServeWithWeb is a non-blocking call that begins serving content
// for grpcWeb (over HTTP) and GRPC on the same port.
// GRPC endpoints must be registered before making this call.
func (c *ProtoComms) ServeWithWeb() {
	grpcServer := c.GetServer()

	// Split netListener into two distinct listeners for GRPC and HTTP
	mux := cmux.New(c.netListener)
	grpcMatcher := cmux.HTTP2HeaderField("content-type", "application/grpc")

	// TODO: do we need this anymore?  the header should match either case...
	if TestingOnlyDisableTLS {
		grpcMatcher = cmux.HTTP2()
	}
	grpcL := mux.Match(grpcMatcher)
	c.mux = mux

	listenGRPC := func(l net.Listener) {
		// This blocks for the lifetime of the listener.
		if err := grpcServer.Serve(l); err != nil {
			// Cannot panic here due to shared net.Listener
			jww.ERROR.Printf("Failed to serve GRPC: %+v", err)
		}
		jww.INFO.Printf("Shutting down GRPC server listener")
	}
	listenPort := func() {
		if err := mux.Serve(); err != nil {
			// Cannot panic here due to shared net.Listener
			jww.ERROR.Printf("Failed to serve port: %+v", err)
		}
		jww.INFO.Printf("Shutting down port server listener")
	}
	go listenGRPC(grpcL)
	go listenPort()
}

// ProvisionHttps provides a tls cert and key to the thread which serves the
// grpcweb endpoints, allowing it to serve with https.  Note that https will
// not be usable until this has been called at least once, unblocking the
// listenHTTP func in ServeWithWeb.  Future calls will be handled by the
// startUpdateCertificate thread.
func (c *ProtoComms) ServeHttps(cert, key []byte) error {
	if c.mux == nil {
		return errors.New("mux does not exist; is https enabled?")
	}

	if c.netListener == nil {
		return errors.New("ProtoComms is closed, call ServeWithWeb to initialize listener")
	}

	var httpL net.Listener
	if TestingOnlyDisableTLS && cert == nil && key == nil {
		httpL = c.mux.Match(cmux.HTTP1())
	} else {
		httpL = c.mux.Match(cmux.TLS())
	}

	grpcServer := c.grpcServer
	var keyPair tls.Certificate
	var err error
	if cert == nil && key == nil {
		if !TestingOnlyDisableTLS {
			return errors.New("cannot serve without tls outside of testing")
		}
	} else {
		keyPair, err = tls.X509KeyPair(
			cert, key)
		if err != nil {
			return errors.WithMessage(err, "cert & key could not be parsed to valid tls certificate")
		}
	}

	listenHTTP := func(l net.Listener) {
		jww.INFO.Printf("Starting HTTP listener on GRPC endpoints: %+v",
			grpcweb.ListGRPCResources(grpcServer))
		httpServer := grpcweb.WrapServer(grpcServer,
			grpcweb.WithOriginFunc(func(origin string) bool { return true }))
		if TestingOnlyDisableTLS && c.privateKey == nil {
			jww.WARN.Printf("Starting HTTP server to without TLS!")
			if err := http.Serve(l, httpServer); err != nil {
				// Cannot panic here due to shared net.Listener
				jww.ERROR.Printf("Failed to serve HTTP: %+v", err)
			}
		} else {
			// Configure TLS for this listener, using the config from
			// http.ServeTLS
			tlsConf := &tls.Config{}
			tlsConf.NextProtos = append(tlsConf.NextProtos, "h2", "http/1.1")

			// We use the GetCertificate field and a function which returns
			// c.httpsCertificate wrapped by a ReadLock to allow for future
			// changes to the certificate without downtime
			c.httpsCertificate = &keyPair
			tlsConf.GetCertificate = c.GetCertificateFunc()

			var serverName string
			cert, err := x509.ParseCertificate(keyPair.Certificate[0])
			if err != nil {
				jww.FATAL.Panicf("Failed to load TLS certificate: %+v", err)
			}
			serverName = cert.DNSNames[0]
			tlsConf.ServerName = serverName

			tlsLis := tls.NewListener(l, tlsConf)
			if err := http.Serve(tlsLis, httpServer); err != nil {
				// Cannot panic here due to shared net.Listener
				jww.ERROR.Printf("Failed to serve HTTP: %+v", err)
			}
		}

		jww.INFO.Printf("Shutting down HTTP server listener")
	}

	go listenHTTP(httpL)

	return nil
}

// Shutdown performs a graceful shutdown of the local server.
func (c *ProtoComms) Shutdown() {
	// Also handles closing of net.Listener
	c.grpcServer.GracefulStop()
	// Close all Manager connections
	c.DisconnectAll()
	c.grpcServer = nil
	c.netListener = nil
	c.mux = nil
	jww.INFO.Printf("Comms server successfully shut down")
}

// Stringer method
func (c *ProtoComms) String() string {
	if c.netListener == nil {
		return c.listeningAddress
	}
	return c.netListener.Addr().String()
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
func (c *ProtoComms) Send(host *Host, f func(conn Connection) (*any.Any,
	error)) (result *any.Any, err error) {

	jww.TRACE.Printf("Attempting to send to host: %s", host)
	fSh := func(conn Connection) (interface{}, error) {
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
func (c *ProtoComms) Stream(host *Host, f func(conn Connection) (
	interface{}, error)) (client interface{}, err error) {

	// Ensure the connection is running
	jww.TRACE.Printf("Attempting to stream to host: %s", host)
	return c.transmit(host, f)
}

// GetCertificateFunc returns a function which serves as the value of
// GetCertificate in the https TLS config, enabling certificate changing
// without downtime
func (c *ProtoComms) GetCertificateFunc() func(*tls.ClientHelloInfo) (*tls.Certificate, error) {
	return func(clientHello *tls.ClientHelloInfo) (*tls.Certificate, error) {
		return c.httpsCertificate, nil
	}
}

// returns true if the connection error is one of the connection errors which
// should be retried
func isConnError(err error) bool {
	return strings.Contains(err.Error(), "context deadline exceeded") ||
		strings.Contains(err.Error(), "connection refused") ||
		strings.Contains(err.Error(), "host disconnected")
}
