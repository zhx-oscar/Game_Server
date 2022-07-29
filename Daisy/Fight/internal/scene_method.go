package internal

import (
	"Cinder/Base/linemath"
	"Daisy/Fight/internal/conf"
	"Daisy/Fight/internal/log"
	"Daisy/Proto"
	"errors"
	"fmt"
	"math"
	"unsafe"
)

// generateUID 生成唯一ID
func (scene *Scene) generateUID() uint32 {
	scene.uidSrc++

	if 0 == scene.uidSrc {
		scene.uidSrc++
	}

	return scene.uidSrc
}

// timeToFrames 时间转换为帧数（毫秒）
func (scene *Scene) timeToFrames(ms uint32) uint32 {
	frames := ms * scene.secFrames / 1000
	return frames
}

// framesToTime 帧数转换为时间（毫秒）
func (scene *Scene) framesToTime(frames uint32) uint32 {
	ms := frames * 1000 / scene.secFrames
	return ms
}

// SimulatorMode 是否是模拟器模式
func (scene *Scene) SimulatorMode() bool {
	return scene.Info.SimulatorMode
}

// TestMode 是否是测试模式
func (scene *Scene) TestMode() bool {
	return scene.Info.TestMode
}

//Summon 召唤接口
func (scene *Scene) Summon(info *PawnInfo) (*Pawn, error) {
	summonPawn := &Pawn{}
	err := summonPawn.init(scene, info)
	if err != nil {
		return nil, err
	}

	summonPawn.Info.IsSummon = true

	err = scene.putPawn(summonPawn)
	if err != nil {
		return nil, err
	}

	return summonPawn, nil
}

func (scene *Scene) GetBackgroundPawn(camp Proto.Camp_Enum) *Pawn {
	return scene.formationList[camp].BackgroundPawn
}

// putPawn 放置pawn
func (scene *Scene) putPawn(pawn *Pawn) error {
	if pawn == nil {
		return errors.New("nil pawn")
	}

	//场景中战斗对象上限检测
	if len(scene.pawnList) > conf.ScenePawnCountMax {
		return errors.New("scene pawnlist Over limit")
	}

	// pawn加入场景
	scene.pawnList = append(scene.pawnList, pawn)

	// 阵型
	formation := scene.formationList[pawn.Info.Camp]

	// pawn加入阵型
	formation.PawnList = append(formation.PawnList, pawn)

	// 刷新生命期
	pawn.RefreshLifeTime()

	//同步客户端召唤消息
	if pawn.Info.IsSummon {
		scene.PushAction(&Proto.SummonPawn{SelfId: pawn.UID})
	}

	// 记录客户端日志
	scene.PushDebugInfo(func() string {
		return fmt.Sprintf("创建${PawnID:%d}，类型：%v，等级：%d，阵营：%v，位置：%v，属性ID：%d，成长后属性：%+v,%+v",
			pawn.UID,
			pawn.Info.Type,
			pawn.Info.Level,
			pawn.Info.Camp,
			pawn.GetPos(),
			pawn.Info.ConfigId,
			*pawn.Attr._Attr.PawnAttr,
			*pawn.Attr._Attr.ExtendAttr)
	})

	// 记录日志
	log.Debugf("战场ID：%v，指令：放置Pawn，时间(ms)：%v，PawnID：%v，类型：%v，配置ID：%v，等级：%v，阵营：%v，位置：%v，基础属性：%+v，扩展属性：%+v",
		uintptr(unsafe.Pointer(scene)),
		scene.NowTime,
		pawn.UID,
		pawn.Info.Type,
		pawn.Info.ConfigId,
		pawn.Info.Level,
		pawn.Info.Camp,
		pawn.GetPos(),
		*pawn.Attr._Attr.PawnAttr,
		*pawn.Attr._Attr.ExtendAttr)

	scene.AddHaloMember(pawn)

	return nil
}

//buildSceneSizeByBoundaryPoints 通过边界点构建场景边界
func (scene *Scene) buildSceneSizeByBoundaryPoints() (x, y float64) {
	//构建战斗场景长宽用于构建 box2d场景世界
	var xmin, xmax, ymin, ymax float32
	for _, val := range scene.Info.BoundaryPoints {
		if val.X > xmax {
			xmax = val.X
		}

		if val.X < xmin {
			xmin = val.X
		}

		if val.Y > ymax {
			ymax = val.Y
		}

		if val.Y < ymin {
			ymin = val.Y
		}
	}

	x = float64(xmax - xmin)
	y = float64(ymax - ymin)
	return
}

// GetTargetBackPos 获取目标后背点
func (scene *Scene) GetTargetBackPos(target, original *Pawn) (*linemath.Vector2, bool) {
	targetPos := &linemath.Vector2{}
	targetRadius := original.Attr.CollisionRadius + target.Attr.CollisionRadius + DistancePawn(original, target)
	angle := target.GetAngle()

	//后方
	{
		targetAngle := angle + math.Pi
		targetPos.X = targetRadius * float32(math.Cos(float64(targetAngle)))
		targetPos.Y = targetRadius * float32(math.Sin(float64(targetAngle)))
		targetPos.AddS(original.GetPos())
		if !scene.checkCircleShapeOverlapWorldBoundary(float64(targetPos.X), float64(targetPos.Y), float64(original.Attr.CollisionRadius)) {
			return targetPos, true
		}
	}

	//右后方
	{
		targetAngle := angle + 0.75*math.Pi
		targetPos.X = targetRadius * float32(math.Cos(float64(targetAngle)))
		targetPos.Y = targetRadius * float32(math.Sin(float64(targetAngle)))
		targetPos.AddS(original.GetPos())
		if !scene.checkCircleShapeOverlapWorldBoundary(float64(targetPos.X), float64(targetPos.Y), float64(original.Attr.CollisionRadius)) {
			return targetPos, true
		}
	}

	//左后方
	{
		targetAngle := angle + 1.25*math.Pi
		targetPos.X = targetRadius * float32(math.Cos(float64(targetAngle)))
		targetPos.Y = targetRadius * float32(math.Sin(float64(targetAngle)))
		targetPos.AddS(original.GetPos())
		if !scene.checkCircleShapeOverlapWorldBoundary(float64(targetPos.X), float64(targetPos.Y), float64(original.Attr.CollisionRadius)) {
			return targetPos, true
		}
	}

	//右方
	{
		targetAngle := angle + 0.5*math.Pi
		targetPos.X = targetRadius * float32(math.Cos(float64(targetAngle)))
		targetPos.Y = targetRadius * float32(math.Sin(float64(targetAngle)))
		targetPos.AddS(original.GetPos())
		if !scene.checkCircleShapeOverlapWorldBoundary(float64(targetPos.X), float64(targetPos.Y), float64(original.Attr.CollisionRadius)) {
			return targetPos, true
		}
	}

	//左方
	{
		targetAngle := angle + 1.5*math.Pi
		targetPos.X = targetRadius * float32(math.Cos(float64(targetAngle)))
		targetPos.Y = targetRadius * float32(math.Sin(float64(targetAngle)))
		targetPos.AddS(original.GetPos())
		if !scene.checkCircleShapeOverlapWorldBoundary(float64(targetPos.X), float64(targetPos.Y), float64(original.Attr.CollisionRadius)) {
			return targetPos, true
		}
	}

	//左前方
	{
		targetAngle := angle + 1.75*math.Pi
		targetPos.X = targetRadius * float32(math.Cos(float64(targetAngle)))
		targetPos.Y = targetRadius * float32(math.Sin(float64(targetAngle)))
		targetPos.AddS(original.GetPos())
		if !scene.checkCircleShapeOverlapWorldBoundary(float64(targetPos.X), float64(targetPos.Y), float64(original.Attr.CollisionRadius)) {
			return targetPos, true
		}
	}

	//右前方
	{
		targetAngle := angle + 0.25*math.Pi
		targetPos.X = targetRadius * float32(math.Cos(float64(targetAngle)))
		targetPos.Y = targetRadius * float32(math.Sin(float64(targetAngle)))
		targetPos.AddS(original.GetPos())
		if !scene.checkCircleShapeOverlapWorldBoundary(float64(targetPos.X), float64(targetPos.Y), float64(original.Attr.CollisionRadius)) {
			return targetPos, true
		}
	}

	//前方
	{
		targetAngle := angle
		targetPos.X = targetRadius * float32(math.Cos(float64(targetAngle)))
		targetPos.Y = targetRadius * float32(math.Sin(float64(targetAngle)))
		targetPos.AddS(original.GetPos())
		if !scene.checkCircleShapeOverlapWorldBoundary(float64(targetPos.X), float64(targetPos.Y), float64(original.Attr.CollisionRadius)) {
			return targetPos, true
		}
	}

	return targetPos, false
}

func (scene *Scene) BuildNPCInfo(id uint32) (*PawnInfo, error) {
	logicCfg, ok := scene.GetMonsterExcelConfig().Logic_ConfigItems[id]
	if !ok {
		return nil, fmt.Errorf("not find MonsterConfig id:%v ", id)
	}

	isBoss := logicCfg.Difficulty != 0

	info := &PawnInfo{
		PawnInfo: &Proto.PawnInfo{
			Type:     Proto.PawnType_Npc,
			ConfigId: id,
			Npc: &Proto.FightNpcInfo{
				IsBoss: isBoss,
			},
			Camp:  Proto.Camp_Blue,
			Level: 1,
		},
		NormalAtkList:             logicCfg.NormalAttack,
		SuperSkillList:            logicCfg.SuperSkill,
		OverDriveNormalAttackList: logicCfg.OverDriveNormalAttack,
		OverDriveSuperSkillList:   logicCfg.OverDriveSuperSkill,
		BornBuffs:                 logicCfg.BornBuffs,
	}

	return info, nil
}

func (scene *Scene) GetPawnList() []*Pawn {
	return scene.pawnList
}
