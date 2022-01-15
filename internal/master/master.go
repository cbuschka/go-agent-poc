package master

import (
	"context"
	"flag"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"

	pb "github.com/cbuschka/go-agent-poc/internal/protocol/generated"
	"google.golang.org/grpc"
)

const (
	addr     = "localhost:50051"
	magicKey = "42"
)

func Run() error {
	flag.Parse()
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewAgentClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	log.Printf("Sending ping with magic key %s...", magicKey)
	r, err := c.Ping(ctx, &pb.PingRequest{MagicKey: magicKey})
	if err != nil {
		return err
	}
	log.Printf("Got a reply with magic key %s", r.MagicKey)

	return nil
}
