package main

import (
	"flag"
	"fmt"
	"os"
)

func runReviewsCmd(args []string) {
	fs := flag.NewFlagSet("reviews", flag.ExitOnError)
	var (
		dsn   string
		email string
	)
	fs.StringVar(&dsn, "dsn", "", "PostgreSQL DSN (or set DATABASE_URL env var)")
	fs.StringVar(&email, "email", "", "User email to generate reviews for (required)")
	fs.Usage = func() {
		fmt.Println(`seed reviews - Generate review history (placeholder)

Usage:
  seed reviews -email <email> [flags]

Flags:`)
		fs.PrintDefaults()
		fmt.Println(`
Note: This command is a placeholder for future implementation.`)
	}
	if err := fs.Parse(args); err != nil {
		os.Exit(1)
	}

	if email == "" {
		fmt.Fprintln(os.Stderr, "error: -email flag is required")
		fs.Usage()
		os.Exit(1)
	}

	fmt.Printf("seed reviews: placeholder - would generate reviews for %s\n", email)
	fmt.Println("This feature is not yet implemented.")
}

func runAllCmd(args []string) {
	fs := flag.NewFlagSet("all", flag.ExitOnError)
	var (
		dsn       string
		notebooks int
		facts     int
	)
	fs.StringVar(&dsn, "dsn", "", "PostgreSQL DSN (or set DATABASE_URL env var)")
	fs.IntVar(&notebooks, "notebooks", 10, "Number of notebooks to create")
	fs.IntVar(&facts, "facts", 15, "Number of facts per notebook")
	fs.Usage = func() {
		fmt.Println(`seed all - Full test agent setup (user + notebooks)

Usage:
  seed all [flags]

Flags:`)
		fs.PrintDefaults()
		fmt.Println(`
This command:
  1. Creates/verifies the test agent account
  2. Seeds notebooks and facts for the test agent`)
	}
	if err := fs.Parse(args); err != nil {
		os.Exit(1)
	}

	// Step 1: Create test user
	fmt.Println("Step 1: Creating test agent account...")
	if err := createTestUser(dsn); err != nil {
		fmt.Fprintf(os.Stderr, "error creating test user: %v\n", err)
		os.Exit(1)
	}

	// Step 2: Seed notebooks
	fmt.Println("Step 2: Seeding notebooks...")
	if err := seedNotebooks(dsn, testAgentEmail, notebooks, facts, false); err != nil {
		fmt.Fprintf(os.Stderr, "error seeding notebooks: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Test agent setup complete!")
}
