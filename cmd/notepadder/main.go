package main

import (
	"flag"
	"log"
	"notepadder/pkg/win"
)

func main() {
	var noNew bool
	flag.BoolVar(&noNew, "no-new", false, "Do not open a new tab")
	flag.BoolVar(&noNew, "n", false, "Shorthand for no-new")
	flag.Parse()
	if err := win.Run(noNew); err != nil {
		log.Fatal(err)
	}
} 