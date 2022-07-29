#pragma once

#include "PxPhysicsAPI.h"
#include "PhysXCWrap.h"
#include "physx_wrap.h"
#include "physx_queryfilter.h"
#include <vector>

class PxSceneWrap;
class PxActorWrap;

// 在范围触发器中数据
struct InTrap
{
	PxActorWrap* TrapActor;  // 范围触发器Actor
	InTrapStat SelfStat;     // 自身状态
};

// PxActor包装
class PxActorWrap : public PxWrap
{
public:
	PxActorWrap();
	virtual ~PxActorWrap();

	// 查询包装类型
	virtual PxWrap::PxType GetPxType();

	// 初始化
	virtual bool Init(PxSceneWrap* pSceneWrap, physx::PxRigidActor* pActor, bool upRight, 
		ActorModeFlags actorModeFlags, physx::PxFilterData hitFilter, void* bindGoObj);

	// 初始化是否成功
	virtual bool InitOK();

	// 释放
	virtual void Release();

	// 获取Actor绑定的Go对象
	virtual void* GetBindGoObj();

	// 设置模式标记
	virtual bool SetActorModeFlags(ActorModeFlags actorModeFlags);

	// 查询模式标记
	virtual ActorModeFlags GetActorModeFlags();

	// 设置过滤器
	virtual void SetHitFilter(physx::PxFilterData hitFilter);

	// 查询过滤器
	virtual physx::PxFilterData GetHitFilter();

	// 获取坐标与角度
	virtual physx::PxTransform GetPose();

	// 设置坐标与角度
	virtual bool SetPose(const physx::PxTransform& pose);

	// 测试单步移动
	virtual StepRes CheckStep(const physx::PxTransform& pose);

	// 单步移动
	virtual StepRes Step(const physx::PxTransform& pose);

	// 扫描进入的范围触发器
	virtual void ScanEnterTraps();

	// 统计进入的范围触发器数量
	virtual int CountInTraps();

	// 访问进入的范围触发器数据
	virtual InTrap* GetInTrap(int index);

	// 删除进入的范围触发器数据
	virtual void DeleteInTrap(int index);

public:	
	physx::PxRigidActor* mPxActor;
	PxSceneWrap* mSceneWrap;
	void* mBindGoObj;
	ActorModeFlags mActorModeFlags;
	bool mInScene;
	PxQueryFilterEnterTraps mPxQFEnterTraps;
	PxQueryFilterStep mPxQFStep;
	std::vector<InTrap> mInTraps;
};


