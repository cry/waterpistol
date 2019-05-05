package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

/**
Architectures which we support
*/
const DEFAULT_NETWORK = "basic_tcp_network"
const DEFAULT_GOOS = "linux"
const DEFAULT_ARCH = "amd64"

var ARCHS = map[string][]string{"linux": []string{"amd64", "386", "arm64", "arm"}, "windows": []string{"386", "amd64"}, "darwin": []string{"386", "amd64", "arm", "arm64"}}

var VALID_ARCHS = []string{"amd64", "386", "arm64", "arm"}

// Project structure
type project struct {
	Name          string
	Srcdir        string
	Srcid         string
	Ip            string
	Port          int
	GOOS          string
	GOARCH        string
	Modules       []string
	Download_url  string
	NetworkModule string
}

// Checks that the arch and os are valid and can be paired
// Then sets it to the projects GOOS and GOARCH
func (project *project) setArch(os string, arch string) {
	found := false
	for _, a := range ARCHS[os] {
		if strings.Compare(a, arch) == 0 {
			found = true
			break
		}
	}

	if !found {
		log.Println(os, arch, "isn't a valid os")
		log.Println("Valids are: ", ARCHS)
		return
	}

	project.GOOS = os
	project.GOARCH = arch
	log.Println("OS and ARCH set to", os, arch)
}

// Login to the c2 if it is up and running
func (project *project) ssh() {
	// CMD :  ssh -t root@ip screen -dr c2
	//	cmd := exec.Command("ssh", "-i", "./id_c2", "-t", "ec2-user@"+waterpistol.ip, "screen", "-dr", "c2")
	// Instead of scren, why not just run it on load
	cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking no", "-i", HOME_DIR+project.Name+"/id_c2", "-t", "ec2-user@"+project.Ip, "./c2 ./cert.pem ./key.pem")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}

func (project *project) setNetworkModule(module string) {
	valid := false
	for _, s := range valid_network_modules() {
		if strings.Compare(s, module) == 0 {
			valid = true
			break
		}
	}

	if !valid {
		log.Println("Please select a valid network module", valid_network_modules())
		return
	}

	project.NetworkModule = module
}

func (project *project) enableModule(module string) {
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

	for _, s := range project.Modules {
		if strings.Compare(s, module) == 0 {
			log.Println("Module already enabled")
			return
		}
	}

	project.Modules = append(project.Modules, module)
}

func (project *project) disableModule(module string) {
	old_modules := project.Modules
	project.Modules = []string{}

	for _, m := range old_modules {
		if strings.Compare(m, module) != 0 {
			project.Modules = append(project.Modules, m)
		}
	}

}

func (waterpistol *waterpistol) destroy_current_project() {
	current_project := waterpistol.current_project()
	if current_project == nil {
		log.Print("No project selected\n")
		return
	}

	log.Println("Destroying project")
	if current_project.Ip != "" {
		log.Println("Destroying c2 && ec2 instance")
		waterpistol.checkCommand(false, "cmd/c2_down", HOME_DIR+current_project.Name)
	}

	os.RemoveAll(HOME_DIR + current_project.Name)

	waterpistol.remove_project(waterpistol.current)
	waterpistol.current = -1
	waterpistol.term.SetPrompt(LIGHTBLUE + "waterpistol% " + RESET)

}

// Saves project specs to disk
func (project *project) saveProject() {
	data, err := json.MarshalIndent(project, "", " ")

	if err != nil {
		log.Println("Error saving project", project.Name)
		return
	}

	err = ioutil.WriteFile(HOME_DIR+project.Name+"/data.json", data, 0644)

	if err != nil {
		log.Println("Error saving project", project.Name)
		return
	}

}

func (project *project) Print() {
	if project.Ip == "" {
		fmt.Println(RESET+"<"+RED+project.Name+RESET+">",
			"@",
			"<"+RED+"NO_IP"+RESET+">",
			":",
			"<"+GREEN+project.GOOS+RESET+"/"+GREEN+project.GOARCH+RESET+">",
			RED+project.NetworkModule+RESET,
			project.Modules, RESET)
	} else {
		fmt.Println(RESET+"<"+GREEN+project.Name+RESET+">",
			"@",
			"<"+GREEN+project.Ip+RESET+">",
			":",
			"<"+GREEN+project.GOOS+RESET+"/"+GREEN+project.GOARCH+RESET+">",
			GREEN+project.NetworkModule+RESET,
			project.Modules, " : "+BLUE+project.Download_url, RESET)
	}
}

func load_projects() []project {
	files, err := ioutil.ReadDir(HOME_DIR)
	log.Println("Loading previous projects")

	if err != nil {
		panic(err)
	}

	ret := []project{}

	for _, f := range files {
		if !f.IsDir() {
			continue
		}

		file, err := ioutil.ReadFile(HOME_DIR + f.Name() + "/data.json")
		if err != nil {
			log.Println("Failed to load", f.Name())
			log.Println(err)
			continue
		}

		proj := project{}
		err = json.Unmarshal([]byte(file), &proj)
		if err != nil {
			log.Println("Failed to load", f.Name())
			log.Println(err)
			continue
		}
		ret = append(ret, proj)
	}

	return ret
}
