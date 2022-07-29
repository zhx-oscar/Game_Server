package matchapi

import (
	"Cinder/Matcher/matchapi/internal"
	"Cinder/Matcher/matchapi/mtypes"
)

// ITeamService 是队伍管理接口
// 匹配前先组队可以保证成员会匹配在同一战斗中，而且一般会在同一阵营。
// 匹配组队与战斗中的队伍(阵营)不同，如30个匹配队(每队1..n个人)匹配成一个50v50的对战，实际战斗队伍分成了2个。
// 进入房间后，队伍变化不会改变房间内的组队。
type ITeamService interface {
	// 创建队伍
	CreateTeam(creatorInfo mtypes.RoleInfo, teamInfo mtypes.TeamInfo) (mtypes.TeamInfo, error)
	// 加入队伍，触发队伍广播 JoinTeamMsg
	JoinTeam(roleInfo mtypes.RoleInfo, teamID mtypes.TeamID, passwd string) (mtypes.TeamInfo, error)
	// 离开队伍, 或踢除队员, 触发队伍广播 LeaveTeamMsg
	LeaveTeam(memberID mtypes.RoleID, teamID mtypes.TeamID) error

	// 变更队长, 触发队伍广播 ChangeTeamLeaderMsg
	ChangeTeamLeader(teamID mtypes.TeamID, newLeader mtypes.RoleID) error
	// 设置队伍数据, 触发队伍广播 SetTeamDataMsg
	SetTeamData(teamID mtypes.TeamID, key string, data interface{}) error

	// 队伍广播。触发队伍广播 BroadcastTeamMsg
	BroadcastTeam(teamID mtypes.TeamID, msg interface{}) error
}

// IRoomService 是匹配房间管理接口
// 匹配即创建一个房间或加入一个房间，同时只能在一个房间。
// 可以先组队, 再创建或加入房间。
type IRoomService interface {
	// 加入随机房间, 触发房间广播 JoinRoomMsg
	RoleJoinRandomRoom(roleInfo mtypes.RoleInfo, matchMode mtypes.MatchMode) (mtypes.RoomInfo, error)
	// 加入随机房间, 带初始化信息, 触发房间广播 JoinRoomMsg
	RoleJoinRandomRoomWithInfo(roleInfo mtypes.RoleInfo, createInfo mtypes.CreateRoomInfo) (mtypes.RoomInfo, error)
	// 创建房间，触发房间广播 JoinRoomMsg
	RoleCreateRoom(creatorInfo mtypes.RoleInfo, matchMode mtypes.MatchMode) (mtypes.RoomInfo, error)
	// 创建房间, 带初始化信息，触发房间广播 JoinRoomMsg
	RoleCreateRoomWithInfo(creatorInfo mtypes.RoleInfo, createInfo mtypes.CreateRoomInfo) (mtypes.RoomInfo, error)

	// 加入指定房间, 触发房间广播 JoinRoomMsg
	RoleJoinRoom(roleInfo mtypes.RoleInfo, roomID mtypes.RoomID) (mtypes.RoomInfo, error)
	// 自己离开房间, 或踢人, 触发房间广播 LeaveRoomMsg
	RoleLeaveRoom(roleID mtypes.RoleID, roomID mtypes.RoomID) (mtypes.RoomInfo, error)

	// 加入随机房间, 触发房间广播 JoinRoomMsg
	TeamJoinRandomRoom(teamID mtypes.TeamID, matchMode mtypes.MatchMode) (mtypes.RoomInfo, error)
	// 加入随机房间, 带初始化信息, 触发房间广播 JoinRoomMsg
	TeamJoinRandomRoomWithInfo(teamID mtypes.TeamID, createInfo mtypes.CreateRoomInfo) (mtypes.RoomInfo, error)
	// 创建房间，触发房间广播 JoinRoomMsg
	TeamCreateRoom(creatorTeamID mtypes.TeamID, matchMode mtypes.MatchMode) (mtypes.RoomInfo, error)
	// 创建房间, 带初始化信息，触发房间广播 JoinRoomMsg
	TeamCreateRoomWithInfo(creatorTeamID mtypes.TeamID, createInfo mtypes.CreateRoomInfo) (mtypes.RoomInfo, error)
	// 加入指定房间, 触发房间广播 JoinRoomMsg
	TeamJoinRoom(teamID mtypes.TeamID, roomID mtypes.RoomID) (mtypes.RoomInfo, error)

	// 房间广播。触发房间广播 BroadcastRoomMsg
	BroadcastRoom(roomID mtypes.RoomID, msg interface{}) error

	// 获取房间列表，不可见房间不会列出。批量获取多个模式请用 ListRooms()。
	GetRoomList(matchMode mtypes.MatchMode) (roomIDs []mtypes.RoomID, err error)
	// 获取房间详情。批量获取多个房间请用 GetRoomInfos()
	GetRoomInfo(roomID mtypes.RoomID) (mtypes.RoomInfo, error)
	// 获取房间列表，不可见房间不会列出
	ListRooms(matchModes []mtypes.MatchMode) (rooms map[mtypes.MatchMode][]mtypes.RoomInfo, err error)
	// 获取房间详情
	GetRoomInfos(roomIDs []mtypes.RoomID) (map[mtypes.RoomID]mtypes.RoomInfo, error)

	// 设置房间数据, 触发房间内广播 SetRoomDataMsg
	SetRoomData(roomID mtypes.RoomID, key string, data interface{}) error
	// 更新房间内角色数据, 触发房间内广播 UpdateRoomRoleDataMsg
	SetRoomRoleFloatData(roomID mtypes.RoomID, roleID mtypes.RoleID, key string, data float64) error
	SetRoomRoleStringData(roomID mtypes.RoomID, roleID mtypes.RoleID, key string, data string) error
	AddRoomRoleTag(roomID mtypes.RoomID, roleID mtypes.RoleID, tag string) error
	DelRoomRoleTag(roomID mtypes.RoomID, roleID mtypes.RoleID, tag string) error
	// 删除房间, 触发所有人离开房间
	DeleteRoom(roomID mtypes.RoomID) error
}

// 获取本区组队服务
func GetTeamService() ITeamService {
	return internal.NewTeamService("", "")
}

// 获取本区房间匹配服务
func GetRoomService() IRoomService {
	return internal.NewRoomService("", "")
}

// 获取跨区组队服务
// 需要指定跨区匹配服的区号(AreaID)和服务器号
func GetGlobalTeamService(matcherAreaID string, matcherServerID string) ITeamService {
	return internal.NewTeamService(matcherAreaID, matcherServerID)
}

// 获取跨区房间匹配服务
// 需要指定跨区匹配服的区号(AreaID)和服务器号
func GetGlobalRoomService(matcherAreaID string, matcherServerID string) IRoomService {
	return internal.NewRoomService(matcherAreaID, matcherServerID)
}
