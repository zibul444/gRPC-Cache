//go:generate protoc -I ../description --go_out=plugins=grpc:../description ../description/descIDL.proto

package main

import (
	"context"
	"log"
	"net"

	pb "gRPC-Cache/cacher/description"

	"google.golang.org/grpc"
)

const (
	port = "localhost:9999"
)

// описание используется для реализации description.serviceCacher
type server struct{}

func (s *server) GetRandomDataStream(ctx context.Context, in *pb.Request) (*pb.Reply, error) {
	//log.Printf("Received: %v", in.Name)
	//log.Printf("Received: %s", in.String())
	log.Println("Received")
	return &pb.Reply{Data: "Hello " + "World"}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterCacherServer(grpcServer, &server{})
	log.Println("Register CacherServer success! matherHucker")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
