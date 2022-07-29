package Core

import (
	"Cinder/Base/CRpc"
	"Cinder/Base/Net"
	"Cinder/Base/Prop"
	"Cinder/Base/SrvNet"
)

type Info struct {
	ServiceType string
	AreaID      string
	ServiceID   string

	PortMin    int
	PortMax    int
	ListenAddr string
	OuterAddr  string

	RpcProc              interface{}
	NetServerMessageProc Net.IProc
}

func NewDefaultInfo() *Info {
	return &Info{
		PortMax:    30000,
		PortMin:    20000,
		ListenAddr: "0.0.0.0",
	}
}

type IServerQuery interface {
	GetServiceID() string
	GetServiceType() string
}

type ICore interface {
	SrvNet.INodeCaller
	SrvNet.INodeQuery
	IServerQuery
	CRpc.IRpcCaller
	Prop.IMgr

	Init(info *Info) error
	Destroy()

	GetNetNode() SrvNet.INode
}

var Inst ICore
