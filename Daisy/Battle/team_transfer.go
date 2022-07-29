package main

import (
	CConst "Cinder/Base/Const"
	"Cinder/Base/Core"
	"Daisy/Const"
	"Daisy/DB"
	"Daisy/DHDB"
	"Daisy/ErrorCode"
	"Daisy/NotifyCode"
	"errors"
	"fmt"
	log "github.com/cihub/seelog"
	"sync"
	"sync/atomic"
	"time"
)

type TransferCallback struct {
	Err  error
	Done chan *TransferCallback
}

var transferCallbackID uint32
var transferCallbacks sync.Map

func genTransferCallback() (uint32, *TransferCallback) {
	callbackID := atomic.AddUint32(&transferCallbackID, 1)
	callback := &TransferCallback{
		Done: make(chan *TransferCallback, 1),
	}
	transferCallbacks.Store(callbackID, callback)

	go func() {
		time.Sleep(5 * time.Minute)

		if value, ok := transferCallbacks.Load(callbackID); ok {
			transferCallbacks.Delete(callbackID)

			cb := value.(*TransferCallback)
			cb.Err = errors.New(fmt.Sprintf("time out %d", callbackID))
			cb.Done <- cb
		}
	}()

	return callbackID, callback
}

func OnRecvTransferCallback(callbackID uint32) int32 {
	if value, ok := transferCallbacks.Load(callbackID); ok {
		transferCallbacks.Delete(value)

		cb := value.(*TransferCallback)
		cb.Done <- cb
	}

	return ErrorCode.Success
}

func Transfer(oldTeamID, newTeamID, userID, notifier string, reason uint8) (err error) {
	log.Infof("Transfer start %s %s->%s %d", userID, oldTeamID, newTeamID, reason)
	curSrvID := Core.Inst.GetServiceID()

	roleCache, err := DHDB.GetRoleCache(userID)
	if err != nil {
		log.Errorf("Transfer GetRoleCache id:%s err:%s", userID, err)
		return err
	}
	pendingRole := roleCache.Base

	newSrvID, err := LoadTeam(newTeamID)
	if err != nil {
		log.Errorf("Transfer load new team:%s err:%s", newTeamID, err)
		return err
	}
	oldSrvID, err := LoadTeam(oldTeamID)
	if err != nil {
		log.Errorf("load old team:%s err:%s", oldTeamID, err)
		return err
	}

	ret := <-Core.Inst.RpcByID(newSrvID, "RPC_HoldPlace", newTeamID, userID, reason)
	if ret.Err != nil {
		log.Errorf("Transfer call RPC_HoldPlace err:%s newTeam:%s user:%s", ret.Err, newTeamID, userID)
		return ret.Err
	}
	if retCode := ret.Ret[0].(int32); retCode != ErrorCode.Success {
		if retCode == ErrorCode.TeamFull {
			// send notify
			if reason == Const.TransferReason_Invite {
				Core.Inst.RpcByID(newSrvID, "RPC_SendNotify", newTeamID, uint32(NotifyCode.JoinTeamFailed), notifier, pendingRole.Name, false, false)
				Core.Inst.RpcByID(oldSrvID, "RPC_SendNotify", oldTeamID, uint32(NotifyCode.TargetTeamFull), userID, "", false, false)
			} else if reason == Const.TransferReason_Apply {
				Core.Inst.RpcByID(newSrvID, "RPC_SendNotify", newTeamID, uint32(NotifyCode.SelfTramFull), notifier, "", false, false)
			}
		}

		log.Errorf("Transfer call RPC_HoldPlace errCode:%d newTeam:%s user:%s", retCode, newTeamID, userID)
		return errors.New("transfer call RPC_HoldPlace err")
	}

	defer func() {
		if err != nil {
			Core.Inst.RpcByID(newSrvID, "RPC_HoldPlaceRollback", newTeamID, userID)
		}
	}()

	callbackID, callback := genTransferCallback()

	ret = <-Core.Inst.RpcByID(oldSrvID, "RPC_StartLeave", oldTeamID, userID, newTeamID, reason, curSrvID, callbackID)
	if ret.Err != nil {
		log.Errorf("Transfer call RPC_StartLeave err:%s oldTeam:%s user:%s", ret.Err, oldTeamID, userID)
		return ret.Err
	}
	if retCode := ret.Ret[0].(int32); retCode != ErrorCode.Success {
		if retCode == ErrorCode.AlreadInTeam {
			// send notify
			if reason == Const.TransferReason_Invite {
				Core.Inst.RpcByID(oldSrvID, "RPC_SendNotify", oldTeamID, uint32(NotifyCode.JoinTeamSuccessWithName), userID, "", false, false)
				Core.Inst.RpcByID(newSrvID, "RPC_SendNotify", newTeamID, uint32(NotifyCode.RefuseFormTeamInvitatoin), notifier, pendingRole.Name, false, false)
			} else if reason == Const.TransferReason_Apply {
				Core.Inst.RpcByID(newSrvID, "RPC_SendNotify", newTeamID, uint32(NotifyCode.JoinTeamSuccessWithName), notifier, pendingRole.Name, false, false)
			}
		}

		if retCode == ErrorCode.RoleTransfering {
			Core.Inst.RpcByID(oldSrvID, "RPC_SendNotify", oldTeamID, uint32(NotifyCode.JoiningTeam), userID, "", false, false)
		}

		log.Errorf("Transfer call RPC_StartLeave errCode:%d oldTeam:%s user:%s", retCode, oldTeamID, userID)
		return errors.New("transfer call RPC_StartLeave err")
	}
	oldBattling := ret.Ret[1].(bool)
	if oldBattling {
		// send notify
		if reason == Const.TransferReason_Apply || reason == Const.TransferReason_Invite || reason == Const.TransferReason_AutoJoin {
			Core.Inst.RpcByID(oldSrvID, "RPC_SendNotify", oldTeamID, uint32(NotifyCode.ApplyTeam), userID, "", true, true)
			defer Core.Inst.RpcByID(oldSrvID, "RPC_SendNotify", oldTeamID, uint32(NotifyCode.ApplyTeam), userID, "", true, false)
		}
		if reason == Const.TransferReason_Apply || reason == Const.TransferReason_Invite {
			Core.Inst.RpcByID(newSrvID, "RPC_SendNotify", newTeamID, uint32(NotifyCode.BeAgreeedJoinTeam), notifier, pendingRole.Name, false, false)
		}
	}

	defer func() {
		if err != nil {
			Core.Inst.RpcByID(oldSrvID, "RPC_StartLeaveRollback", oldTeamID, userID)
		}
	}()

	cb := <-callback.Done
	if cb.Err != nil {
		log.Errorf("Transfer callback err:%s", cb.Err)
		return ret.Err
	}

	ret = <-Core.Inst.RpcByID(newSrvID, "RPC_AddTeamMember", oldTeamID, newTeamID, userID, curSrvID, callbackID)
	if ret.Err != nil {
		log.Errorf("Transfer call RPC_AddTeamMember err:%s newTeam:%s user:%s", ret.Err, newTeamID, userID)
		return ret.Err
	}
	if retCode := ret.Ret[0].(int32); retCode != ErrorCode.Success {
		log.Errorf("Transfer call RPC_AddTeamMember errCode:%d newTeam:%s user:%s", retCode, newTeamID, userID)
		return errors.New("transfer call RPC_AddTeamMember err")
	}
	newBattling := ret.Ret[1].(bool)
	if newBattling {
		// send notify
		if reason == Const.TransferReason_Apply || reason == Const.TransferReason_Invite || reason == Const.TransferReason_AutoJoin {
			Core.Inst.RpcByID(oldSrvID, "RPC_SendNotify", oldTeamID, uint32(NotifyCode.ApplyTeam), userID, "", true, true)
			defer Core.Inst.RpcByID(oldSrvID, "RPC_SendNotify", oldTeamID, uint32(NotifyCode.ApplyTeam), userID, "", true, false)
			Core.Inst.RpcByID(newSrvID, "RPC_SendNotify", newTeamID, uint32(NotifyCode.ApplyTeamWithName), "", pendingRole.Name, true, true)
			defer Core.Inst.RpcByID(newSrvID, "RPC_SendNotify", newTeamID, uint32(NotifyCode.ApplyTeamWithName), "", pendingRole.Name, true, false)
		}
	}

	cb = <-callback.Done
	if cb.Err != nil {
		log.Errorf("Transfer callback err:%s", cb.Err)
		return ret.Err
	}

	ret = <-Core.Inst.RpcByID(oldSrvID, "RPC_QuitTeam", oldTeamID, userID)
	if ret.Err != nil {
		log.Errorf("Transfer call RPC_QuitTeam err:%s oldTeam:%s user:%s", ret.Err, oldTeamID, userID)
		return ret.Err
	}
	if retCode := ret.Ret[0].(int32); retCode != ErrorCode.Success {
		log.Errorf("Transfer call RPC_QuitTeam errCode:%d oldTeam:%s user:%s", retCode, oldTeamID, userID)
		return errors.New("transfer call RPC_QuitTeam err")
	}
	online := ret.Ret[1].(bool)

	if reason == Const.TransferReason_Quit {
		Core.Inst.RpcByID(oldSrvID, "RPC_SendNotify", oldTeamID, uint32(NotifyCode.LeaveTeam), "", pendingRole.Name, false, false)
	} else if reason == Const.TransferReason_Kick {
		Core.Inst.RpcByID(oldSrvID, "RPC_SendNotify", oldTeamID, uint32(NotifyCode.KickOutTeam), "", pendingRole.Name, false, false)
	}

	ret = <-Core.Inst.RpcByID(newSrvID, "RPC_JoinTeam", newTeamID, userID, online)
	if ret.Err != nil {
		log.Errorf("Transfer call RPC_JoinTeam err:%s newTeam:%s user:%s", ret.Err, newTeamID, userID)
		return ret.Err
	}
	if retCode := ret.Ret[0].(int32); retCode != ErrorCode.Success {
		log.Errorf("Transfer call RPC_JoinTeam errCode:%d newTeam:%s user:%s", retCode, newTeamID, userID)
		return errors.New("transfer call RPC_JoinTeam err")
	}
	log.Infof("Transfer success %s %s->%s %d", userID, oldTeamID, newTeamID, reason)

	if reason == Const.TransferReason_Apply || reason == Const.TransferReason_Invite || reason == Const.TransferReason_AutoJoin {
		Core.Inst.RpcByID(newSrvID, "RPC_SendNotify", newTeamID, uint32(NotifyCode.NotifyJoinTeam), "", pendingRole.Name, false, false)
	}

	Core.Inst.CallRpcToUser(userID, CConst.Agent, "RPC_TransferEnd", reason)

	if reason == Const.TransferReason_Invite {
		//清除角色身上的所有邀请
		vals, err := DB.GetApply2InviteUtil().RemoveAllInvitesInRole(userID)
		if err != nil {
			log.Error("Transfer DB ", err)
			return err
		}
		Core.Inst.CallRpcToUser(userID, CConst.Space, "RPC_RemoveRoleInvite", "", true)

		//send notify to all other teams
		for i := 0; i < len(vals); i++ {
			val := vals[i]
			if val.TeamID != newTeamID {
				if srvID, err := DB.TeamUtil().GetSrvID(val.TeamID); err == nil {
					Core.Inst.RpcByID(srvID, "RPC_SendNotify", val.TeamID, uint32(NotifyCode.RefuseFormTeamInvitatoin), val.InviteParam.Instigator, pendingRole.Name, false, false)
				}
			}
		}

	}

	return nil
}
