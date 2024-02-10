package main

import (
	"os"
	"flag"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
)

func main() {
	// Define command-line flags for server IP, username, authentication method, and command
	serverIP := flag.String("server", "", "Server IP address")
	username := flag.String("username", "", "User Name")
	password := flag.String("password", "", "Password for SSH authentication")
	privateKeyFile := flag.String("key", "", "Private key file for SSH authentication")
	command := flag.String("command", "", "Command to execute on the server")

	// Set up custom usage message
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options]\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "Options:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	// Check if server IP, username, authentication method, and command are provided
	if *serverIP == "" || *username == "" || (*password == "" && *privateKeyFile == "") || *command == "" {
		log.Fatal("Server IP, username, authentication method (password or private key file), and command must be provided")
	}

	var authMethods []ssh.AuthMethod

	// Use password authentication if a password is provided
	if *password != "" {
		authMethods = append(authMethods, ssh.Password(*password))
	} else if *privateKeyFile != "" { // Use private key authentication if a private key file is provided
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

		authMethods = append(authMethods, ssh.PublicKeys(signer))
	} else {
		log.Fatal("Either password or private key file must be provided for authentication")
	}

	// SSH client configuration
	config := &ssh.ClientConfig{
		User: *username,
		Auth: authMethods,
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
