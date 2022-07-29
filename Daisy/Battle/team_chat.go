package main

import (
	cConst "Cinder/Base/Const"
	"Cinder/Chat/chatapi"
	"Cinder/Space"
	"Daisy/Const"
)


//initTeamChatChannel 初始化队伍聊天频道
func (team *_Team) initTeamChatChannel() {
	tc := team.getTeamChatChannel()
	if tc == nil {
		team.Errorf("[initTeamChatChannel] 创建队伍频道失败 teamID:%v", team.GetID())
		return
	}

	// 成员全部加入对应频道
	team.TraversalActor(func(actor Space.IActor) {
		role := actor.(*_Role)
		team.teamChatChannelAddMember(role)
	})
}

//teamChatChannelAddMember 队友聊天频道加入新成员
func (team *_Team) teamChatChannelAddMember(role *_Role) {
	if role == nil {
		return
	}
	user := role.GetOwnerUser()
	if user == nil {
		team.Error("[teamChatChannelAddMember] get user failed from role")
		return
	}
	tid := team.getTeamChatChannelKey()
	user.Rpc(cConst.Game, "RPC_TeamChatChannelAddMember", tid)
}

//teamChatChannelDelMember 队友聊天频道删除成员
func (team *_Team) teamChatChannelDelMember(role *_Role) {
	if role == nil {
		return
	}
	user := role.GetOwnerUser()
	if user == nil {
		team.Error("[teamChatChannelDelMember] get user failed from role")
		return
	}
	tid := team.getTeamChatChannelKey()
	user.Rpc(cConst.Game, "RPC_TeamChatChannelDelMember", tid)
}

//getTeamChatChannelKey 获取队伍聊天频道key
func (team *_Team) getTeamChatChannelKey() string {
	return Const.ChatGroupTypeTeam + ":" + team.GetID()
}

//getTeamChatChannel 获取队伍聊天频道
func (team *_Team) getTeamChatChannel() chatapi.IGroup {
	tc, _ := chatapi.GetOrCreateGroup(team.getTeamChatChannelKey())
	return tc
}

//teamChatDestroy 队伍解散的时候调用  队伍成员数量为0
func (team *_Team) teamChatDestroy() {
	if len(team.prop.Data.Base.Members) != 0 {
		return
	}
	//成员全部移除对应频道
	team.TraversalActor(func(actor Space.IActor) {
		role := actor.(*_Role)
		team.teamChatChannelDelMember(role)
	})

	chatapi.DeleteGroup(team.getTeamChatChannelKey())
}
