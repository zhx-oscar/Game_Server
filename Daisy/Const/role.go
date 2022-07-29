package Const

//特工相关常量
const (
	SpecialAgent_buildCountMax         = 1 //特工常量表---build列表上限
	SpecialAgent_buildSuperSlillCount  = 2 //特工常量表---build内超能技槽位数量
	SpecialAgent_upgradeExpCoefficient = 3 //特工常量表---有满级特工之后经验加成系数
	SpecialAgent_upgradeExpBase        = 4 //特工常量表---有满级特工之后经验加成基础值
)

const (
	Gold    = 1 //金币，游戏内流通，可以用钻石兑换  1表示货币表中的ID
	Diamond = 2 //钻石,人民币兑换来的	2表示货币表中的ID
)

const (
	NORMAL_ACTION  uint32 = iota + 1 // 常规货币操作
	SELL_EQUIPMENT                   // 卖装备
	Energize                         //充能
	SpeedUp                          //加速
	SupplyBoxCost                    //补给箱花费
)

// 好友相关数据
const (
	FriendApplyListFull             = 100 // 目标好友申请列表已满
	FriendNumMax                    = 100 // 好友数量上限
	FriendRecommendMax              = 61  // 从数据库里拉取60条推荐好友放到缓存里
	FriendRecommendOnce             = 10  // 每次需要从缓存里取出10条给客户端显示
	FriendRecommendOfflineTimeLimit = 3   // 推荐好友离线时间不超过3天
	FriendRemarkLengthLimit         = 12  // 6个中文字符，12个英文字符
)

type RefreshDataToFriendType uint8

const (
	FriendTypeLevel RefreshDataToFriendType = iota + 1 // 1:等级刷新通知好友和群
	FriendTypeHead                                     // 2:头像刷新通知好友和群
	FriendTypeSeasonScore								//3:赛季积分刷新
)
