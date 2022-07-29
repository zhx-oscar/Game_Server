package physxgo

import (
	"Cinder/plugin/physxgo/internal"
)

// HitModeFlags 碰撞模式标记
type HitModeFlags = internal.HitModeFlags

const (
	HitMode_eNone       = internal.HitMode_eNone       // 默认
	HitMode_eHitStatic  = internal.HitMode_eHitStatic  // 检测静态目标
	HitMode_eHitDynamic = internal.HitMode_eHitDynamic // 检测动态目标
	HitMode_eAnyHit     = internal.HitMode_eAnyHit     // 有碰撞目标就停止检测（可以提高性能，返回的结果不一定是最近的碰撞目标，eAnyHit与eAllHit都设置时，eAnyHit优先级高，都不设置表示返回最近的一个目标）
	HitMode_eAllHit     = internal.HitMode_eAllHit     // 返回所有碰撞目标，不按距离排序（eAnyHit与eAllHit都不设置表示返回最近的一个目标）
	HitMode_eHitTrap    = internal.HitMode_eHitTrap    // 是否能碰撞区域触发器
)

// HitFilter 碰撞过滤器
type HitFilter = internal.HitFilter

// NoHitFilter 不使用碰撞过滤器
var NoHitFilter = HitFilter{}

// Hit 碰撞
type Hit = internal.Hit

// GeomType 几何体类型
type GeomType = internal.GeomType

const (
	GeomType_eSPHERE       = internal.GeomType_eSPHERE       // 球体
	GeomType_ePLANE        = internal.GeomType_ePLANE        // 平面
	GeomType_eCAPSULE      = internal.GeomType_eCAPSULE      // 胶囊体
	GeomType_eBOX          = internal.GeomType_eBOX          // 盒子
	GeomType_eCONVEXMESH   = internal.GeomType_eCONVEXMESH   // 凸面网格
	GeomType_eTRIANGLEMESH = internal.GeomType_eTRIANGLEMESH // 三角面网格
	GeomType_eHEIGHTFIELD  = internal.GeomType_eHEIGHTFIELD  // 高度空间
)

// Geometry 几何体
type Geometry = internal.Geometry

// ActorModeFlags Actor模式标记
type ActorModeFlags = internal.ActorModeFlags

const (
	ActorMode_eNone       = internal.ActorMode_eNone       // 默认
	ActorMode_eNotBeQuery = internal.ActorMode_eNotBeQuery // 不能被场景查询
	ActorMode_eNotBeTrap  = internal.ActorMode_eNotBeTrap  // 不能被区域触发器捕获
	ActorMode_eTrap       = internal.ActorMode_eTrap       // 是区域触发器（Actor创建后不能被设置）
)

// TransForm 位置与角度
type TransForm = internal.TransForm

// StepRes 单步移动结果
type StepRes = internal.StepRes

// ActorEvents Actor事件回调
type ActorEvents = internal.ActorEvents

// IPxActor Actor接口
type IPxActor = internal.IPxActor

// IPxScene Scene接口
type IPxScene = internal.IPxScene

// InitPxSdk 初始化PxSdk
func InitPxSdk(usePvd bool, pvdHost string, port int32) error {
	return internal.InitPxSdk(usePvd, pvdHost, port)
}

// ShutPxSdk 销毁PxSdk
func ShutPxSdk() {
	internal.ShutPxSdk()
}

// CreatePxScene 创建PxScene
func CreatePxScene(multiThread bool) (IPxScene, error) {
	return internal.CreatePxScene(multiThread)
}
