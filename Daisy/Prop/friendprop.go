package Prop

import (
	"Cinder/Base/Message"
	"Cinder/Base/Prop"
	"Daisy/Proto"
)

// SyncAddFriendList 同步添加好友列表
func (u *RoleProp) SyncAddFriendList(friendListArray *Proto.FriendListArray) {
	u.Sync("AddFriendList", Message.PackArgs(friendListArray), false, Prop.Target_Client)
}

// AddFriendList 同步添加好友
func (u *RoleProp) AddFriendList(friendListArray *Proto.FriendListArray) {
	for key, val := range friendListArray.FriendListArray {
		u.Data.Friends.FriendList[key] = val
	}
}

// SyncAddFriend 同步添加好友
func (u *RoleProp) SyncAddFriend(friend *Proto.Friend) {
	u.Sync("AddFriend", Message.PackArgs(friend), false, Prop.Target_Client)
}

// AddFriend 添加好友
func (u *RoleProp) AddFriend(friend *Proto.Friend) {
	_, ok := u.Data.Friends.FriendList[friend.FriendID]
	if ok {
		return
	}
	u.Data.Friends.FriendList[friend.FriendID] = friend
}

// SyncRemoveFriend 同步移除好友
func (u *RoleProp) SyncRemoveFriend(friendID string) {
	u.Sync("RemoveFriend", Message.PackArgs(friendID), false, Prop.Target_Client)
}

// RemoveFriend 移除好友
func (u *RoleProp) RemoveFriend(friendID string) {
	delete(u.Data.Friends.FriendList, friendID)
}

// SyncAddApplyFriend 同步添加到好友申请列表
func (u *RoleProp) SyncAddApplyFriend(friend *Proto.Friend) {
	u.Sync("AddApplyFriend", Message.PackArgs(friend), false, Prop.Target_Client)
}

// AddApplyFriend 添加到好友申请列表
func (u *RoleProp) AddApplyFriend(friend *Proto.Friend) {
	fl := u.Data.Friends.ApplyList
	// 检查是不是已经申请过
	for i := 0; i < len(fl); i++ {
		if fl[i].FriendID == friend.FriendID {
			return
		}
	}
	u.Data.Friends.ApplyList = append(u.Data.Friends.ApplyList, friend)
}

// SyncRemoveApplyFriend 同步移除好友申请列表
func (u *RoleProp) SyncRemoveApplyFriend(friendID string) {
	u.Sync("RemoveApplyFriend", Message.PackArgs(friendID), false, Prop.Target_Client)
}

// RemoveApplyFriend 移除好友申请列表
func (u *RoleProp) RemoveApplyFriend(friendID string) {
	applys := u.Data.Friends.ApplyList
	for i := 0; i < len(applys); i++ {
		if applys[i].FriendID == friendID {
			u.Data.Friends.ApplyList = append(applys[:i], applys[i+1:]...)
			break
		}
	}
}

// SyncUpdateFriendRemark 同步更新好友备注(设置备注、离线时间更新、角色升级等变化)
func (u *RoleProp) SyncUpdateFriendRemark(friendID string, remark string) {
	u.Sync("UpdateFriendRemark", Message.PackArgs(friendID, remark), false, Prop.Target_Client)
}

// UpdateFriendRemark 更新好友备注
func (u *RoleProp) UpdateFriendRemark(friendID string, remark string) {
	f, ok := u.Data.Friends.FriendList[friendID]
	if !ok {
		return
	}
	f.Remark = remark
}

// SyncUpdateFriendIsOnline 同步更新好友在线状态
func (u *RoleProp) SyncUpdateFriendIsOnline(friendID string, online bool) {
	u.Sync("UpdateFriendIsOnline", Message.PackArgs(friendID, online), false, Prop.Target_Client)
}

// UpdateFriendIsOnline 更新好友好友在线状态
func (u *RoleProp) UpdateFriendIsOnline(friendID string, online bool) {
	f, ok := u.Data.Friends.FriendList[friendID]
	if !ok {
		return
	}
	f.IsOnline = online
}

// SyncUpdateFriendOfflineTime 同步更新好友离线时间
func (u *RoleProp) SyncUpdateFriendOfflineTime(friendID string, offlineTime int64) {
	u.Sync("UpdateFriendOfflineTime", Message.PackArgs(friendID, offlineTime), false, Prop.Target_Client)
}

// UpdateFriendOfflineTime 更新好友离线时间
func (u *RoleProp) UpdateFriendOfflineTime(friendID string, offlineTime int64) {
	f, ok := u.Data.Friends.FriendList[friendID]
	if !ok {
		return
	}
	f.OfflineTime = offlineTime
}

// SyncUpdateFriendLevel 同步更新好友等级
func (u *RoleProp) SyncUpdateFriendLevel(friendID string, level uint32) {
	u.Sync("UpdateFriendLevel", Message.PackArgs(friendID, level), false, Prop.Target_Client)
}

// UpdateFriendLevel 更新好友离线等级
func (u *RoleProp) UpdateFriendLevel(friendID string, level uint32) {
	f, ok := u.Data.Friends.FriendList[friendID]
	if !ok {
		return
	}
	f.Level = level
}

// SyncUpdateFriendSeasonScore 更新好友赛季积分
func (u *RoleProp) SyncUpdateFriendSeasonScore(friendID string, score uint32) {
	u.Sync("UpdateFriendSeasonScore", Message.PackArgs(friendID, score), false, Prop.Target_Client)
}
func (u *RoleProp) UpdateFriendSeasonScore(friendID string, score uint32) {
	u.Data.Friends.FriendList[friendID].SeasonTeamScore = score
}
