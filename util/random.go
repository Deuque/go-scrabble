package util

import (
	"math/rand"
	"time"
)

func RandomInt(max int) int {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	return r.Intn(max)
}
