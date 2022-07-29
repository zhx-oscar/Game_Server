package internal

/*
#cgo CFLAGS: -I../physxcwrap/install/incl
#include "PhysXCWrap.h"
#include "stdlib.h"
*/
import "C"
import (
	"Cinder/Base/linemath"
	"errors"
	"sync"
)

// actorMap 所有Actor
var actorMap sync.Map

// PxActor Actor
type PxActor struct {
	mPxActor     C.PxHandle
	mActorEvents ActorEvents
	mUserData    interface{}
}

// init 初始化Actor
func (actor *PxActor) init(cActor C.PxHandle, actorModeFlags ActorModeFlags, actorEvents *ActorEvents) bool {
	if cActor == nil {
		return false
	}

	actor.mPxActor = cActor

	if actorEvents != nil {
		actor.mActorEvents = *actorEvents
	}

	actor.mUserData = nil

	return true
}

// Release 销毁Scene
func (actor *PxActor) Release() {
	if actor.mPxActor == nil {
		return
	}

	inTrapNum := C.PxActorCountInTraps(actor.mPxActor)

	for i := inTrapNum - 1; i >= 0; i-- {
		inTrapData := C.PxActorGetInTrapData(actor.mPxActor, i)

		trap := (*PxActor)(C.PxActorGetBindGoObj(inTrapData.TrapActor))
		if trap == nil {
			continue
		}

		switch inTrapData.SelfStat {
		case C.eStay:
			fallthrough
		case C.eLeave:
			if actor.mActorEvents.OnLeaveTrapCallback != nil {
				actor.mActorEvents.OnLeaveTrapCallback(actor, trap)
			}

			if trap.mActorEvents.OnActorLeaveCallback != nil {
				trap.mActorEvents.OnActorLeaveCallback(trap, actor)
			}

			C.PxActorDeleteInTrapData(actor.mPxActor, i)
		}
	}

	C.ReleasePxActor(actor.mPxActor)
	actor.mPxActor = nil
	actor.mUserData = nil
	actorMap.Delete(actor)
}

// SetUserData 设置用户数据
func (actor *PxActor) SetUserData(data interface{}) {
	actor.mUserData = data
}

// GetUserData 获取用户数据
func (actor *PxActor) GetUserData() interface{} {
	return actor.mUserData
}

// SetActorModeFlags 设置模式标记
func (actor *PxActor) SetActorModeFlags(actorModeFlags ActorModeFlags) error {
	if actor.mPxActor == nil {
		return errors.New("nil PxActor")
	}

	if C.PxActorSetActorModeFlags(actor.mPxActor, C.ActorModeFlags(actorModeFlags)) != C._Bool(true) {
		return errors.New("PxActor SetActorModeFlags failed")
	}

	actor.emitEnterLeaveTrapEvent()

	return nil
}

// GetActorModeFlags 查询模式标记
func (actor *PxActor) GetActorModeFlags() ActorModeFlags {
	if actor.mPxActor == nil {
		return ActorMode_eNone
	}

	cActorModeFlags := C.PxActorGetActorModeFlags(actor.mPxActor)
	return ActorModeFlags(cActorModeFlags)
}

// SetHitFilter 设置过滤器
func (actor *PxActor) SetHitFilter(hitFilter HitFilter) error {
	if actor.mPxActor == nil {
		return errors.New("nil PxActor")
	}

	if C.PxActorSetHitFilter(actor.mPxActor, hitFilter.toC()) != C._Bool(true) {
		return errors.New("PxActor SetHitFilter failed")
	}

	return nil
}

// GetHitFilter 查询过滤器
func (actor *PxActor) GetHitFilter() HitFilter {
	if actor.mPxActor == nil {
		return HitFilter{}
	}

	cHitFilter := C.PxActorGetHitFilter(actor.mPxActor)

	var hitFilter HitFilter
	hitFilter.fromC(&cHitFilter)

	return hitFilter
}

// GetActorEvents 获取事件回调
func (actor *PxActor) GetActorEvents() *ActorEvents {
	return &actor.mActorEvents
}

// GetPose 获取Actor坐标与角度
func (actor *PxActor) GetPose() TransForm {
	cPose := C.PxActorGetPose(actor.mPxActor)
	pose := TransForm{}
	pose.fromC(&cPose)
	return pose
}

// SetPose 设置Actor设置坐标与朝向
func (actor *PxActor) SetPose(pose TransForm) error {
	if actor.mPxActor == nil {
		return errors.New("nil PxActor")
	}

	cPose := pose.toC()
	if C.PxActorSetPose(actor.mPxActor, cPose) != C._Bool(true) {
		return errors.New("PxActor SetPose failed")
	}

	actor.emitEnterLeaveTrapEvent()

	return nil
}

// SetPosition 设置Actor坐标
func (actor *PxActor) SetPosition(pos linemath.Vector3) error {
	if actor.mPxActor == nil {
		return errors.New("nil PxActor")
	}

	cPos := Vector3GoToC(pos)
	if C.PxActorSetPosition(actor.mPxActor, cPos) != C._Bool(true) {
		return errors.New("PxActor SetPosition failed")
	}

	actor.emitEnterLeaveTrapEvent()

	return nil
}

// SetOrientation 设置Actor朝向
func (actor *PxActor) SetOrientation(orient linemath.Quaternion) error {
	if actor.mPxActor == nil {
		return errors.New("nil PxActor")
	}

	cOrient := QuatGoToC(orient)
	if C.PxActorSetOrientation(actor.mPxActor, cOrient) != C._Bool(true) {
		return errors.New("PxActor SetOrientation failed")
	}

	actor.emitEnterLeaveTrapEvent()

	return nil
}

// CheckStep 测试单步移动
func (actor *PxActor) CheckStep(pose TransForm) (StepRes, error) {
	if actor.mPxActor == nil {
		return StepRes{}, errors.New("nil PxActor")
	}

	cPose := pose.toC()
	cStepRes := C.PxActorCheckStep(actor.mPxActor, cPose)
	if cStepRes.Ok != C._Bool(true) {
		return StepRes{}, errors.New("PxActor CheckStep failed")
	}

	var stepRes StepRes
	stepRes.fromC(&cStepRes)

	return stepRes, nil
}

// Step 单步移动
func (actor *PxActor) Step(pose TransForm) (StepRes, error) {
	if actor.mPxActor == nil {
		return StepRes{}, errors.New("nil PxActor")
	}

	cPose := pose.toC()
	cStepRes := C.PxActorStep(actor.mPxActor, cPose)
	if cStepRes.Ok != C._Bool(true) {
		return StepRes{}, errors.New("PxActor Step failed")
	}

	var stepRes StepRes
	stepRes.fromC(&cStepRes)

	actor.emitEnterLeaveTrapEvent()

	return stepRes, nil
}

// emitEnterLeaveTrapEvent 发送离开区域触发器事件
func (actor *PxActor) emitEnterLeaveTrapEvent() {
	inTrapNum := C.PxActorCountInTraps(actor.mPxActor)

	for i := inTrapNum - 1; i >= 0; i-- {
		inTrapData := C.PxActorGetInTrapData(actor.mPxActor, i)

		trap := (*PxActor)(C.PxActorGetBindGoObj(inTrapData.TrapActor))
		if trap == nil {
			continue
		}

		switch inTrapData.SelfStat {
		case C.eEnter:
			if actor.mActorEvents.OnEnterTrapCallback != nil {
				actor.mActorEvents.OnEnterTrapCallback(actor, trap)
			}

			if trap.mActorEvents.OnActorEnterCallback != nil {
				trap.mActorEvents.OnActorEnterCallback(trap, actor)
			}

			C.PxActorSetInTrapStat(actor.mPxActor, i, C.eStay)

		case C.eLeave:
			if actor.mActorEvents.OnLeaveTrapCallback != nil {
				actor.mActorEvents.OnLeaveTrapCallback(actor, trap)
			}

			if trap.mActorEvents.OnActorLeaveCallback != nil {
				trap.mActorEvents.OnActorLeaveCallback(trap, actor)
			}

			C.PxActorDeleteInTrapData(actor.mPxActor, i)
		}
	}
}
