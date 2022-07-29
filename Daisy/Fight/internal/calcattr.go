package internal

import (
	. "Daisy/Fight/attraffix"
	"Daisy/Fight/internal/conf"
	"Daisy/Proto"
	"fmt"
)

// CalcAttr 属性计算器
type CalcAttr struct {
	FightAttr
}

// Init 初始化
func (calcAttr *CalcAttr) Init(pawnType Proto.PawnType_Enum, configID uint32, level int32, attrAffixList []AttrAffix) error {
	// 查询pawn配置
	pawnConf, ok := conf.GetConfigs().GetPawnConfig(pawnType, configID)
	if !ok {
		return fmt.Errorf("not found pawnType %d configID %d config", pawnType, configID)
	}

	// 初始化属性
	return calcAttr.initAttr(pawnConf, level, attrAffixList)
}

// AddAttrAffix 添加属性词缀
func (calcAttr *CalcAttr) AddAttrAffix(attrAffixList []AttrAffix) {
	for _, v := range attrAffixList {
		calcAttr.ChangeAttr(v.Field, v.ParaA, v.ParaB, true)
	}
}

// RemoveAttrAffix 删除属性词缀
func (calcAttr *CalcAttr) RemoveAttrAffix(attrAffixList []AttrAffix) {
	for _, v := range attrAffixList {
		calcAttr.ChangeAttr(v.Field, v.ParaA, v.ParaB, false)
	}
}
