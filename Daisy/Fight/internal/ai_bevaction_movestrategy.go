package internal

import (
	"Cinder/Base/linemath"
	"Daisy/Fight/internal/log"
	"Daisy/Proto"
	b3 "github.com/magicsea/behavior3go"
	b3config "github.com/magicsea/behavior3go/config"
	b3core "github.com/magicsea/behavior3go/core"
	"math"
	"math/rand"
)

//DetourMove 迂回移动
type DetourMove struct {
	b3core.Action
	radiuskey          string //半径
	anglekey           string //半开弧度
	durationkey        string
	targetmovedstopkey string //迂回目标移动是否立马停止 开关
}

func (bev *DetourMove) Initialize(setting *b3config.BTNodeCfg) {
	bev.Action.Initialize(setting)
	bev.radiuskey = setting.GetPropertyAsString("radiuskey")
	bev.anglekey = setting.GetPropertyAsString("anglekey")
	bev.durationkey = setting.GetPropertyAsString("durationkey")
	bev.targetmovedstopkey = setting.GetPropertyAsString("targetmovedstopkey")
}

func (bev *DetourMove) OnOpen(tick *b3core.Tick) {
	pawn, b := tick.GetTarget().(*Pawn)
	if !b {
		log.Error("*************[DetourMove] trans Pawn fail")
		return
	}

	target, ok := getPawnBoard(tick, attackTarget)
	if !ok {
		log.Error("*************[DetourMove] getPawnBoard fail")
		return
	}

	//同步当前时间 当前毫秒
	tick.Blackboard.Set(BeginMoveTime, int64(pawn.Scene.NowTime), "", "")

	//moveMode := tick.Blackboard.GetInt32(detourNextPosMoveMode, "", "")
	//nextPos, ok := getVector2Board(tick, detourNextPos)
	//if ok {
	//	//默认迂回速度
	//	v2 := nextPos.Sub(pawn.GetPos()).Normalized().Mul(pawn.Attr.LookAtSpeed)
	//	if Proto.FightEnum_MoveMode(moveMode) == Proto.FightEnum_MoveMode_LookAtBack {
	//		v2 = nextPos.Sub(pawn.GetPos()).Normalized().Mul(pawn.Attr.LookAtBackSpeed)
	//	}
	//	angle := CalcAngle(pawn.GetPos(), target.GetPos())
	//
	//	pawn.MoveToAndChangeAngle(Proto.FightEnum_MoveMode(moveMode), *nextPos, v2, float32(angle), false)
	//	tick.Blackboard.Set(CurrentHP, pawn.Attr.CurHP, "", "")
	//}
	targetPos := target.GetPos()
	tick.Blackboard.Set(detourMoveTargetLastPos, &targetPos, "", "")
}

// OnTick 循环
func (bev *DetourMove) OnTick(tick *b3core.Tick) b3.Status {
	pawn, b := tick.GetTarget().(*Pawn)
	if !b {
		log.Error("*************[DetourMove] trans Pawn fail")
		return b3.ERROR
	}

	pawn.pushActionTorsionInfo(bev.GetName())

	target, ok := getPawnBoard(tick, attackTarget)
	if !ok {
		log.Error("*************[DetourMove] getPawnBoard fail")
		bev.clearDetourCacheData(tick)
		return b3.FAILURE
	}

	nextPos, okFindNextPos := getVector2Board(tick, detourNextPos)

	//当前血量
	blackHP := tick.Blackboard.GetInt64(CurrentHP, "", "")
	//失血的时候，结束迂回
	if blackHP > pawn.Attr.CurHP {
		if okFindNextPos {
			pawn.Stop()
		}
		bev.clearDetourCacheData(tick)
		return b3.SUCCESS
	}

	//当有可用必杀超能技的时候 停止迂回
	if pawn.getUseableSkill(true) != nil {
		if okFindNextPos {
			pawn.Stop()
		}
		bev.clearDetourCacheData(tick)
		return b3.SUCCESS
	}

	//目标移动停止迂回开关
	targetmovedstop, ok := pawn.getBlackboardValueByKey(bev.targetmovedstopkey).(bool)
	if !ok {
		log.Errorf("外部黑板没有配置迂回 targetmovedstopkey")
		return b3.FAILURE
	}

	//目标移动的话 迂回结束
	targetLastPos, ok := getVector2Board(tick, detourMoveTargetLastPos)
	if targetmovedstop && ok && !targetLastPos.IsEqual(target.GetPos()) {
		if okFindNextPos {
			pawn.Stop()
		}
		bev.clearDetourCacheData(tick)
		return b3.SUCCESS
	}

	duration, ok := pawn.getBlackboardValueByKey(bev.durationkey).(float64)
	if !ok {
		log.Errorf("外部黑板没有配置迂回 Durationkey")
		return b3.FAILURE
	}

	//迂回达到持续时间
	beginTime := tick.Blackboard.GetInt64(BeginMoveTime, "", "")
	if int64(pawn.Scene.NowTime) > beginTime+int64(duration) {
		if okFindNextPos {
			pawn.Stop()
		}
		bev.clearDetourCacheData(tick)
		return b3.SUCCESS
	}

	//目标死亡，停止迂回
	if !target.IsAlive() {
		if okFindNextPos {
			pawn.Stop()
		}
		bev.clearDetourCacheData(tick)
		return b3.SUCCESS
	}

	//当自己处于地方警戒范围以内 迂回失败退出
	if pawn.InCampAlertAOERange(pawn.GetEmemyCamp()) {
		if okFindNextPos {
			pawn.Stop()
		}
		bev.clearDetourCacheData(tick)
		return b3.FAILURE
	}

	//目标技能进入后摇，停止迂回
	//if target.curSkill != nil && target.curSkill.Stat == Proto.FightEnum_SkillFlowStat_Later {
	//	if okFindNextPos {
	//		pawn.Stop()
	//	}
	//	bev.clearDetourCacheData(tick)
	//	return b3.SUCCESS
	//}
	//
	////自己进入被击,停止迂回
	//if pawn.Attr.IsBeingAttack() {
	//	if okFindNextPos {
	//		pawn.Stop()
	//	}
	//	bev.clearDetourCacheData(tick)
	//	return b3.SUCCESS
	//}

	//var dis float32
	//if ok {
	//	dis = Distance(*nextPos, pawn.GetPos())
	//	//一帧距离误差范围内
	//	if dis <= pawn.Attr.FastSpeed/float32(pawn.Scene.Get_SecFrames()) {
	//		pawn.Stop()
	//		bev.clearDetourCacheData(tick)
	//		return b3.SUCCESS
	//	}
	//}

	//fmt.Println("---------00 ", pawn.UID, okFindNextPos && !pawn.IsMoving())
	//fmt.Println("---------11 ", pawn.UID, okFindNextPos && pawn.GetPos().IsEqual(*nextPos))

	//迂回过程中，没有达到迂回目标点就停止，应该是被碰撞检测毙掉了
	if okFindNextPos && !pawn.IsMoving() && !pawn.GetPos().IsEqual(*nextPos) {
		bev.clearDetourCacheData(tick)
		return b3.SUCCESS
	}

	//没有迂回目标点 || 已经停止移动 || 已经到达迂回点   ==>重新计算下一次迂回目标点
	if !okFindNextPos || (okFindNextPos && !pawn.IsMoving()) || (okFindNextPos && pawn.GetPos().IsEqual(*nextPos)) {
		nextPos, moveMode, ok := bev.buildNextPos(tick)
		if !ok {
			bev.clearDetourCacheData(tick)
			return b3.FAILURE
		}
		tick.Blackboard.Set(detourNextPos, nextPos, "", "")
		tick.Blackboard.Set(detourNextPosMoveMode, moveMode, "", "")
		{
			v2 := nextPos.Sub(pawn.GetPos()).Normalized().Mul(pawn.Info.LookAtSpeed)
			//迂回 角度朝向永远对着目标
			angle := CalcAngle(pawn.GetPos(), target.GetPos())
			if Proto.MoveMode_Enum(moveMode) == Proto.MoveMode_LookAtBack {
				v2 = nextPos.Sub(pawn.GetPos()).Normalized().Mul(pawn.Info.LookAtBackSpeed)
				//angle = float64(pawn.GetAngle())
			}
			if pawn.canMove(false) {
				result := pawn.MoveToAndChangeAngle(Proto.MoveMode_Enum(moveMode), *nextPos, v2, float32(angle), true)
				if !result {
					log.Debug("当前处于硬值 迂回失败 ", pawn.State.CantMove)
					return b3.FAILURE
				}
			}
		}
	}

	//running过程中，每帧修正迂回朝向 始终面向目标
	angle := CalcAngle(pawn.GetPos(), target.GetPos())
	pawn.SetAngle(float32(angle), false, false, true)

	tick.Blackboard.Set(CurrentHP, pawn.Attr.CurHP, "", "")
	return b3.RUNNING
}

//clearDetourCacheData 清除迂回缓存数据
func (bev *DetourMove) clearDetourCacheData(tick *b3core.Tick) {
	tick.Blackboard.Remove(detourNextPos)
	tick.Blackboard.Remove(detourNextPosMoveMode)
	tick.Blackboard.Remove(CurrentHP)
}

func (bev *DetourMove) buildNextPos(tick *b3core.Tick) (*linemath.Vector2, int32, bool) {
	pawn, b := tick.GetTarget().(*Pawn)
	if !b {
		log.Error("*************[DetourMove] trans Pawn fail")
		return nil, 0, false
	}

	target, ok := getPawnBoard(tick, attackTarget)
	if !ok {
		return nil, 0, false
	}

	radius, ok := pawn.getBlackboardValueByKey(bev.radiuskey).(float64)
	if !ok {
		log.Errorf("外部黑板没有配置迂回 radiuskey")
		return nil, 0, false
	}

	angle, ok := pawn.getBlackboardValueByKey(bev.anglekey).(float64)
	if !ok {
		log.Errorf("外部黑板没有配置迂回 Anglekey")
		return nil, 0, false
	}
	//角度变弧度
	angle = angle * math.Pi / 180

	dis := DistancePawn(pawn, target)
	if dis < float32(radius)+pawn.Attr.CollisionRadius+target.Attr.CollisionRadius {
		//需要迂回后退
		nextPos := pawn.GetPos().Sub(target.GetPos()).Normalized().Mul(pawn.Attr.CollisionRadius + target.Attr.CollisionRadius + float32(radius) - dis).Add(pawn.GetPos())
		if !nextPos.IsEqual(pawn.GetPos()) {
			return &nextPos, int32(Proto.MoveMode_LookAtBack), true
		}
	}

	//随机数引入
	prop := rand.Float64()
	realAngle := AddAngle(float32(CalcAngle(target.GetPos(), pawn.GetPos())), float32((prop*2-1)*angle))
	v2 := pawn.GetPos().Sub(target.GetPos()).Normalized()
	v2.X = float32(math.Cos(float64(realAngle)))
	v2.Y = float32(math.Sin(float64(realAngle)))
	nextPos := v2.Mul(dis).Add(target.GetPos())
	pawn.setLookAtCenterPos(target.GetPos())
	return &nextPos, int32(Proto.MoveMode_LookAt), true
}

//MoveTo 获取当前技能
type MoveTo struct {
	b3core.Action
	IsMovingFaceTarget bool
	MinDuration        string
	MaxDuration        string
}

func (bev *MoveTo) Initialize(setting *b3config.BTNodeCfg) {
	bev.Action.Initialize(setting)
	bev.IsMovingFaceTarget = setting.GetPropertyAsBool("IsMovingFaceTarget")
	bev.MinDuration = setting.GetPropertyAsString("MinDuration")
	bev.MaxDuration = setting.GetPropertyAsString("MaxDuration")
}

func (bev *MoveTo) OnOpen(tick *b3core.Tick) {
	pawn, b := tick.GetTarget().(*Pawn)
	if !b {
		log.Error("*************[MoveTo] trans Pawn fail")
		return
	}

	target, ok := getPawnBoard(tick, attackTarget)
	if !ok {
		log.Error("*************[MoveTo] getPawnBoard fail")
		return
	}

	if bev.MinDuration != "" && bev.MaxDuration != "" {
		minDuration, ok := pawn.getBlackboardValueByKey(bev.MinDuration).(float64)
		if !ok {
			log.Error("*************[MoveTo] 外部黑板 minDuration 没有配置")
			return
		}
		maxDuration, ok := pawn.getBlackboardValueByKey(bev.MaxDuration).(float64)
		if !ok {
			log.Error("*************[MoveTo] 外部黑板 maxDuration 没有配置")
			return
		}

		//时间区间数据错误
		if minDuration < 0 || maxDuration < 0 {
			return
		}

		var randTime int64
		if minDuration >= maxDuration {
			randTime = int64(minDuration)
		} else {
			randTime = rand.Int63n(int64(maxDuration)-int64(minDuration)+1) + int64(minDuration)
		}

		//同步当前时间 当前毫秒
		tick.Blackboard.Set(BeginMoveTime, int64(pawn.Scene.NowTime), "", "")
		tick.Blackboard.Set(StopMoveTime, randTime, "", "")
	}

	pos, ok := getVector2Board(tick, attackPos)
	if ok && !pos.IsEqual(pawn.GetPos()) {
		v2 := pos.Sub(pawn.GetPos()).Normalized().Mul(pawn.Info.RunSpeed)

		tick.Blackboard.Set(lastDistanceWithTarget, Distance(target.GetPos(), pawn.GetPos()), "", "")
		angle := CalcAngle(pawn.GetPos(), *pos)
		if bev.IsMovingFaceTarget && !pawn.pos.IsEqual(target.pos) {
			angle = CalcAngle(pawn.GetPos(), target.pos)
		}

		result := pawn.MoveToAndChangeAngle(Proto.MoveMode_Run, *pos, v2, float32(angle), false)
		tick.Blackboard.Set(moveToSuccess, result, "", "")
		targetPos := target.GetPos()
		tick.Blackboard.Set(lastTargetPos, &targetPos, "", "")
	}
}

// OnTick 循环
func (bev *MoveTo) OnTick(tick *b3core.Tick) b3.Status {
	pawn, b := tick.GetTarget().(*Pawn)
	if !b {
		log.Error("*************[MoveTo] trans Pawn fail")
		return b3.ERROR
	}

	pawn.pushActionTorsionInfo(bev.GetName())

	pos, ok := getVector2Board(tick, attackPos)
	if !ok {
		log.Error("*************[MoveTo] getVector2Board fail ", attackPos)
		return b3.FAILURE
	}

	//如果目标点和自己当前位置一样，则不需要移动
	if pawn.GetPos().IsEqual(*pos) {
		tick.Blackboard.Remove(lastDistanceWithTarget)
		pawn.Stop()
		return b3.SUCCESS
	}

	if bev.MinDuration != "" && bev.MaxDuration != "" {
		//移动时间到了停止移动返回失败
		currTime := int64(pawn.Scene.NowTime)
		var startTime = tick.Blackboard.GetInt64(BeginMoveTime, "", "")
		var randTime = tick.Blackboard.GetInt64(StopMoveTime, "", "")
		if currTime-startTime > randTime {
			if pawn.IsMoving() {
				pawn.Stop()
			}
			tick.Blackboard.Remove(lastDistanceWithTarget)
			return b3.FAILURE
		}
	}

	moveToSuccess := tick.Blackboard.GetBool(moveToSuccess, "", "")
	if !moveToSuccess {
		//moveto open移动失败
		return b3.FAILURE
	}

	//lastDistance := tick.Blackboard.Get(lastDistanceWithTarget, "", "").(float32)

	target, ok := getPawnBoard(tick, attackTarget)
	if !ok {
		log.Error("*************[MoveTo] getPawnBoard fail")
		return b3.FAILURE
	}

	skill, skillok := getSkillBoard(tick, curSkill)
	//黑板中有可放技能
	//if skillok && skill.InAttackRange(target) {
	//	pawn.Stop()
	//	return b3.SUCCESS
	//}

	//当有可用必杀超能技的时候 停止移动 针对当前没有技能 或者 当前准备释放的技能为普攻的时候
	if pawn.getUseableSkill(true) != nil && (!skillok || !skill.IsUltimateSkill()) {
		if pawn.IsMoving() {
			pawn.Stop()
		}
		tick.Blackboard.Remove(lastDistanceWithTarget)
		return b3.FAILURE
	}

	//todo 等待策划逻辑清晰再调整
	////当自己处于地方警戒范围以内 当前节点失败退出
	//if pawn.InCampAlertAOERange(pawn.GetEmemyCamp()) {
	//	if pawn.IsMoving() {
	//		pawn.Stop()
	//	}
	//	tick.Blackboard.Remove(lastDistanceWithTarget)
	//	return b3.FAILURE
	//}

	nowDis := DistancePawn(pawn, target)

	////一帧距离误差范围内
	//if nowDis <= pawn.Attr.FastSpeed/float32(pawn.Scene.Get_SecFrames()) {
	//	pawn.Stop()
	//	tick.Blackboard.Remove(lastDistanceWithTarget)
	//	return b3.SUCCESS
	//}

	//fmt.Println("========== 00  ", DistancePawn(pawn, target), lastDistance)
	//目标死亡了 或者距离目标越来越远 或者已经处于攻击范围以内了 移动失败
	if !target.IsAlive() /*|| nowDis > lastDistance*/ {
		pawn.Stop()
		tick.Blackboard.Remove(lastDistanceWithTarget)
		return b3.FAILURE
	}

	//处于移动过程中
	if pawn.IsMoving() {
		tick.Blackboard.Set(lastDistanceWithTarget, nowDis, "", "")
		cacheLastTargetPos, ok := getVector2Board(tick, lastTargetPos)

		//重新计算移动目标点
		if !ok || !cacheLastTargetPos.IsEqual(target.GetPos()) {
			var skillAttackMin, skillAttackMax float32
			if skillok {
				skillAttackMin, skillAttackMax = bev.aiGetSkillAttackRange(pawn, target, skill)
			}

			////在射程内,原地不动
			//if nowDis <= skillAttackRange {
			//	pawn.Stop()
			//	return b3.SUCCESS
			//}
			baseAttackDis := tick.Blackboard.GetFloat64(randBaseAttackDis, "", "")
			_, stationPos, ok := target.getFreeNearestStation(pawn, float32(baseAttackDis), skillAttackMin, skillAttackMax, true)
			if !ok {
				//空闲站位点查找失败
				log.Debug("*************[MoveTo] getFreeNearestStation fail")
				return b3.FAILURE
			}

			//得到的就位点就是自己当前位置
			if pawn.GetPos().IsEqual(stationPos) {
				pawn.Stop()
				return b3.SUCCESS
			}

			tick.Blackboard.Set(attackPos, &stationPos, "", "")
			{
				v2 := stationPos.Sub(pawn.GetPos()).Normalized().Mul(pawn.Info.RunSpeed)

				tick.Blackboard.Set(lastDistanceWithTarget, Distance(target.GetPos(), pawn.GetPos()), "", "")
				angle := CalcAngle(pawn.GetPos(), stationPos)

				if bev.IsMovingFaceTarget && !pawn.pos.IsEqual(target.pos) {
					angle = CalcAngle(pawn.GetPos(), target.pos)
				}

				result := pawn.MoveToAndChangeAngle(Proto.MoveMode_Run, stationPos, v2, float32(angle), false)
				if !result {
					log.Debug("当前处于硬值 移动失败 ", pawn.State.CantMove)
					return b3.FAILURE
				}
			}
		} else {
			if bev.IsMovingFaceTarget && !pawn.pos.IsEqual(target.pos) {
				angle := CalcAngle(pawn.GetPos(), target.pos)
				pawn.SetAngle(float32(angle), true, true, false)
			}
		}
		return b3.RUNNING
	}

	tick.Blackboard.Remove(lastDistanceWithTarget)
	return b3.SUCCESS
}

//aiGetSkillAttackRange ai获取攻击范围 因为需要考虑是否引入 冲刺
func (bev *MoveTo) aiGetSkillAttackRange(pawn, target *Pawn, skill *_SkillItem) (float32, float32) {
	if pawn == nil || target == nil || skill == nil {
		return 0, 0
	}

	//如果和目标距离处于 攻击范围以内的话，不考虑冲刺。否则需要考虑冲刺
	dis := DistancePawn(pawn, target)
	var minCastDistance, maxCastDistance float32
	if skill.GetMinCastDistance() > 0 {
		minCastDistance = skill.GetMinCastDistance() + pawn.Attr.CollisionRadius + target.Attr.CollisionRadius
	}
	maxCastDistance = skill.GetMaxCastDistance() + pawn.Attr.CollisionRadius + target.Attr.CollisionRadius

	if dis >= minCastDistance && dis <= maxCastDistance {
		return skill.GetMinCastDistance(), skill.GetMaxCastDistance()
	}

	return skill.GetAttackRange()
}

//GetOutOfPositiveBossRangePos 跑出boss正面范围之外
type GetOutOfPositiveBossRangePos struct {
	b3core.Action
	PositiveBossAngleKey  string
	PositiveBossRadiusKey string
}

func (bev *GetOutOfPositiveBossRangePos) Initialize(setting *b3config.BTNodeCfg) {
	bev.Action.Initialize(setting)
	bev.PositiveBossAngleKey = setting.GetPropertyAsString("PositiveBossAngleKey")
	bev.PositiveBossRadiusKey = setting.GetPropertyAsString("PositiveBossRadiusKey")
}

// OnTick 循环
func (bev *GetOutOfPositiveBossRangePos) OnTick(tick *b3core.Tick) b3.Status {
	pawn, b := tick.GetTarget().(*Pawn)
	if !b {
		log.Error("*************[GetOutOfPositiveBossRangePos] trans Pawn fail")
		return b3.ERROR
	}

	pawn.pushActionTorsionInfo(bev.GetName())

	target, ok := getPawnBoard(tick, attackTarget)
	if !ok {
		log.Error("*************[GetOutOfPositiveBossRangePos] getPawnBoard fail")
		return b3.FAILURE
	}

	angle, ok := pawn.getBlackboardValueByKey(bev.PositiveBossAngleKey).(float64)
	if !ok {
		log.Error("*************[GetOutOfPositiveBossRangePos] 外部黑板 PositiveBossAngleKey 没有配置")
		return b3.ERROR
	}
	radius, ok := pawn.getBlackboardValueByKey(bev.PositiveBossRadiusKey).(float64)
	if !ok {
		log.Error("*************[GetOutOfPositiveBossRangePos] 外部黑板 PositiveBossRadiusKey 没有配置")
		return b3.ERROR
	}
	//角度转化为弧度
	angle = angle * math.Pi / 180

	//randAngle := float64(target.GetAngle()) + rand.Float64()*(math.Pi*2-angle) + angle/2
	//randAngle = math.Mod(randAngle, math.Pi*2)

	//pos := linemath.Vector2{
	//	X: float32(math.Cos(randAngle)),
	//	Y: float32(math.Sin(randAngle)),
	//}
	//pos = pos.Mul(float32(radius)).Add(target.GetPos())
	pos, ok := target.getOutOfPositiveRangePos(float32(radius)+pawn.Attr.CollisionRadius+target.Attr.CollisionRadius, float32(angle), pawn)
	if !ok {
		log.Debug("*************[getOutOfPositiveRangePos] 后背点未找到")
		return b3.FAILURE
	}

	tick.Blackboard.Set(attackPos, &pos, "", "")
	return b3.SUCCESS
}

//GetAttackPos 通过当前技能获取适当的攻击位置
type GetAttackPos struct {
	b3core.Action
}

func (bev *GetAttackPos) Initialize(setting *b3config.BTNodeCfg) {
	bev.Action.Initialize(setting)
}

// OnTick 循环
func (bev *GetAttackPos) OnTick(tick *b3core.Tick) b3.Status {
	pawn, b := tick.GetTarget().(*Pawn)
	if !b {
		log.Error("*************[GetAttackPos] trans Pawn fail")
		return b3.ERROR
	}

	pawn.pushActionTorsionInfo(bev.GetName())

	target, ok := getPawnBoard(tick, attackTarget)
	if !ok {
		log.Error("*************[GetAttackPos] getPawnBoard fail")
		return b3.FAILURE
	}

	//目标死亡了
	if !target.IsAlive() {
		return b3.FAILURE
	}

	dis := DistancePawn(pawn, target)

	skill, ok := getSkillBoard(tick, curSkill)
	if !ok {
		log.Error("*************[GetAttackPos] getSkillBoard fail")
		return b3.FAILURE
	}

	//在射程内,原地不动
	min, max := bev.aiGetSkillAttackRange(pawn, target, skill)
	if dis >= min && dis <= max {
		pos := pawn.GetPos()
		if pawn.OverlapCircleShape(pos) {
			tick.Blackboard.Set(attackPos, &pos, "", "")
			return b3.SUCCESS
		} else {
			//fmt.Println("+++")
		}
	}

	baseAttackDis := rand.Float32() * (max - min)

	station, stationPos, ok := target.getFreeNearestStation(pawn, baseAttackDis, min, max, true)
	if !ok {
		//空闲站位点查找失败
		log.Debug("*************[GetAttackPos] getFreeNearestStation fail ", pawn.UID)
		return b3.FAILURE
	}

	////用于检测 筛选出来的攻击点是否处于射程范围以内验证
	//if !skill.PosInCastDistance(stationPos) && !skill.Config.CanDash {
	//	fmt.Println("++++++++ getFreeNearestStation ", Distance(stationPos, target.pos))
	//	fmt.Println("")
	//}

	//不在射程以内
	//realDis := dis - skill.GetAttackRange()
	//realPos := ((stationPos.Sub(pawn.GetPos())).Normalized()).Mul(realDis).Add(pawn.GetPos())

	//fmt.Println("++++++++ GetAttackPos begin")
	//fmt.Println("++++++++ UID  ", pawn.UID, target.UID)
	//fmt.Println("++++++++ CollisionRadius  ", pawn.Info.CollisionRadius, pawn.Attr.CollisionRadius, target.Info.CollisionRadius, target.Attr.CollisionRadius)
	//fmt.Println("++++++++ SkillAttackRange   ", skill.Config.SkillMain_Config.Name, pawn.Attr.CollisionRadius+target.Attr.CollisionRadius+min, pawn.Attr.CollisionRadius+target.Attr.CollisionRadius+min+max)
	//fmt.Println("++++++++ pos  ", pawn.pos, target.pos)
	//fmt.Println("++++++++ Distance ", Distance(target.pos, stationPos), Distance(pawn.pos, stationPos))
	//fmt.Println("++++++++ GetAttackPos end")

	target.setStationData(station)
	tick.Blackboard.Set(attackPos, &stationPos, "", "")
	tick.Blackboard.Set(randBaseAttackDis, float64(baseAttackDis), "", "")
	return b3.SUCCESS
}

//aiGetSkillAttackRange ai获取攻击范围 因为需要考虑是否引入 冲刺
func (bev *GetAttackPos) aiGetSkillAttackRange(pawn, target *Pawn, skill *_SkillItem) (float32, float32) {
	if pawn == nil || target == nil || skill == nil {
		return 0, 0
	}

	//如果和目标距离处于 攻击范围以内的话，不考虑冲刺。否则需要考虑冲刺
	dis := DistancePawn(pawn, target)
	var minCastDistance, maxCastDistance float32
	if skill.GetMinCastDistance() > 0 {
		minCastDistance = skill.GetMinCastDistance() + pawn.Attr.CollisionRadius + target.Attr.CollisionRadius
	}
	maxCastDistance = skill.GetMaxCastDistance() + pawn.Attr.CollisionRadius + target.Attr.CollisionRadius

	if dis >= minCastDistance && dis <= maxCastDistance {
		return skill.GetMinCastDistance(), skill.GetMaxCastDistance()
	}

	return skill.GetAttackRange()
}

//MoveToOutOfPositivePos 跑开目标正面范围以外
type MoveToOutOfPositivePos struct {
	b3core.Action
	PositiveBossAngleKey  string
	PositiveBossRadiusKey string
}

func (bev *MoveToOutOfPositivePos) Initialize(setting *b3config.BTNodeCfg) {
	bev.Action.Initialize(setting)
	bev.PositiveBossAngleKey = setting.GetPropertyAsString("PositiveBossAngleKey")
	bev.PositiveBossRadiusKey = setting.GetPropertyAsString("PositiveBossRadiusKey")
}

func (bev *MoveToOutOfPositivePos) OnOpen(tick *b3core.Tick) {
	pawn, b := tick.GetTarget().(*Pawn)
	if !b {
		log.Error("*************[MoveToOutOfPositivePos] trans Pawn fail")
		return
	}

	pos, ok := getVector2Board(tick, attackPos)
	if ok && !pos.IsEqual(pawn.GetPos()) {
		v2 := pos.Sub(pawn.GetPos()).Normalized().Mul(pawn.Info.RunSpeed)

		angle := CalcAngle(pawn.GetPos(), *pos)
		result := pawn.MoveToAndChangeAngle(Proto.MoveMode_Run, *pos, v2, float32(angle), false)
		tick.Blackboard.Set(moveToSuccess, result, "", "")
	}
}

// OnTick 循环
func (bev *MoveToOutOfPositivePos) OnTick(tick *b3core.Tick) b3.Status {
	pawn, b := tick.GetTarget().(*Pawn)
	if !b {
		log.Error("*************[MoveToOutOfPositivePos] trans Pawn fail")
		return b3.ERROR
	}

	pawn.pushActionTorsionInfo(bev.GetName())

	pos, ok := getVector2Board(tick, attackPos)
	if !ok {
		log.Error("*************[MoveToOutOfPositivePos] getVector2Board fail ", attackPos)
		return b3.FAILURE
	}

	//fmt.Println("+++++++++ ", pawn.UID, pawn.GetPos(), *pos)

	//如果目标点和自己当前位置一样，则不需要移动
	if pawn.GetPos().IsEqual(*pos) {
		return b3.SUCCESS
	}

	isMoveToSuccess := tick.Blackboard.GetBool(moveToSuccess, "", "")
	if !isMoveToSuccess {
		//moveto open移动失败
		return b3.FAILURE
	}

	//lastDistance := tick.Blackboard.Get(lastDistanceWithTarget, "", "").(float32)

	target, ok := getPawnBoard(tick, attackTarget)
	if !ok {
		log.Error("*************[MoveToOutOfPositivePos] getPawnBoard fail")
		return b3.FAILURE
	}

	//目标死亡了 或者距离目标越来越远 或者已经处于攻击范围以内了 移动失败
	if !target.IsAlive() /*|| nowDis > lastDistance*/ {
		pawn.Stop()
		return b3.FAILURE
	}

	//处于移动过程中
	if pawn.IsMoving() {

		//自己进入被击,停止移动
		if pawn.State.CantMove {
			pawn.Stop()
			return b3.SUCCESS
		}

		return b3.RUNNING
	}

	return b3.SUCCESS
}

//FlashMoveTo 闪现到目标位置
type FlashMoveTo struct {
	b3core.Action
}

func (bev *FlashMoveTo) Initialize(setting *b3config.BTNodeCfg) {
	bev.Action.Initialize(setting)
}

func (bev *FlashMoveTo) OnOpen(tick *b3core.Tick) {
	pawn, b := tick.GetTarget().(*Pawn)
	if !b {
		log.Error("*************[FlashMoveTo] trans Pawn fail")
		return
	}

	pos, ok := getVector2Board(tick, attackPos)
	if !ok {
		return
	}

	target, ok := getPawnBoard(tick, attackTarget)
	if !ok {
		log.Error("*************[FlashMoveTo] getPawnBoard fail")
		return
	}

	v2 := pos.Sub(pawn.GetPos()).Normalized().Mul(pawn.Info.FastSpeed)

	tick.Blackboard.Set(lastDistanceWithTarget, Distance(target.GetPos(), pawn.GetPos()), "", "")
	angle := CalcAngle(pawn.GetPos(), *pos)
	pawn.MoveToAndChangeAngle(Proto.MoveMode_Fast, *pos, v2, float32(angle), false)
}

// OnTick 循环
func (bev *FlashMoveTo) OnTick(tick *b3core.Tick) b3.Status {
	pawn, b := tick.GetTarget().(*Pawn)
	if !b {
		log.Error("*************[FlashMoveTo] trans Pawn fail")
		return b3.ERROR
	}

	pawn.pushActionTorsionInfo(bev.GetName())

	pos, ok := getVector2Board(tick, attackPos)
	if !ok {
		log.Error("*************[FlashMoveTo] getVector2Board fail ", attackPos)
		return b3.FAILURE
	}

	//如果目标点和自己当前位置一样，则不需要移动
	if pawn.GetPos().IsEqual(*pos) {
		//pawn.Stop()
		return b3.SUCCESS
	}

	target, ok := getPawnBoard(tick, attackTarget)
	if !ok {
		log.Error("*************[FlashMoveTo] getPawnBoard fail")
		return b3.FAILURE
	}

	//目标死亡了 或者距离目标越来越远 或者已经处于攻击范围以内了 移动失败
	if !target.IsAlive() {
		pawn.Stop()
		return b3.FAILURE
	}

	//处于移动过程中
	if pawn.IsMoving() {

		//自己进入被击,停止移动
		if pawn.State.CantMove {
			pawn.Stop()
			return b3.SUCCESS
		}
		return b3.RUNNING
	}
	return b3.SUCCESS
}

//SteeringSmoothing 转向平滑处理
type SteeringSmoothing struct {
	b3core.Action
}

func (bev *SteeringSmoothing) Initialize(setting *b3config.BTNodeCfg) {
	bev.Action.Initialize(setting)
}

func (bev *SteeringSmoothing) OnOpen(tick *b3core.Tick) {
	//pawn, b := tick.GetTarget().(*Pawn)
	//if !b {
	//	log.Error("*************[SelectFurthestEnemy] trans Pawn fail")
	//	return
	//}
	//
	//var targetAngle float32
	//if bev.IsTarget {
	//	target, ok := getPawnBoard(tick, attackTarget)
	//	if !ok {
	//		log.Error("*************[SteeringSmoothing] getPawnBoard fail")
	//		return
	//	}
	//
	//	targetAngle = float32(CalcAngle(pawn.pos, target.pos))
	//} else {
	//	pos, ok := getVector2Board(tick, attackPos)
	//	if !ok {
	//		log.Error("*************[SteeringSmoothing] getVector2Board fail")
	//		return
	//	}
	//
	//	targetAngle = float32(CalcAngle(pawn.pos, *pos))
	//}
	//
	//fmt.Println("")
	//fmt.Println("")
	//fmt.Println("")
	//
	//fmt.Println("++++++++ SteeringSmoothing open ", pawn.UID, pawn.angle*180/math.Pi, targetAngle*180/math.Pi)

	tick.Blackboard.Remove(beginTimeSteeringSmoothing)
}

// OnTick 循环
func (bev *SteeringSmoothing) OnTick(tick *b3core.Tick) b3.Status {
	pawn, b := tick.GetTarget().(*Pawn)
	if !b {
		log.Error("*************[SelectFurthestEnemy] trans Pawn fail")
		return b3.ERROR
	}

	pawn.pushActionTorsionInfo(bev.GetName())

	target, ok := getPawnBoard(tick, attackTarget)
	if !ok {
		log.Error("*************[SteeringSmoothing] getPawnBoard fail")
		return b3.FAILURE
	}

	//如果目标是自己 不用转向
	if pawn.UID == target.UID {
		return b3.SUCCESS
	}

	//目标角度
	targetAngle := float32(CalcAngle(pawn.pos, target.pos))

	//目标角度和自己当前角度一致
	if targetAngle == pawn.angle {
		//fmt.Println("----------------- SteeringSmoothing tick success ", pawn.UID, pawn.angle*180/math.Pi, pawn.Scene.NowTime)
		return b3.SUCCESS
	}

	beginTime := tick.Blackboard.GetInt32(beginTimeSteeringSmoothing, "", "")

	tick.Blackboard.Set(beginTimeSteeringSmoothing, int32(pawn.Scene.NowTime), "", "")
	if beginTime == 0 {
		return b3.RUNNING
	}

	//时间差 毫秒转化为秒
	timeDifference := float32(pawn.Scene.NowTime-uint32(beginTime)) / 1000

	//逆时针
	if bev.isAnticlockwise(pawn.angle, targetAngle) {
		diffAngle := timeDifference * pawn.Info.TurnSpeed * math.Pi / 180
		nextAngle := AddAngle(pawn.angle, diffAngle)
		if bev.isAnticlockwise(nextAngle, targetAngle) {
			pawn.SetAngle(nextAngle, true, true, false)
			//fmt.Println("----------------- SteeringSmoothing tick 逆时针  ", pawn.UID, nextAngle*180/math.Pi, targetAngle*180/math.Pi)
			return b3.RUNNING
		} else {
			pawn.SetAngle(targetAngle, true, true, false)
			//fmt.Println("----------------- SteeringSmoothing tick 逆时针 success ", pawn.UID, pawn.angle*180/math.Pi, pawn.Scene.NowTime)
			return b3.SUCCESS
		}
	} else {
		diffAngle := timeDifference * pawn.Info.TurnSpeed * math.Pi / 180
		nextAngle := AddAngle(pawn.angle, -diffAngle)
		if !bev.isAnticlockwise(nextAngle, targetAngle) {
			pawn.SetAngle(nextAngle, true, true, false)
			//fmt.Println("----------------- SteeringSmoothing tick 顺时针  ", pawn.UID, nextAngle*180/math.Pi, targetAngle*180/math.Pi)
			return b3.RUNNING
		} else {
			pawn.SetAngle(targetAngle, true, true, false)
			//fmt.Println("----------------- SteeringSmoothing tick 顺时针 success ", pawn.UID, pawn.angle*180/math.Pi, pawn.Scene.NowTime)
			return b3.SUCCESS
		}
	}
}

//isAnticlockwise 目标角度是处于自己逆时针方向
func (bev *SteeringSmoothing) fixAngle(curAngle, targetAngle float32) bool {
	curV := linemath.Vector2{
		X: float32(math.Cos(float64(curAngle))),
		Y: float32(math.Sin(float64(curAngle))),
	}

	targetV := linemath.Vector2{
		X: float32(math.Cos(float64(targetAngle))),
		Y: float32(math.Sin(float64(targetAngle))),
	}

	return curV.Cross(targetV) > 0
}

//isAnticlockwise 目标角度是处于自己逆时针方向
func (bev *SteeringSmoothing) isAnticlockwise(curAngle, targetAngle float32) bool {
	curV := linemath.Vector2{
		X: float32(math.Cos(float64(curAngle))),
		Y: float32(math.Sin(float64(curAngle))),
	}

	targetV := linemath.Vector2{
		X: float32(math.Cos(float64(targetAngle))),
		Y: float32(math.Sin(float64(targetAngle))),
	}

	return curV.Cross(targetV) > 0
}

//GetTargetBackPos 获取目标后背位置
type GetTargetBackPos struct {
	b3core.Action
}

func (bev *GetTargetBackPos) Initialize(setting *b3config.BTNodeCfg) {
	bev.Action.Initialize(setting)
}

// OnTick 循环
func (bev *GetTargetBackPos) OnTick(tick *b3core.Tick) b3.Status {
	pawn, b := tick.GetTarget().(*Pawn)
	if !b {
		log.Error("*************[GetTargetBackPos] trans Pawn fail")
		return b3.ERROR
	}

	pawn.pushActionTorsionInfo(bev.GetName())

	target, ok := getPawnBoard(tick, attackTarget)
	if !ok {
		log.Error("*************[GetTargetBackPos] getPawnBoard fail")
		return b3.FAILURE
	}

	//目标死亡了
	if !target.IsAlive() {
		return b3.FAILURE
	}

	pos, ok := pawn.Scene.GetTargetBackPos(target, pawn)
	if !ok {
		log.Error("*************[GetTargetBackPos] GetTargetBackPos fail")
		return b3.FAILURE
	}

	tick.Blackboard.Set(attackPos, pos, "", "")

	return b3.SUCCESS
}

//Retreat 后退节点
type Retreat struct {
	b3core.Action
	Durationkey string //时长限制
	Diskey      string //距离
	Anglekey    string //半开弧度
	FaceBosskey string //是否面向boss
}

func (bev *Retreat) Initialize(setting *b3config.BTNodeCfg) {
	bev.Action.Initialize(setting)
	bev.Diskey = setting.GetPropertyAsString("Diskey")
	bev.Anglekey = setting.GetPropertyAsString("Anglekey")
	bev.Durationkey = setting.GetPropertyAsString("Durationkey")
	bev.FaceBosskey = setting.GetPropertyAsString("FaceBosskey")
}

func (bev *Retreat) OnOpen(tick *b3core.Tick) {
	pawn, b := tick.GetTarget().(*Pawn)
	if !b {
		log.Error("*************[Retreat] trans Pawn fail")
		return
	}

	target, ok := getPawnBoard(tick, attackTarget)
	if !ok {
		log.Error("*************[Retreat] getPawnBoard fail")
		return
	}

	angle, ok := pawn.getBlackboardValueByKey(bev.Anglekey).(float64)
	if !ok {
		log.Error("*************[Retreat] 外部黑板 Anglekey 没有配置")
		return
	}
	dis, ok := pawn.getBlackboardValueByKey(bev.Diskey).(float64)
	if !ok {
		log.Error("*************[Retreat] 外部黑板 Diskey 没有配置")
		return
	}

	//随机左右开角
	if rand.Float32() > 0.5 {
		angle *= -1
	}

	//计算后退目标点
	retreatPos := pawn.GetPos()
	curAngle := float32(CalcAngle(pawn.GetPos(), target.GetPos()))
	curAngle = AddAngle(curAngle, float32(angle))

	if Vector2Equal(pawn.GetPos(), target.GetPos()) {
		curAngle = pawn.angle
	}

	curDis := DistancePawn(pawn, target)
	maxDis := float32(dis) + pawn.Attr.CollisionRadius + target.Attr.CollisionRadius
	//当目标距离低于需要后退距离时
	if curDis < maxDis {
		retreatPos.X = float32(math.Cos(float64(curAngle)))*maxDis + pawn.GetPos().X
		retreatPos.Y = float32(math.Sin(float64(curAngle)))*maxDis + pawn.GetPos().Y
	}

	faceBoss, ok := pawn.getBlackboardValueByKey(bev.FaceBosskey).(bool)
	if !ok {
		return
	}

	tick.Blackboard.Set(RetreatPos, &retreatPos, "", "")
	//同步当前时间 当前毫秒
	tick.Blackboard.Set(BeginMoveTime, int64(pawn.Scene.NowTime), "", "")

	//速度
	v2 := retreatPos.Sub(pawn.GetPos()).Normalized().Mul(pawn.Info.RunSpeed)
	retreatAngle := CalcAngle(pawn.GetPos(), retreatPos)
	//移动过程中时刻面对boss
	if faceBoss {
		retreatAngle = CalcAngle(pawn.GetPos(), target.pos)
	}
	result := pawn.MoveToAndChangeAngle(Proto.MoveMode_Run, retreatPos, v2, float32(retreatAngle), false)
	tick.Blackboard.Set(moveToSuccess, result, "", "")
}

// OnTick 循环
func (bev *Retreat) OnTick(tick *b3core.Tick) b3.Status {
	pawn, b := tick.GetTarget().(*Pawn)
	if !b {
		log.Error("*************[Retreat] trans Pawn fail")
		return b3.ERROR
	}

	pawn.pushActionTorsionInfo(bev.GetName())

	target, ok := getPawnBoard(tick, attackTarget)
	if !ok {
		log.Error("*************[Retreat] getPawnBoard fail")
		return b3.FAILURE
	}

	//目标死亡了
	if !target.IsAlive() {
		bev.clearRetreatCache(tick)
		return b3.FAILURE
	}

	////当自己处于地方警戒范围以内 当前节点失败退出
	//if pawn.InCampAlertAOERange(pawn.GetEmemyCamp()) {
	//	bev.clearRetreatCache(tick)
	//	return b3.FAILURE
	//}

	isMoveToSuccess := tick.Blackboard.GetBool(moveToSuccess, "", "")
	if !isMoveToSuccess {
		//Retreat open移动失败
		bev.clearRetreatCache(tick)
		return b3.FAILURE
	}

	retreatPos, ok := getVector2Board(tick, RetreatPos)
	if !ok || Vector2Equal(*retreatPos, pawn.GetPos()) {
		bev.clearRetreatCache(tick)
		return b3.SUCCESS
	}

	duration, ok := pawn.getBlackboardValueByKey(bev.Durationkey).(float64)
	if !ok {
		log.Errorf("*************[Retreat] 外部黑板没有配置 Durationkey")
		bev.clearRetreatCache(tick)
		return b3.FAILURE
	}

	//迂回达到持续时间
	beginTime := tick.Blackboard.GetInt64(BeginMoveTime, "", "")
	if int64(pawn.Scene.NowTime) > beginTime+int64(duration) {
		if pawn.IsMoving() {
			pawn.Stop()
		}
		bev.clearRetreatCache(tick)
		return b3.SUCCESS
	}

	if !pawn.IsMoving() {
		bev.clearRetreatCache(tick)
		return b3.SUCCESS
	}

	faceBoss, ok := pawn.getBlackboardValueByKey(bev.FaceBosskey).(bool)
	if !ok {
		bev.clearRetreatCache(tick)
		return b3.FAILURE
	}

	//需要时刻面对boss
	if faceBoss {
		retreatAngle := CalcAngle(pawn.GetPos(), target.pos)
		pawn.SetAngle(float32(retreatAngle), true, true, false)
	}

	return b3.RUNNING
}

func (bev *Retreat) clearRetreatCache(tick *b3core.Tick) {
	tick.Blackboard.Remove(RetreatPos)
	tick.Blackboard.Remove(BeginMoveTime)
}
