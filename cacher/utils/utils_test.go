package utils

import (
	"testing"
)

func TestGetRandom(t *testing.T) {
	//min, max := 0, 1
	min, max := 10, 100
	i := 0

	for {
		i++
		r := GetRandom(min, max)
		if r < min {
			t.Fatal("Expected random, got", r)
		} else if r > max {
			t.Fatal("Expected random, got", r)
		} else if r == max-1 {
			t.Log(i, r)
			break
		}
	}
	for {
		i++
		r := GetRandom(min, max)
		if r < min {
			t.Fatal("Expected random, got", r)
		} else if r > max {
			t.Fatal("Expected random, got", r)
		} else if r == min {
			t.Log(i, r)
			break
		}
	}
}

// TODO
func TestExecuteCommand(t *testing.T) {
	//defer Conn.Close()

	//fmt.Printf("--- %v\n", utils.ExecuteCommand("EXPIRE", "test:string", 100)) // время жизни значения
	t.Logf("--- %s\n", ExecuteCommand("PING"))
	t.Logf("--- %v\n", ExecuteCommand("EXPIRE", "foo", 100)) // время жизни значения

	t.Logf("--- %d\n", ExecuteCommand("APPEND", "foo", " v"))
	//t.Logf("--- %s\n", ExecuteCommand("GET", "foo"))
	t.Logf("--- %s\n", Execute("GET", "foo"))

	t.Logf("--- %s\n", ExecuteCommand("SET", "foo", "c"))
	t.Logf("--- %s\n", Execute("GET", "foo"))
	t.Logf("--- %v\n", ExecuteCommand("TTL", "foo"))
}
