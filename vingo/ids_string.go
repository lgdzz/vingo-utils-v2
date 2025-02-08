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

type StringSlice []string

func (s StringSlice) Value() (driver.Value, error) {
	return strings.Join(s, ","), nil
}

func (s *StringSlice) Scan(value interface{}) error {
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
		*s = strings.Split(v, ",")
	}
	return nil
}

func (s *StringSlice) Strings() (result []string) {
	CustomOutput(s, &result)
	return
}
