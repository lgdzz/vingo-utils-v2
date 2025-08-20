package vingo

import (
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

func (s *DateAt) StartTime() time.Time {
	t, err := time.ParseInLocation(DatetimeFormat, s.Start(), time.Local)
	if err != nil {
		panic(err.Error())
	}
	return t
}

func (s *DateAt) EndTime() time.Time {
	t, err := time.ParseInLocation(DatetimeFormat, s.End(), time.Local)
	if err != nil {
		panic(err.Error())
	}
	return t
}

// bool切片字符串形态，如：true,false
type BoolString string

func (s *BoolString) ToSlice() []bool {
	result := make([]bool, 0)
	arr := strings.Split(string(*s), ",")
	for _, item := range arr {
		result = append(result, ToBool(item))
	}
	return result
}

// 数字切片字符串形态，如：1,2,3
type IntString string

func (s *IntString) ToSlice() []int {
	result := make([]int, 0)
	arr := strings.Split(string(*s), ",")
	for _, item := range arr {
		result = append(result, ToInt(item))
	}
	return result
}

func (s *IntString) ToUintSlice() []uint {
	result := make([]uint, 0)
	arr := strings.Split(string(*s), ",")
	for _, item := range arr {
		result = append(result, ToUint(item))
	}
	return result
}

// 文本切片字符串形态，如：a,b,c
type TextString string

func (s *TextString) ToSlice() []string {
	return strings.Split(string(*s), ",")
}

type Ids[T any] struct {
	Ids []T `json:"ids"`
}

// 范围字符串形态，如：100,500
type BetweenText string
type BetweenStruct struct {
	Start float64
	End   float64
}

func (s *BetweenText) ToStruct() BetweenStruct {
	arr := strings.Split(string(*s), ",")
	if len(arr) != 2 {
		panic("范围字符串格式错误")
	}
	return BetweenStruct{
		Start: ToFloat(arr[0]),
		End:   ToFloat(arr[1]),
	}
}

// 排序
type Sort[T any] Ids[T]

type KeyValue struct {
	Key   string
	Value string
}

type IdBody struct {
	Id any `json:"id"`
}

type DeleteBody struct {
	IdBody
}

type DetailQuery struct {
	Id    uint       `form:"id"`
	Fetch TextString `form:"fetch"`
}
