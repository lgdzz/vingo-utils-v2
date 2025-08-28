// *****************************************************************************
// 作者: lgdz
// 创建时间: 2024/4/29
// 描述：
// *****************************************************************************

package vingo

import (
	"database/sql/driver"
	"encoding/json"
	"strings"
)

type IntIds []int

func (s IntIds) Value() (driver.Value, error) {
	return strings.Join(SliceIntToString(s), ","), nil
}

func (s *IntIds) Scan(value interface{}) error {
	var v string
	switch t := value.(type) {
	case string:
		v = t
	case []byte:
		v = string(t)
	default:
		panic("未知数据类型")
	}

	if v == "" {
		err := json.Unmarshal([]byte("[]"), s)
		if err != nil {
			panic(err)
		}
	} else {
		CustomOutput(SliceStringToUint(strings.Split(v, ",")), s)
	}
	return nil
}

func (s *IntIds) Ints() (result []int) {
	CustomOutput(s, &result)
	return
}
