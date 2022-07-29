package ErrorCode

// 好友相关 1801-

const (
	FriendRemarkOver              = 1801 // 好友备注超过固定长度
	FriendAddRepeated             = 1802 // 重复添加好友
	FriendListFull                = 1803 // 好友数量达到上限
	FriendTargetApplyListFull     = 1804 // 对方好友申请列表已满
	FriendDBFailed                = 1805 // 聊天服错误
	FriendNotExit                 = 1806 // 好友不存在
	FriendFindSelf                = 1807 // 查找好友找到自己
	FriendFindNoExit              = 1808 // 查找不到该好友
	FriendTargetFriendListFull    = 1809 // 对方好友列表已满
	FriendAllTargetFriendListFull = 1810 // 申请列表里所有申请者的好友都已经满了
	FriendUnknowErr               = 1811 // 好友服未定义的错误
)
