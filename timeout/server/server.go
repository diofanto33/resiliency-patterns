package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"

	"github.com/diofanto33/resiliency-patterns/timeout/product"
	"google.golang.org/grpc"
)

type server struct {
	product.UnimplementedProductServiceServer
}

func (s *server) Create(ctx context.Context, in *product.CreateProductRequest) (*product.CreateProductResponse, error) {
	time.Sleep(2 * time.Second)
	return &product.CreateProductResponse{ProductId: 1243}, nil
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	// Simulate that we have two goroutines that are doing some work in the background
	go func() {
		res := randomFunc(ctx, "a")
		log.Println(res)
		cancel()
	}()
	go func() {
		res := randomFunc(ctx, "b")
		log.Println(res)
		cancel()
	}()

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", 8080))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	product.RegisterProductServiceServer(grpcServer, &server{})
	grpcServer.Serve(listener)
}

func randomFunc(ctx context.Context, name string) string {
	rand.Seed(time.Now().UnixNano())
	min := 3
	max := 7
	sleepTime := rand.Intn(max-min+1) + min
	time.Sleep(time.Duration(sleepTime) * 1000000)
	return "hello from " + name
}
