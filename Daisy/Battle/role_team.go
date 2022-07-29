package main

import (
	CConst "Cinder/Base/Const"
	"Cinder/Base/Core"
	"Cinder/Base/Prop"
	"Cinder/Base/Util"
	"Daisy/Const"
	"Daisy/DB"
	"Daisy/DHDB"
	"Daisy/ErrorCode"
	"Daisy/NotifyCode"
	Prop2 "Daisy/Prop"
	"Daisy/Proto"
)

func (user *_User) RPC_AddRoleInvite(teamID, instigator string) int32 {
	if user.role == nil {
		user.Error("RPC_AddRoleInvite Failed TeamID:", teamID, "instigator:", instigator)
		return ErrorCode.RoleIsNil
	}

	return user.role.AddRoleInvite(teamID, instigator)
}

func (r *_Role) AddRoleInvite(teamID, instigator string) int32 {
	if len(r.prop.Data.Invites) >= Const.MaxInvites {
		r.Error("AddRoleInvite InviteOverLimit", teamID)
		return ErrorCode.InviteOverLimit
	}

	r.prop.SyncAddRoleInviteInfo(teamID, instigator)
	r.Info("AddRoleInvite Success TeamID:", teamID, "instigator:", instigator)
	return ErrorCode.Success
}

func (user *_User) RPC_RemoveRoleInvite(teamID string, clear bool) int32 {
	if user.role == nil {
		user.Error("RPC_RemoveRoleInvite Failed TeamID:", teamID)
		return ErrorCode.RoleIsNil
	}
	return user.role.RemoveRoleInvite(teamID, clear)
}

func (r *_Role) RemoveRoleInvite(teamID string, clear bool) int32 {
	if clear {
		r.prop.SyncClearRoleInviteInfo()
	} else {
		r.prop.SyncRemoveRoleInviteInfo(teamID)
	}

	r.Info("RemoveRoleInvite Success", teamID, clear)
	return ErrorCode.Success
}

func (user *_User) RPC_SendInvite(userID string) int32 {
	if user.role == nil {
		user.Error("RPC_SendInvite Failed Target UserID:", userID)
		return ErrorCode.RoleIsNil
	}
	return user.role.SendInvite(userID)
}

func (r *_Role) SendInvite(beInvitedUserID string) int32 {
	team := r.GetSpace().(*_Team)

	exist, err := DB.GetApply2InviteUtil().InviteIsExist(beInvitedUserID, team.GetID())
	if err != nil {
		r.Error("SendInvite DB Err:", err)
		return ErrorCode.DBOpErr
	}
	if exist {
		r.Error("SendInvite repeat invite ", beInvitedUserID)
		return ErrorCode.RepeatInvite
	}

	beInvitedRole, err := DHDB.GetRoleCache(beInvitedUserID)
	if err != nil {
		r.Error("SendInvite GetRoleError", err)
		return ErrorCode.GetRoleError
	}

	if _, ok := team.prop.Data.Base.Members[beInvitedUserID]; ok {
		team.SendNotify(uint32(NotifyCode.JoinTeamSuccessWithName), r.GetID(), beInvitedRole.Base.Name, false, false)

		r.Error("SendInvite AlreadInTeam ", beInvitedUserID)
		return ErrorCode.AlreadInTeam
	}

	if len(team.prop.Data.Base.Members) >= Const.TeamMaxMemberNum {
		team.SendNotify(uint32(NotifyCode.SelfTramFull), r.GetID(), "", false, false)

		r.Error("SendInvite TeamFull ", beInvitedUserID)
		return ErrorCode.TeamFull
	}

	beInvitedTeamBase, err := DHDB.GetTeamBase(beInvitedRole.Base.TeamID)
	if err != nil {
		r.Error("SendInvite GetTeamBase Err:", err)
		return ErrorCode.GetTeamError
	}
	if len(beInvitedTeamBase.Members) != 1 {
		r.Error("SendInvite beInvited user AlreadInTeam")
		return ErrorCode.AlreadInTeam
	}

	err = DB.GetApply2InviteUtil().AddInvite(beInvitedUserID, team.GetID(), r.GetID())
	if err != nil {
		r.Error("SendInvite DB Err:", err)
		return ErrorCode.DBOpErr
	}

	Core.Inst.CallRpcToUser(beInvitedUserID, CConst.Space, "RPC_AddRoleInvite", team.GetID(), r.GetID())

	r.Info("SendInvite Success Target UserID:", beInvitedUserID)
	return ErrorCode.Success
}

func (user *_User) RPC_AgreeInvite(teamID string) int32 {
	if user.role == nil {
		user.Error("RPC_AgreeInvite Failed TeamID:", teamID)
		return ErrorCode.RoleIsNil
	}
	return user.role.AgreeInvite(teamID)
}

func (r *_Role) AgreeInvite(inviteTeamID string) int32 {
	inviteInfo, ok := r.prop.Data.Invites[inviteTeamID]
	if !ok {
		r.Error("AgreeInvite NotInInvites TeamID:", inviteTeamID)
		return ErrorCode.NotInInvites
	}

	r.prop.SyncRemoveRoleInviteInfo(inviteTeamID)
	if err := DB.GetApply2InviteUtil().RemoveInvite(r.GetID(), inviteTeamID); err != nil {
		r.Error("AgreeInvite DB err ", err)
		return ErrorCode.DBOpErr
	}

	go func() {
		Transfer(r.GetSpace().(*_Team).GetID(), inviteTeamID, r.GetID(), inviteInfo.Instigator, Const.TransferReason_Invite)
	}()

	r.Info("AgreeInvite Success TeamID:", inviteTeamID)
	return ErrorCode.Success
}

func (user *_User) RPC_RefuseInvite(teamID string) int32 {
	if user.role == nil {
		user.Error("RPC_RefuseInvite TeamID:", teamID)
		return ErrorCode.RoleIsNil
	}
	return user.role.RefuseInvite(teamID)
}

func (r *_Role) RefuseInvite(inviteTeamID string) int32 {
	info, ok := r.prop.Data.Invites[inviteTeamID]
	if !ok {
		r.Error("RefuseInvite NotInInvites TeamID:", inviteTeamID)
		return ErrorCode.NotInInvites
	}

	if srvID, err := DB.TeamUtil().GetSrvID(inviteTeamID); err == nil {
		Core.Inst.RpcByID(srvID, "RPC_SendNotify", inviteTeamID, uint32(NotifyCode.RefuseFormTeamInvitatoin), info.Instigator, r.prop.Data.Base.Name, false, false)
	}

	r.prop.SyncRemoveRoleInviteInfo(inviteTeamID)
	DB.GetApply2InviteUtil().RemoveInvite(r.GetID(), inviteTeamID)

	r.Info("RefuseInvite Success TeamID:", inviteTeamID)
	return ErrorCode.Success
}

func (user *_User) RPC_SendApply(applyTeamID, message string) (int32, int32) {
	if len(message) > 20 {
		user.Error("RPC_SendApply message too long")
		return ErrorCode.ArgsWrong, 0
	}

	if applyTeamID == user.GetSpace().GetID() {
		user.GetSpace().(*_Team).SendNotify(uint32(NotifyCode.JoinTeamSuccess), user.GetID(), "", false, false)

		user.Error("RPC_SendApply AlreadInTeam ", applyTeamID)
		return ErrorCode.AlreadInTeam, 0
	}

	if user.role == nil {
		user.Error("RPC_SendApply Failed TeamID:", applyTeamID)
		return ErrorCode.RoleIsNil, 0
	}

	return user.role.SendApply(applyTeamID, message)
}

func (r *_Role) SendApply(applyTeamID, message string) (int32, int32) {
	team := r.GetSpace().(*_Team)
	if team.GetID() == applyTeamID {
		r.Error("SendApply ArgsWrong ", applyTeamID)
		return ErrorCode.ArgsWrong, 0
	}

	if len(team.prop.Data.Base.Members) != 1 {
		team.SendNotify(uint32(NotifyCode.JoinTeamSuccess), r.GetID(), "", false, false)

		r.Error("SendApply AlreadInTeam ", applyTeamID)
		return ErrorCode.AlreadInTeam, 0
	}

	exist, err := DB.GetApply2InviteUtil().ApplyIsExist(r.GetID(), applyTeamID)
	if err != nil {
		r.Error("SendApply DB Err:", err)
		return ErrorCode.DBOpErr, 0
	}
	if exist {
		r.Error("SendApply repeat apply ", applyTeamID)
		return ErrorCode.RepeatApply, 0
	}

	targetTeamBase, err := DHDB.GetTeamBase(applyTeamID)
	if err != nil {
		r.Error("SendApply DB Err:", err)
		return ErrorCode.GetTeamError, 0
	}
	if len(targetTeamBase.Members) >= Const.TeamMaxMemberNum {
		r.Error("SendApply TeamFull TeamID:", applyTeamID)
		return ErrorCode.TeamFull, 0
	}

	if r.checkAutoJoin(targetTeamBase) {
		go func() {
			Transfer(team.GetID(), applyTeamID, r.GetID(), "", Const.TransferReason_AutoJoin)
		}()

		return ErrorCode.Success, 1
	}

	err = DB.GetApply2InviteUtil().AddApply(r.GetID(), applyTeamID, team.GetID(), message)
	if err != nil {
		r.Error("SendApply DB Err:", err)
		return ErrorCode.DBOpErr, 0
	}

	if srvID, err := DB.TeamUtil().GetSrvID(applyTeamID); err == nil {
		Core.Inst.RpcByID(srvID, "RPC_AddTeamApply", applyTeamID, r.GetID(), message)
	}

	r.Info("SendApply Success TeamID:", applyTeamID)
	return ErrorCode.Success, 0
}

func (user *_User) RPC_AgreeApply(userID string) int32 {
	if user.role == nil {
		user.Error("RPC_AgreeApply Failed UserID:", userID)
		return ErrorCode.RoleIsNil
	}
	return user.role.AgreeApply(userID)
}

func (r *_Role) AgreeApply(userID string) int32 {
	team := r.GetSpace().(*_Team)
	status := team.prop.Data.Base.Members[r.GetID()].Status
	if status != Const.TeamStatus_LEADER {
		r.Error("AgreeApply NoPermission ", userID)
		return ErrorCode.NoPermission
	}

	if _, ok := team.prop.Data.Applys[userID]; !ok {
		r.Error("AgreeApply NotInApplys ", userID)
		return ErrorCode.NotInApplys
	}

	team.prop.SyncRemoveApplyInfo(userID)
	if err := DB.GetApply2InviteUtil().RemoveApply(userID, team.GetID()); err != nil {
		r.Error("AgreeApply DB err ", err)
		return ErrorCode.DBOpErr
	}

	beApplyedRole, err := DHDB.GetRoleCache(userID)
	if err != nil {
		r.Error("AgreeApply GetRoleError ", err)
		return ErrorCode.GetRoleError
	}

	go func() {
		Transfer(beApplyedRole.Base.TeamID, team.GetID(), userID, r.GetID(), Const.TransferReason_Apply)
	}()

	r.Info("AgreeApply Success UserID:", userID)
	return ErrorCode.Success
}

func (user *_User) RPC_RefuseApply(userID string) int32 {
	if user.role == nil {
		user.Error("RPC_RefuseApply Failed UserID:", userID)
		return ErrorCode.RoleIsNil
	}
	return user.role.RefuseApply(userID)
}

func (r *_Role) RefuseApply(userID string) int32 {
	team := r.GetSpace().(*_Team)
	status := team.prop.Data.Base.Members[r.GetID()].Status
	if status != Const.TeamStatus_LEADER {
		r.Error("RefuseApply NoPermission ", userID)
		return ErrorCode.NoPermission
	}

	if _, ok := team.prop.Data.Applys[userID]; !ok {
		r.Error("RefuseApply NotInApplys ", userID)
		return ErrorCode.NotInApplys
	}

	team.prop.SyncRemoveApplyInfo(userID)
	DB.GetApply2InviteUtil().RemoveApply(userID, team.GetID())

	applyRole, err := DHDB.GetRoleCache(userID)
	if err != nil {
		r.Error("RefuseApply GetRoleError ", err)
		return ErrorCode.GetRoleError
	}

	if srvID, err := DB.TeamUtil().GetSrvID(applyRole.Base.TeamID); err == nil {
		Core.Inst.RpcByID(srvID, "RPC_SendNotify", applyRole.Base.TeamID, uint32(NotifyCode.RefuseFormTeamApply), userID, r.prop.Data.Base.Name, false, false)
	}

	r.Info("RefuseApply Success UserID:", userID)
	return ErrorCode.Success
}

func (user *_User) RPC_QuitTeam() int32 {
	if user.role == nil {
		user.Error("RPC_QuitTeam Failed")
		return ErrorCode.RoleIsNil
	}
	return user.role.QuitTeam()
}

func (r *_Role) QuitTeam() int32 {
	team := r.GetSpace().(*_Team)
	if len(team.prop.Data.Base.Members) == 1 {
		r.Error("QuitTeam error in single team")
		return ErrorCode.Failure
	}

	go func() {
		Transfer(team.GetID(), Util.GetGUID(), r.GetID(), "", Const.TransferReason_Quit)
	}()

	r.Info("QuitTeam Success")
	return ErrorCode.Success
}

func (user *_User) RPC_KickMember(userID string) int32 {
	if user.GetID() == userID {
		user.Error("RPC_KickMember arg error")
		return ErrorCode.ArgsWrong
	}

	if user.role == nil {
		user.Error("RPC_KickMember Failed UserID:", userID)
		return ErrorCode.RoleIsNil
	}

	return user.role.KickMember(userID)
}

func (r *_Role) KickMember(userID string) int32 {
	team := r.GetSpace().(*_Team)
	if _, ok := team.prop.Data.Base.Members[userID]; !ok {
		r.Error("KickMember userID not in team", userID)
		return ErrorCode.ArgsWrong
	}

	status := team.prop.Data.Base.Members[r.GetID()].Status
	if status != Const.TeamStatus_LEADER {
		r.Error("KickMember NoPermission")
		return ErrorCode.NoPermission
	}

	go func() {
		Transfer(team.GetID(), Util.GetGUID(), userID, "", Const.TransferReason_Kick)
	}()

	r.Info("KickMember Success UserID:", userID)
	return ErrorCode.Success
}

func (user *_User) RPC_ModifyTeamName(name string) int32 {
	if Const.UTF8Width(name) > 10 {
		user.Error("RPC_ModifyTeamName name too long")
		return ErrorCode.ArgsWrong
	}

	if user.role == nil {
		user.Error("RPC_ModifyTeamName Failed Name:", name)
		return ErrorCode.RoleIsNil
	}
	return user.role.ModifyTeamName(name)
}

func (r *_Role) ModifyTeamName(name string) int32 {
	team := r.GetSpace().(*_Team)
	status := team.prop.Data.Base.Members[r.GetID()].Status
	if status != Const.TeamStatus_LEADER {
		r.Error("ModifyTeamName NoPermission")
		return ErrorCode.NoPermission
	}

	team.prop.SyncSetName(name)
	team.FlushToCache()

	r.Info("ModifyTeamName Success Name:", name)
	return ErrorCode.Success
}

func (user *_User) RPC_GetRecruitments() int32 {
	if user.role == nil {
		user.Error("RPC_GetRecruitments Failed")
		return ErrorCode.RoleIsNil
	}
	return user.role.GetRecruitments()
}

func (r *_Role) GetRecruitments() int32 {
	team := r.GetSpace().(*_Team)
	team.GetRecruitments(r)

	r.Info("GetRecruitments Success")
	return ErrorCode.Success
}

func (r *_Role) OnGetRecruitments(res *Proto.Recruitments) {
	iu := r.GetOwnerUser()
	if iu == nil {
		return
	}

	//过滤已申请的队伍-
	vals, err := DB.GetApply2InviteUtil().GetApplysInRole(r.GetID())
	if err != nil {
		r.Error("GetRecruitments DB Err:", err)
	} else {
		applyedTeamIDs := make(map[string]bool)
		for i := 0; i < len(vals); i++ {
			applyedTeamIDs[vals[i].TeamID] = true
		}

		for i := 0; i < len(res.Data); i++ {
			if _, ok := applyedTeamIDs[res.Data[i].TeamID]; ok {
				res.Data[i].Applied = true
			}
		}
	}

	r.Info("OnGetRecruitments Success")
	iu.Rpc(CConst.Agent, "RPC_SendRecruitments", res)
}

func (user *_User) RPC_GetTeamByUID(uid uint64) (int32, *Proto.Recruitments) {
	if user.role == nil {
		user.Error("RPC_GetTeamByUID Failed UID:", uid)
		return ErrorCode.RoleIsNil, nil
	}
	return user.role.GetTeamByUID(uid)
}

func (r *_Role) GetTeamByUID(uid uint64) (int32, *Proto.Recruitments) {
	team := r.GetSpace().(*_Team)
	if team.prop.Data.Base.UID == uid {
		r.Error("GetTeamByUID FindSelf")
		return ErrorCode.FindSelf, nil
	}

	res, err := team.GetTeamByUID(uid)
	if err != nil {
		r.Error("GetTeamByUID FindNone")
		return ErrorCode.FindNone, nil
	}

	r.Info("GetTeamByUID Success UID:", uid)
	return ErrorCode.Success, res
}

func (user *_User) RPC_GetInvites() (int32, *Proto.Inviters) {
	if user.role == nil {
		user.Error("RPC_GetInvites Failed")
		return ErrorCode.RoleIsNil, nil
	}
	return user.role.GetInvites()
}

func (r *_Role) GetInvites() (int32, *Proto.Inviters) {
	team := r.GetSpace().(*_Team)
	invites, err := team.GetInvites(r.prop.Data.Invites)
	if err != nil {
		r.Error("GetInvites err", err)
		return ErrorCode.GetInvitesErr, nil
	}

	r.Info("GetInvites Success")
	return ErrorCode.Success, invites
}

func (user *_User) RPC_GetApplys() (int32, *Proto.Applyers) {
	if user.role == nil {
		user.Error("RPC_GetApplys Failed")
		return ErrorCode.RoleIsNil, nil
	}
	return user.role.GetApplys()
}

func (r *_Role) GetApplys() (int32, *Proto.Applyers) {
	team := r.GetSpace().(*_Team)
	status := team.prop.Data.Base.Members[r.GetID()].Status
	if status != Const.TeamStatus_LEADER {
		r.Error("ModifyTeamBoard NoPermission")
		return ErrorCode.NoPermission, nil
	}

	applys, err := team.GetApplys(team.prop.Data.Applys)
	if err != nil {
		r.Error("GetApplys err ", err)
		return ErrorCode.GetApplysErr, nil
	}

	r.Info("GetApplys Success")
	return ErrorCode.Success, applys
}

func (user *_User) RPC_PublishTeam(board string, needAuth bool, autoJoinIdx uint32) int32 {
	if user.role == nil {
		user.Error("RPC_PublishTeam Failed")
		return ErrorCode.RoleIsNil
	}
	return user.role.PublishTeam(board, needAuth, autoJoinIdx)
}

func (r *_Role) PublishTeam(board string, needAuth bool, autoJoinIdx uint32) int32 {
	team := r.GetSpace().(*_Team)
	status := team.prop.Data.Base.Members[r.GetID()].Status
	if status != Const.TeamStatus_LEADER {
		r.Error("ModifyTeamBoard NoPermission")
		return ErrorCode.NoPermission
	}

	if Const.UTF8Width(board) > 40 {
		r.Error("PublishTeam board too long")
		return ErrorCode.BoardTooLong
	}

	team.prop.SyncPublishTeam(board, needAuth, autoJoinIdx)
	team.FlushToDB()
	team.FlushToCache()

	r.Infof("PublishTeam Success board:%s needAuth:%v autoJoinIdx:%d", board, needAuth, autoJoinIdx)
	return ErrorCode.Success
}

func (r *_Role) checkAutoJoin(teamBase *Proto.TeamBase) bool {
	if teamBase.NeedAuth {
		return false
	}

	teamAvgScore := int32(0)
	members := make([]string, len(teamBase.Members), len(teamBase.Members))
	types := make([]string, len(teamBase.Members), len(teamBase.Members))
	i := 0
	for key := range teamBase.Members {
		members[i] = key
		types[i] = Prop2.RolePropType
		i++
	}

	results, err := Prop.GetBatchCacheProp(types, members)
	if err != nil {
		r.Errorf("GetBatchCacheProp err:%s", err)
		return false
	}

	for i = 0; i < len(results); i++ {
		cache := results[i].(*Proto.RoleCache)
		if fightingBD, ok := cache.BuildMap[cache.FightingBuildID]; !ok {
			r.Errorf("role上阵BD不存在 FightingBuildID:%s", cache.FightingBuildID)
			return false
		} else {
			teamAvgScore += fightingBD.FightAttr.TotalScore
		}
	}
	teamAvgScore = teamAvgScore / int32(len(members))

	autoJoinDiff := map[uint32]float32{0: 0.05, 1: 0.1, 2: 0.15}
	if diff, ok := autoJoinDiff[teamBase.AutoJoinIdx]; !ok {
		r.Errorf("AutoJoinIdx:%d out range", teamBase.AutoJoinIdx)
		return false
	} else {
		realDiff := r.prop.Data.BuildMap[r.prop.Data.FightingBuildID].FightAttr.TotalScore - teamAvgScore
		if realDiff < 0 {
			realDiff *= -1
		}

		if float32(realDiff) <= float32(teamAvgScore)*diff {
			return true
		}
	}

	return false
}
