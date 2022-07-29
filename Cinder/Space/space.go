package Space

import (
	"Cinder/Base/Core"
	"Cinder/Base/Prop"
	BaseUser "Cinder/Base/User"
	"Cinder/Base/Util"
	"context"
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type _ISpace interface {
	ISpace
	EnterSpace(userID string) error
	LeaveSpace(userID string) error
	onUserAgentChanged(userID string, oldAgentID string, newAgentID string)
	//onUserDataRet(pid string, userData []byte)
	onAddActor(actor _IActor)
	onRemoveActor(actor _IActor)

	onAddUser(user IUser)
	onRemoveUser(user IUser)

	onInit(data *_InitInfo)
	setDestroyFlag()
	destroySignal() <-chan struct{}

	GetSafeCall() Util.ISafeCall
}

type Space struct {
	userMgr BaseUser.IUserMgr

	IActorMgr
	Prop.IPropOwner

	Util.ISafeCall

	userData interface{}
	owner    ISpaceServer
	realPtr  interface{}
	looper   BaseUser.ILoop

	id       string
	debugStr string

	ownerUser IUser

	userAgentMap map[string][]string
	time         time.Time
	deltaTime    time.Duration

	fiveSecTicker *time.Ticker

	ctx        context.Context
	cancelFunc context.CancelFunc

	destroyCtx        context.Context
	destroyCancelFunc context.CancelFunc
}

func (space *Space) onInit(data *_InitInfo) {

	space.userAgentMap = make(map[string][]string)

	space.id = data.ID
	space.debugStr = fmt.Sprintf("[S:%s]", space.id)
	space.owner = data.Owner.(ISpaceServer)
	space.realPtr = data.RealPtr
	space.userData = data.UserData

	space.ctx, space.cancelFunc = context.WithCancel(context.Background())
	space.destroyCtx, space.destroyCancelFunc = context.WithCancel(context.Background())

	space.ISafeCall = Util.NewSafeCall(space.realPtr, viper.GetBool("Config.Recover"))

	space.userMgr = BaseUser.NewUserMgr(data.UserPT.(IUser), space.owner, false)
	space.IActorMgr = newActorMgr(space.realPtr.(_ISpace))

	if ii, ok := space.realPtr.(BaseUser.ILoop); ok {
		space.looper = ii
	}

	space.deltaTime = 0
	space.IPropOwner = Core.Inst.NewPropOwner(space.GetRealPtr())
	space.IPropOwner.InitPropOwner(data.PropData)

	if p := space.GetProp(); p != nil {
		p.GetCaller().SetParentCaller(space.ISafeCall)
	}

	space.fiveSecTicker = time.NewTicker(5 * time.Second)

	go space.mainLoop()
}

func (space *Space) mainLoop() {

	if ii, ok := space.realPtr.(BaseUser.IInit); ok {
		ii.Init()
	}

	ticker := time.NewTicker((1000 / 30) * time.Millisecond)
	defer ticker.Stop()
	lastTime := time.Now()
loop:
	for {
		select {
		case <-ticker.C:
			curTime := time.Now()
			dt := curTime.Sub(lastTime)
			space.setTime(curTime)
			space.setDeltaTime(dt)
			space.onLoop()

			lastTime = curTime
		case <-space.CallSignal():
			space.BatchCallMethod()
		case <-space.ctx.Done():
			break loop
		}
	}

	space.onDestroy()
	space.Debug("space mainloop exit")
}

func (space *Space) onDestroy() {

	space.IPropOwner.DestroyPropOwner()
	space.userMgr.Traversal(func(user BaseUser.IUser) bool {
		space.owner.(*_Server).CleanUserToSpaceMap(user.GetID())
		return true
	})

	space.fiveSecTicker.Stop()
	space.userMgr.Destroy()
	space.DestroyAllActor()

	ii, ok := space.realPtr.(BaseUser.IDestroy)
	if ok {
		ii.Destroy()
	}

	space.destroyCancelFunc()
	space.ISafeCall.SafeCallDestroy()
}

func (space *Space) setDestroyFlag() {
	space.cancelFunc()
}

func (space *Space) GetID() string {
	return space.id
}

func (space *Space) GetRealPtr() interface{} {
	return space.realPtr
}

func (space *Space) GetUserData() interface{} {
	return space.userData
}

func (space *Space) destroySignal() <-chan struct{} {
	return space.destroyCtx.Done()
}

func (space *Space) DestroySelf() {
	go func() {
		space.Debug("DestroySelf")
		_ = space.owner.DestroySpace(space.GetID())
	}()
}

func (space *Space) triggerTimerTicker() {
	select {
	case <-space.fiveSecTicker.C:
		space.refreshOwnerUser()
	default:
	}

}

func (space *Space) onLoop() {
	defer func() {
		if err := recover(); err != nil {
			space.Error("Space loop panic", err)
			if !viper.GetBool("Config.Recover") {
				panic(err)
			} else {
				space.Error(Util.GetPanicStackString())
			}
		}
	}()

	space.userMgr.Loop()
	space.UpdateActors()
	space.triggerTimerTicker()

	if space.looper != nil {
		space.looper.Loop()
	}
}

func (space *Space) GetOwnerUser() IUser {
	return space.ownerUser
}

func (space *Space) GetSafeCall() Util.ISafeCall {
	return space.ISafeCall
}
