package internal

import (
	"Cinder/Base/linemath"
	"math"
)

const (
	faultDis                  = 0   //用于容错的距离
	pawnStationNumberOfTurns  = 1   //pawn周边站位点圈数
	pawnStationNumberOfRadius = 1   //pawn周边站位点碰撞半径
	pawnStationOldOrNewDis    = 0.5 //新旧就位点距离阈值
)

//resultPoint 就位点数据
type resultPoint struct {
	station    *stationData
	stationPos linemath.Vector2
}

//stationData 位置信息
type stationData struct {
	stationIndex int     //位置索引
	stationAngle float64 //该位置对应角度
	truns        int     //对应圈数

	radius float32 //面向相对半径
	pawn   *Pawn
}

type _PawnFreeStation struct {
	stationPoint     map[int]*stationData //站位点数据
	outOfPositivePos map[int]*stationData //背向点站位点数据
	owner            *Pawn
}

func (f *_PawnFreeStation) init(pawn *Pawn) {
	f.stationPoint = map[int]*stationData{}
	f.outOfPositivePos = map[int]*stationData{}
	f.owner = pawn
}

//getFreeNearestStation 获取空闲最新站位  技能射程参数
func (f *_PawnFreeStation) getFreeNearestStation(target *Pawn, baseAttackDis, attackDisMin, attackDisMax float32, isAI bool) (station stationData, stationPos linemath.Vector2, ok bool) {
	//如果是自己找自己的攻击点，就是原地不动
	if target.UID == f.owner.UID {
		return stationData{pawn: target}, target.GetPos(), true
	}

	var oldPos linemath.Vector2
	var oldStation *stationData
	var beginIndex int
	var resultList []resultPoint

	//周边圈数
	for turns := 1; turns <= pawnStationNumberOfTurns; {

		stationPointRadius := f.owner.Attr.CollisionRadius + float32(turns)*(pawnStationNumberOfRadius+faultDis)
		angle := math.Asin(float64(pawnStationNumberOfRadius / stationPointRadius))
		maxIndexNum := int(2 * math.Pi / angle)

		////当前圈的外部空闲位置数量
		//maxIndexNum = int(2 * math.Pi / angle)

		//根据计算出来的旋转角度 得出一圈应该有几个点
		for angleIndex := 0; angleIndex < maxIndexNum; beginIndex++ {
			curAngle := angle * float64(angleIndex)
			angleIndex++

			pos := linemath.Vector2{
				X: float32(math.Cos(curAngle)),
				Y: float32(math.Sin(curAngle)),
			}
			realRadius := f.owner.Attr.CollisionRadius + float32(turns)*(target.Attr.CollisionRadius+faultDis)
			pos = pos.Mul(realRadius).Add(f.owner.GetPos())

			//触发边界
			if f.owner.Scene.checkCircleShapeOverlapWorldBoundary(float64(pos.X), float64(pos.Y), float64(target.Attr.CollisionRadius)) {
				continue
			}

			//点是否处于战场内
			if !f.owner.Scene.checkPointInWorldBoundary(pos) {
				continue
			}

			//碰撞检测 此处位置已经被占
			targetList := f.owner.Scene.overlapCircleShape(float64(pos.X), float64(pos.Y), float64(target.Attr.CollisionRadius))
			if len(targetList) > 0 {
				//是否需要过滤当前位置(被他人占用，且不是相切。 需要过滤)
				var filterRequired bool

				//检测如果碰撞其他对象正好相切可以站位
				for _, val := range targetList {
					//是攻击者自己，不管
					if val.UID == target.UID {
						continue
					}

					//是被攻击者自己
					if val.UID == f.owner.UID {
						continue
					}

					if Distance(val.GetPos(), pos) < val.Attr.CollisionRadius+target.Attr.CollisionRadius {
						filterRequired = true
						break
					}
				}

				if filterRequired {
					continue
				}
			}

			ownObj := f.stationPoint[beginIndex]
			//已经被占有，且对象存活
			if ownObj != nil && ownObj.pawn.IsAlive() {
				if ownObj.pawn.UID == target.UID {
					oldStation = ownObj
					oldPos = pos
				}

				continue
			}

			if !f.checkoutAlreadyStayInPos(pos, target, f.stationPoint) {
				continue
			}

			resultList = append(resultList, resultPoint{
				&stationData{
					stationIndex: beginIndex,
					stationAngle: curAngle,
					truns:        turns,
					pawn:         target,
					radius:       realRadius,
				},
				pos,
			})
		}

		turns++
	}

	//攻击距离预处理
	for i := len(resultList) - 1; i >= 0; i-- {
		find, realPos := f.stationPreprocessing(resultList[i].stationPos, target, baseAttackDis, attackDisMin, attackDisMax, isAI)
		if !find {
			resultList = append(resultList[:i], resultList[i+1:]...)
			continue
		}

		resultList[i].stationPos = realPos
	}

	//老位置攻击距离再次判定
	if oldStation != nil {
		find, realPos := f.stationPreprocessing(oldPos, target, baseAttackDis, attackDisMin, attackDisMax, isAI)
		if find {
			oldPos = realPos
		} else {
			oldStation = nil
		}
	}

	minDis := float32(-1)
	//检测可占用点 处于攻击范围内 且距离自己最近
	for _, val := range resultList {
		dis := Distance(val.stationPos, target.GetPos())

		if minDis == -1 {
			minDis = dis
			station = *val.station
			stationPos = val.stationPos
			ok = true
			continue
		}

		if minDis > dis {
			minDis = dis
			station = *val.station
			stationPos = val.stationPos
			ok = true
		}
	}

	if oldStation != nil {
		//找到了新位置
		if ok {
			//那么和新位置 距离对比
			if Distance(oldPos, target.GetPos()) <= Distance(stationPos, target.GetPos())+pawnStationOldOrNewDis {
				stationPos = oldPos
				station = *oldStation
			}
		} else {
			stationPos = oldPos
			station = *oldStation
			ok = true
		}
	}

	//fmt.Println("++++++++++++ getFreeNearestStation ", f.owner.UID, target.UID, station, stationPos)
	return
}

//stationPreprocessing 就位点 对于攻击距离 预处理
func (f *_PawnFreeStation) stationPreprocessing(stationPos linemath.Vector2, target *Pawn, baseAttackDis, attackDisMin, attackDisMax float32, isAI bool) (bool, linemath.Vector2) {
	v2 := stationPos.Sub(target.GetPos())
	minLength := v2.Len()
	if stationPos.IsEqual(target.GetPos()) {
		minLength = 0
		v2 = f.owner.GetPos().Sub(target.GetPos())
	}

	var isBack bool
	//当就位点位于 被攻击者相对后方的时候
	if Distance(target.GetPos(), stationPos) >= Distance(target.GetPos(), f.owner.GetPos()) {
		if isAI {
			//minLength = Distance(target.GetPos(), f.owner.GetPos())
			isBack = true
		}

		if !isAI {
			return true, stationPos
		}
	}

	if minLength >= attackDisMax {
		dis := minLength - attackDisMax + baseAttackDis
		if isBack {
			dis = minLength + attackDisMax + baseAttackDis
		}
		stationPos = v2.Normalized().Mul(dis).Add(target.GetPos())
	} else if attackDisMin > 0 && minLength <= attackDisMin {
		stationPos = v2.Normalized().Mul(-1).Mul(attackDisMin - minLength + baseAttackDis).Add(target.GetPos())
	} else if DistancePawn(f.owner, target) >= f.owner.Attr.CollisionRadius+target.Attr.CollisionRadius {
		if target.OverlapCircleShape(target.GetPos()) {
			stationPos = target.GetPos()
		} else {
			//fmt.Println("自己当前位置有碰撞")
		}
	} else {
		//此处表示 双方部分重叠，需要双方拉开。 就位点是可以保证双方不重叠的。
		//fmt.Println("==========")
	}

	//触发边界
	if f.owner.Scene.checkCircleShapeOverlapWorldBoundary(float64(stationPos.X), float64(stationPos.Y), float64(target.Attr.CollisionRadius)) {
		return false, stationPos
	}

	//点是否处于战场内
	if !f.owner.Scene.checkPointInWorldBoundary(stationPos) {
		return false, stationPos
	}

	//站位点处于攻击范围以外
	realDis := float32(math.Round(float64(Distance(stationPos, f.owner.GetPos()))))
	if realDis < (attackDisMin+target.Attr.CollisionRadius) || realDis > (attackDisMax+target.Attr.CollisionRadius) {
		return false, stationPos
	}

	return true, stationPos
}

//setStationData 设置被占者信息
func (f *_PawnFreeStation) setStationData(station stationData) {
	if station.pawn.UID == f.owner.UID {
		return
	}

	//旧的被占数据删除
	for key, val := range f.stationPoint {
		if val.pawn.UID == station.pawn.UID {
			delete(f.stationPoint, key)
		}
	}
	f.stationPoint[station.stationIndex] = &station

	f.removePartnetStationData(station.pawn)
}

//removePartnetStationData 移除队友被占位置 目标不能同时占有两个怪物的位置
func (f *_PawnFreeStation) removePartnetStationData(target *Pawn) {
	for _, partner := range f.owner.GetPartnerList() {
		//排除死亡队友或者自己
		if !partner.IsAlive() || partner.UID == f.owner.UID {
			continue
		}

		partner.removeStationData(target)
	}
}

//removeStationData 移除位置占有信息
func (f *_PawnFreeStation) removeStationData(target *Pawn) {
	for key, val := range f.stationPoint {
		if val.pawn.UID == target.UID {
			delete(f.stationPoint, key)
		}
	}
}

//checkFilterAngle 过滤角度范围检测
func (f *_PawnFreeStation) checkFilterAngle(minAngle, maxAngle, angle float32) bool {

	if minAngle < 0 {
		minAngle += 2 * math.Pi

		if minAngle < 0 {
			return false
		}
	}

	if maxAngle >= 2*math.Pi {
		maxAngle -= 2 * math.Pi

		if maxAngle >= 2*math.Pi {
			return false
		}
	}

	if minAngle < maxAngle {
		if angle >= minAngle && angle <= maxAngle {
			return false
		}
	} else {
		if (angle >= minAngle && angle <= 2*math.Pi) || (angle >= 0 && angle <= maxAngle) {
			return false
		}
	}

	return true
}

//getOutOfPositiveRangePos 获取当前自己面向范围以外的点
func (f *_PawnFreeStation) getOutOfPositiveRangePos(radius, angle float32, target *Pawn) (targetPos linemath.Vector2, ok bool) {
	var oldPos linemath.Vector2
	var oldStation *stationData
	var station stationData
	var resultList []resultPoint

	baseAngle := math.Asin(float64(pawnStationNumberOfRadius / radius))
	maxIndexNum := int(2 * math.Pi / baseAngle)

	for angleIndex := 0; angleIndex < maxIndexNum; angleIndex++ {
		curAngle := baseAngle * float64(angleIndex)

		if !f.checkFilterAngle(f.owner.GetAngle()-angle/2, f.owner.GetAngle()+angle/2, float32(curAngle)) {
			continue
		}

		pos := linemath.Vector2{
			X: float32(math.Cos(curAngle)),
			Y: float32(math.Sin(curAngle)),
		}
		pos = pos.Mul(radius).Add(f.owner.GetPos())

		//触发边界
		if f.owner.Scene.checkCircleShapeOverlapWorldBoundary(float64(pos.X), float64(pos.Y), float64(target.Attr.CollisionRadius)) {
			continue
		}

		//点是否处于战场内
		if !f.owner.Scene.checkPointInWorldBoundary(pos) {
			continue
		}

		//碰撞检测 此处位置已经被占
		targetList := f.owner.Scene.overlapCircleShape(float64(pos.X), float64(pos.Y), float64(target.Attr.CollisionRadius))
		if len(targetList) > 0 {
			//是否需要过滤当前位置(被他人占用，且不是相切。 需要过滤)
			var filterRequired bool

			//检测如果碰撞其他对象正好相切可以站位
			for _, val := range targetList {
				//是攻击者自己，不管
				if val.UID == target.UID {
					continue
				}

				//是被攻击者自己
				if val.UID == f.owner.UID {
					continue
				}

				if Distance(val.GetPos(), pos) < val.Attr.CollisionRadius+target.Attr.CollisionRadius {
					filterRequired = true
					break
				}
			}

			if filterRequired {
				continue
			}
		}

		ownObj := f.outOfPositivePos[angleIndex]
		//已经被占有，且对象存活
		if ownObj != nil && ownObj.pawn.IsAlive() {
			if ownObj.pawn.UID == target.UID {
				oldStation = ownObj
				oldPos = pos
			}

			continue
		}

		if !f.checkoutAlreadyStayInPos(pos, target, f.outOfPositivePos) {
			continue
		}

		resultList = append(resultList, resultPoint{
			&stationData{
				stationIndex: angleIndex,
				stationAngle: curAngle,
				pawn:         target,
				radius:       radius,
			},
			pos,
		})
	}

	minDis := float32(-1)
	//选择距离最近的点
	for _, val := range resultList {
		dis := Distance(val.stationPos, target.GetPos())

		if minDis == -1 {
			minDis = dis
			station = *val.station
			targetPos = val.stationPos
			ok = true
			continue
		}

		if minDis > dis {
			minDis = dis
			station = *val.station
			targetPos = val.stationPos
			ok = true
		}
	}

	if oldStation != nil {
		//找到了新位置
		if ok {
			//那么和新位置 距离对比
			if Distance(oldPos, target.GetPos())+pawnStationOldOrNewDis <= Distance(targetPos, target.GetPos()) {
				targetPos = oldPos
				station = *oldStation
			}
		} else {
			targetPos = oldPos
			station = *oldStation
			ok = true
		}
	}

	if ok {
		f.setOutOfPositiveStationData(station)
	}

	return
}

//setOutOfPositiveStationData 设置背向点站位点数据
func (f *_PawnFreeStation) setOutOfPositiveStationData(station stationData) {
	if station.pawn.UID == f.owner.UID {
		return
	}

	//旧的被占数据删除
	for key, val := range f.outOfPositivePos {
		if val.pawn.UID == station.pawn.UID {
			delete(f.outOfPositivePos, key)
		}
	}
	f.outOfPositivePos[station.stationIndex] = &station

	f.RemovePartnetOutOfPositiveStationData(station.pawn)
}

//RemovePartnetOutOfPositiveStationData 移除队友被占位置 目标不能同时占有两个怪物的位置
func (f *_PawnFreeStation) RemovePartnetOutOfPositiveStationData(target *Pawn) {
	for _, partner := range f.owner.GetPartnerList() {
		//排除死亡队友或者自己
		if !partner.IsAlive() || partner.UID == f.owner.UID {
			continue
		}

		partner.OutOfPositiveStationData(target)
	}
}

//OutOfPositiveStationData 移除位置占有信息
func (f *_PawnFreeStation) OutOfPositiveStationData(target *Pawn) {
	for key, val := range f.outOfPositivePos {
		if val.pawn.UID == target.UID {
			delete(f.outOfPositivePos, key)
		}
	}
}

//checkoutAlreadyStayInPos 已经占有位置检测
func (f *_PawnFreeStation) checkoutAlreadyStayInPos(pos linemath.Vector2, target *Pawn, data map[int]*stationData) bool {
	for key, val := range data {

		//删除废弃站位信息
		if !val.pawn.IsAlive() {
			delete(data, key)
			continue
		}

		if val.pawn.UID == target.UID {
			continue
		}

		outOfPositivePos := linemath.Vector2{
			X: float32(math.Cos(val.stationAngle)),
			Y: float32(math.Sin(val.stationAngle)),
		}
		outOfPositivePos = outOfPositivePos.Mul(val.radius).Add(f.owner.GetPos())
		//fmt.Println("++++++++ ", Distance(outOfPositivePos, pos))
		if Distance(outOfPositivePos, pos) < target.Attr.CollisionRadius+val.pawn.Attr.CollisionRadius {
			return false
		}
	}

	return true
}
