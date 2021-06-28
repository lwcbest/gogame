package util

import (
	"math/rand"
	"time"
)

func RandomSome(max int) int {
	rand.Seed(time.Now().UnixNano())
	// 隝机数
	r := rand.Intn(max)
	return r
}
