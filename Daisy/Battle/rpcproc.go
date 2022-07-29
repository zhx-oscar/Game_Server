package main

import (
	"Cinder/Space"
	"Daisy/ErrorCode"
	log "github.com/cihub/seelog"
)

type _RPCProc struct {
}

// RPC_LoadTeamFromDB 从DB加载队伍
func (proc *_RPCProc) RPC_LoadTeamFromDB(teamID string, cellServerID uint32) int32 {
	log.Infof("RPC_LoadTeamFromDB Team: %s Cell: %d", teamID, cellServerID)

	Space.Inst.CreateSpace(teamID, []byte{}, cellServerID)

	return ErrorCode.Success
}

func (proc *_RPCProc) RPC_EnterBattle(teamID, userID string) int32 {
	if err := Space.Inst.EnterSpace(userID, teamID); err != nil {
		log.Errorf("RPC_EnterBattle EnterSpace err %s Team: %s User: %s", err, teamID, userID)
		return ErrorCode.EnterTeamSpaceError
	}

	space, err := Space.Inst.GetSpace(teamID)
	if err != nil {
		log.Errorf("RPC_EnterBattle GetSpace err %s Team: %s User: %s", err, teamID, userID)
		return ErrorCode.GetTeamError
	}

	ret := <-space.SafeCall("MemberOnline", userID)
	if ret.Err != nil {
		log.Errorf("RPC_EnterBattle Call MemberOnline err %s Team: %s User: %s", ret.Err, teamID, userID)
		return ErrorCode.TeamMemberOnlineError
	}
	if retCode := ret.Ret[0].(int32); retCode != 0 {
		log.Errorf("RPC_EnterBattle err %d Team: %s User: %s", retCode, teamID, userID)
		return retCode
	}

	log.Infof("RPC_EnterBattle Success Team: %s User: %s", teamID, userID)
	return ErrorCode.Success
}

func (proc *_RPCProc) RPC_LeaveBattle(teamID, userID string) int32 {
	space, err := Space.Inst.GetSpace(teamID)
	if err != nil {
		log.Errorf("RPC_LeaveBattle GetSpace err %s Team: %s User: %s", err, teamID, userID)
		return ErrorCode.GetTeamError
	}

	ret := <-space.SafeCall("MemberOffline", userID)
	if ret.Err != nil {
		log.Errorf("RPC_LeaveBattle Call MemberOffline err %s Team: %s User: %s", ret.Err, teamID, userID)
		return ErrorCode.TeamMemberOfflineError
	}
	if retCode := ret.Ret[0].(int32); retCode != 0 {
		log.Errorf("RPC_LeaveBattle err %d Team: %s User: %s", retCode, teamID, userID)
		return retCode
	}

	if err = Space.Inst.LeaveSpace(userID); err != nil {
		log.Errorf("RPC_LeaveBattle LeaveSpace err %s Team: %s User: %s", err, teamID, userID)
		return ErrorCode.LeaveTeamSpaceError
	}

	log.Infof("RPC_LeaveBattle Success Team: %s User: %s", teamID, userID)
	return ErrorCode.Success
}

func (proc *_RPCProc) RPC_QuitTeam(teamID, userID string) (int32, bool) {
	space, err := Space.Inst.GetSpace(teamID)
	if err != nil {
		log.Error("RPC_QuitTeam GetSpace err", err, teamID)
		return ErrorCode.GetTeamError, false
	}

	online := false
	_, err = space.GetUser(userID)
	if err == nil {
		online = true
		Space.Inst.LeaveSpace(userID)
	}

	ret := <-space.SafeCall("QuitTeam", userID)
	if ret.Err != nil {
		log.Error("RPC_QuitTeam safeCall err", ret.Err, teamID, userID)
		return ErrorCode.Timeout, online
	}
	if retCode := ret.Ret[0].(int32); retCode != 0 {
		log.Error("RPC_QuitTeam safeCall errCode", retCode, teamID, userID)
		return retCode, online
	}

	return ErrorCode.Success, online
}

func (proc *_RPCProc) RPC_JoinTeam(teamID, userID string, online bool) int32 {
	space, err := Space.Inst.GetSpace(teamID)
	if err != nil {
		log.Error("RPC_JoinTeam GetSpace err", err, teamID, userID)
		return ErrorCode.GetTeamError
	}

	ret := <-space.SafeCall("JoinTeam", userID)
	if ret.Err != nil {
		log.Error("RPC_JoinTeam safeCall err", ret.Err, teamID, userID)
		return ErrorCode.Timeout
	}
	if retCode := ret.Ret[0].(int32); retCode != 0 {
		log.Error("RPC_JoinTeam safeCall errCode", retCode, teamID, userID)
		return retCode
	}

	if online {
		Space.Inst.EnterSpace(userID, teamID)
	}

	return ErrorCode.Success
}

func (proc *_RPCProc) RPC_AddTeamMember(sourceTeamID, targetTeamID, userID, callbackSrvID string, callbackID uint32) (int32, bool) {
	space, err := Space.Inst.GetSpace(targetTeamID)
	if err != nil {
		log.Error("RPC_AddTeamMember GetSpace err", err, targetTeamID)
		return ErrorCode.GetTeamError, false
	}

	ret := <-space.SafeCall("AddTeamMember", sourceTeamID, userID, callbackSrvID, callbackID)
	if ret.Err != nil {
		log.Error("RPC_AddTeamMember safeCall err", ret.Err, sourceTeamID, targetTeamID, userID)
		return ErrorCode.Timeout, false
	}
	if retCode := ret.Ret[0].(int32); retCode != 0 {
		log.Error("RPC_AddTeamMember safeCall errCode", retCode, sourceTeamID, targetTeamID, userID)
		return retCode, false
	}

	return ErrorCode.Success, ret.Ret[1].(bool)
}

func (proc *_RPCProc) RPC_HoldPlace(teamID, userID string, reason uint8) int32 {
	space, err := Space.Inst.GetSpace(teamID)
	if err != nil {
		log.Error("RPC_HoldPlace GetSpace err", err, teamID, reason)
		return ErrorCode.GetTeamError
	}

	ret := <-space.SafeCall("HoldPlace", userID, reason)
	if ret.Err != nil {
		log.Error("RPC_HoldPlace safeCall err", ret.Err, teamID, userID, reason)
		return ErrorCode.Timeout
	}
	if retCode := ret.Ret[0].(int32); retCode != 0 {
		log.Error("RPC_HoldPlace safeCall errCode", retCode, teamID, userID, reason)
		return retCode
	}

	return ErrorCode.Success
}

func (proc *_RPCProc) RPC_HoldPlaceRollback(teamID, userID string) int32 {
	space, err := Space.Inst.GetSpace(teamID)
	if err != nil {
		log.Error("RPC_HoldPlaceRollback GetSpace err", err, teamID, userID)
		return ErrorCode.GetTeamError
	}

	ret := <-space.SafeCall("HoldPlaceRollback", userID)
	if ret.Err != nil {
		log.Error("RPC_HoldPlaceRollback safeCall err", ret.Err, teamID, userID)
		return ErrorCode.Timeout
	}
	if retCode := ret.Ret[0].(int32); retCode != 0 {
		log.Error("RPC_HoldPlaceRollback safeCall errCode", retCode, teamID, userID)
		return retCode
	}

	return ErrorCode.Success
}

func (proc *_RPCProc) RPC_AddTeamApply(teamID, userID, message string) int32 {
	space, err := Space.Inst.GetSpace(teamID)
	if err != nil {
		log.Error("RPC_AddTeamApply GetSpace err", err, teamID, userID)
		return ErrorCode.GetTeamError
	}

	ret := <-space.SafeCall("AddTeamApply", userID, message)
	if ret.Err != nil {
		log.Error("RPC_AddTeamApply safeCall err", ret.Err, teamID, userID)
		return ErrorCode.Timeout
	}
	if retCode := ret.Ret[0].(int32); retCode != 0 {
		log.Error("RPC_AddTeamApply safeCall errCode", retCode, teamID, userID)
		return retCode
	}

	return ErrorCode.Success
}

func (proc *_RPCProc) RPC_RemoveTeamApply(teamID, userID string) int32 {
	space, err := Space.Inst.GetSpace(teamID)
	if err != nil {
		log.Error("RPC_RemoveTeamApply GetSpace err", err, teamID, userID)
		return ErrorCode.GetTeamError
	}

	ret := <-space.SafeCall("RemoveTeamApply", userID)
	if ret.Err != nil {
		log.Error("RPC_RemoveTeamApply safeCall err", ret.Err, teamID, userID)
		return ErrorCode.Timeout
	}
	if retCode := ret.Ret[0].(int32); retCode != 0 {
		log.Error("RPC_RemoveTeamApply safeCall errCode", retCode, teamID, userID)
		return retCode
	}

	return ErrorCode.Success
}

func (proc *_RPCProc) RPC_StartLeave(teamID, userID, joinTeamID string, reason uint8, callbackSrvID string, callbackID uint32) (int32, bool) {
	space, err := Space.Inst.GetSpace(teamID)
	if err != nil {
		log.Error("RPC_StartLeave GetSpace err", err, teamID, userID, joinTeamID, reason)
		return ErrorCode.GetTeamError, false
	}

	ret := <-space.SafeCall("StartLeave", userID, joinTeamID, reason, callbackSrvID, callbackID)
	if ret.Err != nil {
		log.Error("RPC_StartLeave safeCall err", ret.Err, teamID, userID, joinTeamID, reason)
		return ErrorCode.Timeout, false
	}
	if retCode := ret.Ret[0].(int32); retCode != 0 {
		log.Error("RPC_StartLeave safeCall errCode", retCode, teamID, userID, joinTeamID, reason)
		return retCode, false
	}

	return ErrorCode.Success, ret.Ret[1].(bool)
}

func (proc *_RPCProc) RPC_StartLeaveRollback(teamID, userID string) int32 {
	space, err := Space.Inst.GetSpace(teamID)
	if err != nil {
		log.Error("RPC_StartLeaveRollback GetSpace err", err, teamID, userID)
		return ErrorCode.GetTeamError
	}

	ret := <-space.SafeCall("StartLeaveRollback", userID)
	if ret.Err != nil {
		log.Error("RPC_StartLeaveRollback safeCall err", ret.Err, teamID, userID)
		return ErrorCode.Timeout
	}
	if retCode := ret.Ret[0].(int32); retCode != 0 {
		log.Error("RPC_StartLeaveRollback safeCall errCode", retCode, teamID, userID)
		return retCode
	}

	return ErrorCode.Success
}

func (proc *_RPCProc) RPC_SendNotify(teamID string, notifyID uint32, notifier, args string, isLong, showLong bool) int32 {
	space, err := Space.Inst.GetSpace(teamID)
	if err != nil {
		log.Error("RPC_SendNotify GetSpace err", err, teamID, notifyID, notifier, args, isLong, showLong)
		return ErrorCode.GetTeamError
	}

	ret := <-space.SafeCall("SendNotify", notifyID, notifier, args, isLong, showLong)
	if ret.Err != nil {
		log.Error("RPC_SendNotify safeCall err", ret.Err, teamID, notifyID, notifier, args, isLong, showLong)
		return ErrorCode.Timeout
	}
	if retCode := ret.Ret[0].(int32); retCode != 0 {
		log.Error("RPC_SendNotify safeCall errCode", retCode, teamID, notifyID, notifier, args, isLong, showLong)
		return retCode
	}

	return ErrorCode.Success
}

func (proc *_RPCProc) RPC_SendTransferCallback(callbackID uint32) int32 {
	OnRecvTransferCallback(callbackID)

	return ErrorCode.Success
}
