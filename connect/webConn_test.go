package connect

import (
	"gitlab.com/xx_network/crypto/csprng"
	"gitlab.com/xx_network/primitives/id"
	"testing"
)

func TestWebConn_Close(t *testing.T) {
	rng := csprng.NewSystemRNG()
	hostId, err := id.NewRandomID(rng, id.User)
	if err != nil {
		t.Fatal(err)
	}
	hostParams := GetDefaultHostParams()
	hostParams.ConnectionType = Web
	h, err := NewHost(hostId, "0.0.0.0", nil, hostParams)
	if err != nil {
		t.Fatal(err)
	}
	wc := webConn{
		h:          h,
		connection: nil,
	}
	err = wc.Connect()
	if err != nil {
		t.Fatal(err)
	}

	err = wc.Close()
	if err != nil {
		t.Fatal(err)
	}
}

func TestWebConn_Connect(t *testing.T) {
	rng := csprng.NewSystemRNG()
	hostId, err := id.NewRandomID(rng, id.User)
	if err != nil {
		t.Fatal(err)
	}
	hostParams := GetDefaultHostParams()
	hostParams.ConnectionType = Web
	h, err := NewHost(hostId, "0.0.0.0", nil, hostParams)
	if err != nil {
		t.Fatal(err)
	}
	wc := webConn{
		h:          h,
		connection: nil,
	}
	err = wc.Connect()
	if err != nil {
		t.Fatal(err)
	}
}

func TestWebConn_GetGrpcConn(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Did not receive expected fatal error")
		}
	}()
	rng := csprng.NewSystemRNG()
	hostId, err := id.NewRandomID(rng, id.User)
	if err != nil {
		t.Fatal(err)
	}
	hostParams := GetDefaultHostParams()
	hostParams.ConnectionType = Web
	h, err := NewHost(hostId, "0.0.0.0", nil, hostParams)
	if err != nil {
		t.Fatal(err)
	}
	wc := webConn{
		h:          h,
		connection: nil,
	}
	err = wc.Connect()
	if err != nil {
		t.Fatal(err)
	}

	conn := wc.GetGrpcConn()
	if conn != nil {
		t.Fatal("Expected panic, received conn instead")
	}

}

func TestWebConn_GetWebConn(t *testing.T) {
	rng := csprng.NewSystemRNG()
	hostId, err := id.NewRandomID(rng, id.User)
	if err != nil {
		t.Fatal(err)
	}
	hostParams := GetDefaultHostParams()
	hostParams.ConnectionType = Web
	h, err := NewHost(hostId, "0.0.0.0", nil, hostParams)
	if err != nil {
		t.Fatal(err)
	}
	wc := webConn{
		h:          h,
		connection: nil,
	}
	err = wc.Connect()
	if err != nil {
		t.Fatal(err)
	}

	conn := wc.GetWebConn()
	if conn == nil {
		t.Fatal("Expected grpcConn, received nil instead")
	}
}

func TestWebConn_IsWeb(t *testing.T) {
	wc := webConn{}
	if !wc.IsWeb() {
		t.Fatal("WebConn is not web")
	}
}
