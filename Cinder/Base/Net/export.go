package Net

import "Cinder/Base/Message"

type ISess interface {
	SetSendSecretKey(key []byte)
	SetRecvSecretKey(key []byte)

	Send(message Message.IMessage , msgNo uint32) error
	Read() (Message.IMessage , uint32 , error)
	SetData(data interface{})
	GetData() interface{}
	Close()
	SetValidate()
	IsValidate() bool
}

type IProc interface {
	OnSessConnected(sess ISess)
	OnSessClosed(sess ISess)
	OnSessMessageHandle(sess ISess,msgNo uint32 , message Message.IMessage)
}

type IService interface {
	Register(handler IProc)
	Init(addr string) error
	Destroy()
}

type IClientPool interface {
	Register(handler IProc)
	Init(interestedNetSrvTypes []string) error
	Destroy()
}
