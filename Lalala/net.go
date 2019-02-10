package main

import (
	"fmt"
	"gRPC-Cache/utils"

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

	conf := utils.GetConfig("resources/config.yml")
	url = utils.GetRandomUrl(conf)
	keys := utils.Execute("KEYS", "*")

	// fixme говнище подчищаю чтоли....
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
		if err != nil {
			//fmt.Fprintf(os.Stderr, "fetch: %v\n", err)
			panic(err)
		}
		//s := fmt.Sprint(resp)
		//fmt.Println(resp)
		//fmt.Println(s[:100])

		timeLife = utils.GetRandomTimeLife(conf)
		utils.ExecuteCommand("SETEX", url, timeLife, resp)
		utils.Execute("GET", url)

		url = utils.GetRandomUrl(conf)
		resp, err = http.Get(url)

		utils.Execute("KEYS", "*")

		time.Sleep(100 * time.Millisecond)
	}

	//for { //в вечном цикле собираем данные
	//	x, err := goquery.ParseUrl(utils.GetRandomUrl(conf))
	//	if err == nil {
	//	fmt.Printf("\nx: %h", x)
	//	if s := strings.TrimSpace(x.Find(".fi_text").Text()); s != "" {
	//	//c <- s //и отправляем их в канал
	//	fmt.Println(s) //и отправляем их в канал
	//}

}
