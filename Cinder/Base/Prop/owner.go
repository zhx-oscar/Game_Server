package Prop

import (
	"Cinder/Base/Const"
	"Cinder/Base/Message"
	"Cinder/Base/SrvNet"
	"Cinder/Base/Util"
	"Cinder/Cache"
	"context"
	"errors"
	"fmt"
	"time"
)

type _Owner struct {
	realPtr         interface{}
	syncer          IPropSync
	prop            IProp
	propInfoFetcher IPropInfoFetcher
	srvNode         SrvNet.INode
	debugStr        string

	mgr *_Mgr

	typ          string
	id           string
	dbSrvID      string
	propDataInit bool
}

func (o *_Owner) InitPropOwner(data []byte) {
	var err error

	if err = o.initProp(data); err != nil {
		o.Error("init prop failed", err)
		return
	}
	if err = o.initPropDBPart(); err != nil {
		o.Error("init prop db part failed", err)
		return
	}
}

func (o *_Owner) DestroyPropOwner() {
	if o.prop != nil {
		o.prop.(_IProp).DestroyProp()
	}
}

func (o *_Owner) initProp(data []byte) error {

	if o.propInfoFetcher == nil {
		o.debugStr = fmt.Sprintf("[Owner]")
		return nil
	}

	propType, propID := o.propInfoFetcher.GetPropInfo()
	if propType == "" {
		o.debugStr = fmt.Sprintf("[O:%s]", propID)
		return nil
	}

	o.debugStr = fmt.Sprintf("[O:%s:%s]", propType, propID)

	prop, err := o.mgr.CreateProp(propType)
	if err != nil {
		return err
	}

	err = prop.UnMarshal(data)
	if err != nil {
		return err
	}

	o.SetProp(prop)
	o.typ = propType
	o.id = propID

	return nil
}

func (o *_Owner) initPropDBPart() error {

	var err error

	if err = o.initDBSrvID(); err != nil {
		return err
	}
	if err = o.initPropFromDBSrv(); err != nil {
		return err
	}

	return nil
}

func (o *_Owner) initDBSrvID() error {
	if o.id == "" {
		return nil
	}

	srvID, err := o.srvNode.GetSrvIDByType(Const.DB)
	if err != nil {
		return errors.New("couldn't find db service available")
	}

	srvID, err = Cache.GetOrSetPropDBSrvID(o.typ, o.id, srvID)
	if err != nil {
		o.Error("couldn't set prop db srvID", err)
		return err
	}

	o.dbSrvID = srvID
	return nil
}

func (o *_Owner) initPropFromDBSrv() error {

	if o.id == "" || o.dbSrvID == "" {
		return nil
	}

	tempID := Util.GetGUID()

	msg := &Message.PropDataReq{
		TempID:      tempID,
		ID:          o.id,
		PropType:    o.typ,
		ServiceType: o.srvNode.GetType(),
	}

	retC, err := o.mgr.fetchRetC(tempID)
	if err != nil {
		return err
	}

	defer o.mgr.removeRetC(tempID)

	if err = o.srvNode.Send(o.dbSrvID, msg); err != nil {
		return err
	}

	ctx, ctxCancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer ctxCancel()

	select {
	case <-ctx.Done():
		return errors.New("get prop data timeout")
	case data := <-retC:
		o.propDataInit = true
		if data != nil {
			if err = o.GetProp().UnMarshal(data.([]byte)); err != nil {
				return err
			}
		}
	}

	return nil
}

func (o *_Owner) SetProp(prop IProp) {
	o.prop = prop

	prop.(_IProp).InitProp(o.realPtr.(IPropOwner), prop)
}

func (o *_Owner) GetProp() IProp {
	return o.prop
}

func (o *_Owner) GetPropID() string {
	return o.id
}

func (o *_Owner) GetPropType() string {
	return o.typ
}

func (o *_Owner) GetDBSrvID() string {
	return o.dbSrvID
}

func (o *_Owner) GetSync() IPropSync {
	return o.syncer
}

func (o *_Owner) GetSrvNode() SrvNet.INode {
	return o.srvNode
}

func (o *_Owner) FlushToDB() chan error {

	ret := make(chan error, 1)
	if o.dbSrvID == "" || o.id == "" {
		ret <- errors.New("couldn't prop save")
		close(ret)
		return ret
	}

	tempID := Util.GetGUID()

	msg := &Message.PropDataFlushReq{
		TempID:      tempID,
		ID:          o.id,
		Type:        o.typ,
		ServiceType: o.srvNode.GetType(),
	}

	err := o.srvNode.Send(o.dbSrvID, msg)
	if err != nil {
		ret <- err
		close(ret)
		return ret
	}

	retC, err := o.mgr.fetchRetC(tempID)
	if err != nil {
		ret <- err
		close(ret)
		return ret
	}

	go func() {
		ctx, ctxCancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer ctxCancel()

		select {
		case <-ctx.Done():
			ret <- errors.New("flush to db timeout")
		case retErr := <-retC:
			if retErr != "" {
				ret <- errors.New(retErr.(string))
			} else {
				ret <- nil
			}
		}

		close(ret)

		o.mgr.removeRetC(tempID)
	}()

	return ret
}

func (o *_Owner) FlushToCache() chan error {
	ret := make(chan error, 1)
	if o.dbSrvID == "" || o.id == "" {
		ret <- errors.New("couldn't prop save")
		close(ret)
		return ret
	}

	tempID := Util.GetGUID()

	msg := &Message.PropCacheFlushReq{
		TempID:      tempID,
		ID:          o.id,
		Type:        o.typ,
		ServiceType: o.srvNode.GetType(),
	}

	err := o.srvNode.Send(o.dbSrvID, msg)
	if err != nil {
		ret <- err
		close(ret)
		return ret
	}

	retC, err := o.mgr.fetchRetC(tempID)
	if err != nil {
		ret <- err
		close(ret)
		return ret
	}

	go func() {
		ctx, ctxCancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer ctxCancel()

		select {
		case <-ctx.Done():
			ret <- errors.New("flush to cache timeout")
		case retErr := <-retC:
			if retErr != "" {
				ret <- errors.New(retErr.(string))
			} else {
				ret <- nil
			}
		}

		close(ret)

		o.mgr.removeRetC(tempID)
	}()

	return ret
}
