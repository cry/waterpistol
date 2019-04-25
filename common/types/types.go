package types

import "malware/common/messages"

// Module interface which will require all modules to implement a
// set of basic operations that can be called by the main core
type Module interface {
	ID() string
	Init()
	HandleMessage(*messages.CheckCmdReply, func(*messages.ImplantReply)) bool
	Shutdown()
}

const (
	ERR_FILE_NOT_FOUND       int32 = 1
	ERR_MODULE_NOT_IMPL      int32 = 2
	ERR_PORTSCAN_RUNNING     int32 = 3
	ERR_IPSCAN_RUNNING       int32 = 4
	ERR_INVALID_RANGE_IPSCAN int32 = 5
)

var ErrorToString = map[int32]string{
	ERR_FILE_NOT_FOUND:       "File not found on implant",
	ERR_MODULE_NOT_IMPL:      "Module is not implemented on implant",
	ERR_PORTSCAN_RUNNING:     "A portscan is already running on this implant",
	ERR_IPSCAN_RUNNING:       "An IPScan is already running on this implant",
	ERR_INVALID_RANGE_IPSCAN: "Invalid IPv4 range",
}
