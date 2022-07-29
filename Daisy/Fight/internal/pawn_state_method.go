package internal

import (
	"Daisy/Proto"
	"fmt"
)

// Dead 死亡
func (state *FightState) Dead(attack *Attack) {
	if !state.pawn.IsAlive() {
		return
	}

	// 设置死亡标记
	state.ChangeStat(Stat_Death, true)

	// 合体必杀技灭灯
	state.pawn.Scene.TurnOffCombineSkillPoint(state.pawn)

	//移除对应刚体
	state.pawn.Scene.destroyPawnShape(state.pawn.UID)

	// 记录debug信息
	if attack != nil {
		state.pawn.Scene.PushDebugInfo(func() string {
			return fmt.Sprintf("${PawnID:%d}击杀${PawnID:%d}",
				attack.Caster.UID,
				state.pawn.UID)
		})
	}

	state.pawn.Scene.PushDebugInfo(func() string {
		return fmt.Sprintf("${PawnID:%d}死亡%s",
			state.pawn.UID,
			func() string {
				if attack != nil {
					return fmt.Sprintf("，被${PawnID:%d}击杀", attack.Caster.UID)
				}
				return ""
			}())
	})

	// 发送击杀成功事件
	if attack != nil {
		attack.Caster.Events.EmitKillTarget(attack, state.pawn)
	}

	// 发送死亡事件
	state.pawn.Events.EmitDead(attack)

	// 打断当前技能
	state.pawn.BreakCurSkill(state.pawn, Proto.SkillBreakReason_Normal)

	// 删除全部buff
	for key, buff := range state.pawn.buffTab {
		if buff.Config.DeadNoDestroy {
			continue
		}

		state.pawn.RemoveBuff(key)
	}

	state.pawn.Scene.RemoveHaloMember(state.pawn)
}
