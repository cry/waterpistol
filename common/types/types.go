package types

// Event type passed between Cores
type Message struct {
	Capability string
	Caller     chan Message // Channel to reply to caller on
	Args       []string
}

type Rx_Tx struct {
	Rx chan Message
	Tx chan Message
}

// Module interface which will require all modules to implement a
// set of basic operations that can be called by the main core
type Module interface {
	ID() string
	Init() Rx_Tx
	Shutdown()
	Capability() string
}
