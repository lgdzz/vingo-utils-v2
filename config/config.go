// *****************************************************************************
// 作者: lgdz
// 创建时间: 2024/3/15
// 描述：
// *****************************************************************************

package config

type Config struct {
}

func (s *Config) StringValue(value *string, defaultValue string) {
	if *value == "" {
		*value = defaultValue
	}
}

func (s *Config) IntValue(value *int, defaultValue int) {
	if *value == 0 {
		*value = defaultValue
	}
}
