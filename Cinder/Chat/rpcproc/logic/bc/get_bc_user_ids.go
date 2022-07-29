package bc

import (
	"Cinder/Chat/rpcproc/logic/types"

	log "github.com/cihub/seelog"
)

// getBcUserIDs 获取需要广播的目标ID: 好友和群成员。
// 其中好友中包含了离线的。
// 排除自身。
func getBcUserIDs(userID types.UserID) userIDs {
	u := userMgr.GetUser(userID)
	if u == nil {
		return nil
	}

	groupIDs := u.GetUserGroupMgr().GetGroupIDs()
	ret := getGroupsOnlineMembers(groupIDs)

	friendIDs := u.GetFriendMgr().GetFriendIDs()
	for _, id := range friendIDs {
		ret[id] = true
	}
	delete(ret, userID)
	return ret
}

// getGroupsOnlineMembers 获取多个群的在线成员
func getGroupsOnlineMembers(groupIDs []types.GroupID) userIDs {
	ret := userIDs{}
	for _, groupID := range groupIDs {
		group, groupErr := groupMgr.GetOrLoadGroup(groupID)
		if groupErr != nil {
			log.Errorf("GetOrLoadGroup(%d): %s", groupID, groupErr)
			continue
		}

		members := group.CopyOnlineMemberIDs()
		for id, _ := range members {
			ret[id] = true
		}
	}
	return ret
}
