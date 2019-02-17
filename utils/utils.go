package utils

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"sync"

	//"github.com/garyburd/redigo/redis"

	//"github.com/garyburd/redigo/redis"
	"io"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"gopkg.in/yaml.v2"
)

var (
	Conn    = NewPool().Get()
	logger  = log.New(os.Stdout, fmt.Sprint(time.Now().Format(time.StampNano))+": ", log.Lshortfile)
	muUtils = sync.Mutex{}
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
func unmarshalConfig(marshal string) (config *Config) {
	err := yaml.Unmarshal([]byte(marshal), &config)
	HandleError(err)
	logger.Println("Конфиг жив(анмаршалинг завершен)")
	return
}

// Для выполнения любых команд
func ExecuteCommand(commandName string, args ...interface{}) interface{} {
	muUtils.Lock()
	defer muUtils.Unlock()
	//logger.Printf("1 %s:%v\n", commandName, args)
	result, err := Conn.Do(commandName, args...)
	//logger.Printf("2 %s:%v\n", commandName, result)

	//if err == io.ErrShortWrite { //fixme
	//	Conn = NewPool().Get()
	//	logger.Println("Error", err)
	//	err = nil
	//	ExecuteCommand(commandName, args)
	//}
	HandleError(err)
	return result
}

// Для выполнения команд "get string"
func Execute(commandName string, args ...interface{}) string {
	muUtils.Lock()
	defer muUtils.Unlock()
	result, err := Conn.Do(commandName, args...)
	//logger.Println(fmt.Sprintf("Execute: %s", result))

	//if err == io.ErrShortWrite { //fixme
	//	Conn = NewPool().Get()
	//	logger.Println("Error", err)
	//	err = nil
	//	Execute(commandName, args)
	//}
	HandleError(err)
	return fmt.Sprintf("%s", result)
}

func HandleError(err error) {
	if err != nil {
		//logger.Print(err)
		panic(err)
		//debug.PrintStack()

		//wgTest.Done()
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

func ToInt64(i1 interface{}) int64 {
	if i1 == nil {
		return 0
	}
	switch i2 := i1.(type) {
	default:
		i3, _ := strconv.ParseInt(ToString(i2), 10, 64)
		return i3
	case *json.Number:
		i3, _ := i2.Int64()
		return i3
	case json.Number:
		i3, _ := i2.Int64()
		return i3
	case int64:
		return i2
	case float64:
		return int64(i2)
	case float32:
		return int64(i2)
	case uint64:
		return int64(i2)
	case int:
		return int64(i2)
	case uint:
		return int64(i2)
	case bool:
		if i2 {
			return 1
		} else {
			return 0
		}
	case *bool:
		if i2 == nil {
			return 0
		}
		if *i2 {
			return 1
		} else {
			return 0
		}
	}
}

func ToString(i1 interface{}) string {
	if i1 == nil {
		return ""
	}
	switch i2 := i1.(type) {
	default:
		return fmt.Sprint(i2)
	case bool:
		if i2 {
			return "true"
		} else {
			return "false"
		}
	case string:
		return i2
	case *bool:
		if i2 == nil {
			return ""
		}
		if *i2 {
			return "true"
		} else {
			return "false"
		}
	case *string:
		if i2 == nil {
			return ""
		}
		return *i2
	case *json.Number:
		return i2.String()
	case json.Number:
		return i2.String()
	}
}
