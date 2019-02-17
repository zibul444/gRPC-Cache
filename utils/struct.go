package utils

import (
	"sync"
	"time"
)

var (
	instance *Config
	once     sync.Once
	muConfig = new(sync.Mutex)
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

	chReturnUrls chan string
	chGetUrls    chan string
}

//func init() {
//	chReturnUrls = make(chan string)
//	chGetUrls = make(chan string)
//}

func (c *Config) TakeURL(chanel chan<- string) { //, name ...int) {
	logger.Println("TakeURL", "1")

	muConfig.Lock()
	logger.Println("TakeURL", "2")
	defer muConfig.Unlock()
	//logger.Println("TakeURL", "3")
	var url string
	//logger.Println("TakeURL", "4", c.LenURLs())
	for {
		//logger.Println("TakeURL", "5")
		if len(c.URLs) > 0 {
			//logger.Println("TakeURL", "6")
			r := GetRandom(0, len(c.URLs))
			//logger.Println("TakeURL", "7")
			url = c.URLs[r]
			//logger.Println("TakeURL", "8")
			c.URLs = append(c.URLs[0:r], c.URLs[r+1:]...)
			logger.Println("TakeURL", "9")
			break
		} else {
			time.Sleep(time.Second)
			logger.Println("TakeURL Sleep")
		}
		//logger.Println("TakeURL", "10")
	}
	logger.Println("TakeURL", "11")
	//debug.PrintStack()
	chanel <- url
	logger.Println("TakeURL", "12")
}

// поднимаеться и работает как goroutine
func (c *Config) ReturnURL(returnCh <-chan string) {
	//mu := sync.Mutex{}
	for {
		logger.Println("ReturnURL")
		//mu.Lock()
		c.URLs = append(c.URLs, <-returnCh)
		//mu.Unlock()
		time.Sleep(time.Second)
	}
}

func GetConfig(configPath string) (config *Config) {
	once.Do(func() {
		configString := ReadFileConfig(configPath)
		instance = unmarshalConfig(configString)
	})
	return instance
}

func (c *Config) LenURLs() (lenURLs int) {
	return len(c.URLs)
}
