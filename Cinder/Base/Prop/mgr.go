package Prop

import (
	"Cinder/Base/CRpc"
	"Cinder/Base/SrvNet"
	"errors"
	"reflect"
	"sync"
)

func NewMgr(srvNode SrvNet.INode, rpcClient CRpc.IClient) IMgr {
	mgr := &_Mgr{
		propProtoTypes:       make(map[string]reflect.Type),
		propObjectProtoTypes: make(map[string]reflect.Type),
		srvNode:              srvNode,
		rpcClient:            rpcClient,
	}

	srvNode.AddMessageProc(mgr)
	defaultMgr = mgr
	return mgr
}

var defaultMgr IMgr

type _Mgr struct {
	propProtoTypes map[string]reflect.Type
	pendingRetC    sync.Map

	propObjectProtoTypes map[string]reflect.Type
	propObjects          sync.Map

	srvNode   SrvNet.INode
	rpcClient CRpc.IClient
}

func (mgr *_Mgr) RegisterProp(propType string, propProto IProp) {
	mgr.propProtoTypes[propType] = reflect.TypeOf(propProto).Elem()
}

func (mgr *_Mgr) CreateProp(propType string) (IProp, error) {

	typ, ok := mgr.propProtoTypes[propType]
	if !ok || typ == nil {
		return nil, errors.New("prototype not exist " + propType)
	}

	return reflect.New(typ).Interface().(IProp), nil
}

func (mgr *_Mgr) NewPropOwner(realPtr interface{}) IPropOwner {
	o := &_Owner{
		realPtr: realPtr,

		srvNode:      mgr.srvNode,
		typ:          "",
		id:           "",
		dbSrvID:      "",
		propDataInit: false,
		mgr:          mgr,
	}

	is, ok := realPtr.(IPropSync)
	if ok {
		o.syncer = is
	}

	ip, ok := realPtr.(IPropInfoFetcher)
	if ok {
		o.propInfoFetcher = ip
	}

	return o
}

type _IInit interface {
	Init()
}

type _IDestroy interface {
	Destroy()
}

func (mgr *_Mgr) RegisterPropObject(propObjectType string, object IPropObject) {
	mgr.propObjectProtoTypes[propObjectType] = reflect.TypeOf(object).Elem()
}

func (mgr *_Mgr) CreatePropObject(propObjectType string, id string, propData []byte, userData interface{}) (IPropObject, error) {
	typ, ok := mgr.propObjectProtoTypes[propObjectType]
	if !ok || typ == nil {
		return nil, errors.New("prototype not exist " + propObjectType)
	}

	if _, ok = mgr.propObjects.Load(id); ok {
		return nil, errors.New("prop object have existed " + id)
	}

	ref := reflect.New(typ).Interface()

	if err := ref.(_IPropObject).initPropObject(propObjectType, id, propData, mgr, ref, userData); err != nil {
		return nil, err
	}

	ii, ok := ref.(_IInit)
	if ok {
		ii.Init()
	}

	mgr.propObjects.Store(id, ref)

	return ref.(IPropObject), nil
}

func (mgr *_Mgr) GetPropObject(id string) (IPropObject, error) {
	ii, ok := mgr.propObjects.Load(id)
	if !ok {
		return nil, errors.New("no prop object")
	}

	return ii.(IPropObject), nil
}

func (mgr *_Mgr) DestroyPropObject(id string) error {

	ii, ok := mgr.propObjects.Load(id)
	if !ok {
		return errors.New("no prop object")
	}

	iid, ok := ii.(_IDestroy)
	if ok {
		iid.Destroy()
	}

	ii.(_IPropObject).destroyObject()
	mgr.propObjects.Delete(id)
	return nil
}

func (mgr *_Mgr) GetCacheProp(typ, id string) (ICacheProp, error) {
	return nil, nil
}

func (mgr *_Mgr) GetBatchCacheProp(typ, id []string) ([]ICacheProp, error) {
	return nil, nil
}
