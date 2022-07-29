package main

import (
	"Cinder/Base/Const"
	"Cinder/Base/User"
	"Cinder/Base/Util"
	"Cinder/Space"
	"Daisy/Fight"
	"Daisy/Prop"
	"Daisy/Proto"
	"time"
)

type _RaidBattleState struct {
	team          *_Team
	battleEndTime time.Time
}

func NewRaidBattleState(team *_Team) *_RaidBattleState {
	return &_RaidBattleState{
		team: team,
	}
}

func (battle *_RaidBattleState) OnEnter(args ...interface{}) {
	if len(args) < 2 {
		return
	}
	isBoss := args[0].(bool)
	battleFieldID := args[1].(uint32)

	// 调用fight接口生成战报
	sceneInfo, err := battle.team.BuildRaidSceneInfo(battle.team.prop.Data.Raid.Progress, isBoss, battleFieldID)
	if err != nil {
		battle.team.Error("EnterRaidBattle 生成对战场景信息失败", err)
		return
	}

	result, err := Fight.Play(sceneInfo)
	if err != nil {
		battle.team.Error("EnterRaidBattle 模拟战斗失败", err)
		return
	}

	result.Id = Util.GetGUID()
	result.BattleFieldID = battleFieldID
	result.Progress = battle.team.prop.Data.Raid.Progress
	result.IsBoss = isBoss
	result.StartTimestamp = time.Now().Unix()
	result.DelayRewards = map[uint32]uint32{}

	battle.team.raidBattle = &_RaidBattleModel{
		lastFight:     result,
		userBattleEnd: make(map[string]bool),
	}

	battle.team.TraversalUser(func(user User.IUser) bool {
		user.Rpc(Const.Agent, "RPC_RaidBattleResult", result)
		return true
	})
	// 获取本次战斗耗时，加上网络延迟的误差ms
	battle.battleEndTime = time.Now().Add(time.Duration(result.Time+10000) * time.Millisecond)
	battle.team.Infof("EnterRaidBattle isBoss:%v battleField:%d realTime:%d FightID:%s", isBoss, battleFieldID, result.Time, result.Id)

	battle.team.EnterRaidBattleDrop(result)
}

func (battle *_RaidBattleState) OnLeave() {
	//进入战斗状态有可能因为战报生成失败而退出状态，raidBattle要判空
	if battle.team.raidBattle == nil {
		return
	}
	battle.team.Infof("LeaveRaidBattle FightID:%s", battle.team.raidBattle.lastFight.Id)

	if battle.team.raidBattle.lastFight.IsBoss {
		battle.team.prop.SyncSetOwnTickets(0)
		if battle.team.raidBattle.lastFight.Outcome[0] == Proto.Camp_Red {
			battle.team.Pass(battle.team.prop.Data.Raid.Progress)
			battle.team.prop.SyncSetRaidProgress(battle.team.prop.Data.Raid.Progress + 1)
			battle.team.TraversalActor(func(ia Space.IActor) {
				if _, ok := battle.team.prop.Data.Base.Members[ia.GetID()]; ok {
					roleProp := ia.GetProp().(*Prop.RoleProp)
					if roleProp.Data.Base.RaidProgress < battle.team.prop.Data.Raid.Progress {
						roleProp.SyncSetRaidProgress(battle.team.prop.Data.Raid.Progress)
					}
				}
			})
		}
	} else {
		if battle.team.raidBattle.lastFight.Outcome[0] == Proto.Camp_Red {
			battle.team.prop.SyncSetOwnTickets(battle.team.prop.Data.Raid.OwnTickets + 1)
		}
	}

	battle.team.LeaveRaidBattleDrop()

	//落地继承的战斗属性
	battle.team.saveInheritAttr(battle.team.raidBattle.lastFight)

	//todo 战斗结束的经验，掉落处理

	battle.team.raidBattle = nil
}

func (battle *_RaidBattleState) OnLoop(delta time.Duration) {
	if (battle.team.raidBattle != nil && time.Now().After(battle.battleEndTime)) || battle.team.raidBattle == nil {
		if battle.team.raidBattle != nil && !battle.team.raidBattle.lastFight.IsBoss {
			battle.team.Infof("非BOSS战斗结束后瞬移到目标路点:%d，避免客户端特工的停顿", battle.team.targetRoadIndex)
			battle.team.SetState(TeamState_Running, battle.team.targetRoadIndex)
		} else {
			battle.team.SetState(TeamState_Running)
		}

		if battle.team.raidBattle != nil {
			battle.team.Errorf("%s 战斗超时", battle.team.raidBattle.lastFight.Id)
		}
	}
}

func (battle *_RaidBattleState) CanChangeTo(state uint8) bool {
	return state == TeamState_Running
}
