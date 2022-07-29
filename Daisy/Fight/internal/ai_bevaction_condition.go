package internal

import (
	"Daisy/Fight/internal/log"
	b3 "github.com/magicsea/behavior3go"
	b3config "github.com/magicsea/behavior3go/config"
	b3core "github.com/magicsea/behavior3go/core"
	"math"
	"math/rand"
	"strconv"
	"strings"
)

// HPLowerThan 血量低于百分比
type HPLowerThan struct {
	b3core.Condition
	HPValue float64
}

// Initialize 初始化
func (bev *HPLowerThan) Initialize(setting *b3config.BTNodeCfg) {
	bev.Condition.Initialize(setting)
	bev.HPValue = setting.GetProperty("HPValue")
}

// OnTick 循环
func (bev *HPLowerThan) OnTick(tick *b3core.Tick) b3.Status {
	pawn, b := tick.GetTarget().(*Pawn)
	if !b {
		log.Error("*************[HPLowerThan] trans Pawn fail")
		return b3.ERROR
	}

	pawn.pushActionTorsionInfo(bev.GetName())

	if float64(pawn.Attr.CurHP)/float64(pawn.Attr.MaxHP) < bev.HPValue {
		return b3.SUCCESS
	}

	return b3.FAILURE
}

// targetDistanceLowestThan 目标距离低于 XX 米 Value：XX
type targetDistanceLowestThan struct {
	b3core.Condition
	DistanceKey string
}

// Initialize 初始化
func (bev *targetDistanceLowestThan) Initialize(setting *b3config.BTNodeCfg) {
	bev.Condition.Initialize(setting)
	bev.DistanceKey = setting.GetPropertyAsString("DistanceKey")
}

// OnTick 循环
func (bev *targetDistanceLowestThan) OnTick(tick *b3core.Tick) b3.Status {
	pawn, b := tick.GetTarget().(*Pawn)
	if !b {
		log.Error("*************[targetDistanceLowestThan] trans Pawn fail")
		return b3.ERROR
	}

	pawn.pushActionTorsionInfo(bev.GetName())

	target, ok := getPawnBoard(tick, attackTarget)
	if !ok {
		log.Error("*************[targetDistanceLowestThan] getPawnBoard fail")
		return b3.FAILURE
	}

	//外部黑板数据处理
	value, ok := pawn.getBlackboardValueByKey(bev.DistanceKey).(float64)
	if !ok {
		log.Error("*************[targetDistanceLowestThan] 外部黑板 DistanceKey 没有配置")
		return b3.ERROR
	}

	if DistancePawn(pawn, target) <= float32(value) {
		return b3.SUCCESS
	}

	return b3.FAILURE
}

// curSkillIsNormalAttack 当前技能是否为普攻
type curSkillIsNormalAttack struct {
	b3core.Condition
}

// Initialize 初始化
func (bev *curSkillIsNormalAttack) Initialize(setting *b3config.BTNodeCfg) {
	bev.Condition.Initialize(setting)
}

// OnTick 循环
func (bev *curSkillIsNormalAttack) OnTick(tick *b3core.Tick) b3.Status {
	pawn, b := tick.GetTarget().(*Pawn)
	if !b {
		log.Error("*************[targetDistanceLowestThan] trans Pawn fail")
		return b3.ERROR
	}

	pawn.pushActionTorsionInfo(bev.GetName())

	target, ok := getSkillBoard(tick, curSkill)
	if !ok {
		log.Error("*************[curSkillIsNormalAttack] getSkillBoard fail")
		return b3.FAILURE
	}

	if target.IsNormalAttack() {
		return b3.SUCCESS
	}

	tick.Blackboard.Remove(curSkill)
	return b3.FAILURE
}

// EnableUltimateSkillIsNil 必杀技可用数量为0
type EnableUltimateSkillIsNil struct {
	b3core.Condition
}

// Initialize 初始化
func (bev *EnableUltimateSkillIsNil) Initialize(setting *b3config.BTNodeCfg) {
	bev.Condition.Initialize(setting)
}

// OnTick 循环
func (bev *EnableUltimateSkillIsNil) OnTick(tick *b3core.Tick) b3.Status {
	pawn, b := tick.GetTarget().(*Pawn)
	if !b {
		log.Error("*************[targetDistanceLowestThan] trans Pawn fail")
		return b3.ERROR
	}

	pawn.pushActionTorsionInfo(bev.GetName())

	if len(pawn.GetUsableUltimateSkills()) > 0 {
		return b3.FAILURE
	}

	return b3.SUCCESS
}

// EnableSuperSkillIsNil 超能可用数量为0
type EnableSuperSkillIsNil struct {
	b3core.Condition
}

// Initialize 初始化
func (bev *EnableSuperSkillIsNil) Initialize(setting *b3config.BTNodeCfg) {
	bev.Condition.Initialize(setting)
}

// OnTick 循环
func (bev *EnableSuperSkillIsNil) OnTick(tick *b3core.Tick) b3.Status {
	pawn, b := tick.GetTarget().(*Pawn)
	if !b {
		log.Error("*************[EnableSuperSkillIsNil] trans Pawn fail")
		return b3.ERROR
	}

	pawn.pushActionTorsionInfo(bev.GetName())

	if len(pawn.GetUsableSuperSkills()) > 0 {
		return b3.FAILURE
	}

	return b3.SUCCESS
}

// SelfRangeHasEnemy 自身 半径 范围内是否有敌人
type SelfRangeHasEnemy struct {
	b3core.Condition
	RadiusKey string
}

// Initialize 初始化
func (bev *SelfRangeHasEnemy) Initialize(setting *b3config.BTNodeCfg) {
	bev.Condition.Initialize(setting)
	bev.RadiusKey = setting.GetPropertyAsString("RadiusKey")
}

// OnTick 循环
func (bev *SelfRangeHasEnemy) OnTick(tick *b3core.Tick) b3.Status {
	pawn, b := tick.GetTarget().(*Pawn)
	if !b {
		log.Error("*************[SelfRangeHasEnemy] trans Pawn fail")
		return b3.ERROR
	}

	pawn.pushActionTorsionInfo(bev.GetName())

	radius, ok := pawn.getBlackboardValueByKey(bev.RadiusKey).(float64)
	if !ok {
		log.Error("*************[SelfRangeHasEnemy] 外部黑板 RadiusKey 没有配置")
		return b3.ERROR
	}

	var enemys []*Pawn
	targetList := pawn.Scene.overlapCircleShape(float64(pawn.GetPos().X), float64(pawn.GetPos().Y), radius)
	for _, enemy := range targetList {
		//同阵营或者就是自己都需要过滤
		if !(enemy.UID == pawn.UID || enemy.GetCamp() == pawn.GetCamp()) {
			enemys = append(enemys, enemy)
		}
	}

	if len(enemys) > 0 {
		tick.Blackboard.Set(enemyList, enemys, "", "")
		return b3.SUCCESS
	}

	return b3.FAILURE
}

// SkillAttackRangeHasEnemy 技能攻击范围内是否有敌人
type SkillAttackRangeHasEnemy struct {
	b3core.Condition
	skillIndexkey string
	enemyCountkey string
}

// Initialize 初始化
func (bev *SkillAttackRangeHasEnemy) Initialize(setting *b3config.BTNodeCfg) {
	bev.Condition.Initialize(setting)
	bev.skillIndexkey = setting.GetPropertyAsString("skillIndexkey")
	bev.enemyCountkey = setting.GetPropertyAsString("enemyCountkey")
}

// OnTick 循环
func (bev *SkillAttackRangeHasEnemy) OnTick(tick *b3core.Tick) b3.Status {
	pawn, b := tick.GetTarget().(*Pawn)
	if !b {
		log.Error("*************[SkillAttackRangeHasEnemy] trans Pawn fail")
		return b3.ERROR
	}

	pawn.pushActionTorsionInfo(bev.GetName())

	skillList := pawn.GetSuperSkillList()
	if len(skillList) == 0 {
		return b3.FAILURE
	}

	var skillIndex, enemyCount int
	//外部黑板数据处理
	{
		value, ok := pawn.getBlackboardValueByKey(bev.skillIndexkey).(float64)
		if !ok {
			log.Error("*************[SkillAttackRangeHasEnemy] 外部黑板 skillIndexkey 没有配置")
			return b3.ERROR
		}
		skillIndex = int(value)

		value, ok = pawn.getBlackboardValueByKey(bev.enemyCountkey).(float64)
		if !ok {
			log.Error("*************[SkillAttackRangeHasEnemy] 外部黑板 enemyCountkey 没有配置")
			return b3.ERROR
		}
		enemyCount = int(value)
	}

	if skillIndex >= len(skillList) {
		log.Error("*************[SkillAttackRangeHasEnemy] superSkillList out of range ", skillIndex)
		return b3.FAILURE
	}
	skill := skillList[skillIndex]
	enemys := skill.SearchTargets(pawn.GetPos())
	tick.Blackboard.Set(enemyList, enemys, "", "")
	if skill != nil && len(enemys) <= enemyCount {
		return b3.SUCCESS
	}

	return b3.FAILURE
}

// IsRage 是否处于狂暴状态
type IsRage struct {
	b3core.Condition
}

// Initialize 初始化
func (bev *IsRage) Initialize(setting *b3config.BTNodeCfg) {
	bev.Condition.Initialize(setting)
}

// OnTick 循环
func (bev *IsRage) OnTick(tick *b3core.Tick) b3.Status {
	pawn, b := tick.GetTarget().(*Pawn)
	if !b {
		log.Error("*************[IsRage] trans Pawn fail")
		return b3.ERROR
	}

	pawn.pushActionTorsionInfo(bev.GetName())

	if pawn.State.Raged {
		return b3.SUCCESS
	}

	return b3.FAILURE
}

// AttrLowerThan 属性值低于某个百分比
type AttrLowerThan struct {
	b3core.Condition
	AttrTypeKey string
}

// Initialize 初始化
func (bev *AttrLowerThan) Initialize(setting *b3config.BTNodeCfg) {
	bev.Condition.Initialize(setting)
	bev.AttrTypeKey = setting.GetPropertyAsString("AttrTypeKey")
}

// OnTick 循环
func (bev *AttrLowerThan) OnTick(tick *b3core.Tick) b3.Status {
	pawn, b := tick.GetTarget().(*Pawn)
	if !b {
		log.Error("*************[IsRage] trans Pawn fail")
		return b3.ERROR
	}

	pawn.pushActionTorsionInfo(bev.GetName())

	attrValue := pawn.getAttrValueByKey(bev.AttrTypeKey)
	blackBoardValue := pawn.getBlackboardValueByKey(bev.AttrTypeKey)
	switch bev.AttrTypeKey {
	case eB_Attr_HP_Per:
		if blackBoardValue != nil && attrValue.(float32) <= float32(blackBoardValue.(float64)) {
			return b3.SUCCESS
		}
	}

	return b3.FAILURE
}

// IsOverDrive 是否处于超载状态
type IsOverDrive struct {
	b3core.Condition
}

// Initialize 初始化
func (bev *IsOverDrive) Initialize(setting *b3config.BTNodeCfg) {
	bev.Condition.Initialize(setting)
}

// OnTick 循环
func (bev *IsOverDrive) OnTick(tick *b3core.Tick) b3.Status {
	pawn, b := tick.GetTarget().(*Pawn)
	if !b {
		log.Error("*************[IsOverDrive] trans Pawn fail")
		return b3.ERROR
	}

	pawn.pushActionTorsionInfo(bev.GetName())

	//处于超载状态
	if pawn.State.OverDrive {
		return b3.SUCCESS
	}

	return b3.FAILURE
}

// IsPartnerLoseHP 是否伙伴处于失血状态
type IsPartnerLoseHP struct {
	b3core.Condition
}

// Initialize 初始化
func (bev *IsPartnerLoseHP) Initialize(setting *b3config.BTNodeCfg) {
	bev.Condition.Initialize(setting)
}

// OnTick 循环
func (bev *IsPartnerLoseHP) OnTick(tick *b3core.Tick) b3.Status {
	pawn, b := tick.GetTarget().(*Pawn)
	if !b {
		log.Error("*************[IsPartnetLoseHP] trans Pawn fail")
		return b3.ERROR
	}

	pawn.pushActionTorsionInfo(bev.GetName())

	partnetList := pawn.GetPartnerList()
	for _, partner := range partnetList {
		if partner.IsAlive() && !pawn.cantBeSelect(partner) && partner.Attr.CurHP != partner.Attr.MaxHP {
			return b3.SUCCESS
		}
	}

	return b3.FAILURE
}

// IsUseableLebelsSkill 是否标签技能可用
type IsUseableLebelsSkill struct {
	b3core.Condition
	SkillLebelsKey string
}

// Initialize 初始化
func (bev *IsUseableLebelsSkill) Initialize(setting *b3config.BTNodeCfg) {
	bev.Condition.Initialize(setting)
	bev.SkillLebelsKey = setting.GetPropertyAsString("SkillLebelsKey")
}

// OnTick 循环
func (bev *IsUseableLebelsSkill) OnTick(tick *b3core.Tick) b3.Status {
	pawn, b := tick.GetTarget().(*Pawn)
	if !b {
		log.Error("*************[IsUseableLebelsSkill] trans Pawn fail")
		return b3.ERROR
	}

	pawn.pushActionTorsionInfo(bev.GetName())

	var lebelList []uint32
	var useableSkills []*_SkillItem

	value, ok := pawn.getBlackboardValueByKey(bev.SkillLebelsKey).(string)
	if !ok {
		log.Error("*************[IsUseableLebelsSkill] 外部黑板 SkillLebelsKey 没有配置")
		return b3.ERROR
	}
	if len(value) > 0 {
		for _, val := range strings.Split(value, "|") {
			lebel, err := strconv.Atoi(val)
			if err != nil {
				log.Error("*************[IsUseableLebelsSkill] 外部黑板 SkillLebelsKey 配置格式不对 正确格式应该是lebel1|lebel2   :", value)
				return b3.FAILURE
			}

			lebelList = append(lebelList, uint32(lebel))
		}
	}

	for _, lebel := range lebelList {
		skills := bev.GetUseablelebelSkills(lebel, pawn)
		useableSkills = append(useableSkills, skills...)
	}

	if len(useableSkills) > 0 {
		skillIndex := rand.Intn(len(useableSkills))
		tick.Blackboard.Set(curSkill, useableSkills[skillIndex], "", "")
		return b3.SUCCESS
	}

	return b3.FAILURE
}

// GetUseablelebelSkills 获取可用的标签技能列表		此处标签 仅仅是针对超能技
func (bev *IsUseableLebelsSkill) GetUseablelebelSkills(lebel uint32, pawn *Pawn) (result []*_SkillItem) {
	for _, skill := range pawn.GetSuperSkillList() {
		for _, lebelValue := range skill.Config.SkillLabel {
			if lebel == lebelValue && pawn.CanUseSkill(skill) {
				result = append(result, skill)
				break
			}
		}
	}

	return
}

// IsPositiveEnemy 是否正面boss
type IsPositiveEnemy struct {
	b3core.Condition
	PositiveBossAngleKey  string
	PositiveBossRadiusKey string
}

// Initialize 初始化
func (bev *IsPositiveEnemy) Initialize(setting *b3config.BTNodeCfg) {
	bev.Condition.Initialize(setting)
	bev.PositiveBossAngleKey = setting.GetPropertyAsString("PositiveBossAngleKey")
	bev.PositiveBossRadiusKey = setting.GetPropertyAsString("PositiveBossRadiusKey")
}

// OnTick 循环
func (bev *IsPositiveEnemy) OnTick(tick *b3core.Tick) b3.Status {
	pawn, b := tick.GetTarget().(*Pawn)
	if !b {
		log.Error("*************[IsPositiveBoss] trans Pawn fail")
		return b3.ERROR
	}

	pawn.pushActionTorsionInfo(bev.GetName())

	target, ok := getPawnBoard(tick, attackTarget)
	if !ok {
		log.Error("*************[IsPositiveBoss] getPawnBoard fail")
		return b3.FAILURE
	}

	angle, ok := pawn.getBlackboardValueByKey(bev.PositiveBossAngleKey).(float64)
	if !ok {
		log.Error("*************[IsPositiveBoss] 外部黑板 PositiveBossAngleKey 没有配置")
		return b3.ERROR
	}
	radius, ok := pawn.getBlackboardValueByKey(bev.PositiveBossRadiusKey).(float64)
	if !ok {
		log.Error("*************[IsPositiveBoss] 外部黑板 PositiveBossRadiusKey 没有配置")
		return b3.ERROR
	}
	//角度转化为弧度
	angle = angle * math.Pi / 180

	relativeAngle := CalcAngle(target.GetPos(), pawn.GetPos())
	//fmt.Println("+++++++++  ", math.Abs(relativeAngle-float64(target.GetAngle()))*180/math.Pi, angle*180/math.Pi/2)
	//fmt.Println("+++++++++  ", DistancePawn(target, pawn), (float32(radius) + pawn.Attr.CollisionRadius + target.Attr.CollisionRadius))
	if math.Abs(relativeAngle-float64(target.GetAngle())) <= angle/2 && DistancePawn(target, pawn) <= (float32(radius)+pawn.Attr.CollisionRadius+target.Attr.CollisionRadius) {
		return b3.SUCCESS
	}

	return b3.FAILURE
}

// InAOERange 处于某方阵营的AOE范围内
type InAOERange struct {
	b3core.Condition
	IsEnemyCamp bool
}

// Initialize 初始化
func (bev *InAOERange) Initialize(setting *b3config.BTNodeCfg) {
	bev.Condition.Initialize(setting)
	bev.IsEnemyCamp = setting.GetPropertyAsBool("IsEnemyCamp")
}

// OnTick 循环
func (bev *InAOERange) OnTick(tick *b3core.Tick) b3.Status {
	pawn, b := tick.GetTarget().(*Pawn)
	if !b {
		log.Error("*************[IsPartnetLoseHP] trans Pawn fail")
		return b3.ERROR
	}

	pawn.pushActionTorsionInfo(bev.GetName())

	//默认己方阵营
	camp := pawn.GetCamp()
	if bev.IsEnemyCamp {
		camp = pawn.GetEmemyCamp()
	}

	if pawn.InCampAlertAOERange(camp) {
		return b3.SUCCESS
	}

	return b3.FAILURE
}
