package main

import (
	"fmt"
	"gRPC-Cache/cacher/utils"
	"gopkg.in/yaml.v2"
	"time"
)

var data2 = `
a: Easy!
b:
  c: 2
  d: [3, 4]
m: [ "2001-01-01T15:04:05Z", "2002-02-02T15:04:05Z" ]

`

type T1 struct {
	A string
	B struct {
		RenamedC int   `yaml:"c"`
		D        []int `yaml:",flow"`
	}
	M []time.Time
}

func main() {
	t := T1{}

	err := yaml.Unmarshal([]byte(data2), &t)
	utils.HandleError(err)
	fmt.Printf("--- t:\n%v\n\n", t)
}
