package internal

import (
	"Daisy/Fight/internal/conf"
	"Daisy/Fight/internal/log"
	"Daisy/Proto"
	"fmt"
)

// _BuffFlow buff流程
type _BuffFlow struct {
	scene *Scene
}

// init 初始化
func (flow *_BuffFlow) init(scene *Scene) {
	flow.scene = scene
}

// update 帧更新
func (flow *_BuffFlow) update() {
	for _, pawn := range flow.scene.pawnList {
		if !pawn.IsAlive() {
			continue
		}

		for _, buff := range pawn.buffList {
			flow.updateOne(buff)
		}

		for index := len(pawn.buffList) - 1; index >= 0; index-- {
			buff := pawn.buffList[index]
			if buff.IsDestroy {
				pawn.buffList = append(pawn.buffList[:index], pawn.buffList[index+1:]...)
			}
		}
	}
}

// updateOne 帧更新一个buff
func (flow *_BuffFlow) updateOne(buff *Buff) {
	if buff == nil || buff.IsDestroy {
		return
	}

	// 检测buff是否结束
	for _, disappearType := range buff.Config.DisappearType {
		if disappearType == conf.BuffDisappearType_Duration {
			if flow.scene.NowTime >= buff.destroyTime {
				flow.removeBuff(buff.Pawn, buff, false)
				return
			}
		}
	}

	// 发送事件
	buff.Pawn.Events.EmitBuffUpdate(buff)
}

// getBuff 查询buff
func (flow *_BuffFlow) getBuff(pawn *Pawn, buffKey BuffKey) (*Buff, bool) {
	if pawn == nil {
		return nil, false
	}

	_, ok := pawn.Scene.GetBuffConfig(buffKey.MainID)
	if !ok {
		return nil, false
	}

	buff, ok := pawn.buffTab[buffKey]

	return buff, ok
}

// getSameBuffs 查询相同id的buff列表
func (flow *_BuffFlow) getSameBuffs(pawn *Pawn, mainID uint32) ([]*Buff, bool) {
	if pawn == nil {
		return nil, false
	}

	_, ok := pawn.Scene.GetBuffConfig(mainID)
	if !ok {
		return nil, false
	}

	if buffList, ok := pawn.sameBuffList[mainID]; ok {
		if len(*buffList) <= 0 {
			return nil, false
		}

		return *buffList, true
	}

	return nil, false
}

// addBuff 添加buff
func (flow *_BuffFlow) addBuff(pawn, caster *Pawn, mainID uint32, extValue int64) (*Buff, bool) {
	if pawn == nil || caster == nil || !pawn.IsAlive() || mainID <= 0 {
		return nil, false
	}

	// 查询buff main配置0.
	buffConf, ok := pawn.Scene.GetBuffConfig(mainID)
	if !ok {
		flow.scene.PushDebugInfo(func() string {
			return fmt.Sprintf("${PawnID:%d}无法向${PawnID:%d}施加Buff${BuffID:%d}，未找到Buff配置",
				caster.UID,
				pawn.UID,
				mainID)
		})
		return nil, false
	}

	// 记录debug信息
	flow.scene.PushDebugInfo(func() string {
		return fmt.Sprintf("${PawnID:%d}向${PawnID:%d}施加Buff${BuffID:%d}",
			caster.UID,
			pawn.UID,
			mainID)
	})

	// buff key
	buffKey := BuffKey{
		MainID: mainID,
		UID:    flow.scene.generateUID(),
	}

	// 处理buff关系逻辑
	switch buffConf.BuffKind {
	case conf.BuffKind_Incr:
		// 不能施加增益类buff
		if pawn.State.CantBeAddIncrBuff {
			// 记录debug信息
			flow.scene.PushDebugInfo(func() string {
				return fmt.Sprintf("${PawnID:%d}取消了被${PawnID:%d}施加的Buff${BuffID:%d}，不能施加增益类Buff",
					pawn.UID,
					caster.UID,
					mainID)
			})
			return nil, false
		}
	case conf.BuffKind_Decr:
		// 不能施加减益类buff
		if pawn.State.CantBeAddDecrBuff {
			// 记录debug信息
			flow.scene.PushDebugInfo(func() string {
				return fmt.Sprintf("${PawnID:%d}取消了被${PawnID:%d}施加的Buff${BuffID:%d}，不能施加减益类Buff",
					pawn.UID,
					caster.UID,
					mainID)
			})
			return nil, false
		}
	}

	// 已存在的同名buff列表
	var existBuffList []*Buff

	// 查询buff是否存在
	existBuffList, _ = flow.getSameBuffs(pawn, mainID)
	// 查询buff是否存在
	if len(existBuffList) > 0 {
		// 能否叠加
		if buffConf.Overlap {
			for _, existBuff := range existBuffList {
				switch conf.BuffOverlapType(buffConf.OverlapCondition) {
				case conf.BuffOverlapType_Self:
					if existBuff.Caster.UID != caster.UID {
						// 记录debug信息
						flow.scene.PushDebugInfo(func() string {
							return fmt.Sprintf("${PawnID:%d}取消了被${PawnID:%d}施加的Buff${BuffID:%d}，buff已存在，且释放者不同，不可叠加",
								pawn.UID,
								caster.UID,
								mainID)
						})
						return nil, false
					}
				case conf.BuffOverlapType_Friend:
					if existBuff.Caster.GetCamp() != caster.GetCamp() {
						// 记录debug信息
						flow.scene.PushDebugInfo(func() string {
							return fmt.Sprintf("${PawnID:%d}取消了被${PawnID:%d}施加的Buff${BuffID:%d}，buff已存在，且释放者阵营不同，不可叠加",
								pawn.UID,
								caster.UID,
								mainID)
						})
						return nil, false
					}
				}
			}

			existBuffNum := int32(len(existBuffList))
			if existBuffNum >= buffConf.OverlapLimit {
				// 记录debug信息
				flow.scene.PushDebugInfo(func() string {
					return fmt.Sprintf("${PawnID:%d}取消了被${PawnID:%d}施加的Buff${BuffID:%d}，达到叠加上限",
						pawn.UID,
						caster.UID,
						mainID)
				})
				return nil, false
			}

			// 设置消失时间
			if buffConf.OverlapRefreshDuration {
				for _, disappearType := range buffConf.DisappearType {
					if disappearType == conf.BuffDisappearType_Duration {
						for _, existBuff := range existBuffList {
							existBuff.destroyTime = flow.scene.NowTime + existBuff.Config.Time
						}
						break
					}
				}
			}
		} else {
			for _, existBuff := range existBuffList {
				flow.removeBuff(pawn, existBuff, false)
			}
		}
	}

	// 处理替换buff
	for _, other := range pawn.buffTab {
		if buffConf.ReplaceGroup != 0 && other.Config.ReplaceGroup == buffConf.ReplaceGroup {
			flow.removeBuff(pawn, other, false)
		}
	}

	// 创建buff
	buff := &Buff{}
	if err := buff.init(buffKey, pawn, caster, extValue); err != nil {
		log.Error(err.Error())
		return nil, false
	}

	// 判断buff是否被免疫
	for _, other := range pawn.buffTab {
		for _, immuneGroup := range other.Config.ImmuneGroupID {
			if immuneGroup != 0 && buff.Config.Group == immuneGroup {
				return nil, false
			}
		}
	}

	// 装载buff
	pawn.buffTab[buffKey] = buff

	if sameBuffList, ok := pawn.sameBuffList[mainID]; ok {
		*sameBuffList = append(*sameBuffList, buff)
	} else {
		pawn.sameBuffList[mainID] = &BuffList{buff}
	}

	pawn.buffList = append(pawn.buffList, buff)

	// 判断已存在buff是否需要清除
	for _, immuneGroupID := range buff.Config.ImmuneGroupID {
		if immuneGroupID == 0 {
			continue
		}

		for _, other := range pawn.buffTab {
			if other.Config.Group == immuneGroupID {
				flow.removeBuff(pawn, other, true)
			}
		}
	}

	// 设置消失时间
	for _, disappearType := range buffConf.DisappearType {
		if disappearType == conf.BuffDisappearType_Duration {
			buff.destroyTime = flow.scene.NowTime + buff.Config.Time
			break
		}
	}

	// 记录回放
	flow.scene.PushAction(&Proto.AddBuff{
		BuffKey:  buffKey.ToUint64(),
		CasterId: caster.UID,
		TargetId: pawn.UID,
		BuffId:   mainID,
	})

	// 注册buff接收事件
	pawn.Events.HookEvent(buffKey.UID, buff.effectTab)

	// 发送事件
	pawn.Events.EmitBuffAdd(buff)

	// 瞬时buff直接卸载
	for _, disappearType := range buffConf.DisappearType {
		if disappearType == conf.BuffDisappearType_Duration && buff.Config.Time <= 0 {
			flow.removeBuff(pawn, buff, false)
			break
		}
	}

	return buff, true
}

// removeBuff 删除buff
func (flow *_BuffFlow) removeBuff(pawn *Pawn, buff *Buff, clear bool) bool {
	if pawn == nil || buff.Config.MainID() <= 0 {
		return false
	}

	// 查询buff
	_, ok := flow.getBuff(pawn, buff.BuffKey)
	if !ok {
		// 记录debug信息
		flow.scene.PushDebugInfo(func() string {
			return fmt.Sprintf("${PawnID:%d}%s找不到Buff${BuffID:%d}，无法移除",
				pawn.UID,
				func() string {
					if clear {
						return "受到驱散，"
					}
					return ""
				}(),
				buff.BuffKey.MainID,
			)
		})
		return false
	}

	// 先删除buff防止重入
	delete(pawn.buffTab, buff.BuffKey)
	buff.IsDestroy = true

	if sameBuffList, ok := pawn.sameBuffList[buff.BuffKey.MainID]; ok {
		for i, sameBuff := range *sameBuffList {
			if sameBuff.BuffKey == buff.BuffKey {
				*sameBuffList = append((*sameBuffList)[:i], (*sameBuffList)[i+1:]...)
				break
			}
		}
	}

	// 记录debug信息
	flow.scene.PushDebugInfo(func() string {
		return fmt.Sprintf("${PawnID:%d}%s移除了被${PawnID:%d}施加的Buff${BuffID:%d}",
			pawn.UID,
			func() string {
				if clear {
					return "受到驱散，"
				}
				return ""
			}(),
			buff.Caster.UID,
			buff.Config.MainID())
	})

	// 记录回放
	flow.scene.PushAction(&Proto.RemoveBuff{
		BuffKey: buff.BuffKey.ToUint64(),
		SelfId:  pawn.UID,
		BuffId:  buff.Config.MainID(),
	})

	// 取消buff接收事件
	pawn.Events.UnhookEvent(buff.BuffKey.UID)

	// 发送事件
	pawn.Events.EmitBuffRemove(buff, clear)

	return true
}

// clearBuffKind 清除指定类型buff
func (flow *_BuffFlow) clearBuffKind(pawn *Pawn, kind conf.BuffKind, num int) int {
	if pawn == nil {
		return 0
	}

	// 记录debug信息
	flow.scene.PushDebugInfo(func() string {
		return fmt.Sprintf("${PawnID:%d}驱散自身%v类型的Buff",
			pawn.UID,
			kind)
	})

	count := 0

	for _, buff := range pawn.buffTab {
		if buff.Config.BuffKind == kind {
			if num > 0 {
				if count >= num {
					break
				}
			}

			if ok := flow.removeBuff(pawn, buff, true); ok {
				count++
			}
		}
	}

	return count
}

// clearBuff 驱散buff
func (flow *_BuffFlow) clearBuff(pawn *Pawn, clearGroup uint32) int {
	if pawn == nil {
		return 0
	}

	// 记录debug信息
	flow.scene.PushDebugInfo(func() string {
		return fmt.Sprintf("${PawnID:%d}驱散自身%v类型的Buff",
			pawn.UID,
			clearGroup)
	})

	count := 0

	for _, buff := range pawn.buffTab {
		if buff.Config.ClearGroup == clearGroup {
			if ok := flow.removeBuff(pawn, buff, true); ok {
				count++
			}
		}
	}

	return count
}

// refreshBuffDuration 刷新buff时长
func (flow *_BuffFlow) refreshBuffDuration(pawn *Pawn, buffKey BuffKey) bool {
	if pawn == nil || buffKey.MainID <= 0 {
		return false
	}

	// 查询buff
	buff, ok := flow.getBuff(pawn, buffKey)
	if !ok {
		return false
	}

	return flow.refreshBuffDurationEx(pawn, buff)
}

// refreshBuffDurationEx 刷新buff时长
func (flow *_BuffFlow) refreshBuffDurationEx(pawn *Pawn, buff *Buff) bool {
	if pawn == nil || buff == nil || !pawn.Equal(buff.Pawn) {
		return false
	}

	// 是否是按时长消失的buff
	if buff.IsDestroy {
		return false
	}

	durationTypeFound := false
	for _, disappearType := range buff.Config.DisappearType {
		if disappearType == conf.BuffDisappearType_Duration {
			durationTypeFound = true
			break
		}
	}

	if !durationTypeFound {
		return false
	}

	buff.destroyTime = flow.scene.NowTime + buff.Config.Time

	// 记录debug信息
	flow.scene.PushDebugInfo(func() string {
		return fmt.Sprintf("${PawnID:%d}刷新自身Buff${BuffID:%d}时长%d(ms), 更新至%d(ms)时销毁",
			pawn.UID,
			buff.Config.MainID(),
			buff.Config.Time,
			buff.destroyTime)
	})

	pawn.Events.EmitBuffChangeDuration(buff, int32(buff.Config.Time))

	return true
}

// extendBuffDuration 延长buff时长
func (flow *_BuffFlow) extendBuffDuration(pawn *Pawn, buffKey BuffKey, duration uint32) bool {
	if pawn == nil || buffKey.MainID <= 0 {
		return false
	}

	// 查询buff
	buff, ok := flow.getBuff(pawn, buffKey)
	if !ok {
		return false
	}

	return flow.extendBuffDurationEx(pawn, buff, duration)
}

// extendBuffDurationEx 延长buff时长
func (flow *_BuffFlow) extendBuffDurationEx(pawn *Pawn, buff *Buff, duration uint32) bool {
	if pawn == nil || buff == nil || !pawn.Equal(buff.Pawn) || duration <= 0 {
		return false
	}

	// 是否是按时长消失的buff
	if buff.IsDestroy {
		return false
	}

	durationTypeFound := false
	for _, disappearType := range buff.Config.DisappearType {
		if disappearType == conf.BuffDisappearType_Duration {
			durationTypeFound = true
			break
		}
	}

	if !durationTypeFound {
		return false
	}

	buff.destroyTime += duration

	// 记录debug信息
	flow.scene.PushDebugInfo(func() string {
		return fmt.Sprintf("${PawnID:%d}延长自身Buff${BuffID:%d}时长%d(ms), 更新至%d(ms)时销毁",
			pawn.UID,
			buff.Config.MainID(),
			duration,
			buff.destroyTime)
	})

	pawn.Events.EmitBuffChangeDuration(buff, int32(duration))

	return true
}

// setBuffDestroyTime 修改buff销毁时间
func (flow *_BuffFlow) setBuffDestroyTime(pawn *Pawn, buff *Buff, time uint32) bool {
	if pawn == nil || buff == nil || !pawn.Equal(buff.Pawn) {
		return false
	}

	// 是否是按时长消失的buff
	if buff.IsDestroy {
		return false
	}

	durationTypeFound := false
	for _, disappearType := range buff.Config.DisappearType {
		if disappearType == conf.BuffDisappearType_Duration {
			durationTypeFound = true
			break
		}
	}

	if !durationTypeFound {
		return false
	}

	delta := int32(time) - int32(buff.destroyTime)

	buff.destroyTime = time

	// 记录debug信息
	flow.scene.PushDebugInfo(func() string {
		return fmt.Sprintf("${PawnID:%d}设置自身Buff${BuffID:%d}销毁时间值%d",
			pawn.UID,
			buff.Config.MainID(),
			buff.destroyTime)
	})

	pawn.Events.EmitBuffChangeDuration(buff, delta)

	return true
}

// batchAddBuffs 批量添加buff
func (flow *_BuffFlow) batchAddBuffs(pawn, caster *Pawn, mainIDs []uint32, extValue int64) []*Buff {
	if pawn == nil || caster == nil || !pawn.IsAlive() || len(mainIDs) <= 0 {
		return nil
	}

	var buffList []*Buff

	for _, mainID := range mainIDs {
		buff, ok := flow.addBuff(pawn, caster, mainID, extValue)
		if !ok {
			continue
		}

		buffList = append(buffList, buff)
	}

	return buffList
}

// batchRemoveBuffs 批量删除buff
func (flow *_BuffFlow) batchRemoveBuffs(pawn *Pawn, mainIDs []uint32, casterID uint32) int {
	if pawn == nil || len(mainIDs) <= 0 {
		return 0
	}

	var count int

	for _, mainID := range mainIDs {
		for buffKey, buff := range pawn.buffTab {
			if buffKey.MainID == mainID {
				if flow.removeBuff(pawn, buff, true) {
					count++
				}
			}
		}

	}

	return count
}
