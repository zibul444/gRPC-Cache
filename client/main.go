package main

import (
	"context"
	"fmt"
	"gRPC-Cache/utils"
	"io"
	"log"
	"os"
	"time"

	pb "gRPC-Cache/description"

	"google.golang.org/grpc"
)

const (
	address = "localhost:9999"
	max     = 100
)

var (
	logger  = log.New(os.Stdout, fmt.Sprint(time.Now().Format(time.StampNano))+" : ", log.Lshortfile)
	ch      = make(chan string, max) //fixme ??????max
	counter = 0
)

func main() {
	var opts []grpc.DialOption

	opts = append(opts, grpc.WithInsecure())

	conn, err := grpc.Dial(address, opts...)
	utils.HandleError(err)
	defer conn.Close()

	client := pb.NewCacherClient(conn)

	for i := 0; i < max; i++ {
		func(client pb.CacherClient, request *pb.Request) {
			logger.Println("Request started")
			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel()
			stream, err := client.GetRandomDataStream(ctx, request)
			utils.HandleError(err)
			for {
				reply, err := stream.Recv()
				//_, err := stream.Recv()
				if err == io.EOF {
					break
				}
				utils.HandleError(err)
				logger.Println(reply.Data[:len(reply.Data)/20]) // fixme никуда не выводить
			}
			counter++

			ch <- "End" + fmt.Sprint(counter)
		}(client, &pb.Request{})
	}
	//handleResponse(client, &pb.Request{})
	//go printer()

	if counter == max {
		logger.Println("Client ended")
	}
}

// handleResponse lists all the features within the given bounding Rectangle.

//func printer() {
//	for {
//		logger.Println(<-ch)
//		time.Sleep(time.Second * 1)
//	}
//}
