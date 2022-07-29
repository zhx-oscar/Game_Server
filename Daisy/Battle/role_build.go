package main

import (
	"Cinder/Base/Util"
	"Daisy/Const"
	"Daisy/Data"
	"Daisy/ErrorCode"
	"Daisy/ItemProto"
	"Daisy/Proto"
	"time"
)

//recalculationAllBuildAttr 重新计算build属性
func (r *_Role) recalculationAllBuildAttr() {
	for buildID := range r.prop.Data.BuildMap {
		r.refreshBuildFightAttr(buildID)
	}
}

//newSpecialAgent 新建特工数据
func (r *_Role) newSpecialAgent(specialAgentID uint32) *Proto.SpecialAgent {
	specialAgent := &Proto.SpecialAgent{
		Base: &Proto.SpecialAgentBase{
			ConfigID: specialAgentID,
			Level:    1,
			Exp:      0,
			GainTime: time.Now().Unix(),
		},
		Talent: r.initTalent(specialAgentID),
	}

	return specialAgent
}

//newBuild 新建build数据
func (r *_Role) newBuild(name string, specialAgentID uint32) *Proto.BuildData {
	return &Proto.BuildData{
		BuildID:        Util.GetGUID(),
		Name:           name,
		SpecialAgentID: specialAgentID,
		Skill: &Proto.BuildSkillData{
			UltimateSkillID: 0,
			SuperSkill:      map[uint32]uint32{},
		},
		EquipmentMap: map[int32]string{},
		CreateTime:   time.Now().Unix(),
		FightAttr:    &Proto.FightAttr{},
	}
}

//RPC_CreateBuild RPC创建build
func (user *_User) RPC_CreateBuild(name string, specialAgentID uint32) int32 {
	if user.role == nil {
		return ErrorCode.RoleIsNil
	}

	return user.role.createBuild(name, specialAgentID)
}

//createBuild 创建build
func (r *_Role) createBuild(name string, specialAgentID uint32) int32 {
	buildCountMaxConf, ok := Data.GetSpecialAgentConfig().SpecialAgentConst_ConfigItems[Const.SpecialAgent_buildCountMax]
	if !ok {
		return ErrorCode.BuildCountMaxNotFind
	}

	//build列表上限判断
	if len(r.prop.Data.BuildMap) >= int(buildCountMaxConf.Value) {
		return ErrorCode.OverLimitBuildCountMaxNotFind
	}

	//特工是否拥有
	if _, find := r.prop.Data.SpecialAgentList[specialAgentID]; !find {
		return ErrorCode.NoSpecialAgent
	}

	//build新建
	build := r.newBuild(name, specialAgentID)
	r.prop.SyncAddBuild(build)
	r.refreshBuildFightAttr(build.BuildID)
	return ErrorCode.Success
}

//RPC_FightingBuild 上阵作战build
func (user *_User) RPC_FightingBuild(buildID string) int32 {
	if user.role == nil {
		return ErrorCode.RoleIsNil
	}

	return user.role.fightingBuild(buildID)
}

//fightingBuild 作战build
func (r *_Role) fightingBuild(buildID string) int32 {
	_, ok := r.prop.Data.BuildMap[buildID]
	if !ok {
		return ErrorCode.NotFindBuild
	}

	r.prop.SyncSetFightingBuildID(buildID)
	r.FlushToDB()
	r.FlushToCache()
	return ErrorCode.Success
}

//RPC_ChangeBuildName 修改build名字
func (user *_User) RPC_ChangeBuildName(buildID, name string) int32 {
	if user.role == nil {
		return ErrorCode.RoleIsNil
	}

	return user.role.changeBuildName(buildID, name)
}

//changeBuildName 设置build名字
func (r *_Role) changeBuildName(buildID, name string) int32 {
	_, ok := r.prop.Data.BuildMap[buildID]
	if !ok {
		return ErrorCode.NotFindBuild
	}

	r.prop.SyncSetBuildName(buildID, name)
	return ErrorCode.Success
}

//itemAddInUse 道具引用计数 ++ 处理
func (r *_Role) itemAddInUse(item ItemProto.IItem, buildID string) {
	if item == nil {
		return
	}

	item.GetData().ExpandData.InUse++
	item.GetData().ExpandData.RelationBuildList[buildID] = true
	r.UpdataItem(item)
}

//itemDelInUse 道具引用计数 -- 处理
func (r *_Role) itemDelInUse(item ItemProto.IItem, buildID string) {
	if item == nil {
		return
	}

	item.GetData().ExpandData.InUse--
	delete(item.GetData().ExpandData.RelationBuildList, buildID)
	r.UpdataItem(item)
}
