package main

import (
	"log"

	"github.com/luncj/mess/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatalf("failed to execute command: %s", err)
	}
}
