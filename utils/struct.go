package utils

import (
	"sync"
	"time"
)

var (
	instance *Config
	once     sync.Once
)

//URLs можно было контроллировать иным способом, но уже сделал так...
type Config struct {
	URLs             []string `yaml:"URLs"`
	MinTimeout       int      `yaml:"MinTimeout"`
	MaxTimeout       int      `yaml:"MaxTimeout"`
	NumberOfRequests int      `yaml:"NumberOfRequests"`

	ChReturnUrls chan string
	ChGetUrls    chan string
}

// Выделяет доступные URL`s
func (c *Config) takeURL() {
	for {
		var (
			url string
			n   = time.Millisecond * 10
		)
		for {
			if len(c.URLs) > 0 {
				r := GetRandom(0, len(c.URLs))
				url = c.URLs[r]
				c.URLs = append(c.URLs[0:r], c.URLs[r+1:]...)
				break
			} else {
				time.Sleep(n)
				n = n + n
				//log.Println("takeURL Sleep", n)
			}
		}
		c.ChGetUrls <- url
	}
}

// Возвращает свободные
func (c *Config) returnURL() {
	for {

		c.URLs = append(c.URLs, <-instance.ChReturnUrls)
	}
}

// Получить объект конфига
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

// Кол-во доступных url-ов
func (c *Config) LenURLs() (lenURLs int) {
	return len(c.URLs)
}
