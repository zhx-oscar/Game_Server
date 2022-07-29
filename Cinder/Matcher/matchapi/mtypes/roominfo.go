package mtypes

import (
	"Cinder/Base/Util"
	"time"
)

// RoomInfo 匹配房间信息。
type RoomInfo struct {
	RoomID     RoomID    // 房间ID
	CreateTime time.Time // 创建时间
	MatchMode  MatchMode // 匹配模式

	// 有的以组队方式加入，有的以角色方式加入。
	// 以角色方式加入 RoleInfo.TeamID 为空；组队的 RoleInfo.TeamID 不为空。
	// 一旦加入房间，TeamID 不变。
	Roles RoleMap

	// 以上字段不可以在房间事件处理函数中更改
	// 以下字段可以在房间事件处理函数中更改

	IsHidden   bool // 是否隐藏, 隐藏后无法列举
	IsFull     bool // 是否人已满
	IsDeleting bool // 正在删除，不可见，不可加人. 可在事件处理函数中设之为true来删除房间。

	// 自定义数据.
	Data M // map[string]interface{}
}

// CreateRoomInfo 创建房间信息。
type CreateRoomInfo struct {
	MatchMode MatchMode // 匹配模式
	IsHidden  bool      // 是否隐藏, 隐藏后无法列举
	// 初始自定义数据.
	Data M // map[string]interface{}
}

type RoomID string

func NewRoomInfo(createRoomInfo CreateRoomInfo) *RoomInfo {
	roomData := createRoomInfo.Data
	if roomData == nil {
		roomData = make(M) // 保证非nil
	}
	return &RoomInfo{
		RoomID:     RoomID(Util.GetGUID()),
		CreateTime: time.Now(),
		MatchMode:  createRoomInfo.MatchMode,
		Data:       roomData, // 保证非 nil
		IsHidden:   createRoomInfo.IsHidden,
		Roles:      make(RoleMap), // 保证非 nil
	}
}

func (r *RoomInfo) GetData(key string) (interface{}, bool) {
	if r.Data == nil {
		return nil, false
	}
	result, ok := r.Data[key]
	return result, ok
}

func (r *RoomInfo) SetData(key string, val interface{}) {
	if r.Data != nil {
		r.Data[key] = val
		return
	}
	r.Data = map[string]interface{}{key: val}
}

func (r *RoomInfo) GetSrvIDToRoleIDs() SrvIDToRoleIDs {
	return r.Roles.GetSrvIDToRoleIDs()
}

func (r *RoomInfo) Copy() RoomInfo {
	cpy := *r
	cpy.Roles = r.Roles.Copy()
	cpy.Data = r.Data.Copy()
	return cpy
}
