package utils

import (
	"fmt"
	"log"
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
	log.Println("Start TestConfig_TakeURL")
	n := length
	var wgTest sync.WaitGroup
	x := 0

	for i := 0; i < n; i++ {
		wgTest.Add(1)
		go func(number int) {
			x++
			log.Println("count:", x)
			var url string
			url = <-ChGetUrls
			buf = append(buf, url)

			time.Sleep(50 * time.Millisecond) // даем поработать go takeURL
			wgTest.Done()
		}(i)
	}

	log.Println("wgTest.Wait()")
	wgTest.Wait()
	log.Println("wgTest.Wait() - ended")
	if len(config.URLs) != 0 {
		t.Fatal("len", len(config.URLs))
	}
	if len(buf) != length {
		t.Fatal("URLs:", len(config.URLs), "length:", length, "buf:", len(buf))
	}
	log.Println("Source:", config.URLs)
	log.Println("Dest:", buf)
}

func TestCheckAvailabilityResources(t *testing.T) {
	var data string
	for _, resource := range config.URLs {
		resp, err := http.Get(resource)
		HandleError(err)
		data = fmt.Sprint(resp)
		log.Println(resource, " - ", len(data), " - ", data)
	}
}
