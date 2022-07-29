#pragma once

#include "PxPhysicsAPI.h"  
#include <vector>

class PxActorWrap;

// 过滤器场景查询忽略范围触发器
class PxQueryFilterIgnoreTrap : public physx::PxQueryFilterCallback
{
public:
	physx::PxQueryHitType::Enum preFilter(const physx::PxFilterData& filterData, const physx::PxShape* shape, const physx::PxRigidActor* actor, physx::PxHitFlags& queryFlags);
	physx::PxQueryHitType::Enum postFilter(const physx::PxFilterData& filterData, const physx::PxQueryHit& hit);
};

// 过滤器Actor扫描进入的范围触发器
class PxQueryFilterEnterTraps : public physx::PxQueryFilterCallback
{
public:
	PxQueryFilterEnterTraps(PxActorWrap* self) : mSelf(self) {}	

	physx::PxQueryHitType::Enum preFilter(const physx::PxFilterData& filterData, const physx::PxShape* shape, const physx::PxRigidActor* actor, physx::PxHitFlags& queryFlags);
	physx::PxQueryHitType::Enum postFilter(const physx::PxFilterData& filterData, const physx::PxQueryHit& hit);

public:			  
	PxActorWrap* mSelf;
	std::vector<PxActorWrap*> mEnterTraps;
};

// 过滤器Actor单步移动
class PxQueryFilterStep : public physx::PxQueryFilterCallback
{
public:
	PxQueryFilterStep(PxActorWrap* self) : mSelf(self) {}

	physx::PxQueryHitType::Enum preFilter(const physx::PxFilterData& filterData, const physx::PxShape* shape, const physx::PxRigidActor* actor, physx::PxHitFlags& queryFlags);
	physx::PxQueryHitType::Enum postFilter(const physx::PxFilterData& filterData, const physx::PxQueryHit& hit);

public:
	PxActorWrap* mSelf;
};