package main

import (
	cConst "Cinder/Base/Const"
	"Cinder/Space"
	"Daisy/ActivityTimer"
	"Daisy/Battle/drop"
	"Daisy/Const"
	"Daisy/DB"
	"Daisy/Data"
	"Daisy/ErrorCode"
	"Daisy/Fight/attraffix"
	"Daisy/ItemProto"
	"Daisy/Prop"
	"Daisy/Proto"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func (user *_User) RPC_GM(cmd, args string) int32 {
	if user.role != nil {
		return user.role.GM(cmd, args)
	}

	return ErrorCode.Failure
}

func (r *_Role) GM(cmd, args string) int32 {
	r.Debug("RPC_GM", cmd, args)

	switch strings.ToLower(cmd) {
	case "additem":
		return r.gm_addItem(args)
	case "removeitem":
		return r.gm_removeItem(args)
	case "addrobot":
		return r.gm_addRobot(args)
	case "drop":
		return r.gm_drop(args)
	case "addspecialagent":
		//todo 临时代码 针对新账号默认解锁 一个特工+对应一个build
		specialAgentID, err := strconv.Atoi(args)
		if err != nil {
			return ErrorCode.Failure
		}

		if _, ok := r.prop.Data.SpecialAgentList[uint32(specialAgentID)]; ok {
			return ErrorCode.Success
		}

		specialAgent := &Proto.SpecialAgent{
			Base: &Proto.SpecialAgentBase{
				ConfigID: uint32(specialAgentID),
				Level:    1,
				Exp:      0,
				GainTime: time.Now().Unix(),
			},
			Talent: r.initTalent(uint32(specialAgentID)),
		}
		r.prop.SyncAddSpecialAgent(specialAgent)
	case "ticket":
		team := r.GetSpace().(*_Team)
		team.prop.SyncSetOwnTickets(100)
	case "addmoney":
		list := strings.Split(args, ",")
		if len(list) == 2 {
			moneyID, err := strconv.Atoi(list[0])
			if err != nil {
				return ErrorCode.Failure
			}

			moneyNum, err := strconv.Atoi(list[1])
			if err != nil {
				return ErrorCode.Failure
			}

			r.AddMoney(uint32(moneyNum), uint32(moneyID), Const.NORMAL_ACTION)
		}
	case "changeattribute":
		list := strings.Split(args, ",")
		if len(list) == 2 {
			attrType, err := strconv.Atoi(list[0])
			if err != nil {
				return ErrorCode.Failure
			}

			attrValue, err := strconv.Atoi(list[1])
			if err != nil {
				return ErrorCode.Failure
			}

			affix := attraffix.AttrAffix{
				Field: attraffix.Field(attrType),
				ParaA: float32(attrValue),
				ParaB: float64(attrValue),
			}
			r.debugAttraffix = append(r.debugAttraffix, affix)
		}
	case "clearallattribute":
		r.debugAttraffix = []attraffix.AttrAffix{}
	case "addredpoint":

		data, ok := r.prop.Data.RedPointsData[args]
		if !ok {
			data = &Proto.RedPointInfo{
				Value:      1,
				CreateTime: time.Now().Unix(),
			}
		} else {
			data.Value++
			data.CreateTime = time.Now().Unix()
		}
		r.prop.SyncAddRedPoint(args, data)

	case "removeredpoint":
		data, ok := r.prop.Data.RedPointsData[args]
		if ok {
			data.Value--
			if data.Value == 0 {
				r.prop.SyncRemoveRedPoint(args)
			} else {
				r.prop.SyncAddRedPoint(args, data)
			}
		}
	case "addcommanderexp":
		exp, err := strconv.Atoi(args)
		if err != nil {
			return ErrorCode.Failure
		}

		r.addCommanderExp(uint32(exp))
	case "addfightingspecialagentexp":
		exp, err := strconv.Atoi(args)
		if err != nil {
			return ErrorCode.Failure
		}

		r.addFightingSpecialAgentExp(uint64(exp))
	case "addspecialagentexp":
		list := strings.Split(args, ",")
		if len(list) == 2 {
			id, err := strconv.Atoi(list[0])
			if err != nil {
				return ErrorCode.Failure
			}

			//是否拥有该特工
			if _, ok := r.prop.Data.SpecialAgentList[uint32(id)]; !ok {
				return ErrorCode.Failure
			}

			exp, err := strconv.Atoi(list[1])
			if err != nil {
				return ErrorCode.Failure
			}

			r.addSpecialAgentExp(uint32(id), uint64(exp))
		}
	case "changespecialagentlv":
		list := strings.Split(args, ",")
		if len(list) == 2 {
			id, err := strconv.Atoi(list[0])
			if err != nil {
				return ErrorCode.Failure
			}

			specialAgent, ok := r.prop.Data.SpecialAgentList[uint32(id)]
			//是否拥有该特工
			if !ok {
				return ErrorCode.Failure
			}

			lv, err := strconv.Atoi(list[1])
			if err != nil {
				return ErrorCode.Failure
			}

			specialAgent.Base.Level = uint32(lv)
			specialAgent.Base.Exp = 0
			specialAgent.Talent = r.initTalent(uint32(id))

			r.prop.SyncAddSpecialAgent(specialAgent)
		}
	case "mail":
		return r.gm_mail(args)
	case "mailto":
		return r.gm_mailto(args)
	case "clearequipbag":
		container := r.GetContainer(int32(Proto.ContainerEnum_EquipBag))
		container.Traversal(func(item ItemProto.IItem) bool {
			container.RemoveItem(item.GetPos())
			return true
		})
		var num uint32
		container.Traversal(func(item ItemProto.IItem) bool {
			num++
			return true
		})
		r.Debug("包裹道具数量：", num)
	case "clearskillbag":
		container := r.GetContainer(int32(Proto.ContainerEnum_SkillBag))
		container.Traversal(func(item ItemProto.IItem) bool {
			container.RemoveItem(item.GetPos())
			return true
		})
		var num uint32
		container.Traversal(func(item ItemProto.IItem) bool {
			num++
			return true
		})
		r.Debug("包裹道具数量：", num)
	case "addtitle":
		return r.gm_addTitle(args)
	case "setroad":
		return r.gm_setRoad(args)
	case "setteampro":
		return r.gm_setTeamProgress(args)
	case "addseasonscore":
		team := r.GetOwnerUser().GetSpace().(*_Team)
		team.ChangeSeasonScore(team.prop.Data.SeasonInfo.TeamScore + 10)
		team.TraversalActor(func(actor Space.IActor) {
			role := actor.(*_Role)
			//重复打不加分
			role.prop.SyncAddSeasonScore(10)
			user := r.GetOwnerUser()
			if user == nil {
				r.Error("[findNewTitle] role's user is nil")
			} else {
				user.Rpc(cConst.Game, "RPC_UpdateChatUserActivateData")
			}
		})
		data := team.GetRankList()
		fmt.Println(data)
	case "endseason":
		act := ActivityTimer.GetActivities().GetActivity(Const.SEASON)
		act.End()
	case "startseason":
		act := ActivityTimer.GetActivities().GetActivity(Const.SEASON)
		act.Start()
	case "initseason":
		act := &SeasonActivity{}
		data := DB.ActivityInfo{
			Step:          1,
			StartTime:     time.Now().Unix(),
			NextStartTime: time.Now().Unix(),
			EndTime:       time.Now().Unix() + int64(act.GetInterval()+act.GetLast())}
		DB.SetActivitiesInfo(Const.SEASON_KEY, &data)
		ActivityTimer.GetActivities().RegisterActivity(&SeasonActivity{})
	case "fienddatachange":
		return r.gm_fiendDataChange(args)
	case "addbuild":
		list := strings.Split(args, ",")
		if len(list) == 2 {
			name := list[0]
			specialAgentID, err := strconv.Atoi(list[1])
			if err != nil {
				return ErrorCode.Failure
			}

			r.createBuild(name, uint32(specialAgentID))
		}
	default:
		r.Error("unknown cmd", cmd, len(cmd))
		return ErrorCode.Failure
	}

	return ErrorCode.Success
}

func (r *_Role) gm_addRobot(args string) int32 {
	r.Debug("gm_addRobot ", args)
	configID, err := strconv.Atoi(args)
	if err != nil {
		return ErrorCode.ArgsWrong
	}

	team := r.GetSpace().(*_Team)
	_, err = team.AddRobot("robot", uint32(configID))
	if err != nil {
		return ErrorCode.ArgsWrong
	}

	return ErrorCode.Success
}

func (r *_Role) gm_setFightingConfigID(args string) int32 {
	r.Debug("gm_setFightingConfigID ", args)
	configID, err := strconv.Atoi(args)
	if err != nil {
		return ErrorCode.ArgsWrong
	}

	if card, ok := r.prop.Data.BuildMap[r.prop.Data.FightingBuildID]; ok && card.SpecialAgentID == uint32(configID) {
		return ErrorCode.Success
	}

	for key, value := range r.prop.Data.BuildMap {
		if value.SpecialAgentID == uint32(configID) {
			r.prop.SyncSetFightingBuildID(key)
			r.FlushToDB()
			r.FlushToCache()
			return ErrorCode.Success
		}
	}

	return ErrorCode.ArgsWrong
}

func (r *_Role) gm_addItem(args string) int32 {
	list := strings.Split(args, ",")
	if len(list) >= 3 {
		id, err := strconv.Atoi(list[0])
		if err != nil {
			return ErrorCode.Success
		}
		typ, err := strconv.Atoi(list[1])
		if err != nil {
			return ErrorCode.Success
		}
		num, err := strconv.Atoi(list[2])
		if err != nil {
			return ErrorCode.Success
		}
		var _args []uint32
		_args = make([]uint32, 0)
		for _, v := range list[3:] {
			u, err1 := strconv.Atoi(v)
			if err1 == nil {
				_args = append(_args, uint32(u))
			}
		}
		addNum := r.AddItem(uint32(id), uint32(typ), uint32(num), _args...)
		r.Debug("gm addItem num=", addNum)
	}
	return ErrorCode.Success
}

func (r *_Role) gm_removeItem(args string) int32 {
	list := strings.Split(args, ",")
	if len(list) == 3 {
		id, err := strconv.Atoi(list[0])
		if err != nil {
			return ErrorCode.Success
		}
		typ, err := strconv.Atoi(list[1])
		if err != nil {
			return ErrorCode.Success
		}
		num, err := strconv.Atoi(list[2])
		if err != nil {
			return ErrorCode.Success
		}
		remove := r.RemoveItem(uint32(id), uint32(typ), uint32(num))
		r.Debug("gm removeitem num=", remove)
	} else if len(list) == 1 {
		typ, err := strconv.Atoi(list[0])
		if err != nil {
			return ErrorCode.Failure
		}
		container := r.GetContainer(int32(typ))
		container.Traversal(func(item ItemProto.IItem) bool {
			container.RemoveItem(item.GetPos())
			return true
		})
		var num uint32
		container.Traversal(func(item ItemProto.IItem) bool {
			num++
			return true
		})
		r.Debug("包裹道具数量：", num)
	}
	return ErrorCode.Success
}

func (r *_Role) gm_equip(args string) int32 {
	list := make([]*Proto.DropMaterial, 0)
	list = append(list, &Proto.DropMaterial{
		MaterialId:   1000500546,
		MaterialType: 10,
		MaterialNum:  1,
	})
	list = append(list, &Proto.DropMaterial{
		MaterialId:   1001001047,
		MaterialType: 10,
		MaterialNum:  1,
	})
	list = append(list, &Proto.DropMaterial{
		MaterialId:   1000500546,
		MaterialType: 10,
		MaterialNum:  1,
	})
	list = append(list, &Proto.DropMaterial{
		MaterialId:   1001,
		MaterialType: 11,
		MaterialNum:  1,
	})
	list = append(list, &Proto.DropMaterial{
		MaterialId:   1002,
		MaterialType: 11,
		MaterialNum:  1,
	})

	var removeList []*Proto.GetItemData
	if _, ok := r.CanAddItemList(list); ok == true {
		removeList = r.AddItemList(list)
	}

	list2 := make([]ItemProto.IItem, 0)
	for _, v := range list {
		item := r.CreateItem(v.MaterialId, v.MaterialType, v.MaterialNum)
		if item != nil {
			list2 = append(list2, item)
		}
	}

	if _, ok := r.CanAddItemIList(list2); ok == true {
		r.AddItems(list2)
	}

	list3 := make([]ItemProto.IItem, 0)
	list4 := make([]string, 0)
	for _, v := range removeList {
		if v.ItemData != nil {
			list3 = append(list3, ItemProto.CreateIItemByData(v.ItemData))
			list4 = append(list4, v.ItemData.Base.ID)
		}
	}
	r.SellEquip(list4)
	//	r.RemoveItemList(list3)
	/*
		r.AddItem(1010001, uint32(Proto.ItemEnum_Equipment), 14)

		itemList := make([]*Proto.Item, 0)
		for _, v := range r.prop.Data.ItemContainerMap[int32(Proto.ContainerEnum_BuildEquipBag)].ItemMap {
			itemList = append(itemList, v)
		}
		for _, v := range r.prop.Data.ItemContainerMap[int32(Proto.ContainerEnum_EquipBag)].ItemMap {
			itemList = append(itemList, v)
		}

		for _, v := range itemList {
			v.ExpandData.InUse += 1
			r.UpdataItem(v)
		}

		for _, v := range itemList {
			v.ExpandData.InUse -= 1
			r.UpdataItem(v)
		}
	*/
	return ErrorCode.Success
}

func (r *_Role) gm_drop(args string) int32 {
	id, err := strconv.Atoi(args)
	if err != nil {
		return ErrorCode.Failure
	}
	r.Debug("gm 测试开始 drop： id=", id)
	drop := drop.Drop{}
	_, list := drop.Drop(uint32(id), 0, 0)
	if len(list) == 0 {
		r.Debug("gm 测试 drop 没有任何掉落, 测试结束")
		return ErrorCode.Success
	}

	for _, d := range list {
		r.Debug("掉落=", d)
	}

	for _, v := range list {
		num := r.AddItem(v.MaterialId, v.MaterialType, v.MaterialNum)
		if num == 0 {
			r.Debug("放入包裹失败")
		}
		/*
			r.Debug(fmt.Sprintf("装备ID=%d, 品质=%d, 孔数=%d, 是否珍品:%t", item.Base.ConfigID, item.EquipmentData.Quality, len(item.EquipmentData.Socket), item.EquipmentData.IsPrecious))
			for _,v1 := range item.EquipmentData.Affixes {
				if v1.Const == false {
					continue
				}
				r.Debug(fmt.Sprintf("固定词缀: id=%d, Value=%f", v1.AffixID, v1.Value))
			}
			for _,v2 := range item.EquipmentData.Affixes {
				if v2.Const == true {
					continue
				}
				r.Debug(fmt.Sprintf("随机词缀: id=%d, Value=%f", v2.AffixID, v2.Value))
			}*/
	}
	r.Debug("gm 测试 drop： 结束")
	return ErrorCode.Success
}

func (r *_Role) gm_addTitle(args string) int32 {
	fmt.Println(args)
	id, err := strconv.Atoi(args)
	fmt.Println(id)
	if err != nil {
		return ErrorCode.Failure
	}
	r.AddTitle(uint32(id))
	return ErrorCode.Success
}

func (r *_Role) gm_setRoad(args string) int32 {
	team := r.GetSpace().(*_Team)
	err := team.SetRoadByName(args)
	if err != nil {
		r.Errorf("gm_setRoad err:%s", err)
		return ErrorCode.Failure
	}
	r.Debugf("gm_setRoad args:%s success", args)
	return ErrorCode.Success
}

func (r *_Role) gm_setTeamProgress(args string) int32 {
	team := r.GetSpace().(*_Team)
	progress, err := strconv.Atoi(args)
	if err != nil {
		r.Errorf("gm_setTeamProgress args:%s err:%s", args, err)
		return ErrorCode.Failure
	}

	if progress <= 0 || progress > len(Data.GetSceneConfig().BattleArea_ConfigItems) {
		r.Errorf("gm_setTeamProgress progress:%d out range", progress)
		return ErrorCode.Failure
	}

	team.prop.SyncSetRaidProgress(uint32(progress))
	team.TraversalActor(func(ia Space.IActor) {
		if _, ok := team.prop.Data.Base.Members[ia.GetID()]; ok {
			roleProp := ia.GetProp().(*Prop.RoleProp)
			if roleProp.Data.Base.RaidProgress < team.prop.Data.Raid.Progress {
				roleProp.SyncSetRaidProgress(team.prop.Data.Raid.Progress)
			}
		}
	})

	r.Debugf("gm_setTeamProgress args:%s success", args)
	return ErrorCode.Success
}
