package internal

import "Daisy/Proto"

// _BeHitFlow 受击流程
type _BeHitFlow struct {
	scene *Scene
}

// init 初始化
func (flow *_BeHitFlow) init(scene *Scene) {
	flow.scene = scene
}

// update 帧更新
func (flow *_BeHitFlow) update() {
	for _, pawn := range flow.scene.pawnList {
		if !pawn.IsAlive() {
			continue
		}

		if pawn.State.Dodging && pawn.BeHit.DodgeEndTime <= pawn.Scene.NowTime {
			pawn.State.ChangeStat(Stat_Dodging, false)

			if pawn.IsMoving() {
				pawn.Stop()
			}

			pawn.State.ChangeStat(Stat_Invincible, false)
			pawn.AIPause(false)
			pawn.AIBackToRoot()

			continue
		}

		if pawn.State.BeHitStat == 0 {
			continue
		}

		for hitStat := Proto.HitType_Hit; hitStat <= Proto.HitType_BlockBreak; hitStat++ {
			if !pawn.State.BeHitStat.Test(int32(hitStat)) {
				continue
			}

			if pawn.BeHit.BeHitEndTime[hitStat] > flow.scene.NowTime {
				continue
			}

			pawn.State.ChangeBeHitStatBit(hitStat, false)
		}
	}
}
