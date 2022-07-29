package drop

import (
	"Daisy/Const"
	"Daisy/Data"
	"Daisy/Proto"
	log "github.com/cihub/seelog"
	"math"
	"math/rand"
)

//MakeEquipment 生成装备数据
func (drop *Drop) MakeEquipment(item *Proto.Item, args ...uint32) *Proto.Item {
	//log.Debug("装备唯一id=", item.Base.ID)

	item.EquipmentData = &Proto.Equipment{}

	var MLevel, MType, locky uint32
	if len(args) == 3 {
		MLevel = args[0]
		MType = args[1]
		locky = args[2]
	} else {
		MLevel = 1
		MType = 0
		locky = 10 //默认
	}

	config, ok := Data.GetEquipConfig().EquipMent_ConfigItems[item.Base.ConfigID]
	if ok == true {
		//品质
		var _quality uint32
		if config.Quality != 0 {
			_quality = config.Quality
			//log.Debug("取配置表品质:", config.Quality)
		} else {
			_quality = drop.CountEquipQulity(config.QualityLevel, config.ItemLevel, MLevel, MType, uint32(locky))
		}

		if _quality == 0 {
			log.Info("生成装备出错，品质计算错误：0， id=", item.Base.ID)
			return item
		}

		if _quality == Const.Quality_Dullgold { //粉色
			//取另外的装备id
			if config.GoldEquipID != 0 {
				_config, _ok := Data.GetEquipConfig().EquipMent_ConfigItems[config.GoldEquipID]
				if _ok == true {
					config = _config
					//log.Debug("粉色装备，切换id = ", config.ID)
				} else {
					//log.Debug("粉色装备，切换失败，找不到 id = ", config.ID)
				}
			}
		} else if _quality == Const.Quality_Green {
			if config.SuitEquipID != 0 {
				_config, _ok := Data.GetEquipConfig().EquipMent_ConfigItems[config.SuitEquipID]
				if _ok == true {
					config = _config
					//log.Debug("橙色装备，切换id = ", config.ID)
				} else {
					//log.Debug("橙色装备，切换失败，找不到 id = ", config.ID)
				}
			}
		}

		item.Base.ConfigID = config.ID
		item.EquipmentData.Quality = _quality
		//log.Debug("生成装备品质:", item.EquipmentData.Quality)
		//词缀
		item.EquipmentData.Affixes = drop.CreateEquipAffixes(item.EquipmentData.Quality, item)
		//孔
		socketNum := drop.CountEquipHoleNum(config.ItemLevel, config.MinSocketNum, config.MaxSocketNum)
		//log.Debug("孔数量=", socketNum)

		item.EquipmentData.Socket = make([]*Proto.Item, 0)
		for i := uint32(0); i < socketNum; i++ {
			item.EquipmentData.Socket = append(item.EquipmentData.Socket, &Proto.Item{Base: &Proto.ItemBase{}})
		}
		//珍品
		if config.IsPrecious == false {
			item.EquipmentData.IsPrecious = drop.CountEquipIsPrecious(item.EquipmentData.Quality)
			//log.Debug("随机珍品,是否珍品：%t", item.EquipmentData.IsPrecious)
		} else {
			item.EquipmentData.IsPrecious = true
			//log.Debug("必然获得珍品")
		}

		item.EquipmentData.Score = drop.Score(item)
		//log.Debug("获得装备评分", item.EquipmentData.Score)
	}

	return item
}

//CountEquipHoleNum 获得装备孔数量
func (drop *Drop) CountEquipHoleNum(itemLevel, min, max uint32) uint32 {
	var addNum, num uint32
	for _, v := range Data.GetEquipConfig().Hole_ConfigItems {
		if itemLevel >= v.MinLevel && itemLevel <= v.MaxLevel {
			addNum = v.HoleAdd
			break
		}
	}

	if min+addNum > max {
		num = max
	} else {
		num = min + addNum
	}
	//log.Debug(fmt.Sprintf("随机孔: min=%d, max=%d", min, num))
	return uint32(RandBetween(int(min), int(num)))
}

//RandBetween 返回随机值 [min,max]
//panic min > max
func RandBetween(min, max int) int {
	return rand.Intn(max-min+1) + min
}
func GetRandValue(precision uint32, min, max float32) float32{
	var value float32
	switch precision {
	case Const.AffixPrecision_0:
		value = float32(RandBetween(int(min), int(max)))
	case Const.AffixPrecision_2:
		value = float32(RandBetween(int(min*100), int(max*100)))/100.0
	case Const.AffixPrecision_4:
		value = float32(RandBetween(int(min*10000), int(max*10000)))/10000.0
	}
	return value
}
//CountEquipIsPrecious 判断是否能成为珍品
func (drop *Drop) CountEquipIsPrecious(quality uint32) bool {
	return Data.GetEquipConfig().Quality_ConfigItems[quality].GemRate >= uint32(rand.Intn(10000))
}

//CountEquipQulity 计算装备品质
//QLevel 装备品质等级
//MLevel 怪物等级
//MType 怪物类型 0：普通 1：精英 2：boss
//ILevel 装备等级
//lucky 人物幸运值
func (drop *Drop) CountEquipQulity(QLevel, ILevel, MLevel, MType, lucky uint32) uint32 {
	//log.Debug(fmt.Sprintf("装备品质等级=%d, 怪物等级=%d, 装备等级=%d, 怪物类型=%d, 人物幸运值=%d", QLevel, MLevel, ILevel, MType, lucky))
	var quality map[uint32]uint32
	quality = make(map[uint32]uint32, len(Data.GetEquipConfig().Quality_ConfigItems))

	config := Data.GetEquipConfig().Quality_ConfigItems
	//获取基础品质概率

	switch MType {
	case 0:
		for k, v := range config {
			quality[k] = v.Common
		}
	case 1:
		for k, v := range config {
			quality[k] = v.Elit
		}
	case 2:
		for k, v := range config {
			quality[k] = v.Boss
		}
	}

	//log.Debugf("基础品质概率=%d,%d,%d,%d,%d,%d", quality[1], quality[2], quality[3], quality[4], quality[5], quality[6])
	//计算幸运概率加成
	for k2, v2 := range config {
		if v2.LuckyCoe == 0 {
			continue
		}
		quality[k2] += lucky * v2.LuckyCoe * 10 / (lucky + v2.LuckyCoe)
	}

	//log.Debug("计算幸运概率加成后的品质概率=", quality)
	//计算 品质等级和怪物等级加成
	//(1-ABS（(MLVL-QLVL)）/(MLVL*QLVL))/100
	add := uint32(100 - math.Abs((float64(MLevel)-float64(QLevel))*float64(100)/float64(MLevel*QLevel)))

	//log.Debug("品质等级和怪物等级加成=", add)
	//概率最高取10000，剔除低品质的概率
	_rand := rand.Intn(10000)
	//log.Debug("品质随机值 _rand =", _rand)

	for i := len(quality); i > 0; i-- {
		//加上等级加成
		if int(quality[uint32(i)]+add) >= _rand {
			return uint32(i)
		} else {
			_rand -= int(quality[uint32(i)] + add)
		}
	}
	return 0
}

func (drop *Drop) CreateEquipAffixes(quality uint32, item *Proto.Item) []*Proto.AffixData {
	affixList := make([]*Proto.AffixData, 0)
	config, has := Data.GetEquipConfig().Quality_ConfigItems[quality]
	if has == false {
		log.Debug("品质数据找不到,Quality=", quality)
		return affixList
	}
	equipConfig, ok := Data.GetEquipConfig().EquipMent_ConfigItems[item.Base.ConfigID]
	if ok == false {
		log.Debug("装备数据找不到,Quality=", quality)
		return affixList
	}

	//固定词缀
	if len(equipConfig.FixedAffix) != 0 {
		for _, v := range equipConfig.FixedAffix {
			affixConfig, ok3 := Data.GetEquipConfig().EquipAffix_ConfigItems[v]
			if ok3 == true {
				data := Proto.AffixData{}
				data.AffixID = affixConfig.AffixID
				data.AffixEffectType = affixConfig.AffixEffectType
				if data.AffixEffectType == Const.AffixEffectType_Attr {
					data.PropertyID = affixConfig.AffixAttID
					data.Value = GetRandValue(affixConfig.AffixPrecision, affixConfig.MinAffixAttValue, affixConfig.MaxAffixAttValue)
				} else {
					data.PropertyID = affixConfig.BuffID
				}

				data.AffixParam = affixConfig.AffixParam
				data.Type = affixConfig.AffixPlace
				data.Const = true
				//log.Debug(fmt.Sprintf("固定词缀: id=%d, Value=%f", data.AffixID, data.Value))
				affixList = append(affixList, &data)
			}
		}
	}

	groupMap := make(map[uint32]bool)
	//词缀等级= max(QualityLevel, ItemLevel)
	var maxALvl uint32
	if equipConfig.QualityLevel > equipConfig.ItemLevel {
		maxALvl = equipConfig.QualityLevel
	} else {
		maxALvl = equipConfig.ItemLevel
	}
	if maxALvl > 99 {
		maxALvl = 99
	}

	//log.Debug("词缀 maxALvl=", maxALvl)

	//前缀
	//log.Debug(fmt.Sprintf("随机词缀: 前缀随机 %d次", config.EquipPrefixNum))
	loop := 1000
	num := uint32(0)

	for {

		if num >= config.EquipPrefixNum{
			break
		}
		if loop--; loop == 0 {
			log.Infof("装备:%s,%s,%u,没有随机到前缀", item.Base.ID, equipConfig.Name, equipConfig.ID)
			break
		}

		_, list := drop.Drop(equipConfig.AffixDropID[0], 0, 0)
		if len(list) == 0 {
			//log.Debug(fmt.Sprintf("随机词缀: 没有掉落任何词缀, dropID = %d", equipConfig.AffixDropID))
			continue
		}
		affixConfig, ok3 := Data.GetEquipConfig().EquipAffix_ConfigItems[list[0].MaterialId]
		if ok3 == false {
			log.Debugf("随机词缀：词缀表没有:%d", list[0].MaterialId)
			continue
		}
		//log.Debug(fmt.Sprintf("随机词缀: 掉落词缀, ID = %d", affixConfig.AffixID))

		//前缀
		if affixConfig.AffixPlace != Const.AffixPlace_Front {
			//log.Debug(fmt.Sprintf("随机词缀: 不是前缀, ID = %d", affixConfig.AffixID))
			continue
		}

		//品质判断
		pass := false
		for _, v := range affixConfig.QualityLimit {
			if v == 0 || v == quality {
				pass = true
				break
			}
		}
		if pass == false {
			//log.Debug(fmt.Sprintf("随机词缀: 词缀品质不符, %d， config=", quality), equipConfig.Quality)
			continue
		}

		//部位判断
		pass = false
		for _, v := range affixConfig.PositionLimit {
			if v == 0 || v == equipConfig.Position {
				pass = true
				break
			}
		}
		if pass == false {
			//log.Debug(fmt.Sprintf("随机词缀: 词缀部位不符, %d , config=", affixConfig.PositionLimit), equipConfig.Position)
			continue
		}

		//等级判断
		if affixConfig.AffixLevel > maxALvl {
			//log.Debug(fmt.Sprintf("随机词缀: 词缀等级太高, %d , %d", affixConfig.AffixLevel, maxALvl))
			continue
		}

		//group 判断
		if affixConfig.AffixGroup != 0 {
			if _, found := groupMap[affixConfig.AffixGroup]; found == true {
				//log.Debug(fmt.Sprintf("随机词缀: 已经有同组, %d ", affixConfig.AffixGroup))
				continue
			} else {
				groupMap[affixConfig.AffixGroup] = true
			}
		}

		data := Proto.AffixData{}
		data.AffixID = affixConfig.AffixID
		data.AffixEffectType = affixConfig.AffixEffectType
		if data.AffixEffectType == Const.AffixEffectType_Attr {
			data.PropertyID = affixConfig.AffixAttID
			data.Value = GetRandValue(affixConfig.AffixPrecision, affixConfig.MinAffixAttValue, affixConfig.MaxAffixAttValue)
		} else {
			data.PropertyID = affixConfig.BuffID
		}
		data.AffixParam = affixConfig.AffixParam
		data.Type = affixConfig.AffixPlace
		data.Const = false

		//log.Debug(fmt.Sprintf("随机词缀: id=%d, Value=%f", data.AffixID, data.Value))
		affixList = append(affixList, &data)
		num++
	}

	//后缀
	//log.Debug(fmt.Sprintf("随机词缀: 后缀随机 %d次", config.EquipPrefixNum))
	loop = 1000
	num = uint32(0)
	for {
		if num >= config.EquipSuffixNum{
			break
		}
		if loop--; loop == 0 {
			log.Infof("装备:%s,%s,%u,没有随机到后缀", item.Base.ID, equipConfig.Name, equipConfig.ID)
			break
		}

		var dropID uint32
		if len(equipConfig.AffixDropID) > 1 {
			dropID = equipConfig.AffixDropID[1]
		}
		_, list := drop.Drop(dropID, 0, 0)
		if len(list) == 0 {
			//log.Debug(fmt.Sprintf("随机词缀: 没有掉落任何词缀, dropID = %d", equipConfig.AffixDropID))
			continue
		}
		affixConfig, ok3 := Data.GetEquipConfig().EquipAffix_ConfigItems[list[0].MaterialId]
		if ok3 == false {
			//log.Debugf("随机词缀：词缀表没有:%d", list[0].MaterialId)
			continue
		}
		//log.Debug(fmt.Sprintf("随机词缀: 掉落词缀, ID = %d", affixConfig.AffixID))

		//后缀
		if affixConfig.AffixPlace != Const.AffixPlace_Back {
			//log.Debug(fmt.Sprintf("随机词缀: 不是后缀, ID = %d", affixConfig.AffixID))
			continue
		}

		//品质判断
		pass := false
		for _, v := range affixConfig.QualityLimit {
			if v == 0 || v == quality {
				pass = true
				break
			}
		}
		if pass == false {
			//log.Debug(fmt.Sprintf("随机词缀: 词缀品质不符, %d， config=", quality), equipConfig.Quality)
			continue
		}

		//部位判断
		pass = false
		for _, v := range affixConfig.PositionLimit {
			if v == 0 || v == equipConfig.Position {
				pass = true
				break
			}
		}
		if pass == false {
			//log.Debug(fmt.Sprintf("随机词缀: 词缀部位不符, %d , config=", affixConfig.PositionLimit), equipConfig.Position)
			continue
		}

		//等级判断
		if affixConfig.AffixLevel > maxALvl {
			//log.Debug(fmt.Sprintf("随机词缀: 词缀等级太高, %d , %d", affixConfig.AffixLevel, maxALvl))
			continue
		}

		//group 判断
		if affixConfig.AffixGroup != 0 {
			if _, found := groupMap[affixConfig.AffixGroup]; found == true {
				//log.Debug(fmt.Sprintf("随机词缀: 已经有同组, %d ", affixConfig.AffixGroup))
				continue
			} else {
				groupMap[affixConfig.AffixGroup] = true
			}
		}

		data := Proto.AffixData{}
		data.AffixID = affixConfig.AffixID
		data.AffixParam = affixConfig.AffixParam
		data.AffixEffectType = affixConfig.AffixEffectType
		if data.AffixEffectType == Const.AffixEffectType_Attr {
			data.PropertyID = affixConfig.AffixAttID
			data.Value = GetRandValue(affixConfig.AffixPrecision, affixConfig.MinAffixAttValue, affixConfig.MaxAffixAttValue)
		} else {
			data.PropertyID = affixConfig.BuffID
		}

		data.Type = affixConfig.AffixPlace
		data.Const = false
		//log.Debug(fmt.Sprintf("随机词缀: id=%d, Value=%f", data.AffixID, data.Value))
		affixList = append(affixList, &data)
		num++
	}
	return affixList
}

func (drop *Drop) Score(item *Proto.Item) uint32 {
	var score float32
	for _, v := range item.EquipmentData.Affixes {
		if v.AffixEffectType == Const.AffixEffectType_Attr {
			//属性
			config, ok := Data.GetEquipConfig().AttEnumEration_ConfigItems[v.PropertyID]
			if ok == false || config.ScoreParam == 0.0 {
				continue
			}
			score += v.Value / config.ScoreParam
		} else if v.AffixEffectType == Const.AffixEffectType_Buff {
			config, ok := Data.GetEquipConfig().EquipAffix_ConfigItems[v.AffixID]
			if ok == false || config.BuffAffixScore == 0 {
				continue
			}
			score += float32(config.BuffAffixScore)
		}
	}
	return uint32(score)
}
