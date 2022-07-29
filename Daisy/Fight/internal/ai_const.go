package internal

//================AI相关
const (
	enemyList    = "enemyList"    //黑板敌人列表key
	attackTarget = "attackTarget" //黑板 当前攻击目标key

	attackPos                  = "attackPos"                  //黑板攻击位置key
	curSkill                   = "curSkill"                   //黑板当前技能key
	detourNextPos              = "detourNextPos"              //迂回下一刻位置
	detourNextPosMoveMode      = "detourNextPosMoveMode"      //迂回下一刻位置移动模式
	lastDistanceWithTarget     = "lastDistanceWithTarget"     //上一次与目标之间的距离
	isAlreadyCastSkill         = "isAlreadyCastSkill"         //释放技能 如果需要技能释放节点running状态的话，需要紧紧释放一次技能
	CurrentHP                  = "CurrentHP"                  //当前血量
	ChangeFormBeginTime        = "ChangeFormBeginTime"        //形态转变开始时间key
	ChangeFormEnd              = "ChangeFormEnd"              //当前形态转变是否结束key
	moveToSuccess              = "moveToSuccess"              //移动open成功
	detourMoveTargetLastPos    = "detourMoveTargetLastPos"    //迂回目标的上一次位置
	dash_target                = "dash_target"                //冲刺-目标
	dash_skill                 = "dash_skill"                 //冲刺-技能
	randBaseAttackDis          = "randBaseAttackDis"          //根据最短最长攻击距离 随机出base攻击距离
	lastTargetPos              = "lastTargetPos"              //上一次目标位置
	beginTimeSteeringSmoothing = "beginTimeSteeringSmoothing" //转向平滑处理开始时间
	RandWaitActionStartTime    = "RandWaitActionStartTime"    //随机等待节点中开始时间
	RandWaitActionRandTime     = "RandWaitActionRandTime"     //随机等待节点中 实际等待的时间
	RetreatPos                 = "RetreatPos"                 //后撤目标点
	BeginMoveTime              = "BeginMoveTime"              //后撤开始时间key
	StopMoveTime               = "StopMoveTime"               //后撤开始时间key
)

const (
	skillType_super    = 1 //超能技
	skillType_ultimate = 2 //必杀技
)

const (
	CastBloadSkillRunning_castSkillEnd = 1 //castSkill直到技能结束
	CastBloadSkillRunning_canCastSkill = 2 //castSkill直到技能命中之后，可以继续释放其他技能
)

//AI externalBlackboardKeys
const (
	eB_Attr_HP_Per = "Attr_HP_Per" //血量百分比

	eB_Attr_HP = "Attr_HP" //血量属性

	eB_SuperSkillInterruptNormalSkill      = "SuperSkillInterruptNormalSkill"      //超能技是否可以中断普攻连击 bool
	eB_UltimateSkillInterruptNormalSkill   = "UltimateSkillInterruptNormalSkill"   //必杀技是否可以中断普攻连击 bool
	eB_DetourMove_radius                   = "DetourMove_radius"                   //迂回后退距离
	eB_DetourMove_Angle                    = "DetourMove_Angle"                    //迂回开角
	eB_DetourMove_Duration                 = "DetourMove_Duration"                 //迂回持续最大时间 毫秒
	eB_DetourMove_RandSuccess              = "DetourMove_RandSuccess"              //触发迂回概率
	eB_SelfRangeHasEnemy_Range             = "SelfRangeHasEnemy_Range"             //自己半径内有敌人-半径
	eB_SkillAttackRangeHasEnemy_Skillindex = "SkillAttackRangeHasEnemy_Skillindex" //周边有敌人使用技能攻击-技能索引
	eB_SkillAttackRangeHasEnemy_enemyCount = "SkillAttackRangeHasEnemy_enemyCount" //技能攻击范围内是否有敌人-敌人数量
	eB_targetDistanceLowestThan_Distance   = "targetDistanceLowestThan_Distance"   //目标距离低于多少-距离
)

//===================================AI end
