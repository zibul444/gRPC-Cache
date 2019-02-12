package main

import (
	"fmt"
	"gRPC-Cache/cacher/resources"
	"gRPC-Cache/cacher/utils"
)

func main() {
	var config resources.Config
	config = utils.GetConfig("cacher/resources/config.yml")

	fmt.Println(config)
}
