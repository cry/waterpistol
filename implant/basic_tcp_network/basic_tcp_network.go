package basic_tcp_network

import (
	"fmt"
	"malware/common"
	"malware/common/types"
	"net"
	"strings"

	"google.golang.org/grpc"
)

// Message to be sent by
type sendText struct {
	text string
}

// SendText generates a `sendCommand` event message
func SendText(args []string) types.Event {
	return sendText{strings.Join(args, "\n")}
}

type state struct {
	running      bool
	eventChannel chan types.Event
	listener     net.Listener
	grpc         *grpc.Server
}

type settings struct {
	state *state
	port  int16
}

func (settings settings) createServer() {
	host := fmt.Sprintf(":%d", settings.port)
	listener, err := net.Listen("tcp", host)

	if err != nil {
		common.Panicf(err, "Hosting on %s failed", host)
	}

	server := grpc.NewServer()
	settings.state.listener = listener
	settings.state.grpc = server

	pb.Register

	fmt.Println("Listening on ", server.GetServiceInfo())
}

func (settings settings) listenServer() {
	for settings.state.running {
		conn, err := settings.state.listener.Accept()
		if err != nil {
			common.Panicf(err, "Received message can't read %V", conn)
		}
		// TODO: deserialize message
		buffer := make([]byte, 1024)
		for n, err := conn.Read(buffer); err == nil; {
			fmt.Println(n, ":", string(buffer))
		}
		conn.Close()
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
	port := int16(8080)
	state := state{eventChannel: nil}

	return settings{&state, port}
}

func (settings settings) Init() chan types.Event {
	settings.state.eventChannel = make(chan types.Event, 1)
	settings.state.running = true
	settings.createServer()

	go settings.listenServer()
	go settings.runEvents()

	return settings.state.eventChannel
}

func (settings settings) Shutdown() {
	settings.state.running = false
	settings.state.listener.Close()
}

func (settings) ID() string { return "Basic TCP Network" }
