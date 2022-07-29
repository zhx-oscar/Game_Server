package SrvNet

import (
	"Cinder/Base/MQNet"
	"Cinder/Base/Message"
)

type INode interface {
	Init(areaID string, srvID string, srvType string) error
	Destroy()
	GetID() string
	GetType() string
	AddMessageProc(proc MQNet.IProc)
	SetLoadBlanceGetter(f LoadGetFunc)
	INodeCaller
	INodeQuery
}

type LoadGetFunc func() float32

type INodeCaller interface {
	Send(srvID string, msg Message.IMessage) error
	Broadcast(srvType string, msg Message.IMessage) error
}

type INodeQuery interface {
	GetSrvIDSByType(srvType string) ([]string, error)
	GetSrvIDByType(srvType string) (string, error)
	GetSrvTypeByID(srvID string) (string, error)
}
