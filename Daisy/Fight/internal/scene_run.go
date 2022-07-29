package internal

import (
	"Daisy/Proto"
)

// Run 运行战斗逻辑
func (scene *Scene) Run() *Proto.FightResult {
	// 战斗结果
	fightResult := &Proto.FightResult{}

	// 添加出生buff
	for _, pawn := range scene.pawnList {
		pawn.BatchAddBuffs(pawn, pawn.Info.InnerBornBuffs, 0)
		pawn.BatchAddBuffs(pawn, pawn.Info.BornBuffs, 0)
	}

	// 开始运行战斗逻辑
	for i := scene.nowFrames; i <= scene.maxFrames; i++ {
		if !scene.update() {
			break
		}
	}

	// 插入结束帧
	scene.PushAction(&Proto.FightEnd{})

	// 战斗继承数据
	fightResult.Inherit = &Proto.FightInherit{}
	fightResult.Inherit.PawnInheritMap = make(map[string]*Proto.PawnInherit)

	// 记录pawn信息
	for _, pawn := range scene.pawnList {
		if pawn.IsRole() {
			fightResult.Inherit.PawnInheritMap[pawn.Info.Role.RoleId] = &Proto.PawnInherit{
				UltimateSkillPower: pawn.Attr.UltimateSkillPower,
			}
		}

		pawn.Info.PawnInfo.FightEndDead = !pawn.IsAlive()

		fightResult.PawnInfos = append(fightResult.PawnInfos, pawn.Info.PawnInfo.PawnInfo)
	}

	// 战斗结果
	fightResult.Outcome = []Proto.Camp_Enum{scene.winCamp, GetEnemyCamp(scene.winCamp)}

	// 记录回放
	fightResult.Time = scene.framesToTime(scene.nowFrames)
	fightResult.RealTime = scene.fightEndTime - scene.fightBeginTime
	fightResult.Replay = &scene.replay

	// 模拟器记录额外信息
	if scene.SimulatorMode() {
		fightResult.DebugSceneInfo = &Proto.DebugSceneInfo{}
		for _, val := range scene.Info.BoundaryPoints {
			fightResult.DebugSceneInfo.BoundaryPoints = append(fightResult.DebugSceneInfo.BoundaryPoints, &Proto.Position{X: val.X, Y: val.Y})
		}
		fightResult.DebugReplay = &scene.debugReplay
	}

	return fightResult
}

// update 更新战斗帧
func (scene *Scene) update() bool {
	// 检查帧进度
	if scene.secFrames <= 0 || scene.nowFrames > scene.maxFrames {
		return false
	}

	switch scene.stage {
	case SceneStage_InitFight:
		scene.inInitFight()
	case SceneStage_WaitFight:
		scene.inWaitFight()
	case SceneStage_Fight:
		scene.inFight()
	case SceneStage_WaitEnd:
		scene.inWaitEnd()
	}

	// 增加帧进度与时间
	scene.nowFrames++
	scene.NowTime = scene.framesToTime(scene.nowFrames)

	return scene.stage != SceneStage_End
}
