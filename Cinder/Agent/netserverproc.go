package main

import (
	"Cinder/Base/Const"
	"Cinder/Base/Message"
	"Cinder/Base/Net"
	"errors"
	"strings"
	"time"

	log "github.com/cihub/seelog"
)

type NetServerProc struct{}

var (
	errUnexpectedMessage = errors.New("unexpected message")
)

func newNetServerProc() *NetServerProc {
	return &NetServerProc{}
}

func (p *NetServerProc) OnSessConnected(sess Net.ISess) {
	log.Info("Sess Connected ", sess.GetData())
}

func (p *NetServerProc) OnSessClosed(sess Net.ISess) {
	log.Info("Sess Closed ", sess.GetData())

	userID, ok := sess.GetData().(string)
	if !ok {
		return
	}

	userMgr.PendLogout(userID)
}

func (p *NetServerProc) OnSessMessageHandle(sess Net.ISess, msgNo uint32, message Message.IMessage) {

	if !sess.IsValidate() && message.GetID() != Message.ID_Client_Validate_Req {
		log.Errorf("OnSessMessageHandle get unexpected message id %d Sess %s", message.GetID(), sess.GetData())
		sess.Close()
		return
	}

	if message.GetID() != Message.ID_Client_Validate_Req {
		if err := p.checkMsgValidate(sess, msgNo); err != nil {
			log.Errorf("OnSessMessageHandle message id %d check err %s Sess %s", message.GetID(), err, sess.GetData())
			return
		}
	}

	switch message.GetID() {
	case Message.ID_Client_Validate_Req:
		p.onClientValidateReq(sess, message.(*Message.ClientValidateReq))
	case Message.ID_Client_Rpc_Req:
		p.onClientRpcReq(sess, message.(*Message.ClientRpcReq))
	case Message.ID_Forward_User_Message:
		p.onClientForward(sess, message.(*Message.ForwardUserMessage))
	case Message.ID_Heart_Beat:
		p.onHeartbeat(sess, message.(*Message.HeartBeat))
	default:
		log.Warn("OnSessMessageHandle ignore message id ", message.GetID())
	}
}

func (p *NetServerProc) checkMsgValidate(sess Net.ISess, msgNo uint32) error {

	if sess.GetData() == nil {
		sess.Close()
		return errUnexpectedMessage
	}

	userID := sess.GetData().(string)
	if err := userMgr.CheckRecvMsgValidate(userID, msgNo); err != nil {
		sess.Close()
		return err
	}

	return nil
}

func (p *NetServerProc) onClientRpcReq(sess Net.ISess, message *Message.ClientRpcReq) {
	if !strings.HasPrefix(message.MethodName, "RPC_") {
		log.Warnf("RPC Method Name invalid %s Sess %s", message.MethodName, sess.GetData())
		return
	}

	message.UserID = sess.GetData().(string)

	user, err := userMgr.GetAgentUser(message.UserID)
	if err != nil {
		log.Error("RPC couldn't find user %s", message.UserID)
		return
	}

	if err = user.SendToPeerServer(message.SrvType, message); err != nil {
		log.Errorf("RPC SendToPeerServer %s err %s", message.SrvType, err)
	}
}

func (p *NetServerProc) onClientValidateReq(sess Net.ISess, message *Message.ClientValidateReq) {
	if err := userMgr.Login(sess, message); err != nil {
		log.Errorf("Client login err %s userID %s", err, message.ID)
		return
	}

	log.Info("Client validate success userID ", message.ID)
}

func (p *NetServerProc) onClientForward(sess Net.ISess, message *Message.ForwardUserMessage) {

	if sess.GetData() == nil {
		log.Warn("onClientForward sess data nil")
		return
	}

	userID := sess.GetData().(string)
	user, err := userMgr.GetUser(userID)
	if err != nil {
		log.Errorf("onClientForward GetUser err %s userID %s", err, userID)
		return
	}

	message.UserID = userID
	if err = user.SendToPeerServer(message.TargetSrv, message); err != nil {
		log.Errorf("onClientForward SendToPeerServer %s err %s userID %s", message.TargetSrv, err, userID)
	}
}

func (p *NetServerProc) onHeartbeat(sess Net.ISess, msg *Message.HeartBeat) {
	if sess.GetData() == nil {
		log.Warn("onHeartbeat sess data nil")
		return
	}

	userID := sess.GetData().(string)
	user, err := userMgr.GetUser(userID)
	if err != nil {
		log.Errorf("onHeartbeat GetUser err %s userID %s", err, userID)
		return
	}

	_, err = user.GetPeerServerID(Const.Space)
	if err != nil {
		msg.ServerTime = time.Now().UnixNano()
		if err = user.SendToClient(msg); err != nil {
			log.Errorf("onHeartbeat SendToClient err %s userID %s", err, userID)
		}
	} else {
		if err = user.SendToPeerUser(Const.Space, msg); err != nil {
			log.Errorf("onHeartbeat SendToPeerUser space err %s userID %s", err, userID)
		}
	}
}
