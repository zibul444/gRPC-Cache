package utils

import (
	"fmt"
	"gRPC-Cache/resources"
	"github.com/garyburd/redigo/redis"
	"gopkg.in/yaml.v2"
	"io"
	"log"
	"math/rand"
	"os"
	"time"
)

// fixme Необходимо проинициализировать до вызова(подумать о синглетоне на менднях)
var C redis.Conn

func init() {
	C = NewPool().Get()
}

func GetConfig(configName string) (config resources.Config) {
	configString := ReadFileConfig(configName)
	return UnmarshalConfig(configString)
}

// читаем файл по имени
func ReadFileConfig(fileName string) (result string) {
	f, err := os.Open(fileName)
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
	defer f.Close()
	buf := make([]byte, 64)
	for {
		n, err := f.Read(buf)
		if err == io.EOF { // если конец файла
			break // выходим из цикла
		}
		//fmt.Print(string(data[:n]))
		result += string(buf[:n])
	}
	fmt.Println("End Reading file config")
	//fmt.Printf("Data from the file: \n%v", result)
	return
}

// Десериализуем конфигурационный файл в объект resources.Config
func UnmarshalConfig(marshal string) (config resources.Config) {
	err := yaml.Unmarshal([]byte(marshal), &config)
	if err != nil {
		log.Printf("error: %v", err)
	}
	//fmt.Printf("\nconfig:\n%v\n", config)
	fmt.Println("End Unmarshaling file config")
	return
}

// Для выполнения любых команд не возвращающих строковое значение
func ExecuteCommand(commandName string, args ...interface{}) interface{} {
	result, err := C.Do(commandName, args...)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%s:%v\n", commandName, result)
	return result
}

// Для выполнения команд возвращающих строку
func Execute(commandName string, args ...interface{}) string {
	result, err := C.Do(commandName, args...)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%s:%s\n", commandName, result)
	return fmt.Sprintf("%s", result)
}

func NewPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:   80,
		MaxActive: 12000, // max number of connections
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", ":6379")
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}
}

// Рандомное значение от min до max включительно.
func GetRandom(min int, max int) int {
	max++
	rand.Seed(time.Now().UnixNano())
	//maxp := flag.Int("max", max-min, "the max value")

	n := rand.Intn(max - min)
	n += min
	return n
}

// Получаем случайный ресурс из config
func GetRandomUrl(config resources.Config) string {
	min, max := 0, len(config.URLs)
	max--
	r := GetRandom(min, max)

	return config.URLs[r]
}

// Получаем случайное время жизни впределах указанных в config
func GetRandomTimeLife(config resources.Config) int {
	min, max := config.MinTimeout, config.MaxTimeout

	return GetRandom(min, max)
}
