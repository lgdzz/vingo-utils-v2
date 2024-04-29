package vingo

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/duke-git/lancet/v2/convertor"
	"math"
	"strings"
)

// JsonToString 结构体转字符串
func JsonToString(data any) string {
	output, err := json.Marshal(data)
	if err != nil {
		panic(err.Error())
	}
	return string(output)
}

// StringToJson 字符串转结构体
func StringToJson(data string, output any) {
	err := json.Unmarshal([]byte(data), &output)
	if err != nil {
		panic(err.Error())
	}
}

func MD5(str string) string {
	data := []byte(str) //切片
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has) //将[]byte转成16进制
	return md5str
}

func SHA256Hash(str string) string {
	hash := sha256.New()
	hash.Write([]byte(str))
	hashValue := hash.Sum(nil)
	return fmt.Sprintf("%x", hashValue)
}

// 自定义输出格式
func CustomOutput(input any, output any) {
	b, err := json.Marshal(input)
	if err != nil {
		panic(err.Error())
	}
	err = json.Unmarshal(b, output)
	if err != nil {
		panic(err.Error())
	}
}

// 转金额保留两位小数
func ToMoney(value float64) float64 {
	return ToDecimal(value)
}

// 浮点数保留两位小数
func ToDecimal(value float64) float64 {
	return math.Round(value*100) / 100
}

// 浮点数转百分比字符串
func ToPercentString(value float64) string {
	return fmt.Sprintf("%v%%", math.Round(value*100))
}

func ToInt(value any) int {
	return int(ToInt64(value))
}

func ToInt64(value any) int64 {
	v, _ := convertor.ToInt(value)
	return v
}

func ToUint(value any) uint {
	return uint(ToInt64(value))
}

func ToFloat(value any) float64 {
	v, _ := convertor.ToFloat(value)
	return v
}

func ToBool(value string) bool {
	v, _ := convertor.ToBool(value)
	return v
}

// 将值转换为字符串，对于数字、字符串、[]byte，将转换为字符串。 对于其他类型（切片、映射、数组、结构体）将调用 json.Marshal
func ToString(value any) string {
	return convertor.ToString(value)
}

func ToBase64(value any) string {
	return convertor.ToStdBase64(value)
}

func ToUrlBase64(value any) string {
	return convertor.ToUrlBase64(value)
}

func ToJson(value any) string {
	v, _ := convertor.ToJson(value)
	return v
}

// 将阿拉伯数字转换为中文数字
func NumberToChinese(number int) string {
	chineseDigits := []string{"零", "一", "二", "三", "四", "五", "六", "七", "八", "九"}
	chineseUnits := []string{"", "十", "百", "千", "万"}

	// 转换为字符串，以便逐个处理每一位数字
	numberStr := fmt.Sprint(number)

	var result strings.Builder

	for i, digit := range numberStr {
		digitInt := int(digit - '0')

		// 处理零
		if digitInt == 0 {
			// 如果上一位已经是零，则不重复添加
			if i > 0 && int(numberStr[i-1]-'0') == 0 {
				continue
			}

			result.WriteString(chineseDigits[digitInt])
		} else {
			// 处理非零数字
			result.WriteString(chineseDigits[digitInt])
			result.WriteString(chineseUnits[len(numberStr)-i-1])
		}
	}

	return result.String()
}
