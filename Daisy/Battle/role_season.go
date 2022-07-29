package main

import (
	"Daisy/Battle/drop"
	"Daisy/DB"
	"Daisy/Data"
	"Daisy/ErrorCode"
	"Daisy/Proto"
	"fmt"
	"strconv"
	"strings"
)

func (r *_Role) SeasonEnd(Season uint32) {

	team := r.GetSpace().(*_Team)
	r.prop.Data.SeasonInfo.SaveTeamID = team.GetID()
	r.prop.Data.SeasonInfo.SaveSeasonID = team.prop.Data.SeasonInfo.SeasonID
	r.prop.Data.SeasonInfo.SaveSeasonScore = team.prop.Data.SeasonInfo.TeamScore
	r.prop.Data.SeasonInfo.SaveSeasonPlace = team.GetLegendRankPlace(r.prop.Data.SeasonInfo.SaveSeasonID)
	r.prop.Data.SeasonInfo.SaveSeasonPlacePercent = DB.DefaultIRankUtil.GetPlacePercentage(DB.SeasonRank, r.prop.Data.SeasonInfo.SaveSeasonID, r.prop.Data.SeasonInfo.SaveTeamID)
	r.Infof("赛季结束:SaveTeamID=%s,SaveSeasonID=%d,SaveSeasonScore=%d,SaveSeasonPlace=%d", r.prop.Data.SeasonInfo.SaveTeamID, r.prop.Data.SeasonInfo.SaveSeasonID, r.prop.Data.SeasonInfo.SaveSeasonScore, r.prop.Data.SeasonInfo.SaveSeasonPlace)
	r.prop.SyncUpdateSeasonInfo(r.prop.Data.SeasonInfo)
	r.SeasonEndAward()
}

func (r *_Role) SeasonEndAward() {
	//给过奖励了
	if r.prop.Data.SeasonInfo.GetAwards == true{
		return
	}

	if r.prop.Data.SeasonInfo.SaveSeasonID == 0 {
		r.Infof("赛季结束,id=0")
		return
	}

	var stage,seasonID uint32
	seasonID = r.prop.Data.SeasonInfo.SaveSeasonID
	team := r.GetSpace().(*_Team)
	level := team.GetSeasonLevelByScore(seasonID, r.prop.Data.SeasonInfo.SaveSeasonScore)

	levelNum := uint32(0)
	for i:=uint32(1); i<=level;i++ {
		//自动发未领取的等级奖励
		_, num := r.GetSeasonLevelAward(i)
		levelNum += num
		}
	r.SendMail(getSeasonGeneralString(SeasonLevelAwardMailTile), getSeasonGeneralString(SeasonLevelAwardMailContent)+fmt.Sprint(levelNum), r.GetID(), r.prop.Data.Base.Name,false, nil)

	config := team.GetSeasonLevelConfig(seasonID, team.GetSeasonLevelByScore(seasonID, r.prop.Data.SeasonInfo.SaveSeasonScore))

	if config != nil {
		stage = config.SeasonStageID
	}
	stageItemList := make([]*Proto.MailAttachment, 0)
	itemList := make([]*Proto.DropMaterial, 0)
	//赛季奖励
	for _,v := range Data.GetSeasonConfig().SeasonReward_ConfigItems {
		if seasonID == v.SeasonID && stage == v.SeasonStageID{
			for _,_v := range v.SeasonReward{
				list := strings.Split(_v, ",")
				if len(list) != 3 {
					continue
				}

				id,_ := strconv.Atoi(list[0])
				typ,_ := strconv.Atoi(list[1])
				num,_ := strconv.Atoi(list[2])
				r.Infof("赛季获得阶段奖励:seasin=%d, stage=%d,道具:%d,%d,%d", seasonID, stage, id, typ, num)
				itemList = append(itemList, &Proto.DropMaterial{
					MaterialId:uint32(id),
					MaterialType:uint32(typ),
					MaterialNum: uint32(num),
				})
			}
			break
		}
	}
	iitemList := r.CreatProp(itemList)
	for _,iv := range iitemList{
		stageItemList = append(stageItemList, &Proto.MailAttachment{Data:iv.GetData()})
	}
	r.SendMail(getSeasonGeneralString(SeasonStageAwardMailTile), getSeasonGeneralString(SeasonStageAwardMailContent), r.GetID(), r.prop.Data.Base.Name,false, stageItemList)

	//传说奖励
	legendItemList := make([]*Proto.MailAttachment, 0)
	itemList = make([]*Proto.DropMaterial, 0)
	place := r.prop.Data.SeasonInfo.SaveSeasonPlace
	for _,v2 := range Data.GetSeasonConfig().SeasonLegendReward_ConfigItems{
		if seasonID != v2.SeasonID {
			continue
		}
		if int32(place) >= v2.LegendStage[0] && int32(place) <= v2.LegendStage[1] {
			for _,_v := range v2.LegendStageReward{
				list := strings.Split(_v, ",")
				if len(list) != 3 {
					continue
				}
				id,_ := strconv.Atoi(list[0])
				typ,_ := strconv.Atoi(list[1])
				num,_ := strconv.Atoi(list[2])
				r.Infof("赛季获得传说排名:seasin=%d, place=%d,道具:%d,%d,%d", seasonID, stage, id, typ, num)
				itemList = append(itemList, &Proto.DropMaterial{
					MaterialId:uint32(id),
					MaterialType:uint32(typ),
					MaterialNum: uint32(num),
				})
			}
			break
		}
	}
	iitemList = r.CreatProp(itemList)
	for _,iv := range iitemList{
		legendItemList = append(legendItemList, &Proto.MailAttachment{Data:iv.GetData()})
	}

	r.SendMail(getSeasonGeneralString(SeasonLegendAwardMailTile), getSeasonGeneralString(SeasonLegendAwardMailContent), r.GetID(), r.prop.Data.Base.Name,false, stageItemList)

	r.prop.SyncGetSeasonEndAward()
	r.prop.SyncSetSeasonAwardNotify(true)
}

func (r *_Role) SeasonStart(Season uint32) {
	r.prop.Data.SeasonInfo.GetAwards = false
	r.prop.Data.SeasonInfo.LevelAwards = make(map[uint32]bool)
	//Todo 是否记录玩家未领取的奖励
	r.prop.SyncUpdateSeasonInfo(r.prop.Data.SeasonInfo)
}

//RPC_GetSeasonRankPercentage请求百分比
//返回 百分比排名，错误码
func (user *_User) RPC_GetSeasonRankPercentage() (uint32,int32) {
	if user.role == nil {
		return 0, ErrorCode.RoleIsNil
	}
	team := user.GetSpace().(*_Team)
	return team.GetSeasonPlace(team.prop.Data.SeasonInfo.SeasonID),ErrorCode.Success
}

//RPC_GetSeasonRank请求排行榜
//返回 排行榜，错误码
func (user *_User) RPC_GetSeasonRank() (*Proto.RankData, int32) {
	if user.role == nil {
		return nil, ErrorCode.RoleIsNil
	}
	team := user.GetSpace().(*_Team)
	return team.GetRankList(),ErrorCode.Success
}

//RPC_GetAwardDetails 获得奖励数据
//返回 经验，道具数量，错误码
func (user *_User) RPC_GetAwardDetails(id uint32) (uint32,uint32,int32){
	if user.role == nil {
		return 0,0, ErrorCode.RoleIsNil
	}
	return user.role.GetAwardDetails(id)
}

func (r *_Role) GetAwardDetails(id uint32) (uint32,uint32,int32){
	config,ok := Data.GetSupplyConfig().SupplyBox_ConfigItems[id]
	if ok == false{
		return 0,0,ErrorCode.Failure
	}
	tmpDrop := &drop.Drop{}
	ok, items := tmpDrop.Drop(config.DropID, 0, 0)
	if !ok {
		return 0,0, ErrorCode.Failure
	}
	return 0, uint32(len(items)), ErrorCode.Success
}

func (user *_User) RPC_RetSeasonEndAwardNotify() int32{
	if user.role == nil {
		return ErrorCode.RoleIsNil
	}
	user.role.prop.SyncSetSeasonAwardNotify(false)
	return ErrorCode.Success
}

func (r *_Role) GetSeasonRankData() *Proto.RankMemberData{
	data := &Proto.RankMemberData{}
	data.Name = r.prop.Data.Base.Name
	data.ID = r.GetID()
	data.Level = r.prop.Data.Base.Level
	data.Head = r.prop.Data.Base.Head
	return data
}

func (user *_User) RPC_GetSeasonLevelAward(level uint32) int32 {
	if user.role == nil {
		return  ErrorCode.RoleIsNil
	}
	result,_ := user.role.GetSeasonLevelAward(level)
	return result
}

func (r _Role) GetSeasonLevelAward(level uint32) (int32,uint32) {
	has,ok := r.prop.Data.SeasonInfo.LevelAwards[level]
	if ok == true && has == true{
		return ErrorCode.Failure,0
	}

	team := r.GetSpace().(*_Team)

	my_level := team.GetSeasonLevelByScore(team.prop.Data.SeasonInfo.SeasonID, r.prop.Data.SeasonInfo.MyScore)
	if level > my_level {
		return ErrorCode.Failure,0
	}
	config := team.GetSeasonLevelConfig(team.prop.Data.SeasonInfo.SeasonID, level)
	if config == nil {
		return ErrorCode.Failure,0
	}
	all := uint32(0)
	for _,v := range config.SeasonStageReward {
		list := strings.Split(v, ",")
		id,_ := strconv.Atoi(list[0])
		num,_ := strconv.Atoi(list[1])
		all += uint32(num)
		r.Infof("赛季获得等级奖励:level=%d,补给箱子:%d,%d", level, id, num)
		r.prop.SyncAddSupplyNum(uint32(id), uint32(num))
	}
	r.prop.SyncGetSeasonLevelAward(level)
	return ErrorCode.Success, all
}