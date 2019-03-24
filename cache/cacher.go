//go:generate protoc -I ../description --go_out=plugins=grpc:../description ../description/descIDL.proto

package cache

import (
	"fmt"
	pb "gRPC-Cache/description"
	"gRPC-Cache/utils"
	"github.com/op/go-logging"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	_ "expvar"

	"google.golang.org/grpc"
)

const (
	address = "localhost:9999"
)

var (
	conf *utils.Config

	logger = logging.MustGetLogger("cache")
	format = utils.GetFormatter()

	chReturnUrls chan<- string
	chUrl        <-chan string
	chData       = make(chan string)
)

// описание используется для реализации description.serviceCacher
type server struct {
	reply []*string
}

func init() {
	filepath.Base()
	conf = utils.GetConfig("config.yml") // получаем конфиг

	chUrl = conf.ChGetUrls
	chReturnUrls = conf.ChReturnUrls

	backend := logging.NewLogBackend(os.Stdout, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)
	logging.SetBackend(backendFormatter)
}

func (s *server) GetRandomDataStream(reply *pb.Request, stream pb.Cacher_GetRandomDataStreamServer) error {
	//logger.Notice("Received:", reply.N)
	var (
		waitGroup sync.WaitGroup
		data      string
	)

	waitGroup.Add(conf.NumberOfRequests)
	for i := 0; i < conf.NumberOfRequests; i++ {
		go func(chData chan<- string, recNum int32, num int) {
			url, ok := <-chUrl
			if !ok {
				return
			}
			ok = checkCashAlive(url)

			if ok {
				chData <- utils.Execute("GET", url)
			} else {
				data := getDataFromResource(url)
				go func() {
					chData <- data
				}()

				timeLife := utils.GetRandomTimeLife(*conf)
				utils.ExecuteCommand("SETEX", url, timeLife, data)
			}

			chReturnUrls <- url // ОЧЕНЬ важная штука(Вернуть URL)
			waitGroup.Done()
		}(chData, reply.N, i) // передали канал для получения данных ресурса, номер клиента, и порядковый новмер

		sendStreamData(data, stream)
	}

	waitGroup.Wait()
	logger.Notice("FOR is End", reply.N)

	return nil
}

func sendStreamData(data string, stream pb.Cacher_GetRandomDataStreamServer) {
	data, ok := <-chData
	if !ok {
		return
	}
	err := stream.Send(&pb.Reply{Data: data})
	utils.HandleError(err)
}

// Получение данных от ресурса
func getDataFromResource(url string) (dataResource string) {
	//logger.Notice("GetDataFromResource url", url)
	resp, err := http.Get(url)
	utils.HandleError(err)
	defer resp.Body.Close()

	dataResource = fmt.Sprint(resp) // преобразовали данные для отправки
	return dataResource
}

// Проверка наличия живого КЭШа в БД
func checkCashAlive(url string) (have bool) {
	keys := utils.Execute("KEYS", url)

	have = checkCash(keys, url)

	if have {
		ttl := utils.ToInt64(utils.ExecuteCommand("TTL", url))
		if ttl < 1 {
			have = !have
			logger.Notice("Не успеваем взять КЭШ. TTL is:", ttl)
		}
	}
	return have
}

// Проверка наличия КЭШа в БД
func checkCash(keys string, url string) (have bool) {
	have = strings.Contains(keys, url)
	return have
}

// Запуск сервера
func StartServerCacher() {
	lis, err := getListener()
	utils.HandleError(err)

	grpcServer := registerCacherServer()

	go grpcServer.Serve(lis)
}

func registerCacherServer() *grpc.Server {
	grpcServer := grpc.NewServer()
	pb.RegisterCacherServer(grpcServer, &server{})
	logger.Notice("Register CacherServer success!")
	return grpcServer
}

func getListener() (net.Listener, error) {
	lis, err := net.Listen("tcp", address)
	utils.HandleError(err)
	return lis, err
}
