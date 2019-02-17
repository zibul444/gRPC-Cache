//go:generate protoc -I ../description --go_out=plugins=grpc:../description ../description/descIDL.proto
// ~/go/src % protoc --go_out=plugins=grpc:. gRPC-Cache/cacher/description/descIDL.proto

package main

import (
	"fmt"
	pb "gRPC-Cache/description"
	"gRPC-Cache/utils"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	address = "localhost:9999"
)

var (
	logger = log.New(os.Stdout, fmt.Sprint(time.Now().Format(time.StampNano))+": ", log.Lshortfile)
	conf   *utils.Config

	chReturnUrls = make(chan string)
	//
	chUrl  = make(chan string)
	chData = make(chan string)
	url    string
	data   string
	//wg1    = sync.WaitGroup{}
	waitGroup = sync.WaitGroup{}
)

// описание используется для реализации description.serviceCacher
type server struct {
	reply []*string
}

func init() {
	conf = utils.GetConfig("config.yml") // получаем конфиг
	go conf.ReturnURL(chReturnUrls)
}

func (s *server) GetRandomDataStream(reply *pb.Request, stream pb.Cacher_GetRandomDataStreamServer) error {
	logger.Println("Received:", reply.N)

	//wg1.Add(conf.NumberOfRequests)
	waitGroup.Add(conf.NumberOfRequests)
	for i := 0; i < conf.NumberOfRequests; i++ {
		//wg1.Add(1)
		go func(ch chan<- string, recNum int32, num int) {

			go conf.TakeURL(chUrl)

			url = <-chUrl

			keys := utils.Execute("KEYS", "*")

			ok := strings.Contains(keys, url)
			logger.Println("Contains keys is :", ok, ":", keys)

			if ok {
				ttl := utils.ToInt64(utils.Execute("TTL", url))
				logger.Println(recNum, num, "-14 Узнали срок годности КЭШа:", ttl)
				if ttl < 2 {
					ok = !ok
					logger.Println(recNum, num, "-14.1 Не успеваем взять КЭШ, keys is:", ok)
				}
			}

			if ok {
				logger.Println(recNum, num, "-14.2 Ключь найден в БД, забросим в ch <-")
				ch <- utils.Execute("GET", url) // получаем закешированные ресурсы
				logger.Println("-BD:", url, data[:len(data)/10])
				logger.Println(recNum, num, "-15 Завершился iF")
			} else {
				logger.Println(recNum, num, "-16 URL не найден в БД, запросили http.Get")
				resp, err := http.Get(url) // получаем данные
				utils.HandleError(err)
				logger.Println(recNum, num, "-17 http.Get отработал")
				//defer resp.Body.Close() // утилизируем ресурсы

				logger.Println(recNum, num, "--18.1 Конвертируем данные от респонса")
				data := fmt.Sprint(resp) // преобразовали данные для отправки
				go func() {
					logger.Println(recNum, num, "--18.2 Запустили ГоПодпрограмму. Отправляем данные ch <-")
					ch <- data
				}()

				logger.Println(recNum, num, "-19 go func() - Запущена")
				logger.Println(recNum, num, "-TCP:", url, data[:len(data)/10])

				timeLife := utils.GetRandomTimeLife(*conf) //  получили случаейное время жизни в пределах заданных в конфиге
				logger.Println(recNum, num, "-20 Получаем время жизни кэша:", timeLife)
				logger.Println(recNum, num, "-21 Отправляем кэш в БД")
				utils.ExecuteCommand("SETEX", url, timeLife, data) // положили данные в БД указали время жизни кэша
				logger.Println(recNum, num, "-22 Завершился ELSE")
			}
			//wg1.Done()
			defer waitGroup.Done()
		}(chData, reply.N, i) // передали канал для получения данных ресурса

		//wg1.Wait()
		//logger.Println(i, "24 Дождались группу wg1")
		data = <-chData
		logger.Println(i, "25 Получили данные из <- chData")
		//err := stream.Send(&pb.Reply{Data: data})
		if err := stream.Send(&pb.Reply{Data: data}); err != nil {
			return err
		}
		logger.Println(i, "26 Отправили данные в stream")
		chReturnUrls <- url // ОЧЕНЬ важная штука
		logger.Println(i, "27 ОЧЕНЬ важная штука отработала.")
	}

	waitGroup.Wait()
	logger.Println("FOR is End")

	return nil
}

func main() {

	defer close(chUrl)
	defer close(chData)
	defer close(chReturnUrls)

	lis, err := net.Listen("tcp", address)
	utils.HandleError(err)

	//grpcServer := grpc.NewServer()
	//pb.RegisterCacherServer(grpcServer, &server{})
	//logger.Println("Register CacherServer success! matherHucker")
	//if err := grpcServer.Serve(lis); err != nil {
	//	logger.Fatalf("failed to serve: %v", err)
	//}

	grpcServer := grpc.NewServer()

	pb.RegisterCacherServer(grpcServer, &server{})
	logger.Println("Register CacherServer success! matherHucker")
	err = grpcServer.Serve(lis)
	utils.HandleError(err)

}

//func newServer() *pb.CacherServer {
//	s := &pb.CacherServer{}
//	s.loadFeatures(*jsonDBFile)
//	return s
//}
