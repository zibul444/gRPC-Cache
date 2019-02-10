package main

import (
	"fmt"
	"gRPC-Cache/resources"
	"gRPC-Cache/utils"
)

func main() {
	var config resources.Config
	config = utils.GetConfig("resources/config.yml")

	fmt.Println(config)
}
