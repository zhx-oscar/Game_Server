# todo

* 标志：不加入别人房间
* 排除特定人，即踢人
* 房间加密码，不允许他人加入，邀请他人即发送密码，房间的匹配条件很松
* 可以同时匹配加入多个房间，直到有个房间先满，自动退出其他房间。
	+ 哪个房间先满就匹配到哪个。但是只显示一个房间，选择最可能成功的房间。
* 小队匹配支持人满和人不满模式
	+ 人满即 5v5 ，则最终匹配结果肯定是 5v5 ,否则会匹配失败；
	+ 人不满即 5v5 ,在超时之前可能只能匹配到 2v3 ,则依然匹配成功，你可以再自行添加机器人。
* 防止Game宕后无法删除房间
* 添加 SetHidden 隐藏房间接口
* 房间 Leader: 选择 Leader，转移
* Role 加入房间的时间和序号，用于客户端显示
* 不允许同时在多个房间。需要角色号查询房间号。与房间管理器之间因为允许并发很难保持数据一致。
* 立即开始，即强制开始
	+ 可实现为自定义事件，该事件原样通知 roomEventHandler, 由其动作。
		- roomEventHandler 需要能够房间广播: RPC 返回时的房间数据可利用为广播数据
* 同一房间加减人消息合并发送，减少数据量
	- 房间有个消息队列，延时发送
* RPC 改为 MessageProc 可提高性能