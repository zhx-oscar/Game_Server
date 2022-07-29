package internal

type _IActivity interface {
	Init() error
	Timer()
	Start()
	End()

	GetID() uint32
	GetStartTime() int64
	GetEndTime() int64
	GetLast() uint32
	GetInterval() uint32
	GetLoop() uint32
	GetKey() string
}

type IActivityItem interface {
	Init()
	Timer()
	IsActive() bool
	GetStep() uint32
	GetEndTime() int64
	GetStartTime() int64

	End()		//活动直接停止，测试用
	Start()		//活动直接开启，测试用
}

type IActivities interface {
	RegisterActivity(act _IActivity)
	GetActivity(id uint32) IActivityItem
}