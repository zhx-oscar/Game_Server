package buffeffect

import (
	. "Daisy/Fight/internal"
	"Daisy/Fight/internal/effects"
	"fmt"
)

type _1001_EnergyShield struct {
	effects.Blank
	buff *Buff

	shield                    *HPShield
	EnergyHPShieldRecoverTime uint32 // 能量血护盾恢复时间
}

func (effect *_1001_EnergyShield) Init(buff *Buff) error {
	effect.buff = buff

	return nil
}

// OnBuffAdd 施加buff后（buff自身能收到）
func (effect *_1001_EnergyShield) OnBuffAdd(buff *Buff) {
	pawn := buff.Pawn
	if !pawn.IsRole() {
		pawn.RemoveBuff(buff.BuffKey)
		return
	}

	if pawn.Info.EnergyShieldCapability <= 0 {
		pawn.RemoveBuff(buff.BuffKey)
		return
	}

	// 能量护盾参数
	energyHPShield := pawn.Attr.PowerShieldHP

	// 护盾最小值为1
	if energyHPShield <= 0 {
		energyHPShield = 1
	}

	effect.shield = pawn.Attr.AddHPShield(energyHPShield, 0, effect)
	pawn.State.ChangeStat(Stat_EnergyShieldOn, true)
}

func (effect *_1001_EnergyShield) OnBuffUpdate(buff *Buff) {
	pawn := effect.buff.Pawn

	if !pawn.IsRole() {
		return
	}

	if pawn.Info.EnergyShieldCapability <= 0 {
		return
	}

	if effect.shield.ShieldHP >= pawn.Attr.PowerShieldHP {
		return
	}

	if effect.EnergyHPShieldRecoverTime > pawn.Scene.NowTime {
		return
	}

	if !effect.buff.Pawn.State.EnergyShieldOn {
		effect.shield = pawn.Attr.AddHPShield(0, 0, effect)
		pawn.State.ChangeStat(Stat_EnergyShieldOn, true)
	}

	oldShieldHP := effect.shield.ShieldHP
	if effect.shield.ShieldHP+int64(float64(pawn.Attr.PowerShieldHP)*float64(pawn.Attr.PowerShieldRecoverSpeed)) > pawn.Attr.PowerShieldHP {
		pawn.Attr.ChangeHPShield(effect.shield.UID, pawn.Attr.PowerShieldHP-oldShieldHP)
	} else {
		pawn.Attr.ChangeHPShield(effect.shield.UID, int64(float64(pawn.Attr.PowerShieldHP)*float64(pawn.Attr.PowerShieldRecoverSpeed)))
	}

	pawn.Scene.PushDebugInfo(func() string {
		if oldShieldHP != effect.shield.ShieldHP {
			return fmt.Sprintf("${PawnID:%d}每秒恢复离子护盾%d，离子护盾变化：%d => %d",
				pawn.UID,
				int64(float64(pawn.Attr.PowerShieldHP)*float64(pawn.Attr.PowerShieldRecoverSpeed)),
				oldShieldHP,
				effect.shield.ShieldHP)
		}
		return ""
	})

	effect.EnergyHPShieldRecoverTime = pawn.Scene.NowTime + 1000
}

// OnBuffRemove buff销毁
func (effect *_1001_EnergyShield) OnBuffRemove(buff *Buff, clear bool) {
	if effect.shield == nil {
		return
	}

	effect.buff.Pawn.Attr.RemoveHPShield(effect.shield.UID)
	effect.buff.Pawn.State.ChangeStat(Stat_EnergyShieldOn, false)
}

// OnShieldBroken 护盾破损
func (effect *_1001_EnergyShield) OnShieldBroken(attack *Attack, shield *HPShield) {
	effect.buff.Pawn.State.ChangeStat(Stat_EnergyShieldOn, false)
}
