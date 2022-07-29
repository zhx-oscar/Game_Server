package ErrorCode

// build 1301~
const (
	//OverLimitBuildCountMaxNotFind 超出build列表上限
	OverLimitBuildCountMaxNotFind = 1301

	//NoSpecialAgent 没有该特工
	NoSpecialAgent = 1302

	//NotFindBuild 对应build没找到
	NotFindBuild = 1303

	//NotFindBuildEquip build内装备的装备道具 不存在
	NotFindBuildEquip = 1304

	//SkillBagNotEnoughSpace 技能背包空间不足
	SkillBagNotEnoughSpace = 1305

	//EquipBagNotEnoughSpace 装备背包空间不足
	EquipBagNotEnoughSpace = 1306

	//SkillItemNotFind 技能道具配置表中不存在
	SkillItemNotFind = 1307

	//SkillNotMatchSpecialAgent 技能不匹配当前特工
	SkillNotMatchSpecialAgent = 1308

	//BuildSuperSkillCountMaxNotFind build内超能技槽位数量限制配置不存在
	BuildSuperSkillCountMaxNotFind = 1309

	//OverLimitBuildSuperSkillMaxNotFind 超出超能技槽位限制
	OverLimitBuildSuperSkillMaxNotFind = 1310

	//BuilldSkillKindNotMatch build内技能类型不匹配
	BuilldSkillKindNotMatch = 1311

	//BuildEquipNotEnoughLevel build内特工等级不够装备
	BuildEquipNotEnoughLevel = 1312

	//BuildEquipPosWrong build装备位置不对
	BuildEquipPosWrong = 1313

	//BuilldSkillSuperPosError build内技能类型不匹配
	BuilldSkillSuperPosError = 1314

	//BuildCountMaxNotFind build列表上限不存在
	BuildCountMaxNotFind = 1315
)

// build skill 1101-
const (
	//BuildUpgradeSkillMaxLv 已经达到最大等级
	BuildUpgradeSkillMaxLv = 1101

	//BuildUpgradeSkillCoinNotEnough 技能升级需要货币数量不满足
	BuildUpgradeSkillCoinNotEnough = 1102

	//BuildUpgradeSkillmaterialsNumNotEnough 技能升级需要 道具数量不满足
	BuildUpgradeSkillmaterialsNumNotEnough = 1103

	//BuildUpgradeSkillmaterialsNumNotMatchMaxLv 技能升级 不同等级材料和等级数量不匹配
	BuildUpgradeSkillmaterialsNumNotMatchMaxLv = 1104

	//BuildSkillNotLearned 技能还未学习
	BuildSkillNotLearned = 1105
)

// name 1503-
const (
	NameInvalid   = 1501 //名字不合法
	NameDuplicate = 1502 //重名
)

// coin 9701-
const (
	DiamondNotEnough = 9701 //钻石不足
	GoldNotEnough    = 9702 //金币不足
)

// fastbattle 2501-
const (
	UnEnergize     = 2501 //未充能
	AlreadySpeedUp = 2502 //加速中
	SpeedUpExhaust = 2503 //加速次数耗尽
)

// chat 1601-
const (
	//ChatMsgContent 聊天内容超越长度限制
	ChatMsgContent = 1601

	//ChatPrivateSendFail 私聊投递失败
	ChatPrivateSendFail = 1602
)

// supply 1801-
const (
	GiftBagUnknown        = 1801 //礼包id未知
	GiftBagOpenUpperLimit = 1802 //礼包开启次数已达上限
	GiftBagNotEnough      = 1803 //礼包数不足
)
