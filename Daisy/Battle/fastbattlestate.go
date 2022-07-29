package main

import (
	"Cinder/Base/Const"
	"Cinder/Base/User"
	"Cinder/Space"
	"Daisy/Battle/sceneconfigdata"
	"Daisy/Data"
	"Daisy/Proto"
	"time"
)

type _FastBattleState struct {
	team    *_Team
	endTime time.Time
	multi   float32
}

func NewFastBattleState(team *_Team) *_FastBattleState {
	return &_FastBattleState{
		team: team,
	}
}

func (fb *_FastBattleState) OnEnter(args ...interface{}) {
	if len(args) > 0 {
		fb.multi = args[0].(float32)
	}

	fb.endTime = time.Now().Add(12 * time.Second)
	fb.team.Infof("EnterFastBattleState multi:%f", fb.multi)

	velocity := Data.GetFastBattleConfig().Variable_ConfigItems[1].AccelSpeed
	duration := Data.GetFastBattleConfig().Variable_ConfigItems[1].AccelDuration
	passbyRoadIndexs, err := fb.team.GetFastBattlePassbyRoadIndexs(velocity, duration)
	if err != nil {
		fb.team.Error("CalcPassbyRoadIndexs err:", err)
		return
	}
	if len(passbyRoadIndexs) > 0 {
		fb.team.fastBattle.finishRoadIndex = passbyRoadIndexs[len(passbyRoadIndexs)-1]
	}

	passbyPoints := &Proto.PVector3Array{
		Data: make([]*Proto.PVector3, len(passbyRoadIndexs), len(passbyRoadIndexs)),
	}
	for i := 0; i < len(passbyRoadIndexs); i++ {
		point := fb.team.paths.Points[passbyRoadIndexs[i]]
		point = sceneconfigdata.UnconvertVector3(point)
		passbyPoints.Data[i] = &Proto.PVector3{
			X: point.X,
			Y: point.Y,
			Z: point.Z,
		}
	}
	fb.team.TraversalUser(func(iu User.IUser) bool {
		iu.Rpc(Const.Agent, "RPC_FastBattleStart", passbyPoints, duration)
		return true
	})
}

func (fb *_FastBattleState) OnLeave() {
	fb.team.Infof("LeaveFastBattleState")
	fb.team.TraversalActor(func(ia Space.IActor) {
		_, ok := fb.team.prop.Data.Base.Members[ia.GetID()]
		if !ok {
			return
		}

		_, ok = fb.team.fastBattle.enjoyableRoles[ia.GetID()]
		if !ok {
			return
		}

		role := ia.(*_Role)
		f, err := role.CalcAwardByDuration(fb.team.prop.Data.Raid.Progress, 2*time.Hour, fb.multi)
		if err != nil {
			role.Error(err)
			return
		}

		role.prop.SyncAddGold(f.Money)
		role.addCommanderExp(f.ActorExp)
		role.addFightingSpecialAgentExp(role.expBonus(uint64(f.SpecialAgentExp)))

		if iu, err := fb.team.GetUser(role.GetID()); err == nil {
			iu.Rpc(Const.Agent, "RPC_FastBattleAward", f, fb.multi)
		}
	})

	fb.team.fastBattle.opening = false
}

func (fb *_FastBattleState) OnLoop(delta time.Duration) {
	if time.Now().After(fb.endTime) {
		fb.team.SetState(TeamState_Running, fb.team.fastBattle.finishRoadIndex)
		fb.team.Error("快速战斗超时")
	}
}

func (fb *_FastBattleState) CanChangeTo(state uint8) bool {
	return state == TeamState_Running
}
