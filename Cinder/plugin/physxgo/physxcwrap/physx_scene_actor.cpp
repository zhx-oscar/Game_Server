#include "physx_scene.h"
#include "physx_actor.h"
#include "physx_actor_trap.h"
#include "common.h"

using namespace physx;

PxActorWrap* CreateActorWrap(ActorModeFlags actorModeFlags)
{
	if ((actorModeFlags & ActorMode::eTrap) != 0)
	{
		return new PxTrapWrap();
	}

	return new PxActorWrap();
}

#define INIT_ACTOR(pActor, upright) \
auto pActorWrap = CreateActorWrap(actorModeFlags); \
CHECK_NULL(pActorWrap, { (pActor)->release(); }, NULL); \
CHECK_FALSE(pActorWrap->Init(this, (pActor), (upright), actorModeFlags, hitFilter, bindGoObj), { delete pActorWrap; (pActor)->release(); }, NULL); 

PxActorWrap* PxSceneWrap::CreatePlane(const PxVec3& normalvec, PxReal distance, PxFilterData hitFilter, void* bindGoObj)
{
	CHECK_FALSE(InitOK(), , NULL);

	auto plane = PxCreatePlane(mPxScene->getPhysics(), PxPlane(normalvec, distance), *mPxMaterial);
	CHECK_NULL(plane, , NULL);

	ActorModeFlags actorModeFlags = eNotBeTrap;
	INIT_ACTOR(plane, false);

	return pActorWrap;
}

PxActorWrap* PxSceneWrap::CreateHeightField(const std::vector<int16_t>& heightmap, unsigned columns, unsigned rows, const PxVec3& scale, PxFilterData hitFilter, void* bindGoObj)
{
	CHECK_FALSE(InitOK(), , NULL);

	PxHeightFieldGeometry hfGeom;
	CHECK_FALSE(BuildHeightFieldGeometry(hfGeom, heightmap, columns, rows, scale), , NULL);

	return CreateHeightField(hfGeom, hitFilter, bindGoObj);
}

PxActorWrap* PxSceneWrap::CreateHeightField(const PxHeightFieldGeometry& hfGeom, PxFilterData hitFilter, void* bindGoObj)
{
	CHECK_FALSE(InitOK(), , NULL);

	auto columns = hfGeom.heightField->getNbColumns();
	auto rows = hfGeom.heightField->getNbRows();
	PxTransform pose = PxTransform(PxIdentity);
	pose.p = PxVec3(-(float(columns) / 2 * hfGeom.columnScale), 0, -(float(rows) / 2 * hfGeom.rowScale));

	auto hfActor = mPxScene->getPhysics().createRigidStatic(pose);
	CHECK_NULL(hfActor, , NULL);

	auto hfShape = PxRigidActorExt::createExclusiveShape(*hfActor, hfGeom, *mPxMaterial);
	CHECK_NULL(hfActor, { hfActor->release(); }, NULL);

	ActorModeFlags actorModeFlags = eNotBeTrap;
	INIT_ACTOR(hfActor, false);

	return pActorWrap;
}

PxActorWrap* PxSceneWrap::CreateBoxKinematic(const PxTransform& pose, const PxVec3& halfExtents, ActorModeFlags actorModeFlags, PxFilterData hitFilter, void* bindGoObj)
{
	CHECK_FALSE(InitOK(), , NULL);

	auto box = PxCreateKinematic(mPxScene->getPhysics(), pose, PxBoxGeometry(halfExtents), *mPxMaterial, 1.0f);
	CHECK_NULL(box, , NULL);

	INIT_ACTOR(box, false);	

	return pActorWrap;
}

PxActorWrap* PxSceneWrap::CreateBoxStatic(const PxTransform& pose, const PxVec3& halfExtents, ActorModeFlags actorModeFlags, PxFilterData hitFilter, void* bindGoObj)
{
	CHECK_FALSE(InitOK(), , NULL);

	auto box = PxCreateStatic(mPxScene->getPhysics(), pose, PxBoxGeometry(halfExtents), *mPxMaterial);
	CHECK_NULL(box, , NULL);

	INIT_ACTOR(box, false);	

	return pActorWrap;
}

PxActorWrap* PxSceneWrap::CreateSphereKinematic(const PxTransform& pose, PxReal radius, ActorModeFlags actorModeFlags, PxFilterData hitFilter, void* bindGoObj)
{
	CHECK_FALSE(InitOK(), , NULL);

	auto sphere = PxCreateKinematic(mPxScene->getPhysics(), pose, PxSphereGeometry(radius), *mPxMaterial, 1.0f);
	CHECK_NULL(sphere, , NULL);

	INIT_ACTOR(sphere, true);

	return pActorWrap;
}

PxActorWrap* PxSceneWrap::CreateSphereStatic(const PxTransform& pose, PxReal radius, ActorModeFlags actorModeFlags, PxFilterData hitFilter, void* bindGoObj)
{
	CHECK_FALSE(InitOK(), , NULL);

	auto sphere = PxCreateStatic(mPxScene->getPhysics(), pose, PxSphereGeometry(radius), *mPxMaterial);
	CHECK_NULL(sphere, , NULL);

	INIT_ACTOR(sphere, true);

	return pActorWrap;
}

PxActorWrap* PxSceneWrap::CreateCapsuleKinematic(const PxTransform& pose, PxReal radius, PxReal halfHeight, ActorModeFlags actorModeFlags, PxFilterData hitFilter, void* bindGoObj)
{
	CHECK_FALSE(InitOK(), , NULL);

	auto capsule = PxCreateKinematic(mPxScene->getPhysics(), pose, PxCapsuleGeometry(radius, halfHeight), *mPxMaterial, 1.0f);
	CHECK_NULL(capsule, , NULL);

	INIT_ACTOR(capsule, true);

	return pActorWrap;
}

PxActorWrap* PxSceneWrap::CreateCapsuleStatic(const PxTransform& pose, PxReal radius, PxReal halfHeight, ActorModeFlags actorModeFlags, PxFilterData hitFilter, void* bindGoObj)
{
	CHECK_FALSE(InitOK(), , NULL);

	auto capsule = PxCreateStatic(mPxScene->getPhysics(), pose, PxCapsuleGeometry(radius, halfHeight), *mPxMaterial);
	CHECK_NULL(capsule, , NULL);

	INIT_ACTOR(capsule, true);

	return pActorWrap;
}

PxActorWrap* PxSceneWrap::CreateMeshKinematic(const PxTransform& pose, const PxVec3& scale, const std::vector<PxReal>& vb, const std::vector<uint16_t>& ib, ActorModeFlags actorModeFlags, PxFilterData hitFilter, void* bindGoObj)
{
	CHECK_FALSE(InitOK(), , NULL);

	PxTriangleMeshGeometry triGeom;
	CHECK_FALSE(BuildMeshGeometry(triGeom, scale, vb, ib), , NULL);

	return CreateMeshKinematic(pose, triGeom, actorModeFlags, hitFilter, bindGoObj);
}

PxActorWrap* PxSceneWrap::CreateMeshKinematic(const PxTransform& pose, const PxTriangleMeshGeometry& triGeom, ActorModeFlags actorModeFlags, PxFilterData hitFilter, void* bindGoObj)
{
	CHECK_FALSE(InitOK(), , NULL);

	auto mesh = PxCreateKinematic(mPxScene->getPhysics(), pose, triGeom, *mPxMaterial, 1.0f);
	CHECK_NULL(mesh, , NULL);

	INIT_ACTOR(mesh, false);

	return pActorWrap;
}

PxActorWrap* PxSceneWrap::CreateMeshStatic(const PxTransform& pose, const PxVec3& scale, const std::vector<PxReal>& vb, const std::vector<uint16_t>& ib, ActorModeFlags actorModeFlags, PxFilterData hitFilter, void* bindGoObj)
{
	CHECK_FALSE(InitOK(), , NULL);

	PxTriangleMeshGeometry triGeom;
	CHECK_FALSE(BuildMeshGeometry(triGeom, scale, vb, ib), , NULL);

	return CreateMeshStatic(pose, triGeom, actorModeFlags, hitFilter, bindGoObj);
}

PxActorWrap* PxSceneWrap::CreateMeshStatic(const PxTransform& pose, const PxTriangleMeshGeometry& triGeom, ActorModeFlags actorModeFlags, PxFilterData hitFilter, void* bindGoObj)
{
	CHECK_FALSE(InitOK(), , NULL);

	auto mesh = PxCreateStatic(mPxScene->getPhysics(), pose, triGeom, *mPxMaterial);
	CHECK_NULL(mesh, , NULL);

	INIT_ACTOR(mesh, false);

	return pActorWrap;
}

// 放置Actor
bool PxSceneWrap::AddActor(PxActorWrap* actor)
{
	CHECK_FALSE(InitOK(), , false);
	CHECK_NULL(actor, , false);
	CHECK_FALSE(actor->InitOK(), , false);
	CHECK_TRUE(actor->mInScene, , false);

	PxSceneWrapWriteLock lock(this); 
	mPxScene->addActor(*actor->mPxActor);	
	actor->mInScene = true;

	actor->ScanEnterTraps();

	return true;
}