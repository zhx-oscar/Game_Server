package internal

import (
	"Daisy/Proto"
	"fmt"
	"sort"
)

// HPShield 护盾
type HPShield struct {
	UID      uint32          // uid
	ShieldHP int64           // 护盾血量
	Weight   uint32          // 权重，数值越大优先级越高
	effect   IEffectCallback // 绑定的效果
}

// AddHPShield 添加护盾
func (attr *FightAttr) AddHPShield(shieldHP int64, weight uint32, effect IEffectCallback) *HPShield {
	if !attr.inFight || attr.pawn.isSnapshot {
		return nil
	}

	shield := &HPShield{
		UID:      attr.pawn.Scene.generateUID(),
		ShieldHP: shieldHP,
		Weight:   weight,
		effect:   effect,
	}

	attr.HPShieldList = append(attr.HPShieldList, shield)

	sort.SliceStable(attr.HPShieldList, func(i, j int) bool {
		return attr.HPShieldList[i].Weight > attr.HPShieldList[j].Weight
	})

	old := attr.AllHPShield
	attr.AllHPShield += shieldHP

	attr.pawn.Scene.PushDebugInfo(func() string {
		return fmt.Sprintf("${PawnID:%d}增加护盾%d，护盾HP：%d，权重：%d", attr.pawn.UID, shield.UID, shieldHP, weight)
	})

	// 记录回放
	attr.AttrSyncToClient(Proto.AttrType_HPShield, float64(old), float64(attr.AllHPShield))

	return shield
}

// RemoveHPShield 删除护盾
func (attr *FightAttr) RemoveHPShield(shieldUID uint32) {
	if !attr.inFight || attr.pawn.isSnapshot {
		return
	}

	var shield *HPShield

	for i := 0; i < len(attr.HPShieldList); i++ {
		if attr.HPShieldList[i].UID == shieldUID {
			shield = attr.HPShieldList[i]
			attr.HPShieldList = append(attr.HPShieldList[0:i], attr.HPShieldList[i+1:]...)
			break
		}
	}

	if shield == nil {
		return
	}

	oldAllHPShield := attr.AllHPShield
	attr.AllHPShield -= shield.ShieldHP

	attr.pawn.Scene.PushDebugInfo(func() string {
		return fmt.Sprintf("${PawnID:%d}删除护盾%d，护盾HP：%d，权重：%d", attr.pawn.UID, shield.UID, shield.ShieldHP, shield.Weight)
	})

	// 记录回放
	attr.AttrSyncToClient(Proto.AttrType_HPShield, float64(oldAllHPShield), float64(attr.AllHPShield))
}

// ChangeHPShield 修改护盾HP
func (attr *FightAttr) ChangeHPShield(shieldUID uint32, deltaHP int64) {
	if !attr.inFight || attr.pawn.isSnapshot {
		return
	}

	var shield *HPShield

	for i := 0; i < len(attr.HPShieldList); i++ {
		if attr.HPShieldList[i].UID == shieldUID {
			shield = attr.HPShieldList[i]
			break
		}
	}

	if shield == nil {
		return
	}

	oldShieldHP := shield.ShieldHP
	shield.ShieldHP += deltaHP
	if shield.ShieldHP < 0 {
		shield.ShieldHP = 0
	}

	deltaHP = shield.ShieldHP - oldShieldHP

	oldAllHPShield := attr.AllHPShield
	attr.AllHPShield += deltaHP

	attr.pawn.Scene.PushDebugInfo(func() string {
		return fmt.Sprintf("${PawnID:%d}修改护盾%d，护盾HP变化：%d => %d", attr.pawn.UID, shield.UID, oldShieldHP, shield.ShieldHP)
	})

	// 记录回放
	attr.AttrSyncToClient(Proto.AttrType_HPShield, float64(oldAllHPShield), float64(attr.AllHPShield))

	return
}

// GetHPShield 查询护盾
func (attr *FightAttr) GetHPShield(shieldUID uint32) *HPShield {
	if !attr.inFight || attr.pawn.isSnapshot {
		return nil
	}

	for i := 0; i < len(attr.HPShieldList); i++ {
		if attr.HPShieldList[i].UID == shieldUID {
			return attr.HPShieldList[i]
		}
	}

	return nil
}
