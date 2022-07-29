package CRpc

import (
	"Cinder/Base/SrvNet"
	"Cinder/Base/Util"
	"github.com/spf13/viper"
)

type _Server struct {
	Util.ISafeCall
}

func (srv *_Server) Init(srvNode SrvNet.INode, rpcProc interface{}) error {

	srv.ISafeCall = Util.NewSafeCall(rpcProc, viper.GetBool("Config.Recover"))

	srvNode.AddMessageProc(&_RpcServerMsgProc{
		srv:  srv,
		node: srvNode,
	})

	return nil
}

func (srv *_Server) Destroy() {
	srv.ISafeCall.SafeCallDestroy()
}
