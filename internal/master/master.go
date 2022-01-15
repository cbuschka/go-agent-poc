//go:generate go run ../generators/embedded_bytes_generator.go
package master

import (
	"context"
	pb "github.com/cbuschka/go-agent-poc/internal/protocol/generated"
	"github.com/cbuschka/go-agent-poc/internal/remote"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/user"
	"time"
)

const (
	magicKey = "42"
)

func runMaster() error {
	log.SetOutput(os.Stderr)

	var session remote.Session
	if len(os.Args) == 2 && os.Args[1] == "ssh" {
		log.Printf("Connecting via ssh...")

		currUser, err := user.Current()
		if err != nil {
			return err
		}

		session, err = remote.OpenSshSession(currUser.Username, "localhost", 22)
		if err != nil {
			return err
		}
	} else {
		var err error
		session, err = remote.OpenLocalSession()
		if err != nil {
			return err
		}
	}
	defer session.Close()

	agentBytes, err := ioutil.ReadFile("./dist/agent")
	if err != nil {
		return err
	}

	err = session.PlaceAgent("/tmp/agent", agentBytes)
	if err != nil {
		return err
	}

	log.Printf("Starting agent...")
	remoteAgent, err := session.StartAgent("/tmp/agent")
	defer remoteAgent.Close()

	log.Printf("Starting grpc session...")
	conn, err := grpc.Dial("", grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return remoteAgent.GetConnection()
		}))

	agentClient := pb.NewAgentClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	log.Printf("Sending ping with magic key %s...", magicKey)
	response, err := agentClient.Ping(ctx, &pb.PingRequest{MagicKey: magicKey})
	if err != nil {
		return err
	}
	log.Printf("Got a reply with magic key %s", response.MagicKey)

	return nil
}
