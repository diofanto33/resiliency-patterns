package server

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/diofanto33/resiliency-patterns/circuit-breaker/user"
	"google.golang.org/grpc"
)

type Server struct {
	user.UnimplementedUserServiceServer
}

func (s *Server) CreateUser(ctx context.Context, in *user.CreateUserRequest) (*user.CreateUserResponse, error) {
	return &user.CreateUserResponse{UserId: 1}, nil
}

func main() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", 8080))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	user.RegisterUserServiceServer(grpcServer, &Server{})
	grpcServer.Serve(listener)
}
