#pragma once

#include "physx_actor.h"
#include <unordered_set>

class PxTrapWrap : public PxActorWrap 
{
public:
	PxTrapWrap();
	virtual ~PxTrapWrap();

	// 初始化
	virtual bool Init(PxSceneWrap* pSceneWrap, physx::PxRigidActor* pActor, bool upRight,
		ActorModeFlags actorModeFlags, physx::PxFilterData hitFilter, void* bindGoObj);

	// 释放
	virtual void Release();

	// 设置模式标记
	virtual bool SetActorModeFlags(ActorModeFlags actorModeFlags);

	// 设置过滤器
	virtual void SetHitFilter(physx::PxFilterData hitFilter);

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
	std::unordered_set<PxActorWrap*> mCatchActors;
};