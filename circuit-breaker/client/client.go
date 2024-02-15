package main

import (
	"context"
	"log"
	"time"

	"github.com/diofanto33/resiliency-patterns/circuit-breaker/middleware"
	"github.com/diofanto33/resiliency-patterns/circuit-breaker/user"
	"github.com/sony/gobreaker"
	"google.golang.org/grpc"
)

var cb *gobreaker.CircuitBreaker

func main() {
	cb = gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        "demo",
		MaxRequests: 3,
		Timeout:     4,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return failureRatio >= 0.1
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			log.Printf("Circuit Breaker: %s, changed from %v, to %v", name, from, to)
		},
	})
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithUnaryInterceptor(middleware.CircuitBreakerClientInterceptor(cb)))
	conn, err := grpc.Dial("localhost:8080", opts...)
	if err != nil {
		log.Fatalf("Failed to connect order service. Err: %v", err)
	}

	defer conn.Close()

	userClient := user.NewUserServiceClient(conn)
	for {
		log.Println("Creating User...")
		userResponse, errCreate := userClient.CreateUser(context.Background(), &user.CreateUserRequest{
			Name:     "Diego",
			Email:    "diofanto33@proton.me",
			Password: "pepe123",
		})

		if errCreate != nil {
			log.Printf("Failed to create user. Err: %v", errCreate)
		} else {
			log.Printf("User %d is created successfully.", userResponse.UserId)
		}
		time.Sleep(1 * time.Second)
	}

}
