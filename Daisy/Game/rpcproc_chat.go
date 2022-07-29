package main

import (
	"Cinder/Base/Const"
	"Cinder/Chat/chatapi/types"
	"Cinder/Game"
	"Daisy/Proto"
	"encoding/json"
	log "github.com/cihub/seelog"
	"time"
)

// RPC_ChatRecvP2PMessage 接收私聊消息，接收者回调
func (proc *_RPCProc) RPC_ChatRecvP2PMessage(targetID string, targetData []byte, fromID string, fromNick string, fromData []byte, msgContent []byte) {
	fromProfile := make(map[string]string)
	err := json.Unmarshal(fromData, &fromProfile)
	if err != nil {
		log.Error("[RPC_ChatRecvGroupMessageV2] 解析fromData失败")
		return
	}

	user, err := Game.UserMgr.GetUser(targetID)
	if err != nil {
		_ = log.Error("RPC_ChatRecvGroupMessage can't find targetUser: ", targetID, err)
		return
	}

	ret := &Proto.ChatMessage{
		Name:       fromNick,
		From:       fromID,
		FromHead:   fromProfile["FromHead"],
		TitleID:    fromProfile["TitleID"],
		To:         targetID,
		MsgType:    Proto.ChatMessage_MsgTypePrivate,
		Level:      fromProfile["Level"],
		SendTime:   time.Now().Unix(),
		MsgContent: string(msgContent),
	}
	user.Rpc(Const.Agent, "RPC_StoCMessage", ret)
	log.Debug("[RPC_ChatRecvP2PMessage] 接收私聊消息，接收者回调", ret)
}

// RPC_ChatRecvGroupMessageV2 接收群聊消息，接收者回调
func (proc *_RPCProc) RPC_ChatRecvGroupMessageV2(groupID string, targetsJson []byte, fromID string, fromNick string, fromData []byte, msgContent []byte) {
	log.Debug("[RPC_ChatRecvGroupMessageV2] 接收群聊消息，接收者回调")

	var hh []types.Target
	if err := json.Unmarshal(targetsJson, &hh); err != nil {
		log.Error("[RPC_ChatRecvGroupMessageV2] 解析targetJson失败")
		return
	}

	for i, val := range hh {
		log.Debug("[RPC_ChatRecvGroupMessageV2] 遍历一个群里所有玩家, 当前玩家编号, id", i, val.ID)
		user, err := Game.UserMgr.GetUser(val.ID)
		if err != nil {
			_ = log.Error("RPC_ChatRecvGroupMessage can't find targetUser: ", val.ID, err)
			return
		}

		fromProfile := make(map[string]string)
		err = json.Unmarshal(fromData, &fromProfile)
		if err != nil {
			log.Error("[RPC_ChatRecvGroupMessageV2] 解析fromData失败")
			return
		}
		message := &Proto.ChatMessage{
			Name:       fromNick,
			From:       fromID,
			FromHead:   fromProfile["FromHead"],
			TitleID:    fromProfile["TitleID"],
			To:         groupID,
			MsgType:    Proto.ChatMessage_MsgTypeGroup,
			Level:      fromProfile["Level"],
			SendTime:   time.Now().Unix(),
			MsgContent: string(msgContent),
		}

		user.Rpc(Const.Agent, "RPC_StoCMessage", message)
	}
}
