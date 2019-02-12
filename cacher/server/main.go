//go:generate protoc -I ../description --go_out=plugins=grpc:../description ../description/descIDL.proto
// ~/go/src % protoc --go_out=plugins=grpc:. gRPC-Cache/cacher/description/descIDL.proto

package main

import (
	"context"
	"fmt"
	pb "gRPC-Cache/cacher/description"
	"gRPC-Cache/cacher/resources"
	"gRPC-Cache/cacher/utils"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"google.golang.org/grpc"
)

const (
	port = "localhost:9999"
)

var (
	timeLife int
	url      string

	conf resources.Config
)

// описание используется для реализации description.serviceCacher
type server struct{}

func (s *server) GetRandomDataStream(ctx context.Context, in *pb.Request) (*pb.Reply, error) {
	log.Println("Received")

	conf = utils.GetConfig("cacher/resources/config.yml") // получаем конфиг
	//utils.HandleError(err)
	var data string
	for i := 0; i < conf.NumberOfRequests; i++ {
		url = utils.GetRandomUrl(conf)        // получаем случайный ресурс из конфига
		keysRAW := utils.Execute("KEYS", "*") // получили доступные ключи

		if strings.Contains(keysRAW, url) {
			data = utils.Execute("GET", url) // получаем закешированные ресурсы
			log.Println(data[:75])
		} else {
			resp, err := http.Get(url) // получаем данные
			utils.HandleError(err)
			defer utils.HandleError(resp.Body.Close()) // утилизируем ресурсы

			data = fmt.Sprint(resp)
			log.Println(data[:40])

			timeLife = utils.GetRandomTimeLife(conf)           //  получили случаейное время жизни в пределах заданных в конфиге
			utils.ExecuteCommand("SETEX", url, timeLife, resp) // положили данные в БД указали время жизни кэша

			time.Sleep(100 * time.Millisecond) // fixme удалить
		}
	}

	return &pb.Reply{Data: data}, nil // fixme пока отдает последний результат
}

func main() {
	lis, err := net.Listen("tcp", port)
	utils.HandleError(err)

	grpcServer := grpc.NewServer()
	pb.RegisterCacherServer(grpcServer, &server{})
	log.Println("Register CacherServer success! matherHucker")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
