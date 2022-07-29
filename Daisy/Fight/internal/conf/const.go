package conf

// AttackType 伤害体类型
type AttackType uint32

const (
	AttackType_Single AttackType = iota // 0: 单体
	AttackType_Aoe                      // 1: 区域
)

// AttackHitMode 伤害体命中模式
type AttackHitMode uint32

const (
	AttackHitMode_None        AttackHitMode = iota // 0: 没有hit
	AttackHitMode_TimeLine                         // 1: 时间轴
	AttackHitMode_FixInterval                      // 2: 固定时间间隔
)

// AttackDestroyType 伤害体消失类型
type AttackDestroyType uint32

const (
	AttackDestroyType_None           AttackDestroyType = iota // 0: 不销毁
	AttackDestroyType_LifeTime                                // 1: 生存时间
	AttackDestroyType_FixTimeMoveEnd                          // 2: 固定时间移动结束时销毁
	AttackDestroyType_HitTimes                                // 3: hit次数
)

// AttackSpawnPos 伤害体出生位置
type AttackSpawnPos uint32

const (
	AttackSpawnPos_Caster  AttackSpawnPos = iota // 0: 施法者位置
	AttackSpawnPos_Target                        // 1: 目标位置
	AttackSpawnPos_Inherit                       // 2: 从上个伤害体继承
)

// AttackTargetCategory 伤害体目标选取策略
type AttackTargetCategory uint32

const (
	AttackTargetCategory_Enemy  AttackTargetCategory = iota // 0: 敌人
	AttackTargetCategory_Friend                             // 1: 友方
)

// AttackMoveMode 伤害体位移模式
type AttackMoveMode uint32

const (
	AttackMoveMode_None                AttackMoveMode = iota // 0: 不移动
	AttackMoveMode_FllowTarget                               // 1: 跟随目标
	AttackMoveMode_FllowCaster                               // 2: 跟随施法者
	AttackMoveMode_FixTimeMoveToTarget                       // 3: 固定时间移向目标
)

// BuffKind buff增益类型
type BuffKind = uint32

const (
	BuffKind_Incr BuffKind = iota // 0：增益
	BuffKind_Decr                 // 1：减益
)

// BuffDisappearType buff销毁类型
type BuffDisappearType = uint32

const (
	BuffDisappearType_Duration   BuffDisappearType = iota // 0：时间到消失
	BuffDisappearType_Dead                                // 1：死亡后消失
	BuffDisappearType_Attack                              // 2: 发动攻击后消失
	BuffDisappearType_BeDamage                            // 3: 受到伤害后消失
	BuffDisappearType_OutOfRange                          // 4: 出区域后消失
	BuffDisappearType_Permanent                           // 5: 永久
	BuffDisappearType_Other                               // 6: 其他
)

// ShieldValueSrc 护盾数值来源定义
type ShieldValueSrc uint32

const (
	ShieldValueSrc_CasterMaxHP ShieldValueSrc = iota // 0：最大血量（施法者）
	ShieldValueSrc_ExtValue                          // 1：外部传入值
)

// BuffOverlapType buff叠加类型
type BuffOverlapType uint32

const (
	BuffOverlapType_Self   BuffOverlapType = iota // 0: 自己释放的可叠加
	BuffOverlapType_Friend                        // 1: 队友释放的可叠加
)

// DamageStep 伤害步骤
type DamageStep uint32

const (
	DamageStep_AttackLuckyJudge          DamageStep = iota // 攻击幸运判定
	DamageStep_CountAttackDamageValue                      // 统计攻击伤害
	DamageStep_IncrAttackDamageValue                       // 攻击增伤修正
	DamageStep_CountCastSkillDamageValue                   // 统计施法伤害
	DamageStep_IncrCastSkillDamageValue                    // 施法增伤修正
	DamageStep_RoundTableJudge                             // 圆桌判定
	DamageStep_WaterfallJudge                              // 瀑布判定
	DamageStep_DefendDeductHPShield                        // 防御扣除护盾
	DamageStep_ResistanceDamageValue                       // 抗性修正
	DamageStep_DecrDamageValue                             // 减伤修正
	DamageStep_Sputtering                                  // 分裂伤害
	DamageStep_DeductHP                                    // 伤害扣血
	DamageStep_Bloodsucker                                 // 吸血
)

// DamageFlow 伤害流程
type DamageFlow uint32

const (
	DamageFlow_Attack    DamageFlow = iota // 0：攻击流程
	DamageFlow_CastSkill                   // 1：施法流程
)

// DamageKind 伤害类型定义
type DamageKind uint32

const (
	DamageKind_Hurt         DamageKind = iota // 0：伤害
	DamageKind_Heal                           // 1：治疗
	DamageKind_DOT                            // 2：DOT
	DamageKind_HOT                            // 3：HOT
	DamageKind_Bloodsucking                   // 4: 吸血
	DamageKind_Thorns                         // 5: 反伤
	DamageKind_Sputtering                     // 6: 溅射
)

// DamageValueKind 伤害数值类型
type DamageValueKind uint32

const (
	DamageValueKind_Normal    DamageValueKind = iota // 0：物理
	DamageValueKind_Fire                             // 1: 火
	DamageValueKind_Cold                             // 2: 冰
	DamageValueKind_Poison                           // 3: 毒
	DamageValueKind_Lightning                        // 4: 电
	DamageValueKind_End
	DamageValueKind_Begin = DamageValueKind_Normal
)

// DamageJudgeRv 伤害判定结果
type DamageJudgeRv uint32

const (
	DamageJudgeRv_Miss  DamageJudgeRv = iota // MISS
	DamageJudgeRv_Dodge                      // 闪避
	DamageJudgeRv_Crit                       // 暴击
	DamageJudgeRv_Block                      // 格挡
	DamageJudgeRv_Hit                        // 命中
	DamageJudgeRv_Count
)

// InnerBuffID 内部buff id
type InnerBuffID = uint32

const (
	InnerBuffID_Begin                     InnerBuffID = 1
	InnerBuffID_OverDrive                 InnerBuffID = iota // 1: 程序-超载
	InnerBuffID_EnergyShield                                 // 2: 程序-能量护盾
	InnerBuffID_Unbalance                                    // 3: 程序-受击失衡
	InnerBuffID_BLockBreak                                   // 4: 程序-格挡失衡
	InnerBuffID_BornAct                                      // 5: 程序-出生表演
	InnerBuffID_RecoverHP                                    // 6: 程序-每秒恢复HP
	InnerBuffID_RecoverUltimateSkillPower                    // 7: 程序-每秒恢复必杀技能量
	InnerBuffID_Thorns                                       // 8: 程序-反弹伤害
	InnerBuffID_StealUltimateSkillPower                      // 9: 程序-偷取必杀技能量
	InnerBuffID_End
)

// AttackTemplID 伤害体模板ID
type AttackTemplID = uint32

const (
	AttackTemplID_Begin        AttackTemplID = 1
	AttackTemplID_SignleRemote AttackTemplID = iota // 1: 单体远程攻击
	AttackTemplID_Aoe                               // 2: 不可移动AOE
	AttackTemplID_Laser                             // 3: 穿透型激光
	AttackTemplID_GravityGun                        // 4: 重力炮
	AttackTemplID_SignleMelee                       // 5: 单体近战攻击
	AttackTemplID_MoveableAoe                       // 6: 可移动AOE
	AttackTemplID_LightChain                        // 7: 闪电链
	AttackTemplID_BounceBall                        // 8: 弹弹球
	AttackTemplID_P2PLaser                          // 9: 点对点激光
	AttackTemplID_Sputtering                        // 10: 分裂伤害
	AttackTemplID_End
)

// PropType 属性类型
type PropType = uint32

const (
	PropType_None          PropType = iota // 0: 无
	PropType_Strength                      // 1: 力量型
	PropType_Agility                       // 2: 敏捷型
	PropType_Psychokinesis                 // 3: 念力型
)

// SkillKind 技能类型
type SkillKind = uint32

const (
	SkillKind_Super     SkillKind = iota + 1 // 1：超能技
	SkillKind_NormalAtk                      // 2：普攻
	SkillKind_Ultimate                       // 3: 必杀技
	SkillKind_Combine                        // 4: 合体必杀技
)

// 战斗常量表配置项定义
const (
	ConstExcel_WeakDebuff               = 1 // 进入虚弱状态获得的debuff
	ConstExcel_SputteringAttackTmplID   = 2 // 分裂伤害体模板ID
	ConstExcel_SputteringAttackTimeLine = 3 // 分裂伤害体时间轴
	ConstExcel_GravityAcceleration      = 4 // 重力加速度
	ConstExcel_HitFloatRiseTime         = 5 // 击飞上升时间(ms)
	ConstExceL_HitReFloatRiseTime       = 6 // 重复被击飞上升时间(ms)
	ConstExcel_HitReFloatDeltaTime      = 7 // 重复被击飞保护间隔时间(ms)
)

// InnerAttackID 内置伤害体定义
type InnerAttackID uint32

const (
	InnerAttackID_SputteringAttack InnerAttackID = iota // 0: 分裂伤害体
)

// FrameRate 帧率
const FrameRate uint32 = 30

// DodgeAngleRange 闪避角度范围
const DodgeAngleRange uint32 = 150

// ScenePawnCountMax 战斗对象数量上限
const ScenePawnCountMax = 100 //沟通策划之后，暂定战斗场景中最多100个战斗对象

type SummonType uint32

const (
	SummonTypeDelay           SummonType = iota //延迟固定时间召唤
	SummonTypeDead                              //宿主死亡时召唤
	SummonTypeOnceDelayOrDead                   //宿主死亡或者延迟固定时间到了的时候召唤 召唤仅仅生效一次
)
