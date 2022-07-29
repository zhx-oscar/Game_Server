package buffeffect

import (
	. "Daisy/Fight/internal"
	. "Daisy/Fight/internal/effects"
	"fmt"
)

type _1005_RecoverHP struct {
	Blank
	lastTime uint32
}

// OnBuffUpdate buff帧更新（buff自身能收到）
func (effect *_1005_RecoverHP) OnBuffUpdate(buff *Buff) {
	pawn := buff.Pawn

	// 检测时间间隔
	if pawn.Scene.NowTime < effect.lastTime+1000 {
		return
	}
	effect.lastTime = pawn.Scene.NowTime

	// 检测恢复血量
	if pawn.Attr.RecoverHP <= 0 || !pawn.IsAlive() {
		return
	}

	oldHP := pawn.Attr.CurHP

	// 调整血量
	pawn.Attr.ChangeHP(pawn.Attr.CurHP + pawn.Attr.RecoverHP)

	pawn.Scene.PushDebugInfo(func() string {
		if oldHP != pawn.Attr.CurHP {
			return fmt.Sprintf("${PawnID:%d}每秒恢复HP${DamageHP:%d}，当前HP变化：%d => %d", pawn.UID, pawn.Attr.RecoverHP, oldHP, pawn.Attr.CurHP)
		}
		return ""
	})
}
