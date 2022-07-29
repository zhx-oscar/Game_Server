#pragma once

#include "PxPhysicsAPI.h"	
#include "PhysXCWrap.h"
#include "physx_wrap.h"
#include <vector>

class PxSdkWrap;
class PxActorWrap;

// PxScene包装
class PxSceneWrap : public PxWrap
{
public:
	PxSceneWrap();
	virtual ~PxSceneWrap();

	// 查询包装类型
	PxWrap::PxType GetPxType();

	// 初始化
	bool Init(PxSdkWrap* pSdkWrap, bool multiThread, void* bindGoObj);

	// 初始化是否成功
	bool InitOK();

	// 释放
	void Release();

	// 帧更新
	void Update(physx::PxReal elapsedTime);	

	// 创建Actor
	PxActorWrap* CreatePlane(const physx::PxVec3& normalvec, physx::PxReal distance, physx::PxFilterData hitFilter, void* bindGoObj);
	PxActorWrap* CreateHeightField(const std::vector<int16_t>& heightmap, unsigned columns, unsigned rows, const physx::PxVec3& scale, physx::PxFilterData hitFilter, void* bindGoObj);
	PxActorWrap* CreateHeightField(const physx::PxHeightFieldGeometry& hfGeom, physx::PxFilterData hitFilter, void* bindGoObj);
	PxActorWrap* CreateBoxKinematic(const physx::PxTransform& pose, const physx::PxVec3& halfExtents, ActorModeFlags actorModeFlags, physx::PxFilterData hitFilter, void* bindGoObj);
	PxActorWrap* CreateBoxStatic(const physx::PxTransform& pose, const physx::PxVec3& halfExtents, ActorModeFlags actorModeFlags, physx::PxFilterData hitFilter, void* bindGoObj);
	PxActorWrap* CreateSphereKinematic(const physx::PxTransform& pose, physx::PxReal radius, ActorModeFlags actorModeFlags, physx::PxFilterData hitFilter, void* bindGoObj);
	PxActorWrap* CreateSphereStatic(const physx::PxTransform& pose, physx::PxReal radius, ActorModeFlags actorModeFlags, physx::PxFilterData hitFilter, void* bindGoObj);
	PxActorWrap* CreateCapsuleKinematic(const physx::PxTransform& pose, physx::PxReal radius, physx::PxReal halfHeight, ActorModeFlags actorModeFlags, physx::PxFilterData hitFilter, void* bindGoObj);
	PxActorWrap* CreateCapsuleStatic(const physx::PxTransform& pose, physx::PxReal radius, physx::PxReal halfHeight, ActorModeFlags actorModeFlags, physx::PxFilterData hitFilter, void* bindGoObj);
	PxActorWrap* CreateMeshKinematic(const physx::PxTransform& pose, const physx::PxVec3& scale, const std::vector<physx::PxReal>& vb, const std::vector<uint16_t>& ib, ActorModeFlags actorModeFlags, physx::PxFilterData hitFilter, void* bindGoObj);
	PxActorWrap* CreateMeshKinematic(const physx::PxTransform& pose, const physx::PxTriangleMeshGeometry& triGeom, ActorModeFlags actorModeFlags, physx::PxFilterData hitFilter, void* bindGoObj);
	PxActorWrap* CreateMeshStatic(const physx::PxTransform& pose, const physx::PxVec3& scale, const std::vector<physx::PxReal>& vb, const std::vector<uint16_t>& ib, ActorModeFlags actorModeFlags, physx::PxFilterData hitFilter, void* bindGoObj);
	PxActorWrap* CreateMeshStatic(const physx::PxTransform& pose, const physx::PxTriangleMeshGeometry& triGeom, ActorModeFlags actorModeFlags, physx::PxFilterData hitFilter, void* bindGoObj);

	// 放置Actor
	bool AddActor(PxActorWrap* actor);

	// 构造复杂几何体
	bool BuildHeightFieldGeometry(physx::PxHeightFieldGeometry& geom, const std::vector<int16_t>& heightmap, unsigned columns, unsigned rows, const physx::PxVec3& scale);
	bool BuildMeshGeometry(physx::PxTriangleMeshGeometry& geom, const physx::PxVec3& scale, const std::vector<float>& vb, const std::vector<uint16_t>& ib);

	// 射线检测
	bool Raycast(const physx::PxVec3& origin, const physx::PxVec3& unitDir, const physx::PxReal distance, 
		physx::PxRaycastCallback& hitBuf, physx::PxHitFlags hitFlags, 
		const physx::PxQueryFilterData& filter, physx::PxQueryFilterCallback* filterCall = NULL,
		const physx::PxQueryCache* cache = NULL);

	// 滑动检测
	bool Sweep(const physx::PxGeometry& geometry, const physx::PxTransform& pose, const physx::PxVec3& unitDir, const physx::PxReal distance,
		physx::PxSweepCallback& hit, physx::PxHitFlags hitFlags, 
		const physx::PxQueryFilterData& filter, physx::PxQueryFilterCallback* filterCall = NULL,
		const physx::PxQueryCache* cache = NULL, const physx::PxReal inflation = 0.0f, bool upright = true);

	// 重叠检测
	bool Overlap(const physx::PxGeometry& geometry, const physx::PxTransform& pose, physx::PxOverlapCallback& hit, 
		const physx::PxQueryFilterData& filter, physx::PxQueryFilterCallback* filterCall = NULL, bool upright = true);

public:
	physx::PxDefaultCpuDispatcher* mPxCpuDispatcher;
	physx::PxScene* mPxScene;
	physx::PxMaterial* mPxMaterial;
	PxSdkWrap* mPxSdkWrap;
	bool mMultiThread;
	void* mBindGoObj;
};

class PxSceneWrapWriteLock
{
public:
	PxSceneWrapWriteLock(PxSceneWrap* pxSceneWrap)
	: mPxScene(NULL)
	{
		if (pxSceneWrap == NULL)
		{
			return;
		}

		if (pxSceneWrap->mMultiThread)
		{
			mPxScene = pxSceneWrap->mPxScene;
		}

		if (mPxScene != NULL)
		{
			mPxScene->lockWrite();
		}
	}
	~PxSceneWrapWriteLock()
	{
		if (mPxScene != NULL)
		{
			mPxScene->unlockWrite();
		}
	}

private:
	physx::PxScene* mPxScene;
};

class PxSceneWrapReadLock
{
public:
	PxSceneWrapReadLock(PxSceneWrap* pxSceneWrap)
	: mPxScene(NULL)
	{
		if (pxSceneWrap == NULL)
		{
			return;
		}

		if (pxSceneWrap->mMultiThread)
		{
			mPxScene = pxSceneWrap->mPxScene;
		}

		if (mPxScene != NULL)
		{
			mPxScene->lockRead();
		}
	}
	~PxSceneWrapReadLock()
	{
		if (mPxScene != NULL)
		{
			mPxScene->unlockRead();
		}
	}

private:
	physx::PxScene* mPxScene;
};