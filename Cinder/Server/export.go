package Server

func Init(srvType string, areaID string, serverID string, rpcProc interface{}) error {
	return _Init(srvType, areaID, serverID, rpcProc)
}
func Destroy() {
	_Destroy()
}
