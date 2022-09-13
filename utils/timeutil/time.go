package timeutil

import (
	"fmt"
	"time"
)

const (
	MonthFormat      = "2006-01"
	DateFormat       = "2006-01-02"
	DateNumberFormat = "20060102"
	TimeFormat       = "2006-01-02 15:04:05"
	MinuteFormat     = "2006-01-02 15:04"
	TimeFormNoSplit  = "20060102150405"
	TimeZoneFormat   = "2006-01-02T15:04:05+08:00"
)

type Time time.Time

// 实现它的json序列化方法
func (t Time) MarshalJSON() ([]byte, error) {
	var stamp = fmt.Sprintf("\"%s\"", time.Time(t).Format(TimeFormat))
	return []byte(stamp), nil
}

func (t *Time) UnmarshalJSON(data []byte) error {
	now, err := time.ParseInLocation(`"`+TimeFormat+`"`, string(data), time.Local)
	*t = Time(now)
	return err
}
