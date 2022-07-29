# 压测

## 编译

编译压测工具 MatcherStress:
```
cd D:\Daisy\Cinder_Server_Cinder\Matcher\MatcherStress
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

以下压测的对象为示例匹配服 MatcherExample。
MatchStress 运行前，先在 bin 目录下运行 MatcherExample。
注意，MatchStress 不要 debug 版，不要加 -race 检测，禁止 Debug 日志输出，不然数据会差别很大。
同时需开启 nsq 和 etcd 服务。

在 bin 目录下运行 MatchStress:
```
D:\Daisy\Server\bin (master -> origin)
λ MatcherStress.exe --goroutines 1500 --test RoleJoinRandomRoom
1587867068270399800 [Info] Logdir:../log/ Loglevel:debug
2020-04-26 10:11:08.337 ERR [logger.go:13] ERR    2 [matcherstress_5ea4edbcacd6174140b7c160/channel_matcherstress_matcherstress_5ea4edbcacd6174140b7c160] error querying nsqlookupd (http://127.0.0.1:4161/lookup?topic=matcherstress_5ea4edbcacd6174140b7c160) - got response 404 Not Found "{\"message\":\"TOPIC_NOT_FOUND\"}"
2020-04-26 10:11:08.34 ERR [logger.go:13] ERR    2 [matcherstress_5ea4edbcacd6174140b7c160/channel_matcherstress_matcherstress_5ea4edbcacd6174140b7c160] error querying nsqlookupd (http://127.0.0.1:4161/lookup?topic=matcherstress_5ea4edbcacd6174140b7c160) - got response 404 Not Found "{\"message\":\"TOPIC_NOT_FOUND\"}"
2020-04-26 10:11:08.342 ERR [logger.go:13] ERR    2 [matcherstress_5ea4edbcacd6174140b7c160/channel_matcherstress_matcherstress_5ea4edbcacd6174140b7c160] error querying nsqlookupd (http://127.0.0.1:4161/lookup?topic=matcherstress_5ea4edbcacd6174140b7c160) - got response 404 Not Found "{\"message\":\"TOPIC_NOT_FOUND\"}"
matcher stress test is running (goroutines=1500 test=RoleJoinRandomRoom time=15s)...
RoleJoinRandomRoom TPS: 17656 / 15.000000 = 1177.066667 (t/s)
```

参数说明
* goroutines 是同时请求的协程数。每个协程不断发送请求，最后计算每秒事务数(TPS)输出。
* test 测试名。请用 `MatcherStreee.exe --list` 列出所有 test 名。非法名字会报错退出。
* time 测试时长。
* list 列出所有测试名。

```
λ MatcherStress.exe -h
1587867304053182800 [Info] Logdir:../log/ Loglevel:debug
Usage of MatcherStress.exe:
  -g, --goroutines int   goroutines count (default 1)
  -l, --list             list all test names
  -t, --test string      test case name (default "RoleJoinRandomRoom")
  -m, --time duration    test time (default 15s)
pflag: help requested
```

测试名：
```
λ MatcherStress.exe --list
1587867767770790100 [Info] Logdir:../log/ Loglevel:debug
Test names:
        BroadcastRoom
        BroadcastTeam
        ChangeTeamLeader
        CreateLeaveTeam
        GetRoomInfo
        GetRoomList
        JoinLeaveTeam
        MyRPCTest
        Ping
        RoleCreateDeleteRoom
        RoleJoinLeaveRoom
        RoleJoinRandomRoom
        RoleJoinRoom
        SetRoomData
        SetRoomRoleData
        SetTeamData
        TeamCreateDeleteRoom
        TeamJoinRandomRoom
        TeamJoinRoom
```

### 压测结果

MatcherStress 输出每秒事务数(TPS). 
多数测试的事务为单次请求，个别测试一个事务为2次请求，如 RoleJoinLeaveRoom。
如果同时在线 100K, 平均每人 100s 请求一次匹配，就需要服务器能提供大于 1K 的 RPS(每秒请求数).

#### 本机测试环境
* Processor: Intel Core i5-7400 CPU @ 3.00GHz
* RAM: 8 GB

#### BenchmarkRPC_RoleJoinRandomRoom

直接调用 RPC_RoleJoinRandomRoom() 的耗时约 108 us，即每秒 9300 次。该数值为 RPC 调用的极限值。
```
D:/Go/bin/go.exe test -test.bench=.* [D:/Daisy/Cinder_Server_Cinder/Matcher/matcherlib/internal/rpcproc]
1587627982140279300 [Error] Read server.json failed, use default value
1587627982143324900 [Info] Logdir:../log/ Loglevel:debug
...
goos: windows
goarch: amd64
pkg: Cinder/Matcher/matcherlib/internal/rpcproc
BenchmarkRPC_RoleJoinRandomRoom-4   	   14632	    107908 ns/op
PASS
ok  	Cinder/Matcher/matcherlib/internal/rpcproc	3.563s
成功: 进程退出代码 0.
```

#### 各接口测试结果

goroutines 大概在 1000-3000 会到达 RPS 的最大值，此时 CPU 满载，利用率高。
继续加大 goroutines 会因为协程消耗使 RPS 略微下降, 并且可能会出现超时错误。
测试时长为 15s.

test名					|goroutines |RPS 	|说明
------------------------|-----------|-------|------------------------------------------------------------------
BroadcastRoom			|1200		|7587	|房间内有3人
BroadcastTeam           |1000       |6286   |队伍内有3人
ChangeTeamLeader        |1300       |7232   |队伍内有3人
CreateLeaveTeam         |3000       |3800*2 |创建队伍并立即离开
GetRoomInfo             |3000       |6300   |
GetRoomList             |3000       |9000   |房间列表为空
JoinLeaveTeam           |3000       |3600*2 |队伍本身有1人，另一人加入并立即离开
MyRPCTest               |3000       |10000  |RPC_MyRPCTest(arg string) uint32
Ping                    |3000       |11000  |RPC_Ping 是 MatcherExample 的一个无参无返回的 RPC
RoleCreateDeleteRoom    |3000       |3200*2 |创建房间并立即删除
RoleJoinLeaveRoom       |2200       |2347*2 |房间本身有2人，第3人加入并立即离开
RoleJoinRandomRoom      |3000       |4334   |进入 "2v2" 房间
RoleJoinRoom            |2000       |4552   |第1人创建，后面3人加入房间
SetRoomData             |2000       |4500   |房间内只有一人
SetRoomRoleData         |2000       |5917   |房间内只有一人
SetTeamData             |2500       |6406   |队伍内有3人
TeamCreateDeleteRoom    |2500       |2666*2 |队伍有1人，创建房间后立即删除
TeamJoinRandomRoom      |2000       |1372*2 |队伍有1人，加入随机房间后立即删除
TeamJoinRoom            |3500       |1817*3 |队伍有1人，1人先创建房间，队伍加入房间，然后删除房间

结果说明：服务器通知广播明显拉低了整体性能，广播量越大性能越差。
如果没有通知消息合并的优化，性能还会降低50%。

## PProf

在压测过程中，可查看性能分析数据 http://localhost:16060/debug/pprof/

### CPU Profile

压测开始后，运行 `go tool pprof http://localhost:16060/debug/pprof/profile` 获取 CPU 分析数据，
30s 后可输入 web 打开 CPU 耗时分析图
```
D:\Daisy\Server\bin (master -> origin)
λ go tool pprof http://localhost:16060/debug/pprof/profile
Fetching profile over HTTP from http://localhost:16060/debug/pprof/profile
...
(pprof) web
```

如 [RoleJoinRandomRoom CPU profile](pprof_RoleJoinRandomRoom.svg)
