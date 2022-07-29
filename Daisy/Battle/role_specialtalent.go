package main

import (
	"Daisy/Const"
	"Daisy/Data"
	"Daisy/ErrorCode"
	"Daisy/Proto"
	log "github.com/cihub/seelog"
	"math"
	"time"
)

const (
	TalentTypeFixed    = 1 // 固定天赋
	TalentTypeUltimate = 2 // 终极天赋
)

// 策划约定的升级公式
const (
	TalentFixedLevelLimit    = 2 // 固定天赋升级需要特工的等级 / 2  向上取整
	TalentUltimateLevelLimit = 6 // 终极天赋升级需要特工等级 / 6    向下取整
)

const TalentPointCost = 1 // 天赋升级/学习消耗科技点点数（1点）,策划约定

// RPC_UpgradeTalent 升级基本天赋
func (user *_User) RPC_UpgradeTalent(sid uint32, tid uint32) int32 {
	return user.role.levelUpTalent(sid, tid)
}

// RPC_StudyAdvancedTalent 学习进阶天赋
func (user *_User) RPC_StudyAdvancedTalent(sid uint32, tid uint32) int32 {
	return user.role.studyTalent(sid, tid)
}

// initTalent 初始化特工天赋列表
func (r *_Role) initTalent(specialAgentID uint32) *Proto.SpecialAgentTalent {

	talent := &Proto.SpecialAgentTalent{
		TalentPoint: 0,
		TalentMap:   make(map[uint32]*Proto.TalentData),
	}

	// 读特工配置表的Talent, 把所有的天赋都存到属性系统里
	tmp, ok := Data.GetSpecialAgentConfig().SpecialAgent_ConfigItems[specialAgentID]
	if !ok {
		r.Error("[initTalent] get special agent talen failed, specialAgentID is ", specialAgentID)
		return nil
	}
	for _, val := range tmp.Talent {
		td := &Proto.TalentData{
			Unlock: false,
			Level:  0,
			Study:  false,
			GiveUp: false,
		}
		talent.TalentMap[val] = td
	}

	return talent
}

// batchCheckRedPoint 批量检查是否可升级或者可学习，只要有任一，就发红点消息
func (r *_Role) batchCheckTalentRedPoint() bool {
	specialAgent := r.getCurSpecialAgent()
	if specialAgent == nil {
		log.Error("[batchUnlockTalent] the special agent no exit ")
		return false
	}

	talent := specialAgent.Talent
	if talent == nil {
		log.Error("[batchUnlockTalent] the special agent has no talent system")
		return false
	}
	talentMap := talent.TalentMap
	if talentMap == nil {
		log.Error("[batchUnlockTalent] the special agent talent system has no talentData")
		return false
	}

	tp := r.getTalentPoint(specialAgent)

	//如果当前特工对应的科技点数为0
	if tp == 0 {
		return false
	}

	for tid, v := range talentMap {

		if r.canLevelUpTalent(v, tid, specialAgent.Base.Level, tp) {
			return true
		}

		ok, _ := r.canStudyTalent(v, specialAgent, tid)
		if ok {
			return true
		}
	}
	return false
}

// batchUnlockTalent 批量解锁上阵特工天赋，在特工升级或者学习成功以后检查
func (r *_Role) batchUnlockTalent() {
	specialAgent := r.getCurSpecialAgent()
	if specialAgent == nil {
		log.Errorf("[batchUnlockTalent] the special agent no exit ")
		return
	}

	talent := specialAgent.Talent
	if talent == nil {
		log.Errorf("[batchUnlockTalent] the special agent has no talent system")
		return
	}
	talentMap := talent.TalentMap
	if talentMap == nil {
		log.Errorf("[batchUnlockTalent] the special agent talent system has no talentData")
		return
	}

	for tid, _ := range talentMap {
		r.unlockTalent(specialAgent, tid)
	}
}

// UnlockTalent 解锁天赋
func (r *_Role) unlockTalent(specialAgent *Proto.SpecialAgent, tid uint32) int32 {
	talent := r.getTalentData(specialAgent, tid)
	if talent == nil {
		return ErrorCode.Failure
	}
	if !r.canUnlockTalent(talent, tid, specialAgent) {
		return ErrorCode.Failure
	}

	talent.Unlock = true
	r.prop.SyncUpdateTalentData(specialAgent.Base.ConfigID, tid, talent)
	return ErrorCode.Success
}

// CanUnlockTalent 能否解锁天赋, 确保talent一定不为空
func (r *_Role) canUnlockTalent(talent *Proto.TalentData, tid uint32, sagent *Proto.SpecialAgent) bool {
	if talent.Unlock {
		return false
	}

	// 确保BaseTalent_ConfigItems 一定存在这个先从基本表里找id
	tmp, ok := Data.GetTalentConfig().BaseTalent_ConfigItems[tid]
	if ok {
		unlockLevel := tmp.UnlockLevel
		if sagent.Base == nil {
			log.Error("[canUnlockTalent] 特工基础数据读取失败", tid)
			return false
		}
		level := sagent.Base.Level
		if level < unlockLevel {
			return false
		}
		return true
	}

	// 上面没有返回表示基本天赋表里没有找到该天赋，再从进阶天赋表里找
	ttmp, tok := Data.GetTalentConfig().AdvancedTalent_ConfigItems[tid]
	if tok {
		unlockLevel := ttmp.UnlockLevel
		needTalent := ttmp.UnlockCondition
		if sagent.Base == nil {
			log.Error("[canUnlockTalent] 特工基础数据读取失败", tid)
			return false
		}
		level := sagent.Base.Level
		if level < unlockLevel {
			return false
		}
		if len(needTalent) == 0 {
			return true
		}

		noNeed := false
		for _, v := range needTalent {
			// 表示不需要学习任何天赋id
			if v == 0 {
				noNeed = true
				break
			}
			myTalent := r.getTalentData(sagent, v)
			if myTalent == nil {
				//log.Error("[canUnlockTalent] 特工没有这个天赋", v)
				return false
			} else {
				if myTalent.Study {
					//log.Error("[canUnlockTalent] 特工学习了需要的这个天赋", v)
					return true
				}
			}
		}

		return noNeed
	}

	// 上面两个表格都没有找到这个天赋id，说明天赋id配置错误
	log.Error("[canUnlockTalent] 天赋id错误，不在任何一张表中", tid)
	return false
}

// levelUpTalent 升级天赋
func (r *_Role) levelUpTalent(sid uint32, tid uint32) int32 {

	sagent := r.getSpecialAgent(sid)
	if sagent == nil {
		return ErrorCode.Failure
	}
	talent := r.getTalentData(sagent, tid)
	if talent == nil {
		return ErrorCode.Failure
	}

	if sagent.Base == nil {
		log.Error("[levelUpTalent] 特工基础数据读取失败")
		return ErrorCode.Failure
	}

	tp := r.getTalentPoint(sagent)
	if !r.canLevelUpTalent(talent, tid, sagent.Base.Level, tp) {
		return ErrorCode.Failure
	}

	if tp < TalentPointCost {
		log.Warnf("[levelUpTalent] 特工科技点不足 当前科技点%d - 需要消耗科技点%d", tp, TalentPointCost)
		return ErrorCode.Failure
	}
	tp -= TalentPointCost
	talent.Level++

	r.prop.SyncUpdateTalentPoint(sid, tp)
	r.prop.SyncUpdateTalentData(sid, tid, talent)
	r.notifyTalentRedPoint()

	return ErrorCode.Success
}

// canLevelUpTalent 天赋是否可以升级
func (r *_Role) canLevelUpTalent(talent *Proto.TalentData, tid uint32, agentLevel uint32, point uint32) bool {
	if talent == nil {
		return false
	}
	if !talent.Unlock {
		log.Warn("[canLevelUpTalent] 天赋没有解锁 ", tid)
		return false
	}

	if point < TalentPointCost {
		log.Warnf("[levelUpTalent] 天赋%d升级失败 特工科技点不足 当前科技点%d - 需要消耗科技点%d", tid, point, TalentPointCost)
		return false
	}

	// 确保BaseTalent_ConfigItems 一定存在这个先从基本表里找id
	tmp, ok := Data.GetTalentConfig().BaseTalent_ConfigItems[tid]
	if !ok {
		return false
	}
	maxLevel := tmp.MaxLevel
	var limit uint32
	var min uint32

	talentType := tmp.TalentType
	switch talentType {
	case TalentTypeFixed:

		limit = uint32(math.Ceil(float64(agentLevel) / float64(TalentFixedLevelLimit)))
	case TalentTypeUltimate:
		limit = agentLevel / TalentUltimateLevelLimit
	}
	if maxLevel < limit {
		min = maxLevel
	} else {
		min = limit
	}

	if talent.Level < min {
		return true
	}
	log.Warnf("[canLevelUpTalent] 天赋%d 等级%d 已达到上限%d，无法升级", tid, talent.Level, min)
	return false
}

// studyTalent 学习天赋
func (r *_Role) studyTalent(sid uint32, tid uint32) int32 {

	sagent := r.getSpecialAgent(sid)
	if sagent == nil {
		return ErrorCode.Failure
	}
	talent := r.getTalentData(sagent, tid)
	if talent == nil {
		return ErrorCode.Failure
	}

	ok, ec := r.canStudyTalent(talent, sagent, tid)
	if !ok {
		return ec
	}

	tp := r.getTalentPoint(sagent)
	if tp < TalentPointCost {
		log.Warn("[levelUpTalent] 特工科技点不足 当前科技点%d - 需要消耗科技点%d", tp, TalentPointCost)
		return ErrorCode.Failure
	}
	tp -= TalentPointCost
	talent.Study = true

	r.prop.SyncUpdateTalentPoint(sid, tp)
	r.prop.SyncUpdateTalentData(sid, tid, talent)

	// 所有除我以外的天赋都置为放弃
	child := r.getAllChild(tid)

	for _, childId := range child {
		if childId == tid {
			continue
		}
		childTalent := r.getTalentData(sagent, childId)
		if childTalent == nil {
			log.Error("[canStudyTalent] 二选一天赋不存在于特工身上 %d， 根天赋表或者特工表配置错误", childId)
			return ErrorCode.Failure
		}
		childTalent.GiveUp = true
		r.prop.SyncUpdateTalentData(sid, childId, childTalent)
	}

	// 学习成功后就检查该特工是否可以解锁新的天赋
	r.batchUnlockTalent()

	r.notifyTalentRedPoint()
	return ErrorCode.Success
}

// canStudyTalent 天赋是否可以学习
func (r *_Role) canStudyTalent(talent *Proto.TalentData, sagent *Proto.SpecialAgent, tid uint32) (bool, int32) {

	if !talent.Unlock {
		log.Warnf("[canStudyTalent] 天赋未解锁，不能学习 %d", tid)
		return false, ErrorCode.TalentUnlock
	}

	// 检查科技点是否足够
	tp := r.getTalentPoint(sagent)
	if tp < TalentPointCost {
		log.Warnf("[levelUpTalent] 特工科技点不足 当前科技点%d - 需要消耗科技点%d", tp, TalentPointCost)
		return false, ErrorCode.TalentPointNotEnough
	}

	// 需要所有子天赋都解锁才能学习
	child := r.getAllChild(tid)

	// 没有子天赋id，不能学习（学习的天赋只能是多选一）
	if child == nil {
		return false, ErrorCode.Failure
	}
	for _, childId := range child {
		if childId == tid {
			continue
		}
		childTalent := r.getTalentData(sagent, childId)
		if childTalent == nil {
			log.Error("[canStudyTalent] 二选一天赋不存在于特工身上 %d, 根天赋表或者特工表配置错误", childId)
			return false, ErrorCode.Failure
		}
		if !childTalent.Unlock {
			log.Warn("[canStudyTalent] 子天赋未解锁，不能学习 %d", childId)
			return false, ErrorCode.Failure
		}
	}

	if (!talent.Study) && (!talent.GiveUp) {
		return true, ErrorCode.Success
	}

	log.Warn("[canStudyTalent] 天赋已学习或已放弃，不能学习 %d", tid)
	return false, ErrorCode.Failure
}

// 工具类函数

// getCurSpecialAgentId 得到当前role身上的特工id
func (r *_Role) getCurSpecialAgentID() uint32 {
	var specialId uint32
	specialId = 0
	buildMap := r.prop.Data.BuildMap
	if buildMap != nil {
		val, ok := buildMap[r.prop.Data.FightingBuildID]
		if ok {
			specialId = val.SpecialAgentID
		}
	}
	return specialId
}

// getCurSpecialAgent 得到当前role身上的特工
func (r *_Role) getCurSpecialAgent() *Proto.SpecialAgent {
	specialId := r.getCurSpecialAgentID()
	if specialId == 0 {
		log.Error("[getCurSpecialAgent] get special agent id failed")
		return nil
	}
	return r.getSpecialAgent(specialId)
}

// getSpecialAgent 根据特工id得到一个特工
func (r *_Role) getSpecialAgent(sid uint32) *Proto.SpecialAgent {
	sagents := r.prop.Data.SpecialAgentList
	if sagents == nil {
		log.Error("[getSpecialAgent] role has no special agent ")
		return nil
	}
	sagent, ok := sagents[sid]
	if !ok {
		log.Error("[getSpecialAgent] the special agent is not exit %d", sid)
		return nil
	}
	return sagent
}

// getTalentData 根据特工和天赋id得到一个天赋
func (r *_Role) getTalentData(sagent *Proto.SpecialAgent, tid uint32) *Proto.TalentData {
	if sagent == nil {
		return nil
	}
	talent := sagent.Talent
	if talent == nil {
		log.Error("[getTalentData] the special agent has no talent system")
		return nil
	}
	talentMap := talent.TalentMap
	if talentMap == nil {
		log.Error("[getTalentData] the special agent talent system has no talentData")
		return nil
	}
	talentData, tok := talentMap[tid]
	if !tok {
		log.Error("[getTalentData] the special agent this talentid error %d", tid)
		return nil
	}
	return talentData
}

// getAllChild 根据天赋id得到所归属的子天赋列表
func (r *_Role) getAllChild(tid uint32) []uint32 {
	ttmp, tok := Data.GetTalentConfig().AdvancedTalent_ConfigItems[tid]
	if !tok {
		//log.Error("[canStudyTalent] 进阶表里没有找到该天赋 ", tid)
		return nil
	}
	root := ttmp.RootAdvancedTalentID
	tmp, ok := Data.GetTalentConfig().RootAdvancedTalent_ConfigItems[root]
	if !ok {
		log.Error("[getAllChild] 根天赋表里没有找到索引 ", root)
		return nil
	}

	return tmp.ChildAdvandedTalent
}

// getTalentPoint 根据特工得到所拥有的科技点数
func (r *_Role) getTalentPoint(sagent *Proto.SpecialAgent) uint32 {
	if sagent == nil {
		return 0
	}
	talent := sagent.Talent
	if talent == nil {
		log.Error("[unlockTalent] the special agent has no talent system")
		return 0
	}
	return talent.TalentPoint
}

// getCurSpecialAgentTalentBuff 拿到当前特工所拥有的天赋buff
func (r *_Role) getCurSpecialAgentTalentBuff() []uint32 {
	var specialId uint32
	specialId = 0
	buildMap := r.prop.Data.BuildMap
	if buildMap != nil {
		val, ok := buildMap[r.prop.Data.FightingBuildID]
		if ok {
			specialId = val.SpecialAgentID
		}
	}
	if specialId == 0 {
		log.Error("[getCurSpecialAgentTalentBuff] get special agent id failed")
		return nil
	}
	sagent := r.getSpecialAgent(specialId)
	if sagent == nil {
		return nil
	}

	talentSys := sagent.Talent
	if talentSys == nil {
		return nil
	}

	buffs := make([]uint32, 0)
	for k, v := range talentSys.TalentMap {
		if !v.Unlock {
			continue
		}
		buffList := r.getTalentBuffid(k, v.Level)

		for _, buffid := range buffList {
			if buffid == 0 {
				continue
			}
			buffs = append(buffs, buffid)
		}
	}
	return buffs
}

// getTalentBuffid 获得单个天赋的buffid
func (r *_Role) getTalentBuffid(tid uint32, level uint32) []uint32 {
	// 确保BaseTalent_ConfigItems 一定存在这个先从基本表里找id
	buffList := make([]uint32, 0)
	tmp, ok := Data.GetTalentConfig().BaseTalent_ConfigItems[tid]
	if ok {
		buffs := tmp.BuffID
		if uint32(len(buffs)) < level {
			log.Error("[getTalentBuffid]基本表配置错误，没有这个等级的buffid", tid, level)
			return nil
		}

		// 天赋等级升级到了1级才可以装备buff
		if level >= 1 {
			buffList = append(buffList, buffs[level-1])
		}
		return buffList
	}
	// 上面没有返回表示基本天赋表里没有找到该天赋，再从进阶天赋表里找
	ttmp, tok := Data.GetTalentConfig().AdvancedTalent_ConfigItems[tid]
	if tok {
		sagent := r.getCurSpecialAgent()
		talentData := r.getTalentData(sagent, tid)
		if !talentData.Study {
			return nil
		}
		return ttmp.SkillID
	}
	log.Error("[getTalentBuffid] 所有表格里都没找到该天赋", tid)
	return nil
}

// notifyTalentRedPoint 增加或者移除红点科技系统提示
func (r *_Role) notifyTalentRedPoint() {
	notify := r.getRedPointKey(Const.RedPointType_notifyTalent, "")
	_, ok := r.prop.Data.RedPointsData[notify]
	canNotify := r.batchCheckTalentRedPoint()

	// 没有红点又可以学习或者升级，增加红点
	if !ok && canNotify {
		redPointData := &Proto.RedPointInfo{
			Value:      1,
			CreateTime: time.Now().Unix(),
		}
		r.prop.SyncAddRedPoint(notify, redPointData)
	}

	// 有红点又不能学习或者升级，移除红点
	if ok && !canNotify {
		r.prop.SyncRemoveRedPoint(notify)
	}
}

// addTalentPointByLevelUp 升级增加天赋科技点
func (r *_Role) addTalentPointByLevelUp(beforeLevel uint32, sid uint32) {
	specialAgent := r.getSpecialAgent(sid)

	if specialAgent == nil {
		return
	}
	var tp uint32
	tp = 0
	for i := beforeLevel; i < specialAgent.Base.Level; i++ {
		//获取当前特工下一个等级的配置
		nextCfg, ok := Data.GetSpecialAgentConfig().Upgrade_ConfigItems[i+1]
		if !ok {
			break
		}
		tp += nextCfg.TechnologyPointNum
	}

	if tp == 0 {
		return
	}

	if specialAgent.Talent != nil {
		r.prop.SyncUpdateTalentPoint(specialAgent.Base.ConfigID, specialAgent.Talent.TalentPoint+tp)
	}
	r.batchUnlockTalent()    // 批量解锁天赋
	r.notifyTalentRedPoint() // 天赋红点提示
}
