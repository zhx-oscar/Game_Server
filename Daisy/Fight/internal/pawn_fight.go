package internal

// _PawnFight pawn战斗模块
type _PawnFight struct {
	pawn        *Pawn
	Attr        FightAttr   // 战斗属性
	State       FightState  // 战斗状态
	BeHit       FightBeHit  // 战斗受击
	Events      FightEvents // 事件系统
	Master      *Pawn       // 主人
	SlaveList   []*Pawn     // 从属npc列表
	BornTime    uint32      // 出生时间
	DestroyTime uint32      // 销毁时间
}

// init 初始化
func (pawnFight *_PawnFight) init(pawn *Pawn) {
	pawnFight.pawn = pawn

	// 初始化战斗属性
	pawn.Attr.init(pawn)

	// 初始化战斗状态
	pawn.State.init(pawn)

	// 初始化战斗受击
	pawn.BeHit.init(pawn)

	// 初始化事件源
	pawn.Events.init(pawn)
}

// RefreshLifeTime 刷新Pawn生命期
func (pawnFight *_PawnFight) RefreshLifeTime() bool {
	if !pawnFight.pawn.IsAlive() || pawnFight.pawn.Info.LifeTime <= 0 {
		return false
	}

	pawnFight.pawn.DestroyTime = pawnFight.pawn.Scene.NowTime + pawnFight.pawn.Info.LifeTime

	return true
}

// ExtendLifeTime 延长Pawn生命期
func (pawnFight *_PawnFight) ExtendLifeTime(duration uint32) bool {
	if !pawnFight.pawn.IsAlive() || pawnFight.pawn.Info.LifeTime <= 0 {
		return false
	}

	pawnFight.pawn.DestroyTime += duration

	return true
}
