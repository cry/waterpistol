package main

import (
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
	name         string
	srcdir       string
	srcid        string
	ip           string
	port         int
	GOOS         string
	GOARCH       string
	modules      []string
	download_url string
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
	cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking no", "-i", HOME_DIR+project.name+"/id_c2", "-t", "ec2-user@"+project.ip, "./c2 ./cert.pem ./key.pem")
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

	for _, s := range project.modules {
		if strings.Compare(s, module) == 0 {
			log.Println("Module already enabled")
			return
		}
	}

	project.modules = append(project.modules, module)
}

func (project *project) disableModule(module string) {
	old_modules := project.modules
	project.modules = []string{}

	for _, m := range old_modules {
		if strings.Compare(m, module) != 0 {
			project.modules = append(project.modules, m)
		}
	}

}
