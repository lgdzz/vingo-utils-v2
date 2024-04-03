package vingo

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
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

func ToUint(v any) uint {
	switch t := v.(type) {
	case uint:
		return t
	case int32:
		return uint(t)
	case int64:
		return uint(t)
	case uint32:
		return uint(t)
	case float32:
		return uint(t)
	case float64:
		return uint(t)
	case string:
		v, _ := strconv.Atoi(t)
		return uint(v)
	default:
		panic(fmt.Sprintf("Cannot convert to uint: %v", v))
	}
	return 0
}

func ToString(v any) string {
	switch value := v.(type) {
	case int:
		return strconv.Itoa(value)
	case int8:
		return strconv.FormatInt(int64(value), 10)
	case int16:
		return strconv.FormatInt(int64(value), 10)
	case int32:
		return strconv.FormatInt(int64(value), 10)
	case int64:
		return strconv.FormatInt(value, 10)
	case uint:
		return strconv.FormatUint(uint64(value), 10)
	case uint8:
		return strconv.FormatUint(uint64(value), 10)
	case uint16:
		return strconv.FormatUint(uint64(value), 10)
	case uint32:
		return strconv.FormatUint(uint64(value), 10)
	case uint64:
		return strconv.FormatUint(value, 10)
	case float32:
		return strconv.FormatFloat(float64(value), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(value, 'f', -1, 64)
	}
	return ""
}

func ToFloat64(value interface{}) float64 {
	switch v := value.(type) {
	case string:
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			panic(fmt.Sprintf("字符串：%v转换float64失败，错误：%v", v, err.Error()))
		}
		return f
	case int:
		return float64(v)
	case int8:
		return float64(v)
	case int16:
		return float64(v)
	case int32:
		return float64(v)
	case int64:
		return float64(v)
	case uint:
		return float64(v)
	case uint8:
		return float64(v)
	case uint16:
		return float64(v)
	case uint32:
		return float64(v)
	case uint64:
		return float64(v)
	case float32:
		return float64(v)
	case float64:
		return v
	default:
		return 0
	}
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
