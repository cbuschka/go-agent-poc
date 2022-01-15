package agent

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "github.com/cbuschka/go-agent-poc/internal/protocol/generated"
	"google.golang.org/grpc"
)

const (
	port = 50051
)

type server struct {
	pb.UnimplementedAgentServer
}

func (s *server) Ping(ctx context.Context, in *pb.PingRequest) (*pb.PingReply, error) {
	log.Printf("Received ping with %s", in.MagicKey)
	return &pb.PingReply{MagicKey: in.GetMagicKey()}, nil
}

func Run() error {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		return err
	}
	s := grpc.NewServer()
	pb.RegisterAgentServer(s, &server{})
	log.Printf("Agent listening at %v...", lis.Addr())
	if err := s.Serve(lis); err != nil {
		return err
	}

	return nil
}
