//go:generate protoc -I ../description --go_out=plugins=grpc:../description ../description/descIDL.proto

package consumer

import (
	"context"
	pb "gRPC-Cache/description"
	"gRPC-Cache/utils"
	"github.com/op/go-logging"
	"io"
	"net"
	"sync"

	_ "expvar"

	"google.golang.org/grpc"
)

const (
	address       = "localhost:8888"
	addressCacher = "localhost:9999"
	max           = 1000
)

type server struct {
	reply []*string
}

var (
	logger = logging.MustGetLogger("consumer")
	format = utils.GetFormatter()
)

func (s *server) CacherRunner(reply *pb.Request, stream pb.Consumer_CacherRunnerServer) error {
	var (
		opts       []grpc.DialOption
		wgConsumer sync.WaitGroup
	)
	opts = append(opts, grpc.WithInsecure())
	wgConsumer.Add(max)
	for i := int32(1); i <= max; i++ {
		go func(request *pb.Request) {
			defer wgConsumer.Done()

			conn, err := grpc.Dial(addressCacher, opts...)
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
	logger.Notice("wgConsumer.Done")

	return nil
}

// Запуск Consumer сервера
func StartConsumerServer() {
	lis, err := getListener()

	grpcServer := registerConsumerServer()

	err = grpcServer.Serve(lis)
	utils.HandleError(err)
}

func registerConsumerServer() *grpc.Server {
	grpcServer := grpc.NewServer()
	pb.RegisterConsumerServer(grpcServer, &server{})
	logger.Notice("Register ConsumerServer success!")
	return grpcServer
}

func getListener() (net.Listener, error) {
	lis, err := net.Listen("tcp", address)
	utils.HandleError(err)
	return lis, err
}
