package main

import (
	"os"

	"github.com/kristofer/kg"
)

func main() {
	argv := os.Args // array of filenames to edit
	argc := len(argv)
	edit := &kg.Editor{}
	edit.StartEditor(argv, argc)
}
