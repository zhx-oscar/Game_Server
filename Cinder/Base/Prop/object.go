package Prop

import (
	"Cinder/Base/Const"
	"Cinder/Base/Message"
	_ "Cinder/Base/ServerConfig"
	"Cinder/Base/Util"
	"Cinder/Cache"
	"context"
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Object struct {
	IPropOwner
	Util.ISafeCall
	mgr      IMgr
	id       string
	typ      string
	realPtr  interface{}
	userData interface{}
	debugStr string

	looper _ILoop

	watcherIDS []string

	ctx        context.Context
	cancelFunc context.CancelFunc
	destroyC   chan struct{}
}

type _IPropObject interface {
	Util.ISafeCall
	GetPropData() []byte
	GetPropType() string
	AddWatcher(userID string)
	RemoveWatcher(userID string)
	DestroySelf()
	destroyObject()
	initPropObject(propObjectType string, id string, propData []byte, mgr IMgr, realPtr interface{}, userData interface{}) error
}

type _ILoop interface {
	Loop()
}

func (obj *Object) initPropObject(propObjectType string, id string, propData []byte, mgr IMgr, realPtr interface{}, userData interface{}) error {
	obj.id = id
	obj.typ = propObjectType
	obj.mgr = mgr
	obj.realPtr = realPtr
	obj.userData = userData
	obj.ISafeCall = Util.NewSafeCall(realPtr, viper.GetBool("Config.Recover"))
	obj.debugStr = fmt.Sprintf("[OB:%s:%s]", obj.typ, obj.id)

	obj.ctx, obj.cancelFunc = context.WithCancel(context.Background())
	obj.destroyC = make(chan struct{}, 1)

	ii, ok := realPtr.(_ILoop)
	if ok {
		obj.looper = ii
	}

	if err := obj.regPropObj(); err != nil {
		return err
	}

	obj.IPropOwner = obj.mgr.NewPropOwner(realPtr)
	obj.IPropOwner.InitPropOwner(propData)

	if p := obj.GetProp(); p != nil {
		p.GetCaller().SetParentCaller(obj.ISafeCall)
	}

	go obj.MainLoop()

	return nil
}

func (obj *Object) DestroyPropObject() {
	if obj.IPropOwner != nil {
		if err := <-obj.FlushToDB(); err != nil {
			obj.Error("DestroyPropObject FlushToDB err", err)
		}
		obj.unRegPropObj()
		obj.IPropOwner.DestroyPropOwner()
	}

	obj.ISafeCall.SafeCallDestroy()
}

// DestroySelf call in self coroutine loop
func (obj *Object) DestroySelf() {
	obj.cancelFunc()
}

func (obj *Object) destroyObject() {
	obj.cancelFunc()
	<-obj.destroyC
	close(obj.destroyC)
}

func (obj *Object) MainLoop() {
	defer func() {
		if err := recover(); err != nil {
			obj.Error("PropObj loop panic", err)
			if !viper.GetBool("Config.Recover") {
				panic(err)
			} else {
				obj.Error(Util.GetPanicStackString())
			}
		}
	}()

	ticker := time.NewTicker(100 * time.Millisecond)
	keepAliveTicker := time.NewTicker(Cache.PropObjectSrvIDExpireTime / 2)

	defer ticker.Stop()
	defer keepAliveTicker.Stop()

	for {
		select {
		case <-ticker.C:
			if obj.looper != nil {
				obj.looper.Loop()
			}
		case <-obj.CallSignal():
			obj.BatchCallMethod()
		case <-keepAliveTicker.C:
			Cache.KeepAlivePropObjectSrvID(obj.GetType(), obj.GetID())
		case <-obj.ctx.Done():

			iid, ok := obj.realPtr.(_IDestroy)
			if ok {
				iid.Destroy()
			}
			obj.DestroyPropObject()
			obj.destroyC <- struct{}{}
			return
		}
	}
}

func (obj *Object) GetID() string {
	return obj.id
}

func (obj *Object) GetType() string {
	return obj.typ
}

func (obj *Object) GetUserData() interface{} {
	return obj.userData
}

func (obj *Object) regPropObj() error {
	return Cache.SetPropObjectSrvID(obj.GetType(), obj.GetID(), obj.mgr.(*_Mgr).srvNode.GetID())
}

func (obj *Object) unRegPropObj() {
	Cache.ClearPropObjectSrvID(obj.GetType(), obj.GetID())
}

func (obj *Object) GetPropData() []byte {

	data, err := obj.GetProp().Marshal()
	if err != nil {
		obj.Error("GetPropData Marshal err", err)
		return []byte{}
	}

	return data
}

func (obj *Object) AddWatcher(userID string) {
	obj.watcherIDS = append(obj.watcherIDS, userID)
}

func (obj *Object) RemoveWatcher(userID string) {

	for i := 0; i < len(obj.watcherIDS); i++ {
		if obj.watcherIDS[i] == userID {
			obj.watcherIDS = append(obj.watcherIDS[0:i], obj.watcherIDS[i+1:]...)
			break
		}
	}
}

func (obj *Object) SyncProp(methodName string, args []byte, targets ...int) {

	as, _ := Message.UnPackArgs(args)
	_, _ = obj.GetProp().GetCaller().CallMethod(methodName, as...)

	msg := &Message.PropObjectPropNotify{
		ObjID:      obj.GetID(),
		MethodName: methodName,
		Args:       args,
	}

	for _, target := range targets {

		switch target {
		case Target_All_Clients:
			obj.sendToWatchers(msg)
		default:
			obj.Debug("no support target type", target)
		}
	}
}

func (obj *Object) sendToWatchers(msg Message.IMessage) {
	step := len(obj.watcherIDS) / 100
	if len(obj.watcherIDS)%100 > 0 {
		step++
	}
	for i := 0; i < step; i++ {
		start := i * 100
		end := (i + 1) * 100
		if end > len(obj.watcherIDS) {
			end = len(obj.watcherIDS)
		}

		obj.SendMessageToUsers(obj.watcherIDS[start:end], Const.Agent, msg)
	}
}

func (obj *Object) SendMessageToUsers(userIDS []string, srvType string, message Message.IMessage) {
	obj.mgr.(*_Mgr).rpcClient.SendMessageToUsers(userIDS, srvType, message)
}
