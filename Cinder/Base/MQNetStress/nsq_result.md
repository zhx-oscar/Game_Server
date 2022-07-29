# MQNet 压测

测试 MQNet 的最大每秒请求数(RPS).

参数：
```
λ MQNetStress.exe -h
Usage of MQNetStress.exe:
  -g, --goroutines int      goroutines count (default 1)
  -n, --nsq_addr string     NSQ address (default "127.0.0.1:4150")
  -a, --nsq_admin string    NSQAdmin address (default "127.0.0.1:4171")
  -l, --nsq_lookup string   NSQLookup address (default "127.0.0.1:4161")
  -m, --time duration       test time (default 15s)
pflag: help requested
```

示例：
```
D:\Daisy\Cinder_Server_Cinder\Base\MQNetStress (7-mqnetstress -> origin)
λ MQNetStress.exe --goroutines 4 --time 3s
1588057570735289800 [Error] ERR    2 [5ea7d5e1acd617121428f8c5#ephemeral/channel_mqnetstress_5ea7d5e1acd617121428f8c5#ephemeral] error querying nsqlookupd (http://127.0.0.1:4161/lookup?topic=5ea7d5e1acd617121428f8c5%23ephemeral) - got response 404 Not Found "{\"message\":\"TOPIC_NOT_FOUND\"}"
1588057570738265400 [Error] ERR    2 [5ea7d5e1acd617121428f8c5#ephemeral/channel_mqnetstress_5ea7d5e1acd617121428f8c5#ephemeral] error querying nsqlookupd (http://127.0.0.1:4161/lookup?topic=5ea7d5e1acd617121428f8c5%23ephemeral) - got response 404 Not Found "{\"message\":\"TOPIC_NOT_FOUND\"}"
1588057570741267700 [Error] ERR    2 [5ea7d5e1acd617121428f8c5#ephemeral/channel_mqnetstress_5ea7d5e1acd617121428f8c5#ephemeral] error querying nsqlookupd (http://127.0.0.1:4161/lookup?topic=5ea7d5e1acd617121428f8c5%23ephemeral) - got response 404 Not Found "{\"message\":\"TOPIC_NOT_FOUND\"}"
1588057570747945000 [Error] ERR    3 [mqnetstress#ephemeral/channel_mqnetstress_5ea7d5e1acd617121428f8c5#ephemeral] error querying nsqlookupd (http://127.0.0.1:4161/lookup?topic=mqnetstress%23ephemeral) - got response 404 Not Found "{\"message\":\"TOPIC_NOT_FOUND\"}"
1588057570752330500 [Error] ERR    3 [mqnetstress#ephemeral/channel_mqnetstress_5ea7d5e1acd617121428f8c5#ephemeral] error querying nsqlookupd (http://127.0.0.1:4161/lookup?topic=mqnetstress%23ephemeral) - got response 404 Not Found "{\"message\":\"TOPIC_NOT_FOUND\"}"
1588057570754331000 [Error] ERR    3 [mqnetstress#ephemeral/channel_mqnetstress_5ea7d5e1acd617121428f8c5#ephemeral] error querying nsqlookupd (http://127.0.0.1:4161/lookup?topic=mqnetstress%23ephemeral) - got response 404 Not Found "{\"message\":\"TOPIC_NOT_FOUND\"}"
MQNet stress test is running (goroutines=4 time=3s)...
1588057572092318500 [Debug] Receive mq hello message, message: this is hello from my broad channel mqnetstress
1588057572213319300 [Debug] Receive mq hello message, message: this is hello from my channel 5ea7d5e1acd617121428f8c5
RPS: 28 / 3.000000 = 9.333333 (r/s), max delay: 251ms
```

## 测试结果

本机测试环境:
* Processor: Intel Core i5-7400 CPU @ 3.00GHz
* RAM: 8 GB

测试时长为 15s

协程数	| RPS (r/s)	| 最大延时 (ms)
--------|-----------|---------------
1       |3.67       |251 
2       |7.33       |253
3       |11.0       |251
4       |14.7       |253
5       |18.3       |253
10      |36.7       |252
20      |73.3       |254
30      |110        |254
50      |180        |260
100     |9444       |153
200     |11045      |94
500     |14543      |89
1000    |17254      |159
2000    |18601      |241
3000    |19702      |365
4000    |18306      |343
5000    |19779      |407
6000    |20782      |510
7000    |19778      |476

