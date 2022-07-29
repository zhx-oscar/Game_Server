package main

import (
	"Cinder/Base/linemath"
	"Cinder/plugin/navmesh"
	"Cinder/plugin/physxgo"
	"Daisy/Battle/sceneconfigdata"
	"Daisy/Data"
	"Daisy/Proto"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

// 移动状态定义
const (
	RunStateIdle uint8 = iota
	RunStateMove
	RunStateJump
)

type _RunState struct {
	team *_Team

	//move
	targetPaths  []linemath.Vector3
	location     linemath.Vector3
	moveLeftTime float32
	velocity     float32

	//jump
	jumpTargetLocation linemath.Vector3
	jumpVelocity       float32

	state uint8

	mapLoaded bool

	curTrigger physxgo.IPxActor
}

func NewRunState(team *_Team) *_RunState {
	run := &_RunState{
		team: team,
	}

	battleAreaCfg, ok := Data.GetSceneConfig().BattleArea_ConfigItems[team.prop.Data.Raid.Progress]
	if !ok {
		panic(fmt.Sprintf("can't find battleArea config progress=%d", team.prop.Data.Raid.Progress))
	}
	run.init(fmt.Sprintf("%d", battleAreaCfg.MapID))

	return run
}

func (run *_RunState) OnEnter(args ...interface{}) {
	roadIdx := uint32(0)
	if len(args) > 0 {
		roadIdx = args[0].(uint32)
	}

	if roadIdx != 0 {
		run.setPosByRoadIndex(roadIdx)
	}
	run.changeRunState(RunStateMove)
}

func (run *_RunState) OnLeave() {
	run.changeRunState(RunStateIdle)
}

func (run *_RunState) OnLoop(deltaTime time.Duration) {
	if !run.mapLoaded {
		return
	}

	//todo 切换地图逻辑

	//GM调试
	if run.team.debugRoadIndex >= 0 {
		run.state = RunStateMove
		run.setPosByRoadIndex(uint32(run.team.debugRoadIndex))

		run.team.debugRoadIndex = -1
	}

	if run.state == RunStateMove {
		run.onStateMove(float32(deltaTime.Seconds()))
	} else if run.state == RunStateJump {
		run.onStateJump(float32(deltaTime.Seconds()))
	}
}

func (run *_RunState) CanChangeTo(state uint8) bool {
	return run.team.runChest == nil
}

/*
 队伍跑图逻辑
*/

func (run *_RunState) init(mapID string) error {
	err := run.loadMap(mapID)
	if err != nil {
		run.team.Errorf("team.run init err:%s", err)
		return err
	}

	run.velocity = run.team.paths.LeaderSpeed
	return nil
}

func (run *_RunState) reload(mapID string) error {
	run.unloadMap()
	err := run.loadMap(mapID)
	if err != nil {
		run.team.Errorf("team.run reload err:%s", err)
		return err
	}

	run.velocity = run.team.paths.LeaderSpeed
	return nil
}

// 卸载当前地图
func (run *_RunState) unloadMap() {
	if run.team.query != nil {
		navmesh.DestroyQuery(run.team.query)
		run.team.query = nil
	}
	run.mapLoaded = false
}

// 根据当前队伍进度, 加载相应地图信息
func (run *_RunState) loadMap(mapID string) error {
	paths, err := sceneconfigdata.LoadPath(fmt.Sprintf("../res/MapData/%s/path.json", mapID))
	if err != nil {
		return err
	}
	run.team.paths = paths

	//navmesh
	query := navmesh.CreateQuery(fmt.Sprintf("../res/MapData/%s/navmesh.bin", mapID), 2048)
	if query == nil {
		return errors.New("carate navmesh query fail")
	}
	run.team.query = query

	//physics
	run.team.pxScene, err = LoadPxScene(mapID)
	if err != nil {
		return err
	}

	run.mapLoaded = true

	return nil
}

// 设置位置到指定的路点
func (run *_RunState) setPosByRoadIndex(idx uint32) {
	run.location = run.team.paths.Points[idx]
	nexts := run.team.GetNextRoadIndexs(idx)
	if len(nexts) == 0 {
		panic(fmt.Sprintf("idx:%d 没有配置下一个路点", idx))
	}
	if len(nexts) > 1 {
		run.team.Errorf("next应该只有一个")
	}
	run.setTarget(nexts[0])
}

func (run *_RunState) findPathToTarget(origin, target linemath.Vector3) ([]linemath.Vector3, error) {
	//convert to y-up
	origin = sceneconfigdata.UnconvertVector3(origin)

	//convert to y-up
	target = sceneconfigdata.UnconvertVector3(target)

	paths, finded, err := run.team.query.FindPath(origin, target)
	if err != nil {
		return nil, err
	} else if !finded {
		return nil, errors.New("not found")
	} else {
		for i := 0; i < len(paths); i++ {
			//convert to z-up
			paths[i] = sceneconfigdata.ConvertVector3(paths[i])
		}
		return paths, nil
	}
}

func (run *_RunState) changeRunState(state uint8) {
	if run.state != state {
		run.state = state
	}
}

func randomFork(separate []uint32) uint32 {
	if len(separate) <= 1 {
		return 0
	}
	total := 0
	for _, v := range separate {
		total += int(v)
	}

	if total == 0 {
		return 0
	}

	ran := rand.Intn(total)
	cur := 0
	for i := 0; i < len(separate); i++ {
		cur += int(separate[i])
		if ran < cur {
			return uint32(i)
		}
	}
	return 0
}

func (run *_RunState) onEnterTrigger(trigger *sceneconfigdata.TriggerCfg) {
	triggerParam := trigger.TriggerParam.Positive
	if !run.team.positive {
		triggerParam = trigger.TriggerParam.Opposite
	}

	//run.team.Debugf("onEnterTrigger type:%d ID:%s", trigger.TriggerType, trigger.ID)
	switch trigger.TriggerType {
	case sceneconfigdata.TriggerType_TeamFork:
		param := triggerParam.(*sceneconfigdata.TriggerParamTeamFork)
		if param.Separate == nil || len(param.Separate) == 0 { //不生效
			return
		}
		ran := randomFork(param.Separate)
		nexts := run.team.GetNextRoadIndexs(run.team.targetRoadIndex)
		if len(nexts) == 0 {
			panic(fmt.Sprintf("idx:%d 没有配置下一个路点", run.team.targetRoadIndex))
		}
		if ran >= uint32(len(nexts)) {
			run.team.Errorf("TeamFork的配置有误，选择第1个")
			ran = 0
		}
		run.setTarget(nexts[ran])

	case sceneconfigdata.TriggerType_Battle:
		param := triggerParam.(*sceneconfigdata.TriggerParamBattle)
		if param.BattleID == 0 { //不生效
			return
		}
		if run.team.runChest != nil {
			run.team.Errorf("宝箱过程中触发了新的战斗，找场景策划改配置 TriggerID:%s", trigger.ID)
		} else {
			run.team.Infof("onEnterTrigger Battle TriggerID:%s", trigger.ID)
			run.team.onTriggerRaidBattle(false, param.BattleID)
		}

	case sceneconfigdata.TriggerType_Speed:
		param := triggerParam.(*sceneconfigdata.TriggerParamSpeed)
		if param.Speed == 0 { //不生效
			return
		}
		if run.velocity != param.Speed {
			run.setTargetLocation(run.targetPaths[0], param.Speed)
		}

	case sceneconfigdata.TriggerType_Jump:
		param := triggerParam.(*sceneconfigdata.TriggerParamJump)
		if !param.IsIn { //不生效
			return
		}
		run.jumpTargetLocation = param.OutLoc
		run.jumpVelocity = trigger.TriggerParam.JumpSpeed
		//run.team.Debugf("start jump targetLoc:%v, speed:%v", run.jumpTargetLocation, run.jumpVelocity)
		run.changeRunState(RunStateJump)

		run.setTargetLocation(run.jumpTargetLocation, run.velocity)

	case sceneconfigdata.TriggerType_Chest:
		if run.team.runChest != nil {
			run.team.Errorf("宝箱过程中触发了新的宝箱，找场景策划改配置 TriggerID:%s", trigger.ID)
		} else {
			param := triggerParam.(*sceneconfigdata.TriggerParamChest)
			run.team.Infof("onEnterTrigger Chest TriggerID:%s", trigger.ID)
			run.team.onTriggerChest(param)
		}

	case sceneconfigdata.TriggerType_SpawnMonster:
		param := triggerParam.(*sceneconfigdata.TriggerParamSpawnMonster)
		run.team.Infof("onTriggerSpawnMonster TriggerID:%s BattleFieldID:%d", trigger.ID, param.BattleID)
		run.team.onTriggerSpawnMonster(param)

	}
}

func (run *_RunState) setTarget(target uint32) {
	//run.team.Debugf(">>>>>set target road index:%v", target)
	run.team.targetRoadIndex = target
	to := run.team.paths.Points[target]
	run.targetPaths = []linemath.Vector3{to}
	run.setTargetLocation(to, run.velocity)
}

func (run *_RunState) onStateMove(delta float32) {
	if len(run.targetPaths) == 0 {
		run.team.Errorf("配置有误，没有目标路点无法继续移动")
		return
	}

	targetPosChanged := false
	location := run.location
	to := run.targetPaths[0]

	t := to.Sub(location).Len() / run.velocity
	run.moveLeftTime += delta
	if t >= run.moveLeftTime {
		direction := to.Sub(location).Normalized()
		length := run.moveLeftTime * run.velocity
		run.location = location.Add(direction.Mul(length))

		run.moveLeftTime = 0
	} else {
		targetPosChanged = true
		run.location = to
		run.targetPaths = run.targetPaths[1:]

		run.moveLeftTime -= t
	}

	//是否切换下一个路点
	if len(run.targetPaths) == 0 {
		nexts := run.team.GetNextRoadIndexs(run.team.targetRoadIndex)
		if len(nexts) == 0 {
			panic(fmt.Sprintf("idx:%d 没有配置下一个路点", run.team.targetRoadIndex))
		}
		run.setTarget(nexts[0])
	} else if targetPosChanged {
		run.setTargetLocation(run.targetPaths[0], run.velocity)
	}

	b, trigger := run.checkTrigger(location, run.location)
	if b {
		run.onEnterTrigger(trigger.GetUserData().(*sceneconfigdata.TriggerCfg))
	}
}

func (run *_RunState) onStateJump(delta float32) {
	location := run.location
	to := run.jumpTargetLocation

	t := to.Sub(location).Len() / run.jumpVelocity
	if t >= delta {
		direction := to.Sub(location).Normalized()
		length := delta * run.jumpVelocity
		run.location = location.Add(direction.Mul(length))
	} else {
		run.location = to

		//跳之后要寻下一个目标点
		//run.team.Debugf("---------jump old targetRoadIndex:%d", run.targetRoadIndex)
		nexts := run.team.GetNextRoadIndexs(run.team.targetRoadIndex)
		if len(nexts) == 0 {
			panic(fmt.Sprintf("idx:%d 没有配置下一个路点", run.team.targetRoadIndex))
		}
		run.setTarget(nexts[0])
		//run.team.Debugf("---------jump new targetRoadIndex:%d", run.targetRoadIndex)

		run.changeRunState(RunStateMove)
	}
}

func (run *_RunState) checkTrigger(origin, target linemath.Vector3) (bool, physxgo.IPxActor) {
	geo := &physxgo.Geometry{
		Type:       physxgo.GeomType_eCAPSULE,
		Radius:     0.3,
		HalfHeight: 0.9,
	}

	var hitActor physxgo.IPxActor
	if origin.IsEqual(target) {
		hit, err := run.team.pxScene.OverlapOne(*geo,
			physxgo.TransForm{P: origin, Q: linemath.Quaternion{X: 0, Y: 0, Z: 0, W: 1}},
			physxgo.HitMode_eHitStatic,
			physxgo.NoHitFilter)

		if err != nil {
			return false, nil
		}

		if hit != nil {
			hitActor = hit.Target
		}
	} else {
		direction := target.Sub(origin).Normalized()
		distance := direction.Len()

		hit, err := run.team.pxScene.SweepOne(*geo,
			physxgo.TransForm{P: origin, Q: linemath.Quaternion{X: 0, Y: 0, Z: 0, W: 1}},
			direction,
			distance,
			0,
			physxgo.HitMode_eHitStatic,
			physxgo.NoHitFilter)

		if err != nil {
			return false, nil
		}

		if hit != nil {
			hitActor = hit.Target
		}
	}

	if hitActor == nil {
		run.curTrigger = nil
		return false, nil
	} else {
		if run.curTrigger == nil {
			run.curTrigger = hitActor
			return true, run.curTrigger
		} else {
			if hitActor == run.curTrigger {
				return false, nil
			} else {
				run.curTrigger = hitActor
				return true, run.curTrigger
			}
		}
	}
}

func (run *_RunState) setTargetLocation(targetLocation linemath.Vector3, velocity float32) {
	run.velocity = velocity
	//run.team.Debugf(">>>>>target pos:%v, veloctiy:%v", targetLocation, velocity)
	//conver to y-up
	curLoc := sceneconfigdata.UnconvertVector3(run.location)
	targetLoc := sceneconfigdata.UnconvertVector3(targetLocation)
	run.team.onSetLocation(run.linemath2ProtoVector3(curLoc), run.linemath2ProtoVector3(targetLoc), velocity, run.team.targetRoadIndex)
}

func (run *_RunState) linemath2ProtoVector3(vec linemath.Vector3) *Proto.PVector3 {
	return &Proto.PVector3{
		X: vec.X,
		Y: vec.Y,
		Z: vec.Z,
	}
}
