///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////
package interconnect

import (
	"errors"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/xx_network/comms/connect"
	"gitlab.com/xx_network/comms/messages"
	"google.golang.org/grpc"
)

// consensus Server -> cmix Server Send Function
func (s *Comms) SendGetNDF(host *connect.Host,
	message *messages.Ping) (*NDF, error) {

	// Create the Send Function
	f := func(conn *grpc.ClientConn) (*any.Any, error) {
		// Set up the context
		ctx, cancel := connect.MessagingContext()
		defer cancel()
		//Format to authenticated message type
		// Send the message

		resultMsg, err := NewInterconnectClient(conn).GetNDF(ctx, message)
		if err != nil {
			return nil, errors.New(err.Error())
		}
		return ptypes.MarshalAny(resultMsg)
	}

	// Execute the Send function
	jww.DEBUG.Printf("Sending Post Phase message: %+v", message)
	resultMsg, err := s.Send(host, f)
	if err != nil {
		return nil, err
	}

	// Marshall the result
	result := &NDF{}
	return result, ptypes.UnmarshalAny(resultMsg, result)
}
