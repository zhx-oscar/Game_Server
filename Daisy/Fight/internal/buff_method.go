package internal

import "Daisy/Fight/internal/conf"

// Destroy 销毁buff自身
func (buff *Buff) Destroy() {
	if buff.IsDestroy {
		return
	}

	buff.Pawn.RemoveBuff(buff.BuffKey)
}

// RefreshDuration 刷新自身时长
func (buff *Buff) RefreshDuration() bool {
	return buff.Pawn.Scene.refreshBuffDurationEx(buff.Pawn, buff)
}

// ExtendDuration 延长自身时长
func (buff *Buff) ExtendDuration(duration uint32) bool {
	return buff.Pawn.Scene.extendBuffDurationEx(buff.Pawn, buff, duration)
}

// SetDestroyTime 设置销毁时间
func (buff *Buff) SetDestroyTime(time uint32) bool {
	return buff.Pawn.Scene.setBuffDestroyTime(buff.Pawn, buff, time)
}

// RemainTime 剩余时长
func (buff *Buff) RemainTime() uint32 {
	if buff.IsDestroy {
		return 0
	}

	durationTypeFound := false
	for _, disappearType := range buff.Config.DisappearType {
		if conf.BuffDisappearType(disappearType) == conf.BuffDisappearType_Duration {
			durationTypeFound = true
			break
		}
	}

	if !durationTypeFound {
		return 0
	}

	if buff.destroyTime < buff.Pawn.Scene.NowTime {
		return 0
	}

	return buff.destroyTime - buff.Pawn.Scene.NowTime
}
