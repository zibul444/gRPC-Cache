package main

import (
	"fmt"
	"gRPC-Cache/resources"
	"gRPC-Cache/utils"
)

func main() {
	var config resources.Config
	marshal := utils.ReadFileConfig("resources/config.yml")
	config = utils.UnmarshalConfig(marshal)

	fmt.Println(config)
}
