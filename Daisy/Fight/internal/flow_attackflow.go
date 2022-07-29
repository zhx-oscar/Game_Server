package internal

import (
	"Cinder/Base/linemath"
	"Daisy/Fight/internal/conf"
	"Daisy/Fight/internal/log"
	"Daisy/Proto"
	"fmt"
	"math/rand"
)

// _DelayCreateNextAttack 被触发后延迟创建的伤害体
type _DelayCreateNextAttack struct {
	delayCreateTime uint32
	frontAttack     *Attack
	index           uint32
}

// _AttackFlow 伤害体流程
type _AttackFlow struct {
	scene                     *Scene
	attackList                []*Attack
	delayCreateNextAttackList []*_DelayCreateNextAttack
}

// init 初始化
func (flow *_AttackFlow) init(scene *Scene) {
	flow.scene = scene
}

// update 帧更新
func (flow *_AttackFlow) update() {
	for _, attack := range flow.attackList {
		// 帧更新一个伤害体
		flow.updateOne(attack)
	}

	for i := len(flow.delayCreateNextAttackList) - 1; i >= 0; i-- {
		createAtk := flow.delayCreateNextAttackList[i]

		if flow.scene.NowTime < createAtk.delayCreateTime {
			continue
		}

		if int(createAtk.index) < len(createAtk.frontAttack.configTab) {
			// 创建伤害体
			nextAttack := &Attack{}
			if err := nextAttack.initNextAttack(createAtk.frontAttack, createAtk.index); err == nil {
				// 放置伤害体
				flow.putAttack(nextAttack)
			} else {
				log.Error(err.Error())
			}
		}

		flow.delayCreateNextAttackList = append(flow.delayCreateNextAttackList[:i], flow.delayCreateNextAttackList[i+1:]...)
	}

	for i := len(flow.attackList) - 1; i >= 0; i-- {
		attack := flow.attackList[i]
		if attack.IsDestroy {
			flow.attackList = append(flow.attackList[:i], flow.attackList[i+1:]...)
		}
	}
}

// updateOne 更新一个伤害体
func (flow *_AttackFlow) updateOne(attack *Attack) {
	// 检测是否已删除
	if attack == nil || attack.IsDestroy {
		return
	}

	// 是否执行hit
	execHit := true

	// 更新伤害体位置
	switch attack.Config.MoveMode {
	case conf.AttackMoveMode_FllowTarget:
		attack.Pos = attack.getTargetPos()

	case conf.AttackMoveMode_FllowCaster:
		attack.Pos = attack.Caster.GetPos()

	case conf.AttackMoveMode_FixTimeMoveToTarget:
		// 计算伤害体位置
		if attack.fixTimeMoveEnd > attack.createTime {
			prog := float32(flow.scene.NowTime-attack.createTime) / float32(attack.fixTimeMoveEnd-attack.createTime)
			if prog < 1 {
				attack.Pos = attack.spawnPos.Add(attack.getTargetPos().Sub(attack.spawnPos).Mul(prog))
			} else {
				attack.Pos = attack.getTargetPos()
			}
		} else {
			attack.Pos = attack.getTargetPos()
		}

		switch attack.Config.Type {
		case conf.AttackType_Single:
			// 单体需要移结束后才能开始命中
			if flow.scene.NowTime < attack.fixTimeMoveEnd {
				execHit = false
			}
		}
	}

	// 记录Aoe移动帧
	flow.saveAttackMoveAoe(attack)

	// 执行命中
	if execHit {
		// 命中模式
		switch attack.Config.HitMode {
		case conf.AttackHitMode_TimeLine:
			if attack.HitTimes < uint32(len(attack.Config.Hits)) {
				// 命中时间轴
				hitTimeLine := attack.Config.Hits[attack.HitTimes]

				// 检测命中时间
				if flow.scene.NowTime-attack.createTime >= attack.ZoomAttackTime(hitTimeLine.Begin+attack.hitExtendTime) {
					// 命中目标
					flow.hitTargets(attack)
					if attack.IsDestroy {
						return
					}

					attack.HitTimes++
				}
			}
		case conf.AttackHitMode_FixInterval:
			if flow.scene.NowTime >= attack.createTime+(attack.HitTimes*attack.Config.FixInterval) {
				// 命中目标
				flow.hitTargets(attack)
				if attack.IsDestroy {
					return
				}

				attack.HitTimes++
			}
		}
	}

	// 位移结束
	switch attack.Config.MoveMode {
	case conf.AttackMoveMode_FllowTarget:
		if flow.scene.NowTime >= attack.fixTimeMoveEnd {
			if !attack.isMoveEnd {
				// 位移结束后创建伤害体
				flow.delayCreateNextAttack(attack, attack.Config.OnMoveEnd)

				attack.isMoveEnd = true
			}
		}
	}

	// 检测销毁伤害体
	switch attack.Config.DestroyType {
	case conf.AttackDestroyType_LifeTime:
		if flow.scene.NowTime-attack.createTime >= attack.ZoomAttackTime(attack.Config.LifeTime+attack.hitExtendTime) {
			flow.removeAttack(attack, false, nil)
			return
		}
	case conf.AttackDestroyType_FixTimeMoveEnd:
		switch attack.Config.Type {
		case conf.AttackType_Single:
			if attack.HitTimes >= uint32(len(attack.Config.Hits)) {
				flow.removeAttack(attack, false, nil)
				return
			}
		case conf.AttackType_Aoe:
			if flow.scene.NowTime >= attack.fixTimeMoveEnd+attack.hitExtendTime {
				flow.removeAttack(attack, false, nil)
				return
			}
		}
	case conf.AttackDestroyType_HitTimes:
		if attack.HitTimes >= attack.Config.FixLimit {
			flow.removeAttack(attack, false, nil)
			return
		}
	}
}

// createSkillAttack 创建技能伤害体
func (flow *_AttackFlow) createSkillAttack(skill *Skill, index uint32) (*Attack, error) {
	attack := &Attack{}

	if err := attack.initSkillAttack(skill, index); err != nil {
		return nil, err
	}

	if err := flow.scene.putAttack(attack); err != nil {
		return nil, err
	}

	return attack, nil
}

// CreateCustomAttack 创建自定义伤害体
func (flow *_AttackFlow) CreateCustomAttack(innerAttackID conf.InnerAttackID, caster, casterSnapshot *Pawn,
	targetList []*Pawn, targetPos linemath.Vector2, callback IEffectCallback, scale float32) (*Attack, error) {
	attackConf, ok := flow.scene.GetInnerAttackTemplate(innerAttackID)
	if !ok {
		return nil, fmt.Errorf("no found inner attack %d", innerAttackID)
	}

	attack := &Attack{}

	if err := attack.init(caster, casterSnapshot, attackConf.AttackConfs, 0,
		targetList, targetPos, callback, nil, scale); err != nil {
		return nil, err
	}

	if err := flow.scene.putAttack(attack); err != nil {
		return nil, err
	}

	return attack, nil
}

// putAttack 放置伤害体
func (flow *_AttackFlow) putAttack(attack *Attack) error {
	if attack == nil {
		return fmt.Errorf("nil attack")
	}

	// 单体技能必须有目标
	switch attack.Config.Type {
	case conf.AttackType_Single:
		if len(attack.castTargets) <= 0 {
			return fmt.Errorf("no target")
		}
	}

	// 创建时间
	attack.createTime = flow.scene.NowTime

	// 当前点
	attack.Pos = attack.spawnPos

	// 固定时间移向目标
	switch attack.Config.MoveMode {
	case conf.AttackMoveMode_FixTimeMoveToTarget:
		if attack.Config.Speed > 0 {
			// 计算移动时间
			attack.fixTimeMoveEnd = attack.createTime + uint32(Distance(attack.Pos, attack.getTargetPos())/attack.Config.Speed)
		} else {
			// 使用生存时间
			attack.fixTimeMoveEnd = attack.createTime + attack.Config.LifeTime
		}
	}

	// 计算自动调整形状大小数值
	if attack.Config.Spawn.AutoExtend {
		switch attack.Config.Type {
		case conf.AttackType_Aoe:
			attack.autoExtendValue = float32(Max(float64(attack.getTargetPos().Sub(attack.spawnPos).Len()-attack.Caster.Attr.CollisionRadius), 0.1))
		}
	}

	// 插入队列
	flow.attackList = append(flow.attackList, attack)

	// 构造回放
	action := &Proto.NewAttack{
		CasterId:    attack.Caster.UID,
		Src:         attack.Src(),
		Index:       attack.index,
		AttackId:    attack.UID,
		MoveEndTime: attack.fixTimeMoveEnd,
		GroupID:     attack.groupID,
	}

	switch attack.Src() {
	case Proto.AttackSrc_Skill:
		action.SkillId = attack.Skill.Config.ValueID()
	case Proto.AttackSrc_Buff:
		action.BuffKey = attack.Buff.BuffKey.ToUint64()
		action.BuffId = attack.Buff.Config.MainID()
	case Proto.AttackSrc_Custom:
		action.ConfigId = attack.Config.ConfigID
	}

	switch attack.Config.Type {
	case conf.AttackType_Single:
		if len(attack.castTargets) > 0 {
			action.TargetId = attack.castTargets[0].UID
		}
	case conf.AttackType_Aoe:
		action.TargetPos = &Proto.Position{
			X: attack.castAoePos.X,
			Y: attack.castAoePos.Y,
		}
	}

	// 位移起点
	switch attack.Config.MoveMode {
	case conf.AttackMoveMode_FixTimeMoveToTarget:
		if attack.spawnPawn != nil {
			action.MoveBeginId = attack.spawnPawn.UID
		}
	}

	// 创建Aoe警示区域
	switch attack.Config.Type {
	case conf.AttackType_Aoe:
		if attack.Config.MoveMode == conf.AttackMoveMode_None {
			// 缩放Aoe范围
			shape := scaleShape(attack.Config.AttackArgs, attack.Scale)

			// 动态调整大小
			autoExtendShape(shape, attack.Config.AttackArgs.Spawn.AutoExtend, attack.autoExtendValue)

			// 创建区域
			region, err := flow.scene.createRegion(&_RegionInfo{
				Creator: attack.Caster,
				Flag:    RegionFlag_AoeWarning,
				Shape: _RegionShape{
					Type:     shape.Type,
					Extend:   shape.Extend,
					Radius:   shape.Radius,
					FanAngle: shape.FanAngle,
				},
				Pos:   attack.Pos,
				Angle: attack.Angle,
			})
			if err != nil {
				var skillID, buffID uint32

				switch attack.Src() {
				case Proto.AttackSrc_Skill:
					skillID = attack.Skill.UID
				case Proto.AttackSrc_Buff:
					buffID = attack.Buff.BuffKey.UID
				}

				log.Errorf("attack create aoe warning region failed, skill %d, buff %d, %s", skillID, buffID, err.Error())
			} else {
				attack.aoeWarnRegionUID = region.UID
			}
		}
	}

	// 记录回放
	flow.scene.PushAction(action)

	// 记录Aoe显示帧
	flow.saveAttackShowAoe(attack)

	// 发送事件
	attack.Caster.Events.EmitAttackInit(attack)

	return nil
}

// removeAttack 删除伤害体
func (flow *_AttackFlow) removeAttack(attack *Attack, isBreak bool, breakCaster *Pawn) {
	// 已删除
	if attack.IsDestroy {
		return
	}

	// 标记已删除
	attack.IsDestroy = true

	// 删除Aoe警示范围
	if attack.aoeWarnRegionUID > 0 {
		flow.scene.destroyRegion(attack.aoeWarnRegionUID)
	}

	// 记录回放
	flow.scene.PushAction(&Proto.DelAttack{
		AttackId: attack.UID,
	})

	// 发送事件
	attack.Caster.Events.EmitAttackDestroy(attack, isBreak, breakCaster)

	// 结束时创建伤害体
	if !isBreak {
		flow.delayCreateNextAttack(attack, attack.Config.OnFinish)
	}
}

// hitTargets 命中目标
func (flow *_AttackFlow) hitTargets(attack *Attack) {
	// aoe伤害每次命中前重新查询目标
	if attack.Config.Type == conf.AttackType_Aoe {
		attack.HitTargets = attack.SearchTargets()
	} else {
		attack.HitTargets = attack.castTargets
	}

	// 发送事件
	attack.Caster.Events.EmitAttackBeforeHitAll(attack)
	if attack.IsDestroy {
		return
	}

	// 是否有hit上限
	if attack.Config.MaxHitTarget > 0 && int32(len(attack.HitTargets)) > attack.Config.MaxHitTarget {
		// 随机hit目标
		randIdxs := rand.Perm(len(attack.HitTargets))[:attack.Config.MaxHitTarget]

		for _, i := range randIdxs {
			target := attack.HitTargets[i]

			if !target.IsAlive() {
				continue
			}

			attack.Caster.Events.EmitAttackHitTarget(attack, target)
			if attack.IsDestroy {
				return
			}
		}

	} else {
		// 循环伤害所有目标
		for _, target := range attack.HitTargets {
			if !target.IsAlive() {
				continue
			}

			attack.Caster.Events.EmitAttackHitTarget(attack, target)
			if attack.IsDestroy {
				return
			}
		}
	}

	// 发送事件
	attack.Caster.Events.EmitAttackAfterHitAll(attack)
}

// breakSkillAttacks 打断技能产生的伤害体
func (flow *_AttackFlow) breakSkillAttacks(skill *Skill, breakCaster *Pawn) {
	for i := len(flow.attackList) - 1; i >= 0; i-- {
		attack := flow.attackList[i]

		if !attack.Skill.Equal(skill) || !attack.Config.CanBreak {
			continue
		}

		flow.removeAttack(attack, true, breakCaster)
	}
}

// delayCreateNextAttack 被触发后延迟创建后续的伤害体
func (flow *_AttackFlow) delayCreateNextAttack(frontAttack *Attack, attackIDs []int) {
	for _, v := range attackIDs {
		if v < 0 || v >= len(frontAttack.configTab) {
			continue
		}

		atkConf := frontAttack.configTab[v]

		// 检测链式伤害体最大链接数量
		if atkConf.IsLink() {
			if frontAttack.linkTimes >= atkConf.MaxLinkTarget {
				continue
			}
		}

		flow.delayCreateNextAttackList = append(flow.delayCreateNextAttackList, &_DelayCreateNextAttack{
			delayCreateTime: flow.scene.NowTime + atkConf.Delay,
			frontAttack:     frontAttack,
			index:           uint32(v),
		})
	}
}

// searchAttackTargets 查询伤害体目标列表
func (flow *_AttackFlow) searchAttackTargets(atkArgs *conf.AttackArgs, pos linemath.Vector2, angle, autoExtendValue float32, casterCamp Proto.Camp_Enum, scale float32) []*Pawn {
	// 检测能查询目标
	if atkArgs.Type != conf.AttackType_Aoe && !atkArgs.IsLink() {
		return nil
	}

	// 目标阵营
	campBit := getTargetCampBit(atkArgs, casterCamp)

	shape := scaleShape(atkArgs, scale)
	if shape == nil {
		return nil
	}

	// 动态调整大小
	autoExtendShape(shape, atkArgs.Spawn.AutoExtend, autoExtendValue)

	// 形状偏移位置
	pos = TransPos(pos, getAutoOffset(atkArgs, shape), angle)

	// 查询目标
	return flow.scene.searchShapeTargets(campBit, pos, angle, shape, atkArgs.TargetCategory == conf.AttackTargetCategory_Enemy,
		atkArgs.TargetCategory == conf.AttackTargetCategory_Friend)
}

// saveAttackShowAoe 记录Aoe显示帧
func (flow *_AttackFlow) saveAttackShowAoe(attack *Attack) {
	if !flow.scene.SimulatorMode() {
		return
	}

	switch attack.Config.Type {
	case conf.AttackType_Aoe:
		if len(attack.Config.AttackTimeLine.Hits) <= 0 {
			return
		}

		shape := scaleShape(attack.Config.AttackArgs, attack.Scale)
		if shape == nil {
			return
		}

		autoExtendShape(shape, attack.Config.AttackArgs.Spawn.AutoExtend, attack.autoExtendValue)

		shapeOffset := getAutoOffset(attack.Config.AttackArgs, shape)
		shapePos := TransPos(attack.spawnPos, shapeOffset, attack.Angle)

		flow.scene.PushDebugAction(&Proto.AttackShowAoe{
			AttackId: attack.UID,
			Shape: &Proto.AttackAoeShape{
				Type: shape.Type,
				Extend: &Proto.Position{
					X: shape.Extend.X,
					Y: shape.Extend.Y,
				},
				Radius:   shape.Radius,
				FanAngle: shape.FanAngle,
			},
			Spawn: &Proto.AttackAoeTrans{
				Pos: &Proto.Position{
					X: shapePos.X,
					Y: shapePos.Y,
				},
				Angle: attack.Angle,
			},
		})
	}
}

// saveAttackShowAoe 记录Aoe移动帧
func (flow *_AttackFlow) saveAttackMoveAoe(attack *Attack) {
	if !flow.scene.SimulatorMode() || attack.isMoveEnd {
		return
	}

	switch attack.Config.Type {
	case conf.AttackType_Aoe:
		if attack.Config.MoveMode == conf.AttackMoveMode_None || len(attack.Config.AttackTimeLine.Hits) <= 0 {
			return
		}

		shape := scaleShape(attack.Config.AttackArgs, attack.Scale)
		if shape == nil {
			return
		}

		autoExtendShape(shape, attack.Config.AttackArgs.Spawn.AutoExtend, attack.autoExtendValue)

		shapeOffset := getAutoOffset(attack.Config.AttackArgs, shape)
		shapePos := TransPos(attack.Pos, shapeOffset, attack.Angle)

		action := &Proto.AttackMoveAoe{
			AttackId: attack.UID,
			Trans: &Proto.AttackAoeTrans{
				Pos: &Proto.Position{
					X: shapePos.X,
					Y: shapePos.Y,
				},
				Angle: attack.Angle,
			},
		}

		if attack.actionMoveAoe != nil {
			if FloatEqual(float64(attack.actionMoveAoe.Trans.Pos.X), float64(action.Trans.Pos.X)) &&
				FloatEqual(float64(attack.actionMoveAoe.Trans.Pos.Y), float64(action.Trans.Pos.Y)) &&
				FloatEqual(float64(attack.actionMoveAoe.Trans.Angle), float64(action.Trans.Angle)) {
				return
			}
		}

		flow.scene.PushDebugAction(action)

		attack.actionMoveAoe = action
	}
}

// getTargetCampBit 获取目标阵营
func getTargetCampBit(atkArgs *conf.AttackArgs, selfCamp Proto.Camp_Enum) (campBit Bits) {
	// 目标选择策略
	switch atkArgs.TargetCategory {
	case conf.AttackTargetCategory_Enemy:
		campBit.TurnOn(int32(GetEnemyCamp(selfCamp)))
	case conf.AttackTargetCategory_Friend:
		campBit.TurnOn(int32(selfCamp))
	}
	return
}

// autoExtendShape 获取动态形状
func autoExtendShape(shape *conf.AttackShape, autoExtend bool, autoExtendValue float32) {
	// 修改形状大小
	if autoExtend {
		shape := shape

		if autoExtendValue <= 0 {
			return
		}

		switch shape.Type {
		case Proto.AttackShapeType_Rect:
			shape.Extend.Y = autoExtendValue
		case Proto.AttackShapeType_Circle:
			fallthrough
		case Proto.AttackShapeType_Fan:
			shape.Radius = autoExtendValue
		}
	}
}

// getAutoOffset 获取形状自动偏移位置
func getAutoOffset(atkArgs *conf.AttackArgs, shape *conf.AttackShape) linemath.Vector2 {
	// 偏移位置
	shapeOffset := atkArgs.Spawn.Offset

	// 修改偏移位置
	if atkArgs.Spawn.AutoRectOffset {
		switch shape.Type {
		case Proto.AttackShapeType_Rect:
			shapeOffset.Y = shape.Extend.Y / 2
		}
	}

	return shapeOffset
}

// scaleShape 缩放形状
func scaleShape(atkArgs *conf.AttackArgs, scale float32) *conf.AttackShape {
	shape := &conf.AttackShape{
		Type:     atkArgs.Shape.Type,
		Extend:   atkArgs.Shape.Extend,
		Radius:   atkArgs.Shape.Radius,
		FanAngle: atkArgs.Shape.FanAngle,
	}

	switch shape.Type {
	case Proto.AttackShapeType_Rect:
		shape.Extend.X *= scale
		shape.Extend.Y *= scale
	case Proto.AttackShapeType_Circle:
		fallthrough
	case Proto.AttackShapeType_Fan:
		shape.Radius *= scale
	}

	return shape
}
