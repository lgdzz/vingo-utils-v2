package vingo

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

/**
 * 检查字符串中是否包含某些字符串(or)
 * @param string $haystack 如：Hello World
 * @param string|[]string $needles 如：Hello 或 []string{"Hello", "World"}
 * @return bool
 */
func StringContainsOr(haystack string, needles interface{}) bool {
	for _, n := range toSlice(needles) {
		if n != "" && strings.Contains(haystack, n) {
			return true
		}
	}
	return false
}

/**
 * 检查字符串中是否包含某些字符串(and)
 * @param string $haystack 如：Hello World
 * @param string|[]string $needles 如：Hello 或 []string{"Hello", "World"}
 * @return bool
 */
func StringContainsAnd(haystack string, needles interface{}) bool {
	for _, n := range toSlice(needles) {
		if n != "" && !strings.Contains(haystack, n) {
			return false
		}
	}
	return true
}

/**
 * 检查字符串是否以某些字符串开头
 * @param string $haystack 如：Hello World
 * @param string|[]string $needles 如：Hello 或 []string{"Hello", "World"}
 * @return bool
 */
func StringStartsWith(haystack string, needles interface{}) bool {
	for _, n := range toSlice(needles) {
		if strings.HasPrefix(haystack, n) {
			return true
		}
	}
	return false
}

/**
 * 检查字符串是否以某些字符串结尾
 * @param string $haystack 如：Hello World
 * @param string|[]string $needles 如：Hello 或 []string{"Hello", "World"}
 * @return bool
 */
func StringEndsWith(haystack string, needles interface{}) bool {
	for _, n := range toSlice(needles) {
		if strings.HasSuffix(haystack, n) {
			return true
		}
	}
	return false
}

// 将传入的参数转换成 []string 类型的切片
func toSlice(v interface{}) []string {
	switch v := v.(type) {
	case string:
		return []string{v}
	case []string:
		return v
	default:
		return nil
	}
}

// 截取字符串
func StringSubstr(s string, start int, length int) string {
	runes := []rune(s)
	l := start + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[start:l])
}

// 解析身份证信息
func IdCard(id string) IdCardInfo {
	if len(id) != 18 {
		panic("身份证号长度不正确")
	}
	now := time.Now()
	info := IdCardInfo{IdCard: id, Gender: Male, RegionCode: id[:6]}
	year := id[6:10]
	month := id[10:12]
	day := id[12:14]
	info.Birthday = fmt.Sprintf("%v-%v-%v", year, month, day)
	info.UniformAge = now.Year() - int(ToUint(year))
	if ToUint(time.Now().Format("01")) < ToUint(month) {
		info.Age = info.UniformAge - 1
	} else {
		info.Age = info.UniformAge
	}
	if i, _ := strconv.Atoi(string(id[16])); i%2 == 0 {
		info.Gender = Female
	}
	return info
}

// 验证身份证号是否正确（校验码验证）
// 身份证号码的最后一位校验码是根据前面的17位数字计算出来的。计算步骤如下：
// 1．将身份证号码的前17位数字分别乘以对应的系数。系数从左到右依次是：{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}
// 2．将上述乘积求和
// 3. 将求和结果除以11，取余数
// 4. 根据余数，从对应表中找出校验码：{0: "1", 1: "0", 2: "X", 3: "9", 4: "8", 5: "7", 6: "6", 7: "5", 8: "4", 9: "3", 10: "2"}
func IdCardCheck(id string) bool {
	if len(id) != 18 {
		return false
	}
	strSlice := strings.Split(id, "")
	var sum uint
	for i, factor := range idCardFactors {
		sum += ToUint(strSlice[i]) * factor
	}
	// 余数
	code := sum % 11
	if last, ok := idCardCodes[code]; ok {
		// 余数取到的效验码等于身份证最后一位则为有效
		return last == strings.ToUpper(strSlice[17])
	}
	return false
}

/**
 * 将字节转换为可读文本
 * @param int size 大小
 * @param int precision 保留小数位数
 * @return string
 */
func FormatBytes(size int64, precision int) string {
	units := []string{"B", "KB", "MB", "GB", "TB", "PB"}
	i := 0
	fsize := big.NewFloat(float64(size))
	for fsize.Cmp(big.NewFloat(1024)) >= 0 && i < 6 {
		fsize.Quo(fsize, big.NewFloat(1024))
		i++
	}
	format := fmt.Sprintf("%%.%df %%s", precision)
	return fmt.Sprintf(format, fsize, units[i])
}

// 手机号加星
func MaskMobilePhone(mobile string) string {
	if len(mobile) != 11 {
		panic("手机号必须为11位")
	}
	return mobile[:3] + "****" + mobile[7:]
}

// 名称加星
func MaskName(name string) string {
	length := utf8.RuneCountInString(name)
	firstRune, _ := utf8.DecodeRuneInString(name)
	lastRune, _ := utf8.DecodeLastRuneInString(name)
	if length == 2 {
		return string(firstRune) + strings.Repeat("*", length-1)
	} else {
		return string(firstRune) + strings.Repeat("*", length-2) + string(lastRune)
	}
}
