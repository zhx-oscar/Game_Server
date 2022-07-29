#include "physx_actor_trap.h"

using namespace physx;

PxTrapWrap::PxTrapWrap()
{
}

PxTrapWrap::~PxTrapWrap()
{
}

// 初始化
bool PxTrapWrap::Init(PxSceneWrap* pSceneWrap, physx::PxRigidActor* pActor, bool upRight,
	ActorModeFlags actorModeFlags, physx::PxFilterData hitFilter, void* bindGoObj)
{
	return PxActorWrap::Init(pSceneWrap, pActor, upRight, ActorMode::eTrap, PxFilterData(PxEmpty), bindGoObj);
}

// 释放
void PxTrapWrap::Release()
{
	if (mPxActor != NULL)
	{
		for (auto it = mCatchActors.begin(); it != mCatchActors.end(); it++)
		{
			auto pActorWrap = (PxActorWrap*)(*it);

			for (auto itj = pActorWrap->mInTraps.begin(); itj != pActorWrap->mInTraps.end(); itj++)
			{
				if (itj->TrapActor == this) 
				{
					pActorWrap->mInTraps.erase(itj);
					break;
				}
			}
		}
	}

	PxActorWrap::Release();
}

// 设置模式标记
bool PxTrapWrap::SetActorModeFlags(ActorModeFlags actorModeFlags)
{
	return false;
}

// 设置过滤器
void PxTrapWrap::SetHitFilter(physx::PxFilterData hitFilter)
{
}

// 设置坐标与角度
bool PxTrapWrap::SetPose(const physx::PxTransform& pose)
{
	return false;
}

// 测试单步移动
StepRes PxTrapWrap::CheckStep(const physx::PxTransform& pose)
{
	return StepRes{ NULL, ZeroTransForm(), Vector3{0,0,0}, false };
}

// 单步移动
StepRes PxTrapWrap::Step(const physx::PxTransform& pose)
{
	return StepRes{ NULL, ZeroTransForm(), Vector3{0,0,0}, false };
}

// 扫描进入的范围触发器
void PxTrapWrap::ScanEnterTraps()
{	
}

// 统计进入的范围触发器数量
int PxTrapWrap::CountInTraps()
{
	return 0;
}

// 访问进入的范围触发器数据
InTrap* PxTrapWrap::GetInTrap(int index)
{
	return NULL;
}

// 删除进入的范围触发器数据
void PxTrapWrap::DeleteInTrap(int index)
{
}
