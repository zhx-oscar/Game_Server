package internal

import (
	"Cinder/Base/linemath"
	"Daisy/Fight/internal/conf"
	"Daisy/Proto"
	"errors"
)

// CreateScene 创建场景
func CreateScene(sceneInfo *SceneInfo) (*Scene, error) {
	// 创建场景
	scene := &Scene{}
	if err := scene.init(sceneInfo); err != nil {
		return nil, err
	}

	return scene, nil
}

// SceneStage 场景阶段
type SceneStage uint32

const (
	SceneStage_InitFight SceneStage = iota // 0: 场景初始化战斗阶段
	SceneStage_WaitFight                   // 1: 场景等待开战阶段 （包含怪物出生CG）
	SceneStage_Fight                       // 2: 场景战斗阶段
	SceneStage_WaitEnd                     // 3: 场景等待结束阶段 （包含技能后摇，子弹命中，怪物死亡CG）
	SceneStage_End                         // 4: 场景结束阶段
)

// SceneInfo 战斗场景信息
type SceneInfo struct {
	Formation       []*FormationInfo    // 对战双方阵型信息
	MaxMilliseconds uint32              // 最大战斗时间（毫秒）
	SimulatorMode   bool                // 模拟器模式
	TestMode        bool                // 测试模式
	BoundaryPoints  []linemath.Vector2  // 边界点
	Inherit         *Proto.FightInherit // 继承战斗数据
}

// Scene 战斗场景
type Scene struct {
	NowTime           uint32          // 当前时间
	nowFrames         uint32          // 当前帧数
	maxFrames         uint32          // 最大帧数
	secFrames         uint32          // 每秒帧数
	uidSrc            uint32          // 唯一ID生成器
	eventDeep         uint32          // 事件递归深度
	Info              *SceneInfo      // 场景信息
	pawnList          []*Pawn         // pawn列表
	formationList     []*Formation    // 双方战斗阵型
	fightBeginTime    uint32          // 战斗开始时间
	fightEndTime      uint32          // 战斗结束结束时间
	winCamp           Proto.Camp_Enum // 胜利阵营
	stage             SceneStage      // 场景阶段
	_SkillFlow                        // 技能流程
	_BuffFlow                         // buff流程
	_CombineSkillFlow                 // 合体必杀技流程
	_AttackFlow                       // 伤害体流程
	_BeHitFlow                        // 受击流程
	_BehaviorFlow                     // 行为树流程
	_MovementFlow                     // 移动控制器流程
	_HaloFlow                         // 光环控制器流程
	_ReplayMaker                      // 回放生成器
	_Terrain                          // 2d地形
	_RegionMgr                        // 区域管理器
	*conf.Configs                     // 战斗配置
}

// init 初始化场景
func (scene *Scene) init(sceneInfo *SceneInfo) error {
	if sceneInfo == nil {
		return errors.New("args invalid")
	}

	scene.Info = sceneInfo

	// 获取配置
	scene.Configs = conf.GetConfigs()

	// 创建地形
	scene._Terrain.init(scene)

	// 创建区域管理器
	scene._RegionMgr.init(scene)

	// 初始化帧与时间
	scene.secFrames = conf.FrameRate
	scene.NowTime = 0
	scene.nowFrames = 0
	scene.maxFrames = sceneInfo.MaxMilliseconds / 1000 * scene.secFrames
	scene.fightEndTime = scene.framesToTime(scene.maxFrames)
	scene.winCamp = Proto.Camp_Blue

	// 初始化回放生成器
	scene._ReplayMaker.init(scene)

	if len(sceneInfo.Formation) <= 1 {
		return errors.New("formation not enough")
	}

	// 初始化对战阵型
	scene.formationList = make([]*Formation, len(sceneInfo.Formation))

	for i := range scene.formationList {
		scene.formationList[i] = &Formation{}
		if err := scene.formationList[i].init(scene, sceneInfo.Formation[i], Proto.Camp_Enum(i)); err != nil {
			return err
		}
	}

	// 所有阵型上阵pawn
	for _, formation := range scene.formationList {
		if err := formation.putPawns(); err != nil {
			return err
		}

		if err := formation.putBackGroundPawns(); err != nil {
			return err
		}
	}

	// 初始化skill流程
	scene._SkillFlow.init(scene)

	// 初始化buff流程
	scene._BuffFlow.init(scene)

	// 初始化合体技流程
	scene._CombineSkillFlow.init(scene)

	// 初始化伤害体流程
	scene._AttackFlow.init(scene)

	// 初始化受击流程
	scene._BeHitFlow.init(scene)

	// 初始化行为树流程
	scene._BehaviorFlow.init(scene)

	// 初始化移动控制器流程
	scene._MovementFlow.init(scene)

	// 初始化光环控制器流程
	scene._HaloFlow.init(scene)

	return nil
}
