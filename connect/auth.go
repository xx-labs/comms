///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

// Handles authentication logic for the top-level comms object

package connect

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/elixxir/crypto/signature/rsa"
	"gitlab.com/elixxir/crypto/xx"
	"gitlab.com/elixxir/primitives/id"
	"gitlab.com/xx_network/comms/connect/token"
	pb "gitlab.com/xx_network/comms/messages"
	"google.golang.org/grpc/metadata"
)

// Auth represents an authorization state for a message or host
type Auth struct {
	// Indicates whether authentication was successful
	IsAuthenticated bool
	// The information about the Host that sent the authenticated communication
	Sender *Host
}

// Perform the client handshake to establish reverse-authentication
func (c *ProtoComms) clientHandshake(host *Host) (err error) {

	// Set up the context
	client := pb.NewGenericClient(host.connection)
	ctx, cancel := MessagingContext()
	defer cancel()

	// Send the token request message
	result, err := client.RequestToken(ctx,
		&pb.Ping{})
	if err != nil {
		return errors.New(err.Error())
	}

	// Pack the authenticated message with signature enabled
	msg, err := c.PackAuthenticatedMessage(&pb.AssignToken{
		Token: result.Token,
	}, host, true)
	if err != nil {
		return errors.New(err.Error())
	}

	// Add special client-specific info to the
	// newly generated authenticated message if needed
	if c.pubKeyPem != nil && c.salt != nil {
		msg.Client.Salt = c.salt
		msg.Client.PublicKey = string(c.pubKeyPem)
	}

	// Set up the context
	ctx, cancel = MessagingContext()
	defer cancel()

	// Send the authenticate token message
	_, err = client.AuthenticateToken(ctx, msg)
	if err != nil {
		return errors.New(err.Error())
	}

	// Assign the host token
	host.transmissionToken.Set(result.Token)

	return
}

// Convert any message type into a authenticated message
func (c *ProtoComms) PackAuthenticatedMessage(msg proto.Message, host *Host,
	enableSignature bool) (*pb.AuthenticatedMessage, error) {

	// Marshall the provided message into an Any type
	anyMsg, err := ptypes.MarshalAny(msg)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	// Build the authenticated message
	authMsg := &pb.AuthenticatedMessage{
		ID:      c.Id.Marshal(),
		Token:   host.transmissionToken.Get(),
		Message: anyMsg,
		Client: &pb.ClientID{
			Salt:      make([]byte, 0),
			PublicKey: "",
		},
	}

	// If signature is enabled, sign the message and add to payload
	if enableSignature && !c.disableAuth {
		authMsg.Signature, err = c.SignMessage(msg)
		if err != nil {
			return nil, err
		}
	}

	return authMsg, nil
}

// Add authentication fields to a given context and return it
func (c *ProtoComms) PackAuthenticatedContext(host *Host,
	ctx context.Context) context.Context {

	ctx = metadata.AppendToOutgoingContext(ctx, "ID", c.Id.String())
	ctx = metadata.AppendToOutgoingContext(ctx, "TOKEN", base64.StdEncoding.EncodeToString(host.transmissionToken.Get()))
	return ctx
}

// Returns authentication packed into a context
func UnpackAuthenticatedContext(ctx context.Context) (*pb.AuthenticatedMessage, error) {
	auth := &pb.AuthenticatedMessage{}
	var err error

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("unable to retrieve meta data / header")
	}

	idStr := md.Get("ID")[0]
	auth.ID, err = base64.StdEncoding.DecodeString(idStr)
	if err != nil {
		return nil, errors.WithMessage(err, "could not decode authentication ID")
	}

	tokenStr := md.Get("TOKEN")[0]
	auth.Token, err = base64.StdEncoding.DecodeString(tokenStr)
	if err != nil {
		return nil, errors.WithMessage(err, "could not decode authentication Live")
	}

	return auth, nil
}

// Generates a new token and adds it to internal state
func (c *ProtoComms) GenerateToken() ([]byte, error) {
	return c.tokens.Generate().Marshal(), nil
}

// Performs the dynamic authentication process such that Hosts that were not
// already added to the Manager can establish authentication
func (c *ProtoComms) dynamicAuth(msg *pb.AuthenticatedMessage) (
	host *Host, err error) {

	// Verify the client is attempting a dynamic authentication
	if msg.Client == nil || msg.Client.PublicKey == "" || msg.Client.Salt == nil {
		return nil, errors.New("Invalid dynamic authentication attempt!")
	}

	// Process the public key
	pubKey, err := rsa.LoadPublicKeyFromPem([]byte(msg.Client.PublicKey))
	if err != nil {
		return nil, errors.New(err.Error())
	}

	// Generate the user's ID from supplied Client information
	uid, err := xx.NewID(pubKey, msg.Client.Salt, id.User)
	if err != nil {
		return nil, err
	}

	// Verify the ID provided correctly matches the generated ID
	if !bytes.Equal(msg.ID, uid.Marshal()) {
		return nil, errors.Errorf(
			"Provided ID does not match. Expected: %s, Actual: %s",
			uid.String(), msg.ID)
	}

	// Create and add the new host to the manager
	host, err = newDynamicHost(uid, []byte(msg.Client.PublicKey))
	if err != nil {
		return
	}
	c.addHost(host)
	return
}

// Validates a signed token using internal state
func (c *ProtoComms) ValidateToken(msg *pb.AuthenticatedMessage) (err error) {
	// Convert EntityID to ID
	msgID, err := id.Unmarshal(msg.ID)
	if err != nil {
		return err
	}

	// Verify the Host exists for the provided ID
	host, ok := c.GetHost(msgID)
	if !ok {

		// If the host does not already exist, attempt dynamic authentication
		jww.DEBUG.Printf("Attempting dynamic authentication: %s", msgID.String())
		host, err = c.dynamicAuth(msg)
		if err != nil {
			return errors.Errorf(
				"Unable to complete dynamic authentication: %+v", err)
		}
	}

	// Get the signed token
	tokenMsg := &pb.AssignToken{}
	err = ptypes.UnmarshalAny(msg.Message, tokenMsg)
	if err != nil {
		return errors.Errorf("Unable to unmarshal token: %+v", err)
	}

	remoteToken, err := token.Unmarshal(tokenMsg.Token)
	if err != nil {
		return errors.Errorf("Unable to unmarshal token: %+v", err)
	}

	// Verify the token signature unless disableAuth has been set for testing
	if !c.disableAuth {
		if err := c.VerifyMessage(tokenMsg, msg.Signature, host); err != nil {
			return errors.Errorf("Invalid token signature: %+v", err)
		}
	}

	ok = c.tokens.Validate(remoteToken)
	if !ok {
		jww.ERROR.Printf("Failed to validate token %v from %s", remoteToken, host)
		return errors.Errorf("Failed to validate token: %v", remoteToken)
	}
	// Token has been validated and can be safely stored
	host.receptionToken.Set(tokenMsg.Token)
	jww.DEBUG.Printf("Live validated: %v", tokenMsg.Token)
	return
}

// AuthenticatedReceiver handles reception of an AuthenticatedMessage,
// checking if the host is authenticated & returning an Auth state
func (c *ProtoComms) AuthenticatedReceiver(msg *pb.AuthenticatedMessage) (*Auth, error) {
	// Convert EntityID to ID
	msgID, err := id.Unmarshal(msg.ID)
	if err != nil {
		return nil, err
	}

	// Try to obtain the Host for the specified ID
	host, ok := c.GetHost(msgID)
	if !ok {
		return &Auth{
			IsAuthenticated: false,
			Sender:          &Host{},
		}, nil
	}

	// Check the token's validity
	receptionToken := host.receptionToken.Get()
	validToken := receptionToken != nil && msg.Token != nil &&
		bytes.Compare(receptionToken, msg.Token) == 0

	// Assemble the Auth object
	res := &Auth{
		IsAuthenticated: validToken,
		Sender:          host,
	}

	jww.TRACE.Printf("Authentication status: %v, ProvidedId: %v ProvidedToken: %v",
		res.IsAuthenticated, msg.ID, msg.Token)
	return res, nil
}

// DisableAuth makes the authentication code skip signing and signature verification if the
// set.  Can only be set while in a testing structure.  Is not thread safe.
func (c *ProtoComms) DisableAuth() {
	jww.WARN.Print("Auth checking disabled, running insecurely")
	c.disableAuth = true
}

// Takes a message and returns its signature
// The message is signed with the ProtoComms RSA PrivateKey
func (c *ProtoComms) SignMessage(msg proto.Message) ([]byte, error) {
	// Hash the message data
	msgBytes, err := proto.Marshal(msg)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	options := rsa.NewDefaultOptions()
	hash := options.Hash.New()
	hash.Write(msgBytes)
	hashed := hash.Sum(nil)

	// Obtain the private key
	key := c.GetPrivateKey()
	if key == nil {
		return nil, errors.Errorf("Cannot sign message: No private key")
	}
	// Sign the message and return the signature
	signature, err := rsa.Sign(rand.Reader, key, options.Hash, hashed, nil)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	return signature, nil
}

// Takes a message and a Host, verifies the signature
// using Host public key, returning an error if invalid
func (c *ProtoComms) VerifyMessage(msg proto.Message, signature []byte, host *Host) error {

	// Get hashed data of the message
	msgBytes, err := proto.Marshal(msg)
	if err != nil {
		return errors.New(err.Error())
	}
	options := rsa.NewDefaultOptions()
	hash := options.Hash.New()
	hash.Write(msgBytes)
	hashed := hash.Sum(nil)

	// Verify signature of message using host public key
	err = rsa.Verify(host.rsaPublicKey, options.Hash, hashed, signature, nil)
	if err != nil {
		return errors.New(err.Error())
	}

	return nil
}
