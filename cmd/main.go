package main

import (
	"os"

	"tioga.co/kristofer/kg"
)

func main() {
	argv := os.Args // array of filenames to edit
	argc := len(argv)

	kg.StartEditor(argv, argc)
}
