package util

import (
	"math/rand"
)

var chars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func GetRandom(l int) string {
	s := make([]rune, l)
	for i := 0; i < l; i++ {
		s[i] = chars[rand.Intn(len(chars))]
	}
	return string(s)
}
