package Prop

import (
	"Cinder/Base/Prop"
)

// 离线队伍处理对象

const OfflineTeamObject = "OfflineTeamObject"

type OfflineTeam struct {
	Prop.Object
	*TeamProp
}

func (t *OfflineTeam) Init() {
	t.TeamProp = t.GetProp().(*TeamProp)
}

func (t *OfflineTeam) GetPropInfo() (string, string) {
	return TeamPropType, t.GetID()
}
