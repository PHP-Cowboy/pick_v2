package timeutil

import "time"

func FormatToDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(DateFormat)
}

func FormatToDateTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(TimeFormat)
}

func FormatToDateTimeMinute(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(MinuteFormat)
}
