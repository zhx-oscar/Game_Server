package User

import (
	"Cinder/Base/Core"
	"Cinder/Base/Message"
	"Cinder/Base/Prop"
	"Cinder/Base/Util"
	"errors"
	"reflect"
	"sync"
	"sync/atomic"
)

type _UserMgr struct {
	protoType reflect.Type
	srvInst   Core.ICore
	coroutine bool
	srvType   string
	userLock  sync.Map

	userMap sync.Map
	userNum int32
}

type _InitInfo struct {
	id        string
	typ       string
	realPtr   interface{}
	userData  interface{}
	mgr       IUserMgr
	coroutine bool
}

type _IUser interface {
	Util.ISafeCall
	Prop.IPropOwner

	onInit(info *_InitInfo)
	onLateInit()
	onDestroy()
	onLoop()

	GetID() string
	GetType() string

	setDestroyFlag()
	GetMgr() IUserMgr
	GetSrvInst() Core.ICore

	GetRealPtr() interface{}

	OnRpcRet(retID string, err string, ret []interface{})
	destroySignal() <-chan struct{}

	syncPeerCreate(srvID string, srvType string)
	syncPeerDestroy(srvID string, srvType string)

	SendToPeerServer(srvType string, msg Message.IMessage) error
	SendToClient(msg Message.IMessage) error
}

func NewUserMgr(protoType IUser, srvInst Core.ICore, coroutine bool) IUserMgr {
	mgr := &_UserMgr{}
	mgr.protoType = reflect.TypeOf(protoType).Elem()

	mgr.srvInst = srvInst
	mgr.coroutine = coroutine
	mgr.srvType = srvInst.GetServiceType()
	mgr.userNum = 0

	return mgr
}

func (mgr *_UserMgr) Destroy() {
	if mgr.coroutine {
		mgr.userMap.Range(func(key, value interface{}) bool {
			value.(_IUser).setDestroyFlag()
			//<-value.(_IUser).destroySignal()
			return true
		})

	} else {
		mgr.userMap.Range(func(key, value interface{}) bool {
			value.(_IUser).onDestroy()
			return true
		})
	}
}

func (mgr *_UserMgr) GetUserNum() int {
	return int(mgr.userNum)
}

func (mgr *_UserMgr) GetOrCreateUser(id string, userData interface{}) (IUser, bool, error) {

	locker := mgr.getUserLocker(id)
	locker.Lock()
	defer locker.Unlock()

	iu, err := mgr.GetUser(id)
	if err != nil {
		u, err := mgr.createUser(id, userData)
		return u, true, err
	} else {
		return iu, false, nil
	}
}

func (mgr *_UserMgr) CreateUser(id string, userData interface{}) (IUser, error) {
	locker := mgr.getUserLocker(id)
	locker.Lock()
	defer locker.Unlock()

	return mgr.createUser(id, userData)
}

func (mgr *_UserMgr) createUser(id string, userData interface{}) (IUser, error) {

	_, ok := mgr.userMap.Load(id)
	if ok {
		return nil, errors.New("the user existed " + id)
	}

	iu := reflect.New(mgr.protoType).Interface().(_IUser)

	info := &_InitInfo{
		id:        id,
		typ:       mgr.GetSrvInst().GetServiceType(),
		realPtr:   iu,
		userData:  userData,
		mgr:       mgr,
		coroutine: mgr.coroutine,
	}

	iu.onInit(info)
	mgr.userMap.Store(iu.GetID(), iu)
	iu.onLateInit()
	atomic.AddInt32(&mgr.userNum, 1)

	return iu.(IUser), nil
}

func (mgr *_UserMgr) DestroyUser(id string) error {

	locker := mgr.getUserLocker(id)
	locker.Lock()
	defer locker.Unlock()

	ii, ok := mgr.userMap.Load(id)
	if !ok {
		return errors.New("the user not existed " + id)
	}

	iu := ii.(_IUser)
	mgr.userMap.Delete(id)
	atomic.AddInt32(&mgr.userNum, -1)

	if mgr.coroutine {
		iu.setDestroyFlag()
		<-iu.destroySignal()
	} else {
		iu.onDestroy()
	}
	return nil
}

func (mgr *_UserMgr) getUserLocker(userID string) sync.Locker {
	ii, _ := mgr.userLock.LoadOrStore(userID, &sync.Mutex{})
	return ii.(sync.Locker)
}

func (mgr *_UserMgr) GetUser(id string) (IUser, error) {

	iu, ok := mgr.userMap.Load(id)
	if !ok {
		return nil, errors.New("the user not existed " + id)
	}

	return iu.(IUser), nil
}

func (mgr *_UserMgr) Loop() {
	if !mgr.coroutine {
		mgr.userMap.Range(func(key, value interface{}) bool {
			value.(_IUser).onLoop()
			return true
		})
	}
}

func (mgr *_UserMgr) GetSrvInst() Core.ICore {
	return mgr.srvInst
}

func (mgr *_UserMgr) Traversal(cb func(user IUser) bool) {

	mgr.userMap.Range(func(key, value interface{}) bool {
		return cb(value.(IUser))
	})

}
