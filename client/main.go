package main

import (
	"context"
	pb "gRPC-Cache/description"
	"gRPC-Cache/utils"
	"io"
	"log"
	"sync"

	"google.golang.org/grpc"
)

const (
	address = "localhost:9999"
	max     = 10
)

var (
	wgConsumer sync.WaitGroup
)

func main() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	wgConsumer.Add(max)
	for i := int32(1); i <= max; i++ {
		go func(request *pb.Request) {
			defer wgConsumer.Done()

			conn, err := grpc.Dial(address, opts...)
			utils.HandleError(err)
			defer conn.Close()

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			client := pb.NewCacherClient(conn)

			stream, err := client.GetRandomDataStream(ctx, request)
			utils.HandleError(err)

			for {
				//reply, err := stream.Recv()
				_, err := stream.Recv()
				if err == io.EOF {
					break
				}
				utils.HandleError(err)
			}
		}(&pb.Request{N: i})
	}
	//go printerRoutine()

	wgConsumer.Wait()
	log.Println("wgConsumer.Done")
}

//func printerRoutine() {
//	for w := range ch {
//		log.Println(w)
//	}
//}
