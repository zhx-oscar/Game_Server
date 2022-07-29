package DBAgent

import (
	"Cinder/Base/Core"
)

var Inst Core.ICore

func Init(areaID string, serverID string, rpcProc interface{}) error {
	return _Init(areaID, serverID, rpcProc)
}
func Destroy() {
	_Destroy()
}
