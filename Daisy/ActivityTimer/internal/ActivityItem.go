package internal

import (
	"Daisy/DB"
	"time"
)

var battleActivities *Activities
func GetActivities() IActivities {
	if battleActivities == nil {
		battleActivities = &Activities{}
		battleActivities.Init()
	}
	return battleActivities
}

type ActivityItem struct {
	activity _IActivity	//活动接口

	startTime int64
	endTime int64
	nextStartTime int64
	active bool
	step uint32
}

func (act *ActivityItem) Init() {
	var needSave bool
	act.Load()
	now := time.Now().Unix()
	if now <= act.activity.GetStartTime() {
		act.startTime = act.activity.GetStartTime()
		act.step = 1
		act.nextStartTime = act.startTime
		act.endTime = act.nextStartTime + int64(act.activity.GetLast()+act.activity.GetInterval())
		needSave = true
	} else if act.step == 0{
		act.step = uint32(now - act.activity.GetStartTime())/(act.activity.GetLast()+act.activity.GetInterval()) + 1
		pass := (now - act.activity.GetStartTime())%int64(act.activity.GetLast()+act.activity.GetInterval())
		act.startTime = now - pass
		act.nextStartTime = act.startTime + int64(act.activity.GetLast()+act.activity.GetInterval())
		act.endTime = act.nextStartTime + int64(act.activity.GetLast()+act.activity.GetInterval())
		needSave = true
	}

	if now >= act.startTime && now < act.startTime+int64(act.activity.GetLast()) {
		act.active = true
		act.endTime = act.startTime + int64(act.activity.GetLast())
		needSave = true
	}

	if needSave == true{
		act.Save()
	}
	act.activity.Init()
}

func (act *ActivityItem) Timer() {
	var needSave bool
	now := time.Now().Unix()
	if act.active == true{
		if now >= act.endTime{
			act.active = false
			act.activity.End()
			act.nextStartTime = act.endTime + int64(act.activity.GetInterval())
			needSave = true
		}
	}
	if act.active == false{
		if act.activity.GetEndTime() != 0 && now >= act.activity.GetEndTime(){
			return
		}
		if act.step < act.activity.GetLoop() {
			if now >= act.nextStartTime {
				act.step++
				act.active = true
				act.endTime = act.nextStartTime + int64(act.activity.GetLast())
				act.activity.Start()
				needSave = true
			}
		}
	}
	if needSave == true{
		act.Save()
	}
	act.activity.Timer()
}

func (act *ActivityItem) Load() {
	data := DB.GetActivitiesInfo(act.activity.GetKey())
	act.startTime = data.StartTime
	act.endTime = data.EndTime
	act.nextStartTime = data.NextStartTime
	act.active = data.Active
	act.step = data.Step
}

func (act *ActivityItem) Save() {
	info := DB.ActivityInfo{
		StartTime:act.startTime,
		EndTime:act.endTime,
		NextStartTime:act.nextStartTime,
		Step:act.step,
		Active:act.active,
	}
	DB.SetActivitiesInfo(act.activity.GetKey(), &info)
}

func (act *ActivityItem) IsActive() bool {
	return act.active
}

func (act *ActivityItem) GetStep() uint32 {
	return act.step
}

func (act *ActivityItem) GetEndTime() int64{
	return act.endTime
}

func (act *ActivityItem) GetStartTime() int64{
	return act.nextStartTime
}

func (act *ActivityItem) End() {
	act.endTime = time.Now().Unix()
}

func (act *ActivityItem) Start() {
	act.nextStartTime = time.Now().Unix()
}

type Activities struct {
	activities map[uint32]*ActivityItem
}

func (acts *Activities) Init() bool {
	acts.activities = make(map[uint32]*ActivityItem)
	go func() {
		acts.Loop()
	}()
	return true
}

func (acts *Activities) Loop() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		<-ticker.C

		acts.Timer()
	}
}

func (acts *Activities) Timer() {
	for _,v := range acts.activities {
		v.Timer()
	}
}

func (acts *Activities) RegisterActivity(act _IActivity){
	act.Init()
	item := &ActivityItem{activity:act}
	acts.activities[act.GetID()] = item
	item.Init()
}

func (acts *Activities) GetActivity(id uint32) IActivityItem{
	return acts.activities[id]
}