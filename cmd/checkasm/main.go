package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mmcloughlin/bugsalsa/finder"
)

func main() {
	for _, filename := range os.Args[1:] {
		searchfile(filename)
	}
}

func searchfile(filename string) {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Parse the instructions.
	instructions, err := finder.ParseAssembly(f)
	if err != nil {
		log.Fatal(err)
	}

	// Search for the bug.
	result := finder.Find(instructions)
	if result == nil {
		return
	}

	fmt.Printf("%s: bug found at line %d\n", filename, result.StartLine())
}
