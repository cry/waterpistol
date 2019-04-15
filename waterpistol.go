package main

import (
	"fmt"
	"github.com/chzyer/readline"
	"io"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"
)

// Contains current options/c2ip/modules
type waterpistol struct {
	srcdir  string
	srcid   string
	ip      string
	port    int
	modules []string
	term    *readline.Instance
}

var valid_modules = []string{"sh", "portscan", "file_extractor", "file_uploader"}

func (waterpistol *waterpistol) writeString(str string) {
	log.Print(str)
}

func checkProgram() {
	// Ensure GO is installed and GOPATH is set
	cmd := exec.Command("go", "version")
	if cmd.Run() != nil {
		fmt.Println("Please make sure Go compiler is installed")
		os.Exit(1)
	}
	if len(os.Getenv("GOPATH")) == 0 {
		fmt.Println("Please make sure ENV GOPATH is set")
		os.Exit(1)
	}

	// Check we are in the go directory
	if ex, err := os.Executable(); err != nil && strings.Contains(path.Dir(ex), os.Getenv("GOPATH")) {
		fmt.Println("Please run this binary in the root directory of the waterpistol")
		os.Exit(1)
	}

	// Check that implant/common/c2 exist
	cmd = exec.Command("ls")
	out, _ := cmd.CombinedOutput()
	outp := string(out)
	if !strings.Contains(outp, "common") ||
		!strings.Contains(outp, "c2") ||
		!strings.Contains(outp, "implant") {
		fmt.Println("Please run this binary in the root directory of the waterpistol")
		os.Exit(1)
	}
}

func (waterpistol *waterpistol) ssh() {
	log.Println("Running c2")

	// CMD :  ssh -t root@ip screen -dr c2
	cmd := exec.Command("ssh", "-t", "root@"+waterpistol.ip, "screen", "-dr", "c2")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		panic(err)
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

func valid_check(string) []string {
	return valid_modules
}

var completer = readline.NewPrefixCompleter(
	readline.PcItem("compile"),
	readline.PcItem("enable", readline.PcItemDynamic(valid_check)),
	readline.PcItem("disable"),
	readline.PcItem("exit"),
	readline.PcItem("ssh"),
)

func (waterpistol *waterpistol) setup_terminal() {
	log.SetFlags(log.Ltime)

	l, err := readline.NewEx(&readline.Config{
		Prompt:          "\033[96mwaterpistol%\033[0m ",
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
	log.SetOutput(l.Stderr())
	waterpistol.term = l
}

func (waterpistol *waterpistol) enable(module string) {
	valid := false
	for _, s := range valid_modules {
		if strings.Compare(s, module) == 0 {
			valid = true
			break
		}
	}

	if !valid {
		fmt.Println("Please select a valid module", valid_modules)
		return
	}

	waterpistol.modules = append(waterpistol.modules, module)
}

func (waterpistol *waterpistol) disable(module string) {

}
func help() {
	log.Println("Commands:")
	log.Println("\tcompile          -> ")
	log.Println("\tssh              -> ")
	log.Println("\tenable <module>  -> ")
	log.Println("\tdisable <module> -> ")
	log.Println("\thelp             -> ")
}

func (waterpistol *waterpistol) handle(line string) {
	parts := strings.Split(strings.TrimSpace(line), " ")
	switch parts[0] {
	case "compile":
		if len(waterpistol.modules) == 0 {
			waterpistol.writeString("Maybe `enable` a few modules\n")
			return
		}
		waterpistol.compile_c2_implant()
	case "ssh":
		if len(waterpistol.ip) == 0 {
			waterpistol.writeString("Maybe `compile` first\n")
			return
		}
		waterpistol.ssh()
	case "enable":
		if len(parts) != 2 {
			help()
			return
		}
		waterpistol.enable(parts[1])
	case "disable":
		if len(parts) != 2 {
			help()
			return
		}
		waterpistol.disable(parts[1])
	case "help":
		help()
	default:
		waterpistol.writeString("What")
	}
}

func (waterpistol *waterpistol) ReadUserInput() {
	for {
		line, err := waterpistol.term.Readline()
		if err == readline.ErrInterrupt {
			continue
		} else if err == io.EOF {
			break
		}

		waterpistol.handle(line)
	}
}

func main() {
	checkProgram()

	rand.Seed(time.Now().UTC().UnixNano())

	waterpistol := &waterpistol{}

	waterpistol.setup_terminal()
	defer waterpistol.term.Close()

	waterpistol.ReadUserInput()
}
