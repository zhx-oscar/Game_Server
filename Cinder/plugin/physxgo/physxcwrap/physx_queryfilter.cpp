#include "physx_queryfilter.h"
#include "physx_actor.h"
#include "common.h"

using namespace physx;

PxQueryHitType::Enum PxQueryFilterIgnoreTrap::preFilter(const PxFilterData& filterData, const PxShape* shape, const PxRigidActor* actor, PxHitFlags& queryFlags)
{
	auto pActorWrap = (PxActorWrap*)actor->userData;
	CHECK_NULL(pActorWrap, , PxQueryHitType::eNONE);

	return (pActorWrap->mActorModeFlags & ActorMode::eTrap) != 0 ? PxQueryHitType::eNONE : PxQueryHitType::eBLOCK;
}

PxQueryHitType::Enum PxQueryFilterIgnoreTrap::postFilter(const PxFilterData& filterData, const PxQueryHit& hit)
{
	return PxQueryHitType::eBLOCK;
}

PxQueryHitType::Enum PxQueryFilterEnterTraps::preFilter(const PxFilterData& filterData, const PxShape* shape, const PxRigidActor* actor, PxHitFlags& queryFlags)
{
	CHECK_NULL(mSelf, , PxQueryHitType::eNONE);

	auto pActorWrap = (PxActorWrap*)(actor->userData);
	CHECK_NULL(pActorWrap, , PxQueryHitType::eNONE);

	// 排除自己
	if (mSelf == actor->userData)
	{
		return PxQueryHitType::eNONE;
	}

	// 排除不是范围触发器
	if ((pActorWrap->mActorModeFlags & ActorMode::eTrap) == 0)
	{
		return PxQueryHitType::eNONE;
	}

	return PxQueryHitType::eBLOCK;
}

PxQueryHitType::Enum PxQueryFilterEnterTraps::postFilter(const PxFilterData& filterData, const PxQueryHit& hit)
{
	CHECK_NULL(mSelf, , PxQueryHitType::eNONE);

	auto pActorWrap = (PxActorWrap*)(hit.actor->userData);
	CHECK_NULL(pActorWrap, , PxQueryHitType::eNONE);

	// 记录进入的范围触发器
	mEnterTraps.push_back(pActorWrap);

	return PxQueryHitType::eNONE;
}

physx::PxQueryHitType::Enum PxQueryFilterStep::preFilter(const physx::PxFilterData& filterData, const physx::PxShape* shape, const physx::PxRigidActor* actor, physx::PxHitFlags& queryFlags)
{
	CHECK_NULL(mSelf, , PxQueryHitType::eNONE);

	auto pActorWrap = (PxActorWrap*)(actor->userData);
	CHECK_NULL(pActorWrap, , PxQueryHitType::eNONE);

	// 排除自己
	if (mSelf == actor->userData)
	{
		return PxQueryHitType::eNONE;
	}

	// 排除范围触发器
	if ((pActorWrap->mActorModeFlags & ActorMode::eTrap) != 0)
	{
		return PxQueryHitType::eNONE;
	}

	return PxQueryHitType::eBLOCK;
}

physx::PxQueryHitType::Enum PxQueryFilterStep::postFilter(const physx::PxFilterData& filterData, const physx::PxQueryHit& hit)
{
	return PxQueryHitType::eBLOCK;
}

