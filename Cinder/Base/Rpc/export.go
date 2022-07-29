package Rpc

type IService interface {
	Init(addr string, rpcProc interface{}) error
	Destroy()
}

type IClientPool interface {
	IRpcCall
	Init() error
	Destroy()
}

type IRpcCall interface {
	CallBySrvType(srvType string, serviceMethod string, args interface{}, reply interface{}) error
	CallBySrvID(srvType, srvID string, serviceMethod string, args interface{}, reply interface{}) error
}

type GameRpcReq struct {
	UserID     string
	MethodName string
	Req        string
}

type GameRpcRet struct {
	Err string
	Ret string
}
