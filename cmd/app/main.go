package main

import (
	"log"
	"net"

	"assignment/internal/repository"
	"assignment/internal/transport/grpcserver"		
	
	http "assignment/internal/transport/http"
	"assignment/internal/usecase"
	orderpb "assignment/proto/orderpb"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	paymentClient := usecase.NewPaymentGRPC(conn)

	repo, err := repository.NewOrderPostgresRepository()
	if err != nil {
		log.Fatal(err)
	}

	uc := usecase.NewOrderUsecase(repo, paymentClient)

	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()
	orderpb.RegisterOrderServiceServer(grpcServer, grpcserver.NewOrderGRPCServer(uc))

	go func() {
		log.Println("Order gRPC streaming server running on :50052")
		if err := grpcServer.Serve(lis); err != nil {
			log.Println("gRPC server error:", err)
		}
	}()

	r := gin.Default()
	http.NewOrderHandler(r, uc)

	log.Println("Order Service running on :8080")
	r.Run(":8080")
}
