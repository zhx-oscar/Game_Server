package main

import (
	"Daisy/Battle/drop"
	"Daisy/Data"
	"Daisy/Proto"
	"errors"
)

//Drop 根据当前上阵的特工的类型获取掉落
func (r *_Role) Drop(dropID uint32) ([]*Proto.DropMaterial, error) {

	fightingBuild, ok := r.prop.Data.BuildMap[r.prop.Data.FightingBuildID]
	if !ok {
		return nil, errors.New("找不到出战build")
	}

	specialAgent, ok := r.prop.Data.SpecialAgentList[fightingBuild.SpecialAgentID]
	if !ok {
		return nil, errors.New("找不到出战特工")
	}

	specialAgentCfg, ok := Data.GetSpecialAgentConfig().SpecialAgent_ConfigItems[specialAgent.Base.ConfigID]
	if !ok {
		return nil, errors.New("找不到出战特工的静态配置")
	}

	propCfg, ok := Data.GetPropConfig().PropValue_ConfigItems[specialAgentCfg.PropValueID]
	if !ok {
		return nil, errors.New("找不到出战特工的prop配置")
	}

	tmpDrop := &drop.Drop{}
	ok, items := tmpDrop.Drop(dropID, propCfg.Type, specialAgent.Base.Level)
	if !ok {
		return nil, errors.New("生成掉落品失败")
	}

	return items, nil
}

func (r *_Role) DropMaterialsToItems(materials []*Proto.DropMaterial) *Proto.Items {
	items := &Proto.Items{}
	tmpDrop := &drop.Drop{}
	for _, value := range materials {
		items.Items = append(items.Items, tmpDrop.CreateItem(value.MaterialId, value.MaterialType, value.MaterialNum))
	}

	return items
}
