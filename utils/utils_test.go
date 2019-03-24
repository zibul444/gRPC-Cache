package utils

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strings"
	"sync"
	"testing"
)

func TestGetRandom(t *testing.T) {
	min, max := 10, 100
	i := 0

	for {
		i++
		r := GetRandom(min, max)
		if r < min && r > max {
			t.Fatal("Expected random, got", r)
		} else if r == max-1 {
			t.Log(i, r)
			break
		}
	}
	for {
		i++
		r := GetRandom(min, max)
		if r < min && r > max {
			t.Fatal("Expected random, got", r)
		} else if r == min {
			t.Log(i, r)
			break
		}
	}
}

func TestExecuteCommand(t *testing.T) {

	t.Logf("--- %s\n", ExecuteCommand("PING"))
	t.Logf("--- %v\n", ExecuteCommand("EXPIRE", "foo", 100))

	t.Logf("--- %d\n", ExecuteCommand("APPEND", "foo", " v"))
	t.Logf("--- %s\n", Execute("GET", "foo"))

	t.Logf("--- %s\n", ExecuteCommand("SETEX", "foo", 10, "c"))
	t.Logf("--- %s\n", Execute("GET", "foo"))
	TTL := ExecuteCommand("TTL", "foo")
	ty := reflect.TypeOf(TTL).Name()
	t.Logf("--- TTL %v, %v\n", TTL, ty)
}

func TestExecute2(t *testing.T) {
	keys := Execute("KEYS", "*")
	keys = keys[1 : len(keys)-1]

	keysSlays := strings.Split(keys, " ")

	logger.Notice("Kount", len(keysSlays))
	logger.Notice("keysSlays:", keysSlays)

	KEYS1 := Execute("KEYS", config.URLs[0])
	t.Log("keys:", KEYS1)
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
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: transCfg}

	var wgTest sync.WaitGroup
	for i := 0; i < 20; i++ {
		wgTest.Add(1)
		go func() {
			for _, url := range URLs {
				response, err := client.Get(url)
				HandleError(err)
				data, err := ioutil.ReadAll(response.Body)
				HandleError(err)

				ExecuteCommand("setex", url, 50, fmt.Sprint(data))
				err = response.Body.Close()
				HandleError(err)
			}
			wgTest.Done()
		}()
	}

	wgTest.Wait()
	for i, url := range URLs {
		logger.Notice(i, ExecuteCommand("TTL", url), url)
	}

	keys := Execute("KEYS", "*")
	logger.Notice("Contains keys :", keys)
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
	resp, err := http.Get("https://www.microsoft.com")
	HandleError(err)

	log.Printf("%v", resp)
}
