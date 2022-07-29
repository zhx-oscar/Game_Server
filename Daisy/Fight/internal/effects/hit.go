package effects

import (
	"Cinder/Base/linemath"
	"Daisy/Data"
	. "Daisy/Fight/internal"
	"Daisy/Fight/internal/conf"
	"Daisy/Proto"
	"fmt"
	"math"
	"math/rand"
)

// _HitMutex 命中互斥类型
type _HitMutex uint32

const (
	_HitMutex_Replace          _HitMutex = iota // 替换
	_HitMutex_Discard                           // 丢弃
	_HitMutex_Union                             // 共存
	_HitMutex_ExtendTime                        // 延长时间
	_HitMutex_ReplaceToReFloat                  // 替换成再次击飞
)

// _HitMutexTab 命中类型互斥表定义
type _HitMutexTab [][]_HitMutex

// hitMutexTab 命中类型互斥表（行：新状态，列：旧状态）
var hitMutexTab = _HitMutexTab{
	//			 	空状态				击中			     击破			 击倒				     击飞				      再次击飞			   击晕				 格挡				格挡崩溃
	/*空状态*/ {_HitMutex_Discard, _HitMutex_Discard, _HitMutex_Discard, _HitMutex_Discard, _HitMutex_ReplaceToReFloat, _HitMutex_ReplaceToReFloat, _HitMutex_Discard, _HitMutex_Discard, _HitMutex_Discard},
	/*击中*/ {_HitMutex_Replace, _HitMutex_ExtendTime, _HitMutex_Discard, _HitMutex_Discard, _HitMutex_ReplaceToReFloat, _HitMutex_ReplaceToReFloat, _HitMutex_Union, _HitMutex_Discard, _HitMutex_Replace},
	/*击破*/ {_HitMutex_Replace, _HitMutex_Replace, _HitMutex_ExtendTime, _HitMutex_Discard, _HitMutex_ReplaceToReFloat, _HitMutex_ReplaceToReFloat, _HitMutex_Union, _HitMutex_Discard, _HitMutex_Replace},
	/*击倒*/ {_HitMutex_Replace, _HitMutex_Replace, _HitMutex_Replace, _HitMutex_ExtendTime, _HitMutex_ReplaceToReFloat, _HitMutex_ReplaceToReFloat, _HitMutex_Union, _HitMutex_Discard, _HitMutex_Replace},
	/*击飞*/ {_HitMutex_Replace, _HitMutex_Replace, _HitMutex_Replace, _HitMutex_Replace, _HitMutex_ReplaceToReFloat, _HitMutex_ReplaceToReFloat, _HitMutex_Union, _HitMutex_Discard, _HitMutex_Replace},
	/*再次击飞*/ {_HitMutex_Discard, _HitMutex_Discard, _HitMutex_Discard, _HitMutex_Discard, _HitMutex_Discard, _HitMutex_Discard, _HitMutex_Discard, _HitMutex_Discard, _HitMutex_Discard},
	/*击晕*/ {_HitMutex_Replace, _HitMutex_Union, _HitMutex_Union, _HitMutex_Union, _HitMutex_Union, _HitMutex_Union, _HitMutex_Union, _HitMutex_Discard, _HitMutex_Replace},
	/*格挡*/ {_HitMutex_Replace, _HitMutex_Discard, _HitMutex_Discard, _HitMutex_Discard, _HitMutex_Discard, _HitMutex_Discard, _HitMutex_Discard, _HitMutex_ExtendTime, _HitMutex_Discard},
	/*格挡崩溃*/ {_HitMutex_Discard, _HitMutex_Discard, _HitMutex_Discard, _HitMutex_Discard, _HitMutex_Discard, _HitMutex_Discard, _HitMutex_Discard, _HitMutex_Discard, _HitMutex_Discard},
}

// Overlap 状态叠加
func (tab *_HitMutexTab) Overlap(new, old Proto.HitType_Enum) _HitMutex {
	return (*tab)[new][old]
}

// Hit 击打模板
type Hit struct {
}

// Hit 击打
func (hit *Hit) Hit(attack *Attack, target *Pawn, damageBit Bits) {
	caster := attack.Caster
	if attack.Src() != Proto.AttackSrc_Skill {
		return
	}

	// 伤害类型为闪避，执行闪避逻辑
	if damageBit.Test(int32(Proto.DamageType_Dodge)) {
		hit.Dodge(caster, target)
		return
	}

	if damageBit.Test(int32(Proto.DamageType_Miss)) || damageBit.Test(int32(Proto.DamageType_ExemptionDamage)) {
		return
	}

	hitState := Proto.HitType_None

	if int(attack.HitTimes) >= len(attack.Config.Hits) {
		return
	}

	// 处理顿帧
	attack.Skill.ExtendAttackTime(attack.Config.Hits[attack.HitTimes].Pause)
	attack.ExtendHitTime(attack.Config.Hits[attack.HitTimes].Pause)

	if !target.IsAlive() {
		return
	}

	// 处理受击状态
	if damageBit.Test(int32(Proto.DamageType_Block)) {
		hitState = hit.changeBeHitState(target, Proto.HitType_Block)
	} else if !target.State.CantBeHitControl {
		hitState = hit.changeBeHitState(target, attack.Config.Hits[attack.HitTimes].HitType)
	}

	// 发送击打事件
	caster.Events.EmitHitTarget(attack, target, damageBit, hitState)

	// 发送受击事件
	target.Events.EmitBeHit(attack, damageBit, hitState)

	if hitState == Proto.HitType_Block && target.State.BeHitStat.Test(int32(Proto.HitType_BlockBreak)) {
		hitState = Proto.HitType_BlockBreak
	}

	if hitState == Proto.HitType_None {
		return
	}

	target.BreakCurSkill(caster, Proto.SkillBreakReason_Normal)
	if target.IsMoving() {
		target.Stop()
	}

	caster.Scene.PushAction(&Proto.BeHit{
		TargetId: target.UID,
		AttackId: attack.UID,
		HitType:  hitState,
	})

	caster.Scene.PushDebugInfo(func() string {
		return fmt.Sprintf("${PawnID:%d}%s${PawnID:%d}",
			caster.UID,
			func() string {
				switch hitState {
				case Proto.HitType_Hit:
					return "击中"
				case Proto.HitType_Broken:
					return "击破"
				case Proto.HitType_Down:
					return "击倒"
				case Proto.HitType_Float:
					return "击飞"
				case Proto.HitType_ReFloat:
					return "重复击飞"
				case Proto.HitType_Stun:
					return "击晕"
				case Proto.HitType_Block:
					return "格挡"
				case Proto.HitType_BlockBreak:
					return "格挡崩溃"
				}
				return ""
			}(),
			target.UID,
		)
	})

	if hitState == Proto.HitType_Block || hitState == Proto.HitType_BlockBreak {
		return
	}

	// 处理击退
	hit.hitBack(caster, target, attack.Config.Hits[attack.HitTimes].HitBackDuring, attack.Config.Hits[attack.HitTimes].HitBackDis)
}

// changeBeHitState 修改受击状态
func (hit *Hit) changeBeHitState(pawn *Pawn, newHitState Proto.HitType_Enum) Proto.HitType_Enum {
	// 检查状态
	if pawn.State.CantBeHitControl {
		return Proto.HitType_None
	}

	// 查询质量配置
	massConfig, ok := Data.GetMassConfig().Mass_ConfigItems[pawn.Attr.Mass]
	if !ok {
		return Proto.HitType_None
	}

	switch newHitState {
	case Proto.HitType_Hit:
		if massConfig.ImmuneLightHit {
			return Proto.HitType_None
		}
	case Proto.HitType_Broken:
		if massConfig.ImmuneBreak {
			return Proto.HitType_None
		}
	case Proto.HitType_Down:
		if massConfig.ImmuneDown {
			return Proto.HitType_None
		}
	case Proto.HitType_Float:
		fallthrough
	case Proto.HitType_ReFloat:
		if massConfig.ImmuneAir {
			return Proto.HitType_None
		}
	case Proto.HitType_Stun:
		if massConfig.ImmuneStun {
			return Proto.HitType_None
		}
	}

	for oldHitState := Proto.HitType_Hit; oldHitState <= Proto.HitType_BlockBreak; oldHitState++ {
		if !pawn.State.BeHitStat.Test(int32(oldHitState)) {
			continue
		}

		hitMutex := hitMutexTab.Overlap(newHitState, oldHitState)

		switch hitMutex {
		case _HitMutex_Replace:
			pawn.State.ChangeBeHitStatBit(oldHitState, false)
			pawn.State.ChangeBeHitStatBit(newHitState, true)

			hitTime := pawn.Info.BeHitConf.GetBeHitTime(newHitState)
			pawn.BeHit.BeHitStartTime[newHitState] = pawn.Scene.NowTime
			pawn.BeHit.BeHitEndTime[newHitState] = pawn.Scene.NowTime + hitTime

			return newHitState

		case _HitMutex_Discard:
			return Proto.HitType_None

		case _HitMutex_Union:
			continue

		case _HitMutex_ExtendTime:
			hitTime := pawn.Info.BeHitConf.GetBeHitTime(oldHitState)
			pawn.BeHit.BeHitEndTime[oldHitState] = pawn.Scene.NowTime + hitTime

			return oldHitState
		case _HitMutex_ReplaceToReFloat:
			var ok bool
			var riseTime, reRiseTime, deltaTime uint32
			var gravity float64
			riseTime, ok = pawn.Scene.GetConstConUint32Value(conf.ConstExcel_HitFloatRiseTime)
			if !ok {
				return Proto.HitType_None
			}

			reRiseTime, ok = pawn.Scene.GetConstConUint32Value(conf.ConstExceL_HitReFloatRiseTime)
			if !ok {
				return Proto.HitType_None
			}

			deltaTime, ok = pawn.Scene.GetConstConUint32Value(conf.ConstExcel_HitReFloatDeltaTime)
			if !ok {
				return Proto.HitType_None
			}

			gravity, ok = pawn.Scene.GetConstConFloat64Value(conf.ConstExcel_GravityAcceleration)
			if !ok {
				return Proto.HitType_None
			}

			riseSpeed := gravity * float64(riseTime) / float64(1000)
			reRiseSpeed := gravity * float64(reRiseTime) / float64(1000)

			var hitTime uint32
			switch oldHitState {
			case Proto.HitType_Float:
				if deltaTime > riseTime {
					return Proto.HitType_None
				}

				if pawn.Scene.NowTime > pawn.BeHit.BeHitEndTime[oldHitState]-pawn.Info.BeHitConf.Float.AfterTime ||
					pawn.Scene.NowTime < pawn.BeHit.BeHitStartTime[oldHitState]+riseTime+deltaTime {
					return Proto.HitType_None
				}

				oldStateTime := pawn.Scene.NowTime - pawn.BeHit.BeHitStartTime[oldHitState]
				curHeight := Max(riseSpeed*float64(oldStateTime)/float64(1000)-0.5*gravity*math.Pow(float64(oldStateTime)/float64(1000), 2), 0)
				maxHeight := curHeight + Max(reRiseSpeed*float64(reRiseTime)/float64(1000)-0.5*gravity*math.Pow(float64(reRiseTime)/float64(1000), 2), 0)
				pawn.BeHit.ReFloatStartHeight = curHeight
				hitTime = uint32(math.Pow(maxHeight*2/gravity, 0.5)*float64(1000)) + reRiseTime + pawn.Info.BeHitConf.Float.AfterTime
			case Proto.HitType_ReFloat:
				if pawn.BeHit.BeHitStartTime[oldHitState]+reRiseTime+deltaTime > pawn.BeHit.BeHitEndTime[oldHitState]-pawn.Info.BeHitConf.Float.AfterTime {
					return Proto.HitType_None
				}

				if pawn.Scene.NowTime > pawn.BeHit.BeHitEndTime[oldHitState]-pawn.Info.BeHitConf.Float.AfterTime ||
					pawn.Scene.NowTime < pawn.BeHit.BeHitStartTime[oldHitState]+reRiseTime+deltaTime {
					return Proto.HitType_None
				}

				oldStateTime := pawn.Scene.NowTime - pawn.BeHit.BeHitStartTime[oldHitState]
				curHeight := Max(pawn.BeHit.ReFloatStartHeight+reRiseSpeed*float64(oldStateTime)/float64(1000)-0.5*gravity*math.Pow(float64(oldStateTime)/float64(1000), 2), 0)
				maxHeight := curHeight + Max(reRiseSpeed*float64(reRiseTime)/float64(1000)-0.5*gravity*math.Pow(float64(reRiseTime)/float64(1000), 2), 0)
				pawn.BeHit.ReFloatStartHeight = curHeight
				hitTime = uint32(math.Pow(maxHeight*2/gravity, 0.5)*float64(1000)) + reRiseTime + pawn.Info.BeHitConf.Float.AfterTime
			default:
				return Proto.HitType_None
			}

			newHitState = Proto.HitType_ReFloat
			pawn.State.ChangeBeHitStatBit(oldHitState, false)
			pawn.State.ChangeBeHitStatBit(newHitState, true)

			pawn.BeHit.BeHitStartTime[newHitState] = pawn.Scene.NowTime
			pawn.BeHit.BeHitEndTime[newHitState] = pawn.Scene.NowTime + hitTime

			return newHitState
		}
	}

	if newHitState == Proto.HitType_None {
		return Proto.HitType_None
	}

	pawn.State.ChangeBeHitStatBit(newHitState, true)

	hitTime := pawn.Info.BeHitConf.GetBeHitTime(newHitState)
	pawn.BeHit.BeHitStartTime[newHitState] = pawn.Scene.NowTime
	pawn.BeHit.BeHitEndTime[newHitState] = pawn.Scene.NowTime + hitTime

	return newHitState
}

// Dodge 闪避
func (hit *Hit) Dodge(caster, target *Pawn) {
	if target.State.Dodging || target.Info.DodgeDist <= 0 || target.Info.DodgeTime <= 0 {
		return
	}

	if target.IsSkillRunning() {
		target.BreakCurSkill(caster, Proto.SkillBreakReason_Normal)
	}

	hitAngle := float32(CalcAngle(caster.GetPos(), target.GetPos()))
	dodgeAngleRange := float32(conf.DodgeAngleRange) / float32(180) * math.Pi
	dodgeAngle := (rand.Float32()-0.5)*dodgeAngleRange + hitAngle
	dodgeTarget := target.GetPos().Add(linemath.Vector2{X: target.Info.DodgeDist * float32(math.Cos(float64(dodgeAngle))), Y: target.Info.DodgeDist * float32(math.Sin(float64(dodgeAngle)))})
	target.BeHit.DodgeEndTime = target.Scene.NowTime + target.Info.DodgeTime
	target.State.ChangeStat(Stat_Invincible, true)
	target.AIPause(true)
	target.State.ChangeStat(Stat_Dodging, true)
	target.MoveToAndChangeAngle(Proto.MoveMode_Dodge, dodgeTarget, target.GetAngleVelocity(dodgeAngle, target.Info.DodgeDist/float32(target.Info.DodgeTime)*float32(1000)), dodgeAngle, true)
}

// hitBack 击退
func (hit *Hit) hitBack(caster, target *Pawn, time uint32, distance float32) bool {
	if time <= 0 || FloatEqual(float64(distance), 0) {
		return false
	}

	angle := CalcAngle(caster.GetPos(), target.GetPos())
	targetPos := target.GetPos().Add(linemath.Vector2{X: distance * float32(math.Cos(angle)), Y: distance * float32(math.Sin(angle))})
	target.SetHitAngle(caster.GetAngle())
	initialVelocity := target.GetAngleVelocity(float32(angle), 2*distance/float32(time)*float32(1000))
	acceleration := initialVelocity.Mul(-1000 / float32(time))
	target.UniformlyVariableMoveTo(Proto.MoveMode_HitBack, targetPos, initialVelocity, acceleration, float32(time)/1000, target.GetAngle(), true)
	return true
}
