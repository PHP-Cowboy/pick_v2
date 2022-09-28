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

//MyTime 自定义时间
type MyTime time.Time

func (t *MyTime) UnmarshalJSON(data []byte) error {
	now, err := time.ParseInLocation(`"`+TimeFormat+`"`, string(data), time.Local)
	*t = MyTime(now)
	return err
}

func (t MyTime) MarshalJSON() ([]byte, error) {
	var stamp = fmt.Sprintf("\"%s\"", time.Time(t).Format(TimeFormat))
	return []byte(stamp), nil
}

func (t MyTime) Value() (driver.Value, error) {
	// MyTime 转换成 time.Time 类型
	tTime := time.Time(t)
	return tTime.Format(TimeFormat), nil
}

func (t *MyTime) Scan(v interface{}) error {
	switch vt := v.(type) {
	case time.Time:
		// 字符串转成 time.Time 类型
		*t = MyTime(vt)
	default:
		return errors.New("类型处理错误")
	}
	return nil
}

func (t *MyTime) String() string {
	return time.Time(*t).Format(TimeFormat)
}
