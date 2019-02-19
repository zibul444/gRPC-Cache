//go:generate protoc -I ../description --go_out=plugins=grpc:../description ../description/descIDL.proto

package main

import (
	"context"
	"gRPC-Cache/cache"
	"gRPC-Cache/consumer"
	pb "gRPC-Cache/description"
	"gRPC-Cache/utils"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
)

const (
	address = "localhost:8888"
)

func main() {
	go cache.StartServerCacher()
	time.Sleep(time.Second)

	go consumer.StartConsumerServer()
	time.Sleep(time.Second)

	runnerConsumerClient()

	time.Sleep(time.Second)
	log.Println("Done")
}

func runnerConsumerClient() {
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
}
