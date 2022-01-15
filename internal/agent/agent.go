package agent

import (
	"context"
	"github.com/cbuschka/go-grpc-agent-poc/internal/local"
	pb "github.com/cbuschka/go-grpc-agent-poc/internal/protocol/generated"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
)

type server struct {
	pb.UnimplementedAgentServer
}

func (s *server) Ping(ctx context.Context, in *pb.PingRequest) (*pb.PingReply, error) {
	log.Printf("Received ping with %s", in.MagicKey)
	return &pb.PingReply{MagicKey: in.GetMagicKey()}, nil
}

func runAgent() error {
	lis := local.NewListener()
	s := grpc.NewServer()
	pb.RegisterAgentServer(s, &server{})
	reflection.Register(s)
	log.Printf("Agent listening at %v...", lis.Addr())
	if err := s.Serve(lis); err != nil {
		return err
	}

	return nil
}
