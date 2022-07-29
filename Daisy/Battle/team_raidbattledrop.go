package main

import (
	"Cinder/Space"
	"Daisy/Data"
	"Daisy/ErrorCode"
	"Daisy/ItemProto"
	"Daisy/Prop"
	"Daisy/Proto"
	"errors"
	log "github.com/cihub/seelog"
	"time"
)

type enemyDrop map[uint32]*Spoil

func (ed enemyDrop) AddDrop(pawnID uint32, drops []*Proto.DropMaterial, money uint32, actorExp uint32, specialExp uint32) {
	s := &Spoil{
		MaterialDrop:   drops,
		MoneyDrop:      money,
		ActorExpDrop:   actorExp,
		SpecialExpdrop: specialExp,
	}
	ed[pawnID] = s
}

func (ed enemyDrop) RemoveDrop(pawnID uint32) (*Spoil, error) {
	if drops, ok := ed[pawnID]; ok {
		delete(ed, pawnID)
		return drops, nil
	} else {
		return nil, errors.New("not exist")
	}
}

func (ed enemyDrop) TraversalDrop(cb func(item *Proto.DropMaterial)) {
	for _, value := range ed {
		for _, item := range value.MaterialDrop {
			cb(item)
		}
	}
}

type _RaidBattleDropModel struct {
	items map[string]enemyDrop
}

func (team *_Team) EnterRaidBattleDrop(fight *Proto.FightResult) {
	if fight.Outcome[0] == Proto.Camp_Blue {
		log.Debug("[EnterRaidBattleDrop]进入战斗 玩家战斗输了 没有掉落")
		return
	}

	team.raidBattleDrop = &_RaidBattleDropModel{
		items: make(map[string]enemyDrop),
	}
	checkpoint := team.prop.Data.Raid.Progress

	cfg, ok := Data.GetSceneConfig().BattleArea_ConfigItems[checkpoint]
	if !ok {
		log.Error("[EnterRaidBattleDrop] teamHangUpAward failed, 找不到对应关卡表，掉落失败", checkpoint)
		return
	}
	team.TraversalActor(func(actor Space.IActor) {
		role := actor.(*_Role)
		ed := make(enemyDrop)
		team.raidBattleDrop.items[role.GetID()] = ed
		if team.raidBattle.lastFight.IsBoss {
			items, err := role.Drop(cfg.BossDropID)
			if err != nil {
				role.Error("[EnterRaidBattleDrop] boss战 生成掉落品失败 ", cfg.DropID)
			}
			ed.AddDrop(checkpoint, items, cfg.BossDropMoney, cfg.BossActorDropExp, cfg.BossSpecialDropExp)
		} else {
			items, err := role.Drop(cfg.GetDropID())
			if err != nil {
				role.Error("[EnterRaidBattleDrop] 非boss战 生成掉落品失败 ", cfg.DropID)
			}
			ed.AddDrop(checkpoint, items, cfg.DropMoney, cfg.ActorDropExp, cfg.SpecialDropExp)
		}
	})
}

// LeaveRaidBattleDrop 离开战斗时，结算所有离线玩家的收益
func (team *_Team) LeaveRaidBattleDrop() {

	// 判断胜负
	if team.GetCurRaidFight().Outcome[0] == Proto.Camp_Blue {
		return
	}
	if team.raidBattleDrop == nil {
		log.Error("[LeaveRaidBattleDrop] 玩家没有掉落")
		return
	}

	// 创建队伍成员信息列表
	var memberIDList []*Proto.OwnerInfo
	for id, _ := range team.prop.Data.Base.Members {
		owner := &Proto.OwnerInfo{id, team.GetActorByUserID(id).GetProp().(*Prop.RoleProp).Data.Base.Name}
		memberIDList = append(memberIDList, owner)
	}

	team.TraversalActor(func(actor Space.IActor) {
		role := actor.(*_Role)

		ed, ok := team.raidBattleDrop.items[role.GetID()]
		if ok {
			offlineUpdate := false
			if role.GetOwnerUser() == nil {
				// 玩家不在线 并且超过24小时，那么掉落不进背包
				roleLastOfflineTime := time.Unix(role.prop.Data.Base.LastLogoutTime, 0)
				maxOff := roleLastOfflineTime.Add(AccumulateTime * time.Hour)

				if maxOff.Before(time.Now()) {
					log.Debug("[LeaveRaidBattleDrop] 玩家离线时间超过24小时，不能领取战斗奖励")
					return
				} else {
					offlineUpdate = true
				}
			}

			awardData := &Proto.OfflineAwardData{}

			for _, spoil := range ed {
				itemData := role.AddItemList(spoil.MaterialDrop)
				for _, val := range itemData {
					if val.GetNum == 0 {
						continue
					}
					awardItem := &Proto.OfflineAwardItem{
						ID:       val.ItemData.Base.ID,
						Type:     val.ItemData.Base.Type,
						Num:      val.GetNum,
						ConfigID: val.ItemData.Base.ConfigID,
					}
					awardData.OfflineAwardItems = append(awardData.OfflineAwardItems, awardItem)

					if val.ItemData.Base.Type == Proto.ItemEnum_Equipment {
						val.ItemData.EquipmentData.OwnerTeamMemberList = memberIDList
						role.UpdataItem(ItemProto.CreateIItemByData(val.ItemData))
					}
				}
				awardData.Money += spoil.MoneyDrop
				awardData.ActorExp += spoil.ActorExpDrop
				awardData.SpecialAgentExp += spoil.SpecialExpdrop
			}

			role.prop.SyncAddGold(awardData.Money)
			role.addCommanderExp(awardData.ActorExp) // 更新指挥官经验
			role.addFightingSpecialAgentExp(role.expBonus(uint64(awardData.SpecialAgentExp)))
			// 离线数据更新 保证在线的玩家离线数据不更新
			if offlineUpdate {
				// 为了防止玩家在线但是客户端延迟，没有发送RPC_BattleEnd。在统一结算的时候判断玩家是否在线
				role.prop.SyncAddOfflineAward(awardData)
				log.Debug("[LeaveRaidBattleDrop] 客户端延迟，离线结算个人挂机收益", awardData)
			}
		}
	})
	team.raidBattleDrop = nil
}

// OnRaidBattleOnlineDrop 给在线的玩家发放奖励
func (team *_Team) OnRaidBattleOnlineDrop(roleID string) (int32, *Proto.OfflineAwardData) {
	// 判断胜负
	if team.GetCurRaidFight().Outcome[0] == Proto.Camp_Blue {
		return ErrorCode.Success, nil
	}

	ia, err := team.GetActor(roleID)
	if err != nil {
		return ErrorCode.GetRoleError, nil
	}

	role := ia.(*_Role)

	if team.raidBattleDrop == nil {
		role.Error("[OnRaidBattleOnlineDrop] 玩家没有掉落")
		return ErrorCode.Success, nil
	}

	ed, ok := team.raidBattleDrop.items[roleID]
	if !ok {
		role.Error("玩家没有掉落")
		return ErrorCode.Success, nil
	}
	awardData := &Proto.OfflineAwardData{}

	// 创建队伍成员信息列表
	var memberInfoList []*Proto.OwnerInfo
	for id, _ := range team.prop.Data.Base.Members {
		owner := &Proto.OwnerInfo{id, team.GetActorByUserID(id).GetProp().(*Prop.RoleProp).Data.Base.Name}
		memberInfoList = append(memberInfoList, owner)
	}

	// 遍历每一只怪的掉落
	for _, spoil := range ed {
		itemData := role.AddItemList(spoil.MaterialDrop)
		for _, val := range itemData {
			if val.GetNum == 0 {
				continue
			}
			awardItem := &Proto.OfflineAwardItem{
				ID:       val.ItemData.Base.ID,
				Type:     val.ItemData.Base.Type,
				Num:      val.GetNum,
				ConfigID: val.ItemData.Base.ConfigID,
			}
			awardData.OfflineAwardItems = append(awardData.OfflineAwardItems, awardItem)

			if val.ItemData.Base.Type == Proto.ItemEnum_Equipment {
				val.ItemData.EquipmentData.OwnerTeamMemberList = memberInfoList
				role.UpdataItem(ItemProto.CreateIItemByData(val.ItemData))
			}
		}
		awardData.Money += spoil.MoneyDrop
		awardData.ActorExp += spoil.ActorExpDrop
		awardData.SpecialAgentExp += spoil.SpecialExpdrop
	}
	role.prop.SyncAddGold(awardData.Money)
	role.addCommanderExp(awardData.ActorExp) // 更新指挥官经验
	role.addFightingSpecialAgentExp(role.expBonus(uint64(awardData.SpecialAgentExp)))
	delete(team.raidBattleDrop.items, roleID)
	return ErrorCode.Success, awardData
}

func (team *_Team) getDeadEnemy(fight *Proto.FightResult) []*Proto.PawnInfo {
	deads := make([]*Proto.PawnInfo, 0)
	for _, value := range fight.PawnInfos {
		if value.Camp == Proto.Camp_Blue && value.FightEndDead {
			deads = append(deads, value)
		}
	}
	return deads
}
