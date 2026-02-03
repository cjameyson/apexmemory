package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "user":
		runUserCmd(os.Args[2:])
	case "notebooks":
		runNotebooksCmd(os.Args[2:])
	case "reviews":
		runReviewsCmd(os.Args[2:])
	case "all":
		runAllCmd(os.Args[2:])
	case "help", "-h", "--help":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`seed - Apex Memory data seeding CLI

Usage:
  seed <command> [flags]

Commands:
  user       Create/verify the test agent account
  notebooks  Seed notebooks and facts for a user
  reviews    Generate review history (placeholder)
  all        Run user + notebooks for test agent
  help       Show this help message

Examples:
  seed user                                          # Create test agent account
  seed notebooks -email test-agent@apexmemory.ai    # Seed notebooks for user
  seed notebooks -email user@example.com -clear     # Clear and reseed
  seed all                                           # Full test agent setup

Run 'seed <command> -help' for command-specific help.`)
}
