package Prop

import (
	"Cinder/Base/SrvNet"
	"Cinder/Base/Util"
	"Cinder/Cache"
	"Cinder/DB"
	"errors"
	log "github.com/cihub/seelog"
	"github.com/go-redis/redis/v7"
)

type IMgr interface {
	RegisterProp(propType string, propProto IProp)
	CreateProp(propType string) (IProp, error)

	RegisterPropObject(propObjectType string, object IPropObject)
	CreatePropObject(propObjectType string, id string, propData []byte, userData interface{}) (IPropObject, error)
	DestroyPropObject(id string) error

	NewPropOwner(realPtr interface{}) IPropOwner
}

type IProp interface {
	GetCaller() Util.ISafeCall

	Marshal() ([]byte, error)
	UnMarshal(data []byte) error

	MarshalPart() ([]byte, error)
	UnMarshalPart(data []byte) error
}

type IPropOwner interface {
	InitPropOwner(data []byte)
	DestroyPropOwner()

	FlushToDB() chan error
	FlushToCache() chan error

	GetProp() IProp
	GetPropID() string
	GetPropType() string
	GetDBSrvID() string
	GetSync() IPropSync
	GetSrvNode() SrvNet.INode
}

type IPropInfoFetcher interface {
	GetPropInfo() (propType string, propID string)
}

const (
	Target_Game          = 1
	Target_Space         = 3
	Target_Client        = 4
	Target_All_Clients   = 5
	Target_Other_Clients = 6
)

type IPropSync interface {
	SyncProp(methodName string, args []byte, targets ...int)
}

type IPropObject interface {
	IPropOwner
	GetID() string
	GetType() string
}

type IBsonMarshaler interface {
	MarshalToBson() ([]byte, error)
	UnMarshalFromBson(data []byte) error
}

// ICacheProp 缓存属性接口
type ICacheProp interface {
	MarshalCache() ([]byte, error)
	UnMarshalCache(data []byte) (interface{}, error)
}

var ErrCantCreateCacheProp = errors.New("can't create cache prop")
var ErrTypeIDMismatch = errors.New("types and ids mismatch")

func GetCacheProp(typ, id string) (interface{}, error) {
	prop, err := defaultMgr.CreateProp(typ)
	if err != nil {
		return nil, err
	}

	if im, ok := prop.(ICacheProp); ok {
		var data []byte
		if data, err = Cache.GetPropCache(typ, id); err == nil {
			return im.UnMarshalCache(data)
		} else if err == redis.Nil {
			// load from db
			var propUtil *DB.PropUtil
			if propUtil, err = DB.NewPropUtil(id, typ); err != nil {
				return nil, err
			}

			if ib, ibok := prop.(IBsonMarshaler); ibok {
				if data, err = propUtil.GetBsonData(); err != nil {
					return nil, err
				}
				if err = ib.UnMarshalFromBson(data); err != nil {
					return nil, err
				}
			} else {
				if data, err = propUtil.GetData(); err != nil {
					return nil, err
				}
				if err = prop.UnMarshal(data); err != nil {
					return nil, err
				}
			}

			// save to cache
			if data, err = im.MarshalCache(); err == nil {
				Cache.SetPropCache(typ, id, data)
			} else {
				log.Error("GetCacheProp MarshalCache err ", err)
			}

			return im.UnMarshalCache(data)

		} else {
			return nil, err
		}
	} else {
		return nil, ErrCantCreateCacheProp
	}
}

func GetBatchCacheProp(typ, id []string) ([]interface{}, error) {
	if len(typ) != len(id) {
		return nil, ErrTypeIDMismatch
	}

	if len(typ) == 0 {
		return nil, nil
	}

	propList := make([]IProp, 0, len(typ))
	for _, t := range typ {
		p, err := defaultMgr.CreateProp(t)
		if err != nil {
			return nil, err
		}

		if _, ok := p.(ICacheProp); !ok {
			return nil, ErrCantCreateCacheProp
		}

		propList = append(propList, p)
	}

	datas, err := Cache.GetPropCacheList(typ, id)
	if err != nil {
		return nil, err
	}

	cachePropList := make([]interface{}, 0, len(typ))
	cachePropTypes := make([]string, 0, 1)
	cachePropIDs := make([]string, 0, 1)
	cachePropDatas := make([][]byte, 0, 1)
	for i, data := range datas {
		propType := typ[i]
		propID := id[i]
		prop := propList[i]
		propCache := prop.(ICacheProp)

		if data == nil {
			// load from db
			var propUtil *DB.PropUtil
			if propUtil, err = DB.NewPropUtil(propID, propType); err != nil {
				return nil, err
			}

			if ib, ibok := prop.(IBsonMarshaler); ibok {
				if data, err = propUtil.GetBsonData(); err != nil {
					return nil, err
				}
				if err = ib.UnMarshalFromBson(data); err != nil {
					return nil, err
				}
			} else {
				if data, err = propUtil.GetData(); err != nil {
					return nil, err
				}
				if err = prop.UnMarshal(data); err != nil {
					return nil, err
				}
			}

			if data, err = propCache.MarshalCache(); err != nil {
				return nil, err
			}

			cachePropTypes = append(cachePropTypes, propType)
			cachePropIDs = append(cachePropIDs, propID)
			cachePropDatas = append(cachePropDatas, data)
		}

		var cacheProp interface{}
		if cacheProp, err = propCache.UnMarshalCache(data); err != nil {
			return nil, err
		}

		cachePropList = append(cachePropList, cacheProp)
	}

	if len(cachePropTypes) != 0 {
		if err = Cache.SetPropCacheList(cachePropTypes, cachePropIDs, cachePropDatas); err != nil {
			log.Error("GetBatchCacheProp SetPropCacheList err ", err)
		}
	}

	return cachePropList, nil
}
