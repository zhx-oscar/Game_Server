package internal

import (
	"errors"
)

// Pawn 战斗对象
type Pawn struct {
	// 基础信息
	UID        uint32    // pawn id
	Info       _PawnInfo // pawn配置
	Scene      *Scene    // 当前所在场景
	isSnapshot bool      // 是否是数据快照

	// 战斗模块
	_PawnFight // 战斗模块
	_PawnSkill // 技能栏模块
	_PawnBuff  // buff模块

	// 位移模块
	_PawnMovement    // 移动控制器
	_PawnFreeStation // 空闲站位

	// AI模块
	_PawnBehavior // AI行为树
}

// init Pawn初始化
func (pawn *Pawn) init(scene *Scene, pawnInfo *PawnInfo) error {
	if scene == nil || pawnInfo == nil {
		return errors.New("args invalid")
	}

	pawn.Scene = scene
	pawn.UID = scene.generateUID()

	// 初始化pawn信息
	if err := pawn.Info.init(pawn, pawnInfo); err != nil {
		return err
	}

	// 初始化战斗模块
	pawn._PawnFight.init(pawn)

	// 初始化技能栏模块
	if err := pawn._PawnSkill.init(pawn); err != nil {
		return err
	}

	// 初始化buff模块
	pawn._PawnBuff.init(pawn)

	//初始化移动木块
	pawn._PawnMovement.init(pawn)

	//初始化就位点模块
	pawn._PawnFreeStation.init(pawn)

	//初始化行为树
	pawn._PawnBehavior.init(pawn)

	return nil
}
