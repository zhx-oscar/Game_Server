package ltypes

import (
	"Cinder/Matcher/matchapi/mtypes"
)

// IRoomEventHandler 房间事件处理器.
// 可以在各种事件通知中执行动作，修改 roomInfo.
// RoomInfo 中的以下字段禁止更改：RoomID, CreateTime, MatchMode, Roles
type IRoomEventHandler interface {
	// 创建房间。可向 Space 请求创建场景。
	OnCreate(roomInfo *mtypes.RoomInfo) error
	// 销毁房间。可进行清理动作。
	OnDestroy(roomInfo *mtypes.RoomInfo)

	// 加人前判断是否可加，即执行匹配算法
	OnAddingRoles(roomInfo *mtypes.RoomInfo, roles mtypes.RoleMap) bool
	// 加人后动作
	OnAddedRoles(roomInfo *mtypes.RoomInfo, roles mtypes.RoleMap)
	// 删人前动作
	OnDeletingRole(roomInfo *mtypes.RoomInfo, roleID mtypes.RoleID)
}
