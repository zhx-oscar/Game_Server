package internal

import (
	. "Daisy/Fight/attraffix"
	"Daisy/Fight/internal/log"
	b3 "github.com/magicsea/behavior3go"
	b3config "github.com/magicsea/behavior3go/config"
	b3core "github.com/magicsea/behavior3go/core"
	"math/rand"
)

//ChangeMass 改变质量
type ChangeMass struct {
	b3core.Action
	massId int //百分比
}

func (bev *ChangeMass) Initialize(setting *b3config.BTNodeCfg) {
	bev.Action.Initialize(setting)
	bev.massId = setting.GetPropertyAsInt("massId")
}

// OnTick 循环
func (bev *ChangeMass) OnTick(tick *b3core.Tick) b3.Status {
	pawn, b := tick.GetTarget().(*Pawn)
	if !b {
		log.Error("*************[ChangeMass] trans Pawn fail")
		return b3.ERROR
	}

	pawn.pushActionTorsionInfo(bev.GetName())

	//质量为0 不做调整
	if bev.massId == 0 {
		return b3.SUCCESS
	}

	pawn.Attr.OverrideAttr(Field_Mass, float64(bev.massId))

	return b3.SUCCESS
}

//RandSuccess 随机节点
type RandSuccess struct {
	b3core.Action
	blackboardkey string //百分比
}

func (bev *RandSuccess) Initialize(setting *b3config.BTNodeCfg) {
	bev.Action.Initialize(setting)
	bev.blackboardkey = setting.GetPropertyAsString("blackboardkey")
}

// OnTick 循环
func (bev *RandSuccess) OnTick(tick *b3core.Tick) b3.Status {
	pawn, b := tick.GetTarget().(*Pawn)
	if !b {
		log.Error("*************[RandSuccess] trans Pawn fail")
		return b3.ERROR
	}

	pawn.pushActionTorsionInfo(bev.GetName())

	value := rand.Float64()
	blackBoardValue, ok := pawn.getBlackboardValueByKey(bev.blackboardkey).(float64)
	if ok && blackBoardValue >= value {
		return b3.SUCCESS
	}

	return b3.FAILURE
}

//ChangeForm 通过填写的技能索引找到对应技能并且写入黑板
type ChangeForm struct {
	b3core.Action
	form int
}

func (bev *ChangeForm) Initialize(setting *b3config.BTNodeCfg) {
	bev.Action.Initialize(setting)
	bev.form = setting.GetPropertyAsInt("form")
}

func (bev *ChangeForm) OnOpen(tick *b3core.Tick) {
	/*pawn, b := tick.GetTarget().(*Pawn)
	if !b {
		log.Error("*************[ChangeForm] trans Pawn fail")
		return
	}

	//是否可以转态
	if !bev.canChangeForm(tick) {
		return
	}

	//当前形态和行为树对应形态一致
	if pawn.Attr.Form == uint32(bev.form) {
		return
	}

	//当前形态大于0，说明需要形态转变
	if bev.form > 0 {
		formInfo, ok := pawn.Scene.GetFormDataConf(uint32(bev.form))
		if !ok {
			//查找不到对应配置，形态转变失败
			log.Error("*************[ChangeForm] fail not find formConfigJson ", bev.form)
			return
		}

		//形态转变有持续时间，那么就会触发running状态
		if formInfo.FormTime > 0 {

			//同步当前时间 当前毫秒
			tick.Blackboard.Set(ChangeFormBeginTime, int64(pawn.Scene.NowTime), "", "")
			tick.Blackboard.Set(ChangeFormEnd, false, "", "")
		}

		pawn.Attr.OverrideAttr(Field_Form, float64(bev.form))
	}*/
}

// OnTick 循环
func (bev *ChangeForm) OnTick(tick *b3core.Tick) b3.Status {
	/*pawn, b := tick.GetTarget().(*Pawn)
	if !b {
		log.Error("*************[ChangeForm] trans Pawn fail")
		return b3.ERROR
	}

	//是否可以转态
	if !bev.canChangeForm(tick) {
		return b3.FAILURE
	}

	//当前形态大于0，说明需要形态转变
	if bev.form > 0 {
		formInfo, ok := pawn.Scene.GetFormDataConf(uint32(bev.form))
		if !ok {
			//查找不到对应配置，形态转变失败
			log.Error("*************[ChangeForm] fail not find formConfigJson ", bev.form)
			return b3.FAILURE
		}

		//形态转变有持续时间，那么就会触发running状态
		if formInfo.FormTime > 0 {
			//ChangeForm达到持续时间
			beginTime := tick.Blackboard.GetInt64(ChangeFormBeginTime, "", "")
			if int64(pawn.Scene.NowTime) > beginTime+formInfo.FormTime {
				if tick.Blackboard.GetBool(ChangeFormEnd, "", "") {
					return b3.SUCCESS
				}

				pawn.blackboardSetValue(ChangeFormEnd, true)
				return b3.SUCCESS
			}

			return b3.RUNNING
		}
	}*/

	return b3.SUCCESS
}

//canChangeForm 是否可以转态
func (bev *ChangeForm) canChangeForm(tick *b3core.Tick) bool {
	//pawn, b := tick.GetTarget().(*Pawn)
	//if !b {
	//	log.Error("*************[ChangeForm] trans Pawn fail")
	//	return false
	//}
	//
	////处于被击状态，不可转态
	//if pawn.State.BeHitStat != 0 {
	//	return false
	//}

	return true
}

//ResetAI 切换AI
type ResetAI struct {
	b3core.Action
	AIID uint32
}

func (bev *ResetAI) Initialize(setting *b3config.BTNodeCfg) {
	bev.Action.Initialize(setting)
	bev.AIID = uint32(setting.GetPropertyAsInt("AIID"))
}

// OnTick 循环
func (bev *ResetAI) OnTick(tick *b3core.Tick) b3.Status {
	pawn, b := tick.GetTarget().(*Pawn)
	if !b {
		log.Error("*************[ResetAI] trans Pawn fail")
		return b3.ERROR
	}

	pawn.pushActionTorsionInfo(bev.GetName())

	//AIID为0 不做调整
	if bev.AIID == 0 {
		return b3.SUCCESS
	}

	pawn.ResetBehavior(bev.AIID)

	return b3.SUCCESS
}

//WaitAction 等待节点
type WaitAction struct {
	b3core.Action
	endTimeKey string
}

func (bev *WaitAction) Initialize(setting *b3config.BTNodeCfg) {
	bev.Action.Initialize(setting)
	bev.endTimeKey = setting.GetPropertyAsString("endtimekey")
}

func (bev *WaitAction) OnOpen(tick *b3core.Tick) {
	pawn, b := tick.GetTarget().(*Pawn)
	if !b {
		log.Error("*************[WaitAction] trans Pawn fail")
		return
	}
	startTime := int64(pawn.Scene.NowTime)
	tick.Blackboard.Set("startTime", startTime, "", "")
}

func (bev *WaitAction) OnTick(tick *b3core.Tick) b3.Status {
	pawn, b := tick.GetTarget().(*Pawn)
	if !b {
		log.Error("*************[WaitAction] trans Pawn fail")
		return b3.ERROR
	}

	pawn.pushActionTorsionInfo(bev.GetName())

	endTime, ok := pawn.getBlackboardValueByKey(bev.endTimeKey).(float64)
	if !ok {
		log.Errorf("外部黑板没有配置迂回 endTimeKey")
		return b3.FAILURE
	}

	currTime := int64(pawn.Scene.NowTime)
	var startTime = tick.Blackboard.GetInt64("startTime", "", "")
	//fmt.Println("wait:",this.GetTitle(),tick.GetLastSubTree(),"=>", currTime-startTime)
	if currTime-startTime > int64(endTime) {
		return b3.SUCCESS
	}

	return b3.RUNNING
}

//RandWaitAction 等待节点
type RandWaitAction struct {
	b3core.Action
	beginTimeKey string
	endTimeKey   string
}

func (bev *RandWaitAction) Initialize(setting *b3config.BTNodeCfg) {
	bev.Action.Initialize(setting)
	bev.beginTimeKey = setting.GetPropertyAsString("beginTimeKey")
	bev.endTimeKey = setting.GetPropertyAsString("endTimeKey")
}

func (bev *RandWaitAction) OnOpen(tick *b3core.Tick) {
	pawn, b := tick.GetTarget().(*Pawn)
	if !b {
		log.Error("*************[RandWaitAction] trans Pawn fail")
		return
	}

	//节点开启时间记录
	startTime := int64(pawn.Scene.NowTime)
	tick.Blackboard.Set(RandWaitActionStartTime, startTime, "", "")

	beginTime, endTime, ok := bev.getWaitTimeInterval(tick)
	if !ok {
		return
	}

	//时间区间数据错误
	if beginTime < 0 || endTime < 0 {
		return
	}

	//节点随机持续时间计算
	var randTime int64
	if beginTime >= endTime {
		randTime = beginTime
	} else {
		randTime = rand.Int63n(endTime-beginTime+1) + beginTime
	}
	tick.Blackboard.Set(RandWaitActionRandTime, randTime, "", "")
}

func (bev *RandWaitAction) OnTick(tick *b3core.Tick) b3.Status {
	pawn, b := tick.GetTarget().(*Pawn)
	if !b {
		log.Error("*************[RandWaitAction] trans Pawn fail")
		return b3.ERROR
	}

	pawn.pushActionTorsionInfo(bev.GetName())

	beginTime, endTime, ok := bev.getWaitTimeInterval(tick)
	if !ok {
		return b3.FAILURE
	}

	//时间区间配置错误
	if beginTime > endTime {
		return b3.FAILURE
	}

	//当有可用必杀超能技的时候 停止等待
	if pawn.getUseableSkill(true) != nil {
		return b3.SUCCESS
	}

	currTime := int64(pawn.Scene.NowTime)
	var startTime = tick.Blackboard.GetInt64(RandWaitActionStartTime, "", "")
	var randTime = tick.Blackboard.GetInt64(RandWaitActionRandTime, "", "")
	if currTime-startTime > randTime {
		return b3.SUCCESS
	}

	return b3.RUNNING
}

//getWaitTimeInterval 获取等待时间
func (bev *RandWaitAction) getWaitTimeInterval(tick *b3core.Tick) (int64, int64, bool) {
	pawn, b := tick.GetTarget().(*Pawn)
	if !b {
		log.Error("*************[RandWaitAction] trans Pawn fail")
		return 0, 0, false
	}

	beginTime, ok := pawn.getBlackboardValueByKey(bev.beginTimeKey).(float64)
	if !ok {
		log.Errorf("外部黑板没有配置迂回 beginTimeKey")
		return 0, 0, false
	}

	endTime, ok := pawn.getBlackboardValueByKey(bev.endTimeKey).(float64)
	if !ok {
		log.Errorf("外部黑板没有配置迂回 endTimeKey")
		return 0, 0, false
	}

	return int64(beginTime), int64(endTime), true
}
