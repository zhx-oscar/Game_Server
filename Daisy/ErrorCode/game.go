package ErrorCode

// 挑战BOSS错误定义 9401-
const (
	// 策划未配置该进度
	RaidProgressInvalid = 9401
	// 挑战BOSS所需的门票不足
	TicketsNotEnough = 9402
)

// 跑图宝箱错误定义 9501-
const (
	// 非法拾取
	InvalidLootRunChest = 9501
)

// 装备赠送错误定义(战利品分享) 2601-
const (
	// 装备容器未找到
	EquipContainerNotFound = 2601
	// 装备不存在
	EquipNotExist = 2602
	// 不能赠送给自己
	CantGiveToSelf = 2603
	// 目标不在赠送列表内
	TargetNotInList = 2604
	// 留言长度超过限制
	CharacterCountExceed = 2605
	// 装备已在build中使用
	EquipUsedInBuild = 2606
	// 今日赠送次数已经超过上限
	GiveCountOverLimit = 2607
	// 删除装备失败
	RemoveEquipFailed = 2608
)

// 技能乞求错误定义(战利品分享） 2609-
const (
	// 今日乞求次数已经超过上限
	RequestCountOverLimit = 2609
	// 无法向自己捐赠技能
	CantGiveSkillToSelf = 2610
	// 没有技能可以捐赠
	DontHaveSkill = 2611
	// 赠送对象已不在队伍中
	TargetNotInTeam = 2612
	// 赠送的技能和目标乞求的技能不一样
	GiveSkillNotRequest = 2613
	// 目标收到的技能道具数量已达到上限
	ReceiveSkillCountReachLimit = 2614
	// 赠送技能数量超过上限
	GiveSkillCountReachLimit = 2615
	// 赠送给目标的技能数量超过上限
	GiveSKillToTargetCountReachLimit = 2616
	// 请先加入队伍再发起乞求
	FirstJoinATeam = 2617
	// 技能乞求已失效
	SkillRequestNotValid = 2618
	// 技能容器未找到
	SkillContainerNotFound = 2619
	// 扣除技能道具失败
	RemoveSkillFailed = 2620
)
