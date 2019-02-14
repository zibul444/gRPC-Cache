package main

import (
	"fmt"

	"gRPC-Cache/cacher/utils"

	"net/http"
	"strings"
	"time"
)

var (
	timeLife int
	url      string
)

func main() {
	//c := make(chan string)

	//conf := resources.GetConfig("cacher/resources/config.yml")
	//url = utils.GetRandomUrl(conf)
	keys := utils.Execute("KEYS", "*")

	strings.Replace(keys, "[", "", 1)
	strings.Replace(keys, "]", "", 1)
	keysSlise := strings.Split(keys, " ")

	if len(keysSlise) > 0 {
		for _, k := range keysSlise {
			fmt.Println(k)
		}
	}
	resp, err := http.Get(url)
	defer resp.Body.Close()

	for {
		utils.HandleError(err)

		timeLife = utils.GetRandomTimeLife(conf)
		utils.ExecuteCommand("SETEX", url, timeLife, resp)
		utils.Execute("GET", url)

		//url = utils.GetRandomUrl(conf)
		resp, err = http.Get(url)

		utils.Execute("KEYS", "*")

		time.Sleep(100 * time.Millisecond)
	}

}
