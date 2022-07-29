# 说明

设计前提条件：
* 每个用户的群个数有限(但群内成员数可以很大)
* 每个用户的关注有限(但关注某个用户的粉丝数可能会巨量)
* 每个用户的好友有限。
* 每个用户的黑名单有限(但拉黑某个用户的可能会巨量)

## Mongodb 集合设计

### 用户集合

每个用户一条记录。
Friends, FollowingList和Blacklist个数有限，以数组内嵌在用户记录中。
用户的群虽然个数有限，但因为还需要记录群已读聊天记录序号，所以没有内嵌。

chat.users 用户集合有以下字段：
* userID: string, 用户ID,
	+ $hashed:userID 片键，唯一
* nick: string, 昵称
* activeData: bin, 主动数据，用户自定义数据，变化时主动广播
* passiveData: bin, 被动数据，用户自定义数据，变化时不会主动广播
* friendIDs: []string, 好友ID
* followIDs: []string, 关注ID
* blacklist: []string, 黑名单
* offlineTime: time, 离线时间
* followerNumber: int, 粉丝数

因为需要查询粉丝列表，且粉丝数量不限，所以独立粉丝集合.
chat.users.followers 粉丝集合有以下字段：
* userID: string, 用户ID
	+ $hashed:userID 片键，非唯一
* follower: string, 粉丝ID
	+ (userID, follower) 唯一索引

关注和被关注数据不一致时，以 chat.users.followIDs 的为准，在双方同时上线时检查并修正。

反向黑名单列表，即哪些人拉黑了我，这个不需要查询，所以没有对应集合。

### 加好友请求和返回

加好友请求在对方离线时暂存DB, 对方上线时发送并删除。
同意或拒绝的返回在申请人离线时，也会暂存DB。
因为没有数量限制，所以独立集合。

* chat.users.friend_reqeusts, 用户离线时收到的加好友请求
	+ userID: string, 用户ID, 被加人
		- $hashed:userID 片键，非唯一
	+ fromID: string, 加好友请求人ID
		- (userID, fromID) 唯一索引
	+ reqInto: []byte, 请求信息
* chat.users.friend_responses, 用户离线时收到的加好友应答
	+ userID: string, 用户ID, 加好友请求人
		- $hashed:userID 片键，非唯一
	+ responderID: string, 响应人ID
		- (userID, responderID) 唯一索引
	+ ok: bool, 是否同意

### 聊天群
* chat.users.groups, 用户聊天群集合
	+ userID: string, 用户ID
		- $hashed:userID 片键, 非唯一
	+ groupID: string, 用户的聊天群ID
		- (userID, groupID) 唯一索引
	+ sequenceID: NumberLong, 已收群消息序号
* chat.groups.members, 聊天群成员
	+ groupID: string, 聊天群ID
		- $hashed:groupID 片键，非唯一
	+ memberID: string, 成员ID
		- (groupID, memberID) 唯一索引

允许聊天群成员为0，所以无法从 chat.group.members 集合判断群是否存在。
允许加载一个成员为0的群。

因为群成员数可能很大，所以不用数组表示。

用户<->聊天群 是多对多关系，为了双向查询并支持分片，数据在2个集合中重复。
* chat.users.groups 保存了用户有哪些群
* chat.groups.members 保存了群有哪些成员
因为操作中间出错，造成数据不一致时，以 chat.groups.members 为准。
当用户记录群中已读序号时，可以修正不一致，使之符合群成员的记录。

### 消息记录

群聊记录不删。离线私聊上线获取之后即删除。如果离线私聊有限，可以考虑将他嵌到用户记录中去。

* chat.users.offline_messages, 用户离线私聊记录
	+ userID: string, 用户ID
		- $hashed:userID 片键, 非唯一
	+ fromID: string, 私聊来自ID
	+ fromNick: string, 来自昵称
	+ fromData: Bin, 发送者的自定义数据
	+ sendTime: Date, 发送时间
		- (userID, sendTime) 非唯一索引，用于排序
	+ msgContent: Bin, 聊天信息
* chat.groups.messages, 群聊记录
	+ groupID: string, 群ID
		- $hashed:groupID 片键，非唯一
	+ fromID: string, 来自ID
	+ fromNick: string, 来自昵称
	+ fromData: Bin, 发送者的自定义数据
	+ sendTime: Date, 发送时间
	+ msgContent: Bin, 聊天信息
	+ sequenceID: NumberLong, 消息序号，群内从1开始递增
		- (groupID, sequenceID) 唯一索引
