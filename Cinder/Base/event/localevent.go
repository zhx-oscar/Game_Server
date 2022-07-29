package event

import (
	"errors"
	"reflect"
)

/*
	localEvent 本地事件分发，非线程安全
*/

var (
	errNotFindEventName = errors.New("not find eventName")
	errNotFindEventID   = errors.New("not find eventID")
)

type localEventDispatcher struct {
	eventNameMap map[string][]int
	eventMap     map[int]interface{}
	index        int
}

func (e *localEventDispatcher) init() {
	e.eventNameMap = make(map[string][]int)
	e.eventMap = make(map[int]interface{})
	e.index = 1
}

//RegLocalEvent 注册本地事件
func (e *localEventDispatcher) RegLocalEvent(event string, callBack interface{}) int {
	eventID := e.index
	e.index++
	e.eventMap[eventID] = callBack
	e.eventNameMap[event] = append(e.eventNameMap[event], eventID)
	return eventID
}

//FireLocalEvent 触发本地事件
func (e *localEventDispatcher) FireLocalEvent(event string, args ...interface{}) {
	eventIDList, ok := e.eventNameMap[event]
	if !ok {
		return
	}

	if len(eventIDList) == 0 {
		return
	}

	var callArgs []reflect.Value
	for _, arg := range args {
		callArgs = append(callArgs, reflect.ValueOf(arg))
	}

	for _, id := range eventIDList {
		f, ok := e.eventMap[id]
		if !ok {
			continue
		}

		reflect.ValueOf(f).Call(callArgs)
	}
}

//UnRegLocalEvent 注销本地事件
func (e *localEventDispatcher) UnRegLocalEvent(event string, handle int) error {
	_, ok := e.eventNameMap[event]
	if !ok {
		return errNotFindEventName
	}

	for i := range e.eventNameMap[event] {
		if e.eventNameMap[event][i] == handle {
			e.eventNameMap[event] = append(e.eventNameMap[event][:i], e.eventNameMap[event][i+1:]...)
			return nil
		}
	}

	return errNotFindEventID
}
