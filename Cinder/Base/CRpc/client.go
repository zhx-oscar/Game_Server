package CRpc

import (
	"Cinder/Base/Message"
	"Cinder/Base/SrvNet"
	"errors"
	"sync"
	"time"

	log "github.com/cihub/seelog"
)

type _Client struct {
	retList sync.Map
	node    SrvNet.INode
}

func (cli *_Client) Init(srvNode SrvNet.INode) error {

	cli.node = srvNode
	srvNode.AddMessageProc(&_RpcClientMsgProc{
		cli:  cli,
		node: srvNode,
	})

	return nil
}

func (cli *_Client) Destroy() {

}

func (cli *_Client) RpcByID(srvID string, methodName string, args ...interface{}) chan *RpcRet {

	ret, retID := NewRpcRet()

	msg := &Message.RpcReq{
		MethodName: methodName,
		Args:       Message.PackArgs(args...),
		RetID:      retID,
	}

	err := cli.node.Send(srvID, msg)
	if err != nil {
		ret.Err = err
		ret.Done <- ret
		close(ret.Done)

		return ret.Done
	}

	cli.retList.Store(retID, ret)

	go func() {
		time.Sleep(3 * time.Second)

		ii, _ := cli.retList.LoadOrStore(retID, nil)
		cli.retList.Delete(retID)
		if ii != nil {
			r := ii.(*RpcRet)
			r.Err = errors.New("time out " + retID)
			select {
			case r.Done <- r:
			default:
			}
		}
	}()

	return ret.Done
}

func (cli *_Client) RpcByType(srvType string, methodName string, args ...interface{}) chan *RpcRet {

	srvID, err := cli.node.GetSrvIDByType(srvType)
	if err != nil {
		ret, _ := NewRpcRet()
		ret.Err = err
		ret.Done <- ret
		close(ret.Done)
		return ret.Done
	}

	return cli.RpcByID(srvID, methodName, args...)
}

func (cli *_Client) onRpcRet(retID string, err string, ret []interface{}) {

	ii, _ := cli.retList.LoadOrStore(retID, nil)
	cli.retList.Delete(retID)

	if ii == nil {
		log.Error("onRpcRet fail , get nil info ", retID)
		return
	}

	info := ii.(*RpcRet)
	if err != "" {
		info.Err = errors.New(err)
	}
	info.Ret = ret
	select {
	case info.Done <- info:
	default:
	}

	close(info.Done)
}

func (cli *_Client) CallRpcToUser(userID string, srvType string, methodName string, args ...interface{}) {

	msg := &Message.UserRpcReq{
		UserID:     userID,
		MethodName: methodName,
		Args:       Message.PackArgs(args...),
		RetID:      "",
	}

	_ = cli.node.Broadcast(srvType, msg)
}

func (cli *_Client) CallRpcToUsers(userIDS []string, srvType string, methodName string, args ...interface{}) {

	msg := &Message.UsersRpcReq{
		UserIDS:    userIDS,
		MethodName: methodName,
		Args:       Message.PackArgs(args...),
	}

	_ = cli.node.Broadcast(srvType, msg)
}

func (cli *_Client) CallRpcToAllUsers(srvType string, methodName string, args ...interface{}) {

	msg := &Message.AllUsersRpcReq{
		MethodName: methodName,
		Args:       Message.PackArgs(args...),
	}

	_ = cli.node.Broadcast(srvType, msg)
}

func (cli *_Client) SendMessageToUser(userID string, srvType string, message Message.IMessage) {
	if message == nil {
		return
	}

	buf, err := Message.Pack(message)
	if err != nil {
		log.Error("SendMessageToUser Message.Pack err ", err, " userID: ", userID, " MessageID: ", message.GetID())
		return
	}

	m := &Message.ForwardUserMessage{
		TargetSrv: srvType,
		UserID:    userID,
		MsgData:   buf,
	}

	_ = cli.node.Broadcast(srvType, m)
}

func (cli *_Client) SendMessageToUsers(userIDs []string, srvType string, message Message.IMessage) {
	if message == nil {
		return
	}

	buf, err := Message.Pack(message)
	if err != nil {
		log.Error("SendMessageToUsers Message.Pack err ", err, " userID: ", userIDs, " MessageID: ", message.GetID())
		return
	}
	m := &Message.ForwardUsersMessage{
		TargetSrv: srvType,
		UserIDS:   userIDs,
		MsgData:   buf,
	}

	_ = cli.node.Broadcast(srvType, m)
}

func (cli *_Client) SendMessageToAllUsers(srvType string, message Message.IMessage) {
	if message == nil {
		return
	}

	buf, err := Message.Pack(message)
	if err != nil {
		log.Error("SendMessageToAllUsers Message.Pack err ", err, " MessageID: ", message.GetID())
		return
	}

	m := &Message.ForwardAllUsersMessage{
		TargetSrv: srvType,
		MsgData:   buf,
	}

	_ = cli.node.Broadcast(srvType, m)
}
