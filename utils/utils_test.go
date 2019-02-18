package utils

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"sync"
	"testing"
)

func TestGetRandom(t *testing.T) {
	//min, max := 0, 1
	min, max := 10, 100
	i := 0

	for {
		i++
		r := GetRandom(min, max)
		if r < min {
			t.Fatal("Expected random, got", r)
		} else if r > max {
			t.Fatal("Expected random, got", r)
		} else if r == max-1 {
			t.Log(i, r)
			break
		}
	}
	for {
		i++
		r := GetRandom(min, max)
		if r < min {
			t.Fatal("Expected random, got", r)
		} else if r > max {
			t.Fatal("Expected random, got", r)
		} else if r == min {
			t.Log(i, r)
			break
		}
	}
}

func TestExecuteCommand(t *testing.T) {

	//defer dbConn.Close()

	//fmt.Printf("--- %v\n", utils.ExecuteCommand("EXPIRE", "test:string", 100))
	t.Logf("--- %s\n", ExecuteCommand("PING"))
	t.Logf("--- %v\n", ExecuteCommand("EXPIRE", "foo", 100))

	t.Logf("--- %d\n", ExecuteCommand("APPEND", "foo", " v"))
	//t.Logf("--- %s\n", ExecuteCommand("GET", "foo"))
	t.Logf("--- %s\n", Execute("GET", "foo"))

	t.Logf("--- %s\n", ExecuteCommand("SETEX", "foo", 10, "c"))
	t.Logf("--- %s\n", Execute("GET", "foo"))
	TTL := ExecuteCommand("TTL", "foo")
	//i := interface{}(TTL)
	ty := reflect.TypeOf(TTL).Name()
	t.Logf("--- TTL %v, %v\n", TTL, ty)
}

func TestExecuteCommand2(t *testing.T) {
	URLs := []string{
		"https://www.microsoft.com",
		"https://juliadates.com",
		"https://www.nato.int",
		"http://discovery-romance.com",
		"http://gitarre.ru",
		"https://www.ed.gov",
		"https://partner.edarling.ru",
		"https://minjust.ru",
		"https://partner.edarling.ru",
		"https://www.rosminzdrav.ru",
		"http://www.mkrf.ru",
		"https://www.mos.ru",
		"http://pereborom.ru",
		"http://www.calorizator.ru",
		"https://black-star.ru",
	}

	transCfg := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // disable verify
	}
	// Create Http Client
	client := &http.Client{Transport: transCfg}

	var wgTest sync.WaitGroup
	for i := 0; i < 20; i++ {
		wgTest.Add(1)
		go func() {
			for _, url := range URLs {
				//if time.Now().Second()%2 == 0 {
				// Request
				response, err := client.Get(url)
				// Check Error
				HandleError(err)
				// Close After Read Body
				// Read Body
				data, err := ioutil.ReadAll(response.Body)
				// Check Error
				HandleError(err)
				// Print response html : conver byte to string
				//fmt.Println(string(data))

				//data, err := http.Get(url)
				//HandleError(err)
				ExecuteCommand("setex", url, 50, fmt.Sprint(data))
				response.Body.Close()
				//}
			}
			wgTest.Done()
		}()
	}

	wgTest.Wait()

	for i, url := range URLs {
		logger.Println(i, ExecuteCommand("TTL", url), url)
	}

	keys := Execute("KEYS", "*")
	logger.Println("Contains keys :", keys)
}

func TestExecute2(t *testing.T) {
	//var URL = make([]string, 9)

	keys := Execute("KEYS", "*")
	keys = keys[1 : len(keys)-1]

	keysSlays := strings.Split(keys, " ")

	logger.Println("Kount", len(keysSlays))
	logger.Println("keysSlays:", keysSlays)

}

func TestExecute(t *testing.T) {
	test := "TestLiter"
	t.Logf("--- %s\n", ExecuteCommand("SETEX", "lalala", 10, test))
	dest := Execute("GET", "lalala")

	if !strings.EqualFold(test, dest) {
		t.Fatal("dest:", dest)
	}

}

func TestHTTPGet(t *testing.T) {
	for {
		recp, err := http.Get("http://localhost:8000/l")
		HandleError(err)

		logger.Printf("%v", recp)
	}
}
