package main

import (
	"context"
	"fmt"
	"log"

	orderpb "assignment/proto/orderpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	conn, err := grpc.Dial("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	client := orderpb.NewOrderServiceClient(conn)

	stream, err := client.SubscribeToOrderUpdates(context.Background(), &orderpb.OrderRequest{
		OrderId: "20260417022516",
	})
	if err != nil {
		log.Fatal(err)
	}

	for {
		res, err := stream.Recv()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("STATUS:", res.Status)
	}
}
