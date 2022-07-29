package effects

import (
	. "Daisy/Fight/internal"
)

// GetAttackLogicData 获取伤害体逻辑层数据
func GetAttackLogicData(attack *Attack) *AttackLogicData {
	if attack.LogicData == nil {
		attack.LogicData = &AttackLogicData{}
	}
	return attack.LogicData.(*AttackLogicData)
}

// AttackLogicData 伤害体逻辑层数据
type AttackLogicData struct {
	FirstHitDamageBit      Bits // 首次hit伤害类型
	firstHitDamageBitSaved bool // 是否已记录首次hit伤害类型
}

// SetFirstHitDamageBit 设置首次hit伤害类型
func (data *AttackLogicData) SetFirstHitDamageBit(damageBit Bits) {
	if data.firstHitDamageBitSaved {
		return
	}

	data.FirstHitDamageBit = damageBit
	data.firstHitDamageBitSaved = true
}
