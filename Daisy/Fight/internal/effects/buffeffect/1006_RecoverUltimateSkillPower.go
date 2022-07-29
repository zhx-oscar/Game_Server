package buffeffect

import (
	. "Daisy/Fight/internal"
	. "Daisy/Fight/internal/effects"
	"fmt"
)

type _1006_RecoverUltimateSkillPower struct {
	Blank
	lastTime uint32
}

// OnBuffUpdate buff帧更新（buff自身能收到）
func (effect *_1006_RecoverUltimateSkillPower) OnBuffUpdate(buff *Buff) {
	pawn := buff.Pawn

	// 检测时间间隔
	if pawn.Scene.NowTime < effect.lastTime+1000 {
		return
	}
	effect.lastTime = pawn.Scene.NowTime

	// 检测恢复必杀技能量
	recoverPower := int32(pawn.Attr.RecoverUltimateSkillPowerRate * float32(pawn.Attr.SkillPowerLimit))
	if recoverPower <= 0 || !pawn.IsAlive() {
		return
	}

	oldPower := pawn.Attr.UltimateSkillPower

	pawn.Attr.ChangeUltimateSkillPower(pawn.Attr.UltimateSkillPower + recoverPower)

	pawn.Scene.PushDebugInfo(func() string {
		if oldPower != pawn.Attr.UltimateSkillPower {
			return fmt.Sprintf("${PawnID:%d}每秒恢复必杀技能量%d，必杀技能量变化：%d => %d", pawn.UID, recoverPower, oldPower, pawn.Attr.UltimateSkillPower)
		}
		return ""
	})
}
