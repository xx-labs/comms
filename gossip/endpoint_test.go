package gossip

import (
	"context"
	"gitlab.com/elixxir/primitives/id"
	"gitlab.com/xx_network/comms/connect"
	"testing"
	"time"
)

// Test endpoint when manager has a protocol
func TestManager_Endpoint_toProtocol(t *testing.T) {
	m := NewManager(&connect.ProtoComms{}, DefaultManagerFlags())

	var received bool
	r := func(msg *GossipMsg) error {
		received = true
		return nil
	}
	v := func(*GossipMsg) error {
		return nil
	}
	m.NewGossip("test", DefaultProtocolFlags(), r, v,
		[]*id.ID{id.NewIdFromString("test", id.Node, t)})

	_, err := m.Endpoint(context.Background(), &GossipMsg{
		Tag:       "test",
		Origin:    []byte("origin"),
		Payload:   []byte("payload"),
		Signature: []byte("signature"),
	})
	if err != nil {
		t.Errorf("Failed to send: %+v", err)
	}

	if !received {
		t.Errorf("Didn't receive message in protocol")
	}
}

// Test endpoint function when there is no protocol and no buffer record
func TestManager_Endpoint_toNewBuffer(t *testing.T) {
	m := NewManager(&connect.ProtoComms{}, DefaultManagerFlags())
	_, err := m.Endpoint(context.Background(), &GossipMsg{
		Tag:       "test",
		Origin:    []byte("origin"),
		Payload:   []byte("payload"),
		Signature: []byte("signature"),
	})
	if err != nil {
		t.Errorf("Failed to send message: %+v", err)
	}
	r, ok := m.buffer["test"]
	if !ok {
		t.Error("Did not create expected message record")
	}
	if len(r.Messages) != 1 {
		t.Errorf("Did not add message to buffer")
	}
}

// Test endpoint function when there is no protocol, but an existing buffer
func TestComms_Endpoint_toExistingBuffer(t *testing.T) {
	m := NewManager(&connect.ProtoComms{}, DefaultManagerFlags())
	now := time.Now()
	m.buffer["test"] = &MessageRecord{
		Timestamp: now,
		Messages:  []*GossipMsg{{Tag: "test"}},
	}
	_, err := m.Endpoint(context.Background(), &GossipMsg{
		Tag:       "test",
		Origin:    []byte("origin"),
		Payload:   []byte("payload"),
		Signature: []byte("signature"),
	})
	if err != nil {
		t.Errorf("Failed to send message: %+v", err)
	}
	r, ok := m.buffer["test"]
	if !ok {
		t.Error("Did not create expected message record")
	}
	if len(r.Messages) != 2 {
		t.Errorf("Did not add message to buffer")
	}
}

func TestComms_Stream(t *testing.T) {
	// TODO: Implement test once streaming is enabled
}
