package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	// Parse command line flags
	serverURL := flag.String("server", "wss://ctrlapp.merith.xyz", "Server WebSocket URL")
	username := flag.String("username", "", "Username for authentication")
	password := flag.String("password", "", "Password for authentication")
	register := flag.Bool("register", false, "Register new account instead of login")
	flag.Parse()

	if *username == "" || *password == "" {
		fmt.Println("Error: username and password are required")
		flag.Usage()
		os.Exit(1)
	}

	// Create client
	client := NewClient(*serverURL)

	// Authenticate
	var token string
	var err error

	if *register {
		fmt.Printf("Registering new user: %s\n", *username)
		token, err = client.Register(*username, *password)
		if err != nil {
			log.Fatalf("Registration failed: %v", err)
		}
		fmt.Println("Registration successful!")
	} else {
		fmt.Printf("Logging in as: %s\n", *username)
		token, err = client.Login(*username, *password)
		if err != nil {
			// Try to register if login fails
			fmt.Printf("Login failed, attempting registration...\n")
			token, err = client.Register(*username, *password)
			if err != nil {
				log.Fatalf("Authentication failed: %v", err)
			}
			fmt.Println("Registration successful!")
		} else {
			fmt.Println("Login successful!")
		}
	}

	// Connect to WebSocket
	fmt.Println("Connecting to WebSocket...")
	if err := client.Connect(token); err != nil {
		log.Fatalf("WebSocket connection failed: %v", err)
	}

	fmt.Println("Connected! Starting client interface...")

	// Run TUI
	if err := RunTUI(client); err != nil {
		log.Fatalf("TUI error: %v", err)
	}
}
