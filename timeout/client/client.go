package main

import (
	"context"
	"log"

	"github.com/diofanto33/resiliency-patterns/timeout/product"
	"google.golang.org/grpc"
)

func main() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial("localhost:8080", opts...)
	if err != nil {
		log.Fatalf("Failed to connect product service. Err: %v", err)
	}

	defer conn.Close()

	produtClient := product.NewProductServiceClient(conn)

	log.Println("Creating product...")
	_, errCreate := produtClient.Create(context.Background(), &product.CreateProductRequest{Name: "diofanto33", Code: 2424, Price: 12.3})
	if errCreate != nil {
		log.Printf("Failed to create product. Err: %v", errCreate)
	} else {
		log.Println("Product is created successfully.")
	}
}
