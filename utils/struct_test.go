package utils

import (
	"fmt"
	"net/http"
	"sync"
	"testing"
	"time"
)

var (
	config       = GetConfig("../config.yml")
	ChGetUrls    = config.ChGetUrls
	ChReturnUrls = config.ChReturnUrls
	length       = config.LenURLs()
	buf          = make([]string, 0)
)

func TestGetConfig(t *testing.T) {
	if len(config.URLs) <= 0 {
		t.Fatal(len(config.URLs))
	} else if config.MinTimeout < 0 {
		t.Fatal()
	} else if config.MinTimeout > config.MaxTimeout {
		t.Fatal()
	} else if config.MaxTimeout == 0 {
		t.Fatal()
	} else if config.NumberOfRequests <= 0 {
		t.Fatal()
	}
	//FIXME
	//re := regexp.MustCompile(`https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{2,256}\.[a-z]{2,6}\b([-a-zA-Z0-9@:%_\+.~#?&//=]*)`)
	//for _, url := range config.URLs {
	//	if re.FindAllStringSubmatch(url, -1) == nil {
	//		t.Fatal()
	//	}
	//}

	instance2 := GetConfig("../config.yml")

	if config != instance2 {
		t.Fatal("Objects are not equal!\n")
	}

}

func TestConfig_TakeURL(t *testing.T) {
	logger.Println("Start TestConfig_TakeURL")
	n := length
	var wgTest sync.WaitGroup
	x := 0
	//mu := new(sync.Mutex)

	for i := 0; i < n; i++ {
		wgTest.Add(1) // fixme может нужно будет поднять
		go func(number int) {
			//logger.Println(number, "go func starting")
			//for len(config.URLs) > 0 {
			x++
			logger.Println("count:", x)
			//mu.Lock()
			//go config.takeURL(ch) //, number)
			var url string
			url = <-ChGetUrls
			//logger.Println("url", url)
			//ChReturnUrls <- url
			buf = append(buf, url)

			time.Sleep(50 * time.Millisecond) // даем поработать go takeURL
			//mu.Unlock()
			//}
			wgTest.Done()
		}(i)
	}

	logger.Println("wgTest.Wait()")
	wgTest.Wait()
	logger.Println("wgTest.Wait() - ended")
	if len(config.URLs) != 0 {
		t.Fatal("len", len(config.URLs))
	}
	if len(buf) != length {
		//logger.Println(buf)
		t.Fatal("URLs:", len(config.URLs), "length:", length, "buf:", len(buf))
	}
	logger.Println("Source:", config.URLs)
	logger.Println("Dest:", buf)
}

func TestCheckAvailabilityResources(t *testing.T) {
	var data string
	for _, resource := range config.URLs {
		resp, err := http.Get(resource)
		HandleError(err)
		data = fmt.Sprint(resp)
		logger.Println(resource, " - ", len(data), " - ", data) //data[:len(data)/15])
	}
}

//
//func TestBenchmarcBD(t *testing.T) {
//	t.Skip()
//}
