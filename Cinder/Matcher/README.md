# 匹配服

## 匹配服功能

匹配服有2个功能：组队和匹配。
详细需求说明见：[doc/requirement.md](doc/requirement.md)

### 组队

组队功能有：创建队伍，加入队伍，离开队伍，变更队长, 设置队伍数据, 队伍广播。

匹配前先组队可以保证成员会匹配在同一战斗中，而且一般会在同一阵营。

### 匹配

匹配即创建一个房间或加入一个房间。
功能有：加入随机房间，创建房间，加入指定房间，离开房间，更改房间Leader，房间广播，设置房间数据，设置角色数据。
可以先组队, 再创建或加入房间。

匹配服需要自定义匹配算法来决定一个或一队角色是否可以加入一个房间。

## 如何创建匹配服

匹配服需要自定义匹配算法和回调接口 `IRoomEventHandler`, 应用 matcherlib 创建匹配服应用。
参见示例：MatcherExample.

### 初始化和退出
```
func main() {
	// ...
	
	if err := matcherlib.Init(serverID, &RPCProc{}, &RoomEventHandler{}); err != nil {
		log.Error("start server failed ", err)
		return
	}
	defer matcherlib.Destroy()
	
	// ...
}
```

### `RPCProc`

`RPCProc`可以直接使用 `matcherlib.RPCProc`:
```
matcherlib.Init(serverID, &matcherlib.RPCProc{}, ...)
```
也可以自定义，并内嵌 `matcherlib.RPCProc`:
```
type RPCProc struct {
	matcherlib.RPCProc
}

func (r *RPCProc) RPC_MyRPCTest(arg string) int {
	return 0
}

func main() {
	...
    matcherlib.Init(serverID, &RPCProc{}, ...)
    ...
}
```

### 匹配算法和事件处理

`matcherlib.Init` 还需要一个 `IRoomEventHandler` 参数，该接口需要应用自定义。
`IRoomEventHandler` 是房间事件处理器，当房间有加人，减人等变化时，会调用相应的事件处理函数，
可以在各种事件处理函数中执行动作，读取并修改输入的房间信息。

匹配算法可实现于 `OnAddingRoles`, 在加人前判断是否可加。

示例：
```
func (r *RoomEventHandler) OnAddingRoles(roomInfo *mtypes.RoomInfo, roles mtypes.RoleMap) bool {
	if len(roomInfo.Roles)+len(roles) > getMaxRoleCount(roomInfo.MatchMode) {
		return false // 人已满
	}
	maxLevel := float64(getMaxRoleLevel(roomInfo.MatchMode))
	for _, role := range roles {
		level, ok := role.GetFloatData("level")
		if !ok || level > maxLevel {
			return false // 等级太高了，不允许参加
		}
	}
	return true
}
```

## 如何接入匹配服

Game服调用匹配服功能可使用 matchapi 包。使用说明见：[matchapi/match_usage.md](matchapi/match_usage.md)

## 目录说明

名字				| 说明
----------------|---------------------------------------------------------------
doc				| 文档
matchapi		| 匹配服API
MatcherExample	| 匹配服示例
MatcherStress   | 匹配压测工具
matchlib		| 匹配服库，用于实现匹配服
rpcmsg			| matchapi 调用匹配服 RPC 所用的消息

