package event

//ILocalEventDispatcher 本地事件分发结构
type ILocalEventDispatcher interface {
	FireLocalEvent(event string, args ...interface{})
	RegLocalEvent(event string, callBack interface{}) int
	UnRegLocalEvent(event string, handle int) error
}

//GetLocalEventDispatcher 获取本地事件分发器
func GetLocalEventDispatcher() ILocalEventDispatcher {
	l := &localEventDispatcher{}
	l.init()
	return l
}
