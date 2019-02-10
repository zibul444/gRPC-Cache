package main

import (
	"fmt"
	"gRPC-Cache/utils"
)

func main() {
	//utils.C = utils.NewPool().Get()
	defer utils.Conn.Close()

	//fmt.Printf("--- %v\n", utils.ExecuteCommand("EXPIRE", "test:string", 100)) // время жизни значения
	fmt.Printf("--- %s\n", utils.ExecuteCommand("PING"))
	fmt.Printf("--- %v\n", utils.ExecuteCommand("EXPIRE", "foo", 100)) // время жизни значения

	fmt.Printf("--- %d\n", utils.ExecuteCommand("APPEND", "foo", " v"))
	//fmt.Printf("--- %s\n", utils.ExecuteCommand("GET", "foo"))
	fmt.Printf("--- %s\n", utils.Execute("GET", "foo"))

	fmt.Printf("--- %s\n", utils.ExecuteCommand("SET", "foo", "c"))
	fmt.Printf("--- %s\n", utils.Execute("GET", "foo"))
	fmt.Printf("--- %v\n", utils.ExecuteCommand("TTL", "foo"))

}
