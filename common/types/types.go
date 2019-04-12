package types

import "malware/common/messages"

// Module interface which will require all modules to implement a
// set of basic operations that can be called by the main core
type Module interface {
	ID() string
	Init()
	HandleMessage(*messages.CheckCmdReply, func(*messages.ImplantReply))
	Shutdown()
}

type Error int

const (
	ERR_FILE_NOT_FOUND Error = 0
)
