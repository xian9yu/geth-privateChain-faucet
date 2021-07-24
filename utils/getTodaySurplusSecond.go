package utils

import "time"

// GetTodaySurplusSecond 获取现在到晚上0点的时间戳秒
func GetTodaySurplusSecond() int64 {
	t := "2006-01-02"
	loc, _ := time.LoadLocation("Asia/Shanghai")
	ti, _ := time.ParseInLocation(t, time.Now().Format(t), loc)
	return 86400 - (time.Now().Unix() - ti.Unix())
}
