package utils

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/garyburd/redigo/redis"
	"gopkg.in/yaml.v2"
)

var (
	Conn   = NewPool().Get()
	logger = log.New(os.Stdout, fmt.Sprint(time.Now().Format(time.StampNano))+": ", log.Lshortfile)
	//urls   []string
	mutex = new(sync.Mutex)
)

// Читаем файл по имени
func ReadFileConfig(filePath string) (fileContents string) {
	f, err := os.Open(filePath)
	if err != nil {
		logger.Fatalln(err.Error())
		os.Exit(1)
	}
	defer f.Close()
	buf := make([]byte, 32)
	for {
		n, err := f.Read(buf)
		if err == io.EOF { // если конец файла
			break // выходим из цикла
		}
		HandleError(err)
		fileContents += string(buf[:n])
	}

	logger.Println("Чтение конфига завершено")
	return
}

// Десериализуем конфигурационный файл в объект resources.Config
func UnmarshalConfig(marshal string) (config *Config) {
	err := yaml.Unmarshal([]byte(marshal), &config)
	HandleError(err)
	logger.Println("Конфиг жив(анмаршалинг завершен)")
	return
}

// Для выполнения любых команд
func ExecuteCommand(commandName string, args ...interface{}) interface{} {
	//mutex.Lock()
	//defer mutex.Unlock()
	result, err := Conn.Do(commandName, args...)
	HandleError(err)
	//logger.Printf("%s:%v\n", commandName, result)
	return result
}

// Для выполнения команд "get string"
func Execute(commandName string, args ...interface{}) string {
	//mutex.Lock()
	//defer mutex.Unlock()
	result, err := Conn.Do(commandName, args...)
	HandleError(err)
	return fmt.Sprintf("%s", result)
}

func HandleError(err error) {
	if err != nil {
		panic(err)
	}
}

func NewPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:   1000,
		MaxActive: 10000,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", ":6379")
			HandleError(err)
			return c, err
		},
	}
}

// Рандомное значение от min до max.
func GetRandom(min int, max int) (randInt int) {
	rand.Seed(time.Now().UnixNano())

	randInt = rand.Intn(max - min)
	randInt += min
	return
}

// Получаем случайное время жизни впределах указанных в config
func GetRandomTimeLife(config Config) (timeLife int) {
	min, max := config.MinTimeout, config.MaxTimeout
	return GetRandom(min, max)
}

//func GetLogger() (stdLogger *log.Logger) {
//	stdLogger = log.New(os.Stdout, fmt.Sprint(time.Now().Format(time.StampNano)) + " INFO: ", log.Lshortfile)
//	return stdLogger
//}
