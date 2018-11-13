package tkg

import (
	"log"
	"os"
)

// SetupLogFile logs all the log junk to a file
func SetupLogFile() {
	f, err := os.OpenFile("logfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
	log.Println("Start of Log...")
}
