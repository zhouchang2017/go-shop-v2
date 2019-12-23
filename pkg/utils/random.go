package utils

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func RandomInt(max int64) int64 {
	return rand.Int63n(max)
}