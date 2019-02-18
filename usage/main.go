package main

import (
	"context"
	pb "gRPC-Cache/description"
	"gRPC-Cache/utils"
	"google.golang.org/grpc"
	"io"
	"log"
)

const (
	address = "localhost:8888"
	max     = 100
)

func main() {
	runerConsumerClient()
}

func runerConsumerClient() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	conn, err := grpc.Dial(address, opts...)
	utils.HandleError(err)
	defer conn.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client := pb.NewConsumerClient(conn)

	stream, err := client.CacherRunner(ctx, &pb.Request{})
	utils.HandleError(err)

	for {
		_, err := stream.Recv()
		if err == io.EOF {
			break
		}
		utils.HandleError(err)
	}

	log.Println("wgConsumer.Done")
}
