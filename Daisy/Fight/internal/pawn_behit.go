package internal

import "Daisy/Proto"

// FightBeHit 战斗受击
type FightBeHit struct {
	pawn *Pawn

	BeHitStartTime         []uint32 // 受击开始时间
	BeHitEndTime           []uint32 // 受击结束时间
	DodgeEndTime           uint32   // 闪避结束时间
	CurBlockBreakValue     int32    // 当前格挡失衡值
	BlockBreakValueDecTime uint32   // 格挡失衡值下降时间点
	ReFloatStartHeight     float64  // 再次被击飞起始高度
}

//init 战斗状态初始化
func (beHit *FightBeHit) init(pawn *Pawn) {
	beHit.pawn = pawn
	beHit.BeHitStartTime = make([]uint32, len(Proto.HitType_Enum_name))
	beHit.BeHitEndTime = make([]uint32, len(Proto.HitType_Enum_name))
}
