package main

import (
	"Cinder/Base/Core"
	"Cinder/Base/Message"

	log "github.com/cihub/seelog"
)

type _SrvMsgProc struct{}

func (proc *_SrvMsgProc) MessageProc(srcAddr string, message Message.IMessage) {
	switch message.GetID() {
	case Message.ID_User_Login_Ret:
		msg := message.(*Message.UserLoginRet)
		<-userMgr.onLoginRet(msg.UserID, msg.PropType, msg.PropDef, msg.UserData)
	case Message.ID_User_Logout_Req:
		msg := message.(*Message.UserLogoutReq)
		userMgr.Logout(msg.UserID)
	case Message.ID_User_Destroy_Req:
		msg := message.(*Message.UserDestroyReq)
		userMgr.DestroyUser(msg.UserID)
		if err := Core.Inst.Send(srcAddr, &Message.UserDestroyRet{UserID: msg.UserID}); err != nil {
			log.Error("Send UserDestroyRet err ", err, " UserID: ", msg.UserID)
		}
	case Message.ID_Client_Rpc_Ret:
		msg := message.(*Message.ClientRpcRet)

		user, err := userMgr.GetAgentUser(msg.UserID)
		if err != nil {
			log.Error("MessageProc ClientRpcRet couldn't find user ", msg.UserID)
			return
		}

		if err = user.SendToClient(msg); err != nil {
			log.Errorf("SendToClient err %s userID %s", err, msg.UserID)
		}

	case Message.ID_Space_Broadcast_To_Client:
		proc.forwardMessageToClients(message.(*Message.SpaceBroadcastToClient))

	case Message.ID_User_Destroy_Ret:
		userMgr.onUserDestroyRet(message.(*Message.UserDestroyRet))
	}
}

func (proc *_SrvMsgProc) forwardMessageToClients(broadcastMsg *Message.SpaceBroadcastToClient) {
	innerMsg, err := Message.Unpack(broadcastMsg.MsgData)
	if err != nil {
		log.Error("forwardMessageToClients Unpack innerMsg err ", err)
		return
	}

	var user *_User
	for _, userID := range broadcastMsg.UserList {
		user, err = userMgr.GetAgentUser(userID)
		if err != nil {
			log.Errorf("forwardMessageToClients GetAgentUser err %s userID %s", err, userID)
			continue
		}

		if userID != broadcastMsg.ExceptUserID {
			if err = user.SendToClient(innerMsg); err != nil {
				log.Errorf("forwardMessageToClients SendToClient err %s userID %s", err, userID)
			}
		}
	}
}
