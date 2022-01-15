package main

import (
	"fmt"
	"github.com/cbuschka/go-agent-poc/internal/master"
	"os"
)

func main() {
	err := master.Run()
	if err != nil {
		fmt.Printf("failure: %s", err.Error())
		os.Exit(1)
	}
}
