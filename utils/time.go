package utils

import (
	"github.com/hero1s/gotools/db"
	"math"
	"time"
)

//时间转换的模板，golang里面只能是 "2006-01-02 15:04:05" （go的诞生时间）
const TimeTemplateSecond = "2006-01-02 15:04:05" //常规类型
const TimeTemplate2 = "2006/01/02 15:04:05"      //其他类型
const TimeTemplateDay = "2006-01-02"             //其他类型
const TimeTemplate4 = "15:04:05"                 //其他类型

//获取今天的开始时间秒
func GetTodayUTCTime() int64 {
	timeStr := time.Now().Format(TimeTemplateDay)
	t, _ := time.Parse(TimeTemplateDay, timeStr)
	timeNumber := t.Unix()
	return timeNumber
}

//获取当天指定时间点时间戳
func GetCurDayTime(hour, min, second int64) int64 {
	curtime := time.Now().Unix()
	tm := time.Unix(curtime, 0)
	daytime := tm.Hour()*3600 + tm.Minute()*60 + tm.Second()
	addtime := hour*3600 + min*60 + second
	return curtime + int64(addtime) - int64(daytime)
}

//获得两个时间戳的天数差
func DiffDays(t1, t2 int64) int {
	s1, _ := db.StartEndTimeByTimestamp(t1)
	s2, _ := db.StartEndTimeByTimestamp(t2)
	return int(math.Abs(float64(s2-s1)) / (24 * 3600))
}

//是否在时间区间内
func IsInTimeRang(tstr1, tstr2 string) bool {
	t1, err1 := time.Parse(TimeTemplateSecond, tstr1)
	t2, err2 := time.Parse(TimeTemplateSecond, tstr1)
	if err1 != nil || err2 != nil {
		return false
	}
	if t1.Unix() < time.Now().Unix() && time.Now().Unix() < t2.Unix() {
		return true
	}
	return false
}
