package Message

import (
	"errors"
	log "github.com/cihub/seelog"
	"reflect"
)

type _IMessage interface {
	GetID() uint16
}

type _Def struct {
	msgMap map[uint16]_IMessage
	nameMap map[string]uint16
}

func newDef() *_Def {
	return &_Def{
		msgMap: make(map[uint16]_IMessage),
		nameMap: make(map[string]uint16),
	}
}

var def *_Def = newDef()

func (d *_Def) addDef(protoType _IMessage) {

	if _, ok := d.msgMap[protoType.GetID()]; ok {
		log.Error("message define had existed ", protoType.GetID())
		return
	}

	if protoType == nil {
		log.Error("message should implement _IMessage interface ")
		return
	}
	d.nameMap[reflect.TypeOf(protoType).Name()] = protoType.GetID()
	d.msgMap[protoType.GetID()] = protoType
}

func (d *_Def) fetchMessage(id uint16) (IMessage, error) {

	m, ok := d.msgMap[id]

	if !ok {
		return nil, errors.New("no message found ")
	}

	v := reflect.New(reflect.TypeOf(m).Elem())
	r, _ := v.Interface().(_IMessage)

	return r, nil
}

//FetchMessageByName 通过名字new消息结构
func FetchMessageByName(name string) (IMessage, error) {
	iter,ok := def.nameMap[name]
	if !ok {
		return nil,errors.New("not found msg name")
	}

	msg,err := def.fetchMessage(iter)
	return msg,err
}
