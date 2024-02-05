package main

import (
	"flag"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
)

func main() {
	// Define command-line flags for server IP, private key file, and command
	serverIP := flag.String("server", "", "Server IP address")
	username := flag.String("username","","User Name")
	privateKeyFile := flag.String("key", "", "Private key file")
	command := flag.String("command", "", "Command to execute on the server")
	flag.Parse()

	// Check if server IP, private key file, and command are provided
	if *serverIP == "" || *privateKeyFile == "" || *command == "" || *username == "" {
		log.Fatal("Server IP, private key file, and command must be provided")
	}

	// Read private key file
	key, err := ioutil.ReadFile(*privateKeyFile)
	if err != nil {
		log.Fatalf("Failed to read private key file: %s", err)
	}

	// Parse private key
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatalf("Failed to parse private key: %s", err)
	}

	// SSH client configuration
	config := &ssh.ClientConfig{
		User: *username, // Update with your username
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		// Specify HostKeyCallback function
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // InsecureIgnoreHostKey ignores all host key checks
		// You should use a more secure HostKeyCallback in a production environment.
	}

	// Connect to the SSH server
	client, err := ssh.Dial("tcp", *serverIP+":22", config)
	if err != nil {
		log.Fatalf("Failed to dial: %s", err)
	}
	defer client.Close()

	// Create a session
	session, err := client.NewSession()
	if err != nil {
		log.Fatalf("Failed to create session: %s", err)
	}
	defer session.Close()

	// Run the command on the server
	output, err := session.Output(*command)
	if err != nil {
		log.Fatalf("Failed to run command: %s", err)
	}

	// Print the output
	fmt.Println("Server output for command:", *command)
	fmt.Println(string(output))
}
