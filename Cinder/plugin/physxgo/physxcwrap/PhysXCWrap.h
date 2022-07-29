#pragma once

#include <stddef.h>
#include <stdint.h>
#include <stdbool.h>

#ifdef _MSC_VER
#if PHYSXCWRAP_EXPORTS
#define PHYSXCWRAP_API __declspec(dllexport)
#else
#define PHYSXCWRAP_API __declspec(dllimport)
#endif
#else 
#define PHYSXCWRAP_API
#endif

#if __cplusplus
extern "C"
{
#endif

// PX句柄
typedef void* PxHandle;

// vector3
typedef struct Vector3 {
	float X;
	float Y;
	float Z;
} Vector3;

// 四元数  
typedef struct Quat {
	float X;
	float Y;
	float Z;
	float W;
} Quat;

// 位置和角度
typedef struct TransForm {
	Vector3 P;
	Quat Q;
} TransForm;

PHYSXCWRAP_API TransForm ZeroTransForm();

// 碰撞检测模式
typedef enum HitMode {
	eHitStatic  = 1 << 0,	// 只检测静态目标
	eHitDynamic = 1 << 1,   // 只检测动态目标
	eAnyHit     = 1 << 2,   // 有碰撞目标就停止检测（可以提高性能，返回的结果不一定是最近的碰撞目标，eAnyHit与eAllHit都设置时，eAnyHit优先级高，都不设置表示返回最近的一个目标）
	eAllHit     = 1 << 3,   // 返回所有碰撞目标（eAnyHit与eAllHit都不设置表示返回最近的一个目标）
	eHitTrap    = 1 << 4,   // 是否能碰撞区域触发器
} HitMode;

// 碰撞检测模式标记
typedef uint32_t HitModeFlags;

// 碰撞过滤器
typedef struct HitFilter {
	uint32_t Word0;
	uint32_t Word1;
	uint32_t Word2;
	uint32_t Word3;
} HitFilter;

// 碰撞
typedef struct Hit {
	PxHandle Target;  // 碰撞Actor
	Vector3 Position; // 碰撞点
	Vector3 Normal;   // 碰撞点法线
	float Distance;   // 碰撞距离，小于等于原点与hit点间的距离	
} Hit;

// 几何体类型
typedef enum GeomType {
	eSPHERE,			// 球体
	ePLANE,             // 平面
	eCAPSULE,           // 胶囊体
	eBOX,               // 盒子
	eCONVEXMESH,        // 凸面网格
	eTRIANGLEMESH,      // 三角面网格
	eHEIGHTFIELD,       // 高度空间
	eGEOMETRY_COUNT,	//!< internal use only!
	eINVALID = -1		//!< internal use only!
} GeomType;

// 几何体
typedef struct Geometry {
	GeomType Type; // 类型（只支持eSPHERE，eCAPSULE，eBOX）
	Vector3 HalfExtents;
	float Radius;
	float HalfHeight;
} Geometry;

// Actor模式
typedef enum ActorMode {
	eNotBeQuery = 1 << 0,  // 不能被场景查询		 
	eNotBeTrap  = 1 << 1,  // 不能被区域触发器捕获
	eTrap       = 1 << 2,  // 是区域触发器（Actor创建后不能被设置）
} ActorMode;

// Actor模式标记
typedef uint32_t ActorModeFlags;

// 单步移动结果
typedef struct StepRes {	
	PxHandle BlockActor;   // 阻挡移动的Actor
	TransForm BlockPose;   // 受到阻挡停止位置
	Vector3 BlockNormal;   // 受到阻挡碰撞点法线
	bool Ok;               // 执行成功
} StepRes;

// 在范围触发器中状态
typedef enum InTrapStat
{
	eNone,   // 无
	eEnter,  // 进入
	eStay,   // 停留
	eLeave   // 离开
} InTrapStat;

// 在范围触发器中数据 
typedef struct InTrapData
{
	PxHandle TrapActor;
	InTrapStat SelfStat;
} InTrapData;

// 创建Sdk
PHYSXCWRAP_API PxHandle CreatePxSdk();

// 创建Sdk并连接pvd工具
PHYSXCWRAP_API PxHandle CreatePxSdkConnectPvd(const char* pvdHost, int32_t pvdPort, uint32_t timeoutInMs);

// 销毁Sdk
PHYSXCWRAP_API bool ReleasePxSdk(PxHandle sdk);

// 创建Scene
PHYSXCWRAP_API PxHandle PxSdkCreatePxScene(PxHandle sdk, bool multiThread, void* bindGoObj);

// 销毁Scene
PHYSXCWRAP_API bool ReleasePxScene(PxHandle scene);

// Scene帧更新
PHYSXCWRAP_API void PxSceneUpdate(PxHandle scene, float elapsedTime);

// 创建Kinematic Box
PHYSXCWRAP_API PxHandle PxSceneCreateBoxKinematic(PxHandle scene, TransForm pose, Vector3 halfExtents, ActorModeFlags actorModeFlags, 
	HitFilter hitFilter, void* bindGoObj);

// 创建Static Box
PHYSXCWRAP_API PxHandle PxSceneCreateBoxStatic(PxHandle scene, TransForm pose, Vector3 halfExtents, ActorModeFlags actorModeFlags, 
	HitFilter hitFilter, void* bindGoObj);

// 创建Kinematic Sphere
PHYSXCWRAP_API PxHandle PxSceneCreateSphereKinematic(PxHandle scene, TransForm pose, float radius, ActorModeFlags actorModeFlags, 
	HitFilter hitFilter, void* bindGoObj);

// 创建Static Sphere
PHYSXCWRAP_API PxHandle PxSceneCreateSphereStatic(PxHandle scene, TransForm pose, float radius, ActorModeFlags actorModeFlags, 
	HitFilter hitFilter, void* bindGoObj);

// 创建Kinematic Capsule
PHYSXCWRAP_API PxHandle PxSceneCreateCapsuleKinematic(PxHandle scene, TransForm pose, float radius, float halfHeight, ActorModeFlags actorModeFlags, 
	HitFilter hitFilter, void* bindGoObj);

// 创建Static Capsule
PHYSXCWRAP_API PxHandle PxSceneCreateCapsuleStatic(PxHandle scene, TransForm pose, float radius, float halfHeight, ActorModeFlags actorModeFlags, 
	HitFilter hitFilter, void* bindGoObj);

// Scene添加Actor
PHYSXCWRAP_API bool PxSceneAddActor(PxHandle scene, PxHandle actor);

// 射线检测（返回值-1：失败，[0,N]：碰撞数量）
PHYSXCWRAP_API int32_t PxSceneRaycast(PxHandle scene, Vector3 origin, Vector3 unitDir, float distance, 
	HitModeFlags hitModeFlags, HitFilter hitFilter, Hit* hitBuf, size_t hitBufLen);

// 滑动检测（返回值-1：失败，[0,N]：碰撞数量）
PHYSXCWRAP_API int32_t PxSceneSweep(PxHandle scene, Geometry geom, TransForm pose, Vector3 unitDir, float distance, float inflation,
	HitModeFlags hitModeFlags, HitFilter hitFilter, Hit* hitBuf, size_t hitBufLen);

// 重叠检测（返回值-1：失败，[0,N]：碰撞数量）
PHYSXCWRAP_API int32_t PxSceneOverlap(PxHandle scene, Geometry geom, TransForm pose,
	HitModeFlags hitModeFlags, HitFilter hitFilter, Hit* hitBuf, size_t hitBufLen);

// 销毁Actor
PHYSXCWRAP_API bool ReleasePxActor(PxHandle actor);

// 获取Actor绑定的Go对象
PHYSXCWRAP_API void* PxActorGetBindGoObj(PxHandle actor);

// 设置Actor模式标记
PHYSXCWRAP_API bool PxActorSetActorModeFlags(PxHandle actor, ActorModeFlags actorModeFlags);

// 获取Actor模式标记
PHYSXCWRAP_API ActorModeFlags PxActorGetActorModeFlags(PxHandle actor);

// 设置Actor过滤器
PHYSXCWRAP_API bool PxActorSetHitFilter(PxHandle actor, HitFilter hitFilter);

// 获取Actor过滤器
PHYSXCWRAP_API HitFilter PxActorGetHitFilter(PxHandle actor);

// 获取Actor坐标与朝向
PHYSXCWRAP_API TransForm PxActorGetPose(PxHandle actor);

// 设置Actor设置坐标与朝向
PHYSXCWRAP_API bool PxActorSetPose(PxHandle actor, TransForm pose);

// 设置Actor坐标
PHYSXCWRAP_API bool PxActorSetPosition(PxHandle actor, Vector3 pos);

// 设置Actor朝向
PHYSXCWRAP_API bool PxActorSetOrientation(PxHandle actor, Quat orient);

// 测试单步移动
PHYSXCWRAP_API StepRes PxActorCheckStep(PxHandle actor, TransForm pose);

// 单步移动
PHYSXCWRAP_API StepRes PxActorStep(PxHandle actor, TransForm pose);

// 统计进入的范围触发器数量
PHYSXCWRAP_API int32_t PxActorCountInTraps(PxHandle actor);

// 查询进入范围触发器状态
PHYSXCWRAP_API InTrapData PxActorGetInTrapData(PxHandle actor, int32_t index);

// 设置进入范围触发器状态
PHYSXCWRAP_API bool PxActorSetInTrapStat(PxHandle actor, int32_t index, InTrapStat stat);

// 删除进入的范围触发器数据
PHYSXCWRAP_API bool PxActorDeleteInTrapData(PxHandle actor, int32_t index);

#if __cplusplus
}
#endif																				    











