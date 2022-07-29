package internal

/*
#cgo CFLAGS: -I../physxcwrap/install/incl
#include "PhysXCWrap.h"
*/
import "C"
import "Cinder/Base/linemath"

// ActorModeFlags Actor模式标记
type ActorModeFlags uint32

const (
	ActorMode_eNotBeQuery ActorModeFlags = 1 << iota // 不能被场景查询
	ActorMode_eNotBeTrap                             // 不能被区域触发器捕获
	ActorMode_eTrap                                  // 是区域触发器（Actor创建后不能被修改）
	ActorMode_eNone       ActorModeFlags = 0         // 默认
)

// StepRes 单步移动结果
type StepRes struct {
	BlockActor  IPxActor         // 阻挡移动的Actor（为nil表示到达目标点）
	BlockPose   TransForm        // 受到阻挡停止位置
	BlockNormal linemath.Vector3 // 受到阻挡碰撞点法线
}

func (sr *StepRes) fromC(cSr *C.StepRes) {
	if cSr == nil {
		return
	}

	sr.BlockActor = (*PxActor)(C.PxActorGetBindGoObj(cSr.BlockActor))
	sr.BlockPose.fromC(&cSr.BlockPose)
	sr.BlockNormal = Vector3CToGo(cSr.BlockNormal)
}

// ActorEvents Actor事件回调
type ActorEvents struct {
	OnActorEnterCallback func(self, actor IPxActor) // 当区域触发器有Actor进入
	OnActorLeaveCallback func(self, actor IPxActor) // 当区域触发器有Actor离开
	OnEnterTrapCallback  func(self, trap IPxActor)  // 当Actor进入区域触发器
	OnLeaveTrapCallback  func(self, trap IPxActor)  // 当Actor离开区域触发器
}

// IPxActor Actor接口
type IPxActor interface {
	// Release 销毁Scene
	Release()

	// SetUserData 设置用户数据
	SetUserData(dta interface{})

	// GetUserData 得到用户数据
	GetUserData() interface{}

	// SetActorModeFlags 设置模式标记
	SetActorModeFlags(actorModeFlags ActorModeFlags) error

	// GetActorModeFlags 查询模式标记
	GetActorModeFlags() ActorModeFlags

	// SetHitFilter 设置过滤器
	SetHitFilter(hitFilter HitFilter) error

	// GetHitFilter 查询过滤器
	GetHitFilter() HitFilter

	// GetActorEvents 获取事件回调
	GetActorEvents() *ActorEvents

	// GetPose 获取Actor坐标与朝向
	GetPose() TransForm

	// SetPose 设置Actor设置坐标与朝向
	SetPose(pose TransForm) error

	// SetPosition 设置Actor坐标
	SetPosition(pos linemath.Vector3) error

	// SetOrientation 设置Actor朝向
	SetOrientation(orient linemath.Quaternion) error

	// CheckStep 测试单步移动
	CheckStep(pose TransForm) (StepRes, error)

	// Step 单步移动
	Step(pose TransForm) (StepRes, error)
}
