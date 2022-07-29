package bc

import (
	"Cinder/Chat/chatapi/types"
	ltypes "Cinder/Chat/rpcproc/logic/types"

	assert "github.com/arl/assertgo"
)

type userIDs = map[ltypes.UserID]bool
type srvIDToTargetsMap map[string][]types.Target

// getSrvIDToTargetsMap 将 UserID 按 srvID 拆分
func getSrvIDToTargetsMap(ids userIDs) srvIDToTargetsMap {
	result := map[string][]types.Target{}
	assert.True(userMgr != nil) // 须初始化设置
	for userID, _ := range ids {
		user := userMgr.GetUser(userID)
		if user == nil {
			continue // 有可能并发，UserMgr中已下线，但Group中还未下线
		}

		target := types.Target{
			ID:   string(userID),
			Data: user.GetActiveData(),
		}
		srvID := user.GetSrvID()
		result[srvID] = append(result[srvID], target)
	}
	return result
}
