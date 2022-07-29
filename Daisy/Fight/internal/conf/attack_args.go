package conf

import (
	"Cinder/Base/linemath"
	"Daisy/Proto"
)

// AttackSpawn 伤害体出生配置
type AttackSpawn struct {
	Pos            AttackSpawnPos   // 出生位置
	Rotate         float32          // 转向角度
	Offset         linemath.Vector2 // 偏移位置
	AutoExtend     bool             // 自动调整形状大小（矩形调整长度，圆形和扇形调整半径）
	AutoRectOffset bool             // 自动调整长方形偏移位置
}

// AttackShape 伤害体形状
type AttackShape struct {
	Type     Proto.AttackShapeType_Enum // 伤害体形状类型
	Extend   linemath.Vector2           // 矩形区域全长全高
	Radius   float32                    // 圆形或扇形区域半径
	FanAngle float32                    // 扇形区域夹角
}

// AttackArgs 伤害体参数
type AttackArgs struct {
	Type           AttackType           // 类型
	MoveMode       AttackMoveMode       // 移动模式
	HitMode        AttackHitMode        // hit模式
	DestroyType    AttackDestroyType    // 销毁方式
	Spawn          AttackSpawn          // 出生配置
	Shape          AttackShape          // 伤害体形状
	TargetCategory AttackTargetCategory // 目标选取策略
	MaxHitTarget   int32                // 每次hit最大目标数量
	MaxLinkTarget  int32                // 最大连接目标数量
	LinkDistance   float32              // 最大连接距离
	RepeatLink     bool                 // 能否重复连接
	Grouping       bool                 // 是否分组
	CanBreak       bool                 // 是否随技能打断
	CastRange      bool                 // 作为技能施法范围判定
	NoAuto         bool                 // 不自动开始
	OnFinish       []int                // 结束后创建伤害体
	OnMoveEnd      []int                // 移动结束时创建伤害体
}

// IsLink 是否是链式伤害体
func (attackArgs *AttackArgs) IsLink() bool {
	return attackArgs.MaxLinkTarget > 0 && attackArgs.LinkDistance > 0
}
