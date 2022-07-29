package main

import (
	"Cinder/Chat/chatapi"
	"Cinder/Chat/chatapi/types"
	DConst "Daisy/Const"
	"Daisy/DB"
	"Daisy/DHDB"
	"Daisy/ErrorCode"
	"Daisy/Proto"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"strconv"
)

func (u *_User) friendOnline() {
	u.refreshFriendList()
}

// ------------------------------------客户端RPC调用---------------------------------------------------------

// RPC_AgreeFriendReq 同意添加好友
func (u *_User) RPC_AgreeFriendReq(friendID string) int32 {
	u.Debug("RPC_AgreeFriendReq roleId ", u.GetID(), " ", friendID)
	return u.agreeFriend(friendID)
}

// RPC_RefuseFriendReq 拒绝添加好友
func (u *_User) RPC_RefuseFriendReq(friendID string) int32 {
	u.Debug("RPC_RefuseFriendReq roleId ", u.GetID(), " ", friendID)
	return u.refuseFriend(friendID)
}

// RPC_OneKeyAgreeFriendReq 请求一键同意申请列表
func (u *_User) RPC_OneKeyAgreeFriendReq() int32 {
	u.Debug("RPC_OneKeyAgreeFriendReq roleId ", u.GetID(), " ")
	return u.batchAgreeFriend()
}

// RPC_OneKeyRefuseFriendReq 请求一键同意申请列表
func (u *_User) RPC_OneKeyRefuseFriendReq() int32 {
	u.Debug("RPC_OneKeyAgreeFriendReq roleId ", u.GetID(), " ")
	return u.batchRefuseFriend()
}

// RPC_ApplyAddFriendReq 申请添加好友
func (u *_User) RPC_ApplyAddFriendReq(targetID string) int32 {
	u.Debug("RPC_ApplyAddFriendReq roleId ", u.GetID(), " ", targetID)
	return u.ApplyAddFriend(targetID)
}

// RPC_DeleteFriend 删除好友
func (u *_User) RPC_DeleteFriend(friendID string) int32 {
	u.Debug("RPC_DeleteFriend roleId ", u.GetID(), " ", friendID)
	return u.RemoveFriend(friendID)
}

// RPC_FriendRemarkReq 好友备注
func (u *_User) RPC_FriendRemarkReq(friendID string, remark string) int32 {
	u.Debug("RPC_FriendRemarkReq roleId ", u.GetID(), " ", friendID, remark)
	return u.friendRemarkReq(friendID, remark)
}

// RPC_RecommendedListReq 请求推荐好友列表
func (u *_User) RPC_RecommendedListReq() *Proto.FriendArray {
	u.Debug("RPC_ShowRecommendedListReq roleId ", u.GetID(), " ")
	return u.recommendedListReq()
}

// RPC_FindFriendReq 请求搜索好友
func (u *_User) RPC_FindFriendReq(uid uint64) (int32, *Proto.Friend) {
	u.Debug("RPC_FindFriend roleId ", u.GetID(), " ")
	return u.findFriendReq(uid)
}

// ----------------------------------------------------工具函数----------------------------------

// refreshFriendList 刷新好友列表
func (u *_User) refreshFriendList() {
	if u.chatUser == nil {
		u.Error("[refreshFriendList] role's chatUser is nil")
		return
	}
	friendInfos := u.chatUser.GetFriendList()
	// 从聊天服拿到的好友列表
	friendList := make(map[string]*Proto.Friend, 0)

	for _, val := range friendInfos {
		// 拿到转化后的数据，刷新数据结构体
		f := u.transFriendInfoToFriend(val)
		friendList[f.FriendID] = f
	}

	friendListArray := &Proto.FriendListArray{
		FriendListArray: make(map[string]*Proto.Friend, 0),
	}
	// 填充每个好友的备注
	remarks := u.getMyAllFriendRemarks()

	// 拼接备注，属性刷新
	for key, val := range remarks {
		friend, ok := friendList[key]
		if !ok {
			delete(remarks, key)
			u.setMyAllFriendRemark(remarks)
			// todo remarks 删除这一条, 这是多余的备注
		} else {
			friend.Remark = val
		}
	}
	friendListArray.FriendListArray = friendList

	u.prop.SyncAddFriendList(friendListArray)
}

// 下线期间好友数据发生变化，我需要知道
func (u *_User) transFriendInfoToFriend(friendInfos types.FriendInfo) *Proto.Friend {
	f := &Proto.Friend{
		FriendID:    friendInfos.ID,
		Name:        friendInfos.Nick,
		IsOnline:    friendInfos.IsOnline,
		OfflineTime: friendInfos.OfflineTime.Unix(),
	}

	// 反解析Data
	fromProfile := make(map[string]string)
	err := json.Unmarshal(friendInfos.Data, &fromProfile)
	if err != nil {
		u.Error("[GetHistoryMessages] 解析fromData失败")
	}
	head, _ := strconv.Atoi(fromProfile["FromHead"])
	f.Head = uint32(head)
	level, _ := strconv.Atoi(fromProfile["Level"])
	f.Level = uint32(level)
	seasonTeamScore, _ := strconv.Atoi(fromProfile["SeasonTeamScore"])
	f.SeasonTeamScore = uint32(seasonTeamScore)

	return f
}

// getMyAllFriendRemarks 返回我对所有好友的备注
func (u *_User) getMyAllFriendRemarks() map[string]string {
	// 填充每个好友的备注
	if u.chatUser == nil {
		u.Error("[getMyAllFriendRemarks] role's chatUser is nil")
		return nil
	}
	passiveData, err := u.chatUser.GetPassiveData()
	if err != nil {
		u.Error("[getMyAllFriendRemarks] getPassiveData failed form chat ", err)
		return nil
	}

	fromProfile := make(map[string]string)
	if len(passiveData) == 0 {
		return fromProfile
	}
	err = json.Unmarshal(passiveData, &fromProfile)
	if err != nil {
		u.Error("[getMyAllFriendRemarks] ummarshal PassiveData failed form chat ", err)
		return nil
	}
	return fromProfile
}

func (u *_User) setMyAllFriendRemark(remarks map[string]string) int32 {
	newPassiveData, err := json.Marshal(remarks)
	if err != nil {
		u.Error("[setMyAllFriendRemark] marshal remarks failed ", err)
		return ErrorCode.FriendDBFailed
	}

	if u.chatUser == nil {
		u.Error("[getFriendFromChat] role's chatUser is nil")
		return ErrorCode.FriendUnknowErr
	}

	err = u.chatUser.SetPassiveData(newPassiveData)
	if err != nil {
		u.Error("[setMyAllFriendRemark] set passiveData to chat failed ", err)
		return ErrorCode.FriendDBFailed
	}
	return ErrorCode.Success
}

// AddFriend 同意好友申请
func (u *_User) agreeFriend(fromId string) int32 {
	// todo 其实这个不用加，能发过来证明他不是我的好友。好友是双向的。后期再删除
	// 是否已经是我的好友，如果是，返回错误码
	_, k := u.prop.Data.Friends.FriendList[fromId]
	if k {
		u.prop.SyncRemoveApplyFriend(fromId)
		return ErrorCode.FriendAddRepeated
	}

	// 检查我的好友列表是否已满
	if len(u.prop.Data.Friends.FriendList) >= DConst.FriendNumMax {
		return ErrorCode.FriendListFull
	}

	if u.chatUser == nil {
		u.Error("[agreeFriend] role's chatUser is nil")
		return ErrorCode.FriendUnknowErr
	}
	// 向聊天服发出确认
	err := u.chatUser.ReplyAddFriendReq(fromId, true)
	if err != nil {
		u.Error("[ApplyAddFriend] AddFriendReq fail", err)
		// 对方好友列表已满
		if err == chatapi.ErrPeerReachedMaxFriendCount {
			return ErrorCode.FriendTargetFriendListFull
		}
		return ErrorCode.FriendDBFailed
	}
	// 聊天服确认对方好友数量是否达到上限

	// 移除这个申请从申请列表里
	u.prop.SyncRemoveApplyFriend(fromId)
	// 将对方加入到我的好友列表里
	friend := u.getFriendFromChat(fromId)
	u.prop.SyncAddFriend(friend)

	return ErrorCode.Success
}

func (u *_User) getFriendFromChat(friendID string) *Proto.Friend {
	if u.chatUser == nil {
		u.Error("[getFriendFromChat] role's chatUser is nil")
		return nil
	}
	friendInfos := u.chatUser.GetFriendList()
	friendInfo := u.findFriendFromReqInfos(friendInfos, friendID)
	if friendInfo == nil {
		u.Error("[getFriendFromChat] can't find friendUser")
		return nil
	}

	friend := u.transFriendInfoToFriend(*friendInfo)
	return friend
}

func (u *_User) findFriendFromReqInfos(friendInfos []chatapi.FriendInfo, targetID string) *chatapi.FriendInfo {
	if len(friendInfos) == 0 {
		return nil
	}
	for _, val := range friendInfos {
		if val.ID == targetID {
			return &val
		}
	}
	return nil
}

// refuseFriend 拒绝好友申请
func (u *_User) refuseFriend(fromId string) int32 {
	if u.chatUser == nil {
		u.Error("[refuseFriend] role's chatUser is nil")
		return ErrorCode.FriendUnknowErr
	}
	// 向聊天服发出确认
	err := u.chatUser.ReplyAddFriendReq(fromId, false)
	if err != nil {
		u.Error("[ApplyAddFriend] AddFriendReq fail", err)
		return ErrorCode.FriendDBFailed
	}

	// 移除这个申请从申请列表里
	u.prop.SyncRemoveApplyFriend(fromId)
	return ErrorCode.Success
}

// batchAgreeFriend 一键同意
func (u *_User) batchAgreeFriend() int32 {
	applyList := u.prop.Data.Friends.ApplyList
	// 检查我的好友列表是否已满
	num := 0
	if len(u.prop.Data.Friends.FriendList) >= DConst.FriendNumMax {
		return ErrorCode.FriendListFull
	}
	for i := len(applyList) - 1; i >= 0; i-- {
		err := u.agreeFriend(applyList[i].FriendID)
		if err == ErrorCode.FriendListFull {
			break
		}
		if err == ErrorCode.FriendTargetFriendListFull {
			num++
		}
	}
	if num == len(applyList) {
		return ErrorCode.FriendAllTargetFriendListFull
	}

	// 申请条目中所有玩家的好友数量都达到上限时。提示 对方好友数量已经达到上限。
	return ErrorCode.Success
}

// batchRefuseFriend 一键拒绝
func (u *_User) batchRefuseFriend() int32 {
	applyList := u.prop.Data.Friends.ApplyList
	for i := 0; i < len(applyList); i++ {
		u.refuseFriend(applyList[i].FriendID)
		// 不对错误码做处理
	}
	return ErrorCode.Success
}

// ApplyAddFriend 申请加好友
func (u *_User) ApplyAddFriend(friendId string) int32 {
	// 是否已经是我的好友，如果是，返回错误码
	_, ok := u.prop.Data.Friends.FriendList[friendId]
	if ok {
		return ErrorCode.FriendAddRepeated
	}
	// 检查我的好友列表是否已满
	if len(u.prop.Data.Friends.FriendList) >= DConst.FriendNumMax {
		return ErrorCode.FriendListFull
	}

	// 聊天服 检查对方申请列表是否已满
	// 聊天服 检查对方申请列表是否已经有这个申请了。
	// 聊天服发送申请
	if u.chatUser == nil {
		u.Error("[ApplyAddFriend] role's ChatUser is nil")
		return ErrorCode.FriendUnknowErr
	}
	err := u.chatUser.AddFriendReq(friendId, u.compressInfo())
	if err != nil {
		u.Error("[ApplyAddFriend] AddFriendReq fail", err)
		// 对方申请列表已经满
		if err == chatapi.ErrPeerReachedMaxAddFriendReq {
			return ErrorCode.FriendTargetApplyListFull
		}
		return ErrorCode.FriendDBFailed
	}

	return ErrorCode.Success
}

// compressInfo 压缩信息  此信息为好友申请列表显示的信息
func (u *_User) compressInfo() []byte {
	profile_ := map[string]string{
		"FriendID":    u.GetID(),
		"Name":        u.prop.Data.Base.Name,
		"Head":        strconv.Itoa(int(u.prop.Data.Base.Head)),
		"Level":       strconv.Itoa(int(u.prop.Data.Base.Level)),
		"IsOnline":    strconv.FormatBool(true), // 发送添加好友申请时，一定在线
		"OfflineTime": strconv.Itoa(int(u.prop.Data.Base.LastLogoutTime)),
	}
	u.fillupFriendApplyData(&profile_)

	profile, err := json.Marshal(profile_)
	if err != nil {
		u.Errorf("[sendInfoAboutChat] error ", err)
		return nil
	}
	return profile
}

// fillupFriendApplyData 填充关卡、战力。 用于申请列表数据，不从邮件服去取主动数据，是因为这个数据不用即时刷新
func (u *_User) fillupFriendApplyData(profile_ *map[string]string) {
	roleData, err := DHDB.GetRoleCache(u.GetID())
	if err != nil {
		u.Error("[fillupFriendApplyProgress] get role cache failed", err)
		return
	}
	// 填充关卡数据
	roleBase := roleData.Base
	if roleBase != nil {
		(*profile_)["Progress"] = strconv.Itoa(int(roleBase.RaidProgress))
	} else {
		u.Error("[fillupFriendApplyData] get roleCache base failed")
	}

	// 填充战力数据
	roleBuildMap := roleData.BuildMap
	if roleBuildMap != nil {
		roleBuildData, ok := roleBuildMap[roleData.FightingBuildID]
		if ok {
			if roleBuildData.FightAttr != nil {
				(*profile_)["Power"] = strconv.Itoa(int(roleBuildData.FightAttr.TotalScore))
			} else {
				u.Error("[fillupFriendApplyData] roleCache FightAtttr is nil")
			}
		} else {
			u.Error("[fillupFriendApplyData] get roleCache BuildData failed, fightingBuildID is ", roleData.FightingBuildID)
		}
	} else {
		u.Error("[fillupFriendApplyData] roleCache buildMap is nil")
	}
}

// RemoveFriend 删除好友
func (u *_User) RemoveFriend(friendId string) int32 {
	// 查找好友是否在列表中
	_, ok := u.prop.Data.Friends.FriendList[friendId]
	if !ok {
		return ErrorCode.FriendNotExit
	}
	if u.chatUser == nil {
		u.Error("[RemoveFriend] role's chatUser is nil")
		return ErrorCode.FriendUnknowErr
	}
	// 判断聊天服是否移除成功
	err := u.chatUser.DeleteFriend(friendId)
	if err != nil {
		u.Error("[RemoveFriend] removeFriend fail")
		return ErrorCode.FriendDBFailed
	}
	// 聊天服移除,通知自己客户端移除
	u.prop.SyncRemoveFriend(friendId)
	// 填充每个好友的备注
	remarks := u.getMyAllFriendRemarks()
	delete(remarks, friendId)
	u.setMyAllFriendRemark(remarks)
	// 邮件服移除是双方都移除

	return ErrorCode.Success
}

// friendRemarkReq 给好友设置备注
func (u *_User) friendRemarkReq(friendID string, remark string) int32 {
	// 查找好友是否在列表中
	_, ok := u.prop.Data.Friends.FriendList[friendID]
	if !ok {
		return ErrorCode.FriendNotExit
	}
	num := DConst.UTF8Width(remark)
	if num > DConst.FriendRemarkLengthLimit {
		return ErrorCode.FriendRemarkOver
	}

	// 聊天服设置备注是否成功
	ec := u.setMyFriendRemark(friendID, remark)
	if ec != ErrorCode.Success {
		return ec
	}

	u.prop.SyncUpdateFriendRemark(friendID, remark)
	return ErrorCode.Success
}

func (u *_User) setMyFriendRemark(friendID string, remark string) int32 {
	remarks := u.getMyAllFriendRemarks()
	if remarks == nil {
		return ErrorCode.FriendDBFailed
	}
	remarks[friendID] = remark

	ec := u.setMyAllFriendRemark(remarks)
	return ec
}

// 获得推荐好友列表返回给客户端
func (u *_User) recommendedListReq() *Proto.FriendArray {
	fl := &Proto.FriendArray{}
	// 判断缓存里推荐好友人数是否大于10，不是就从数据库中拉取。
	if len(u.prop.Data.Friends.RecommendedList) < DConst.FriendRecommendOnce {
		// 从数据库里拉60条
		recommendFromDB := u.getRecommendfromDB()
		for i := 0; i < len(recommendFromDB); i++ {
			u.prop.Data.Friends.RecommendedList = append(u.prop.Data.Friends.RecommendedList, recommendFromDB[i])
		}
	}

	// 裁剪推荐数据, 裁剪掉我的朋友
	for _, v := range u.prop.Data.Friends.FriendList {
		for i := 0; i < len(u.prop.Data.Friends.RecommendedList); {
			if v.FriendID == u.prop.Data.Friends.RecommendedList[i].FriendID {
				u.prop.Data.Friends.RecommendedList = append(u.prop.Data.Friends.RecommendedList[:i], u.prop.Data.Friends.RecommendedList[i+1:]...)
			} else {
				i++
			}
		}
	}

	recommendLen := len(u.prop.Data.Friends.RecommendedList)
	// 判断缓存里数据是否大于10，大于则扣10个，小于就发全部
	var min int
	min = 0
	if recommendLen < DConst.FriendRecommendOnce {
		min = recommendLen
	} else {
		min = DConst.FriendRecommendOnce
	}

	limit := recommendLen - min
	if limit < 0 {
		u.Error("[recommendedListReq] 获取推荐好友数量出错")
		return nil
	}
	for i := recommendLen - 1; i >= limit; i-- {
		fl.FriendArray = append(fl.FriendArray, u.prop.Data.Friends.RecommendedList[i])
		// todo 属性系统里同步删除这个,每次从尾部删。或者批量删除
	}
	u.prop.Data.Friends.RecommendedList = append(u.prop.Data.Friends.RecommendedList[:limit], u.prop.Data.Friends.RecommendedList[recommendLen:]...)

	return fl
}

// getRecommendfromDB 从数据库里拉取数据
func (u *_User) getRecommendfromDB() []*Proto.Friend {
	// 推荐好友的标准是 在线或者离线，但离线时间不超过3天,不能是我,不能是我的好友
	//tmp := int64(24 * DConst.FriendRecommendOfflineTimeLimit * 3600)
	//now := time.Now().Unix()
	//Subtime := now - tmp
	// 得到3天前的时间
	//limitTime := time.Now().Sub(tmp)
	param := bson.M{"json_data.base.online": true}
	//param1 := bson.M{"json_data.base.lastlogouttime": bson.M{"$gte": Subtime}}

	//roles, err := DB.GetRoleHallUitl().Find(param, DConst.FriendRecommendMax, 0)
	pipeline := make([]bson.M, 0)
	pipeline = append(pipeline, bson.M{"$match": param})

	roles, err := DB.GetRoleHallUitl().Aggregate(pipeline)
	if err != nil {
		return nil
	}
	fs := u.transBatchProtoWrapToFriend(roles)

	return fs
}

func (u *_User) transBatchProtoWrapToFriend(wraps []*DB.RoleProtoWrap) []*Proto.Friend {
	fs := make([]*Proto.Friend, 0)
	for i := 0; i < len(wraps); i++ {
		// 推荐列表裁剪掉自己
		if wraps[i].ID.Hex() == u.GetID() {
			continue
		}
		f := u.transProtoWrapToFriend(wraps[i])
		fs = append(fs, f)
	}
	return fs
}

func (u *_User) transProtoWrapToFriend(wrap *DB.RoleProtoWrap) *Proto.Friend {
	f := &Proto.Friend{
		FriendID:    wrap.ID.Hex(),
		Name:        wrap.Role.Base.Name,
		Head:        wrap.Role.Base.Head,
		Level:       wrap.Role.Base.Level,
		IsOnline:    wrap.Role.Base.Online,
		OfflineTime: wrap.Role.Base.LastLogoutTime,
		Progress:    wrap.Role.Base.RaidProgress,
	}
	u.fillUpFriendPower(&f, wrap.Role.BuildMap, wrap.Role.FightingBuildID)

	return f
}

// fillUpFriendPower 填充推荐好友结构体里的战力值
func (u *_User) fillUpFriendPower(f **Proto.Friend, mapdata map[string]*Proto.BuildData, fid string) {
	if mapdata != nil {
		build, ok := mapdata[fid]
		if ok {
			fightAttr := build.FightAttr
			if fightAttr != nil {
				(*f).Power = fightAttr.TotalScore
			} else {
				u.Error("[fillUpFriendPower] 推荐好友列表 buildMap[FightingBuildID].fightAttr 为空", fid)
			}
		} else {
			u.Error("[fillUpFriendPower] 推荐好友列表 buildMap里没有找到 fightBuildID", fid)
		}
	} else {
		u.Error("[fillUpFriendPower] 推荐好友列表 buildMap为空")
	}
}

// findFriendReq 用uid查找对应的玩家
func (u *_User) findFriendReq(uid uint64) (int32, *Proto.Friend) {
	if u.prop.Data.Base.UID == uid {
		u.Error("[findFriendReq] find self uid is ", uid)
		return ErrorCode.FriendFindSelf, nil
	}
	f, err := u.findFriendByUid(uid)
	if err != nil {
		u.Error("[findFriendReq] find friend failed uid -err ", uid, err)
		return ErrorCode.FriendFindNoExit, nil
	}
	// 判断是不是我的好友
	_, ok := u.prop.Data.Friends.FriendList[f.FriendID]
	if ok {
		f.IsFriend = true
	}
	return ErrorCode.Success, f
}

func (u *_User) findFriendByUid(uid uint64) (*Proto.Friend, error) {
	wrap, err := DB.GetRoleHallUitl().FindUID(uid)
	if err != nil {
		return nil, err
	}
	f := u.transProtoWrapToFriend(wrap)
	return f, nil
}
