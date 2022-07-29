#include "PhysXCWrap.h"
#include "physx_actor.h"
#include "common.h"

using namespace physx;

// 销毁Actor
PHYSXCWRAP_API bool ReleasePxActor(PxHandle actor)
{
	auto pActorWrap = ConvertPxWrap<PxActorWrap>(actor);
	CHECK_NULL(pActorWrap, , false);

	pActorWrap->Release();

	return true;
}

// 获取Actor绑定的Go对象
PHYSXCWRAP_API void* PxActorGetBindGoObj(PxHandle actor)
{
	auto pActorWrap = ConvertPxWrap<PxActorWrap>(actor);
	CHECK_NULL(pActorWrap, , NULL);

	return pActorWrap->GetBindGoObj();
}

// 设置Actor模式标记
PHYSXCWRAP_API bool PxActorSetActorModeFlags(PxHandle actor, ActorModeFlags actorModeFlags)
{
	auto pActorWrap = ConvertPxWrap<PxActorWrap>(actor);
	CHECK_NULL(pActorWrap, , false);

	return pActorWrap->SetActorModeFlags(actorModeFlags);
}

// 获取Actor模式标记
PHYSXCWRAP_API ActorModeFlags PxActorGetActorModeFlags(PxHandle actor)
{
	auto pActorWrap = ConvertPxWrap<PxActorWrap>(actor);
	CHECK_NULL(pActorWrap, , 0);

	return pActorWrap->mActorModeFlags;
}

// 设置Actor过滤器
PHYSXCWRAP_API bool PxActorSetHitFilter(PxHandle actor, HitFilter hitFilter)
{
	auto pActorWrap = ConvertPxWrap<PxActorWrap>(actor);
	CHECK_NULL(pActorWrap, , false);

	pActorWrap->SetHitFilter(PxFilterData(hitFilter.Word0, hitFilter.Word1, hitFilter.Word2, hitFilter.Word3));

	return true;
}

// 获取Actor过滤器
PHYSXCWRAP_API HitFilter PxActorGetHitFilter(PxHandle actor)
{
	HitFilter hitFilter = HitFilter{ 0, 0, 0, 0 };

	auto pActorWrap = ConvertPxWrap<PxActorWrap>(actor);
	CHECK_NULL(pActorWrap, , hitFilter);

	auto pxHitFilter = pActorWrap->GetHitFilter();

	hitFilter.Word0 = pxHitFilter.word0;
	hitFilter.Word1 = pxHitFilter.word1;
	hitFilter.Word2 = pxHitFilter.word2;
	hitFilter.Word3 = pxHitFilter.word3;

	return hitFilter;
}

// 获取Actor坐标与朝向
PHYSXCWRAP_API TransForm PxActorGetPose(PxHandle actor)
{
	auto pose = ZeroTransForm();

	auto pActorWrap = ConvertPxWrap<PxActorWrap>(actor);
	CHECK_NULL(pActorWrap, , pose);

	auto pxPose = pActorWrap->GetPose();
	
	pose.P.X = pxPose.p.x;
	pose.P.Y = pxPose.p.y;
	pose.P.Z = pxPose.p.z;
	pose.Q.X = pxPose.q.x;
	pose.Q.Y = pxPose.q.y;
	pose.Q.Z = pxPose.q.z;
	pose.Q.W = pxPose.q.w;

	return pose;
}

// 设置Actor设置坐标与朝向
PHYSXCWRAP_API bool PxActorSetPose(PxHandle actor, TransForm pose)
{
	auto pActorWrap = ConvertPxWrap<PxActorWrap>(actor);
	CHECK_NULL(pActorWrap, , false);

	PxTransform pxPose;
	pxPose.p.x = pose.P.X;
	pxPose.p.y = pose.P.Y;
	pxPose.p.z = pose.P.Z;
	pxPose.q.x = pose.Q.X;
	pxPose.q.y = pose.Q.Y;
	pxPose.q.z = pose.Q.Z;
	pxPose.q.w = pose.Q.W;

	pActorWrap->SetPose(pxPose);

	return true;
}

// 设置Actor坐标
PHYSXCWRAP_API bool PxActorSetPosition(PxHandle actor, Vector3 pos)
{
	auto pActorWrap = ConvertPxWrap<PxActorWrap>(actor);
	CHECK_NULL(pActorWrap, , false);

	auto pxPose = pActorWrap->GetPose();
	pxPose.p.x = pos.X;
	pxPose.p.y = pos.Y;
	pxPose.p.z = pos.Z;

	pActorWrap->SetPose(pxPose);

	return true;
}

// 设置Actor朝向
PHYSXCWRAP_API bool PxActorSetOrientation(PxHandle actor, Quat orient)
{
	auto pActorWrap = ConvertPxWrap<PxActorWrap>(actor);
	CHECK_NULL(pActorWrap, , false);

	auto pxPose = pActorWrap->GetPose();
	pxPose.q.x = orient.X;
	pxPose.q.y = orient.Y;
	pxPose.q.z = orient.Z;
	pxPose.q.w = orient.W;

	pActorWrap->SetPose(pxPose);

	return true;
}

// 测试单步移动
PHYSXCWRAP_API StepRes PxActorCheckStep(PxHandle actor, TransForm pose)
{
	auto pActorWrap = ConvertPxWrap<PxActorWrap>(actor);
	CHECK_NULL(pActorWrap, , (StepRes{ NULL, ZeroTransForm(), Vector3{0,0,0}, false }));

	PxTransform pxPose;
	pxPose.p.x = pose.P.X;
	pxPose.p.y = pose.P.Y;
	pxPose.p.z = pose.P.Z;
	pxPose.q.x = pose.Q.X;
	pxPose.q.y = pose.Q.Y;
	pxPose.q.z = pose.Q.Z;
	pxPose.q.w = pose.Q.W;

	return pActorWrap->CheckStep(pxPose);
}

// 单步移动
PHYSXCWRAP_API StepRes PxActorStep(PxHandle actor, TransForm pose)
{
	auto pActorWrap = ConvertPxWrap<PxActorWrap>(actor);
	CHECK_NULL(pActorWrap, , (StepRes{ NULL, ZeroTransForm(), Vector3{0,0,0}, false }));

	PxTransform pxPose;
	pxPose.p.x = pose.P.X;
	pxPose.p.y = pose.P.Y;
	pxPose.p.z = pose.P.Z;
	pxPose.q.x = pose.Q.X;
	pxPose.q.y = pose.Q.Y;
	pxPose.q.z = pose.Q.Z;
	pxPose.q.w = pose.Q.W;

	return pActorWrap->Step(pxPose);
}

// 统计进入的范围触发器数量
PHYSXCWRAP_API int32_t PxActorCountInTraps(PxHandle actor)
{
	auto pActorWrap = ConvertPxWrap<PxActorWrap>(actor);
	CHECK_NULL(pActorWrap, , 0);

	return pActorWrap->CountInTraps();
}

// 查询进入范围触发器状态
PHYSXCWRAP_API InTrapData PxActorGetInTrapData(PxHandle actor, int32_t index)
{
	auto pActorWrap = ConvertPxWrap<PxActorWrap>(actor);
	CHECK_NULL(pActorWrap, , (InTrapData{NULL, InTrapStat::eNone}));

	auto pInTrap = pActorWrap->GetInTrap(index);
	if (NULL == pInTrap)
	{
		return InTrapData{ NULL, InTrapStat::eNone };
	}

	return InTrapData{ pInTrap->TrapActor, pInTrap->SelfStat };
}

// 设置进入范围触发器状态
PHYSXCWRAP_API bool PxActorSetInTrapStat(PxHandle actor, int32_t index, InTrapStat stat)
{
	auto pActorWrap = ConvertPxWrap<PxActorWrap>(actor);
	CHECK_NULL(pActorWrap, , false);

	auto pInTrap = pActorWrap->GetInTrap(index);
	if (NULL == pInTrap)
	{
		return false;
	}

	pInTrap->SelfStat = stat;

	return true;

}

// 删除进入的范围触发器数据
PHYSXCWRAP_API bool PxActorDeleteInTrapData(PxHandle actor, int32_t index)
{
	auto pActorWrap = ConvertPxWrap<PxActorWrap>(actor);
	CHECK_NULL(pActorWrap, , false);

	pActorWrap->DeleteInTrap(index);

	return true;
}
