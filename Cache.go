package main

import (
	"fmt"
	"gRPC-Cache/resources"
	"gRPC-Cache/utils"
)

var config resources.Config

func main() {
	marshal := utils.ReadFileConfig("resources/config.yml")

	config = utils.UnmarshalConfig(marshal)

	fmt.Println("NumberOfRequests:", config.NumberOfRequests)

}
