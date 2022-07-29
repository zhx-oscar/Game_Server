package Prop

import (
	"Cinder/Base/Message"
	"Cinder/Base/Prop"
	"Daisy/Proto"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

type RoleProp struct {
	Data *Proto.Role

	Prop.Prop
}

// todo 填充实现

func (u *RoleProp) Marshal() ([]byte, error) {
	return u.Data.Marshal()
}

func (u *RoleProp) UnMarshal(data []byte) error {
	u.Data = &Proto.Role{}

	if err := u.Data.Unmarshal(data); err != nil {
		return err
	}

	u.fillDefault()

	return nil
}

func (u *RoleProp) MarshalCache() ([]byte, error) {
	cache := &Proto.RoleCache{
		Base:            u.Data.Base,
		BuildMap:        u.Data.BuildMap,
		FightingBuildID: u.Data.FightingBuildID,
	}

	return cache.Marshal()
}

func (u *RoleProp) UnMarshalCache(data []byte) (interface{}, error) {
	cache := &Proto.RoleCache{}
	if err := cache.Unmarshal(data); err != nil {
		return nil, err
	}

	return cache, nil
}

func (u *RoleProp) MarshalPart() ([]byte, error) {
	part := &Proto.Role{
		Base:             u.Data.Base,
		SpecialAgentList: u.Data.SpecialAgentList,
		BuildMap:         u.Data.BuildMap,
		FightingBuildID:  u.Data.FightingBuildID,
		SkillsLearned:    u.Data.SkillsLearned,
		ItemContainerMap: map[int32]*Proto.ItemContainer{
			int32(Proto.ContainerEnum_BuildEquipBag): u.Data.ItemContainerMap[int32(Proto.ContainerEnum_BuildEquipBag)],
		},
		InheritAttr: u.Data.InheritAttr,
		ShareSpoils: &Proto.ShareSpoils{
			RequestSkillData: u.Data.ShareSpoils.RequestSkillData,
		},
		FastBattle: u.Data.FastBattle,
	}

	return part.Marshal()
}

func (u *RoleProp) UnMarshalPart(data []byte) error {
	return u.UnMarshal(data)
}

func (u *RoleProp) MarshalToBson() ([]byte, error) {
	return bson.Marshal(u.Data)
}

func (u *RoleProp) UnMarshalFromBson(data []byte) error {
	u.Data = &Proto.Role{}

	if len(data) != 0 {
		if err := bson.Unmarshal(data, u.Data); err != nil {
			return err
		}
	}

	u.fillDefault()

	return nil
}

func (u *RoleProp) fillDefault() {
	if u.Data.Base == nil {
		u.Data.Base = &Proto.RoleBase{
			RaidProgress: 1,
		}
	}

	if u.Data.SpecialAgentList == nil {
		u.Data.SpecialAgentList = map[uint32]*Proto.SpecialAgent{}
	}

	if u.Data.BuildMap == nil {
		u.Data.BuildMap = map[string]*Proto.BuildData{}
	}

	for _, build := range u.Data.BuildMap {
		if build.Skill == nil {
			build.Skill = &Proto.BuildSkillData{
				UltimateSkillID: 0,
				SuperSkill:      map[uint32]uint32{},
			}
		} else {
			if build.Skill.SuperSkill == nil {
				build.Skill.SuperSkill = map[uint32]uint32{}
			}
		}
	}

	if u.Data.ItemContainerMap == nil {
		u.Data.ItemContainerMap = make(map[int32]*Proto.ItemContainer)
	}

	for _, bag := range u.Data.ItemContainerMap {
		for _, item := range bag.ItemMap {
			if item.ExpandData == nil {
				item.ExpandData = &Proto.ItemExpand{
					InUse:             0,
					RelationBuildList: map[string]bool{},
				}
			}

			if item.ExpandData.RelationBuildList == nil {
				item.ExpandData.RelationBuildList = map[string]bool{}
			}
		}
	}

	if u.Data.Invites == nil {
		u.Data.Invites = make(map[string]*Proto.RoleInviteInfo)
	}

	if u.Data.Chat == nil {
		u.Data.Chat = &Proto.RoleChat{}
	}

	if u.Data.Chat.ChatChannels == nil {
		u.Data.Chat.ChatChannels = map[string]bool{}
	}

	if u.Data.InheritAttr == nil {
		u.Data.InheritAttr = &Proto.PawnInherit{}
	}

	if u.Data.OfflineAwardDatas == nil {
		u.Data.OfflineAwardDatas = &Proto.OfflineAwardData{
			OfflineAwardItems: make([]*Proto.OfflineAwardItem, 0),
		}
	}

	if u.Data.SkillsLearned == nil {
		u.Data.SkillsLearned = map[uint32]*Proto.SkillData{}
	}

	if u.Data.RedPointsData == nil {
		u.Data.RedPointsData = map[string]*Proto.RedPointInfo{}
	}
	if u.Data.Title == nil {
		u.Data.Title = &Proto.Title{}
	}
	if u.Data.Title.TitleList == nil {
		u.Data.Title.TitleList = map[uint32]*Proto.TitleInfo{}
	}

	if u.Data.MailBox == nil {
		u.Data.MailBox = &Proto.MailBox{
			Mails: make([]*Proto.Mail, 0),
		}
	}
	if u.Data.ShareSpoils == nil {
		u.Data.ShareSpoils = &Proto.ShareSpoils{
			RequestSkillData:       map[string]*Proto.RequestSkill{},
			GiveSkillToTargetCount: map[string]*Proto.GiveSkillCount{},
		}
	}

	if u.Data.ShareSpoils.RequestSkillData == nil {
		u.Data.ShareSpoils.RequestSkillData = map[string]*Proto.RequestSkill{}
	}

	for _, requestSkill := range u.Data.ShareSpoils.RequestSkillData {
		if requestSkill.TeamMateGiveSkillCount == nil {
			requestSkill.TeamMateGiveSkillCount = map[string]uint32{}
		}
	}

	if u.Data.ShareSpoils.GiveSkillToTargetCount == nil {
		u.Data.ShareSpoils.GiveSkillToTargetCount = map[string]*Proto.GiveSkillCount{}
	}

	for _, giveSkillCount := range u.Data.ShareSpoils.GiveSkillToTargetCount {
		if giveSkillCount.GiveSkillToTargetCount == nil {
			giveSkillCount.GiveSkillToTargetCount = map[uint32]uint32{}
		}
	}

	if u.Data.SupplyInfo == nil {
		u.Data.SupplyInfo = &Proto.SupplyInfo{
			BoxMold:     map[uint32]*Proto.State{},
			DiscountNum: map[uint32]uint32{},
		}
	}

	if u.Data.SupplyInfo.BoxMold == nil {
		u.Data.SupplyInfo.BoxMold = map[uint32]*Proto.State{}
	}

	if u.Data.SupplyInfo.DiscountNum == nil {
		u.Data.SupplyInfo.DiscountNum = map[uint32]uint32{}
	}

	if u.Data.Friends == nil {
		u.Data.Friends = &Proto.FriendListData{
			ApplyList:       make([]*Proto.Friend, 0),
			FriendList:      map[string]*Proto.Friend{},
			RecommendedList: make([]*Proto.Friend, 0),
		}
	}
	if u.Data.Friends.FriendList == nil {
		u.Data.Friends.FriendList = make(map[string]*Proto.Friend)
	}
	if u.Data.Friends.ApplyList == nil {
		u.Data.Friends.ApplyList = make([]*Proto.Friend, 0)
	}
	if u.Data.Friends.RecommendedList == nil {
		u.Data.Friends.RecommendedList = make([]*Proto.Friend, 0)
	}
	if u.Data.SeasonInfo == nil {
		u.Data.SeasonInfo = &Proto.RoleSeasonInfo{}
	}
	if u.Data.SeasonInfo.LevelAwards == nil {
		u.Data.SeasonInfo.LevelAwards = map[uint32]bool{}
	}
}

func (u *RoleProp) SyncInitedOnce() {
	u.Sync("InitedOnce", Message.PackArgs(), true)
}

func (u *RoleProp) InitedOnce() {
	u.Data.InitedOnce = true
}

func (u *RoleProp) SyncSetName(name string, head uint32) {
	u.Sync("SetName", Message.PackArgs(name, head), true)
}

func (u *RoleProp) SetName(name string, head uint32) {
	u.Data.Base.Name = name
	u.Data.Base.Head = head
}

func (u *RoleProp) SyncSetTeamID(teamID string) {
	u.Sync("SetTeamID", Message.PackArgs(teamID), true, Prop.Target_Client)
}

func (u *RoleProp) SetTeamID(teamID string) {
	u.Data.Base.TeamID = teamID
}

func (u *RoleProp) SyncSetOnline(online bool, lastLogoutTime, lastLoginTime int64) {
	u.Sync("SetOnline", Message.PackArgs(online, lastLogoutTime, lastLoginTime), true, Prop.Target_Other_Clients)
}

func (u *RoleProp) SetOnline(online bool, lastLogoutTime, lastLoginTime int64) {
	u.Data.Base.Online = online
	u.Data.Base.LastLogoutTime = lastLogoutTime
	u.Data.Base.LastLoginTime = lastLoginTime
}

//SyncAddSpecialAgent 同步增加特工
func (u *RoleProp) SyncAddSpecialAgent(specialAgent *Proto.SpecialAgent) {
	u.Sync("AddSpecialAgent", Message.PackArgs(specialAgent), true, Prop.Target_All_Clients)
}

//AddSpecialAgent 增加特工
func (u *RoleProp) AddSpecialAgent(specialAgent *Proto.SpecialAgent) {
	u.Data.SpecialAgentList[specialAgent.Base.ConfigID] = specialAgent
}

//SyncUpdateSpecialAgentLv 同步更新特工等级
func (u *RoleProp) SyncUpdateSpecialAgentLv(id, lv uint32, exp uint64) {
	u.Sync("UpdateSpecialAgentLv", Message.PackArgs(id, lv, exp), true, Prop.Target_All_Clients)
}

//UpdateSpecialAgentLv 更新特工等级数据
func (u *RoleProp) UpdateSpecialAgentLv(id, lv uint32, exp uint64) {
	u.Data.SpecialAgentList[id].Base.Level = lv
	u.Data.SpecialAgentList[id].Base.Exp = exp
}

func (u *RoleProp) SyncSetFightingBuildID(buildID string) {
	u.Sync("SetFightingBuildID", Message.PackArgs(buildID), true, Prop.Target_All_Clients)
}

func (u *RoleProp) SetFightingBuildID(buildID string) {
	u.Data.FightingBuildID = buildID
}

func (u *RoleProp) SyncAddRoleInviteInfo(teamID, instigator string) {
	info := &Proto.RoleInviteInfo{
		Instigator: instigator,
		Time:       time.Now().Unix(),
	}
	u.Sync("AddRoleInviteInfo", Message.PackArgs(teamID, info), false, Prop.Target_Client)
}
func (u *RoleProp) AddRoleInviteInfo(teamID string, info *Proto.RoleInviteInfo) {
	u.Data.Invites[teamID] = info
}

func (u *RoleProp) SyncRemoveRoleInviteInfo(teamID string) {
	u.Sync("RemoveRoleInviteInfo", Message.PackArgs(teamID), false, Prop.Target_Client)
}
func (u *RoleProp) RemoveRoleInviteInfo(teamID string) {
	delete(u.Data.Invites, teamID)
}

func (u *RoleProp) SyncClearRoleInviteInfo() {
	u.Sync("ClearRoleInviteInfo", Message.PackArgs(), false, Prop.Target_Client)
}
func (u *RoleProp) ClearRoleInviteInfo() {
	u.Data.Invites = make(map[string]*Proto.RoleInviteInfo)
}

func (u *RoleProp) GetContainerData(typ int32) (*Proto.ItemContainer, error) {
	data, ok := u.Data.ItemContainerMap[typ]
	if ok == false {
		return nil, fmt.Errorf("not found")
	}
	return data, nil
}

func (u *RoleProp) AddContainerData(typ int32, maxNum uint32) *Proto.ItemContainer {
	data, ok := u.Data.ItemContainerMap[typ]
	if ok == true {
		return data
	}
	data = &Proto.ItemContainer{}
	data.ItemMap = make(map[int32]*Proto.Item)
	data.Id2Pos = make(map[string]int32)
	data.Type = Proto.ContainerEnum_Type(typ)
	data.MaxNum = maxNum
	u.SyncAddContainer(typ, data)

	data2, _ := u.Data.ItemContainerMap[typ]
	return data2
}

func (u *RoleProp) SyncAddContainer(typ int32, data *Proto.ItemContainer) {
	u.Sync("SetAddContainer", Message.PackArgs(typ, data), true, Prop.Target_Client)
}
func (u *RoleProp) SetAddContainer(typ int32, data *Proto.ItemContainer) {
	u.Data.ItemContainerMap[typ] = data
}

func (u *RoleProp) SyncAddItemToContainer(typ int32, items *Proto.Items, target int) {
	u.Sync("SetAddItemToContainer", Message.PackArgs(typ, items), true, target)
}
func (u *RoleProp) SetAddItemToContainer(typ int32, items *Proto.Items) {
	data, err := u.GetContainerData(typ)
	if err == nil {
		if data.Id2Pos == nil {
			data.Id2Pos = make(map[string]int32)
		}
		if data.ItemMap == nil {
			data.ItemMap = make(map[int32]*Proto.Item)
		}

		for _, item := range items.Items {
			if item.ExpandData.RelationBuildList == nil {
				item.ExpandData.RelationBuildList = map[string]bool{}
			}

			data.Id2Pos[item.Base.ID] = item.Base.Pos
			data.ItemMap[item.Base.Pos] = item
		}
	}
}

func (u *RoleProp) SyncDelItemFromContainer(typ int32, Pos *Proto.Int32Array, target int) {
	u.Sync("SetDelItemFromContainer", Message.PackArgs(typ, Pos), true, target)
}
func (u *RoleProp) SetDelItemFromContainer(typ int32, Pos *Proto.Int32Array) {
	data, err := u.GetContainerData(typ)
	if err == nil {
		for _, v := range Pos.Data {
			item, ok := data.ItemMap[v]
			if ok == true {
				delete(data.ItemMap, v)
				delete(data.Id2Pos, item.Base.ID)
			}
		}
	}
}

func (u *RoleProp) SyncUpdateItemToContainer(typ int32, items *Proto.Items, target int) {
	u.Sync("SetUpdateItemToContainer", Message.PackArgs(typ, items), true, target)
}
func (u *RoleProp) SetUpdateItemToContainer(typ int32, items *Proto.Items) {
	data, err := u.GetContainerData(typ)
	if err == nil {
		for _, item := range items.Items {
			if item.ExpandData.RelationBuildList == nil {
				item.ExpandData.RelationBuildList = map[string]bool{}
			}

			data.ItemMap[item.Base.Pos] = item
		}
	}
}

//SyncAddBuild 同步增加build
func (u *RoleProp) SyncAddBuild(build *Proto.BuildData) {
	u.Sync("AddBuild", Message.PackArgs(build), true, Prop.Target_All_Clients)
}

//AddBuild 增加build
func (u *RoleProp) AddBuild(build *Proto.BuildData) {
	if build.Skill == nil {
		build.Skill = &Proto.BuildSkillData{
			UltimateSkillID: 0,
			SuperSkill:      map[uint32]uint32{},
		}
	}

	if build.Skill.SuperSkill == nil {
		build.Skill.SuperSkill = map[uint32]uint32{}
	}

	u.Data.BuildMap[build.BuildID] = build
}

//SyncSetBuildName 同步设置buildName
func (u *RoleProp) SyncSetBuildName(buildID, name string) {
	u.Sync("SetBuildName", Message.PackArgs(buildID, name), true, Prop.Target_All_Clients)
}

//SetBuildName 设置buildName
func (u *RoleProp) SetBuildName(buildID, name string) {
	u.Data.BuildMap[buildID].Name = name
}

//SyncUpdateBuildUltimateSkillID 同步设置build内必杀技
func (u *RoleProp) SyncUpdateBuildUltimateSkillID(buildID string, skillID uint32) {
	u.Sync("UpdateBuildUltimateSkillID", Message.PackArgs(buildID, skillID), true, Prop.Target_All_Clients)
}

//UpdateBuildUltimateSkillID 设置build内必杀技
func (u *RoleProp) UpdateBuildUltimateSkillID(buildID string, skillID uint32) {
	u.Data.BuildMap[buildID].Skill.UltimateSkillID = skillID
}

//SyncSwapBuildSuper 同步交换build内超能技能
func (u *RoleProp) SyncSwapBuildSuper(buildID string, aSkillID, aPos, bSkillID, bPos uint32) {
	u.Sync("SwapBuildSuper", Message.PackArgs(buildID, aSkillID, aPos, bSkillID, bPos), true, Prop.Target_All_Clients)
}

//SwapBuildSuper 交换build内超能技能
func (u *RoleProp) SwapBuildSuper(buildID string, aSkillID, aPos, bSkillID, bPos uint32) {
	//A技能 当前build已经使用
	if aPos != 0 {
		//目标槽位已经有B技能 则 ab两技能交换位置
		if bSkillID != 0 {
			u.Data.BuildMap[buildID].Skill.SuperSkill[bPos] = aSkillID
			u.Data.BuildMap[buildID].Skill.SuperSkill[aPos] = bSkillID
		} else {
			//目标槽位为空，则A技能挪移到目标槽位 原来A的槽位为空
			u.Data.BuildMap[buildID].Skill.SuperSkill[bPos] = aSkillID
			delete(u.Data.BuildMap[buildID].Skill.SuperSkill, aPos)
		}
	} else {
		//A技能 当前build未使用 则直接覆盖到目标槽位
		u.Data.BuildMap[buildID].Skill.SuperSkill[bPos] = aSkillID
	}
}

//SyncBuildEquipItem 同步build装备
func (u *RoleProp) SyncBuildEquipItem(buildID string, pos int32, itemID string) {
	u.Sync("SetBuildEquipItem", Message.PackArgs(buildID, pos, itemID), true, Prop.Target_All_Clients)
}

func (u *RoleProp) SetBuildEquipItem(buildID string, pos int32, itemID string) {
	if u.Data.BuildMap[buildID].EquipmentMap == nil {
		u.Data.BuildMap[buildID].EquipmentMap = make(map[int32]string)
	}
	u.Data.BuildMap[buildID].EquipmentMap[pos] = itemID
}

// SyncSetBuildFightAttr 同步设置build面板显示属性
func (u *RoleProp) SyncSetBuildFightAttr(buildID string, fightAttr *Proto.FightAttr) {
	u.Sync("SetBuildFightAttr", Message.PackArgs(buildID, fightAttr), true, Prop.Target_All_Clients)
}

// SetBuildFightAttr 设置build面板显示属性
func (u *RoleProp) SetBuildFightAttr(buildID string, fightAttr *Proto.FightAttr) {
	if fightAttr == nil {
		fightAttr = &Proto.FightAttr{}
	}
	u.Data.BuildMap[buildID].FightAttr = fightAttr
}

// SyncAddToChatChannel 同步添加聊天频道列表
func (u *RoleProp) SyncAddToChatChannel(groupID string) {
	u.Sync("SetAddToChatChannel", Message.PackArgs(groupID), true)
}

// SetAddToChatChannel 添加到聊天频道列表
func (u *RoleProp) SetAddToChatChannel(groupID string) {
	u.Data.Chat.ChatChannels[groupID] = true
}

// SyncDelFromChatChannel 同步移出聊天频道列表
func (u *RoleProp) SyncDelFromChatChannel(groupID string) {
	u.Sync("SetDelFromChatChannel", Message.PackArgs(groupID), true)
}

// SetDelFromChatChannel 移出聊天频道列表
func (u *RoleProp) SetDelFromChatChannel(groupID string) {
	delete(u.Data.Chat.ChatChannels, groupID)
}

func (u *RoleProp) SyncSetInheritAttr(attr *Proto.PawnInherit) {
	u.Sync("SetInheritAttr", Message.PackArgs(attr), true, Prop.Target_All_Clients)
}
func (u *RoleProp) SetInheritAttr(attr *Proto.PawnInherit) {
	u.Data.InheritAttr = attr
}

func (u *RoleProp) SyncAddGold(num uint32) {
	u.Sync("AddGold", Message.PackArgs(num), true, Prop.Target_Client)
}

func (u *RoleProp) AddGold(num uint32) {
	u.Data.Gold += uint64(num)
}

func (u *RoleProp) SyncRemoveGold(num uint32) {
	u.Sync("RemoveGold", Message.PackArgs(num), true, Prop.Target_Client)
}

func (u *RoleProp) RemoveGold(num uint32) {
	if u.Data.Gold >= uint64(num) {
		u.Data.Gold -= uint64(num)
	} else {
		u.Data.Gold = 0
	}
}

func (u *RoleProp) SyncAddDiamond(num uint32) {
	u.Sync("AddDiamond", Message.PackArgs(num), true, Prop.Target_Client)
}

func (u *RoleProp) AddDiamond(num uint32) {
	u.Data.Diamond += uint64(num)
}

func (u *RoleProp) SyncRemoveDiamond(num uint32) {
	u.Sync("RemoveDiamond", Message.PackArgs(num), true, Prop.Target_Client)
}

func (u *RoleProp) RemoveDiamond(num uint32) {
	if u.Data.Diamond >= uint64(num) {
		u.Data.Diamond -= uint64(num)
	} else {
		u.Data.Diamond = 0
	}
}

// SyncAddOfflineAward 同步刷新离线收益缓存
func (u *RoleProp) SyncAddOfflineAward(award *Proto.OfflineAwardData) {
	u.Sync("AddOfflineAward", Message.PackArgs(award), true)
}

// AddOfflineAward 刷新离线收益缓存
func (u *RoleProp) AddOfflineAward(award *Proto.OfflineAwardData) {
	u.Data.OfflineAwardDatas.ActorExp += award.ActorExp
	u.Data.OfflineAwardDatas.SpecialAgentExp += award.SpecialAgentExp
	u.Data.OfflineAwardDatas.Money += award.Money
	for _, val := range award.OfflineAwardItems {
		u.Data.OfflineAwardDatas.OfflineAwardItems = append(u.Data.OfflineAwardDatas.OfflineAwardItems, val)
	}
}

// SyncResetOfflineAward 同步重置离线收益
func (u *RoleProp) SyncResetOfflineAward() {
	u.Sync("ResetOfflineAward", Message.PackArgs(), true)
}

// ResetOfflineAward 重置离线收益
func (u *RoleProp) ResetOfflineAward() {
	u.Data.OfflineAwardDatas.ActorExp = 0
	u.Data.OfflineAwardDatas.SpecialAgentExp = 0
	u.Data.OfflineAwardDatas.Money = 0
	u.Data.OfflineAwardDatas.OfflineAwardItems = nil
}

func (u *RoleProp) SyncRemoveExp(num uint32) {
	u.Sync("RemoveExp", Message.PackArgs(num), true, Prop.Target_Client)
}

func (u *RoleProp) RemoveExp(num uint32) {
	if u.Data.Base.Exp >= uint64(num) {
		u.Data.Base.Exp -= uint64(num)
	} else {
		u.Data.Base.Exp = 0
	}
}

//SyncUpdateAlreadyLearnedSkill 同步更新已经学习的技能
func (u *RoleProp) SyncUpdateAlreadyLearnedSkill(skillData *Proto.SkillData) {
	u.Sync("UpdateAlreadyLearnedSkill", Message.PackArgs(skillData), true, Prop.Target_All_Clients)
}

//UpdateAlreadyLearnedSkill 更新已经学习的技能
func (u *RoleProp) UpdateAlreadyLearnedSkill(skillData *Proto.SkillData) {
	if skillData == nil {
		return
	}

	u.Data.SkillsLearned[skillData.SkillMainID] = skillData
}

//SyncAddRedPoint 同步增加小红点数据
func (u *RoleProp) SyncAddRedPoint(redPointKey string, data *Proto.RedPointInfo) {
	u.Sync("AddRedPoint", Message.PackArgs(redPointKey, data), true, Prop.Target_Client)
}

//AddRedPoint 增加小红点
func (u *RoleProp) AddRedPoint(redPointKey string, data *Proto.RedPointInfo) {
	u.Data.RedPointsData[redPointKey] = data
}

//SyncRemoveRedPoint 同步移除小红点数据
func (u *RoleProp) SyncRemoveRedPoint(redPointKey string) {
	u.Sync("RemoveRedPoint", Message.PackArgs(redPointKey), true, Prop.Target_Client)
}

//RemoveRedPoint 移除小红点数据
func (u *RoleProp) RemoveRedPoint(redPointKey string) {
	delete(u.Data.RedPointsData, redPointKey)
}

//SyncRemoveRedPointList 同步批量移除小红点数据
func (u *RoleProp) SyncRemoveRedPointList(redPointKeyList *Proto.StringArry) {
	u.Sync("RemoveRedPointList", Message.PackArgs(redPointKeyList), true, Prop.Target_Client)
}

//RemoveRedPointList 批量移除小红点数据
func (u *RoleProp) RemoveRedPointList(redPointKeyList *Proto.StringArry) {
	for _, key := range redPointKeyList.Data {
		delete(u.Data.RedPointsData, key)
	}
}

// SyncUpdateTalentData 同步更新特工天赋
func (u *RoleProp) SyncUpdateTalentData(sid uint32, tid uint32, talent *Proto.TalentData) {
	u.Sync("UpdateTalentData", Message.PackArgs(sid, tid, talent), true, Prop.Target_Client)
}

// UpdateTalentData 更新特工天赋
func (u *RoleProp) UpdateTalentData(sid uint32, tid uint32, talent *Proto.TalentData) {
	sagents := u.Data.SpecialAgentList
	if sagents == nil {
		return
	}
	sagent, ok := sagents[sid]
	if !ok {
		return
	}

	if sagent.Talent == nil {
		return
	}
	myTalent, ok := sagent.Talent.TalentMap[tid]
	if !ok {
		return
	}
	myTalent.Study = talent.Study
	myTalent.Level = talent.Level
	myTalent.Unlock = talent.Unlock
	myTalent.GiveUp = talent.GiveUp
}

// SyncUpdateTalentPoint 同步更新特工科技点数
func (u *RoleProp) SyncUpdateTalentPoint(sid uint32, point uint32) {
	u.Sync("UpdateTalentPoint", Message.PackArgs(sid, point), true, Prop.Target_Client)
}

// UpdateTalentPoint 更新特工科技点数
func (u *RoleProp) UpdateTalentPoint(sid uint32, point uint32) {
	sagents := u.Data.SpecialAgentList
	if sagents == nil {
		return
	}
	sagent, ok := sagents[sid]
	if !ok {
		return
	}
	if sagent.Talent == nil {
		return
	}
	sagent.Talent.TalentPoint = point
}

func (u *RoleProp) SyncSetFastBattle(fb *Proto.FastBattle) {
	u.Sync("SetFastBattle", Message.PackArgs(fb), true, Prop.Target_Client)
}

func (u *RoleProp) SetFastBattle(fb *Proto.FastBattle) {
	u.Data.FastBattle = fb
}

func (u *RoleProp) SyncAddFastBattle(stage uint32, isSelf bool) {
	u.Sync("AddFastBattle", Message.PackArgs(stage, isSelf), true, Prop.Target_Client)
}

func (u *RoleProp) AddFastBattle(stage uint32, isSelf bool) {
	if isSelf {
		if u.Data.FastBattle.StageInfo[stage].MyTimes < u.Data.FastBattle.StageInfo[stage].MaxMyTimes {
			u.Data.FastBattle.StageInfo[stage].MyTimes++
		}
	} else {
		if u.Data.FastBattle.StageInfo[stage].OtherTimes < u.Data.FastBattle.StageInfo[stage].MaxOtherTimes {
			u.Data.FastBattle.StageInfo[stage].OtherTimes++
		}
	}
}

func (u *RoleProp) SyncAddEnergize(awardMulti float32) {
	u.Sync("AddEnergize", Message.PackArgs(awardMulti), true, Prop.Target_Client)
}

func (u *RoleProp) AddEnergize(awardMulti float32) {
	u.Data.FastBattle.EnergizeNum++
	u.Data.FastBattle.AwardMulti = awardMulti
}

func (u *RoleProp) SyncClearFastBattleMulti() {
	u.Sync("ClearFastBattleMulti", Message.PackArgs(), true, Prop.Target_Client)
}

func (u *RoleProp) ClearFastBattleMulti() {
	u.Data.FastBattle.AwardMulti = 0
}

func (u *RoleProp) SyncResetRoleDaily(time int64, requestSkillCountLimit, giveSkillCountLimit uint32) {
	u.Sync("ResetRoleDaily", Message.PackArgs(time, requestSkillCountLimit, giveSkillCountLimit), true)
}

func (u *RoleProp) ResetRoleDaily(time int64, requestSkillCountLimit, giveSkillCountLimit uint32) {
	u.Data.ShareSpoils.GiveEquipCount = 0
	u.Data.ShareSpoils.RequestSkillResetTimeStamp = time
	u.Data.ShareSpoils.RequestSkillCount = 0
	u.Data.ShareSpoils.RequestSkillCountLimit = requestSkillCountLimit
	u.Data.ShareSpoils.GiveSkillResetTimeStamp = time
	u.Data.ShareSpoils.GiveSkillCount = 0
	u.Data.ShareSpoils.GiveSkillCountLimit = giveSkillCountLimit
	u.Data.ShareSpoils.RequestSkillData = make(map[string]*Proto.RequestSkill)
	u.Data.ShareSpoils.GiveSkillToTargetCount = make(map[string]*Proto.GiveSkillCount)
}

func (u *RoleProp) SyncSetGiveEquipCount(giveEquipCount uint32) {
	u.Sync("SetGiveEquipCount", Message.PackArgs(giveEquipCount), true)
}

func (u *RoleProp) SetGiveEquipCount(giveEquipCount uint32) {
	u.Data.ShareSpoils.GiveEquipCount = giveEquipCount
}

func (u *RoleProp) SyncSetGiveSkillResetTimeStamp(time int64) {
	u.Sync("SetGiveSkillResetTimeStamp", Message.PackArgs(time), true, Prop.Target_Client)
}

func (u *RoleProp) SetGiveSkillResetTimeStamp(time int64) {
	u.Data.ShareSpoils.GiveSkillResetTimeStamp = time
}

func (u *RoleProp) SyncSetRequestSkillResetTimeStamp(time int64) {
	u.Sync("SetRequestSkillResetTimeStamp", Message.PackArgs(time), true, Prop.Target_Client)
}

func (u *RoleProp) SetRequestSkillResetTimeStamp(time int64) {
	u.Data.ShareSpoils.RequestSkillResetTimeStamp = time
}

func (u *RoleProp) SyncAddRequestSkill(uid string, requestSkill *Proto.RequestSkill, requestSkillCount uint32) {
	u.Sync("AddRequestSkill", Message.PackArgs(uid, requestSkill, requestSkillCount), true, Prop.Target_All_Clients)
}

func (u *RoleProp) AddRequestSkill(uid string, requestSkill *Proto.RequestSkill, requestSkillCount uint32) {
	if requestSkill.TeamMateGiveSkillCount == nil {
		requestSkill.TeamMateGiveSkillCount = map[string]uint32{}
	}

	u.Data.ShareSpoils.RequestSkillData[uid] = requestSkill
	u.Data.ShareSpoils.RequestSkillCount = requestSkillCount
}

func (u *RoleProp) SyncClearRequestSkill() {
	u.Sync("ClearRequestSkill", Message.PackArgs(), true, Prop.Target_Client)
}

func (u *RoleProp) ClearRequestSkill() {
	u.Data.ShareSpoils.RequestSkillData = make(map[string]*Proto.RequestSkill)
}

func (u *RoleProp) SyncSetReceivedSkillCount(uid string, receivedSkillCount uint32, teamMateId string, teamMateSkillCount uint32) {
	u.Sync("SetReceivedSkillCount", Message.PackArgs(uid, receivedSkillCount, teamMateId, teamMateSkillCount), true, Prop.Target_All_Clients)
}

func (u *RoleProp) SetReceivedSkillCount(uid string, receivedSkillCount uint32, teamMateId string, teamMateSkillCount uint32) {
	if requestSkill, ok := u.Data.ShareSpoils.RequestSkillData[uid]; ok {
		requestSkill.ReceivedSkillCount = receivedSkillCount
		requestSkill.TeamMateGiveSkillCount[teamMateId] = teamMateSkillCount
	}
}

func (u *RoleProp) SyncSetGiveSkillCount(giveSkillCount uint32, targetId string, giveSkillItemId, giveTargetSkillCount uint32) {
	u.Sync("SetGiveSkillCount", Message.PackArgs(giveSkillCount, targetId, giveSkillItemId, giveTargetSkillCount), true, Prop.Target_Client)
}

func (u *RoleProp) SetGiveSkillCount(giveSkillCount uint32, targetId string, giveSkillItemId, giveTargetSkillCount uint32) {
	u.Data.ShareSpoils.GiveSkillCount = giveSkillCount
	if giveSkillCount, ok := u.Data.ShareSpoils.GiveSkillToTargetCount[targetId]; ok {
		if giveSkillCount.GiveSkillToTargetCount == nil {
			giveSkillCount.GiveSkillToTargetCount = map[uint32]uint32{}
		}
		giveSkillCount.GiveSkillToTargetCount[giveSkillItemId] = giveTargetSkillCount
	} else {
		giveSkillCount := &Proto.GiveSkillCount{
			GiveSkillToTargetCount: map[uint32]uint32{},
		}

		giveSkillCount.GiveSkillToTargetCount[giveSkillItemId] = giveTargetSkillCount

		u.Data.ShareSpoils.GiveSkillToTargetCount[targetId] = giveSkillCount
	}
}

func (u *RoleProp) SyncChangeSupplyNum(id, num uint32) {
	u.Sync("ChangeSupplyNum", Message.PackArgs(id, num), true, Prop.Target_Client)
}
func (u *RoleProp) ChangeSupplyNum(id, num uint32) {
	state, ok := u.Data.SupplyInfo.BoxMold[id]
	if !ok {
		state = &Proto.State{}
		u.Data.SupplyInfo.BoxMold[id] = state
	}
	state.HoldNum -= num
	state.OpenNum += num
}

func (u *RoleProp) SyncAddSupplyNum(id, num uint32) {
	u.Sync("AddSupplyNum", Message.PackArgs(id, num), true, Prop.Target_Client)
}
func (u *RoleProp) AddSupplyNum(id, num uint32) {
	state, ok := u.Data.SupplyInfo.BoxMold[id]
	if !ok {
		state = &Proto.State{}
		u.Data.SupplyInfo.BoxMold[id] = state
	}
	state.HoldNum += num
}

func (u *RoleProp) SyncAddDiscountNum(id, num uint32) {
	u.Sync("AddDiscountNum", Message.PackArgs(id, num), true, Prop.Target_Client)
}
func (u *RoleProp) AddDiscountNum(id, num uint32) {
	u.Data.SupplyInfo.DiscountNum[id] += num
}

func (u *RoleProp) SyncResetBoxOpen() {
	u.Sync("ResetBoxOpen", Message.PackArgs(), true, Prop.Target_Client)
}
func (u *RoleProp) ResetBoxOpen() {
	for k := range u.Data.SupplyInfo.BoxMold {
		u.Data.SupplyInfo.BoxMold[k].OpenNum = 0
	}
}

func (u *RoleProp) SyncResetDiscountNum() {
	u.Sync("ResetDiscountNum", Message.PackArgs(), true, Prop.Target_Client)
}
func (u *RoleProp) ResetDiscountNum() {
	for k := range u.Data.SupplyInfo.DiscountNum {
		u.Data.SupplyInfo.DiscountNum[k] = 0
	}
}

//同步清空所有超时称号
func (u *RoleProp) SyncClearOverTimeTitle(deleteTitleList *Proto.Int32Array) {
	u.Sync("ClearOverTimeTitle", Message.PackArgs(deleteTitleList), true, Prop.Target_All_Clients)
}
func (u *RoleProp) ClearOverTimeTitle(deleteTitleList *Proto.Int32Array) {
	for _, ID := range deleteTitleList.Data {
		delete(u.Data.Title.TitleList, uint32(ID))
	}
}

//同步增加玩家已获得称号
func (u *RoleProp) SyncAddTitle(titleID uint32, startTime int64, lostTime int64) {
	u.Sync("AddTitle", Message.PackArgs(titleID, startTime, lostTime), true, Prop.Target_All_Clients)
}
func (u *RoleProp) AddTitle(titleID uint32, startTime int64, lostTime int64) {
	title := &Proto.TitleInfo{titleID, startTime, lostTime}
	u.Data.Title.TitleList[titleID] = title

}

//改变使用称号
func (u *RoleProp) SyncUpdateTitle(titleID uint32) {
	u.Sync("UpdateTitle", Message.PackArgs(titleID), true, Prop.Target_All_Clients)
}

func (u *RoleProp) UpdateTitle(titleID uint32) {
	u.Data.Title.TitleID = titleID
}

// SyncSetUID 同步设置玩家uid
func (u *RoleProp) SyncSetUID(uid uint64) {
	if u.Data.Base.UID != uid {
		u.Sync("SetUID", Message.PackArgs(uid), true, Prop.Target_Client)
	}
}

// SetUID 设置玩家uid
func (u *RoleProp) SetUID(uid uint64) {
	u.Data.Base.UID = uid
}

//SyncGetSeasonEndAward 拿到赛季结束奖励
func (u *RoleProp) SyncGetSeasonEndAward() {
	u.Sync("GetSeasonEndAward", Message.PackArgs(), true)
}
func (u *RoleProp) GetSeasonEndAward() {
	u.Data.SeasonInfo.GetAwards = true
}

//SyncAddSeasonScore 增加赛季积分
func (u *RoleProp) SyncAddSeasonScore(add uint32) {
	u.Sync("AddSeasonScore", Message.PackArgs(add), true, Prop.Target_Client)
}
func (u *RoleProp) AddSeasonScore(add uint32) {
	u.Data.SeasonInfo.MyScore += add
}

//SyncUpdateSeasonInfo 更新赛季信息
func (u *RoleProp) SyncUpdateSeasonInfo(info *Proto.RoleSeasonInfo) {
	u.Sync("UpdateSeasonInfo", Message.PackArgs(info), true, Prop.Target_Client)
}
func (u *RoleProp) UpdateSeasonInfo(info *Proto.RoleSeasonInfo) {
	u.Data.SeasonInfo = info
}

//SyncGetSeasonLevelAward 领取赛季等级奖励
func (u *RoleProp) SyncGetSeasonLevelAward(id uint32) {
	u.Sync("GetSeasonLevelAward", Message.PackArgs(id), true, Prop.Target_Client)
}
func (u *RoleProp) GetSeasonLevelAward(id uint32) {
	if u.Data.SeasonInfo.LevelAwards == nil {
		u.Data.SeasonInfo.LevelAwards = make(map[uint32]bool)
	}
	u.Data.SeasonInfo.LevelAwards[id] = true
}

func (u *RoleProp) SyncSetSeasonAwardNotify(notify bool) {
	u.Sync("SetSeasonAwardNotify", Message.PackArgs(notify), true, Prop.Target_Client)
}
func (u *RoleProp) SetSeasonAwardNotify(notify bool) {
	u.Data.SeasonInfo.Notify = notify
}

func (u *RoleProp) SyncSetRaidProgress(progress uint32) {
	if u.Data.Base.RaidProgress != progress {
		u.Sync("SetRaidProgress", Message.PackArgs(progress), true, Prop.Target_All_Clients)
	}
}

func (u *RoleProp) SetRaidProgress(progress uint32) {
	u.Data.Base.RaidProgress = progress
}

//SyncUpdateCommanderLv 同步更新指挥官等级
func (u *RoleProp) SyncUpdateCommanderLv(lv uint32, exp uint64) {
	u.Sync("UpdateCommanderLv", Message.PackArgs(lv, exp), true, Prop.Target_All_Clients)
}

//UpdateCommanderLv 更新指挥官等级数据
func (u *RoleProp) UpdateCommanderLv(lv uint32, exp uint64) {
	u.Data.Base.Level = lv
	u.Data.Base.Exp = exp
}
