# 压测

## 编译

编译压测工具 ChatStress:
```
cd D:\Daisy\Cinder_Server_Cinder\Chat\ChatStress
go build -v -o D:\Daisy\Server\bin
```

`Base/MQNet/service.go`中 nsq lookupd 查询间隔需要从 10s 改成 1s, 
不然启动时未能连上 nsq 会造成 RPC 3s 超时。
```
func (srv *_Service) initConsumer(addr string) error {
	cfg.LookupdPollInterval = 1 * time.Second
	...
}
```

## 运行

ChatStress 运行前，先在 bin 目录下运行 Chat。
注意，Chat 不要 debug 版，不要加 -race 检测，禁止 Debug 日志输出，不然数据会差别很大。
同时需开启 nsq 和 etcd 服务。

在 bin 目录下运行 ChatStress:
```
D:\Daisy\Server\bin (master -> origin)
λ ChatStress --goroutines 400 --test loginLogout
1587897344675272700 [Info] Logdir:../log/ Loglevel:debug
2020-04-26 18:35:44.835 ERR [logger.go:13] ERR    2 [chatstress_5ea56400acd6174df8f36c38#ephemeral/channel_chatstress_chatstress_5ea56400acd6174df8f36c38#ephemeral] error querying nsqlookupd (http://127.0.0.1:4161/lookup?topic=chatstress_5ea56400acd6174df8f36c38%23ephemeral) - got response 404 Not Found "{\"message\":\"TOPIC_NOT_FOUND\"}"
2020-04-26 18:35:44.838 ERR [logger.go:13] ERR    2 [chatstress_5ea56400acd6174df8f36c38#ephemeral/channel_chatstress_chatstress_5ea56400acd6174df8f36c38#ephemeral] error querying nsqlookupd (http://127.0.0.1:4161/lookup?topic=chatstress_5ea56400acd6174df8f36c38%23ephemeral) - got response 404 Not Found "{\"message\":\"TOPIC_NOT_FOUND\"}"
2020-04-26 18:35:44.841 ERR [logger.go:13] ERR    2 [chatstress_5ea56400acd6174df8f36c38#ephemeral/channel_chatstress_chatstress_5ea56400acd6174df8f36c38#ephemeral] error querying nsqlookupd (http://127.0.0.1:4161/lookup?topic=chatstress_5ea56400acd6174df8f36c38%23ephemeral) - got response 404 Not Found "{\"message\":\"TOPIC_NOT_FOUND\"}"
2020-04-26 18:35:44.849 ERR [logger.go:13] ERR    3 [chatstress#ephemeral/channel_chatstress_chatstress_5ea56400acd6174df8f36c38#ephemeral] error querying nsqlookupd (http://127.0.0.1:4161/lookup?topic=chatstress%23ephemeral) - got response 404 Not Found "{\"message\":\"TOPIC_NOT_FOUND\"}"
2020-04-26 18:35:44.852 ERR [logger.go:13] ERR    3 [chatstress#ephemeral/channel_chatstress_chatstress_5ea56400acd6174df8f36c38#ephemeral] error querying nsqlookupd (http://127.0.0.1:4161/lookup?topic=chatstress%23ephemeral) - got response 404 Not Found "{\"message\":\"TOPIC_NOT_FOUND\"}"
2020-04-26 18:35:44.855 ERR [logger.go:13] ERR    3 [chatstress#ephemeral/channel_chatstress_chatstress_5ea56400acd6174df8f36c38#ephemeral] error querying nsqlookupd (http://127.0.0.1:4161/lookup?topic=chatstress%23ephemeral) - got response 404 Not Found "{\"message\":\"TOPIC_NOT_FOUND\"}"
chat stress test is running (goroutines=400 test=loginLogout time=15s)...
loginLogout TPS: 3884 / 15.000000 = 258.933333 (t/s)
```

参数说明
* goroutines 是同时请求的协程数。每个协程不断发送请求，最后计算每秒事务数(TPS)输出。
* test 测试名。请用 `ChatStress.exe --list` 列出所有 test 名。非法名字会报错退出。
* time 测试时长。
* list 列出所有测试名。

```
λ ChatStress -h
1587897418552929600 [Info] Logdir:../log/ Loglevel:debug
Usage of ChatStress:
  -g, --goroutines int   goroutines count (default 1)
  -l, --list             list all test names
  -t, --test string      test case name (default "loginLogout")
  -m, --time duration    test time (default 15s)
pflag: help requested
```

测试名：
```
λ ChatStress --list
1588128983860724000 [Info] Logdir:../log/ Loglevel:debug
Test names:
        addFriendReqApply
        addFriendToBlacklistAndRemove
        addIntoGroupAndKick
        createDeleteGroup
        createSameGroup
        deleteFriend
        followUnfollowFriendReq
        getFollowerList
        getFollowingList
        getFriendBlacklist
        getFriendInfos
        getFriendInfosNil
        getFriendList
        getGroupMemberCounts
        getGroupMembers
        getHistoryMessages
        getOfflineMessage
        loginLogout
        loginLogoutSame
        loginSame
        sendGroupMessage0
        sendGroupMessage1
        sendGroupMessage10
        sendGroupMessage100
        sendGroupMessage2
        sendGroupMessage3
        sendMessage
```

### 压测结果

ChatStress 输出每秒事务数(TPS). 
多数测试的事务为单次请求，个别测试一个事务为2次请求，如 loginLogout。
如果同时在线 100K, 平均每人 100s 请求一次匹配，就需要服务器能提供大于 1K 的 RPS(每秒请求数).

本机测试环境:
* Processor: Intel Core i5-7400 CPU @ 3.00GHz
* RAM: 8 GB

多数接口当 goroutines 大概在 1000 会到达 RPS 的最大值，此时 CPU 满载，利用率高。
继续加大 goroutines 会因为协程消耗使 RPS 略微下降, 并且可能会出现超时错误。
测试时长为 15s.

各接口测试结果:

test名							|goroutines |RPS 	|说明
--------------------------------|-----------|-------|------------------------------------------------------------------
addFriendReqApply               |500        |1522*2 |
addFriendToBlacklistAndRemove   |500        |1755*2 |
addIntoGroupAndKick             |2000       |800*2  |
createDeleteGroup               |500        |972*2  |
createSameGroup                 |3000       |9512   |
deleteFriend                    |1500       |8034   |
followUnfollowFriendReq         |1000       |533*2  |
getFollowerList                 |1500       |7778   |
getFollowingList                |1000       |6992   |
getFriendBlacklist              |1000       |7222   |
getFriendInfos                  |500        |2878   |一次获取10个ID
getFriendInfosNil               |1000       |7476   |获取ID列表为空
getFriendList                   |1500       |7469   |
getGroupMemberCounts            |500        |3929   |一次获取10个群，每个群都是空的
getGroupMembers                 |1000       |7656   |
getHistoryMessages              |1000       |7637   |
getOfflineMessage               |1000       |1530   |
loginLogout                     |1000       |520*2  |
loginLogoutSame                 |1000       |896*2  |
loginSame                       |1000       |7538   |
sendGroupMessage0               |1000       |7893   |群内无人在线,忽略消息
sendGroupMessage1               |1000       |1729   |群内1人在线
sendGroupMessage2               |500        |1705   |群内2人在线，同一服，--time=120s
sendGroupMessage3               |500        |1722   |群内3人在线，同一服，--time=120s
sendGroupMessage10              |500        |1571   |群内10人在线，同一服，--time=120s
sendGroupMessage100             |100        |929    |群内100人在线，同一服，--time=120s
sendMessage                     |1000       |5215   |

有些接口主要是 mongoDB 操作，CPU较空，需要大一点的 goroutines 才能 CPU 占满。
