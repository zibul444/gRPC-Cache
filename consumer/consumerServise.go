//go:generate protoc -I ../description --go_out=plugins=grpc:../description ../description/descIDL.proto
// ~/go/src % protoc --go_out=plugins=grpc:. gRPC-Cache/cacher/description/descIDL.proto

package main

import (
	"context"
	pb "gRPC-Cache/description"
	"gRPC-Cache/utils"
	"io"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"
)

const (
	address = "localhost:8888"
	max     = 1000
)

type server struct {
	reply []*string
}

func (s *server) CacherRunner(reply *pb.Request, stream pb.Consumer_CacherRunnerServer) error {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	var wgConsumer sync.WaitGroup
	wgConsumer.Add(max)
	for i := int32(1); i <= max; i++ {
		go func(request *pb.Request) {
			defer wgConsumer.Done()

			conn, err := grpc.Dial("localhost:9999", opts...)
			utils.HandleError(err)
			defer conn.Close()

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			client := pb.NewCacherClient(conn)

			stream, err := client.GetRandomDataStream(ctx, request)
			utils.HandleError(err)

			for {
				_, err := stream.Recv()
				if err == io.EOF {
					break
				}
				utils.HandleError(err)
			}
		}(&pb.Request{N: i})
	}
	wgConsumer.Wait()
	log.Println("wgConsumer.Done")

	return nil
}

// Запуск Consumer сервера
func StartConsumerServer() {
	lis, err := getListener()

	grpcServer := RegisterConsumerServer()

	err = grpcServer.Serve(lis)
	utils.HandleError(err)
}

func RegisterConsumerServer() *grpc.Server {
	grpcServer := grpc.NewServer()
	pb.RegisterConsumerServer(grpcServer, &server{})
	log.Println("Register ConsumerServer success!")
	return grpcServer
}

func getListener() (net.Listener, error) {
	lis, err := net.Listen("tcp", address)
	utils.HandleError(err)
	return lis, err
}

func main() {
	StartConsumerServer()

}
