package utils

import (
	"regexp"
	"sync"
	"testing"
	"time"
)

var (
	config = GetConfig("../config.yml")
	ch     = make(chan string)
	length = config.LenURLs()
	//wgTest  = sync.WaitGroup{}
	buf []string
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
	re := regexp.MustCompile(`https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{2,256}\.[a-z]{2,6}\b([-a-zA-Z0-9@:%_\+.~#?&//=]*)`)
	for _, url := range config.URLs {
		if re.FindAllStringSubmatch(url, -1) == nil {
			t.Fatal()
		}
	}
}

func TestConfig_TakeURL(t *testing.T) {
	logger.Println("Start TestConfig_TakeURL")
	n := 5
	wgTest := sync.WaitGroup{}

	mu := new(sync.Mutex)

	for i := 0; i < n; i++ {
		wgTest.Add(1)
		go func(number int) {
			logger.Println(number, "go func start")
			for len(config.URLs) > 0 {
				mu.Lock()
				go config.TakeURL(ch) //, number)

				time.Sleep(time.Millisecond * 50) // даем поработать go TakeURL
				mu.Unlock()
			}
			wgTest.Done()
		}(i)
	}

	logger.Println("go func is run")
	go func() {
		logger.Println("Starting reading chan")
		for url := range ch {
			logger.Println("url", url)
			buf = append(buf, url)
		}
	}()

	logger.Println("wgTest.Wait()")
	wgTest.Wait()
	logger.Println("wgTest.Wait() - ended")
	if len(config.URLs) != 0 {
		t.Fatal("len", len(config.URLs))
	}
	if len(buf) != length {
		t.Fatal(len(config.URLs))
	}
	logger.Println("Source:", config.URLs)
	logger.Println("Dest:", buf)
}
