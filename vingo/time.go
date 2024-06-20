package vingo

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"github.com/duke-git/lancet/v2/datetime"
	"time"
)

type LocalTime time.Time

func NewLocalTime(t time.Time) (l LocalTime) {
	l.To(t)
	return
}

func (t LocalTime) MarshalJSON() ([]byte, error) {
	tTime := time.Time(t).Local()
	return []byte(fmt.Sprintf("\"%v\"", tTime.Format(DatetimeFormat))), nil
}

func (t *LocalTime) UnmarshalJSON(data []byte) error {
	var err error
	var parsedTime time.Time
	if string(data) == "null" {
		*t = LocalTime{}
		return nil
	}

	parsedTime, err = time.ParseInLocation(`"`+DatetimeFormat+`"`, string(data), time.Local)
	if err != nil {
		return err
	}

	*t = LocalTime(parsedTime.Local())
	return nil
}

func (t LocalTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	tlt := time.Time(t).Local()
	//判断给定时间是否和默认零时间的时间戳相同
	if tlt.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return tlt, nil
}

func (t *LocalTime) Scan(v interface{}) error {
	if value, ok := v.(time.Time); ok {
		*t = LocalTime(value.Local())
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

func (t LocalTime) Now() LocalTime {
	return LocalTime(time.Now().Local())
}

func (t *LocalTime) SetNow() {
	*t = LocalTime(time.Now().Local())
}

func (t *LocalTime) To(value time.Time) {
	*t = LocalTime(value.Local())
}

func (t LocalTime) String() string {
	return time.Time(t).Format(DatetimeFormat)
}

func (t LocalTime) Time() time.Time {
	return time.Time(t)
}

func (t LocalTime) Format(layout string) string {
	return t.Time().Format(layout)
}

func (t *LocalTime) ScanFromRow(rows *sql.Rows, columnName string) error {
	var tmp time.Time
	err := rows.Scan(&tmp)
	if err != nil {
		return err
	}
	*t = LocalTime(tmp.Local())
	return nil
}

func (t LocalTime) ValueFromRow(rows *sql.Rows, columnName string) (interface{}, error) {
	return t, nil
}

func TimeAddDays(t time.Time, days int, hour int, min int, sec int) time.Time {
	// Add the specified number of days
	t = t.AddDate(0, 0, days)

	// Set the time to midnight
	year, month, day := t.Date()

	midnight := time.Date(year, month, day, hour, min, sec, 0, t.Location())

	return midnight
}

// 判断当前时间是否大于指定的时间
func TimeIsAfterNow(t time.Time) bool {
	// 获取当前时间
	now := time.Now()

	// 判断当前时间是否大于指定时间
	return now.After(t)
}

// 获取昨日开始时间
func GetYesterdayStartTime() time.Time {
	now := time.Now().Local()
	yesterday := now.AddDate(0, 0, -1)
	return time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, now.Location())
}

// 是否超过时间
func IsTimeExceeded(t time.Time, days int) bool {
	duration := time.Since(t)
	return duration > time.Duration(days)*24*time.Hour
}

// 是否是未来日期
func IsAfterDate(date string) bool {
	targetDate, _ := time.ParseInLocation(DateFormat, date, time.Local)
	return targetDate.After(time.Now())
}

// 是否是过去日期
func IsBeforeDate(date string) bool {
	targetDate, _ := time.ParseInLocation(DateFormat, date, time.Local)
	return targetDate.Before(time.Now())
}

// 是否是未来时间
func IsAfterDatetime(datetime string) bool {
	targetDate, _ := time.ParseInLocation(DatetimeFormat, datetime, time.Local)
	return targetDate.After(time.Now())
}

// 是否是过去时间
func IsBeforeDatetime(datetime string) bool {
	targetDate, _ := time.ParseInLocation(DatetimeFormat, datetime, time.Local)
	return targetDate.Before(time.Now())
}

// 获取最近一周日期
func GetLastWeekDates(startDay ...time.Time) []string {
	var t time.Time
	if len(startDay) == 0 {
		t = time.Now()
	} else {
		t = startDay[0]
	}
	t = t.AddDate(0, 0, -6)

	var lastWeek []string
	for i := 0; i < 7; i++ {
		day := t.AddDate(0, 0, i).Format("2006-01-02")
		lastWeek = append(lastWeek, day)
	}
	return lastWeek
}

// 获取最近一个月日期
func GetLastMonthDates(startDay ...time.Time) []string {
	var t time.Time
	if len(startDay) == 0 {
		t = time.Now()
	} else {
		t = startDay[0]
	}
	t = t.AddDate(0, 0, -29)

	var lastWeek []string
	for i := 0; i < 30; i++ {
		day := t.AddDate(0, 0, i).Format("2006-01-02")
		lastWeek = append(lastWeek, day)
	}
	return lastWeek
}

// 获取未来一个月的日期
func GetNextMonthDates(startDay ...time.Time) []string {
	var t time.Time
	if len(startDay) == 0 {
		t = time.Now()
	} else {
		t = startDay[0]
	}
	t = t.AddDate(0, 0, 1) // 将初始时间设置为明天

	var nextMonth []string
	for i := 0; i < 30; i++ {
		day := t.AddDate(0, 0, i).Format("2006-01-02")
		nextMonth = append(nextMonth, day)
	}
	return nextMonth
}

// 获取当前时间值指针
func GetNowTime() *LocalTime {
	t := LocalTime{}.Now()
	return &t
}

// 获取当前时间值
func GetNowTimeValue() LocalTime {
	return LocalTime{}.Now()
}

// 本月时间范围
func GetThisMonthRange() (r DateRange) {
	now := time.Now()            // 获取当前时间
	year, month, _ := now.Date() // 获取当前年份、月份
	r.Start = time.Date(year, month, 1, 0, 0, 0, 0, now.Location())
	r.End = time.Date(year, month, r.Start.AddDate(0, 1, -1).Day(), 23, 59, 59, 0, now.Location())
	return
}

// 本季时间范围
func GetThisQuarterRange() (r DateRange) {
	now := time.Now()
	quarter := (now.Month() - 1) / 3
	r.Start = time.Date(now.Year(), time.Month(quarter*3+1), 1, 0, 0, 0, 0, now.Location())
	r.End = r.Start.AddDate(0, 3, -1).Add(time.Hour * 23).Add(time.Minute * 59).Add(time.Second * 59) // 设置到本季度最后一天的最后一秒钟
	return
}

// 今年时间范围
func GetThisYearRange() (r DateRange) {
	now := time.Now()  // 获取当前时间
	year := now.Year() // 获取当前年份
	r.Start = time.Date(year, 1, 1, 0, 0, 0, 0, now.Location())
	r.End = time.Date(year, 12, 31, 23, 59, 59, 0, now.Location())
	return
}

// 获取日期范围
func GetDateDayRange(date string) (r DateRange) {
	var err error

	switch date {
	case "yesterday":
		r.Start = GetYesterdayStartTime()
	case "today":
		now := time.Now()
		r.Start = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	default:
		r.Start, err = time.ParseInLocation(DateFormat, date, time.Local)
		if err != nil {
			panic(err.Error())
		}
	}
	r.End = r.Start.Add((DaySec - 1) * time.Second)

	return
}

// 获取指定月日期范围
func GetMonthRange(monthStr string) (r DateRange) {
	t, err := time.ParseInLocation("2006-01", monthStr, time.Local)
	if err != nil {
		panic(fmt.Sprintf("日期解析失败：%v", err.Error()))
	}
	r.Start = time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
	r.End = time.Date(t.Year(), t.Month(), 1, 23, 59, 59, 0, t.Location()).AddDate(0, 1, -1)
	return
}

// 获取昨日开始时间和结束时间
func GetLastDayBetween() DateRange {
	between := DateRange{}
	between.Start = datetime.BeginOfDay(time.Now().AddDate(0, 0, -1))
	between.End = datetime.EndOfDay(between.Start)
	return between
}

// 获取上周开始时间和结束时间
func GetLastWeekBetween() DateRange {
	now := time.Now()
	beforeDay := int(now.Weekday())
	if beforeDay == 0 {
		beforeDay = 7
	}
	now = now.AddDate(0, 0, -beforeDay)
	between := DateRange{}
	between.Start = datetime.BeginOfWeek(now, time.Monday)
	between.End = datetime.EndOfWeek(between.Start, time.Sunday)
	return between
}

// 获取上月开始时间和结束时间
func GetLastMonthBetween() DateRange {
	between := DateRange{}
	between.Start = datetime.BeginOfMonth(time.Now().AddDate(0, -1, 0))
	between.End = datetime.EndOfMonth(between.Start)
	return between
}

// 获取去年开始时间和结束时间
func GetLastYearBetween() DateRange {
	between := DateRange{}
	between.Start = datetime.BeginOfYear(time.Now().AddDate(-1, 0, 0))
	between.End = datetime.EndOfYear(between.Start)
	return between
}

// 获取时间范围内的所有月份数据，格式：YYYY-MM
func GenerateMonths(startDate, endDate string) ([]string, error) {
	// 将传入的日期字符串解析为时间对象
	startTime, err := time.ParseInLocation("2006-01-02", startDate, time.Local)
	if err != nil {
		return nil, err
	}

	endTime, err := time.ParseInLocation("2006-01-02", endDate, time.Local)
	if err != nil {
		return nil, err
	}

	// 存储生成的月份字符串的切片
	months := []string{}

	// 循环生成月份
	currentMonth := startTime
	for currentMonth.Before(endTime) || currentMonth.Equal(endTime) {
		months = append(months, currentMonth.Format("2006-01"))
		currentMonth = currentMonth.AddDate(0, 1, 0)
	}

	return months, nil
}

// 获取时间范围内的所有日期数据，格式：YYYY-MM-DD
func GenerateDatesOfTime(dateRange DateRange) []string {
	dates := []string{}
	currentDate := dateRange.Start
	for !currentDate.After(dateRange.End) {
		dates = append(dates, currentDate.Format(DateFormat))
		currentDate = currentDate.AddDate(0, 0, 1)
	}
	return dates
}

func IsToday(t time.Time) bool {
	now := time.Now()
	return t.Year() == now.Year() && t.YearDay() == now.YearDay()
}

func GetTodayOfLocalTime() LocalTime {
	t := time.Now()
	r := LocalTime{}
	r.To(time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local))
	return r
}

// 获取指定年开始时间和结束时间
func GetYearBetween(year string) DateRange {
	t, err := time.ParseInLocation("2006", year, time.Local)
	if err != nil {
		panic(err.Error())
	}
	between := DateRange{}
	between.Start = datetime.BeginOfYear(t)
	between.End = datetime.EndOfYear(t)
	return between
}
