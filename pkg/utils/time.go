package utils

import "time"

var TimeFormat = "2006-01-02 15:04:05"
var LocationShangHai, _ = time.LoadLocation("Asia/Shanghai")

func TimeJsonOut(t time.Time) string {
	return t.In(LocationShangHai).Format(TimeFormat)
}

func StrToTime(t string) time.Time {
	parse, e := time.ParseInLocation("2006-01-02 15:04:05", t, LocationShangHai)
	if e != nil {
		return time.Now()
	}
	return parse
}

func TodayStart() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, LocationShangHai)
}

func TodayEnd() time.Time  {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, LocationShangHai)
}
