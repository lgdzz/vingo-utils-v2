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

type UintIds []uint

func (s UintIds) Value() (driver.Value, error) {
	return strings.Join(SliceUintToString(s), ","), nil
}

func (s *UintIds) Scan(value interface{}) error {
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

func (s *UintIds) Uints() (result []uint) {
	CustomOutput(s, &result)
	return
}
