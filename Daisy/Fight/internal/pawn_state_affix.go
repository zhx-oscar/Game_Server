package internal

// Stat 状态定义
type Stat uint32

// 一级状态
const (
	_Stat_First_Begin Stat = 1

	Stat_CantUseNormalAtk Stat = iota
	Stat_CantUseSuperSkill
	Stat_CantUseUltimateSkill
	Stat_CantMove
	Stat_CantBeDamage
	Stat_CantBeEnemySelect
	Stat_CantBeHitControl
	Stat_CantBeAddIncrBuff
	Stat_CantBeAddDecrBuff
	Stat_CantBeFriendlySelect

	_Stat_First_End
)

// 二级状态
const (
	_Stat_Second_Begin Stat = 21

	Stat_Invincible Stat = iota + 20
	Stat_Weak
	Stat_OverDrive
	Stat_OverDriveStartCG
	Stat_Balance
	Stat_Dodging
	Stat_Death
	Stat_Raged
	Stat_EnergyShieldOn
	Stat_Frozen

	_Stat_Second_End
)
