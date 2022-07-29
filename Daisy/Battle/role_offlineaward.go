package main

import (
	"Daisy/Data"
	"Daisy/ItemProto"
	"Daisy/Prop"
	"Daisy/Proto"
	"time"
)

type Spoil struct {
	MaterialDrop   []*Proto.DropMaterial
	MoneyDrop      uint32
	ActorExpDrop   uint32
	SpecialExpdrop uint32
}

// RPC_ReqOfflineAward 客户端请求得到离线收益
func (user *_User) RPC_ReqOfflineAward() *Proto.RspOfflineAwardData {

	if user.role.prop.Data.Base.LastLogoutTime == 0 {
		user.role.Error("[RPC_ReqOfflineAward] 新账号，没有离线收益")
		return &Proto.RspOfflineAwardData{}
	}

	ret := &Proto.RspOfflineAwardData{
		OfflineTime:       time.Now().Unix() - user.role.prop.Data.Base.LastLogoutTime,
		OfflineAwardDatas: &Proto.OfflineAwardData{},
	}
	*ret.OfflineAwardDatas = *user.role.prop.Data.OfflineAwardDatas

	user.role.prop.SyncResetOfflineAward()
	return ret
}

func (r *_Role) CalcAwardByDuration(checkpoint uint32, offlineTime time.Duration, multi float32) (*Proto.OfflineAwardData, error) {
	// 读取关卡表，计算该时间段内经验、金币、物品掉落几次。
	tmp := Data.GetSceneConfig().BattleArea_ConfigItems[checkpoint]
	offlineAddFreq := tmp.GetOfflineDropFreq()                            // 表格配置每分钟掉落多少次
	offlineAdd := offlineTime.Minutes() * offlineAddFreq * float64(multi) // 当前离线时间段内掉落多少次， float64
	itemTime := uint32(offlineAdd)                                        // 当前离线时间段内掉落多少次， uint32

	// 根据关卡表得到掉落的战利品
	s := r.getAwardFromNpcList(checkpoint, itemTime)
	// 将战利品转化为离线收益， 并将掉落道具放入背包
	return r.fillUpOfflineAwardData(s, offlineAdd), nil
}

// addOfflineAward 计算队伍挂机奖励 （关卡， 队伍上次销毁时间，队伍本次初始化时间）
func (r *_Role) addOfflineAward(checkpoint uint32, offlineTime time.Duration) {
	f, err := r.CalcAwardByDuration(checkpoint, offlineTime, 1)
	if err != nil {
		r.Error(err)
		return
	}

	// 属性刷新
	r.prop.SyncAddGold(f.Money)
	r.addCommanderExp(f.ActorExp)
	r.addFightingSpecialAgentExp(r.expBonus(uint64(f.SpecialAgentExp)))
	r.prop.SyncAddOfflineAward(f)
}

// fillUpOfflineAwardData 填充离线收益结构体
func (r *_Role) fillUpOfflineAwardData(mySoil Spoil, offlineAddMultiple float64) *Proto.OfflineAwardData {
	award := &Proto.OfflineAwardData{}
	award.Money = uint32(float64(mySoil.MoneyDrop) * offlineAddMultiple) // 掉落经验金币 = 每次掉落数量 * 掉落次数
	award.ActorExp = uint32(float64(mySoil.ActorExpDrop) * offlineAddMultiple)
	award.SpecialAgentExp = uint32(float64(mySoil.SpecialExpdrop) * offlineAddMultiple)

	team := r.GetSpace().(*_Team)
	// 创建队伍成员信息列表
	var memberInfoList []*Proto.OwnerInfo
	for id, _ := range team.prop.Data.Base.Members {
		owner := &Proto.OwnerInfo{id, team.GetActorByUserID(id).GetProp().(*Prop.RoleProp).Data.Base.Name}
		memberInfoList = append(memberInfoList, owner)
	}

	itemData := r.AddItemList(mySoil.MaterialDrop)
	for _, val := range itemData {
		// 表示加入背包失败
		if val.GetNum == 0 {
			continue
		}
		awardItem := &Proto.OfflineAwardItem{
			ID:       val.ItemData.Base.ID,
			Type:     val.ItemData.Base.Type,
			Num:      val.GetNum,
			ConfigID: val.ItemData.Base.ConfigID,
		}
		award.OfflineAwardItems = append(award.OfflineAwardItems, awardItem)

		if val.ItemData.Base.Type == Proto.ItemEnum_Equipment {
			val.ItemData.EquipmentData.OwnerTeamMemberList = memberInfoList
			r.UpdataItem(ItemProto.CreateIItemByData(val.ItemData))
		}
	}
	return award
}

// getAwardFromNpcList 根据怪物列表得到奖励保存在spoil结构中
func (r *_Role) getAwardFromNpcList(checkpoint uint32, DropTimeMultiple uint32) Spoil {
	s := Spoil{}

	cfg, ok := Data.GetSceneConfig().BattleArea_ConfigItems[checkpoint]
	if !ok {
		r.Error("[getAwardFromNpcList] teamHangUpAward failed, 找不到对应关卡表，掉落失败", checkpoint)
		return s
	}
	s.MoneyDrop += cfg.DropMoney
	s.ActorExpDrop += cfg.ActorDropExp
	s.SpecialExpdrop += cfg.SpecialDropExp

	var j uint32
	for j = 0; j != DropTimeMultiple; j++ {
		items, err := r.Drop(cfg.DropID)
		if err != nil {
			r.Error("生成掉落品失败 ", cfg.DropID)
			return s
		}

		for _, val := range items {
			s.MaterialDrop = append(s.MaterialDrop, val)
		}
	}

	return s
}
