// *****************************************************************************
// 作者: lgdz
// 创建时间: 2025/6/13
// 描述：时间格式转换
// *****************************************************************************

package moment

import (
	"github.com/lgdzz/vingo-utils-v2/vingo"
	"strings"
	"time"
)

// 时间字符串转time.Time
func TimeTextToTime(timeText string, layouts ...string) time.Time {
	for _, layout := range layouts {
		if t, err := time.ParseInLocation(layout, timeText, time.Local); err == nil {
			return t
		}
	}
	panic("invalid date format")
}

func (s DateText) ToString() string {
	return strings.TrimSpace(string(s))
}

// 日期字符串转time.Time
func (s DateText) ToTime() time.Time {
	layouts := []string{vingo.DateFormat, "2006-1-2"}
	return TimeTextToTime(s.ToString(), layouts...)
}

// 日期字符串转vingo.LocalTime
func (s DateText) ToLocalTime() vingo.LocalTime {
	return vingo.NewLocalTime(s.ToTime())
}

func (s DateTimeText) ToString() string {
	return strings.TrimSpace(string(s))
}

// 日期时间字符串转time.Time
func (s DateTimeText) ToTime() time.Time {
	layouts := []string{vingo.DatetimeFormat, "2006-1-2 15:4:5"}
	return TimeTextToTime(s.ToString(), layouts...)
}

// 日期时间字符串转vingo.LocalTime
func (s DateTimeText) ToLocalTime() vingo.LocalTime {
	return vingo.NewLocalTime(s.ToTime())
}
