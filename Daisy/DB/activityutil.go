package DB

import (
	"fmt"
	"strconv"
)

type ActivityInfo struct {
	StartTime int64
	EndTime int64
	NextStartTime int64
	Active bool
	Step uint32
}
const (
	StartKey = "start"
	EndKey = "end"
	NextKey = "next"
	ActiveKey = "active"
	StepKey = "step"
)
func GetActivitiesInfo(key string) *ActivityInfo{
	//Todo 处理多个Battle服的情况
	values,err := RedisDB.HMGet(key, StartKey, EndKey, NextKey, ActiveKey, StepKey).Result()
	result := &ActivityInfo{}
	if err != nil {
		return result
	}
	if values[0] != nil {
		v := values[0].(string)
		_v,_ := strconv.ParseInt(v, 10, 64)
		result.StartTime = _v
	}
	if values[1] != nil {
		v := values[1].(string)
		_v,_ := strconv.ParseInt(v, 10, 64)
		result.EndTime = _v
	}
	if values[2] != nil {
		v := values[2].(string)
		_v,_ := strconv.ParseInt(v, 10, 64)
		result.NextStartTime = _v
	}
	if values[3] != nil {
		v := values[3].(string)
		_v,_ := strconv.Atoi(v)
		if _v == 0 {
			result.Active = false
		} else {
			result.Active = true
		}
	}
	if values[4] != nil {
		v := values[4].(string)
		_v,_ := strconv.Atoi(v)
		result.Step = uint32(_v)
	}
	return result
}

func SetActivitiesInfo(key string, info *ActivityInfo){
	//Todo 处理多个Battle服的情况

	RedisDB.HMSet(key, StartKey, fmt.Sprint(info.StartTime))
	RedisDB.HMSet(key, EndKey, fmt.Sprint(info.EndTime))
	RedisDB.HMSet(key, NextKey, fmt.Sprint(info.NextStartTime))
	RedisDB.HMSet(key, StepKey, fmt.Sprint(info.Step))
	if info.Active == false{
		RedisDB.HMSet(key, ActiveKey, fmt.Sprint(0))
	}else {
		RedisDB.HMSet(key, ActiveKey, fmt.Sprint(1))
	}
}