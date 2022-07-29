# MQNet mqnats 压测

测试 MQNet/mqnats 的最大每秒请求数(RPS).
测试结果见：[nats_result.md](nats_result.md)

旧的 mqnsq 测试结果见：[nsq_result.md](nsq_result.md)

参数：
```
λ MQNetStress.exe -h
Usage of MQNetStress.exe:
  -g, --goroutines int     goroutines count (default 1)
  -a, --nats.addr string   NATS address (default "127.0.0.1:4222")
  -m, --time duration      test time (default 15s)
pflag: help requested
```
