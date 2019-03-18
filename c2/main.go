package main

import (
	"context"
	"fmt"
	"malware/common"
	"net"

	pb "malware/common/messages"

	"google.golang.org/grpc"
)

type c2 struct {
	port int
}

func (c2 *c2) Exec(context.Context, *pb.ExecRequest) (*pb.ExecReply, error) {
	fmt.Println("Hi")
	return &pb.ExecReply{Reply: "hi"}, nil
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
	pb.RegisterImplantServer(server, c2)
	server.Serve(listener)
	fmt.Println("Listening on s", server.GetServiceInfo())

}
