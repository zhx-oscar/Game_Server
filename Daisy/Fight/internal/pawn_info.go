package internal

import (
	. "Daisy/Fight/attraffix"
	"Daisy/Fight/internal/conf"
	"Daisy/Proto"
	"errors"
	"fmt"
)

// PawnInfo pawn信息
type PawnInfo struct {
	*Proto.PawnInfo                              // pawn信息
	AttrAffixList             []AttrAffix        // 属性词缀列表
	NormalAtkList             []uint32           // 普攻列表
	SuperSkillList            []uint32           // 超能技列表
	UltimateSkillList         []uint32           // 必杀技列表
	OverDriveNormalAttackList []uint32           // 超载普攻池
	OverDriveSuperSkillList   []uint32           // 超载超能池
	CombineSkillList          []uint32           // 合体技表
	AddComboAttack            uint32             // 普攻追加连段
	BornBuffs                 []uint32           // 出生buff列表
	LifeTime                  uint32             // 生存时间（单位ms）
	SimulatorModeInfo         *SimulatorModeInfo //模拟器下专用数据 nil不生效 内部数值 非0生效
}

// SimulatorModeInfo 模拟器内pawninfo
type SimulatorModeInfo struct {
	MaxHP           int64
	Attack          float64
	ExtendDodgeRate float32
	ExtendBlockRate float32
	ExtendHitRate   float32
	ExtendCritRate  float32
}

// _PawnInfo pawn信息
type _PawnInfo struct {
	*PawnInfo
	*conf.PawnConfig // pawn配置
}

// init 初始化
func (pawnInfo *_PawnInfo) init(pawn *Pawn, info *PawnInfo) error {
	if info == nil || info.PawnInfo == nil {
		return errors.New("pawn info is nil")
	}

	info.Id = pawn.UID
	pawn.Info.PawnInfo = info

	pawnConfig, ok := pawn.Scene.GetPawnConfig(info.Type, info.ConfigId)
	if !ok {
		return fmt.Errorf("not found pawn type %v configid %d config", info.Type, info.ConfigId)
	}
	pawnInfo.PawnConfig = pawnConfig

	return nil
}
