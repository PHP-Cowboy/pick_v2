package timeutil

import "time"

func FormatToDate(t time.Time) (date string) {
	if !t.IsZero() {
		date = t.Format(DateFormat)
	}
	return
}

func FormatToDateTime(t time.Time) (dateTime string) {
	if !t.IsZero() {
		dateTime = t.Format(TimeFormat)
	}
	return
}

func FormatToDateTimeMinute(t time.Time) (minuteTime string) {
	if !t.IsZero() {
		minuteTime = t.Format(MinuteFormat)
	}
	return
}
