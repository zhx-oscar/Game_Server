package internal

import (
	"Cinder/Base/linemath"
	"Daisy/Proto"
	"math"
	"reflect"
)

//移动同步阈值
const (
	moveSyncThresholdDis   = 0.1                 //目标点相距0.1米才会同步action
	moveSyncThresholdAngle = 5 / 180 * math.Pi   //角度变化5度 才会同步action
	moveAccuracy           = 100                 //移动精度 小数点后两位
	hitBackSectorThreshold = 120 * math.Pi / 180 //击退扇形阈值 120度 对应弧度
)

//_PawnMovement 移动模块
type _PawnMovement struct {
	moveMode Proto.MoveMode_Enum //移动状态
	pos      linemath.Vector2    //当前位置

	targetPos    linemath.Vector2 //目标移动点
	realVelocity linemath.Vector2 //现实世界真实的速度 每帧叠加计算  v = v0+at

	velocity     linemath.Vector2 //向目标位置移动的速度   v0
	acceleration linemath.Vector2 //加速度

	angle float32 //目标朝向	0--2PI

	owner              *Pawn
	isMoving           bool    //处于移动中
	isPause            bool    //是否处于停顿中
	isPassive          bool    //是否被动强制  技能驱使就是被动  AI调用移动就是主动
	lastTargetDistance float32 //上一次目标距离

	//迂回相对中心点
	lookAtCenterPos linemath.Vector2

	hitAngle float32 //击退角度 用于 击退碰撞检测

	//移动停止的回调列表
	callBackList []interface{}

	//匀变速运动
	isUniformlyVariable bool             //是否是匀变速运动
	moveBeginTime       uint32           //移动开始时间
	moveBeginPos        linemath.Vector2 //移动开始位置
	expectMoveEndTime   uint32           //期望移动结束时间
	lastMoveAction      *Proto.MoveBegin //上一次移动开始帧数据
}

const DelDiscardMoveEndActionFlag = false //移动消息过滤开关

func (c *_PawnMovement) init(pawn *Pawn) {
	c.pos = linemath.Vector2{X: pawn.Info.BornPos.X, Y: pawn.Info.BornPos.Y}
	c.angle = pawn.Info.BornAngle
	c.owner = pawn

	//pawn 创建刚体
	pawn.Scene.createPawnShape(float64(pawn.GetPos().X), float64(pawn.GetPos().Y), float64(pawn.Attr.CollisionRadius), pawn)
}

//SetHitAngle 设置被击退方向角度
func (c *_PawnMovement) SetHitAngle(angle float32) {
	c.hitAngle = angle
}

//AddMoveEndCallBack 增加移动结束回调
func (c *_PawnMovement) AddMoveEndCallBack(f interface{}) {
	c.callBackList = append(c.callBackList, f)
}

//GetMoveMode 获取移动模式
func (c *_PawnMovement) GetMoveMode() Proto.MoveMode_Enum {
	return c.moveMode
}

//GetAngle 获取对象面部朝向
func (c *_PawnMovement) GetAngle() float32 {
	return c.angle
}

//SetAngle 设置角色面部朝向
func (c *_PawnMovement) SetAngle(angle float32, isSync, isPassive, isSteeringSmoothing bool) {
	if !c.canMove(isPassive) {
		return
	}

	c.angle = angle
	if isSync {
		c.owner.Scene.PushAction(&Proto.FixMoveData{
			SelfId:                c.owner.UID,
			Angle:                 c.angle,
			NeedSteeringSmoothing: isSteeringSmoothing,
		})
	}
}

//GetVelocity 获取速度
func (c *_PawnMovement) GetVelocity() linemath.Vector2 {
	return c.velocity
}

//setVelocity 设置速度
func (c *_PawnMovement) setVelocity(v linemath.Vector2) {
	//是否超越当前对象最大速度限制
	length := v.Len()
	if length > c.owner.Info.MoveMaxSpeed {
		v = v.Mul(c.owner.Info.MoveMaxSpeed / length)
	}

	c.velocity = v
	c.realVelocity = v
	//同步更新实际速度
	c.updateRealVelocity()
}

//GetAngleVelocity 获取朝向对应的速度向量
func (c *_PawnMovement) GetAngleVelocity(angle, speed float32) linemath.Vector2 {
	if speed > c.owner.Info.MoveMaxSpeed {
		speed = c.owner.Info.MoveMaxSpeed
	}

	return linemath.Vector2{
		X: float32(math.Cos(float64(angle))) * speed,
		Y: float32(math.Sin(float64(angle))) * speed,
	}
}

//addVelocity 增加速度
func (c *_PawnMovement) addVelocity(v linemath.Vector2) {
	c.velocity.AddS(v)
	c.setVelocity(c.velocity)
}

//resetVelocity 重置速度
func (c *_PawnMovement) resetVelocity() {
	//以默认速度
	c.velocity = linemath.Vector2{
		X: float32(math.Cos(float64(c.angle))) * c.owner.Info.MoveMaxSpeed,
		Y: float32(math.Sin(float64(c.angle))) * c.owner.Info.MoveMaxSpeed,
	}
}

//GetAcceleration 获取加速度
func (c *_PawnMovement) GetAcceleration() linemath.Vector2 {
	return c.acceleration
}

//SetAcceleration 设置加速度
func (c *_PawnMovement) SetAcceleration(v linemath.Vector2) {
	//是否超越当前对象最大加速度限制
	//length := v.Len()
	//if length > c.owner.Info.MoveMaxAcceleration {
	//	v = v.Mul(c.owner.Info.MoveMaxAcceleration / length)
	//}

	c.acceleration = v
}

//AddAcceleration 增加加速度
func (c *_PawnMovement) AddAcceleration(v linemath.Vector2) {
	c.acceleration.AddS(v)
	c.SetAcceleration(c.acceleration)
}

//MoveTo 移动到目标点
func (c *_PawnMovement) MoveTo(moveMode Proto.MoveMode_Enum, targetPos, velocity linemath.Vector2, isPassive bool) bool {
	if !c.canMove(isPassive) {
		return false
	}

	if Vector2Equal(c.pos, targetPos) {
		return true
	}

	isSync := c.checkSyncAction(moveMode, targetPos, velocity, c.angle)

	c.setVelocity(velocity)
	c.targetPos = targetPos
	c.moveMode = moveMode
	c.isMoving = true
	c.isPause = false
	c.isPassive = isPassive
	c.lastTargetDistance = Distance(c.pos, targetPos)

	if isSync {
		c.syncMoveBeginAction()
	}

	return true
}

//vector2Equal vector2精度引入 相等判断
func (c *_PawnMovement) vector2Equal(v1, v2 linemath.Vector2) bool {
	if int(v1.X*moveAccuracy) == int(v2.X*moveAccuracy) && int(v1.Y*moveAccuracy) == int(v2.Y*moveAccuracy) {
		return true
	}
	return false
}

//checkSyncAction 同步action检测
func (c *_PawnMovement) checkSyncAction(moveMode Proto.MoveMode_Enum, targetPos, velocity linemath.Vector2, angle float32) bool {
	//阈值 检测
	if c.targetPos.Sub(targetPos).Len() <= moveSyncThresholdDis && math.Abs(float64(c.angle-angle)) <= moveSyncThresholdAngle {
		return false
	}

	if c.IsMoving() && c.moveMode == moveMode {
		//目的点、速度、移动模式、一样的话，不需要重复调用移动
		if c.vector2Equal(c.targetPos, targetPos) && c.vector2Equal(c.velocity, velocity) && int(c.angle*moveAccuracy) == int(angle*moveAccuracy) {
			return false
		}
	}
	return true
}

//MoveToAndChangeAngle 移动到目标点并改变朝向
func (c *_PawnMovement) MoveToAndChangeAngle(moveMode Proto.MoveMode_Enum, targetPos, velocity linemath.Vector2, angle float32, isPassive bool) bool {
	if !c.canMove(isPassive) {
		return false
	}

	if Vector2Equal(c.pos, targetPos) {
		return true
	}

	isSync := c.checkSyncAction(moveMode, targetPos, velocity, angle)

	c.setVelocity(velocity)
	c.angle = angle
	c.targetPos = targetPos
	c.moveMode = moveMode
	c.isMoving = true
	c.isPause = false
	c.isPassive = isPassive
	c.lastTargetDistance = Distance(c.pos, targetPos)

	if isSync {
		c.syncMoveBeginAction()
	}

	return true
}

//JumpTo 瞬移到某个点
func (c *_PawnMovement) JumpTo(moveMode Proto.MoveMode_Enum, targetPos linemath.Vector2, angle float32, isPassive bool) {

	c.moveMode = moveMode
	c.setPos(targetPos)
	c.angle = angle

	c.owner.Scene.PushAction(&Proto.FixMoveData{
		SelfId:                c.owner.UID,
		Angle:                 c.angle,
		Pos:                   &Proto.Position{X: targetPos.X, Y: targetPos.Y},
		NeedSteeringSmoothing: !isPassive,
	})

	c.Stop()
}

//PauseNav 暂停移动
func (c *_PawnMovement) PauseNav() {
	c.isPause = true
}

//ResumeNav 恢复移动
func (c *_PawnMovement) ResumeNav() {
	c.isPause = false
}

//IsMoving 是否处于移动中
func (c *_PawnMovement) IsMoving() bool {
	return c.isMoving && !c.isPause
}

//setLookAtCenterPos 设置迂回相对中心点
func (c *_PawnMovement) setLookAtCenterPos(pos linemath.Vector2) {
	c.lookAtCenterPos = pos
}

//syncMoveBeginAction 同步移动开始action
func (c *_PawnMovement) syncMoveBeginAction() {
	//状态变化才需要同步
	actionMove := &Proto.MoveBegin{
		SelfId:   c.owner.UID,
		MoveMode: c.moveMode,
		Angle:    c.angle,
		Pos: &Proto.Position{
			X: c.targetPos.X,
			Y: c.targetPos.Y,
		},
		Speed:                 c.realVelocity.Len(),
		NeedSteeringSmoothing: !c.isPassive,
	}

	if c.moveMode == Proto.MoveMode_LookAt {
		actionMove.LookAtPos = &Proto.Position{
			X: c.lookAtCenterPos.X,
			Y: c.lookAtCenterPos.Y,
		}
	}

	if c.moveMode == Proto.MoveMode_HitBack {
		actionMove.ExpectMoveEndTime = c.expectMoveEndTime
		c.lastMoveAction = actionMove
	}

	c.owner.Scene.PushAction(actionMove)
}

//syncMoveEndAction 同步移动结束action
func (c *_PawnMovement) syncMoveEndAction() {
	//停止移动action
	c.owner.Scene.PushAction(&Proto.MoveEnd{
		SelfId:   c.owner.UID,
		MoveMode: c.moveMode,
		Pos:      &Proto.Position{X: c.pos.X, Y: c.pos.Y},
	})
}

//controllerUpdate 控制器每帧调用
func (c *_PawnMovement) controllerUpdate() {
	c.updatePos()
	c.updateRealVelocity()
}

//delDiscardMoveEndAction 删除废弃的移动停止action
func (c *_PawnMovement) delDiscardMoveEndAction() {
	if DelDiscardMoveEndActionFlag {
		return
	}

	isEligibleAction := func(action *Proto.FightAction) bool {
		if action == nil {
			return false
		}

		//过滤其他人
		if action.ActionMoveEnd != nil && action.ActionMoveEnd.SelfId != c.owner.UID {
			return false
		}
		if action.ActionMoveBegin != nil && action.ActionMoveBegin.SelfId != c.owner.UID {
			return false
		}

		//击退不需要过滤
		if action.ActionMoveEnd != nil && action.ActionMoveEnd.MoveMode == Proto.MoveMode_HitBack {
			return false
		}
		if action.ActionMoveBegin != nil && action.ActionMoveBegin.MoveMode == Proto.MoveMode_HitBack {
			return false
		}

		return true
	}

	if len(c.owner.Scene.replay.FrameList) > 0 {
		frame := c.owner.Scene.replay.FrameList[len(c.owner.Scene.replay.FrameList)-1]
		if len(frame.ActionList) > 0 {
			var lastMoveAction *Proto.FightAction
			//保留最后一个移动action
			for i := len(frame.ActionList) - 1; i >= 0; i-- {
				action := frame.ActionList[i]
				if action != nil && (action.Type == Proto.ActionType_MoveBegin || action.Type == Proto.ActionType_MoveEnd) {
					//是否符合条件
					if !isEligibleAction(action) {
						continue
					}

					lastMoveAction = action
					break
				}
			}

			//当第一个移动action是个moveend 且移动模式不同时保留。其他全部移除
			removeActions := map[int]bool{}
			var firstEnd bool
			if lastMoveAction != nil {
				for i := 0; i < len(frame.ActionList); i++ {
					action := frame.ActionList[i]
					if action != nil && (action.Type == Proto.ActionType_MoveBegin || action.Type == Proto.ActionType_MoveEnd) {
						//是否符合条件
						if !isEligibleAction(action) {
							continue
						}

						if action == lastMoveAction {
							return
						}

						if !firstEnd {
							firstEnd = true

							if action.ActionMoveEnd != nil {
								if lastMoveAction.ActionMoveBegin != nil && lastMoveAction.ActionMoveBegin.MoveMode != action.ActionMoveEnd.MoveMode {
									//保留
									continue
								}

								if lastMoveAction.ActionMoveEnd != nil && lastMoveAction.ActionMoveEnd.MoveMode != action.ActionMoveEnd.MoveMode {
									//保留
									continue
								}
							}
						}

						//符合移除条件
						removeActions[i] = true
					}
				}
			}

			if len(removeActions) > 0 {
				for i := len(frame.ActionList) - 1; i >= 0; i-- {
					if _, ok := removeActions[i]; ok {
						frame.ActionList = append(frame.ActionList[:i], frame.ActionList[i+1:]...)
						delete(removeActions, i)
					}
				}
			}
		}
	}
}

//updateRealVelocity 更新实际速度
func (c *_PawnMovement) updateRealVelocity() {
	c.realVelocity = c.realVelocity.Add(c.acceleration.Mul(1 / float32(c.owner.Scene.secFrames)))

	//超过当前角色最大速度限制
	if c.realVelocity.Len() > c.owner.Info.MoveMaxSpeed {
		c.realVelocity = c.realVelocity.Normalized().Mul(c.owner.Info.MoveMaxSpeed)
	}
}

//updatePos 更新位置
func (c *_PawnMovement) updatePos() {
	//处于移动状态中
	if c.IsMoving() {
		//当速度为0的时候，即停下来
		if !c.isUniformlyVariable && IsZero(float64(c.realVelocity.Len())) {
			c.Stop()
		}

		nextPos := c.getNextPos()

		//下个位置是否可以移动
		if c.isBlocked(nextPos) {
			//fmt.Println("+++++++++++++++++++++++ not canMove   ", c.owner.UID, c.pos, c.isMoving, c.isPause, c.owner.Scene.nowFrames, c.owner.Scene.maxFrames)
			c.Stop()
			return
		}

		//迂回做碰撞检测
		if c.moveMode == Proto.MoveMode_LookAtBack || c.moveMode == Proto.MoveMode_LookAt {
			targetList := c.owner.Scene.overlapCircleShape(float64(nextPos.X), float64(nextPos.Y), float64(c.owner.Attr.CollisionRadius))

			//没有碰撞目标或者碰到的是自己
			if len(targetList) > 1 || (len(targetList) == 1 && targetList[0] != nil && targetList[0].UID != c.owner.UID) {
				c.Stop()
				return
			}
		}

		//击退
		if c.moveMode == Proto.MoveMode_HitBack {
			//targetList := c.owner.Scene.overlapSectorShape(c.pos, c.owner.Attr.CollisionRadius, hitBackSectorThreshold/2, c.hitAngle)
			targetList := c.owner.Scene.overlapCircleShape(float64(nextPos.X), float64(nextPos.Y), float64(c.owner.Attr.CollisionRadius))

			//没有碰撞目标或者碰到的是自己
			if len(targetList) > 1 || (len(targetList) == 1 && targetList[0] != nil && targetList[0].UID != c.owner.UID) {
				c.Stop()
				return
			}
		}

		//技能冲刺做 针对冲刺目标的碰撞检测
		//if c.moveMode == Proto.FightEnum_MoveMode_Fast && c.moveToTarget != nil {
		//	targetList := c.owner.Scene.overlapCircleShape(float64(nextPos.X), float64(nextPos.Y), float64(c.owner.Attr.CollisionRadius))
		//	for _, pawn := range targetList {
		//		if pawn.UID == c.moveToTarget.UID {
		//			c.Stop()
		//			c.moveToTarget = nil
		//			return
		//		}
		//	}
		//}

		//当与目标点距离在一帧误差以内，则认为已经到达目标点直接设位置
		if Distance(c.pos, c.targetPos) <= c.realVelocity.Len()/float32(c.owner.Scene.secFrames) {
			//目标点可能在边界外部
			if c.isBlocked(c.targetPos) {
				c.setPos(nextPos)
			} else {
				c.setPos(c.targetPos)
			}
			c.Stop()
			return
		}

		//已经达到目的点
		if c.isReachTarget(nextPos) {
			//是否达到目标点
			if c.isBlocked(c.targetPos) {
				c.setPos(nextPos)
			} else {
				c.setPos(c.targetPos)
			}
			c.Stop()
			//fmt.Println("+++++++++++++++++++++++ isReachTarget   ", c.owner.UID, c.pos, c.isMoving, c.isPause, c.owner.Scene.nowFrames, c.owner.Scene.maxFrames)
			return
		}

		//受限不能移动
		if !c.canMove(c.isPassive) {
			c.Stop()
			return
		}

		//还在移动路程中
		c.setPos(nextPos)
	}
}

//isReachTarget 是否已经达到目标点
func (c *_PawnMovement) isReachTarget(nextPos linemath.Vector2) bool {
	curTargetDistance := Distance(nextPos, c.targetPos)
	//达到目标点
	if c.lastTargetDistance <= curTargetDistance {
		return true
	}

	c.lastTargetDistance = curTargetDistance
	return false
}

//getNextPos 通过初始位置+朝向旋转角
func (c *_PawnMovement) getNextPos() linemath.Vector2 {
	t := 1 / float32(c.owner.Scene.secFrames)
	var pos linemath.Vector2
	//匀变速运动
	if c.isUniformlyVariable {
		t = float32(c.owner.Scene.NowTime-c.moveBeginTime) / 1000
		pos.X = c.moveBeginPos.X + c.velocity.X*t + 0.5*c.acceleration.X*t*t
		pos.Y = c.moveBeginPos.Y + c.velocity.Y*t + 0.5*c.acceleration.Y*t*t
		return pos
	}

	pos = c.pos.Add(c.realVelocity.Mul(t))

	//迂回的时候下一刻位置 计算 迂回呈现圆形路径，故而取切线下一刻位置
	if c.moveMode == Proto.MoveMode_LookAt {
		//顺时针旋转
		v1 := c.lookAtCenterPos.Sub(c.pos).Normalized().Rotation(90)
		v := linemath.Vector2{}
		if v1.Dot(c.targetPos.Sub(c.pos)) > 0 {
			v = v1.Mul(c.owner.Info.LookAtSpeed)
		} else {
			v = v1.Rotation(180).Mul(c.owner.Info.LookAtSpeed)
		}

		c.setVelocity(v)
		pos = c.pos.Add(c.realVelocity.Mul(t))
	}

	//fmt.Println("getNextPos ", c.realVelocity.IsEqual(c.velocity))
	return pos
}

//Stop 停止移动
func (c *_PawnMovement) Stop() {

	//已经停止,不需要重复停止
	if !c.IsMoving() {
		return
	}

	//fmt.Printf("Stop =====  uid:%v ,MoveMode:%v,  targetPos:%v, pos:%v\n", c.owner.UID, c.moveMode, c.targetPos, c.pos)

	//移动标记重置
	c.isMoving = false
	c.isPause = false
	c.isPassive = false

	//移动目标位置
	c.targetPos = linemath.Vector2{}

	c.lastTargetDistance = 0
	c.velocity = linemath.Vector2{}
	c.acceleration = linemath.Vector2{}
	c.realVelocity = linemath.Vector2{}
	c.lookAtCenterPos = linemath.Vector2{}

	//匀变速运动相关数据重置
	c.actualMoveEnd()
	c.isUniformlyVariable = false
	c.moveBeginTime = 0
	c.moveBeginPos = linemath.Vector2{}
	c.expectMoveEndTime = 0
	c.lastMoveAction = nil

	c.syncMoveEndAction()

	c.fireMoveEndCallBack()

	//移动状态重置
	c.moveMode = Proto.MoveMode_None
}

//fireMoveEndCallBack 触发 移动结束回调
func (c *_PawnMovement) fireMoveEndCallBack() {
	if len(c.callBackList) > 0 {
		for _, callBack := range c.callBackList {
			var callArgs []reflect.Value
			callArgs = append(callArgs, reflect.ValueOf(c.owner))
			reflect.ValueOf(callBack).Call(callArgs)
		}
		c.callBackList = []interface{}{}
	}
}

//isBlocked 是否被阻挡
func (c *_PawnMovement) isBlocked(nextPos linemath.Vector2) bool {

	//targetList := c.owner.Scene.overlapCircleShape(float64(nextPos.X), float64(nextPos.Y), float64(c.owner.Attr.CollisionRadius))
	//touchWorldBoundary := c.owner.Scene.checkCircleShapeOverlapWorldBoundary(float64(nextPos.X), float64(nextPos.Y), float64(c.owner.Attr.CollisionRadius))
	//
	////没有碰撞目标或者碰到的是自己
	//if !touchWorldBoundary && (len(targetList) == 0 || (len(targetList) == 1 && targetList[0] != nil && targetList[0].UID == c.owner.UID)) {
	//	return false
	//}
	//
	//return true

	touchWorldBoundary := c.owner.Scene.checkCircleShapeOverlapWorldBoundary(float64(nextPos.X), float64(nextPos.Y), float64(c.owner.Attr.CollisionRadius))

	//没有碰撞边界
	if !touchWorldBoundary {
		return false
	}

	if !c.owner.Scene.checkPointInWorldBoundary(nextPos) {
		return false
	}

	return true
}

//canMove 是否可以移动到目标位置
func (c *_PawnMovement) canMove(isPassive bool) bool {
	//如果是个背景pawn 不应该移动
	if c.owner.IsBackground() {
		return false
	}

	//处于受限不能不主动移动。
	if !isPassive && c.owner.State.CantMove {
		return false
	}

	return true
}

//IsPassive 是否处于被动移动状态
func (c *_PawnMovement) IsPassive() bool {
	return c.isPassive
}

//setPos 设置位置，同步更新刚体位置
func (c *_PawnMovement) setPos(pos linemath.Vector2) {
	c.pos = pos
	c.owner.Scene.setPawnShapePos(c.owner.UID, float64(pos.X), float64(pos.Y))

	//fmt.Printf("指令：移动 ==========setPos =====  uid:%v ,MoveMode:%v,  targetPos:%v, pos:%v, acceleration: %v, speed: %v\n", c.owner.UID, c.moveMode, c.targetPos, c.pos, c.acceleration, c.realVelocity.Len())
}

func (c *_PawnMovement) GetPos() linemath.Vector2 {
	return c.pos
}

//UniformlyVariableMoveTo 匀变速移动
//moveMode	移动模式
//targetPos	移动目标点
//initialVelocity	初速度
//acceleration 加速度
//time	移动固定时间
//angle	移动时角色朝向
//isPassive	是否强制移动
func (c *_PawnMovement) UniformlyVariableMoveTo(moveMode Proto.MoveMode_Enum, targetPos, initialVelocity, acceleration linemath.Vector2, time, angle float32, isPassive bool) bool {
	if !c.canMove(isPassive) {
		return false
	}

	if Vector2Equal(c.pos, targetPos) {
		return true
	}

	isSync := c.checkSyncAction(moveMode, targetPos, initialVelocity, angle)

	c.setVelocity(initialVelocity)
	c.SetAcceleration(acceleration)
	c.angle = angle
	c.targetPos = targetPos
	c.moveMode = moveMode
	c.isMoving = true
	c.isPause = false
	c.isPassive = isPassive
	c.lastTargetDistance = Distance(c.pos, targetPos)
	c.moveBeginTime = c.owner.Scene.NowTime
	c.moveBeginPos = c.pos
	c.isUniformlyVariable = true
	c.expectMoveEndTime = c.owner.Scene.NowTime + uint32(time*1000)

	//fmt.Println("+++++++ UniformlyVariableMoveTo   ", c.owner.Scene.NowTime, time, Distance(c.moveBeginPos, c.targetPos), initialVelocity.Len(), acceleration.Len(), initialVelocity.Dot(acceleration))
	if isSync {
		c.syncMoveBeginAction()
	}

	return true
}

//actualMoveEnd 实际移动停止
func (c *_PawnMovement) actualMoveEnd() {
	if c.lastMoveAction != nil {
		c.lastMoveAction.ActualMoveEndTime = c.owner.Scene.NowTime
	}
}
