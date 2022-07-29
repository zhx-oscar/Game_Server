package DBAgent

import (
	"errors"
	log "github.com/cihub/seelog"
	"sync"
)

var propMgr *_PropMgr

func newPropMgr() *_PropMgr {
	return &_PropMgr{}
}

type _PropMgr struct {
	props   sync.Map
	lockers sync.Map
}

func (mgr *_PropMgr) FetchOrCreate(id string, typ string) (*_Prop, error) {
	key := mgr.genKey(id, typ)

	locker := mgr.getLocker(key)
	locker.Lock()
	defer locker.Unlock()

	ii, ok := mgr.props.Load(key)
	if ok {
		return ii.(*_Prop), nil
	}

	prop := &_Prop{}
	if err := prop.Init(id, typ); err != nil {
		return nil, err
	}
	mgr.props.Store(key, prop)

	return prop, nil
}

func (mgr *_PropMgr) Get(id, typ string) (*_Prop, error) {
	ii, ok := mgr.props.Load(mgr.genKey(id, typ))
	if ok {
		return ii.(*_Prop), nil
	}

	return nil, errors.New("no prop")
}

func (mgr *_PropMgr) getLocker(key string) sync.Locker {
	i, _ := mgr.lockers.LoadOrStore(key, &sync.Mutex{})
	return i.(sync.Locker)
}

func (mgr *_PropMgr) Destroy() {

	mgr.props.Range(func(key, value interface{}) bool {
		prop := value.(*_Prop)
		prop.Destroy()
		<-prop.destroySignal()
		return true
	})

	log.Info("prop mgr destroy")
}

func (mgr *_PropMgr) genKey(id, typ string) string {
	return typ + "_" + id
}
