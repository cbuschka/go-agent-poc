package agent

import (
	"log"
	"os"
)

func Run() error {
	log.SetOutput(os.Stderr)
	log.SetFlags(0)

	return runAgent()
}
