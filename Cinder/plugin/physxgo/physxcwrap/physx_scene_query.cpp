#include "physx_scene.h"
#include "physx_actor.h"
#include "common.h"

using namespace physx;

// 射线检测
bool PxSceneWrap::Raycast(const PxVec3& origin, const PxVec3& unitDir, const PxReal distance, PxRaycastCallback& hitBuf, PxHitFlags hitFlags, 
	const PxQueryFilterData& filter, PxQueryFilterCallback* filterCall, const PxQueryCache* cache)
{
	CHECK_FALSE(InitOK(), , false);

	PxSceneWrapReadLock lock(this);
	return mPxScene->raycast(origin, unitDir, distance, hitBuf, hitFlags, filter, filterCall, cache);
}

// 滑动检测
bool PxSceneWrap::Sweep(const PxGeometry& geometry, const PxTransform& pose, const PxVec3& unitDir, const PxReal distance,
	PxSweepCallback& hitBuf, PxHitFlags hitFlags, const PxQueryFilterData& filter, PxQueryFilterCallback* filterCall,
	const PxQueryCache* cache, const PxReal inflation, bool upright)
{
	CHECK_FALSE(InitOK(), , false);

	if (upright)
	{
		switch (geometry.getType())
		{
		case PxGeometryType::eCAPSULE:
		case PxGeometryType::eSPHERE:
			pose.rotate(PxVec3(0.0f, 0.0f, PxHalfPi));
			break;
		default:
			break;
		}
	}

	PxSceneWrapReadLock lock(this);
	return mPxScene->sweep(geometry, pose, unitDir, distance, hitBuf, hitFlags, filter, filterCall, cache, inflation);
}

// 重叠检测
bool PxSceneWrap::Overlap(const PxGeometry& geometry, const PxTransform& pose, PxOverlapCallback& hitBuf,
	const PxQueryFilterData& filter, PxQueryFilterCallback* filterCall, bool upright)
{
	CHECK_FALSE(InitOK(), , false);

	if (upright)
	{
		switch (geometry.getType())
		{
		case PxGeometryType::eCAPSULE:
		case PxGeometryType::eSPHERE:
			pose.rotate(PxVec3(0.0f, 0.0f, PxHalfPi));
			break;
		default:
			break;
		}
	}

	PxSceneWrapReadLock lock(this);
	return mPxScene->overlap(geometry, pose, hitBuf, filter, filterCall);
}