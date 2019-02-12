package utils

import (
	"testing"
)

func TestGetRandom(t *testing.T) {
	min, max := 10, 100
	i := 0

	for {
		r := GetRandom(min, max)
		if r < min && r > max {
			t.Error("Expected random, got", r)
		} else if i == max {
			t.Log(i, r)
			break
		}
		i++
	}
	for {
		r := GetRandom(min, max)

		i++
		if r < min && r > max {
			t.Error("Expected random, got", r)
		} else if r == min || r == max {
			t.Log(i, r)
			break
		}
	}
}
