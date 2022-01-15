package main

import (
	"fmt"
	"github.com/cbuschka/go-grpc-agent-poc/internal/agent"
	"os"
)

func main() {
	err := agent.Run()
	if err != nil {
		fmt.Printf("failure: %s", err.Error())
		os.Exit(1)
	}
}
