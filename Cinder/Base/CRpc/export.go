package CRpc

import (
	"Cinder/Base/Message"
	"Cinder/Base/SrvNet"
	"Cinder/Base/Util"
)

type IService interface {
	Init(srvNode SrvNet.INode, rpcProc interface{}) error
	Destroy()
}

type IClient interface {
	Init(srvNode SrvNet.INode) error
	Destroy()

	IRpcCaller
}

type RpcRet struct {
	Ret  []interface{}
	Err  error
	Done chan *RpcRet
}

func NewRpcRet() (*RpcRet, string) {
	ret := &RpcRet{
		Done: make(chan *RpcRet, 1),
	}
	return ret, Util.GetGUID()
}

type IRpcCaller interface {
	RpcByID(srvID string, methodName string, args ...interface{}) chan *RpcRet
	RpcByType(srvType string, methodName string, args ...interface{}) chan *RpcRet

	CallRpcToUser(userID string, srvType string, methodName string, args ...interface{})
	CallRpcToUsers(userIDS []string, srvType string, methodName string, args ...interface{})
	CallRpcToAllUsers(srvType string, methodName string, args ...interface{})

	SendMessageToUser(userID string, srvType string, message Message.IMessage)
	SendMessageToUsers(userIDS []string, srvType string, message Message.IMessage)
	SendMessageToAllUsers(srvType string, message Message.IMessage)
}

func NewService() IService {
	return &_Server{}
}

func NewClient() IClient {
	return &_Client{}
}
