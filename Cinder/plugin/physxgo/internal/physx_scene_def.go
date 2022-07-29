package internal

/*
#cgo CFLAGS: -I../physxcwrap/install/incl
#include "PhysXCWrap.h"
*/
import "C"
import "Cinder/Base/linemath"

// HitModeFlags 碰撞模式标记
type HitModeFlags uint32

const (
	HitMode_eHitStatic  HitModeFlags = 1 << iota // 检测静态目标
	HitMode_eHitDynamic                          // 检测动态目标
	HitMode_eAnyHit                              // 有碰撞目标就停止检测（可以提高性能，返回的结果不一定是最近的碰撞目标，eAnyHit与eAllHit都设置时，eAnyHit优先级高，都不设置表示返回最近的一个目标）
	HitMode_eAllHit                              // 返回所有碰撞目标，不按距离排序（eAnyHit与eAllHit都不设置表示返回最近的一个目标）
	HitMode_eHitTrap                             // 是否能碰撞区域触发器
	HitMode_eNone       HitModeFlags = 0         // 默认
)

// HitFilter 碰撞过滤器
type HitFilter struct {
	Word0, Word1, Word2, Word3 uint32
}

func (hitFilter *HitFilter) toC() C.HitFilter {
	var cHitFilter C.HitFilter
	cHitFilter.Word0 = C.uint32_t(hitFilter.Word0)
	cHitFilter.Word1 = C.uint32_t(hitFilter.Word1)
	cHitFilter.Word2 = C.uint32_t(hitFilter.Word2)
	cHitFilter.Word3 = C.uint32_t(hitFilter.Word3)
	return cHitFilter
}

func (hitFilter *HitFilter) fromC(cHitFilter *C.HitFilter) {
	if cHitFilter == nil {
		return
	}

	hitFilter.Word0 = uint32(cHitFilter.Word0)
	hitFilter.Word1 = uint32(cHitFilter.Word1)
	hitFilter.Word2 = uint32(cHitFilter.Word2)
	hitFilter.Word3 = uint32(cHitFilter.Word3)
}

// Hit 碰撞
type Hit struct {
	Target   IPxActor         // 碰撞Actor
	Position linemath.Vector3 // 碰撞点
	Normal   linemath.Vector3 // 碰撞点法线
	Distance float32          // 碰撞距离，小于等于原点与hit点间的距离
}

func (hit *Hit) fromC(cHit *C.Hit) {
	if cHit == nil {
		return
	}

	hit.Target = (*PxActor)(C.PxActorGetBindGoObj(cHit.Target))
	hit.Position = Vector3CToGo(cHit.Position)
	hit.Normal = Vector3CToGo(cHit.Normal)
	hit.Distance = float32(cHit.Distance)
}

// GeomType 几何体类型
type GeomType int32

const (
	GeomType_eSPHERE         GeomType = iota // 球体
	GeomType_ePLANE                          // 平面
	GeomType_eCAPSULE                        // 胶囊体
	GeomType_eBOX                            // 盒子
	GeomType_eCONVEXMESH                     // 凸面网格
	GeomType_eTRIANGLEMESH                   // 三角面网格
	GeomType_eHEIGHTFIELD                    // 高度空间
	GeomType_eGEOMETRY_COUNT                 //!< internal use only!
	GeomType_eINVALID        GeomType = -1   //!< internal use only!
)

// Geometry 几何体
type Geometry struct {
	Type        GeomType // 类型（只支持eSPHERE，eCAPSULE，eBOX）
	HalfExtents linemath.Vector3
	Radius      float32
	HalfHeight  float32
}

func (geom *Geometry) toC() C.Geometry {
	var cGeom C.Geometry
	cGeom.Type = C.GeomType(geom.Type)
	cGeom.HalfExtents = Vector3GoToC(geom.HalfExtents)
	cGeom.Radius = C.float(geom.Radius)
	cGeom.HalfHeight = C.float(geom.HalfHeight)
	return cGeom
}

// TransForm 位置与角度
type TransForm struct {
	P linemath.Vector3
	Q linemath.Quaternion
}

func (tf *TransForm) toC() C.TransForm {
	var cTf C.TransForm
	cTf.P = Vector3GoToC(tf.P)
	cTf.Q = QuatGoToC(tf.Q)
	return cTf
}

func (tf *TransForm) fromC(cTf *C.TransForm) {
	if cTf == nil {
		return
	}

	tf.P = Vector3CToGo(cTf.P)
	tf.Q = QuatCToGo(cTf.Q)
}

func Vector3GoToC(vec3 linemath.Vector3) C.Vector3 {
	var cVec3 C.Vector3
	cVec3.X = C.float(vec3.X)
	cVec3.Y = C.float(vec3.Y)
	cVec3.Z = C.float(vec3.Z)
	return cVec3
}

func Vector3CToGo(cVec3 C.Vector3) linemath.Vector3 {
	var vec3 linemath.Vector3
	vec3.X = float32(cVec3.X)
	vec3.Y = float32(cVec3.Y)
	vec3.Z = float32(cVec3.Z)
	return vec3
}

func QuatGoToC(quat linemath.Quaternion) C.Quat {
	var cQuat C.Quat
	cQuat.X = C.float(quat.X)
	cQuat.Y = C.float(quat.Y)
	cQuat.Z = C.float(quat.Z)
	cQuat.W = C.float(quat.W)
	return cQuat
}

func QuatCToGo(cQuat C.Quat) linemath.Quaternion {
	var quat linemath.Quaternion
	quat.X = float32(cQuat.X)
	quat.Y = float32(cQuat.Y)
	quat.Z = float32(cQuat.Z)
	quat.W = float32(cQuat.W)
	return quat
}

// IPxScene Scene接口
type IPxScene interface {
	// Release 销毁Scene
	Release()

	// Update Scene帧更新
	Update(elapsedTime float32)

	// AddBoxKinematic 创建Kinematic Box
	AddBoxKinematic(pose TransForm, halfExtents linemath.Vector3, actorModeFlags ActorModeFlags, hitFilter HitFilter,
		actorEvents *ActorEvents) (IPxActor, error)

	// AddBoxStatic 创建Static Box
	AddBoxStatic(pose TransForm, halfExtents linemath.Vector3, actorModeFlags ActorModeFlags, hitFilter HitFilter,
		actorEvents *ActorEvents) (IPxActor, error)

	// AddSphereKinematic 创建Kinematic Sphere
	AddSphereKinematic(pose TransForm, radius float32, actorModeFlags ActorModeFlags, hitFilter HitFilter,
		actorEvents *ActorEvents) (IPxActor, error)

	// AddSphereStatic 创建Static Sphere
	AddSphereStatic(pose TransForm, radius float32, actorModeFlags ActorModeFlags, hitFilter HitFilter,
		actorEvents *ActorEvents) (IPxActor, error)

	// AddCapsuleKinematic 创建Kinematic Capsule
	AddCapsuleKinematic(pose TransForm, radius float32, halfHeight float32, actorModeFlags ActorModeFlags, hitFilter HitFilter,
		actorEvents *ActorEvents) (IPxActor, error)

	// AddCapsuleStatic 创建Static Capsule
	AddCapsuleStatic(pose TransForm, radius float32, halfHeight float32, actorModeFlags ActorModeFlags, hitFilter HitFilter,
		actorEvents *ActorEvents) (IPxActor, error)

	// RaycastOne 射线检测单个
	RaycastOne(origin linemath.Vector3, unitDir linemath.Vector3, distance float32,
		hitModeFlags HitModeFlags, hitFilter HitFilter) (*Hit, error)

	// RaycastMany 射线检测多个
	RaycastMany(origin linemath.Vector3, unitDir linemath.Vector3, distance float32,
		hitModeFlags HitModeFlags, hitFilter HitFilter, maxHit int32) ([]Hit, error)

	// SweepOne 滑动检测单个
	SweepOne(geom Geometry, pose TransForm, unitDir linemath.Vector3, distance float32, inflation float32,
		hitModeFlags HitModeFlags, hitFilter HitFilter) (*Hit, error)

	// SweepMany 滑动检测多个
	SweepMany(geom Geometry, pose TransForm, unitDir linemath.Vector3, distance float32, inflation float32,
		hitModeFlags HitModeFlags, hitFilter HitFilter, maxHit int32) ([]Hit, error)

	// OverlapOne 重叠检测单个
	OverlapOne(geom Geometry, pose TransForm, hitModeFlags HitModeFlags, hitFilter HitFilter) (*Hit, error)

	// OverlapMany 重叠检测多个
	OverlapMany(geom Geometry, pose TransForm, hitModeFlags HitModeFlags, hitFilter HitFilter, maxHit int32) ([]Hit, error)
}
