package internal

import (
	"Daisy/Proto"
	"fmt"
)

// changeStage 切换场景阶段
func (scene *Scene) changeStage(stage SceneStage) {
	switch stage {
	case SceneStage_InitFight:
	case SceneStage_WaitFight:
	case SceneStage_Fight:
		// 记录战斗开始时间
		scene.fightBeginTime = scene.NowTime

		// 计算狂暴时间
		for _, formation := range scene.formationList {
			formation.RageBeginTime = scene.NowTime + formation.Info.RageTime
		}

		for _, formation := range scene.formationList {
			for _, pawn := range formation.PawnList {
				if !pawn.IsRole() {
					continue
				}

				superSkillList := pawn.GetSuperSkillList()
				for _, skillItem := range superSkillList {
					skillItem.cdEndTime = scene.NowTime + skillItem.Config.CoolDown
				}
			}
		}

		// 插入开始帧
		scene.PushAction(&Proto.FightBegin{})

	case SceneStage_WaitEnd:
		// 记录战斗结束时间
		scene.fightEndTime = scene.NowTime

		// 暂停所有AI
		scene.AllAIPause(true)

		// 遍历所有pawn，打断不在后摇阶段的技能
		for _, pawn := range scene.pawnList {
			if pawn.curSkill != nil && pawn.curSkill.Stat == Proto.SkillState_Later {
				pawn.BreakCurSkill(pawn, Proto.SkillBreakReason_Normal)
			}
		}

	case SceneStage_End:
	}

	scene.stage = stage
}

// inInitFight 初始化战斗阶段
func (scene *Scene) inInitFight() {
	// 切换至准备开战阶段
	scene.changeStage(SceneStage_WaitFight)
}

// inWaitFight 等待开战阶段
func (scene *Scene) inWaitFight() {
	// 切换至战斗阶段
	scene.changeStage(SceneStage_Fight)
}

// inFight 战斗阶段
func (scene *Scene) inFight() {
	//更新行为树流程
	scene._BehaviorFlow.update()

	//更新移动控制器流程
	scene._MovementFlow.update()

	// 更新技能流程
	scene._SkillFlow.update()

	// 更新buff流程
	scene._BuffFlow.update()

	// 更新合体必杀技流程
	scene._CombineSkillFlow.update()

	// 更新受击流程
	scene._BeHitFlow.update()

	// 更新伤害体流程
	scene._AttackFlow.update()

	// 召唤物消失流程
	for _, pawn := range scene.pawnList {
		// 检测生存时间
		if pawn.Info.LifeTime > 0 {
			if pawn.Scene.NowTime >= pawn.DestroyTime {
				if pawn.IsNpc() && pawn.Master != nil {
					// 记录debug信息
					pawn.Scene.PushDebugInfo(func() string {
						return fmt.Sprintf("${PawnID:%d}的召唤物${PawnID:%d}消失了",
							pawn.Master.UID,
							pawn.UID)
					})
				}

				pawn.State.Dead(nil)
				return
			}
		}
	}

	// 狂暴buff流程
	for _, formation := range scene.formationList {
		// 开始狂暴
		if scene.NowTime >= formation.RageBeginTime && !formation.Raged {
			bossFound := false
			for _, pawn := range formation.PawnList {
				if pawn.IsBoss() && pawn.IsAlive() {
					bossFound = true
				}
			}

			if bossFound {
				for _, pawn := range formation.PawnList {
					pawn.State.ChangeStat(Stat_Raged, true)
				}
			}

			formation.Raged = true
		}
	}

	// 检测是否是队伍中最后一个boss
	for _, pawn := range scene.pawnList {
		if !pawn.IsAlive() && pawn.IsBoss() {
			formation := pawn.Scene.formationList[pawn.GetCamp()]

			// 检测是否有boss存活
			for _, friend := range formation.PawnList {
				if friend.IsBoss() && friend.IsAlive() {
					break
				}
			}

			// 队伍中最后一个boss死亡后，将己方所有npc致死
			for _, friend := range formation.PawnList {
				friend.State.Dead(nil)
			}
		}
	}

	// 帧末删除废弃action
	for _, pawn := range scene.pawnList {
		pawn.delDiscardMoveEndAction()
	}

	// 检测战斗结果
	for camp, formation := range scene.formationList {
		hasPawnAlive := false

		// 是否全部死亡
		for _, pawn := range formation.PawnList {
			if pawn.IsBackground() {
				continue
			}

			if pawn.IsAlive() {
				hasPawnAlive = true
				break
			}
		}

		// 全部死亡
		if !hasPawnAlive {
			scene.winCamp = GetEnemyCamp(Proto.Camp_Enum(camp))

			// 切换至等待结束阶段
			scene.changeStage(SceneStage_WaitEnd)
			break
		}
	}
}

// inWaitEnd 等待结束阶段
func (scene *Scene) inWaitEnd() {
	//更新行为树流程
	scene._BehaviorFlow.update()

	//更新移动控制器流程
	scene._MovementFlow.update()

	// 更新技能流程
	scene._SkillFlow.update()

	// 更新buff流程
	scene._BuffFlow.update()

	// 更新合体必杀技流程
	scene._CombineSkillFlow.update()

	// 更新受击流程
	scene._BeHitFlow.update()

	// 更新伤害体流程
	scene._AttackFlow.update()

	// 等待所有技能流程结束
	if scene.isSkillFlowRunning() {
		return
	}

	// 切换至结束阶段
	scene.changeStage(SceneStage_End)
}
