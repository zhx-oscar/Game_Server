package main

import (
	"Cinder/Base/Const"
	"Cinder/Base/User"
	"Cinder/Space"
	"Daisy/Data"
	"Daisy/ErrorCode"
	"Daisy/Prop"
	"Daisy/Proto"
	"errors"
	"time"
)

type _FastBattleModel struct {
	userBattleEnd  map[string]bool
	opening        bool
	enjoyableRoles map[string]bool

	finishRoadIndex uint32
}

func NewFastBattleModel() *_FastBattleModel {
	return &_FastBattleModel{}
}

func (team *_Team) ResetFastBattle(resetTime time.Time) {
	team.Info("ResetFastBattle ", resetTime)
	fb := &Proto.FastBattle{
		StageInfo: make(map[uint32]*Proto.FastBattleStageInfo),
	}
	for key, value := range Data.GetFastBattleConfig().FastBattle_ConfigItems {
		fb.StageInfo[key] = &Proto.FastBattleStageInfo{
			MaxMyTimes:    value.MyTimes,
			MaxOtherTimes: value.OtherTimes,
		}
	}

	team.TraversalActor(func(actor Space.IActor) {
		role := actor.(*_Role)
		role.prop.SyncSetFastBattle(fb)
	})

	team.prop.SyncSetFastBattleResetTimestamp(resetTime.Unix())
}

func (team *_Team) StartFastBattle(roleID string, stage uint32, multi float32) int32 {
	if team.fastBattle.opening {
		team.Error("StartFastBattle AlreadySpeedUp")
		return ErrorCode.AlreadySpeedUp
	}

	team.fastBattle.opening = true
	team.fastBattle.userBattleEnd = make(map[string]bool)
	team.fastBattle.enjoyableRoles = map[string]bool{roleID: true}

	team.TraversalActor(func(ia Space.IActor) {
		if _, ok := team.prop.Data.Base.Members[ia.GetID()]; ok && ia.GetID() != roleID {
			roleProp := ia.GetProp().(*Prop.RoleProp)
			if info, ok := roleProp.Data.FastBattle.StageInfo[stage]; ok {
				if info.OtherTimes < info.MaxOtherTimes {
					team.fastBattle.enjoyableRoles[ia.GetID()] = true
				} else {
					team.Infof("StartFastBattle role:%s 无法享受 stage:%d 加速", ia.GetID(), stage)
				}
			} else {
				team.Errorf("StartFastBattle role:%s state:%d not exist", ia.GetID(), stage)
			}
		}
	})

	team.PendingStateEvent = append(team.PendingStateEvent, &_PendingStateEvent{
		State: TeamState_FastBattleing,
		Event: func() {
			team.SetState(TeamState_FastBattleing, multi)
		},
	})

	team.TraversalUser(func(iu User.IUser) bool {
		iu.Rpc(Const.Agent, "RPC_PendingFastBattle", roleID, stage)
		return true
	})

	return ErrorCode.Success
}

func (team *_Team) OnUserFastBattleEnd(userID string) {
	team.fastBattle.userBattleEnd[userID] = true

	finish := true
	team.TraversalUser(func(user User.IUser) bool {
		if _, ok := team.fastBattle.userBattleEnd[user.GetID()]; !ok {
			finish = false
			return false
		}
		return true
	})

	if finish {
		team.SetState(TeamState_Running, team.fastBattle.finishRoadIndex)
	}
}

func (team *_Team) GetFastBattlePassbyRoadIndexs(velocity, time float32) ([]uint32, error) {
	roadIndexs := []uint32{team.targetRoadIndex}
	remainDis := velocity * time

	for {
		if remainDis <= 0 {
			break
		}

		cur := roadIndexs[len(roadIndexs)-1]
		nexts := team.GetNextRoadIndexs(cur)
		if len(nexts) == 0 {
			return nil, errors.New("无法找到下一个路点配置")
		}

		target := nexts[0]
		dis := team.paths.Points[target].Sub(team.paths.Points[cur]).Len()
		remainDis -= dis
		roadIndexs = append(roadIndexs, target)
	}

	return roadIndexs, nil
}
