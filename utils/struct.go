package utils

import (
	"sync"
	"time"
)

var (
	instance *Config
	once     sync.Once
	//muConfig = new(sync.Mutex)
)

//URLs можно было контроллировать иным способом, чуть менее брутФорс,
// превратив URLs в структуру типа - Значение string:ЯвляетьсяЛиЗанятым bool.
// Тогда метод получения и освобождения ресурса будут более понятными.
// Но я уже реализовал, а тогда все нужно переписывать, на это нужно время, чуть позже реализую!
type Config struct {
	URLs             []string `yaml:"URLs"`
	MinTimeout       int      `yaml:"MinTimeout"`
	MaxTimeout       int      `yaml:"MaxTimeout"`
	NumberOfRequests int      `yaml:"NumberOfRequests"`

	ChReturnUrls chan string
	ChGetUrls    chan string
}

//func init() {
//	ChReturnUrls = make(chan string)
//	ChGetUrls = make(chan string)
//}

func (c *Config) takeURL() { //, name ...int) {
	//logger.Println("takeURL", "1")

	for {
		//muConfig.Lock()
		//logger.Println("takeURL", "2")
		//defer muConfig.Unlock()
		//logger.Println("takeURL", "3")
		var url string
		//logger.Println("takeURL", "4", c.LenURLs())
		for {
			//logger.Println("takeURL", "5")
			if len(c.URLs) > 0 {
				//logger.Println("takeURL", "6")
				r := GetRandom(0, len(c.URLs))
				//logger.Println("takeURL", "7")
				url = c.URLs[r]
				//logger.Println("takeURL", "8")
				c.URLs = append(c.URLs[0:r], c.URLs[r+1:]...)
				//logger.Println("takeURL", "9")
				break
			} else {
				time.Sleep(time.Second)
				logger.Println("takeURL Sleep")
			}
			//logger.Println("takeURL", "10")
		}
		//logger.Println("takeURL", "11")
		//debug.PrintStack()
		c.ChGetUrls <- url
		//logger.Println("takeURL", "12")
	}
}

// поднимаеться и работает как goroutine
//func (c *Config) ReturnURL(returnCh <-chan string) {
func (c *Config) returnURL() {
	//mu := sync.Mutex{}
	for {
		logger.Println("ReturnURL")
		//mu.Lock()
		//c.URLs = append(c.URLs, <-returnCh)
		c.URLs = append(c.URLs, <-instance.ChReturnUrls)
		//mu.Unlock()
		time.Sleep(500 * time.Millisecond)
	}
}

func GetConfig(configPath string) (config *Config) {
	once.Do(func() {
		configString := ReadFileConfig(configPath)
		instance = unmarshalConfig(configString)

		instance.ChReturnUrls = make(chan string)
		instance.ChGetUrls = make(chan string)
		go instance.takeURL()
		go instance.returnURL()
	})
	return instance
}

func (c *Config) LenURLs() (lenURLs int) {
	return len(c.URLs)
}
