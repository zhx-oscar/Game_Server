package Space

import (
	BaseUser "Cinder/Base/User"
	"Cinder/Base/Util"
	"errors"
	"fmt"
	"reflect"
	"sync"
)

type _ActorPT struct {
	actorPT reflect.Type
}

type _ActorMgr struct {
	actorPTMap map[string]*_ActorPT
	actorMap   sync.Map
	space      _ISpace
}

func newActorMgr(space _ISpace) IActorMgr {
	mgr := &_ActorMgr{}

	mgr.actorPTMap = make(map[string]*_ActorPT)
	mgr.space = space

	return mgr
}

func (mgr *_ActorMgr) constructActor(actorType string) (interface{}, error) {

	pt, ok := mgr.actorPTMap[actorType]
	if !ok {
		return nil, errors.New("no actor prototype " + actorType)
	}

	d := reflect.New(pt.actorPT).Interface()
	return d, nil
}

func (mgr *_ActorMgr) AddActor(actorType string, actorID string, ownerUserID string, actorPropData []byte, userData interface{}) (string, error) {
	ii, err := mgr.constructActor(actorType)
	if err != nil {
		return "", err
	}

	var id string
	if actorID == "" {
		id = fmt.Sprint("actor_", Util.GetGUID())
	} else {
		id = actorID
	}

	if _, ok := mgr.actorMap.Load(id); ok {
		return "", errors.New("the id have existed " + id)
	}

	data := &_InitInfo{}
	data.ID = id
	data.Type = actorType
	data.Owner = mgr
	data.OwnerRealPtr = mgr.space
	data.RealPtr = ii
	data.OwnerUserID = ownerUserID
	data.PropData = actorPropData
	data.UserData = userData

	ii.(_IActor).InitBase(data)

	mgr.actorMap.Store(id, ii)

	oa, ok := ii.(BaseUser.IInit)
	if ok {
		oa.Init()
	}

	ii.(_IActor).SetReady()

	mgr.space.onAddActor(ii.(_IActor))

	is, ok := ii.(BaseUser.IStart)
	if ok {
		is.Start()
	}

	return id, nil
}

func (mgr *_ActorMgr) RemoveActor(id string) error {
	var actor interface{}
	var ok bool
	if actor, ok = mgr.actorMap.Load(id); !ok {
		return errors.New("no actor exist " + id)
	}

	mgr.space.onRemoveActor(actor.(_IActor))

	mgr.actorMap.Delete(id)

	od, ok := actor.(_IActor).GetRealPtr().(BaseUser.IDestroy)
	if ok {
		od.Destroy()
	}

	actor.(_IActor).DestroyBase()
	return nil
}

func (mgr *_ActorMgr) GetActor(actorID string) (IActor, error) {

	ii, ok := mgr.actorMap.Load(actorID)
	if !ok {
		return nil, errors.New("no existed")
	}

	return ii.(IActor), nil
}

func (mgr *_ActorMgr) UpdateActors() {
	//	mgr.removePendingActor()

	mgr.actorMap.Range(func(k, actor interface{}) bool {
		iu, ok := actor.(BaseUser.ILoopBase)
		if ok {
			iu.LoopBase()
		}

		ou, ok := actor.(BaseUser.ILoop)
		if ok {
			ou.Loop()
		}
		return true
	})
}

func (mgr *_ActorMgr) RegisterActor(actorType string, protoType IActor) {
	typ := reflect.Indirect(reflect.ValueOf(protoType)).Type()
	mgr.actorPTMap[actorType] = &_ActorPT{actorPT: typ}
}

func (mgr *_ActorMgr) TraversalActor(cb func(user IActor)) {

	mgr.actorMap.Range(func(key, value interface{}) bool {
		cb(value.(IActor))
		return true
	})

}

func (mgr *_ActorMgr) DestroyAllActor() {
	mgr.actorMap.Range(func(key, value interface{}) bool {
		mgr.RemoveActor(key.(string))
		return true
	})
}
