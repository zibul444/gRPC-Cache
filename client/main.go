package main

import (
	"context"
	"fmt"
	"gRPC-Cache/utils"
	"io"
	"log"
	"os"
	"sync"
	"time"

	pb "gRPC-Cache/description"

	"google.golang.org/grpc"
)

const (
	address = "localhost:9999"
	max     = 2
)

var (
	logger     = log.New(os.Stdout, fmt.Sprint(time.Now().Format(time.StampNano))+" : ", log.Lshortfile)
	ch         = make(chan string, 2)
	counter    = 0
	wgConsumer = sync.WaitGroup{}
)

func main() {
	var opts []grpc.DialOption

	opts = append(opts, grpc.WithInsecure())

	conn, err := grpc.Dial(address, opts...)
	utils.HandleError(err)
	defer conn.Close()

	client := pb.NewCacherClient(conn)

	wgConsumer.Add(max)
	for i := 0; i < max; i++ {
		go func(client pb.CacherClient, request *pb.Request) {
			defer wgConsumer.Done()
			logger.Println("Request started")
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			stream, err := client.GetRandomDataStream(ctx, request)
			utils.HandleError(err)
			for {
				reply, err := stream.Recv()
				if err == io.EOF {
					break
				}
				utils.HandleError(err)
				logger.Println(reply.Data[:len(reply.Data)/20]) // fixme никуда не выводить
			}
			counter++
			//wgConsumer.Done()
			ch <- "End: " + fmt.Sprint(counter)
		}(client, &pb.Request{N: int32(i)})
		logger.Println("End For:", i)
	}
	logger.Println("For is end")

	go printerRoutine()
	logger.Println("wgConsumer")
	wgConsumer.Wait()
	logger.Println("wgConsumer.Done")
	err = conn.Close()
	utils.HandleError(err)

}

func printerRoutine() {
	for w := range ch {
		logger.Println(w)
	}
}

// handleResponse lists all the features within the given bounding Rectangle.

//func printer() {
//	for {
//		logger.Println(<-ch)
//		time.Sleep(time.Second * 1)
//	}
//}
