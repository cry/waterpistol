package types

// Event type passed between C2 and Cores, probably should make this a
// seriazable struct or something
type Event interface{}

// Module interface which will require all modules to implement a
// set of basic operations that can be called by the main core + C2
type Module interface {
	ID() string
	Init() chan Event
	Shutdown()
}

// Message : Modules communicate with each other via sending a Message with the
// capability they required and a list of string args
type Message struct {
	capability string
	args       []string
}
