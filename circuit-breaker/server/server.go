package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"

	"github.com/diofanto33/resiliency-patterns/circuit-breaker/user"
	"google.golang.org/grpc"
)

type Server struct {
	user.UnimplementedUserServiceServer
}

func (s *Server) CreateUser(ctx context.Context, in *user.CreateUserRequest) (*user.CreateUserResponse, error) {
	var err error
	rand.Seed(time.Now().UnixNano())
	if rand.Intn(2) == 1 {
		err = errors.New("create user error")
	}
	return &user.CreateUserResponse{UserId: 11111}, err
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
