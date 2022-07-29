package Fight

import (
	"Daisy/Fight/internal"
	"Daisy/Fight/internal/conf"
	_ "Daisy/Fight/internal/effects/buffeffect"
	_ "Daisy/Fight/internal/effects/skilleffect"
	"Daisy/Fight/internal/log"
	"Daisy/Proto"
)

// PawnInfo pawn信息
type PawnInfo = internal.PawnInfo

// SimulatorModeInfo 模拟器内pawninfo
type SimulatorModeInfo = internal.SimulatorModeInfo

// SceneInfo 战斗场景信息
type SceneInfo = internal.SceneInfo

// FormationInfo 战斗阵形信息
type FormationInfo = internal.FormationInfo
type BornPoint = internal.BornPoint

// BattleMaxMilliseconds 场景战斗最大时间 沟通策划之后 暂定5分钟 单位:毫秒
const BattleMaxMilliseconds = 5 * 60 * 1000

// Play 运行战斗
func Play(sceneInfo *SceneInfo) (*Proto.FightResult, error) {
	defer log.Flush()

	// 开启测试模式
	if !sceneInfo.TestMode {
		sceneInfo.TestMode = log.GetLogLevel() != log.Level_Off
	}

	// 质检要求测试模式下开启模拟器模式
	sceneInfo.SimulatorMode = sceneInfo.TestMode

	// 创建场景
	scene, err := internal.CreateScene(sceneInfo)
	if err != nil {
		return nil, err
	}

	// 运行战斗
	return scene.Run(), nil
}

// LoadConfig 加载战斗模块配置 提供给战斗模拟器使用
func LoadConfig(path string) {
	internal.LoadBehaviorTreeFile(path + "res/ai/chaos.b3")
	conf.SetSkillTimeLinePath(path + "res/Timeline/Skill")
	conf.SetBeHitPath(path + "res/AnimatorController")
	conf.SetActTimelinePath(path + "res/Timeline")
	conf.SetAttackTimelinePath(path + "res/Timeline/Attack")
	conf.HotUpdateConfs(true, true)
}

// LogLevel 日志等级
const (
	Level_Off     = log.Level_Off
	Level_Console = log.Level_Console
	Level_File    = log.Level_File
)

// SetLogLevel 设置log等级
func SetLogLevel(level int32) {
	log.SetLogLevel(level)
}

// SetLogPath 设置log路径
func SetLogPath(logDir string) {
	log.SetLogPath(logDir)
}

// CalcAttr 属性计算器
type CalcAttr = internal.CalcAttr

// SkillKind 技能类型
type SkillKind = conf.SkillKind

const (
	SkillKind_Super     = conf.SkillKind_Super     // 1：超能技
	SkillKind_NormalAtk = conf.SkillKind_NormalAtk // 2：普攻
	SkillKind_Ultimate  = conf.SkillKind_Ultimate  // 3: 必杀技
	SkillKind_Combine   = conf.SkillKind_Combine   // 4: 合体必杀技
)
