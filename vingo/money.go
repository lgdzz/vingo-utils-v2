// *****************************************************************************
// 作者: lgdz
// 创建时间: 2024/4/17
// 描述：金额转大写
// *****************************************************************************

package vingo

import "strings"

var cnNums = [...]string{"零", "壹", "贰", "叁", "肆", "伍", "陆", "柒", "捌", "玖"}

var cnIntUnits = [...]string{"", "拾", "佰", "仟"}

var cnIntRadice = [...]string{"", "万", "亿", "兆"}

func MoneyToChinese(number string) string {
	// 将金额字符串分割成整数部分和小数部分
	parts := strings.Split(number, ".")
	integerPart := parts[0]
	decimalPart := ""
	if len(parts) > 1 {
		decimalPart = parts[1]
	}

	// 转换整数部分
	result := ""
	zero := false
	needZero := false
	for i := 0; i < len(integerPart); i++ {
		pos := len(integerPart) - i - 1
		num := integerPart[i] - '0'
		if num == 0 {
			zero = true
			needZero = true
		} else {
			if zero {
				result += cnNums[0]
				zero = false
			}
			if needZero {
				result += cnNums[0]
				needZero = false
			}
			result += cnNums[num] + cnIntUnits[pos%4]
			if pos%4 == 0 {
				result += cnIntRadice[pos/4]
			}
		}
	}

	// 转换小数部分
	if decimalPart != "" {
		result += "点"
		for i := 0; i < len(decimalPart); i++ {
			num := decimalPart[i] - '0'
			result += cnNums[num]
		}
	}

	return result + "元"
}
