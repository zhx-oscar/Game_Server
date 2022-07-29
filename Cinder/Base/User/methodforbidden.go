package User

import (
	"Cinder/Base/Config"
	log "github.com/cihub/seelog"
	"sync"
)

var _InstMethodForbidden *_MethodForbidden

type _MethodForbidden struct {
	methodForbiddenMap sync.Map
}

func init() {
	if _InstMethodForbidden != nil {
		return
	}

	_InstMethodForbidden = &_MethodForbidden{
		methodForbiddenMap: sync.Map{},
	}
	_InstMethodForbidden.init()
}

const (
	ForbiddenListPrefix = "ForbiddenListPrefix/"
)

func (p *_MethodForbidden) getForbiddenPrefix() string {
	return ForbiddenListPrefix
}

func (p *_MethodForbidden) getFuncNameByForbiddenKey(key string) string {
	if len(key) <= len(ForbiddenListPrefix) {
		log.Error("ForbiddenListPrefix key is Illegal :", key)
		return ""
	}
	return key[len(ForbiddenListPrefix):]
}

func (p *_MethodForbidden) init() {
	keys, _, err := Config.Inst.GetValuesByPrefix(p.getForbiddenPrefix())
	if err != nil {
		panic("GetValuesByPrefix methodForbiddenList failed")
	}
	for _, key := range keys {
		funcName := p.getFuncNameByForbiddenKey(key)
		p.methodForbiddenMap.Store(funcName, true)
	}

	_, err = Config.Inst.WatchKeys(p.getForbiddenPrefix(), p.watchMethodForbidden)
	if err != nil {
		panic("watch methodForbiddenList failed")
	}
}

func (p *_MethodForbidden) watchMethodForbidden(opType int, rid string, value string) {
	funcName := p.getFuncNameByForbiddenKey(rid)
	if opType == Config.KeyAdd {
		p.methodForbiddenMap.Store(funcName, true)
	} else if opType == Config.KeyDelete {
		p.methodForbiddenMap.Delete(funcName)
	}
}

//isMethodForbidden 当前方法是否被禁止
func (p *_MethodForbidden) isMethodForbidden(methodName string) bool {
	_, ok := p.methodForbiddenMap.Load(methodName)
	return ok
}
