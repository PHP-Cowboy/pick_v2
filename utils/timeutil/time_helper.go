package timeutil

import (
	"math"
	"time"
)

func GetDateTime() string {
	return time.Now().Format(TimeFormat)
}

//获取传入的时间所在月份的第一天，即某月第一天的0点。如传入time.Now(), 返回当前月份的第一天0点时间。
func GetFirstDateOfMonth(d time.Time) time.Time {
	d = d.AddDate(0, 0, -d.Day()+1)
	return GetZeroTime(d)
}

//获取传入的时间所在月份的最后一天，即某月最后一天的0点。如传入time.Now(), 返回当前月份的最后一天0点时间。
func GetLastDateOfMonth(d time.Time) time.Time {
	return GetFirstDateOfMonth(d).AddDate(0, 1, -1)
}

//获取某一天的0点时间
func GetZeroTime(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
}

//获取某一天的最后一秒时间
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

//获取当前时间还有多少秒到下一天
func DayLeftSeconds() int64 {
	now := time.Now()
	return now.Unix() - GetLastTime(now).Unix()
}

//日期格式的字符串转当天 00:00:00 的日期时间格式字符串
//例如 2020-10-26 转 2020-10-26 00:00:00
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

//日期格式的字符串转当天 23:59:59 的日期时间格式字符串
//例如 2020-10-26 转 2020-10-26 23:59:59
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

//计算两个日期之间差多少天
func DiffDays(startDateTime, endDateTime time.Time) int64 {
	return int64(math.Abs(endDateTime.Sub(startDateTime).Hours() / 24))
}
