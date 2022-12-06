package timeutil

import (
	"fmt"
	"math"
	"time"
)

func GetDateTime() string {
	return time.Now().Format(TimeFormat)
}

// 获取传入的时间所在月份的第一天，即某月第一天的0点。如传入time.Now(), 返回当前月份的第一天0点时间。
func GetFirstDateOfMonth(d time.Time) time.Time {
	d = d.AddDate(0, 0, -d.Day()+1)
	return GetZeroTime(d)
}

// 获取传入的时间所在月份的最后一天，即某月最后一天的0点。如传入time.Now(), 返回当前月份的最后一天0点时间。
func GetLastDateOfMonth(d time.Time) time.Time {
	return GetFirstDateOfMonth(d).AddDate(0, 1, -1)
}

// 获取某一天的0点时间
func GetZeroTime(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
}

// 获取某一天的最后一秒时间
func GetLastTime(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 23, 59, 59, 0, d.Location())
}

// 格式化为一天最小时间格式
func FormatTimeToMinDateTime(d time.Time, format string) string {
	d = GetZeroTime(d)
	return d.Format(format)
}

// 格式化为一天最大时间格式
func FormatTimeToMaxDateTime(d time.Time, format string) string {
	d = GetLastTime(d)
	return d.Format(format)
}

// 时间字符串转换为时间戳
// timeStr 时间字符串
// timeStringFormat 时间字符串的格式 不传递默认使用 2006-01-02 15:04:05
func StandardStr2Time(timeStr string, timeStringFormat ...string) int64 {
	timeFormat := TimeFormat
	if len(timeStringFormat) > 0 {
		timeFormat = timeStringFormat[0]
	}
	t, err := time.Parse(timeFormat, timeStr)
	if err != nil {
		return int64(0)
	}
	return t.Unix()
}

// 获取当前时间还有多少秒到下一天
func DayLeftSeconds() int64 {
	now := time.Now()
	return now.Unix() - GetLastTime(now).Unix()
}

// 日期格式的字符串转当天 00:00:00 的日期时间格式字符串
// 例如 2020-10-26 转 2020-10-26 00:00:00
// dateTimeSting 时间
// formFormat 时间格式
// toFormat 目标格式
func GetDateTimeStartByDate(dateTimeString string, formFormat, toFormat string) (string, error) {
	t, err := time.Parse(formFormat, dateTimeString)
	if err != nil {
		return "", err
	}
	return FormatTimeToMinDateTime(t, toFormat), nil
}

func GetDateTimeByDateTimeString(dateTimeString string, formFormat, toFormat string) (string, error) {
	t, err := time.Parse(formFormat, dateTimeString)
	if err != nil {
		return "", err
	}
	return t.Format(toFormat), nil
}

// 日期格式的字符串转当天 23:59:59 的日期时间格式字符串
// 例如 2020-10-26 转 2020-10-26 23:59:59
// dateTimeSting 时间
// formFormat 时间格式
// toFormat 目标格式
func GetDateTimeEndByDate(dateTimeString, formFormat, toFormat string) (string, error) {
	t, err := time.Parse(formFormat, dateTimeString)
	if err != nil {
		return "", err
	}
	return FormatTimeToMaxDateTime(t, toFormat), nil
}

// 计算两个日期之间差多少天
func DiffDays(startDateTime, endDateTime time.Time) int64 {
	return int64(math.Abs(endDateTime.Sub(startDateTime).Hours() / 24))
}

// 获取当前日期前后 days 天 months月 years年 日期
func GetTimeAroundByNum(days, months, years int) string {
	t := time.Now().AddDate(years, months, days)
	return GetZeroTime(t).Format(TimeFormat)
}

// 获取当前日期前后 days 天 日期
func GetTimeAroundByDays(days int) string {
	t := time.Now().AddDate(0, 0, days)
	return GetZeroTime(t).Format(TimeFormat)
}

// 获取当前日期前后 days 天 日期
func GetTimeAroundByMonths(months int) string {
	t := time.Now().AddDate(0, months, 0)
	return GetZeroTime(t).Format(TimeFormat)
}

func GetCurrentMonth() string {
	return time.Now().Format("01")
}

// 获取当前时间前后时间段时间
func GetDurationDateTime(s string) (string, error) {

	duration, err := time.ParseDuration(s) //"-30m"

	if err != nil {
		return "", err
	}

	dateTime := time.Now().Add(duration).Format("2006-01-02 15:04:05")

	return dateTime, nil
}

func GetDiffDateTime(t time.Time) string {

	if t.IsZero() {
		return ""
	}

	ms := int(time.Now().Sub(t).Minutes())

	diffStr := ""

	if ms >= 1440 { //超过一天
		day := ms / 1440 //天
		hour := (ms % 1440) / 60
		diffStr = fmt.Sprintf("%d天%d小时", day, hour)
	} else if ms > 60 { //超过一小时
		minute := ms % 60
		hour := ms / 60
		diffStr = fmt.Sprintf("%d小时%d分钟", hour, minute)
	} else {
		diffStr = fmt.Sprintf("%d分钟", ms)
	}

	return diffStr
}

// 这个方法有特殊用途，不建议修改或使用
// dateStr == "" 时的返回不是很合理，但符合使用场景的业务逻辑
func DateStrToTime(dateStr string) (*time.Time, error) {
	if dateStr == "" {
		return nil, nil
	}

	t, err := time.ParseInLocation(TimeFormat, dateStr, time.Local)

	if err != nil {
		return nil, err
	}

	return &t, nil
}
