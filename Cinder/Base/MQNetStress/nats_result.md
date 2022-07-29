# nats 压测

测试 MQNet/mqnats 的最大每秒请求数(RPS).

## 测试结果

旧版的发送很容易出现 sendChan 满造成丢包的错误，所以删除了。
以下列出了新旧版本的测试结果，

```
Revision: b4b931da9d055e9b4be480979c87ecc8d00d53ea
Author: 江剑 <jiangjian@ztgame.com>
Date: 2020/9/18 16:21:09
Message:
1. 去除MQ的消息发送协程, 在调用Post的时候直接发送消息
2. 增加Post方法出错的错误日志
```

本机测试环境:
* Processor: Intel Core i5-7400 CPU @ 3.00GHz
* RAM: 8 GB

测试时长为 15s

### 新版本直接发送

RPS 和旧版相比，略有提高，大概有 10%, 延时多数下降了。
协程数到 2000 仍会报 "msg chan full"，出现丢包。
但和旧版的报错不同，是 onRecvMsg() 接收队列不够大。

2000以上协程数的测试是加大 msgChan 缓存到 1000000 的结果。
2000-5000协程下，rps 差不多相同，延时加大。
最大 rps 约 3.5K.

协程数	| RPS (r/s)	| 最大延时 (ms)
--------|-----------|---------------
1       |1170       |18
2       |2109       |22
3       |3070       |24
4       |3852       |26
5       |3919       |34
10      |6795       |18
20      |10546      |19
30      |12629      |19
50      |15458      |23
100     |19401      |28
200     |22302      |45
500     |27884      |52
1000    |30479      |72
2000    |34772      |113
3000    |34082      |131
4000    |34576      |237
5000    |34948      |202
10000   |30032      |522

### 旧版带 sendChan 

协程数到 2000 就会报 "msg chan full"，出现丢包。

协程数	| RPS (r/s)	| 最大延时 (ms)
--------|-----------|---------------
1       |1125       |16
2       |1761       |30
3       |2980       |11
4       |3722       |20
5       |4247       |23
10      |6149       |24
20      |9904       |21
30      |11786      |30
50      |15020      |18
100     |17330      |103
200     |22134      |28
500     |23275      |102
1000    |28745      |89

2000: msg chan full