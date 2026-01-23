package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
)

func main() {
	var (
		serverURL = flag.String("url", "https://ctrlapp.merith.xyz", "Server base URL")
		username  = flag.String("username", "", "Login username")
		password  = flag.String("password", "", "Login password")
		register  = flag.Bool("register", false, "Register account if it doesn't exist")
	)
	flag.Parse()

	if *username == "" || *password == "" {
		log.Fatal("Username and password are required. Use -username and -password flags")
	}

	commander := &Commander{
		serverURL: *serverURL,
	}

	// Login first
	fmt.Println("🔐 Authenticating...")
	if err := commander.Login(*username, *password); err != nil {
		if *register {
			// If login failed and register flag is set, try to register the account
			fmt.Println("⚠️  Login failed, attempting to register new account...")
			screenName := strings.Title(*username) + " (Commander)"
			if err := commander.Register(screenName, *username, *password); err != nil {
				log.Fatalf("Failed to register: %v", err)
			}
			fmt.Println("✅ Account registered successfully!")

			// Now try to login again
			if err := commander.Login(*username, *password); err != nil {
				log.Fatalf("Failed to login after registration: %v", err)
			}
		} else {
			log.Fatalf("Failed to authenticate: %v\nTip: Use --register flag to create account if it doesn't exist", err)
		}
	}

	fmt.Printf("✅ Logged in as: %s (%s)\n", commander.user.ScreenName, commander.user.LoginName)
	fmt.Printf("   User ID: %s\n\n", commander.user.ID)

	// Connect to WebSocket
	fmt.Println("🔌 Connecting to WebSocket...")
	if err := commander.Connect(); err != nil {
		log.Fatalf("Failed to connect to WebSocket: %v", err)
	}
	defer commander.Close()

	fmt.Println("✅ Connected to WebSocket server!")
	fmt.Println()

	// Start TUI
	if err := RunTUI(commander); err != nil {
		log.Fatalf("TUI error: %v", err)
	}
}
