package main

import (
	"Cinder/Matcher/matchapi/mtypes"

	log "github.com/cihub/seelog"
)

type RoomEventHandler struct {
}

func (r *RoomEventHandler) OnCreate(roomInfo *mtypes.RoomInfo) error {
	log.Debugf("OnCreate, RoomID=%s", roomInfo.RoomID)
	// 一般应该向Space申请开房间，将返回的房间信息记录在 roomInfo.Data 某个字段中
	return nil
}

func (r *RoomEventHandler) OnDestroy(roomInfo *mtypes.RoomInfo) {
	log.Debugf("OnDestroy, RoomID=%s", roomInfo.RoomID)
}

func (r *RoomEventHandler) OnAddingRoles(roomInfo *mtypes.RoomInfo, roles mtypes.RoleMap) bool {
	log.Debugf("OnAddingRoles: RoomID=%s, Roles=%d", roomInfo.RoomID, len(roles))
	if len(roomInfo.Roles)+len(roles) > getMaxRoleCount(roomInfo.MatchMode) {
		return false // 人已满
	}
	maxLevel := float64(getMaxRoleLevel(roomInfo.MatchMode))
	for _, role := range roles {
		level, ok := role.GetFloatData("level")
		if !ok {
			continue // 没有等级就忽略等级限制，用于测试
		}
		if level > maxLevel {
			return false // 等级太高了，不允许参加
		}
	}
	return true
}

func (r *RoomEventHandler) OnAddedRoles(roomInfo *mtypes.RoomInfo, roles mtypes.RoleMap) {
	log.Debugf("OnAddedRoles: RoomID=%s, Roles=%d", roomInfo.RoomID, len(roles))
	roomInfo.IsFull = (len(roomInfo.Roles) >= getMaxRoleCount(roomInfo.MatchMode))
	roomInfo.SetData("roleCount", len(roomInfo.Roles)) // 任意更新数据，用于后续匹配，或通知房间成员
	if roomInfo.IsFull {
		roomInfo.IsDeleting = true
	}
}

func (r *RoomEventHandler) OnDeletingRole(roomInfo *mtypes.RoomInfo, roleID mtypes.RoleID) {
	log.Debugf("OnDeletingRole: RoomID=%s, RoleID=%s", roomInfo.RoomID, roleID)
}

func getMaxRoleCount(matchMode mtypes.MatchMode) int {
	switch matchMode {
	case "1v1":
		return 2
	case "2v2":
		return 4
	}
	return 0
}

// 获取等级限制
func getMaxRoleLevel(matchMode mtypes.MatchMode) int {
	switch matchMode {
	case "1v1":
		return 200
	case "2v2":
		return 400
	}
	return 0
}
