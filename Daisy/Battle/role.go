package main

import (
	"Cinder/Chat/chatapi"
	"Cinder/Mail/mailapi"
	"Cinder/Space"
	"Daisy/Const"
	"Daisy/DB"
	"Daisy/Data"
	"Daisy/Fight/attraffix"
	"Daisy/Item"
	"Daisy/ItemProto"
	"Daisy/Prop"
	"Daisy/Proto"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const RoleActorType = "Role"

type _Role struct {
	Space.Actor
	prop *Prop.RoleProp
	Item.ContainerOwner
	chatUser chatapi.IUser

	ms mailapi.IMailService
	//称号超时计时器
	titleChan  chan bool
	titleTimer *time.Timer

	//用于跑图调试 debug属性词缀
	debugAttraffix []attraffix.AttrAffix
}

func (r *_Role) Init() {
	r.prop = r.GetProp().(*Prop.RoleProp)

	r.InitItemContainer()

	if !r.prop.Data.InitedOnce {
		r.prop.SyncInitedOnce()
		r.initOnce()
		r.Debugf("[role init] first uid%v-name%v", r.prop.Data.Base.UID, r.prop.Data.Base.Name)
	} else {
		r.Debugf("[role init] not first uid%v-name%v", r.prop.Data.Base.UID, r.prop.Data.Base.Name)
	}

	//上线 重新计算所有build属性  因为build属性不存盘
	r.recalculationAllBuildAttr()

	r.eventRegist()
	r.mailOnline()

	//初始化计时器，上线即进行检查
	r.titleChan = make(chan bool)
	r.onTitleOverTime()

	r.Info("Role Init")
}

func (r *_Role) Destroy() {
	r.mailOffline()
	if r.titleTimer != nil {
		r.titleTimer.Stop()
	}
	r.Info("Role Destroy")
}

func (r *_Role) Loop() {
	select {
	case <-r.titleChan:
		r.onTitleOverTime()
	default:
	}
}

//initOnce 终身只执行一次
func (r *_Role) initOnce() {
	{
		//todo 临时代码 针对新账号默认解锁 一个特工+对应一个build
		specialAgentID := uint32(2)

		specialAgent := r.newSpecialAgent(specialAgentID)
		r.prop.SyncAddSpecialAgent(specialAgent)

		//默认解锁特工对应列表
		build := r.newBuild("预设 1", specialAgentID)
		r.prop.SyncAddBuild(build)
		r.prop.SyncSetFightingBuildID(build.BuildID)
		//当前build默认添加一把武器
		defaultEquipMent := uint32(60001)
		equipConfig, found := Data.GetEquipConfig().EquipMent_ConfigItems[defaultEquipMent]
		if found {
			item := r.CreateItem(defaultEquipMent, uint32(Proto.ItemEnum_Equipment), 1)
			r.AddItemToPack(item)
			r.BuildEquipItem(build.BuildID, int32(equipConfig.Position), item.GetID())
		}
		r.batchUnlockTalent()
		r.notifyTalentRedPoint()
	}

	{

		// 初始化每日任务数据
		team := r.GetSpace().(*_Team)
		team.doDailyResetJob()
	}

	r.FlushToDB()
	r.FlushToCache()
}

func (r *_Role) GetPropInfo() (string, string) {
	return Prop.RolePropType, r.GetID()
}

func (r *_Role) OnOnline() {
	t := time.Now().Unix()
	r.prop.SyncSetOnline(true, r.prop.Data.Base.LastLogoutTime, t)
	r.FlushToCache()
	r.initInvites()
}

func (r *_Role) OnOffline() {
	t := time.Now().Unix()
	r.prop.SyncSetOnline(false, t, r.prop.Data.Base.LastLoginTime)
	r.FlushToCache()
}

func (r *_Role) initInvites() {
	invites, err := DB.GetApply2InviteUtil().GetInvitesInRole(r.GetID())
	if err != nil {
		r.Error(err)
		return
	}

	if len(invites) == 0 {
		return
	}

	expireIDs := make([]primitive.ObjectID, 0)
	infos := make(map[string]*Proto.RoleInviteInfo)

	for i := 0; i < len(invites); i++ {
		invite := invites[i]
		if time.Now().Unix()-invite.Time >= Const.InviteMaxLife {
			invites = append(invites[:i], invites[i+1:]...)
			i--

			expireIDs = append(expireIDs, invite.ID)
		} else {
			infos[invite.TeamID] = &Proto.RoleInviteInfo{
				Instigator: invite.InviteParam.Instigator,
				Time:       invite.Time,
			}
		}
	}

	r.prop.Data.Invites = infos
	DB.GetApply2InviteUtil().RemoveByID(expireIDs)
}

func (r *_Role) GetGold() uint64 {
	return r.prop.Data.Gold
}

func (r *_Role) AddGold(num uint32, action uint32) bool {
	r.prop.SyncAddGold(num)
	r.FireLocalEvent(Const.Event_updateGold)
	r.Infof("[AddGold] %d,%d,%d", num, r.GetGold(), action)
	return true
}

func (r *_Role) RemoveGold(num uint32, action uint32) bool {
	if r.GetGold() < uint64(num) {
		return false
	}
	r.prop.SyncRemoveGold(num)
	r.FireLocalEvent(Const.Event_updateGold)
	r.Infof("[RemoveGold] %d,%d,%d", num, r.GetGold(), action)
	return true
}

func (r *_Role) GetDiamond() uint64 {
	return r.prop.Data.Diamond
}

func (r *_Role) AddDiamond(num uint32, action uint32) bool {
	r.prop.SyncAddDiamond(num)
	r.FireLocalEvent(Const.Event_updateDiamond)
	r.Infof("[AddDiamond] %d,%d,%d", num, r.GetDiamond(), action)
	return true
}

func (r *_Role) RemoveDiamond(num uint32, action uint32) bool {
	if r.GetDiamond() < uint64(num) {
		return false
	}
	r.prop.SyncRemoveDiamond(num)
	r.FireLocalEvent(Const.Event_updateDiamond)
	r.Infof("[RemoveDiamond] %d,%d,%d", num, r.GetDiamond(), action)
	return true
}

//addFightingSpecialAgentExp 增加出战特工经验
func (r *_Role) addFightingSpecialAgentExp(exp uint64) {
	buildData, ok := r.prop.Data.BuildMap[r.prop.Data.FightingBuildID]
	if !ok {
		return
	}

	r.addSpecialAgentExp(buildData.SpecialAgentID, exp)
}

// addSpecialAgentExp 特工增加经验
func (r *_Role) addSpecialAgentExp(specialAgentID uint32, exp uint64) {
	specialAgent, ok := r.prop.Data.SpecialAgentList[specialAgentID]
	if !ok {
		r.Error("没有对应特工 ", specialAgentID)
		return
	}

	//获取当前特工下一个等级的配置
	_, ok = Data.GetSpecialAgentConfig().Upgrade_ConfigItems[specialAgent.Base.Level+1]
	if !ok {
		//已经达到最大等级 不在获取经验值
		r.Debug("特工已经达到满级  ", specialAgent.Base.Level)
		return
	}

	specialAgent.Base.Exp += exp
	r.calcSpecialAgentLv(specialAgent)
}

//expBonus 经验红利加成
func (r *_Role) expBonus(exp uint64) uint64 {
	//只要拥有满级特工 那么获得经验享受加成
	if r.isSpecialAgentMaxlevel() {
		var coefficient uint64 = 100 //加成系数 value/100 便于取整
		var baseExp uint64
		coefficientCfg, ok := Data.GetSpecialAgentConfig().SpecialAgentConst_ConfigItems[Const.SpecialAgent_upgradeExpCoefficient]
		if ok {
			coefficient = uint64(coefficientCfg.Value)
		}
		baseCfg, ok := Data.GetSpecialAgentConfig().SpecialAgentConst_ConfigItems[Const.SpecialAgent_upgradeExpBase]
		if ok {
			baseExp = uint64(baseCfg.Value)
		}

		exp = exp*coefficient/100 + baseExp
	}

	return exp
}

// isSpecialAgentMaxlevel 是否已经有满级特工
func (r *_Role) isSpecialAgentMaxlevel() bool {
	for _, specialAgent := range r.prop.Data.SpecialAgentList {
		_, ok := Data.GetSpecialAgentConfig().Upgrade_ConfigItems[specialAgent.Base.Level+1]
		if !ok {
			return true
		}
	}

	return false
}

// calcSpecialAgentLv 计算特工等级
func (r *_Role) calcSpecialAgentLv(specialAgent *Proto.SpecialAgent) {
	if specialAgent == nil {
		return
	}

	// 是否升级
	leveluped := false
	beforeLevel := specialAgent.Base.Level
	for i := specialAgent.Base.Level; i <= uint32(len(Data.GetSpecialAgentConfig().Upgrade_ConfigItems)); i++ {
		//获取当前特工下一个等级的配置
		nextCfg, ok := Data.GetSpecialAgentConfig().Upgrade_ConfigItems[specialAgent.Base.Level+1]
		if !ok {
			break
		}

		//经验不满足升级需求 退出循环
		if specialAgent.Base.Exp < uint64(nextCfg.Exp) {
			break
		}

		leveluped = true
		specialAgent.Base.Level++
		specialAgent.Base.Exp -= uint64(nextCfg.Exp)
	}

	r.prop.SyncUpdateSpecialAgentLv(specialAgent.Base.ConfigID, specialAgent.Base.Level, specialAgent.Base.Exp)
	r.recalculationAllBuildAttr()

	if leveluped {
		// 发送升级事件
		r.FireLocalEvent(Const.Event_specialAgentLevelUp, beforeLevel, specialAgent.Base.ConfigID)
	}
}

func (r *_Role) GetMoney(typ uint32) uint64 {
	switch typ {
	case Const.Gold:
		return r.GetGold()
	case Const.Diamond:
		return r.GetDiamond()
	}
	return 0
}

func (r *_Role) AddMoney(num, typ, action uint32) bool {
	switch typ {
	case Const.Gold:
		return r.AddGold(num, action)
	case Const.Diamond:
		return r.AddDiamond(num, action)
	}
	return false
}

func (r *_Role) RemoveMoney(num, typ, action uint32) bool {
	switch typ {
	case Const.Gold:
		return r.RemoveGold(num, action)
	case Const.Diamond:
		return r.RemoveDiamond(num, action)
	}
	return false
}

//eventRegist 事件注册
func (r *_Role) eventRegist() {
	r.RegLocalEvent(Const.Event_addItem, r.onItemAdded)
	r.RegLocalEvent(Const.Event_removeItem, r.onItemRemoved)
	r.RegLocalEvent(Const.Event_updateItem, r.onItemUpdated)
	r.RegLocalEvent(Const.Event_updateGold, r.onGoldChanged)
	r.RegLocalEvent(Const.Event_updateDiamond, r.onDiamondChanged)
	r.RegLocalEvent(Const.Event_specialAgentLevelUp, r.onSpecialAgentLevelUp)
	r.RegLocalEvent(Const.Event_commanderLevelUp, r.onCommanderLevelUp)
}

//onItemAdded 获取道具事件处理
func (r *_Role) onItemAdded(items []ItemProto.IItem) {
	for _, iItem := range items {
		item := iItem.GetData()
		if item == nil {
			continue
		}

		switch item.Base.Type {
		case Proto.ItemEnum_SkillItem:
			r.onSkillItemAdded(item)
		default:

		}
	}
}

//onItemRemoved 消耗道具事件处理
func (r *_Role) onItemRemoved(items []ItemProto.IItem) {
	for _, iItem := range items {
		item := iItem.GetData()
		if item == nil {
			continue
		}

		switch item.Base.Type {
		case Proto.ItemEnum_SkillItem:
			r.onSkillItemRemoved(item)
		default:

		}
	}
}

//onItemUpdated 更新道具事件处理  例如 更新技能等级  或者更新道具数量 等等
func (r *_Role) onItemUpdated(items []ItemProto.IItem) {
	for _, iItem := range items {
		item := iItem.GetData()
		if item == nil {
			continue
		}

		switch item.Base.Type {
		case Proto.ItemEnum_SkillItem:
			r.onSkillItemUpdated(item)
		default:

		}
	}
}

//onGoldChanged 更新金币事件处理
func (r *_Role) onGoldChanged() {
	r.checkAllLearnedSkillRedPoint()
}

//onDiamondChanged 更新钻石事件处理
func (r *_Role) onDiamondChanged() {
}

// onSpecialAgentLevelUp 特工升级情况处理
func (r *_Role) onSpecialAgentLevelUp(beforeLevel uint32, sid uint32) {
	r.addTalentPointByLevelUp(beforeLevel, sid)
}

// onCommanderLevelUp 指挥官升级事件处理
func (r *_Role) onCommanderLevelUp(beforeLevel uint32) {
}

//addCommanderExp 增加指挥官经验
func (r *_Role) addCommanderExp(exp uint32) {
	r.prop.Data.Base.Exp += uint64(exp)
	r.calcCommanderLv()
}

//计算指挥官等级
func (r *_Role) calcCommanderLv() {
	// 是否升级
	leveluped := false
	beforeLevel := r.prop.Data.Base.Level
	for i := beforeLevel; i <= uint32(len(Data.GetPlayerUpgradeConfig().PlayerUpgrade_ConfigItems)); i++ {
		//获取下一个等级的配置
		nextCfg, ok := Data.GetPlayerUpgradeConfig().PlayerUpgrade_ConfigItems[r.prop.Data.Base.Level+1]
		if !ok {
			break
		}

		//经验不满足升级需求 退出循环
		if r.prop.Data.Base.Exp < uint64(nextCfg.Exp) {
			break
		}

		leveluped = true
		r.prop.Data.Base.Level++
		r.prop.Data.Base.Exp -= uint64(nextCfg.Exp)
	}

	r.prop.SyncUpdateCommanderLv(r.prop.Data.Base.Level, r.prop.Data.Base.Exp)

	if leveluped {
		r.FireLocalEvent(Const.Event_commanderLevelUp, beforeLevel)
	}
}
