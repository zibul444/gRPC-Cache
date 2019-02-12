package main

import (
	"context"
	"gRPC-Cache/cacher/utils"
	"log"
	"time"

	pb "gRPC-Cache/cacher/description"

	"google.golang.org/grpc"
)

const (
	address = "localhost:9999"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	utils.HandleError(err)
	defer utils.HandleError(conn.Close())
	client := pb.NewCacherClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := client.GetRandomDataStream(ctx, &pb.Request{})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Received: %s", r.Data)
}
