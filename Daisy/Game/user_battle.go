package main

import (
	"Cinder/Base/Const"
	"Cinder/Base/Core"
	"Cinder/Base/DistributeLock"
	"Cinder/Base/Util"
	DConst "Daisy/Const"
	"Daisy/DB"
	"Daisy/ErrorCode"
	"Daisy/Prop"
	"github.com/go-redis/redis/v7"
)

func (u *_User) RPC_EnterBattleWorld() int32 {
	if u.prop.Data.Base.Name == "" {
		u.Error("RPC_EnterBattleWorld 创角流程未走完")
		return ErrorCode.NameInvalid
	}

	if u.curTeamSrvID != "" {
		u.Error("RPC_EnterBattleWorld 重复操作")
		return ErrorCode.EnterSpaceError
	}

	teamID := u.prop.Data.Base.TeamID

	var srvID string
	var err error
	needCreateTeam := false

	var createTeamLock DistributeLock.ILocker
	defer func() {
		if createTeamLock != nil {
			createTeamLock.Unlock()
		}
	}()

	if teamID == "" {
		u.Errorf("RPC_EnterBattleWorld 角色所属队伍不存在, 创建新的, UserID:%s", u.GetID())

		//team的离线对象并将role加入team
		teamObj, err := Core.Inst.CreatePropObject(Prop.OfflineTeamObject, Util.GetGUID(), nil, nil)
		if err != nil {
			return ErrorCode.CreateTeamOfflinePropFailed
		}
		defer Core.Inst.DestroyPropObject(teamObj.GetID())
		defer u.FlushToDB()

		u.prop.SyncSetTeamID(teamObj.GetID())

		teamProp := teamObj.(*Prop.OfflineTeam).TeamProp
		teamProp.AddTeamMember(u.GetID(), DConst.TeamStatus_LEADER)

		teamID = teamObj.GetID()
		needCreateTeam = true
	} else {
		createTeamLock = DistributeLock.New(DConst.LoadTeamDLockPrefix + teamID)
		createTeamLock.Lock()
		srvID, err = DB.TeamUtil().GetSrvID(teamID)
		if err == redis.Nil {
			needCreateTeam = true
		} else if err != nil {
			u.Error("RPC_EnterBattleWorld 获取队伍信息失败", err)
			return ErrorCode.GetTeamError
		}
	}

	if needCreateTeam {
		// 通过负载均衡获取一个战斗服
		srvID, err = Core.Inst.GetSrvIDByType(Const.Space)
		if err != nil {
			u.Error("RPC_EnterBattleWorld 获取Battle服失败", err)
			return ErrorCode.TeamSrvFull
		}

		// 在目标战斗服上创建队伍
		ret := <-Core.Inst.RpcByID(srvID, "RPC_LoadTeamFromDB", teamID, uint32(0))
		if ret.Err != nil {
			u.Error("RPC_EnterBattleWorld call RPC_LoadTeamFromDB err", ret.Err)
			return ErrorCode.CreateTeamError
		}
		if retCode := ret.Ret[0].(int32); retCode != 0 {
			u.Error("RPC_EnterBattleWorld call RPC_LoadTeamFromDB failed", retCode)
			return ErrorCode.CreateTeamError
		}
	}

	// 队伍已经存在了, 上线
	ret := <-Core.Inst.RpcByID(srvID, "RPC_EnterBattle", teamID, u.GetID())
	if ret.Err != nil {
		u.Error("RPC_EnterBattleWorld call RPC_EnterBattle err", ret.Err)
		return ErrorCode.EnterSpaceError
	}
	if retCode := ret.Ret[0].(int32); retCode != 0 {
		u.Error("RPC_EnterBattleWorld call RPC_EnterBattle failed", retCode)
		return ErrorCode.EnterSpaceError
	}

	u.curTeamSrvID = srvID

	u.Info("RPC_EnterBattleWorld Success", u.GetID())

	return ErrorCode.Success
}

// RPC_LeaveBattleWorld 下线当前角色，回到选服界面
func (u *_User) RPC_LeaveBattleWorld() int32 {
	if u.curTeamSrvID == "" {
		u.Error("RPC_LeaveBattleWorld 重复操作")
		return ErrorCode.LeaveSpaceError
	}

	ret := <-Core.Inst.RpcByID(u.curTeamSrvID, "RPC_LeaveBattle", u.prop.Data.Base.TeamID, u.GetID())
	if ret.Err != nil {
		u.Error("RPC_LeaveBattleWorld call RPC_LeaveBattle err", ret.Err)
		return ErrorCode.LeaveSpaceError
	}
	if retCode := ret.Ret[0].(int32); retCode != 0 {
		u.Error("RPC_LeaveBattleWorld call RPC_LeaveBattle failed", retCode)
		return ErrorCode.LeaveSpaceError
	}

	u.Info("RPC_LeaveBattleWorld Success")
	u.curTeamSrvID = ""

	return ErrorCode.Success
}
