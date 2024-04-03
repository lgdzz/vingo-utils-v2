package vingo

import (
	"database/sql/driver"
	"encoding/json"
	"strings"
	"time"
)

// 定位坐标
type Location struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

// 单个字段修改请求体
type PatchBody struct {
	Field string `json:"field"`
	Value any    `json:"value"`
}

// 文件信息
type FileInfo struct {
	Name      string `json:"name"`
	Type      string `json:"type"`
	Mimetype  string `json:"mimetype"`
	Extension string `json:"extension"`
	Size      int64  `json:"size"`
	Realpath  string `json:"realpath"`
}

// 文件信息(简单)
type FileInfoSimple struct {
	Name     string `json:"name"`
	Realpath string `json:"realpath"`
}

// 身份证信息
// IdCardInfo{IdCard: ""}
type IdCardInfo struct {
	IdCard     string // 身份证号码
	RegionCode string // 6位行政区域编码
	Birthday   string // 2006-01-02 格式日期
	Age        int    // 年龄：精确到月份
	UniformAge int    // 年龄：按年份计算
	Gender     string // 性别
}

// 身份验证异常
type AuthException struct {
	Message string
}

// 数据库事务异常
type DbException struct {
	Message string
}

// 时间范围
type DateRange struct {
	Start time.Time
	End   time.Time
}

// 字符串时间范围
type DateRangeString struct {
	Start string
	End   string
}

func (s *DateRange) OfString() DateRangeString {
	return DateRangeString{
		Start: s.Start.Format(DatetimeFormat),
		End:   s.End.Format(DatetimeFormat),
	}
}

type DateAt [2]string

func (s *DateAt) Start() string {
	return s[0]
}

func (s *DateAt) End() string {
	return s[1]
}

type UintIds []uint

func (s UintIds) Value() (driver.Value, error) {
	return strings.Join(SliceUintToString(s), ","), nil
}

func (s *UintIds) Scan(value interface{}) error {
	v := string(value.([]byte))
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

type StringSlice []string

func (s StringSlice) Value() (driver.Value, error) {
	return strings.Join(s, ","), nil
}

func (s *StringSlice) Scan(value interface{}) error {
	v := string(value.([]byte))
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
