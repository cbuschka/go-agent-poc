//go:generate go run ../generators/embedded_bytes_generator.go
package master

import (
	"log"
	"os"
)

func Run() error {
	log.SetOutput(os.Stderr)
	return runMaster()
}
