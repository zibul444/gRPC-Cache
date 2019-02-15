//go:generate protoc -I ../description --go_out=plugins=grpc:../description ../description/descIDL.proto
// ~/go/src % protoc --go_out=plugins=grpc:. gRPC-Cache/cacher/description/descIDL.proto

package main

import (
	"fmt"
	pb "gRPC-Cache/description"
	"gRPC-Cache/utils"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"google.golang.org/grpc"
)

const (
	address = "localhost:9999"
)

var (
	timeLife int
	url      string
	logger   = log.New(os.Stdout, fmt.Sprint(time.Now().Format(time.StampNano))+": ", log.Lshortfile)
	conf     *utils.Config

	chReturn = make(chan string)
)

// описание используется для реализации description.serviceCacher
type server struct {
	reply []*string
}

func init() {
	conf = utils.GetConfig("config.yml") // получаем конфиг
	go conf.ReturnURL(chReturn)
}

func (s *server) GetRandomDataStream(reply *pb.Request, stream pb.Cacher_GetRandomDataStreamServer) error {
	logger.Println("Received", reply)
	//defer func() { chReturn <- url }() // fixme
	logger.Println("1 Объявляем переменные")

	var (
		chUrl  = make(chan string)
		chData = make(chan string)
		data   string // обьявили ответ
		wg1    = sync.WaitGroup{}
		wg2    = sync.WaitGroup{}
		//wg3 = sync.WaitGroup{}
	)

	logger.Println("2 Объявили переменные переменные, начинаем цикл")
	for i := 0; i < conf.NumberOfRequests; i++ {
		logger.Println("2.1 for started")
		//wg3.Add(1)

		go func(ch chan<- string) { // fixme cache go!!!
			logger.Println("-3 Добавили +1 ожидание в группу wg1")
			wg1.Add(1)
			logger.Println("-4 Ждем получение случайного ресурса из конфига")
			go conf.TakeURL(chUrl) // получаем случайный ресурс из конфига
			logger.Println("-4.1 \"Недождались\" идем дальше")
			logger.Println("-4.2 Добавили +1 ожидание в группу wg2")
			wg2.Add(1)
			go func() {
				logger.Println("--5 Запустили ГоПодпрограмму для ожидания получения случайного ресурса из конфига <- chUrl")
				//debug.PrintStack()
				url = <-chUrl
				logger.Println("--5.1 Получили случайный ресурс от конфига")
				wg2.Done()
				logger.Println("--6 Убрали ожидание из группы wg2  -1")
			}()
			logger.Println("-7 Ждем группу wg2")
			wg2.Wait()
			logger.Println("-7.1 Дождались группу wg2")
			keysRAW := utils.Execute("KEYS", "*") // получили доступные кэши
			logger.Println("-8 Получили доступные ключи")

			if strings.Contains(keysRAW, url) {
				logger.Println("-9 Ключь найден в БД, забросим в ch <-")
				ttl := utils.ToInt64(utils.ExecuteCommand("TTL", url))
				if ttl < 3 {
					logger.Println("-9.1 Не успеваем мы взять КЭШ") // TODO
				}
				ch <- utils.Execute("GET", url) // получаем закешированные ресурсы
				logger.Println("-BD:", url, data[:len(data)/10])
				wg1.Done()
				logger.Println("-10 Завершился iF ELSE")
			} else {
				logger.Println("-11 URL не найден в БД, запросили http.Get")
				resp, err := http.Get(url) // получаем данные
				utils.HandleError(err)
				logger.Println("-12 http.Get отработал")
				defer resp.Body.Close() // утилизируем ресурсы
				logger.Println("-13 defer resp.Body.Close()")
				go func() {
					logger.Println("--13.1 Запустили ГоПодпрограмму, конвертируем данные от респонса") // fixme
					data := fmt.Sprint(resp)                                                           // преобразовали данные для отправки
					logger.Println("--13.2 Отправляем данные ch <-")
					ch <- fmt.Sprint(data)
					logger.Println("--13.2 Кто-то забрал данные из <- ch")
				}()

				logger.Println("-13.3 go func() - Запущена")
				logger.Println("-TCP:", url, data[:len(data)/10])
				logger.Println("-14 Получаем время жизни кэша")
				timeLife = utils.GetRandomTimeLife(*conf) //  получили случаейное время жизни в пределах заданных в конфиге
				logger.Println("-14.1 Отправляем кэш в БД")
				utils.ExecuteCommand("SETEX", url, timeLife, resp) // положили данные в БД указали время жизни кэша
				wg1.Done()
				logger.Println("-15 Завершился iF ELSE")
			}
		}(chData) // передали канал для получения данных ресурса
		logger.Println("16 Ждем группу wg1")
		wg1.Wait()
		logger.Println("16.1 Дождались группу wg1")
		data = <-chData
		logger.Println("17 Получили данные из <- chData")
		err := stream.Send(&pb.Reply{Data: data})
		logger.Println("17.1 Отправили данные в stream")
		chReturn <- url // ОЧЕНЬ важная штука
		logger.Println("18 ОЧЕНЬ важная штука отработала.")
		time.Sleep(time.Millisecond * 400) // fixme удалить, ожидание для наглядного вывода
		utils.HandleError(err)
		logger.Println("18.1 Выспались...")
	}

	logger.Println("FOR is End")

	return nil
}

func main() {
	lis, err := net.Listen("tcp", address)
	utils.HandleError(err)

	grpcServer := grpc.NewServer()
	pb.RegisterCacherServer(grpcServer, &server{})
	logger.Println("Register CacherServer success! matherHucker")
	if err := grpcServer.Serve(lis); err != nil {
		logger.Fatalf("failed to serve: %v", err)
		panic(err)
	}
}
