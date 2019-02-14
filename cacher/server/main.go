//go:generate protoc -I ../description --go_out=plugins=grpc:../description ../description/descIDL.proto
// ~/go/src % protoc --go_out=plugins=grpc:. gRPC-Cache/cacher/description/descIDL.proto

package main

import (
	"fmt"
	pb "gRPC-Cache/cacher/description"
	"gRPC-Cache/cacher/utils"
	"log"
	"net"
	"net/http"
	"os"
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
	logger   = log.New(os.Stdout, fmt.Sprint(time.Now().Format(time.StampNano))+": ", log.Lshortfile)
	conf     *utils.Config
	ch       chan string
	chReturn chan string
)

// описание используется для реализации description.serviceCacher
type server struct {
	reply []*string
}

func init() {
	conf = utils.GetConfig("cacher/resources/config.yml") // получаем конфиг
	go conf.ReturnURL(chReturn)
}

func (s *server) GetRandomDataStream(reply *pb.Request, stream pb.Cacher_GetRandomDataStreamServer) error {
	logger.Println("Received", reply)
	defer func() { chReturn <- url }()
	for i := 0; i < conf.NumberOfRequests; i++ {
		handle(stream) // fixme server go!!!
		logger.Println("End")
	}

	return nil
}

func handle(stream pb.Cacher_GetRandomDataStreamServer) {
	var data string // обьявили ответ

	go conf.TakeURL(ch) // получаем случайный ресурс из конфига
	url = <-ch

	keysRAW := utils.Execute("KEYS", "*") // получили доступные кэши

	if strings.Contains(keysRAW, url) {
		data = utils.Execute("GET", url) // получаем закешированные ресурсы
		logger.Println("BD:", url, data[:75])
	} else {
		resp, err := http.Get(url) // получаем данные
		utils.HandleError(err)

		data = fmt.Sprint(resp) // преобразовали данные
		logger.Println("TCP:", url, data[:40])

		timeLife = utils.GetRandomTimeLife(*conf)          //  получили случаейное время жизни в пределах заданных в конфиге
		utils.ExecuteCommand("SETEX", url, timeLife, resp) // положили данные в БД указали время жизни кэша

		resp.Body.Close() // утилизируем ресурсы
	}
	err := stream.Send(&pb.Reply{Data: data})
	time.Sleep(time.Millisecond * 400) // fixme удалить, ожидание для наглядного вывода
	utils.HandleError(err)
}

func main() {
	lis, err := net.Listen("tcp", port)
	utils.HandleError(err)

	grpcServer := grpc.NewServer()
	pb.RegisterCacherServer(grpcServer, &server{})
	logger.Println("Register CacherServer success! matherHucker")
	if err := grpcServer.Serve(lis); err != nil {
		logger.Fatalf("failed to serve: %v", err)
		panic(err)
	}
}
