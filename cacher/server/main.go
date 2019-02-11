//go:generate protoc -I ../description --go_out=plugins=grpc:../description ../description/descIDL.proto
// ~/go/src % protoc --go_out=plugins=grpc:. gRPC-Cache/cacher/description/descIDL.proto

package main

import (
	"context"
	"fmt"
	"gRPC-Cache/utils"
	"log"
	"net"
	"net/http"
	"time"

	pb "gRPC-Cache/cacher/description" //fixme забирать pd.go из ../resources

	"google.golang.org/grpc"
)

const (
	port = "localhost:9999"
)

var (
	timeLife int
	url      string
)

// описание используется для реализации description.serviceCacher
type server struct{}

func (s *server) GetRandomDataStream(ctx context.Context, in *pb.Request) (*pb.Reply, error) {
	log.Println("Received")

	conf := utils.GetConfig("resources/config.yml") // получаем конфиг
	url = utils.GetRandomUrl(conf)                  // получаем случайный ресурс из конфига
	keys := utils.GetListReadyKeys()                // получаем список доступных лючей

	if len(keys) > 0 {
		for _, k := range keys {
			fmt.Println(k)
		}
	}

	resp, err := http.Get(url) // получаем данные
	defer resp.Body.Close()    // утилизируем ресурсы
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	data := fmt.Sprint(resp)
	//log.Println(resp)
	log.Println(data[:100])

	timeLife = utils.GetRandomTimeLife(conf)
	utils.ExecuteCommand("SETEX", url, timeLife, resp)
	utils.Execute("GET", url)

	url = utils.GetRandomUrl(conf)
	resp, err = http.Get(url)

	utils.Execute("KEYS", "*")

	time.Sleep(100 * time.Millisecond)

	//log.Printf("Received: %s", in._())
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
