package vingo

import (
	"fmt"
	"github.com/google/uuid"
	"math/rand"
	"strings"
	"time"
)

// 生成UUID
func GetUUID() string {
	return uuid.NewString()
}

// 生成随机字符串
func RandomString(length int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// 生成随机数
func RandomNumber(length int) string {
	digits := []rune("0123456789")
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, length)
	for i := range b {
		b[i] = digits[rand.Intn(len(digits))]
	}
	return string(b)
}

// 生成按时间+随机数的单号
func OrderNo(length int, check func(string) bool) string {
	if length <= 14 {
		panic("编号长度不少于15位")
	}
	orderNo := fmt.Sprintf("%v%v", time.Now().Format("20060102150405"), RandomNumber(length-14))
	if check != nil && check(orderNo) {
		// 已存在，重新生成
		return OrderNo(length, check)
	}
	return strings.ToUpper(orderNo)
}
