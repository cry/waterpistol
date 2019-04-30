package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

/**
Architectures which we support
*/
const DEFAULT_GOOS = "linux"
const DEFAULT_ARCH = "amd64"

var ARCHS = map[string][]string{"linux": []string{"amd64", "386", "arm64", "arm"}, "windows": []string{"386", "amd64"}, "darwin": []string{"386", "amd64", "arm", "arm64"}}

var VALID_ARCHS = []string{"amd64", "386", "arm64", "arm"}

// Project structure
type project struct {
	Name         string
	Srcdir       string
	Srcid        string
	Ip           string
	Port         int
	GOOS         string
	GOARCH       string
	Modules      []string
	Download_url string
}

// Checks that the arch and os are valid and can be paired
// Then sets it to the projects GOOS and GOARCH
func (project *project) set_arch(os string, arch string) {
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

func load_projects() []project {
	files, err := ioutil.ReadDir(HOME_DIR)
	log.Println("Loading previous projects")

	checkError(err)

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
