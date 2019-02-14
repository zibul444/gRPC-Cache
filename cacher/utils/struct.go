package utils

import (
	"fmt"
)

type Config struct {
	URLs             []string `yaml:"URLs"`
	MinTimeout       int      `yaml:"MinTimeout"`
	MaxTimeout       int      `yaml:"MaxTimeout"`
	NumberOfRequests int      `yaml:"NumberOfRequests"`
}

func (c *Config) TakeURL(chanel chan<- string) { //, name ...int) {
	mutex.Lock()
	defer mutex.Unlock()
	var url string
	for {
		if len(c.URLs) > 0 {
			r := GetRandom(0, len(c.URLs))
			url = c.URLs[r]
			c.URLs = append(c.URLs[0:r], c.URLs[r+1:]...)
			break
		}
	}
	//if len(name) > 0 {
	//	fmt.Println("TakeURL job is:", name[0])
	//chanel <- fmt.Sprint(name[0], ": ", url)
	//} else {
	chanel <- fmt.Sprint(url)
	//}
}

// поднимаеться и работает как goroutine
func (c *Config) ReturnURL(returnCh <-chan string) {
	for {
		mutex.Lock()
		c.URLs = append(c.URLs, <-returnCh) // fixme Может быть ошибка при одновременном извлечении данных  методом TakeURL()
		mutex.Unlock()
	}
}

func GetConfig(configPath string) (config *Config) {
	configString := ReadFileConfig(configPath)
	return UnmarshalConfig(configString)
}
