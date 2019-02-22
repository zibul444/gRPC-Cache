//go:generate protoc -I ../description --go_out=plugins=grpc:../description ../description/descIDL.proto

package main

import (
	"context"
	"gRPC-Cache/cache"
	"gRPC-Cache/consumer"
	pb "gRPC-Cache/description"
	"gRPC-Cache/utils"
	"github.com/op/go-logging"
	"io"
	"os"
	"time"

	_ "expvar"

	"google.golang.org/grpc"
)

const (
	address = "localhost:8888"
)

var (
	logger = logging.MustGetLogger("utils")
	format = logging.MustStringFormatter(
		`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:04x}%{color:reset} %{message}`,
	)
)

func main() {
	backend := logging.NewLogBackend(os.Stdout, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)
	logging.SetBackend(backendFormatter)

	go cache.StartServerCacher()
	time.Sleep(400 * time.Millisecond)

	go consumer.StartConsumerServer()
	time.Sleep(400 * time.Millisecond)

	runnerConsumerClient()

	time.Sleep(400 * time.Millisecond)
	logger.Notice("Done")
}

func runnerConsumerClient() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	conn, err := grpc.Dial(address, opts...)
	utils.HandleError(err)
	defer conn.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client := pb.NewConsumerClient(conn)

	stream, err := client.CacherRunner(ctx, &pb.Request{})
	utils.HandleError(err)

	for {
		_, err := stream.Recv()
		if err == io.EOF {
			break
		}
		utils.HandleError(err)
	}
}
