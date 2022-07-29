package DBAgent

import (
	"Cinder/Base/Core"
	"Cinder/Base/Prop"
	"Cinder/Base/Util"
	"Cinder/Cache"
	"Cinder/DB"
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	log "github.com/cihub/seelog"
)

const (
	WriteToDBInterval = 5 * time.Minute
)

var (
	errNoProp              = errors.New("no prop")
	errCantCreateCacheData = errors.New("can't create cache data")
)

type _Prop struct {
	id       string
	typ      string
	debugStr string

	ctx        context.Context
	cancelFunc context.CancelFunc

	destroyCtx        context.Context
	destroyCancelFunc context.CancelFunc

	isModified bool

	ownerMap sync.Map

	Util.ISafeCall
	Prop.IPropOwner
}

func (u *_Prop) Init(id string, typ string) error {

	u.ctx, u.cancelFunc = context.WithCancel(context.Background())
	u.destroyCtx, u.destroyCancelFunc = context.WithCancel(context.Background())

	u.id = id
	u.typ = typ
	u.debugStr = fmt.Sprintf("[P:%s:%s] ", typ, id)
	u.isModified = false

	u.ISafeCall = Util.NewSafeCall(u, true)

	u.IPropOwner = Core.Inst.NewPropOwner(u)
	u.IPropOwner.InitPropOwner(nil)

	if p := u.GetProp(); p != nil {
		p.GetCaller().SetParentCaller(u.ISafeCall)
	}

	err := u.readFromDB()
	if err != nil {
		log.Error(u.debugStr, "Init readFromDB err ", err)
		return err
	}

	u.WriteToCache()

	go u.MainLoop()
	return nil
}

func (u *_Prop) GetPropInfo() (string, string) {
	return u.typ, ""
}

func (u *_Prop) Destroy() {
	u.cancelFunc()
	u.ISafeCall.SafeCallDestroy()
	u.IPropOwner.DestroyPropOwner()
}

func (u *_Prop) destroySignal() <-chan struct{} {
	return u.destroyCtx.Done()
}

func (u *_Prop) MainLoop() {

	saveTicker := time.NewTicker(WriteToDBInterval)
	keepAliveTimer := time.NewTicker(Cache.PropObjectSrvIDExpireTime / 2)

	defer func() {
		saveTicker.Stop()
		keepAliveTimer.Stop()
	}()

forLoop:
	for {
		select {
		case <-saveTicker.C:
			if u.isModified {
				err := u.WriteToDB()
				if err != nil {
					log.Error(u.debugStr, "WriteToDB err ", err)
				}
				u.isModified = false

				u.WriteToCache()
			}
		case <-u.CallSignal():
			u.BatchCallMethod()
		case <-keepAliveTimer.C:
			Cache.KeepAlivePropDBSrvID(u.typ, u.id)
		case <-u.ctx.Done():
			break forLoop
		}
	}

	if u.isModified {
		err := u.WriteToDB()
		if err != nil {
			log.Error(u.debugStr, "WriteToDB err ", err)
		}
		u.isModified = false

		u.WriteToCache()

		log.Info(u.debugStr, "WriteToDB success")
	}
	u.destroyCancelFunc()
}

func (u *_Prop) Touch() {
	u.isModified = true
}

func (u *_Prop) GetPropData() ([]byte, error) {
	if u.GetProp() == nil {
		return nil, errors.New("no prop")
	}
	return u.GetProp().Marshal()
}

func (u *_Prop) readFromDB() error {

	util, err := DB.NewPropUtil(u.id, u.typ)
	if err != nil {
		return err
	}

	if im, ok := u.GetProp().(Prop.IBsonMarshaler); ok {
		var data []byte
		if data, err = util.GetBsonData(); err != nil {
			return err
		}

		return im.UnMarshalFromBson(data)

	} else {
		var data []byte
		if data, err = util.GetData(); err != nil {
			return err
		}

		return u.GetProp().UnMarshal(data)
	}
}

func (u *_Prop) WriteToDB() error {

	util, err := DB.NewPropUtil(u.id, u.typ)
	if err != nil {
		return err
	}

	if u.GetProp() == nil {
		return errNoProp
	}

	var jsonData []byte
	var binaryData []byte

	if im, ok := u.GetProp().(Prop.IBsonMarshaler); ok {
		jsonData, err = im.MarshalToBson()
	} else {
		binaryData, err = u.GetProp().Marshal()
	}

	if err != nil {
		return err
	}

	if jsonData != nil {
		err = util.SetBsonData(jsonData)
		if err != nil {
			return err
		}
	}

	if binaryData != nil {
		err = util.SetData(binaryData)
		if err != nil {
			return err
		}
	}

	return nil
}

func (u *_Prop) WriteToCache() error {
	if u.GetProp() == nil {
		return errNoProp
	}

	im, ok := u.GetProp().(Prop.ICacheProp)
	if !ok {
		return errCantCreateCacheData
	}

	data, err := im.MarshalCache()
	if err != nil {
		return err
	}

	return Cache.SetPropCache(u.typ, u.id, data)
}

func (u *_Prop) SyncProp(string, []byte, ...int) {}

/*

保存最后一次FetchPropData的服务器ID, 只处理这台服务器的属性同步消息
防止异常情况出现多个GameUser/SpaceUser存在的时候, 数据被写坏

*/

func (u *_Prop) SetCurrentOwner(typ, id string) {
	u.ownerMap.Store(typ, id)
}

func (u *_Prop) VerifyOwner(typ, id string) bool {
	cur, ok := u.ownerMap.Load(typ)
	if !ok {
		return false
	}
	return cur.(string) == id
}
