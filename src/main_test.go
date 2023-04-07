package main

import (
	"flag"
	"os"
	"testing"
)

func TestMainFunction(t *testing.T) {
	// Back up the original command-line arguments and restore them after the test.
	origArgs := os.Args
	defer func() { os.Args = origArgs }()

	// Set up test command-line arguments.
	testArgs := []string{"cmd", "-start-dir", ".", "-output", "test_output.json"}
	os.Args = testArgs

	// Reset the flag set to avoid "flag provided but not defined" error.
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	// Capture any panic that might occur in the main function.
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("The main function panicked: %v", r)
		}
	}()

	// Call the main function.
	main()
}
