package Const

const (
	AffixValueType_A = 1
	AffixValueType_B = 0
)

const (
	AffixPrecision_0 = 0
	AffixPrecision_2 = 1
	AffixPrecision_4 = 2
)

const (
	AffixEffectType_Attr = 0
	AffixEffectType_Buff = 1
)

const (
	Quality_White    = 1 //白
	Quality_Gray     = 2 //灰
	Quality_Blue     = 3 //蓝
	Quality_Yellow   = 4 //黄
	Quality_Dullgold = 5 //暗金
	Quality_Green    = 6 //绿色
)

const (
	AffixPlace_Front = 1 //前缀
	AffixPlace_Back  = 2 //后缀
)

const (
	AffixScoreType_Extra         = iota // 其他
	AffixScoreType_Attack               // 攻击
	AffixScoreType_Defence              // 防御
	AffixScoreType_Psychokinesis        // 超自然
)
