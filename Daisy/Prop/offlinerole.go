package Prop

import (
	"Cinder/Base/Prop"
)

// 离线角色信息

const OfflineRoleObject = "OfflineRoleObject"

type OfflineRole struct {
	Prop.Object
	*RoleProp
}

func (r *OfflineRole) Init() {
	r.RoleProp = r.GetProp().(*RoleProp)
}

func (r *OfflineRole) GetPropInfo() (string, string) {
	return RolePropType, r.GetID()
}
