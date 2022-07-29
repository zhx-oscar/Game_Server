package Prop

import (
	"Cinder/Base/Core"
)

const (
	RolePropType    = "RoleProp"
	TeamPropType    = "TeamProp"
	MailBoxPropType = "MailProp"
)

func RegisterPropType() {
	Core.Inst.RegisterProp(RolePropType, &RoleProp{})
	Core.Inst.RegisterProp(TeamPropType, &TeamProp{})

	Core.Inst.RegisterPropObject(OfflineRoleObject, &OfflineRole{})
	Core.Inst.RegisterPropObject(OfflineTeamObject, &OfflineTeam{})
}
