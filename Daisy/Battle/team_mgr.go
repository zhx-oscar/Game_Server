package main

import (
	"Cinder/Base/Core"
	"Cinder/Space"
	"Daisy/Const"
	"Daisy/DB"
	"Daisy/ErrorCode"
	"Daisy/Prop"
	"math"
)

func (team *_Team) GetActorByUserID(userID string) Space.IActor {
	var ia Space.IActor
	team.TraversalActor(func(actor Space.IActor) {
		if actor.GetID() == userID {
			ia = actor
		}
	})
	return ia
}

func (team *_Team) MemberOnline(userID string) int32 {
	// 取消下线定时器
	if team.offlineTimer != nil {
		team.offlineTimer.Stop()
		team.offlineTimer = nil
	}

	return ErrorCode.Success
}

func (team *_Team) MemberOffline(userID string) int32 {
	return ErrorCode.Success
}

func (team *_Team) JoinTeam(userID string) int32 {
	if userID == "" {
		team.Error("JoinTeam UserIDInvalid")
		return ErrorCode.UserIDInvalid
	}
	if _, ok := team.prop.Data.Base.Members[userID]; ok {
		team.Errorf("JoinTeam AlreadInTeam user:%s", userID)
		return ErrorCode.AlreadInTeam
	}
	if len(team.prop.Data.Base.Members) >= Const.TeamMaxMemberNum {
		team.Errorf("JoinTeam TeamFull user:%s", userID)
		return ErrorCode.TeamFull
	}

	team.RemovePlaceHolder(userID)

	actorID, err := team.AddActor(RoleActorType, userID, userID, nil, nil)
	if err != nil {
		team.Errorf("JoinTeam add actor id:%s err:%s", userID, err)
		return ErrorCode.JoinTeamError
	}

	ia, err := team.GetActor(actorID)
	if err != nil {
		team.Errorf("JoinTeam get actor id:%s err:%s", userID, err)
		return ErrorCode.GetRoleError
	}

	ia.GetProp().(*Prop.RoleProp).SyncSetTeamID(team.GetID())
	status := Const.TeamStatus_NORMAL
	if len(team.prop.Data.Base.Members) == 0 {
		status = Const.TeamStatus_LEADER
	}
	team.prop.AddTeamMember(userID, status)

	//队伍满员 清空申请列表
	if len(team.prop.Data.Base.Members) == Const.TeamMaxMemberNum {
		DB.GetApply2InviteUtil().RemoveAllApplysInTeam(team.GetID())
		team.prop.SyncClearApplyInfo()
	}

	team.FlushToDB()
	team.FlushToCache()
	ia.FlushToDB()
	ia.FlushToCache()

	team.onMemberChange(userID, true)

	team.Infof("JoinTeam success user:%s", userID)
	return ErrorCode.Success
}

func (team *_Team) QuitTeam(userID string) int32 {
	if userID == "" {
		team.Error("QuitTeam UserIDInvalid")
		return ErrorCode.UserIDInvalid
	}

	team.ActorStopTransfer(userID)

	isLeader := false
	if member, ok := team.prop.Data.Base.Members[userID]; !ok {
		team.Errorf("QuitTeam user:%s NotInTeam", userID)
		return ErrorCode.NotInTeam
	} else {
		isLeader = member.Status == Const.TeamStatus_LEADER
	}

	ia := team.GetActorByUserID(userID)
	if ia == nil {
		team.Errorf("QuitTeam get actor err user:%s", userID)
		return ErrorCode.GetRoleError
	}

	//清空该role的所有申请
	vals, err := DB.GetApply2InviteUtil().RemoveAllApplysInRole(ia.GetID())
	if err == nil {
		go func() {
			for i := 0; i < len(vals); i++ {
				val := vals[i]
				if srvID, err := DB.TeamUtil().GetSrvID(val.TeamID); err == nil {
					Core.Inst.RpcByID(srvID, "RPC_RemoveTeamApply", val.TeamID, ia.GetID())
				}
			}
		}()
	}

	ia.GetProp().(*Prop.RoleProp).SyncClearRequestSkill()

	ia.GetProp().(*Prop.RoleProp).SyncSetTeamID("")
	team.prop.RemoveTeamMember(userID)
	team.teamChatChannelDelMember(ia.(*_Role))

	if isLeader {
		team.electLeader()
	}

	err = team.RemoveActor(ia.GetID())
	if err != nil {
		team.Errorf("QuitTeam remove actor id:%s err:%s", userID, err)
		return ErrorCode.QuitTeamError
	}

	team.onMemberChange(userID, false)

	team.FlushToDB()
	team.FlushToCache()
	ia.FlushToDB()
	ia.FlushToCache()

	team.Infof("QuitTeam success user:%s", userID)
	return ErrorCode.Success
}

func (team *_Team) onMemberChange(userID string, join bool) {
	minProgress := uint32(math.MaxUint32)
	team.TraversalActor(func(ia Space.IActor) {
		if _, ok := team.prop.Data.Base.Members[ia.GetID()]; ok {
			progress := ia.GetProp().(*Prop.RoleProp).Data.Base.RaidProgress
			if minProgress > progress {
				minProgress = progress
			}
		}
	})
	if minProgress == 0 {
		minProgress = 1
	}
	team.prop.SyncSetRaidProgress(minProgress)

	team.UpdateSeasonDataOnMemberChange()
}

func (team *_Team) StartLeave(userID, joinTeamID string, reason uint8, callbackSrvID string, callbackID uint32) (int32, bool) {
	ia := team.GetActorByUserID(userID)
	if ia == nil {
		team.Errorf("StartLeave NotInTeam user:%s joinTeam:%s reason:%d", userID, joinTeamID, reason)
		return ErrorCode.NotInTeam, false
	}

	if (reason == Const.TransferReason_Apply || reason == Const.TransferReason_Invite) && len(team.prop.Data.Base.Members) != 1 {
		team.Errorf("StartLeave AlreadInTeam user:%s joinTeam:%s reason:%d", userID, joinTeamID, reason)
		return ErrorCode.AlreadInTeam, false
	}

	if _, ok := team.PendingLeaveUser[userID]; ok {
		team.Errorf("StartLeave RoleTransfering user:%s joinTeam:%s reason:%d", userID, joinTeamID, reason)
		return ErrorCode.RoleTransfering, false
	}

	team.ActorStartTransfer(userID)

	//发送通知

	team.PendingTransferEvent = append(team.PendingTransferEvent, func() {
		Core.Inst.RpcByID(callbackSrvID, "RPC_SendTransferCallback", callbackID)
	})

	team.Infof("StartLeave success user:%s joinTeam:%s reason:%d", userID, joinTeamID, reason)
	return ErrorCode.Success, team.GetState() == TeamState_Raidbattling
}

func (team *_Team) StartLeaveRollback(userID string) int32 {
	team.ActorStopTransfer(userID)

	return ErrorCode.Success
}

func (team *_Team) HoldPlace(userID string, reason uint8) int32 {
	if _, ok := team.PlaceHolder[userID]; ok {
		team.Errorf("HoldPlace RepeatPlaceHolder user:%s reason:%d", userID, reason)
		return ErrorCode.RepeatPlaceHolder
	}

	if len(team.prop.Data.Base.Members)+len(team.PlaceHolder) >= Const.TeamMaxMemberNum {
		team.Errorf("HoldPlace PlaceHolderFull user:%s reason:%d", userID, reason)
		return ErrorCode.TeamFull
	}

	if _, ok := team.prop.Data.Base.Members[userID]; ok {
		team.Errorf("HoldPlace AlreadInTeam user:%s reason:%d", userID, reason)
		return ErrorCode.AlreadInTeam
	}

	team.AddPlaceHolder(userID)

	team.Infof("HoldPlace success user:%s reason:%d", userID, reason)
	return ErrorCode.Success
}

func (team *_Team) HoldPlaceRollback(userID string) int32 {
	team.RemovePlaceHolder(userID)

	return ErrorCode.Success
}

func (team *_Team) AddTeamMember(sourceTeamID, userID, callbackSrvID string, callbackID uint32) (int32, bool) {
	if _, ok := team.PlaceHolder[userID]; !ok {
		team.Error("AddTeamMember ArgsWrong ", sourceTeamID, userID)
		return ErrorCode.ArgsWrong, false
	}

	team.PendingTransferEvent = append(team.PendingTransferEvent, func() {
		Core.Inst.RpcByID(callbackSrvID, "RPC_SendTransferCallback", callbackID)
	})

	return ErrorCode.Success, team.GetState() == TeamState_Raidbattling
}

func (team *_Team) AddTeamApply(userID, message string) int32 {
	if len(team.prop.Data.Applys) >= Const.MaxApplys {
		team.Error("AddTeamApply ApplyOverLimit ", userID)
		return ErrorCode.ApplyOverLimit
	}

	team.prop.SyncAddApplyInfo(userID, message)

	return ErrorCode.Success
}

func (team *_Team) RemoveTeamApply(userID string) int32 {
	team.prop.SyncRemoveApplyInfo(userID)

	return ErrorCode.Success
}

func (team *_Team) electLeader() {
	if len(team.prop.Data.Base.Members) == 0 {
		return
	}

	minJoinTime := int64(math.MaxInt64)
	var tmpUserID string
	for key, value := range team.prop.Data.Base.Members {
		if value.JoinTime < minJoinTime {
			tmpUserID = key
			minJoinTime = value.JoinTime
		}
	}

	team.prop.SyncSetStatus(tmpUserID, Const.TeamStatus_LEADER)
}
