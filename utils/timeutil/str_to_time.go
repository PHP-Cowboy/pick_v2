package timeutil

import (
	"time"
)

// 将字符串时间转换成当天的结束时间
// 时间格式为 DateFormat
func StrToLastTime(timeStr string) (err error, t time.Time) {
	t, err = time.ParseInLocation(DateFormat, timeStr, time.Local)

	if err != nil {
		return
	}

	t = GetLastTime(t)

	return
}
