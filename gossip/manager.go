package gossip

import (
	"github.com/spf13/jwalterweatherman"
	"gitlab.com/elixxir/primitives/id"
	"gitlab.com/xx_network/comms/connect"
	"sync"
	"time"
)

// Structure holding messages for a given tag, if the tag does not yet exist
// If the tag is not created in 5 minutes, the record should be deleted
type MessageRecord struct {
	Timestamp time.Time
	Messages  []*GossipMsg
}

type ManagerFlags struct {
	// How long a message record should last in the buffer
	BufferExpirationTime time.Duration

	// Frequency with which to check the buffer.
	// Should be long, since the thread takes a lock each time it checks the buffer
	MonitorThreadFrequency time.Duration
}

func DefaultManagerFlags() ManagerFlags {
	return ManagerFlags{
		BufferExpirationTime:   300 * time.Second,
		MonitorThreadFrequency: 150 * time.Second,
	}
}

// Manager for various GossipProtocols that are accessed by tag
type Manager struct {
	comms *connect.ProtoComms

	// Stored map of GossipProtocols
	protocols    map[string]*Protocol
	protocolLock sync.RWMutex // Lock for protocols object

	// Buffer messages with tags that do not have a protocol created yet
	buffer     map[string]*MessageRecord // TODO: should this be sync.Map?
	bufferLock sync.RWMutex              // Lock for buffers object

	flags ManagerFlags
}

// Passed into NewGossip to specify how Gossip messages will be handled
type Receiver func(*GossipMsg) error

// Passed into NewGossip to specify how Gossip message signatures will be verified
type SignatureVerification func(*GossipMsg) error

// Creates a new Gossip Manager struct
func NewManager(comms *connect.ProtoComms, flags ManagerFlags) *Manager {
	m := &Manager{
		comms:     comms,
		protocols: map[string]*Protocol{},
		buffer:    map[string]*MessageRecord{},
		flags:     flags,
	}
	_ = m.bufferMonitor()
	return m
}

// Creates and stores a new Protocol in the Manager
func (m *Manager) NewGossip(tag string, flags ProtocolFlags,
	receiver Receiver, verifier SignatureVerification, peers []*id.ID) {
	m.protocolLock.Lock()
	defer m.protocolLock.Unlock()

	tmp := &Protocol{
		fingerprints: map[Fingerprint]uint64{},
		comms:        m.comms,
		peers:        peers,
		flags:        flags,
		receiver:     receiver,
		verify:       verifier,
		IsDefunct:    false,
	}

	m.protocols[tag] = tmp

	m.bufferLock.Lock()
	if record, ok := m.buffer[tag]; ok {
		for _, msg := range record.Messages {
			err := tmp.receive(msg)
			if err != nil {
				jwalterweatherman.WARN.Printf("Failed to receive message: %+v", msg)
			}
		}
		delete(m.buffer, tag)
	}
	m.bufferLock.Unlock()
}

// Returns the Gossip object for the provided tag from the Manager
func (m *Manager) Get(tag string) (*Protocol, bool) {
	m.protocolLock.RLock()
	defer m.protocolLock.RUnlock()

	p, ok := m.protocols[tag]
	return p, ok
}

// Deletes a Protocol from the Manager
func (m *Manager) Delete(tag string) {
	m.protocolLock.Lock()
	defer m.protocolLock.Unlock()

	delete(m.protocols, tag)
}

// Long-running thread to delete any messages in buffer older than 5m
func (m *Manager) bufferMonitor() chan bool {
	killChan := make(chan bool, 0)
	bufferExpirationTime := m.flags.BufferExpirationTime // Time in seconds that a record in the buffer should last
	frequency := m.flags.MonitorThreadFrequency

	go func() {
		for {
			// Loop through each tag in the buffer
			m.bufferLock.Lock()
			for tag, record := range m.buffer {
				if time.Since(record.Timestamp) > bufferExpirationTime {
					delete(m.buffer, tag)
				}
			}
			m.bufferLock.Unlock()

			select {
			case <-killChan:
				return
			default:
				time.Sleep(frequency)
			}
		}
	}()

	return killChan
}
