# 匹配服使用说明

匹配服功能包括组队和开房间，详见：[../doc/requirement.md](../doc/requirement.md)

请使用 Cinder/Match/matchapi 包接入匹配服。仅支持 go 语言并应用 Cinder 框架开发的游戏。

matchapi 一般由 Game 服调用，要求 Game 实现特定的一组 RPC 回调。

## 匹配步骤

可以个人匹配，也可以组队后匹配。组队后匹配可以保证成员进入同一战斗，并且一般在同一阵营。
个人匹配可以看成是先组成一个单人队伍后请求匹配。

请求匹配后，会根据匹配模式加入一个现有房间，或者创建一个新房间，等待其他玩家加入。
房间满时即完成匹配。
匹配过程中可以看到有其他玩家加入或退出。

匹配过程中，房间创建，加减人等会通过 RPC 通知房间内所有成员，也会调用 Notifier 接口通知匹配服的上层应用做自定义处理。

## Game 服回调

`RPC_MatchNotify(msgsJson []byte)`
msgsJson 是 `mtypes.NotifyMsgsToOneSrv` 的 json 打包, 解包及处理示例如下：

```
func (r *RPCProc) RPC_MatchNotify(msgJson []byte) {
	msg, errJson := mtypes.UnmarshalNotifyOneSrvMsg(msgJson)
	if errJson != nil {
		log.Error(errJson)
		return
	}
	// ...
}
```

## 跨区匹配

缺省为本区匹配。为了跨区匹配，需要一个专用的区运行一个跨区匹配服。
所有区公用一个 nsq 集群，但是有各自的 etcd 和 DB.

`GetGlobalTeamService(area, id)` 和 `GetGlobalRoomService(area, id)` 获取的是跨区服务，
需要指定跨区匹配服的区号(AreaID)和服务器号。
