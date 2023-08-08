package rand

import (
	"math/rand"
	"time"
)

const (
	charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	numSet  = "0123456789"
)

func init() {
	rand.Seed(time.Now().UnixMilli())
}

// New 生成指定字符数的随机字符串
func New(bits int) (randString string) {
	for i := 0; i < bits; i++ {
		randString += string(charset[rand.Intn(len(charset))])
	}
	return
}

// Default 生成16位随机字符串
func Default() string { return New(16) }

// Sms 生成6位随机短信验证码
func Sms() (randString string) {
	for i := 0; i < 6; i++ {
		randString += string(numSet[rand.Intn(len(numSet))])
	}
	return
}
