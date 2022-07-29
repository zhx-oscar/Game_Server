package main

import (
	"Cinder/Space"
	"Daisy/Const"
	"Daisy/Data"
	"Daisy/ErrorCode"
	"Daisy/Prop"
	"Daisy/Proto"
	log "github.com/cihub/seelog"
	"math/rand"
	"time"
)

func (user *_User) RPC_RaidBattleEnd(fightID string) (int32, *Proto.OfflineAwardData) {
	team := user.GetSpace().(*_Team)
	if team.state == TeamState_Raidbattling {
		ret, it := team.OnRaidBattleOnlineDrop(user.role.GetID())
		team.OnUserRaidBattleEnd(user.GetID(), fightID)
		return ret, it
	}

	log.Debug("RPC_RaidBattleEnd")
	return ErrorCode.Failure, nil
}

func (user *_User) RPC_SearchBoss() int32 {
	team := user.GetSpace().(*_Team)
	if len(team.PendingLeaveUser) > 0 {
		team.Debug("队伍有成员在转移中，不触发战斗")
		return ErrorCode.RoleTransfering
	}

	if !team.CanSetState(TeamState_Raidbattling) {
		team.Errorf("state can't change form:%d to:%d", team.GetState(), TeamState_Raidbattling)
		return ErrorCode.Failure
	}

	battleAreaCfg, ok := Data.GetSceneConfig().BattleArea_ConfigItems[team.prop.Data.Raid.Progress]
	if !ok {
		team.Errorf("推图进度%d未配置", team.prop.Data.Raid.Progress)
		return ErrorCode.RaidProgressInvalid
	}

	if team.prop.Data.Raid.OwnTickets < battleAreaCfg.BossTickets {
		team.Error("挑战BOSS所需的门票不足")
		return ErrorCode.TicketsNotEnough
	}

	// BOSS战场ID目前暂定为0
	team.SetState(TeamState_Raidbattling, true, battleAreaCfg.BossBattleArea)

	return ErrorCode.Success
}

func (user *_User) RPC_GetCurFightResult() (int32, bool, *Proto.FightResult) {
	team := user.GetSpace().(*_Team)
	if fight := team.GetCurRaidFight(); fight != nil {
		return ErrorCode.Success, true, fight
	}

	return ErrorCode.Success, false, nil
}

func (user *_User) RPC_Energize() int32 {
	if user.role == nil {
		user.Error("RPC_Energize Failed")
		return ErrorCode.RoleIsNil
	}
	return user.role.Energize()
}

func (r *_Role) Energize() int32 {
	cost := uint32(0)
	finded := false
	for i := 1; i <= len(Data.GetFastBattleConfig().Energize_ConfigItems); i++ {
		if item, ok := Data.GetFastBattleConfig().Energize_ConfigItems[uint32(i)]; ok {
			if r.prop.Data.FastBattle.EnergizeNum < item.Times {
				cost = item.Cost
				finded = true
				break
			}
		}
	}
	if !finded {
		r.Error("Energize_Config 查不到本次充能所需的消耗")
		return ErrorCode.ConfigNotExist
	}

	multi := float32(1.0)
	total := uint32(0)
	rnd := uint32(rand.New(rand.NewSource(time.Now().Unix())).Intn(10000))
	for i := 1; i <= len(Data.GetFastBattleConfig().SpeedUp_ConfigItems); i++ {
		if item, ok := Data.GetFastBattleConfig().SpeedUp_ConfigItems[uint32(i)]; ok {
			total += item.Probability
			if rnd < total {
				multi = item.Multi
				break
			}
		}
	}

	if r.RemoveDiamond(cost, Const.Energize) {
		r.prop.SyncAddEnergize(multi)
		r.Info("Energize success")
		return ErrorCode.Success
	} else {
		r.Error("充能钱不够")
		return ErrorCode.DiamondNotEnough
	}
}

func (user *_User) RPC_FastBattle() int32 {
	if user.role == nil {
		user.Error("RPC_FastBattle Failed")
		return ErrorCode.RoleIsNil
	}
	return user.role.FastBattle()
}

func (r *_Role) FastBattle() int32 {
	if r.prop.Data.FastBattle.AwardMulti == 0 {
		r.Error("加速前先充能")
		return ErrorCode.UnEnergize
	}

	stage := uint32(0)
	for i := 1; i <= len(r.prop.Data.FastBattle.StageInfo); i++ {
		if info, ok := r.prop.Data.FastBattle.StageInfo[uint32(i)]; ok {
			if info.MyTimes < info.MaxMyTimes {
				stage = uint32(i)
				break
			}
		}
	}

	if stage == 0 {
		r.Error("加速次数耗尽")
		return ErrorCode.SpeedUpExhaust
	}

	cost := uint32(0)
	if item, ok := Data.GetFastBattleConfig().FastBattle_ConfigItems[stage]; ok {
		cost = item.Cost
	} else {
		r.Error("FastBattle_Config 查不到加速所需的消耗")
		return ErrorCode.ConfigNotExist
	}

	if r.GetDiamond() < uint64(cost) {
		r.Error("加速钱不够")
		return ErrorCode.DiamondNotEnough
	}

	if !r.HasEnoughSpace(1, int32(Proto.ContainerEnum_EquipBag), int32(Proto.ContainerEnum_SkillBag)) {
		r.Error("加速背包空间不足")
		return ErrorCode.EquipBagNotEnoughSpace
	}

	team := r.GetSpace().(*_Team)
	errCode := team.StartFastBattle(r.GetID(), stage, r.prop.Data.FastBattle.AwardMulti)
	if errCode != ErrorCode.Success {
		r.Error("StartFastBattle fail,errCode:", errCode)
		return errCode
	}

	r.RemoveDiamond(cost, Const.SpeedUp)
	r.prop.SyncClearFastBattleMulti()
	team.TraversalActor(func(ia Space.IActor) {
		if _, ok := team.prop.Data.Base.Members[ia.GetID()]; ok {
			roleProp := ia.GetProp().(*Prop.RoleProp)
			roleProp.SyncAddFastBattle(stage, r.GetID() == ia.GetID())
		}
	})

	r.Info("FastBattle success")
	return ErrorCode.Success
}

func (user *_User) RPC_FastBattleEnd() int32 {
	team := user.GetSpace().(*_Team)
	if team.state == TeamState_FastBattleing {
		team.OnUserFastBattleEnd(user.GetID())
		user.Info("RPC_FastBattleEnd")
		return ErrorCode.Success
	}

	user.Error("RPC_FastBattleEnd 不在快速战斗状态中")
	return ErrorCode.Failure
}
