package routes

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os/user"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

func runCommands() {
	hostKeyCallback, err := knownhosts.New("/Users/benjamindelacruzmartinez/.ssh/known_hosts")
	if err != nil {
		log.Fatal(err)
	}
	usr, _ := user.Current()
	file := usr.HomeDir + "/.ssh/deployment"
	key, errFile := ioutil.ReadFile(file)

	if errFile != nil {
		log.Fatalf("unable to read private key: %v", errFile)
	}

	signer, errSecond := ssh.ParsePrivateKey(key)
	if errSecond != nil {
		log.Fatalf("unable to parse private key: %v", errSecond)
	}

	if err != nil {
		log.Fatal(err)
	}
	config := &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: hostKeyCallback,
		Timeout:         0,
	}

	client, err := ssh.Dial("tcp", "kimvex.com:22", config)

	if err != nil {
		panic("Failed to dial: " + err.Error())
	}
	session, err := client.NewSession()
	if err != nil {
		panic("Failed to create session: " + err.Error())
	}
	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b
	errRun := session.Run("cd /home/kimvex-pages && sudo git status")
	if errRun != nil {
		fmt.Println("Failed to run: " + errRun.Error())
	}
	fmt.Println(b.String())
}
