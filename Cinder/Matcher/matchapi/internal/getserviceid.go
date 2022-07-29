package internal

import (
	"Cinder/Base/Core"
	"Cinder/Matcher/matchapi/mtypes"
)

func getServiceID() mtypes.SrvID {
	return mtypes.SrvID(Core.Inst.GetServiceID())
}
