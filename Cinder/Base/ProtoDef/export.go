package ProtoDef

import (
	"encoding/json"
	"errors"
	"github.com/golang/protobuf/proto"
	"reflect"
	"sync"
)

var protoMapInst sync.Map
var nameMapInst sync.Map
var protoDefData []byte

func AddDef(id int, proto interface{}) {
	protoMapInst.Store(id, proto)
	nameMapInst.Store(reflect.TypeOf(proto).Elem().Name(), id)
}

func GetIDByName(typeName string) (int, error) {

	i, ok := nameMapInst.Load(typeName)
	if !ok {
		return 0, errors.New("no type name exist " + typeName)
	}

	return i.(int), nil
}

func GetProtoMessageByID(id int) (proto.Message, error) {

	i, ok := protoMapInst.Load(id)

	if !ok {
		return nil, errors.New("no proto exist")
	}

	m, _ := reflect.New(reflect.TypeOf(i).Elem()).Interface().(proto.Message)
	return m, nil
}

func InitProtoDefData() {
	defMap := make(map[int]string)
	nameMapInst.Range(func(key, value interface{}) bool {
		defMap[value.(int)] = key.(string)
		return true
	})

	var err error
	protoDefData, err = json.Marshal(defMap)
	if err != nil {
		panic(err)
	}
}

func GetProtoDefData() []byte {
	return protoDefData
}
