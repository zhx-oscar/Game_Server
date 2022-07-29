# 通用邮件服务器

通用邮件服务器仅对游戏服务器提供服务，不会与客户端直连。

任意启动一个或多个邮件服务器，可以任意关停或新开。

通过 mailapi 包操作邮件。可以任选一服进行请求，并任意切换。

邮件有2类：普通邮件和系统广播邮件。邮件可带道具附件。

接入说明见：[mailapi/mail_usage.md](mailapi/mail_usage.md)

## 设计

### 用户在线状态记录
用户调用Login()时，将以下信息记录DB：
* 记录用户的服务器ID, 用于新邮件通知
* 初次 Login 记录时间作为创建时间，用于过滤历史邮件

用户Login()时所在的服称为 UserSrv(用户服), 是用户实例所在的服务器。

#### 普通邮件通知
因为Mail服多开，所以用户登录的Mail服和触发新邮件通知的Mail服不是同一服，
如何知道新邮件通知到哪个用户服呢? 此时需要查询 DB.

#### 系统邮件广播
触发系统邮件广播时，如何知道新邮件通知到哪些服务器呢？

* 每个邮件服内存记录所有用户的用户服，并相互同步。
* 用户服列表只增不减，除非判断为已关而删除。
	+ 存在临时断开而被删的情况，直到有用户从该服Login才会恢复.
* 用户服列表必须存DB, 启动时先读DB, 再初始化网络。
* 广播时对所有记录的 UserSrv 按 srvID 进行通知

## MongoDB 设计

采用 MongoDB 存储数据，不用 redis.

* mail.broadcast_mails (broadcast)全服广播邮件，以shard为键并分片
	+ _id ID 自动ID, 无用忽略
	+ shard int 分片号
		- $hashed:shard 片键
	+ originalID ID 原始ID
		- (shard, originalID), 唯一
	+ mail Mail结构体
		- (shard, mail.sendTime, originalID), 唯一, 用于查询邮件列表
		- mail.expireTime, TTL 用于自动删除
* mail.user_srv_ids 用户服ID列表，不分片
	+ srvID: string, 唯一键
* mail.users 用户记录，以用户为键并分片
	+ _id: 自动ID, 用来取创建时间
	+ userID: string, 用户ID
		- $hashed:userID 分片键
		- userID 唯一索引
	+ srvID: string, 用户服ID
* mail.users.bc_mail_states 全服邮件状态，以收件人为键并分片
	+ _id 自动ID, 忽略无用
	+ to string 收件人ID
		+ $hashed:to 片键
	+ originalID ID 全服邮件的原始ID
		+ (to, originalID) 唯一, 用于设置状态
	+ state State结构
	+ sendTime time 用于检索时分页
		+ (to, sendTime) 非唯一, 用于列举全服邮件时，按时间范围查询状态
	+ expireTime time
		+ TTL 自动删除
	+ deleted bool 表示邮件已删除
* mail.users.mails 用户邮件(非全服邮件)，以收件人为键并分片
	+ _id ID
	+ mail Mail结构体
		+ $hashed:mail.to 片键
		+ (mail.to, _id), 唯一
		+ (mail.to, mail.sendTime, _id), 唯一, 用于查询邮件列表
		+ mail.expireTime, TTL 用于自动删除

参见：[maildoc.go](rpcproc/handler/internal/maildoc/maildoc.go)
