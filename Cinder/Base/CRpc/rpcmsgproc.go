package CRpc

import (
	"Cinder/Base/Message"
	"Cinder/Base/SrvNet"

	log "github.com/cihub/seelog"
)

type _RpcServerMsgProc struct {
	srv  *_Server
	node SrvNet.INode
}

func (sproc *_RpcServerMsgProc) MessageProc(srcAddr string, message Message.IMessage) {
	if message == nil {
		return
	}

	switch message.GetID() {
	case Message.ID_Rpc_Req:
		go func() {
			m := message.(*Message.RpcReq)
			retMsg := &Message.RpcRet{
				RetID: m.RetID,
				Ret:   nil,
				Err:   "",
			}

			args, err := Message.UnPackArgs(m.Args)
			if err != nil {
				log.Error("MessageProc rpc call failed, unpack arg failed ", m.MethodName, "  ", err)
				retMsg.Err = "unpack arg failed " + err.Error()
			} else {
				ret, err := sproc.srv.CallMethod(m.MethodName, args...)

				if err != nil {
					log.Error("MessageProc rpc call failed ", m.MethodName, " ", err)
					retMsg.Err = "rpc call failed " + err.Error()
				} else {
					retMsg.Ret = Message.PackArgs(ret...)
				}
			}

			if err = sproc.node.Send(srcAddr, retMsg); err != nil {
				log.Error("MessageProc Send retMsg err ", err, " RetMsg ", retMsg)
			}
		}()
	}
}

type _RpcClientMsgProc struct {
	cli  *_Client
	node SrvNet.INode
}

func (cproc *_RpcClientMsgProc) MessageProc(srcAddr string, message Message.IMessage) {
	if message == nil {
		return
	}

	switch message.GetID() {
	case Message.ID_Rpc_Ret:
		m := message.(*Message.RpcRet)

		args, err := Message.UnPackArgs(m.Ret)
		if err != nil {
			log.Error("MessageProc rpc call failed, unpack arg error ", err)
			m.Err = m.Err + " " + err.Error()
		}

		cproc.cli.onRpcRet(m.RetID, m.Err, args)
	}
}
