package main

import (
	"flag"
	"fmt"
	"log"
	"notepadder/pkg/win"
	"os"
)

func main() {
	var noNew bool
	var debug bool
	flag.BoolVar(&noNew, "no-new", false, "Do not open a new tab")
	flag.BoolVar(&noNew, "n", false, "Shorthand for no-new")
	flag.BoolVar(&debug, "debug", false, "Print debug output to console (requires build without -H=windowsgui)")
	flag.Parse()

	if debug {
		fmt.Println("Debug mode enabled.")
	}

	if err := win.Run(noNew, debug); err != nil {
		if debug {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1) // Exit directly if debugging
		} else {
			log.Fatal(err) // Use log.Fatal for non-debug GUI mode (might not be visible)
		}
	}
	if debug {
		fmt.Println("Operation completed successfully.")
	}
} 