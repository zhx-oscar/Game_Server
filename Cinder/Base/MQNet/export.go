package MQNet

import (
	"Cinder/Base/Message"
)

type IService interface {
	Init(opts ...Option) error
	AddProc(proc IProc)
	Destroy()
	Post(addr string, message Message.IMessage) error
}

type IProc interface {
	MessageProc(srcAddr string, message Message.IMessage)
}
