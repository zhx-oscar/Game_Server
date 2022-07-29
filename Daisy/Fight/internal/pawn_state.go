package internal

// FightState 战斗状态
type FightState struct {
	pawn *Pawn

	// 一级状态
	CantUseNormalAtk     bool // 不能使用普攻
	CantUseSuperSkill    bool // 不能使用超能技
	CantUseUltimateSkill bool // 不能使用必杀技
	CantMove             bool // 不能移动
	CantBeDamage         bool // 免疫伤害
	CantBeEnemySelect    bool // 不能被敌方锁定
	CantBeFriendlySelect bool // 不能被友方锁定
	CantBeHitControl     bool // 不能被受击状态控制
	CantBeAddIncrBuff    bool // 不能被施加增益buff
	CantBeAddDecrBuff    bool // 不能被施加减益buff

	// 二级状态
	Invincible       bool // 无敌
	Weak             bool // 虚弱
	OverDrive        bool // 超载
	OverDriveStartCG bool // 播放超载开始CG
	Balance          bool // 平衡状态
	Dodging          bool // 正在闪避
	Death            bool // 是否死亡
	Raged            bool // 是否狂暴
	EnergyShieldOn   bool // 能量护盾生效
	Frozen           bool // 冰冻

	// changedTab 状态变化中间值表
	changedTab _StatChangedTab

	// 受击状态
	BeHitStat Bits // 受击状态
}

//init 战斗状态初始化
func (state *FightState) init(pawn *Pawn) {
	state.pawn = pawn
	state.changedTab = make(_StatChangedTab, _Stat_Second_End)
}

// copy 拷贝
func (state FightState) copy(pawn *Pawn) *FightState {
	state.pawn = pawn
	state.changedTab = nil
	return &state
}
