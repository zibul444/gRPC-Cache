package utils

import (
	"log"
	"regexp"
	"sync"
	"testing"
	"time"
)

var (
	config  = GetConfig("../resources/config.yml")
	ch      = make(chan string)
	lengthy = len(config.URLs)
	//wgTest  = sync.WaitGroup{}
	buf []string
)

func TestGetConfig(t *testing.T) {
	if len(config.URLs) <= 0 {
		t.Fatal()
	} else if config.MinTimeout < 0 {
		t.Fatal()
	} else if config.MinTimeout > config.MaxTimeout {
		t.Fatal()
	} else if config.MaxTimeout == 0 {
		t.Fatal()
	} else if config.NumberOfRequests <= 0 {
		t.Fatal()
	}
	re := regexp.MustCompile(`https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{2,256}\.[a-z]{2,6}\b([-a-zA-Z0-9@:%_\+.~#?&//=]*)`)
	for _, url := range config.URLs {
		if re.FindAllStringSubmatch(url, -1) == nil {
			t.Fatal()
		}
	}
}

func TestConfig_TakeURL(t *testing.T) {
	log.Println("Start TestConfig_TakeURL")
	n := 5
	wgTest := sync.WaitGroup{}

	mu := new(sync.Mutex)

	for i := 0; i < n; i++ {
		wgTest.Add(1)
		go func(number int) {
			log.Println(number, "go func start")
			for len(config.URLs) > 0 {
				mu.Lock()
				go config.TakeURL(ch) //, number)

				time.Sleep(time.Millisecond * 50) // даем поработать go TakeURL
				mu.Unlock()
			}
			wgTest.Done()
		}(i)
	}

	log.Println("go func is run")
	go func() {
		log.Println("Starting reading chan")
		for url := range ch {
			log.Println("url", url)
			buf = append(buf, url)
		}
	}()

	log.Println("wgTest.Wait()")
	wgTest.Wait()
	log.Println("wgTest.Wait() - ended")
	if len(config.URLs) != 0 {
		t.Fatal("len", len(config.URLs))
	}
	if len(buf) != lengthy {
		t.Fatal(len(config.URLs))
	}
	log.Println("Source:", config.URLs)
	log.Println("Dest::", buf)
}

//func TestConfig_ReturnURL(t *testing.T) {
//	/*for i := 0; ; {
//
//	}*/
//	t.Skip()
//}
