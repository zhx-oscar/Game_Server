package main

import (
	"Cinder/Base/User"
	"Daisy/Proto"
)

type _RaidBattleModel struct {
	userBattleEnd map[string]bool
	lastFight     *Proto.FightResult
}

func (team *_Team) OnUserRaidBattleEnd(userID, fightID string) {
	team.Infof("OnUserRaidBattleEnd UserID:%s FightID:%s", userID, fightID)
	if team.raidBattle != nil && team.raidBattle.lastFight.Id == fightID {
		team.raidBattle.userBattleEnd[userID] = true

		finish := true
		team.TraversalUser(func(user User.IUser) bool {
			if _, ok := team.raidBattle.userBattleEnd[user.GetID()]; !ok {
				finish = false
				return false
			}
			return true
		})

		if finish {
			if team.raidBattle.lastFight.IsBoss {
				team.SetState(TeamState_Running)
			} else {
				team.Infof("非BOSS战斗结束后瞬移到目标路点:%d，避免客户端特工的停顿", team.targetRoadIndex)
				team.SetState(TeamState_Running, team.targetRoadIndex)
			}
		}
	}
}

func (team *_Team) GetCurRaidFight() *Proto.FightResult {
	if team.raidBattle != nil {
		return team.raidBattle.lastFight
	}
	return nil
}
