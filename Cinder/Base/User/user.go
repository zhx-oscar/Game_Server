package User

import (
	"Cinder/Base/Const"
	"Cinder/Base/Core"
	"Cinder/Base/Message"
	"Cinder/Base/Prop"
	_ "Cinder/Base/ServerConfig"
	"Cinder/Base/Util"
	"Cinder/Cache"
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/spf13/viper"
)

type User struct {
	Util.ISafeCall
	id        string
	typ       string
	realPtr   interface{}
	userData  interface{}
	coroutine bool
	debugStr  string
	isReady   bool // 执行完LateInit之后, 才认为初始化完成

	Prop.IPropOwner

	mgr         IUserMgr
	ilooper     ILoop
	ilooperBase ILoopBase
	peers       sync.Map

	retList sync.Map

	ctx        context.Context
	cancelFunc context.CancelFunc

	destroyCtx        context.Context
	destroyCancelFunc context.CancelFunc

	keepSrvCtx        context.Context
	keepSrvCancelFunc context.CancelFunc
}

type ILoop interface {
	Loop()
}

type IInit interface {
	Init()
}

type IStart interface {
	Start()
}

type IDestroy interface {
	Destroy()
}

type ILoopBase interface {
	LoopBase()
}

type IInitBase interface {
	InitBase()
}

type ILateInitBase interface {
	LateInitBase()
}

type IDestroyBase interface {
	DestroyBase()
}

var (
	errNotInited = errors.New("init not finished")
)

func (u *User) onInit(info *_InitInfo) {

	u.id = info.id
	u.typ = info.typ
	u.mgr = info.mgr
	u.realPtr = info.realPtr
	u.userData = info.userData
	u.coroutine = info.coroutine

	u.ISafeCall = Util.NewSafeCall(info.realPtr, viper.GetBool("Config.Recover"))
	u.debugStr = fmt.Sprintf("[U:%s]", u.id)

	u.ctx, u.cancelFunc = context.WithCancel(context.Background())
	u.destroyCtx, u.destroyCancelFunc = context.WithCancel(context.Background())
	u.keepSrvCtx, u.keepSrvCancelFunc = context.WithCancel(context.Background())

	u.IPropOwner = Core.Inst.NewPropOwner(u.GetRealPtr())
	u.IPropOwner.InitPropOwner(nil)

	if p := u.GetProp(); p != nil {
		p.GetCaller().SetParentCaller(u.ISafeCall)
	}

	u.initPeers()
	u.regSrvInfo()
	u.notifyCreate()

	ii, ok := u.GetRealPtr().(ILoop)
	if ok {
		u.ilooper = ii
	}

	ib, ok := u.GetRealPtr().(ILoopBase)
	if ok {
		u.ilooperBase = ib
	}

	iib, ok := u.GetRealPtr().(IInitBase)
	if ok {
		iib.InitBase()
	}

	iii, ok := u.GetRealPtr().(IInit)
	if ok {
		iii.Init()
	}

	u.isReady = true
}

func (u *User) onLateInit() {
	ib, ok := u.GetRealPtr().(ILateInitBase)
	if ok {
		ib.LateInitBase()
	}

	is, ok := u.GetRealPtr().(IStart)
	if ok {
		is.Start()
	}

	if u.coroutine {
		go u.mainLoop()
	}
}

func (u *User) mainLoop() {
	defer func() {
		if err := recover(); err != nil {
			u.Error(err)
			if !viper.GetBool("Config.Recover") {
				panic(err)
			} else {
				u.Error(Util.GetPanicStackString())
			}
		}
	}()

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

loop:
	for {
		select {
		case <-u.ctx.Done():
			break loop
		case <-ticker.C:
			u.onLoop()
		case <-u.CallSignal():
			u.BatchCallMethod()
		}
	}

	u.onDestroy()
}

func (u *User) onDestroy() {

	u.unRegSrvInfo()
	u.notifyDestroy()

	u.IPropOwner.DestroyPropOwner()

	ii, ok := u.GetRealPtr().(IDestroy)
	if ok {
		ii.Destroy()
	}

	ib, ok := u.GetRealPtr().(IDestroyBase)
	if ok {
		ib.DestroyBase()
	}

	u.destroyCancelFunc()
	u.SafeCallDestroy()
}

func (u *User) regSrvInfo() {
	err := Cache.SetUserPeerRedis(u.GetID(), u.GetType(), u.GetSrvInst().GetServiceID())
	if err != nil {
		u.Error("regSrvInfo err", err)
	} else {
		go u.keepSrvAliveCo()
	}
	u.peers.Store(u.GetType(), u.GetSrvInst().GetServiceID())
}

func (u *User) unRegSrvInfo() {
	Cache.ClearUserPeerSrvID(u.GetID(), u.GetType(), u.GetSrvInst().GetServiceID())
	u.keepSrvCancelFunc()
}

func (u *User) keepSrvAliveCo() {
	ticker := time.NewTicker(Cache.UserPeerSrvIDExpireTime / 2)
	defer ticker.Stop()

	for {
		select {
		case <-u.keepSrvCtx.Done():
			return
		case <-ticker.C:
			Cache.KeepAliveUserPeerSrvID(u.GetID(), u.GetType())
		}
	}

}

func (u *User) onLoop() {
	if u.ilooperBase != nil {
		u.ilooperBase.LoopBase()
	}

	if u.ilooper != nil {
		u.ilooper.Loop()
	}
}

func (u *User) setDestroyFlag() {
	u.cancelFunc()
}

func (u *User) destroySignal() <-chan struct{} {
	return u.destroyCtx.Done()
}

func (u *User) GetID() string {
	return u.id
}

func (u *User) GetType() string {
	return u.typ
}

func (u *User) GetRealPtr() interface{} {
	return u.realPtr
}

func (u *User) GetPeers() *sync.Map {
	return &u.peers
}

func (u *User) GetMgr() IUserMgr {
	return u.mgr
}

func (u *User) GetUserData() interface{} {
	return u.userData
}

func (u *User) GetSrvInst() Core.ICore {
	return u.mgr.GetSrvInst()
}

func (u *User) Offline() {
	_ = u.SendToPeerServer(Const.Agent, &Message.UserLogoutReq{UserID: u.GetID()})
}

func (u *User) IsReady() bool {
	return u.isReady
}
