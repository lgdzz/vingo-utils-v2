// *****************************************************************************
// 作者: lgdz
// 创建时间: 2024/4/29
// 描述：
// *****************************************************************************

package vingo

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type JsonObject[T any] struct {
	Data T
}

func (s JsonObject[T]) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *JsonObject[T]) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("Scan source is not []byte")
	}
	return json.Unmarshal(b, s)
}
