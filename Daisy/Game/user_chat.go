package main

import (
	"Cinder/Base/Const"
	"Cinder/Chat/chatapi"
	DConst "Daisy/Const"
	"Daisy/DHDB"
	"Daisy/ErrorCode"
	"Daisy/Proto"
	"encoding/json"
	"strconv"
	"time"
)

type MessageBody struct {
	Type    DConst.MessageType
	Content interface{}
}

type RequestSkillMsg struct {
	RoleId string
	Uid    string
}

// chatOnline 聊天系统上线流程
func (u *_User) chatOnline() {
	u.Debug("User chatOnline")
	u.chatUserLogin()
}

// chatOffline 聊天系统下线流程
func (u *_User) chatOffline() {
	u.Debug("User chatOffline")
	_ = chatapi.Logout(u.GetID())
}

// initWorldChannel 初始化世界聊天
func (u *_User) initWorldChannel() {
	w, err := chatapi.GetOrCreateGroup(DConst.ChatGroupTypeWorld)
	if w == nil || err != nil {
		u.Error("[initWorldChannel] 创建世界频道失败", err)
		return
	}
	err = w.AddIntoGroup(u.GetID())
	if err != nil {
		u.Error("[initWorldChannel] 用户加入世界频道失败, id 是 ", u.GetID())
		return
	}

	u.prop.SyncAddToChatChannel(DConst.ChatGroupTypeWorld)
}

// ------------------------------------battle服RPC调用------------------------------------------------

// RPC_UpdateChatUserActivateData 刷新我的主动数据推送给我的朋友
func (u *_User) RPC_UpdateChatUserActivateData() int32 {
	return u.updateChatUserActivateData()
}

// todo 这个还没有测过

// RPC_RequestSkillMsg 技能乞求消息
func (u *_User) RPC_RequestSkillMsg(roleId, uid string) int32 {
	requestSkillMsg := RequestSkillMsg{
		RoleId: roleId,
		Uid:    uid,
	}

	var err error
	var subJsonBody []byte
	var jsonBody []byte
	subJsonBody, err = json.Marshal(requestSkillMsg)
	if err != nil {
		return ErrorCode.MarshalJsonErr
	}

	messageBody := MessageBody{
		Type:    DConst.ChatTypeRequestSkill,
		Content: string(subJsonBody),
	}

	jsonBody, err = json.Marshal(messageBody)
	if err != nil {
		return ErrorCode.MarshalJsonErr
	}

	ret := &Proto.ChatMessage{
		Name:       u.prop.Data.Base.Name,
		From:       u.GetID(),
		FromHead:   strconv.Itoa(int(u.prop.Data.Base.Head)),
		TitleID:    strconv.Itoa(int(u.prop.Data.Title.TitleID)),
		To:         u.prop.Data.Base.TeamID,
		MsgType:    Proto.ChatMessage_MsgTypeGroup,
		Level:      strconv.Itoa(int(u.prop.Data.Base.Level)),
		SendTime:   time.Now().Unix(),
		MsgContent: string(jsonBody),
	}
	if u.chatUser == nil {
		// 报错
		// todo 为空了为什么还要发送？
		u.Rpc(Const.Agent, "RPC_StoCMessage", ret)
		return ErrorCode.Timeout
	}

	err = u.chatUser.SendGroupMessage(u.getTeamChatChannelKey(), jsonBody)
	if err != nil {
		u.Error("[CtoSMessage] chatapi.UserSendGroupMessage err != nil", err)
	}
	return ErrorCode.Success
}

// todo 需要有返回值吗？

// RPC_TeamChatChannelAddMember 队友聊天频道加入新成员
func (u *_User) RPC_TeamChatChannelAddMember(tid string) {
	tc, _ := chatapi.GetOrCreateGroup(tid)
	if tc == nil {
		u.Error("[RPC_TeamChatChannelAddMember] 创建队伍失败 ", tid)
		return
	}

	// 加入频道
	err := tc.AddIntoGroup(u.GetID())
	if err != nil {
		u.Errorf("[RPC_TeamChatChannelAddMember] 玩家加入队伍频道 %v失败, %v ", tid, err)
	}
	u.prop.SyncAddToChatChannel(tid)
}

// RPC_TeamChatChannelDelMember 队友聊天频道删除新成员
func (u *_User) RPC_TeamChatChannelDelMember(tid string) {
	tc, _ := chatapi.GetOrCreateGroup(tid)
	if tc == nil {
		u.Error("[RPC_TeamChatChannelAddMember] 创建队伍失败 ", tid)
		return
	}
	// 踢出频道
	err := tc.KickFromGroup(u.GetID())
	if err != nil {
		u.Errorf("[RPC_TeamChatChannelDelMember] 玩家退出队伍频道 %v失败, %v ", tid, err)
		return
	}
	u.prop.SyncDelFromChatChannel(tid)
}

// ---------------------------------------------RPC调用 & 和客户端通信函数 ----------------------------------------

func (u *_User) RPC_CToSMessage(msgType Proto.ChatMessage_ChatMsgType, to string, msgContent string) int32 {
	return u.CtoSMessage(msgType, to, msgContent)
}

func (u *_User) RPC_HistoryMessage() (message *Proto.ChannelChatHistoryMessage, int32 int32) {
	u.Debug("RPC_HistoryMessage ", u.GetID())

	return u.sendAllChatHistory(), ErrorCode.Success
}

// SendSystemMessage 系统发送消息时推送给客户端
func (u *_User) SendSystemMessage(msgContent string) {
	u.Rpc(Const.Agent, "RPC_SendSystemMessage", msgContent)
}

// -------------------------------------------- 工具函数 ----------------------------------------------------------

// updateChatUserActivateData 刷新聊天服主动数据（聊天、好友列表里需要即时刷新的数据）
func (u *_User) updateChatUserActivateData() int32 {
	// todo 所有主动数据封装一个结构体
	// 主动数据改变没办法改变一个。因为邮件服是直接覆盖的方式去改的。所以必须是整体。

	profile_ := map[string]string{
		"UserID":   u.GetID(),
		"Name":     u.prop.Data.Base.Name,
		"Level":    strconv.Itoa(int(u.prop.Data.Base.Level)),
		"FromHead": strconv.Itoa(int(u.prop.Data.Base.Head)),
		"TitleID":  strconv.Itoa(int(u.prop.Data.Title.TitleID)),
	}

	teamData, err := DHDB.GetTeamPart(u.GetID())
	if err == nil && teamData != nil {
		if teamData.TeamInfo != nil {
			if teamData.TeamInfo.SeasonInfo != nil {
				profile_["SeasonTeamScore"] = strconv.Itoa(int(teamData.TeamInfo.SeasonInfo.TeamScore))
			}
		}
	} else {
		u.Error("[refreshActivateDataToMyFriends] GetTeamPart failed", err)
	}

	profile, err := json.Marshal(profile_)
	if err != nil {
		u.Error("[sendInfoAboutChat] activate data marshal error ", err)
		return ErrorCode.Failure
	}
	if u.chatUser == nil {
		u.Error("[refreshActiveDataToMyFriends] role's chatUser is nil")
		return ErrorCode.Failure
	}

	if u.chatUser != nil {
		err = u.chatUser.SetActiveData(profile)
	}
	if err != nil {
		return ErrorCode.Failure
	}
	return ErrorCode.Success
}

// chatUserLogin 聊天服登录
func (u *_User) chatUserLogin() bool {
	u.Debug("User chatUserLogin")
	teamID := u.prop.Data.Base.TeamID
	// 新号上线，里面数据都是空。为了确保数据有效。必须在每个主动数据修改的时候同步到聊天服。
	// 老号上线， name、Fromhead都有吗？ 这个时候有做属性系统的Serilize吗？ 有
	profile_ := map[string]string{
		"UserID":   u.GetID(),
		"Name":     u.prop.Data.Base.Name,
		"Level":    strconv.Itoa(int(u.prop.Data.Base.Level)),
		"FromHead": strconv.Itoa(int(u.prop.Data.Base.Head)),
		"TitleID":  strconv.Itoa(int(u.prop.Data.Title.TitleID)),
	}

	if teamID != "" {
		teamData, err := DHDB.GetTeamPart(u.GetID())
		if err != nil || teamData == nil {
			u.Error("[refreshActivateDataToMyFriends] GetTeamPart failed", err)
		} else {
			if teamData.TeamInfo != nil {
				if teamData.TeamInfo.SeasonInfo != nil {
					profile_["SeasonTeamScore"] = strconv.Itoa(int(teamData.TeamInfo.SeasonInfo.TeamScore))
				}
			}
		}
	}

	profile, err := json.Marshal(profile_)
	if err != nil {
		u.Error("[chatUserLogin] Marshal origin profile error ", err)
		return false
	}

	u.chatUser, err = chatapi.Login(u.GetID(), u.prop.Data.Base.Name, profile)
	if err != nil || u.chatUser == nil {
		u.Error("[chatUserLogin] 登录chat 出错 ", err)
		return false
	}
	return true
}

// CtoSMessage 给聊天服务器发消息
func (u *_User) CtoSMessage(msgType Proto.ChatMessage_ChatMsgType, to string, msgContent string) int32 {
	messageBody := MessageBody{
		Type:    DConst.ChatTypeNormal,
		Content: msgContent,
	}

	// 客户端发送上来的消息做长度检测
	if DConst.UTF8Width(msgContent) > DConst.ChatMsgMaxCount {
		return ErrorCode.ChatMsgContent
	}

	jsonBody, err := json.Marshal(messageBody)
	if err != nil {
		return ErrorCode.MarshalJsonErr
	}

	ret := &Proto.ChatMessage{
		Name:       u.prop.Data.Base.Name,
		From:       u.GetID(),
		FromHead:   strconv.Itoa(int(u.prop.Data.Base.Head)),
		TitleID:    strconv.Itoa(int(u.prop.Data.Title.TitleID)),
		To:         to,
		MsgType:    msgType,
		Level:      strconv.Itoa(int(u.prop.Data.Base.Level)),
		SendTime:   time.Now().Unix(),
		MsgContent: string(jsonBody),
	}
	if u.chatUser == nil {
		// 报错
		//u.Rpc(Const.Agent, "RPC_StoCMessage", ret)
		return ErrorCode.Timeout
	}

	switch msgType {
	case Proto.ChatMessage_MsgTypePrivate:
		err = u.chatUser.SendMessage(to, jsonBody)
		if err != nil {
			u.Error("[CtoSMessage] chatapi.UserSendMessage err != nil", err)
			return ErrorCode.ChatPrivateSendFail
		}
		u.Rpc(Const.Agent, "RPC_StoCMessage", ret)
	case Proto.ChatMessage_MsgTypeGroup:
		err = u.chatUser.SendGroupMessage(to, jsonBody)
		if err != nil {
			u.Error("[CtoSMessage] chatapi.UserSendGroupMessage err != nil", err)
		}
	}

	u.Debug("[role] CToSMessage msg: ", jsonBody)
	return ErrorCode.Success
}

// sendAllChatHistory 拉取所有频道历史记录 推送固定数量的世界、工会频道的历史消息  (在上线时使用)
func (u *_User) sendAllChatHistory() *Proto.ChannelChatHistoryMessage {
	if u == nil {
		u.Error("[sendAllChatHistory] r.chatUser 获取失败")
		return nil
	}

	channelHisMsgs := &Proto.ChannelChatHistoryMessage{
		ChannelHisMsg: make(map[string]*Proto.ChatHistoryMessage, 0),
	}

	u.Debug("[lcc] game服获取聊天频道", u.prop.Data.Chat.ChatChannels)
	for groupid := range u.prop.Data.Chat.ChatChannels {
		historyMsgs := &Proto.ChatHistoryMessage{
			HistoryMsg: make([]*Proto.ChatMessage, 0),
		}

		group, _ := chatapi.GetOrCreateGroup(groupid)
		arr := group.GetHistoryMessages(u.getChatHistoryMessageCount(groupid))

		for _, val := range arr {
			if val == nil {
				continue
			}

			fromProfile := make(map[string]string)
			err := json.Unmarshal(val.FromData, &fromProfile)
			if err != nil {
				u.Error("[GetHistoryMessages] 解析fromData失败")
				continue
			}

			hMsg := &Proto.ChatMessage{
				From:       val.From,
				Name:       val.FromNick,
				FromHead:   fromProfile["FromHead"],
				TitleID:    fromProfile["TitleID"],
				MsgType:    Proto.ChatMessage_MsgTypeGroup,
				SendTime:   val.SendTime,
				MsgContent: string(val.MsgContent),
				To:         groupid,
				Level:      strconv.Itoa(int(u.prop.Data.Base.Level)),
			}
			historyMsgs.HistoryMsg = append(historyMsgs.HistoryMsg, hMsg)
		}
		channelHisMsgs.ChannelHisMsg[groupid] = historyMsgs
	}

	//个人的离线消息处理
	OfflineMessage := u.chatUser.GetOfflineMessage()
	if OfflineMessage != nil {
		//私聊处理
		for _, val := range OfflineMessage.P2PMessages {
			fromProfile := make(map[string]string)
			err := json.Unmarshal(val.FromData, &fromProfile)
			if err != nil {
				u.Error("[GetOfflineMessage] 解析fromData失败")
				continue
			}

			hMsg := &Proto.ChatMessage{
				From:       val.From,
				Name:       val.FromNick,
				FromHead:   fromProfile["FromHead"],
				TitleID:    fromProfile["TitleID"],
				MsgType:    Proto.ChatMessage_MsgTypeGroup,
				SendTime:   val.SendTime,
				MsgContent: string(val.MsgContent),
				To:         u.GetID(),
				Level:      strconv.Itoa(int(u.prop.Data.Base.Level)),
			}

			channelHisMsgs.OfflinePrivateMsg = append(channelHisMsgs.OfflinePrivateMsg, hMsg)
		}
	}

	u.Debug("[role_chat] SendHistoryMessage", channelHisMsgs)
	return channelHisMsgs
}

// getChatHistoryMessageCount 获取聊天历史消息限制数量
func (u *_User) getChatHistoryMessageCount(groupID string) uint16 {
	// 私聊和队伍聊天
	switch groupID {
	case u.getTeamChatChannelKey():
		return DConst.ChatTeamHistoryMessageCount
	default:
		return DConst.GroupHistoryMessageCount
	}
}

func (u *_User) getTeamChatChannelKey() string {
	return DConst.ChatGroupTypeTeam + ":" + u.prop.Data.Base.TeamID
}
