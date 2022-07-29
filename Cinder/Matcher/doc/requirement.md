# 需求说明

匹配服实现以下功能：

* 暂为单服匹配，不支持多个匹配服
* 组队功能
* 撮合多个队伍为一场战斗
* 撮合多人为一场战斗
* 匹配算法需要上层自定义
* 可以单个角色(Role)请求匹配，也可以组队(Team)后请求匹配
* 组队后请求匹配，保证匹配在同一场战斗中
	
## 拳皇项目的需求

只需要匹配接口，不需要组队接口。

* SRPC_SetRoomConfig(uid uint64, key string, v int32)
	+ 设置房间属性
	+ 改为 SetRoomData
* SRPC_OnPlayerOutline(uid uint64)
	+ 玩家下线
	+ 改为 RoleLeaveRoom
* SRPC_QuickJoin(stageID uint32, uid uint64, roleBase *protodef.RoleBase, userProxy *entity.EntityProxy)
	+ 快速加入
	+ 改为 RoleJoinRandomRoom
* SRPC_LeaveTeam(uid uint64, userProxy *entity.EntityProxy)
	+ 退出
	+ 改为 RoleLeaveRoom
* SRPC_KickTeam(uid uint64, targetId uint64, userProxy *entity.EntityProxy)
	+ 踢人
	+ 改为 RoleLeaveRoom
* SRPC_ChangeTeamLeader(uid uint64, targetId uint64, userProxy *entity.EntityProxy)
	+ 更改房间主人
	+ 改为 ChangeRoomLeader
* SRPC_ChangeSlotIndex(uid uint64, slotIndex int32, userProxy *entity.EntityProxy)
	+ 更改槽号
	+ 改为更改角色属性 SetRoomRoleFloatData
* SRPC_EnterRoom(uid uint64, userProxy *entity.EntityProxy)
	+ 进入
	+ 改为立即开始，人未满立即开始 StartRoomNow, 无此功能？
* SRPC_JoinPVP(stageID uint32, uid uint64, roleBase *protodef.RoleBase, userProxy *entity.EntityProxy)
	+ 加入
	+ 改为 RoleJoinRandomRoom
* SRPC_CancelJoinPVP(uid uint64, userProxy *entity.EntityProxy)
	+ 取消
	+ 改为取消匹配 RoleLeaveRoom
	