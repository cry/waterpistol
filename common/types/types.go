package types

import "malware/common/messages"

// Module interface which will require all modules to implement a
// set of basic operations that can be called by the main core
type Module interface {
	ID() string
	HandleMessage(message *messages.CheckCmdReply, reply_function func(*messages.CheckCmdRequest)) bool
	Shutdown()
}

/**
HandleMessage is called from basic_tcp_network every time a message is received for all cores
*/

const (
	ERR_FILE_NOT_FOUND int32 = iota + 1
	ERR_MODULE_NOT_IMPL
	ERR_PORTSCAN_RUNNING
	ERR_IPSCAN_RUNNING
	ERR_INVALID_RANGE_IPSCAN
	ERR_CMD_TIMEOUT
)

var ErrorToString = map[int32]string{
	ERR_FILE_NOT_FOUND:       "File not found on implant",
	ERR_MODULE_NOT_IMPL:      "Module is not implemented on implant",
	ERR_PORTSCAN_RUNNING:     "A portscan is already running on this implant",
	ERR_IPSCAN_RUNNING:       "An IPScan is already running on this implant",
	ERR_INVALID_RANGE_IPSCAN: "Invalid IPv4 range",
	ERR_CMD_TIMEOUT:          "Command timed out",
}
