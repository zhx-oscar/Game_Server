package User

import (
	"Cinder/Base/CRpc"
	"Cinder/Base/Core"
	"Cinder/Base/Message"
	"Cinder/Base/Prop"
)

type IUser interface {
	GetID() string
	GetType() string

	Rpc(srvType string, methodName string, args ...interface{}) chan *CRpc.RpcRet
	//RpcR(srvType string, methodName string, retCallback interface{}, args ...interface{})

	GetUserData() interface{}
	GetRealPtr() interface{}
	GetMgr() IUserMgr
	GetSrvInst() Core.ICore

	GetPropType() string
	GetProp() Prop.IProp

	SendToClient(msg Message.IMessage) error
	SendToPeerServer(srvType string, msg Message.IMessage) error
	SendToPeerUser(srvType string, msg Message.IMessage) error
	GetPeerServerID(srvType string) (string, error)

	Offline()
}

type IUserMgr interface {
	Destroy()
	GetOrCreateUser(id string, userData interface{}) (IUser, bool, error)
	CreateUser(id string, userData interface{}) (IUser, error)
	DestroyUser(id string) error
	GetUser(id string) (IUser, error)
	GetSrvInst() Core.ICore

	GetUserNum() int
	Loop()

	Traversal(cb func(user IUser) bool)
}

type IClientMessageSender interface {
	SendMessageToClient(msg Message.IMessage) error
}
