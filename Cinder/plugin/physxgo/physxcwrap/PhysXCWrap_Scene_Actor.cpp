#include "PhysXCWrap.h"
#include "physx_sdk.h"
#include "physx_scene.h"
#include "physx_actor.h"
#include "common.h"

using namespace physx;

// 创建Kinematic Box
PHYSXCWRAP_API PxHandle PxSceneCreateBoxKinematic(PxHandle scene, TransForm pose, Vector3 halfExtents, ActorModeFlags actorModeFlags, 
	HitFilter hitFilter, void* bindGoObj)
{
	auto pSceneWrap = ConvertPxWrap<PxSceneWrap>(scene);
	CHECK_NULL(pSceneWrap, , NULL);

	return pSceneWrap->CreateBoxKinematic(PxTransform(PxVec3(pose.P.X, pose.P.Y, pose.P.Z), PxQuat(pose.Q.X, pose.Q.Y, pose.Q.Z, pose.Q.W)), PxVec3(halfExtents.X, halfExtents.Y, halfExtents.Z),
		actorModeFlags, PxFilterData(hitFilter.Word0, hitFilter.Word1, hitFilter.Word2, hitFilter.Word3), bindGoObj);
}

// 创建Static Box
PHYSXCWRAP_API PxHandle PxSceneCreateBoxStatic(PxHandle scene, TransForm pose, Vector3 halfExtents, ActorModeFlags actorModeFlags, 
	HitFilter hitFilter, void* bindGoObj)
{
	auto pSceneWrap = ConvertPxWrap<PxSceneWrap>(scene);
	CHECK_NULL(pSceneWrap, , NULL);


	return pSceneWrap->CreateBoxStatic(PxTransform(PxVec3(pose.P.X, pose.P.Y, pose.P.Z), PxQuat(pose.Q.X, pose.Q.Y, pose.Q.Z, pose.Q.W)), PxVec3(halfExtents.X, halfExtents.Y, halfExtents.Z),
		actorModeFlags, PxFilterData(hitFilter.Word0, hitFilter.Word1, hitFilter.Word2, hitFilter.Word3), bindGoObj);
}

// 创建Kinematic Sphere
PHYSXCWRAP_API PxHandle PxSceneCreateSphereKinematic(PxHandle scene, TransForm pose, float radius, ActorModeFlags actorModeFlags, 
	HitFilter hitFilter, void* bindGoObj)
{
	auto pSceneWrap = ConvertPxWrap<PxSceneWrap>(scene);
	CHECK_NULL(pSceneWrap, , NULL);

	return pSceneWrap->CreateSphereKinematic(PxTransform(PxVec3(pose.P.X, pose.P.Y, pose.P.Z), PxQuat(pose.Q.X, pose.Q.Y, pose.Q.Z, pose.Q.W)), radius,
		actorModeFlags, PxFilterData(hitFilter.Word0, hitFilter.Word1, hitFilter.Word2, hitFilter.Word3), bindGoObj);
}

// 创建Static Sphere
PHYSXCWRAP_API PxHandle PxSceneCreateSphereStatic(PxHandle scene, TransForm pose, float radius, ActorModeFlags actorModeFlags, 
	HitFilter hitFilter, void* bindGoObj)
{
	auto pSceneWrap = ConvertPxWrap<PxSceneWrap>(scene);
	CHECK_NULL(pSceneWrap, , NULL);

	return pSceneWrap->CreateSphereStatic(PxTransform(PxVec3(pose.P.X, pose.P.Y, pose.P.Z), PxQuat(pose.Q.X, pose.Q.Y, pose.Q.Z, pose.Q.W)), radius,
		actorModeFlags, PxFilterData(hitFilter.Word0, hitFilter.Word1, hitFilter.Word2, hitFilter.Word3), bindGoObj);
}

// 创建Kinematic Capsule
PHYSXCWRAP_API PxHandle PxSceneCreateCapsuleKinematic(PxHandle scene, TransForm pose, float radius, float halfHeight, ActorModeFlags actorModeFlags, 
	HitFilter hitFilter, void* bindGoObj)
{
	auto pSceneWrap = ConvertPxWrap<PxSceneWrap>(scene);
	CHECK_NULL(pSceneWrap, , NULL);

	return pSceneWrap->CreateCapsuleKinematic(PxTransform(PxVec3(pose.P.X, pose.P.Y, pose.P.Z), PxQuat(pose.Q.X, pose.Q.Y, pose.Q.Z, pose.Q.W)), radius, halfHeight,
		actorModeFlags, PxFilterData(hitFilter.Word0, hitFilter.Word1, hitFilter.Word2, hitFilter.Word3), bindGoObj);
}

// 创建Static Capsule
PHYSXCWRAP_API PxHandle PxSceneCreateCapsuleStatic(PxHandle scene, TransForm pose, float radius, float halfHeight, ActorModeFlags actorModeFlags, 
	HitFilter hitFilter, void* bindGoObj)
{
	auto pSceneWrap = ConvertPxWrap<PxSceneWrap>(scene);
	CHECK_NULL(pSceneWrap, , NULL);

	return pSceneWrap->CreateCapsuleStatic(PxTransform(PxVec3(pose.P.X, pose.P.Y, pose.P.Z), PxQuat(pose.Q.X, pose.Q.Y, pose.Q.Z, pose.Q.W)), radius, halfHeight,
		actorModeFlags, PxFilterData(hitFilter.Word0, hitFilter.Word1, hitFilter.Word2, hitFilter.Word3), bindGoObj);
}

// Scene放置Actor
PHYSXCWRAP_API bool PxSceneAddActor(PxHandle scene, PxHandle actor)
{
	auto pSceneWrap = ConvertPxWrap<PxSceneWrap>(scene);
	CHECK_NULL(pSceneWrap, , NULL);

	auto pActorWrap = ConvertPxWrap<PxActorWrap>(actor);
	CHECK_NULL(pActorWrap, , NULL);

	return pSceneWrap->AddActor(pActorWrap);
}