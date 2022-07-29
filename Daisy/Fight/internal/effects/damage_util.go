package effects

import (
	. "Daisy/Fight/internal/conf"
	"math/rand"
)

// roundTabPrec 圆桌表精度
const roundTabPrec int32 = 1000000

// roundFloat2Int 圆桌表概率数值转换
func roundFloat2Int(value float32) int32 {
	return int32(value * float32(roundTabPrec))
}

// roundFloat2Int 圆桌表概率数值转换
func roundInt2Float(value int32) float32 {
	return float32(value) / float32(roundTabPrec)
}

// _RoundItem 圆桌项
type _RoundItem struct {
	Max, Min int32
}

// _RoundTable 圆桌表
type _RoundTable [DamageJudgeRv_Count]_RoundItem

// random 随机
func (tab *_RoundTable) random() (DamageJudgeRv, int32) {
	var count int32

	for i := 0; i < len(*tab); i++ {
		item := &(*tab)[i]

		if item.Max < item.Min {
			item.Max = item.Min
		}

		count += item.Max

		if count > roundTabPrec {
			item.Max -= count - roundTabPrec
			count = roundTabPrec
		}
	}

	for i := len(*tab) - 1; i >= 0; i-- {
		item := &(*tab)[i]

		if item.Max < item.Min {
			delta := item.Min - item.Max

			if i > 0 {
				(*tab)[i-1].Max -= delta
			}

			item.Max = item.Min
		}
	}

	for i := 1; i < len(*tab); i++ {
		(*tab)[i].Max += (*tab)[i-1].Max
	}

	n := rand.Int31n(roundTabPrec)

	for i := 0; i < len(*tab); i++ {
		if n < (*tab)[i].Max {
			return DamageJudgeRv(i), n
		}
	}

	return DamageJudgeRv_Count - 1, n
}

// getRoundRvText 获取圆桌运算结果文本
func getRoundRvText(roundRv DamageJudgeRv) string {
	switch roundRv {
	case DamageJudgeRv_Miss:
		return "丢失"
	case DamageJudgeRv_Dodge:
		return "闪避"
	case DamageJudgeRv_Crit:
		return "暴击"
	case DamageJudgeRv_Block:
		return "格挡"
	case DamageJudgeRv_Hit:
		return "正常"
	}
	return ""
}

// getDamageValueKindText 获取伤害数值类型文本
func getDamageValueKindText(damageValueKind DamageValueKind) string {
	switch damageValueKind {
	case DamageValueKind_Normal:
		return "物理"
	case DamageValueKind_Fire:
		return "火焰"
	case DamageValueKind_Cold:
		return "冰霜"
	case DamageValueKind_Poison:
		return "毒素"
	case DamageValueKind_Lightning:
		return "闪电"
	}
	return ""
}
