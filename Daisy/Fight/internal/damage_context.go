package internal

import (
	. "Daisy/Fight/internal/conf"
	"fmt"
	"math"
)

// DamageContext 伤害上下文
type DamageContext struct {
	DamageKind           DamageKind                   // 伤害类型
	DamageFlow           DamageFlow                   // 伤害流程
	AttackLucky          bool                         // 攻击幸运
	DamageValueTab       [DamageValueKind_End]float64 // 伤害值表
	DamageBit            Bits                         // 伤害bit
	DamageHP             int64                        // 伤害HP值
	DamageHPShield       int64                        // 伤害HP护盾值
	BloodsuckerValue     int64                        // 吸血值
	BloodsuckerRecoverHP int64                        // 吸血实际恢复HP
	Attack               *Attack                      // 攻方伤害体
	Target               *Pawn                        // 攻击目标
	ExtendValue          float64                      // 外部数值
	StepIndex            int32                        // 步骤索引
	debugInfo            string                       // 调试信息
}

// Init 初始化
func (damageCtx *DamageContext) Init(damageKind DamageKind, damageFlow DamageFlow, attack *Attack, target *Pawn, extendValue float64) error {
	if attack == nil || target == nil {
		return fmt.Errorf("args invalid")
	}

	damageCtx.DamageKind = damageKind
	damageCtx.DamageFlow = damageFlow
	damageCtx.Attack = attack
	damageCtx.Target = target
	damageCtx.ExtendValue = extendValue

	return nil
}

// DamageValue 总伤害值
func (damageCtx *DamageContext) DamageValue() (damageValue float64) {
	for i := DamageValueKind_Begin; i < DamageValueKind_End; i++ {
		damageValue += damageCtx.DamageValueTab[i]
	}

	if FloatEqual(damageValue, 0) {
		damageValue = 0
	}

	damageValue = math.Abs(damageValue)

	return
}

// DamageValueWithInvert 返回总伤害值和伤害值是否反转
func (damageCtx *DamageContext) DamageValueWithInvert() (damageValue float64, damageInvert bool) {
	for i := DamageValueKind_Begin; i < DamageValueKind_End; i++ {
		damageValue += damageCtx.DamageValueTab[i]
	}

	if FloatEqual(damageValue, 0) {
		damageValue = 0
	}

	damageInvert = damageValue < 0

	damageValue = math.Abs(damageValue)

	return
}

// PushDebugInfo 记录debug信息
func (damageCtx *DamageContext) PushDebugInfo(fun func() string) {
	if !damageCtx.Attack.Caster.Scene.SimulatorMode() && !damageCtx.Attack.Caster.Scene.TestMode() {
		return
	}
	damageCtx.debugInfo += fun()
}

// FlushDebugInfo 刷新debug信息
func (damageCtx *DamageContext) FlushDebugInfo() {
	if !damageCtx.Attack.Caster.Scene.SimulatorMode() && !damageCtx.Attack.Caster.Scene.TestMode() {
		return
	}
	damageCtx.Attack.Caster.Scene.PushDetailDebugInfo(func() string {
		return damageCtx.debugInfo
	})
}
