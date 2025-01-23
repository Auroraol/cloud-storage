package time

import (
	"strconv"
	"time"
)

const TimeLayout = "2006-01-02 15:04:05"

var loc *time.Location

func init() {
	cnLoc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		panic(err)
	}
	loc = cnLoc
}

func LocalTimeNow() time.Time {
	return time.Now().In(loc)
}

func StringToTime(stringTime string) time.Time {
	time, err := time.ParseInLocation(TimeLayout, stringTime, loc)
	if err == nil {
		return time
	}
	return LocalTimeNow()
}

func TimestampToTime(timestamp string) (time.Time, error) {
	i, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return LocalTimeNow(), err
	}
	tm := time.Unix(i, 0)
	return tm, nil
}

// 2024-09-01 00:00:00 转成 1725120000
func StringTimeToTimestamp(stringTime string) (int64, error) {
	timeInt64, err := time.ParseInLocation(TimeLayout, stringTime, loc)
	if err != nil {
		return 0, err
	}
	timeUnix := timeInt64.Unix()
	return timeUnix, nil
}

// 1725120000 转成 2024-09-01 00:00:00
func TimestampToStringTime(timestamp int64) string {
	return time.Unix(timestamp, 0).In(loc).Format(TimeLayout)
}
