package internal

import (
	"Daisy/Fight/internal/conf"
	"Daisy/Fight/internal/log"
	"Daisy/Proto"
	"fmt"
	b3core "github.com/magicsea/behavior3go/core"
	"math/rand"
)

//_PawnBehavior 行为树
type _PawnBehavior struct {
	owner            *Pawn
	board            *b3core.Blackboard //记录行为状态
	bevTree          *b3core.BehaviorTree
	aiData           *conf.AIInfo
	blackBoardKeys   *conf.BlackBoardKeys //行为树黑板
	pauseCount       uint32               //暂停计数
	dashEnd          bool
	lastAttackTarget *Pawn //上一次目标

	simulatorActionTorsionList []string //战斗模拟器行为树节点扭转列表
}

func (ai *_PawnBehavior) init(owner *Pawn) {
	ai.owner = owner
	ai.blackBoardKeys = &conf.BlackBoardKeys{Data: map[string]interface{}{}}

	ai.initBehavior(owner.Info.AIID)
}

func (ai *_PawnBehavior) initBehavior(aiid uint32) {
	if aiid == 0 {
		return
	}

	aiConf, ok := ai.owner.Scene.GetAIDataConf(aiid)
	if !ok {
		return
	}

	ai.aiData = aiConf
	ai.resetBlackboardValue()
	tree := getBevTree(aiConf.TreeName)
	if tree == nil {
		panic(fmt.Sprintln("initBehavior tree is nil. TreeName is:", aiConf.TreeName))
		return
	}

	ai.board = b3core.NewBlackboard()
	ai.bevTree = tree
}

//getBlackboardValueByKey 获取外部黑板内的值
func (ai *_PawnBehavior) getBlackboardValueByKey(key string) interface{} {
	val, ok := ai.blackBoardKeys.Data[key]
	if ok {
		return val
	}

	return nil
}

//getAttrValueByKey 获取对应属性数值，通过黑板定义的key
func (ai *_PawnBehavior) getAttrValueByKey(key string) interface{} {
	switch key {
	case eB_Attr_HP_Per:
		return float32(ai.owner.Attr.CurHP) / float32(ai.owner.Attr.MaxHP)
	default:
		return 0
	}
}

//resetBlackboardValue 重置行为树对应的外部黑板数据
func (ai *_PawnBehavior) resetBlackboardValue() {
	defaultBlackboard := ai.aiData.BlackBoardKeys

	//默认黑板处理
	for key, val := range defaultBlackboard.Data {
		ai.blackBoardKeys.Data[key] = val
	}

	//外部黑板重写
	for key, val := range ai.owner.Info.BlackBoardKeyData {
		ai.blackBoardKeys.Data[key] = val
	}
}

// behaviorUpdate 行为树循环
func (ai *_PawnBehavior) behaviorUpdate() {
	if ai.bevTree == nil {
		return
	}

	//AI 是否被暂停
	if ai.isAIPause() {
		return
	}

	ai.bevTree.Tick(ai.owner, ai.board)
}

//isAIPause AI是否暂停
func (ai *_PawnBehavior) isAIPause() bool {
	//AI 是否被外部主动暂停
	if ai.pauseCount > 0 {
		return true
	}

	//不能移动 + 不能使用技能(普攻+超能技+必杀技)  == AI暂停
	if ai.owner.State.CantMove && ai.owner.State.CantUseNormalAtk && ai.owner.State.CantUseSuperSkill && ai.owner.State.CantUseUltimateSkill {
		return true
	}

	return false
}

// getBoard 获取黑板
func (ai *_PawnBehavior) getBoard() *b3core.Blackboard {
	return ai.board
}

// GetAIConf 获取AI配置
func (ai *_PawnBehavior) GetAIConf() *conf.AIInfo {
	return ai.aiData
}

//ResetBehavior 重置行为树
func (ai *_PawnBehavior) ResetBehavior(aiid uint32) {
	ai.initBehavior(aiid)
}

//blackboardSetValue 黑板设置值
func (ai *_PawnBehavior) blackboardSetValue(key string, value interface{}) {
	ai.board.Set(key, value, "", "")
}

//AIPause AI暂停
func (ai *_PawnBehavior) AIPause(pause bool) {
	if pause {
		ai.pauseCount++

		//当前是否处于AI控制的移动中，需要主动停止
		if ai.owner.IsMoving() && !ai.owner.IsPassive() {
			ai.owner.Stop()
		}

		ai.owner.BreakCurSkill(ai.owner, Proto.SkillBreakReason_AIPause)

		return
	}

	if ai.pauseCount > 0 {
		ai.pauseCount--
	}
}

//AIBackToRoot 回到根节点
func (ai *_PawnBehavior) AIBackToRoot() {
	tree := getBevTree(ai.aiData.TreeName)
	if tree == nil {
		panic(fmt.Sprintln("AIBackToRoot tree is nil. TreeName is:", ai.aiData.TreeName))
	}

	ai.board = b3core.NewBlackboard()
	ai.bevTree = tree
}

//startDashing 冲刺
func (ai *_PawnBehavior) startDashing() bool {
	pawn := ai.owner
	target, ok := getPawnBoardByBlackboard(ai.board, attackTarget)
	if !ok {
		log.Error("*************[startDashing] getPawnBoardByBlackboard fail")
		return false
	}

	skill, ok := getSkillBoardByBlackboard(ai.board, curSkill)
	if !ok {
		log.Error("*************[startDashing] getSkillBoardByBlackboard fail")
		return false
	}

	//当前技能是否可以冲刺
	if !skill.Config.CanDash {
		return false
	}

	// 查询冲刺目标点
	station, moveToPos, ok := target.getFreeNearestStation(pawn, rand.Float32()*(skill.GetMaxCastDistance()-skill.GetMinCastDistance()), skill.GetMinCastDistance(), skill.GetMaxCastDistance(), false)
	if !ok {
		return false
	}

	//当冲刺点 和自己位置一致，不需要冲刺
	if Vector2Equal(ai.owner.GetPos(), moveToPos) {
		ai.dashEnd = true
		return true
	} else {
		angle := CalcAngle(pawn.GetPos(), moveToPos)
		target.setStationData(station)

		// 向目标移动
		if !pawn.MoveToAndChangeAngle(Proto.MoveMode_Fast, moveToPos, pawn.GetAngleVelocity(float32(angle), pawn.Info.FastSpeed), float32(angle), true) {
			return false
		}

		ai.owner.AddMoveEndCallBack(ai.dashingEndCallBack)
		ai.dashEnd = false
		ai.AIPause(true)
	}

	return true
}

//dashingEndCallBack 冲刺结束回调
func (ai *_PawnBehavior) dashingEndCallBack(pawn *Pawn) {
	pawn.stopDashing()
}

//stopDashing 停止冲刺
func (ai *_PawnBehavior) stopDashing() {
	if !ai.dashEnd {
		ai.dashEnd = true
		ai.AIPause(false)

		//如果当前正在冲刺的话，停止冲刺
		if ai.owner.moveMode == Proto.MoveMode_Fast && ai.owner.IsMoving() {
			ai.owner.Stop()
		}
	}
}

//isDashingEnd 是否冲刺结束
func (ai *_PawnBehavior) isDashingEnd() bool {
	return ai.dashEnd
}

//setAttackTarget 设置AI 攻击目标
func (ai *_PawnBehavior) setAttackTarget(tick *b3core.Tick, target *Pawn) {
	if tick == nil || target == nil {
		return
	}

	tick.Blackboard.Set(attackTarget, target, "", "")
	ai.syncSetTargetAction(target)
}

//syncSetTargetAction 同步SetTargetAction
func (ai *_PawnBehavior) syncSetTargetAction(target *Pawn) {
	if target == nil {
		return
	}

	//目标是自己时
	if ai.owner.UID == target.UID {
		return
	}

	//非角色时
	if !ai.owner.IsRole() {
		return
	}

	//目标和自己时同阵营时
	if ai.owner.GetCamp() == target.GetCamp() {
		return
	}

	//目标一样过滤重复发送
	if ai.lastAttackTarget == target {
		return
	}

	//fmt.Println("+++++++++++++ syncSetTargetAction  ", ai.owner.UID, target.UID, ai.owner.Scene.nowFrames, ai.owner.State.BeHitStat)
	ai.owner.Scene.PushAction(&Proto.SetTarget{
		SelfId:   ai.owner.UID,
		TargetId: target.UID,
	})
	ai.lastAttackTarget = target
}

//cantBeSelect 是否可以被选择
func (ai *_PawnBehavior) cantBeSelect(pawn *Pawn) bool {
	if pawn == nil {
		return true
	}

	//不能被敌方锁定
	if ai.owner.GetCamp() != pawn.GetCamp() && pawn.State.CantBeEnemySelect {
		return true
	}

	//不能被友方锁定
	if ai.owner.GetCamp() == pawn.GetCamp() && pawn.State.CantBeFriendlySelect {
		return true
	}

	return false
}

//pushActionTorsionInfo push 节点名字  只有模拟器模式下才会生效
func (ai *_PawnBehavior) pushActionTorsionInfo(actionName string) {
	if ai.owner == nil {
		return
	}

	if !ai.owner.Scene.SimulatorMode() {
		return
	}

	ai.simulatorActionTorsionList = append(ai.simulatorActionTorsionList, actionName)

	ai.owner.Scene.PushProgrammerDebugInfo(func() string {
		result := fmt.Sprintf("behaviorActionTorsion ${PawnID:%d} => %v", ai.owner.UID, actionName)
		return result
	})
}

//InCampAlertAOERange 是否处于某方阵营的AOE警戒范围内
func (ai *_PawnBehavior) InCampAlertAOERange(camp Proto.Camp_Enum) bool {
	result := ai.owner.Scene.queryPosInRegionList(ai.owner.GetPos(), ai.owner.Attr.CollisionRadius, camp)
	if len(result) > 0 {
		return true
	}

	return false
}
