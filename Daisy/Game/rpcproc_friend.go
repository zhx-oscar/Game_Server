package main

import (
	"Cinder/Chat/chatapi/types"
	"Cinder/Game"
	DConst "Daisy/Const"
	"Daisy/Proto"
	"encoding/json"
	log "github.com/cihub/seelog"
	"strconv"
	"time"
)

// RPC_AddFriendReq 请求加好友，被加方回调
func (proc *_RPCProc) RPC_AddFriendReq(targetID string, fromID string, reqInfo []byte) {
	log.Debugf("[RPC_AddFriendReq] 请求加好友，被加方回调 %v %v %v", targetID, fromID, reqInfo)
	// 发送申请列表
	user, err := Game.UserMgr.GetUser(targetID)
	if err != nil {
		_ = log.Error("RPC_AddFriendReq can't find targetUser: ", targetID, err)
		return
	}
	// 被加方的好友申请列表里添加这个玩家
	u, ok := user.(*_User)
	if !ok || u == nil {
		return
	}
	// 检查我的好友申请列表是否已满
	if len(u.prop.Data.Friends.ApplyList) >= DConst.FriendApplyListFull {
		return
	}
	// 检查是否已经申请过了
	for _, val := range u.prop.Data.Friends.ApplyList {
		if val.FriendID == fromID {
			log.Warnf("[RPC_AddFriendReq] 此玩家已经申请过 %v", fromID)
			return
		}
	}
	f := transByteToFriend(reqInfo)
	fl := u.prop.Data.Friends.ApplyList
	// 检查是不是已经申请过
	for i := 0; i < len(fl); i++ {
		if fl[i].FriendID == f.FriendID {
			return
		}
	}
	u.prop.SyncAddApplyFriend(f)
	// 分在线调用和上线后调用。无法区分调用时间，只有每次都同步给客户端了
}

func transByteToFriend(reqInfo []byte) *Proto.Friend {
	fromProfile := make(map[string]string)
	err := json.Unmarshal(reqInfo, &fromProfile)
	if err != nil {
		log.Error("[GetHistoryMessages] 解析fromData失败")
		return nil
	}
	head, _ := strconv.Atoi(fromProfile["Head"])
	level, _ := strconv.Atoi(fromProfile["Level"])
	isOnline, _ := strconv.ParseBool(fromProfile["IsOnline"])
	offlinetime, _ := strconv.Atoi(fromProfile["OfflineTime"])
	power, _ := strconv.Atoi(fromProfile["Power"])
	progress, _ := strconv.Atoi(fromProfile["Progress"])

	f := &Proto.Friend{
		FriendID:    fromProfile["FriendID"],
		Name:        fromProfile["Name"],
		Head:        uint32(head),
		Level:       uint32(level),
		IsOnline:    isOnline,
		OfflineTime: int64(offlinetime),
		Power:       int32(power),
		Progress:    uint32(progress),
	}

	return f
}

// RPC_AddFriendRet 加好友结果返回，请求方回调
func (proc *_RPCProc) RPC_AddFriendRet(targetID string, fromID string, ok bool) {
	log.Debugf("[RPC_AddFriendRet] 加好友结果返回，请求方回调 %v - %v", targetID, fromID, ok)
	// todo 好友同意了以后需要发红点消息吗？需要提示申请,以后再做
	if !ok {
		return
	}
	user, err := Game.UserMgr.GetUser(targetID)
	if err != nil {
		_ = log.Error("RPC_AddFriendRet can't find myself: ", targetID, err)
		return
	}
	// 被加方的好友申请列表里添加这个玩家
	u, ok := user.(*_User)
	if !ok || u == nil {
		log.Error("[RPC_AddFriendRet]找不到这个玩家")
		return
	}
	friend := u.getFriendFromChat(fromID)
	u.prop.SyncAddFriend(friend)
}

// RPC_FriendDeleted 删除好友返回，被动删除方回调
func (proc *_RPCProc) RPC_FriendDeleted(targetID string, friendID string) {
	log.Debugf("[RPC_FriendDeleted] 删除好友返回，被动删除方回调 %v - %v", targetID, friendID)

	user, err := Game.UserMgr.GetUser(targetID)
	if err != nil {
		_ = log.Error("RPC_FriendDeleted can't find myself: ", targetID, err)
		return
	}
	// 被加方的好友申请列表里添加这个玩家
	u, ok := user.(*_User)
	if !ok || u == nil {
		log.Error("[RPC_FriendDeleted]找不到这个玩家")
		return
	}
	u.prop.SyncRemoveFriend(friendID)
	// 填充每个好友的备注
	remarks := u.getMyAllFriendRemarks()
	delete(remarks, friendID)
	u.setMyAllFriendRemark(remarks)
}

// RPC_ChatUpdateUser 有玩家数据变更时，通知该服务器此玩家的好友
func (proc *_RPCProc) RPC_ChatUpdateUser(updateUserMsgJson []byte) {
	var msg types.UpdateUserMsg
	if err := json.Unmarshal(updateUserMsgJson, &msg); err != nil {
		log.Error("[RPC_ChatUpdateUser] updateUserMsgJson unmarshl failed")
		return
	}

	// 变更玩家的id数据
	friendID := msg.UserID

	log.Debug("[RPC_ChatUpdateUser] 主动数据更新 ", friendID)
	for _, val := range msg.Targets {
		// 跳过自己
		//log.Debugf("[RPC_ChatUpdateUser] 我的目标seq %d - id %v ", i, val.ID)
		if val.ID == msg.UserID {
			continue
		}

		user, err := Game.UserMgr.GetUser(val.ID)
		if err != nil {
			//log.Debugf("RPC_ChatUpdateUser can't find targetUser: %v - %v %v", val.ID, msg.UserID, err)
			continue
		}

		u, ok := user.(*_User)
		if !ok || u == nil {
			return
		}
		// 如果不是好友那么不进行好友数据变化通知
		_, ok2 := u.prop.Data.Friends.FriendList[friendID]
		if !ok2 {
			continue
		}
		switch msg.Type {
		case types.UUT_LOGIN_WITH_NICK_AND_DATA:
			// 更新我的在线状态
			log.Debug("[RPC_ChatUpdateUser] 谁更新登录")
			u.prop.SyncUpdateFriendIsOnline(friendID, true)
		case types.UUT_LOGOUT:
			// 更新我的不在线状态
			log.Debug("[RPC_ChatUpdateUser] 谁下线")
			u.prop.SyncUpdateFriendIsOnline(friendID, false)
			// 更新我的下线时间
			now := time.Now().Unix()
			u.prop.SyncUpdateFriendOfflineTime(friendID, now)
		case types.UUT_NICK:
			// 更新我的昵称
		case types.UUT_DATA:
			// todo 更新我的数据,读出什么就更新什么
			updateUserActivate(msg.Data, u, friendID)
			log.Debug("[RPC_ChatUpdateUser] 谁更新主动数据")
		}
	}
}

func updateUserActivate(profile []byte, user *_User, friendID string) {
	profile_ := map[string]string{}
	err := json.Unmarshal(profile, &profile_)
	if err != nil {
		log.Error("[updateUserActivate] 解析friendActivateData失败", user.GetID())
		return
	}

	score, ok := profile_["SeasonTeamScore"]
	if ok {
		// 好友属性系统更新level
		_score, _ := strconv.Atoi(score)
		user.prop.SyncUpdateFriendSeasonScore(friendID, uint32(_score))
	}
}
