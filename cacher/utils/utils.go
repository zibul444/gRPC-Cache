package utils

import (
	"fmt"
	"gRPC-Cache/cacher/resources"
	"github.com/garyburd/redigo/redis"
	"gopkg.in/yaml.v2"
	"io"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

// fixme Необходимо проинициализировать до вызова(подумать о синглетоне на менднях)
var Conn redis.Conn

func init() {
	Conn = NewPool().Get()
}

func GetConfig(configPath string) (config resources.Config) {
	configString := ReadFileConfig(configPath)
	return UnmarshalConfig(configString)
}

// Читаем файл по имени
func ReadFileConfig(filePath string) (fileContents string) {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatalln(err.Error())
		os.Exit(1)
	}
	defer HandleError(f.Close())
	buf := make([]byte, 64)
	for {
		n, err := f.Read(buf)
		if err == io.EOF { // если конец файла
			break // выходим из цикла
		}
		//log.Print(string(data[:n]))
		fileContents += string(buf[:n])
	}
	log.Println("End Reading file config")
	//log.Printf("Data from the file: \n%v", result)
	return
}

// Десериализуем конфигурационный файл в объект resources.Config
func UnmarshalConfig(marshal string) (config resources.Config) {
	err := yaml.Unmarshal([]byte(marshal), &config)
	HandleError(err)
	log.Println("End Unmarshaling file config")
	return
}

// Для выполнения любых команд
func ExecuteCommand(commandName string, args ...interface{}) interface{} {
	result, err := Conn.Do(commandName, args...)
	HandleError(err)
	//log.Printf("%s:%v\n", commandName, result)
	return result
}

// Для выполнения команд "get string"
func Execute(commandName string, args ...interface{}) string {
	result, err := Conn.Do(commandName, args...)
	HandleError(err)
	//log.Printf("%s:%s\n", commandName, result)
	//return result.(string)
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

// Рандомное значение от min до max включительно.
func GetRandom(min int, max int) (r int) {
	max++
	rand.Seed(time.Now().UnixNano())

	r = rand.Intn(max - min)
	r += min
	return r
}

// Получаем случайный ресурс из config
func GetRandomUrl(config resources.Config) (url string) {
	min, max := 0, len(config.URLs)
	max-- // что бы не выйти за пределы доступных ресурсов конфига
	r := GetRandom(min, max)

	return config.URLs[r]
}

// Получаем случайное время жизни впределах указанных в config
func GetRandomTimeLife(config resources.Config) (timeLife int) {
	min, max := config.MinTimeout, config.MaxTimeout
	GetListReadyKeys()
	return GetRandom(min, max)
}

// Deprecated Получить список доступных ключей
func GetListReadyKeys() (urls []string) {
	keysRAW := Execute("KEYS", "*")
	// fixme подумать как корректно преобразовывать из БД
	//  (или разобраться почему база выдает не в том формате)
	strings.Replace(keysRAW, "[", "", 1)
	strings.Replace(keysRAW, "]", "", 1)
	urls = strings.Split(keysRAW, " ")
	return urls
}
