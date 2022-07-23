package timeutil

import "time"

const (
	Date     = "2006-01-02"
	DateTime = "2006-01-02 15:04:05"
)

func GetDateTime() string {
	return time.Now().Format(DateTime)
}
