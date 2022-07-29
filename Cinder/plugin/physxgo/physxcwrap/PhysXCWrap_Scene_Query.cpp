#include "PhysXCWrap.h"
#include "physx_sdk.h"
#include "physx_scene.h"
#include "physx_queryfilter.h"
#include "common.h"

using namespace physx;

template<class T>
void ExportHit(Hit& hit, const T& pxHit)
{
	hit.Target = pxHit.actor->userData; 
	hit.Position.X = pxHit.position.x; 
	hit.Position.Y = pxHit.position.y; 
	hit.Position.Z = pxHit.position.z; 
	hit.Normal.X = pxHit.normal.x; 
	hit.Normal.Y = pxHit.normal.y; 
	hit.Normal.Z = pxHit.normal.z; 
	hit.Distance = pxHit.distance;
}

template<>
void ExportHit(Hit& hit, const PxOverlapHit& pxHit)
{
	hit.Target = pxHit.actor->userData;
}

void InitPxFilter(const HitModeFlags& hitModeFlags, const HitFilter& hitFilter, PxQueryFilterData& pxFilter)
{
	if ((hitModeFlags & HitMode::eHitStatic) != 0)
	{
		pxFilter.flags |= PxQueryFlag::eSTATIC;
	}
	else
	{
		pxFilter.flags &= ~PxQueryFlag::eSTATIC;
	}

	if ((hitModeFlags & HitMode::eHitDynamic) != 0)
	{
		pxFilter.flags |= PxQueryFlag::eDYNAMIC;
	}
	else
	{
		pxFilter.flags &= ~PxQueryFlag::eDYNAMIC;
	}

	if ((hitModeFlags & HitMode::eAnyHit) != 0)
	{
		pxFilter.flags |= PxQueryFlag::eANY_HIT;
	} 
	else
	{
		pxFilter.flags &= ~PxQueryFlag::eANY_HIT;
	}

	if ((hitModeFlags & HitMode::eAllHit) != 0)
	{
		pxFilter.flags |= PxQueryFlag::eNO_BLOCK;
	}
	else
	{
		pxFilter.flags &= ~PxQueryFlag::eNO_BLOCK;
	}

	if ((hitModeFlags & HitMode::eHitTrap) == 0)
	{
		pxFilter.flags |= PxQueryFlag::ePREFILTER;
	} 
	else
	{
		pxFilter.flags &= ~PxQueryFlag::ePREFILTER;
	}

	// 过滤类型
	pxFilter.data.word0 = hitFilter.Word0;
	pxFilter.data.word1 = hitFilter.Word1;
	pxFilter.data.word2 = hitFilter.Word2;
	pxFilter.data.word3 = hitFilter.Word3;
}

int GeometrySweep(PxSceneWrap* pSceneWrap, const Geometry& geom, const TransForm& pose, const Vector3& unitDir, float distance, float inflation, 
	PxHitFlags pxHitFlags, const PxQueryFilterData& pxFilter, PxSweepBuffer& pxHitBuf)
{
	// 过滤器场景查询忽略范围触发器
	static PxQueryFilterIgnoreTrap pxQueryFilterIgnoreTrap;

	// 滑动检测
	switch (geom.Type)
	{
	case eSPHERE:
		if (!pSceneWrap->Sweep(PxSphereGeometry(geom.Radius), PxTransform(PxVec3(pose.P.X, pose.P.Y, pose.P.Z), PxQuat(pose.Q.X, pose.Q.Y, pose.Q.Z, pose.Q.W)),
			PxVec3(unitDir.X, unitDir.Y, unitDir.Z), distance, pxHitBuf, pxHitFlags, pxFilter, &pxQueryFilterIgnoreTrap, NULL, inflation, true))
		{
			return 0;
		}
		break;
	case eCAPSULE:
		if (!pSceneWrap->Sweep(PxCapsuleGeometry(geom.Radius, geom.HalfHeight), PxTransform(PxVec3(pose.P.X, pose.P.Y, pose.P.Z), PxQuat(pose.Q.X, pose.Q.Y, pose.Q.Z, pose.Q.W)),
			PxVec3(unitDir.X, unitDir.Y, unitDir.Z), distance, pxHitBuf, pxHitFlags, pxFilter, &pxQueryFilterIgnoreTrap, NULL, inflation, true))
		{
			return 0;
		}
		break;
	case eBOX:
		if (!pSceneWrap->Sweep(PxBoxGeometry(PxVec3(geom.HalfExtents.X, geom.HalfExtents.Y, geom.HalfExtents.Z)), PxTransform(PxVec3(pose.P.X, pose.P.Y, pose.P.Z), PxQuat(pose.Q.X, pose.Q.Y, pose.Q.Z, pose.Q.W)),
			PxVec3(unitDir.X, unitDir.Y, unitDir.Z), distance, pxHitBuf, pxHitFlags, pxFilter, &pxQueryFilterIgnoreTrap, NULL, inflation))
		{
			return 0;
		}
		break;
	default:
		return -1;
	}

	return 1;
}

int GeometryOverlap(PxSceneWrap* pSceneWrap, const Geometry& geom, const TransForm& pose, const PxQueryFilterData& pxFilter, PxOverlapCallback& pxHitBuf)
{
	// 过滤器场景查询忽略范围触发器
	static PxQueryFilterIgnoreTrap pxQueryFilterIgnoreTrap;

	// 重叠检测
	switch (geom.Type)
	{
	case eSPHERE:
		if (!pSceneWrap->Overlap(PxSphereGeometry(geom.Radius), PxTransform(PxVec3(pose.P.X, pose.P.Y, pose.P.Z), PxQuat(pose.Q.X, pose.Q.Y, pose.Q.Z, pose.Q.W)),
			pxHitBuf, pxFilter, &pxQueryFilterIgnoreTrap, true))
		{
			return 0;
		}
		break;
	case eCAPSULE:
		if (!pSceneWrap->Overlap(PxCapsuleGeometry(geom.Radius, geom.HalfHeight), PxTransform(PxVec3(pose.P.X, pose.P.Y, pose.P.Z), PxQuat(pose.Q.X, pose.Q.Y, pose.Q.Z, pose.Q.W)),
			pxHitBuf, pxFilter, &pxQueryFilterIgnoreTrap, true))
		{
			return 0;
		}
		break;
	case eBOX:
		if (!pSceneWrap->Overlap(PxBoxGeometry(PxVec3(geom.HalfExtents.X, geom.HalfExtents.Y, geom.HalfExtents.Z)), PxTransform(PxVec3(pose.P.X, pose.P.Y, pose.P.Z), PxQuat(pose.Q.X, pose.Q.Y, pose.Q.Z, pose.Q.W)),
			pxHitBuf, pxFilter, &pxQueryFilterIgnoreTrap))
		{
			return 0;
		}
		break;
	default:
		return -1;
	}

	return 1;
}

// 射线检测（返回值-1：失败，[0,N]：碰撞数量）
PHYSXCWRAP_API int32_t PxSceneRaycast(PxHandle scene, Vector3 origin, Vector3 unitDir, float distance,
	HitModeFlags hitModeFlags, HitFilter hitFilter, Hit* hitBuf, size_t hitBufLen)
{
	memset(hitBuf, 0, sizeof(Hit) * hitBufLen);

	auto pSceneWrap = ConvertPxWrap<PxSceneWrap>(scene);
	CHECK_NULL(pSceneWrap, , -1);

	// 过滤器
	PxQueryFilterData pxFilter;
	InitPxFilter(hitModeFlags, hitFilter, pxFilter);

	// 碰撞参数 
	PxHitFlags pxHitFlags;
	pxHitFlags = PxHitFlag::ePOSITION | PxHitFlag::eNORMAL;

	// 过滤器场景查询忽略范围触发器
	static PxQueryFilterIgnoreTrap pxQueryFilterIgnoreTrap;

	// 只需单个碰撞结果
	if (pxFilter.flags.isSet(PxQueryFlag::eANY_HIT) || !pxFilter.flags.isSet(PxQueryFlag::eNO_BLOCK) || hitBufLen <= 1)
	{
		// 碰撞结果
		PxRaycastBuffer pxHitBuf;		

		// 射线检测
		if (!pSceneWrap->Raycast(PxVec3(origin.X, origin.Y, origin.Z), PxVec3(unitDir.X, unitDir.Y, unitDir.Z), distance,
			pxHitBuf, pxHitFlags, pxFilter, &pxQueryFilterIgnoreTrap))
		{
			return 0;
		}

		if (pxHitBuf.getNbAnyHits() <= 0)
		{
			return 0;
		}

		if (hitBuf != NULL && hitBufLen > 0)
		{			
			ExportHit(hitBuf[0], pxHitBuf.getAnyHit(0));
		}

		return 1;
	}

	// 需要多个碰撞结果
	PxRaycastBuffer pxHitBuf(new PxRaycastHit[hitBufLen], (PxU32)hitBufLen);

	// 射线检测
	if (!pSceneWrap->Raycast(PxVec3(origin.X, origin.Y, origin.Z), PxVec3(unitDir.X, unitDir.Y, unitDir.Z), distance,
		pxHitBuf, pxHitFlags, pxFilter, &pxQueryFilterIgnoreTrap))
	{
		delete pxHitBuf.touches;
		return 0;
	}

	auto pxHitCount = pxHitBuf.getNbAnyHits();

	if (hitBuf != NULL && hitBufLen > 0)
	{
		for (int i = 0; i < (int)hitBufLen && i < (int)pxHitCount; i++)
		{
			ExportHit(hitBuf[i], pxHitBuf.getAnyHit(i));
		}
	}

	delete pxHitBuf.touches;
	return (int32_t)(pxHitCount > hitBufLen ? hitBufLen : pxHitCount);
}

// 滑动检测（返回值-1：失败，[0,N]：碰撞数量）
PHYSXCWRAP_API int32_t PxSceneSweep(PxHandle scene, Geometry geom, TransForm pose, Vector3 unitDir, float distance, float inflation,
	HitModeFlags hitModeFlags, HitFilter hitFilter, Hit* hitBuf, size_t hitBufLen)
{
	memset(hitBuf, 0, sizeof(Hit) * hitBufLen);

	auto pSceneWrap = ConvertPxWrap<PxSceneWrap>(scene);
	CHECK_NULL(pSceneWrap, , -1);

	// 过滤器
	PxQueryFilterData pxFilter;
	InitPxFilter(hitModeFlags, hitFilter, pxFilter);

	// 碰撞参数 
	PxHitFlags pxHitFlags;
	pxHitFlags = PxHitFlag::ePOSITION | PxHitFlag::eNORMAL;

	// 只需单个碰撞结果
	if (pxFilter.flags.isSet(PxQueryFlag::eANY_HIT) || !pxFilter.flags.isSet(PxQueryFlag::eNO_BLOCK) || hitBufLen <= 1)
	{
		// 碰撞结果
		PxSweepBuffer pxHitBuf;

		// 滑动检测
		int rv = GeometrySweep(pSceneWrap, geom, pose, unitDir, distance, inflation, pxHitFlags, pxFilter, pxHitBuf);
		if (rv <= 0)
		{
			return rv;
		}

		if (pxHitBuf.getNbAnyHits() <= 0)
		{
			return 0;
		}

		if (hitBuf != NULL && hitBufLen > 0)
		{
			ExportHit(hitBuf[0], pxHitBuf.getAnyHit(0));
		}

		return 1;
	}

	// 需要多个碰撞结果
	PxSweepBuffer pxHitBuf(new PxSweepHit[hitBufLen], (PxU32)hitBufLen);

	// 滑动检测
	int rv = GeometrySweep(pSceneWrap, geom, pose, unitDir, distance, inflation, pxHitFlags, pxFilter, pxHitBuf);
	if (rv <= 0)
	{
		return rv;
	}

	auto pxHitCount = pxHitBuf.getNbAnyHits();

	if (hitBuf != NULL && hitBufLen > 0)
	{
		for (int i = 0; i < (int)hitBufLen && i < (int)pxHitCount; i++)
		{
			ExportHit(hitBuf[i], pxHitBuf.getAnyHit(i));
		}
	}

	delete pxHitBuf.touches;
	return (int32_t)(pxHitCount > hitBufLen ? hitBufLen : pxHitCount);
}

// 重叠检测（返回值-1：失败，[0,N]：碰撞数量）
PHYSXCWRAP_API int32_t PxSceneOverlap(PxHandle scene, Geometry geom, TransForm pose,
	HitModeFlags hitModeFlags, HitFilter hitFilter, Hit* hitBuf, size_t hitBufLen)
{
	memset(hitBuf, 0, sizeof(Hit) * hitBufLen);

	auto pSceneWrap = ConvertPxWrap<PxSceneWrap>(scene);
	CHECK_NULL(pSceneWrap, , -1);

	// 过滤器
	PxQueryFilterData pxFilter;
	InitPxFilter(hitModeFlags, hitFilter, pxFilter);

	// 只需单个碰撞结果
	if (pxFilter.flags.isSet(PxQueryFlag::eANY_HIT) || !pxFilter.flags.isSet(PxQueryFlag::eNO_BLOCK) || hitBufLen <= 1)
	{
		// 重叠检测单个必须设置eANY_HIT
		pxFilter.flags |= PxQueryFlag::eANY_HIT;

		// 碰撞结果
		PxOverlapBuffer pxHitBuf;

		// 重叠检测
		int rv = GeometryOverlap(pSceneWrap, geom, pose, pxFilter, pxHitBuf);
		if (rv <= 0)
		{
			return rv;
		}

		if (pxHitBuf.getNbAnyHits() <= 0)
		{
			return 0;
		}

		if (hitBuf != NULL && hitBufLen > 0)
		{	
			ExportHit(hitBuf[0], pxHitBuf.getAnyHit(0));
		}

		return 1;
	}

	// 需要多个碰撞结果
	PxOverlapBuffer pxHitBuf(new PxOverlapHit[hitBufLen], (PxU32)hitBufLen);

	// 重叠检测
	int rv = GeometryOverlap(pSceneWrap, geom, pose, pxFilter, pxHitBuf);
	if (rv <= 0)
	{
		return rv;
	}

	auto pxHitCount = pxHitBuf.getNbAnyHits();

	if (hitBuf != NULL && hitBufLen > 0)
	{
		for (int i = 0; i < (int)hitBufLen && i < (int)pxHitCount; i++)
		{
			ExportHit(hitBuf[i], pxHitBuf.getAnyHit(i));
		}
	}

	delete pxHitBuf.touches;
	return (int32_t)(pxHitCount > hitBufLen ? hitBufLen : pxHitCount);
}