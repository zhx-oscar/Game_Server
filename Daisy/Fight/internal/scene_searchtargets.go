package internal

import (
	"Cinder/Base/linemath"
	"Daisy/Fight/internal/conf"
	"Daisy/Proto"
	"sort"
)

// searchShapeTargets 获得区域内的目标列表
func (scene *Scene) searchShapeTargets(campBit Bits, pos linemath.Vector2, angle float32, shape *conf.AttackShape,
	checkCantBeEnemySelect, checkCantBeFriendlySelect bool) []*Pawn {
	if campBit == 0 || shape == nil {
		return nil
	}

	var targetList []*Pawn

	switch shape.Type {
	case Proto.AttackShapeType_Rect:
		targetList = scene.overlapRectangleShape(float64(pos.X), float64(pos.Y), float64(shape.Extend.X/2), float64(shape.Extend.Y/2), float64(angle))

	case Proto.AttackShapeType_Circle:
		targetList = scene.overlapCircleShape(float64(pos.X), float64(pos.Y), float64(shape.Radius))

	case Proto.AttackShapeType_Fan:
		targetList = scene.overlapSectorShape(pos, shape.Radius, shape.FanAngle/2, angle)
	}

	for i := len(targetList) - 1; i >= 0; i-- {
		if !targetList[i].IsAlive() {
			targetList = append(targetList[:i], targetList[i+1:]...)
			continue
		}

		if !campBit.Test(int32(Proto.Camp_Red)) {
			if targetList[i].GetCamp() == Proto.Camp_Red {
				targetList = append(targetList[:i], targetList[i+1:]...)
				continue
			}
		}

		if !campBit.Test(int32(Proto.Camp_Blue)) {
			if targetList[i].GetCamp() == Proto.Camp_Blue {
				targetList = append(targetList[:i], targetList[i+1:]...)
				continue
			}
		}

		if checkCantBeEnemySelect {
			if targetList[i].State.CantBeEnemySelect {
				targetList = append(targetList[:i], targetList[i+1:]...)
				continue
			}
		}

		if checkCantBeFriendlySelect {
			if targetList[i].State.CantBeFriendlySelect {
				targetList = append(targetList[:i], targetList[i+1:]...)
				continue
			}
		}
	}

	sort.Slice(targetList, func(i, j int) bool {
		return Distance(pos, targetList[i].GetPos()) < Distance(pos, targetList[j].GetPos())
	})

	return targetList
}
