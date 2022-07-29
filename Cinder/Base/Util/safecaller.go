package Util

import (
	"errors"
	"fmt"
	log "github.com/cihub/seelog"
	"reflect"
	"sync"
)

type SafeCallRet struct {
	Err error
	Ret []interface{}
}

type ISafeCall interface {
	SafeCall(methodName string, args ...interface{}) chan *SafeCallRet
	SafeCallDestroy()
	CallMethod(methodName string, args ...interface{}) ([]interface{}, error)

	SetParentCaller(caller ISafeCall)
	CallSignal() <-chan bool
	BatchCallMethod()
}

func NewSafeCall(realPtr interface{}, recover bool) ISafeCall {
	s := &_SafeCall{
		realPtr:     realPtr,
		pendingMsgs: NewSafeList(),
		recover:     recover,
		isDestroy:   false,
	}
	return s
}

type _SafeCall struct {
	realPtr     interface{}
	parent      *_SafeCall
	cacheMethod sync.Map
	pendingMsgs ISafeList
	recover     bool
	isDestroy   bool
}

type _CallMsg struct {
	caller     ISafeCall
	methodName string
	args       []interface{}
	ret        chan *SafeCallRet
}

func newCallMsg(caller *_SafeCall, methodName string, args []interface{}) *_CallMsg {
	return &_CallMsg{
		caller:     caller,
		methodName: methodName,
		args:       args,
		ret:        make(chan *SafeCallRet, 1),
	}
}

func (caller *_SafeCall) CallMethod(methodName string, args ...interface{}) ([]interface{}, error) {

	defer func() {
		if err := recover(); err != nil {
			str := fmt.Sprintf("CallMethod %s err %v", methodName, err)
			log.Error(str)
			if !caller.recover {
				panic(str)
			} else {
				log.Error(GetPanicStackString())
			}
		}
	}()

	var m reflect.Value

	im, ok := caller.cacheMethod.Load(methodName)
	if ok {
		m = im.(reflect.Value)
	} else {
		m = reflect.ValueOf(caller.realPtr).MethodByName(methodName)
		if m.IsValid() {
			caller.cacheMethod.Store(methodName, m)
		}
	}

	if !m.IsValid() {
		return nil, errors.New("no method found")
	}

	argValues := make([]reflect.Value, 0, 3)

	//argValues = append(argValues, reflect.ValueOf(caller.realPtr))
	for _, arg := range args {
		argValues = append(argValues, reflect.ValueOf(arg))
	}

	retValues := m.Call(argValues)

	rets := make([]interface{}, 0, 3)

	if retValues != nil {
		for _, v := range retValues {
			rets = append(rets, v.Interface())
		}
	}

	return rets, nil
}

func (caller *_SafeCall) SafeCall(methodName string, args ...interface{}) chan *SafeCallRet {

	if caller.isDestroy {

		ret := make(chan *SafeCallRet, 1)
		sret := &SafeCallRet{
			Err: errors.New("the object had destroyed"),
			Ret: nil,
		}
		ret <- sret
		close(ret)
		return ret

	} else {
		return caller.pushCallMsg(caller, methodName, args...)
	}
}

func (caller *_SafeCall) pushCallMsg(srcCaller *_SafeCall, methodName string, args ...interface{}) chan *SafeCallRet {

	if caller.parent == nil {
		msg := newCallMsg(srcCaller, methodName, args)
		caller.pendingMsgs.Put(msg)
		return msg.ret
	} else {
		return caller.parent.pushCallMsg(srcCaller, methodName, args...)
	}
}

func (caller *_SafeCall) SafeCallDestroy() {
	caller.isDestroy = true
	caller.BatchCallMethod()
}

func (caller *_SafeCall) SetParentCaller(c ISafeCall) {
	caller.parent = c.(*_SafeCall)
}

func (caller *_SafeCall) CallSignal() <-chan bool {
	if caller.parent != nil {
		return caller.parent.pendingMsgs.Signal()
	}
	return caller.pendingMsgs.Signal()
}

func (caller *_SafeCall) BatchCallMethod() {
	for {
		im, err := caller.pendingMsgs.Pop()
		if err != nil {
			break
		}

		msg := im.(*_CallMsg)
		ret, err := msg.caller.CallMethod(msg.methodName, msg.args...)

		sret := &SafeCallRet{
			Err: err,
			Ret: ret,
		}

		msg.ret <- sret
		close(msg.ret)
	}
}
