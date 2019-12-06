package algo

import "time"

const (
	NEW_DAY_TIME = 7 // 每天7点重置
)

//获取当天指定时间点时间戳
func GetCurDayTime(hour, min, second int64) int64 {
	curtime := time.Now().Unix()
	tm := time.Unix(curtime, 0)
	daytime := tm.Hour()*3600 + tm.Minute()*60 + tm.Second()
	addtime := hour*3600 + min*60 + second
	return curtime + int64(addtime) - int64(daytime)
}

//是否需要刷新数据(每天hour时刷新数据)
func IsNeedFlush(lastFlushTime int64, hour int64) bool {
	tm1 := time.Unix(lastFlushTime, 0)
	curtime := time.Now().Unix()
	tm2 := time.Unix(curtime, 0)
	if tm2.Hour() >= int(hour) {
		if tm1.Day() != tm2.Day() {
			return true // 不同天
		} else {
			if tm1.Hour() < tm2.Hour() {
				return true
			}
		}
		return false
	} else {
		if GetCurDayTime(hour, 0, 0)-24*3600 > lastFlushTime {
			return true
		}
	}
	return false
}

//获取到下次重置时间
func GetNewDayExpireTime() int64 {
	tm := time.Unix(time.Now().Unix(), 0)
	daytime := tm.Hour()*3600 + tm.Minute()*60 + tm.Second()
	if tm.Hour() >= NEW_DAY_TIME {
		return int64(NEW_DAY_TIME*3600 + 24*3600 - daytime)
	}
	return int64(NEW_DAY_TIME*3600 - daytime)
}
