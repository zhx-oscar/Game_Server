#include "physx_actor.h"
#include "physx_scene.h"
#include "physx_queryfilter.h"
#include "physx_actor_trap.h"
#include "common.h"
#include <algorithm>

using namespace physx;

PxActorWrap::PxActorWrap()
: mPxActor(NULL), mSceneWrap(NULL), mBindGoObj(NULL), mActorModeFlags(ActorModeFlags(0)), mInScene(false),
mPxQFEnterTraps(this), mPxQFStep(this)
{
}

PxActorWrap::~PxActorWrap()
{
	Release();
}

// 查询包装类型
PxWrap::PxType PxActorWrap::GetPxType()
{
	return PxWrap::PxActor;
}

//  初始化
bool PxActorWrap::Init(PxSceneWrap* pSceneWrap, PxRigidActor* pActor, bool upRight, ActorModeFlags actorModeFlags, PxFilterData hitFilter, void* bindGoObj)
{
	CHECK_TRUE(InitOK(), , true);
	CHECK_NULL(pSceneWrap, , false);
	CHECK_NULL(pActor, , false);

	mSceneWrap = pSceneWrap;
	mPxActor = pActor;
	mPxActor->userData = this;
	mBindGoObj = bindGoObj;
	mActorModeFlags = actorModeFlags;

	pActor->setActorFlag(PxActorFlag::eDISABLE_GRAVITY, true); 
	pActor->setActorFlag(PxActorFlag::eDISABLE_SIMULATION, true); 
		
	PxU32 nbShapes = pActor->getNbShapes(); 
	for (PxU32 i = 0; i < nbShapes; i++) 
	{ 
		PxShape* shape; 
		PxU32 n = pActor->getShapes(&shape, 1, i); 
		if (n != 1) 
		{ 
			continue; 
		} 
	
		if (((actorModeFlags & ActorMode::eNotBeQuery) != 0) && ((actorModeFlags & ActorMode::eNotBeTrap) != 0))
		{
			shape->setFlag(PxShapeFlag::eSCENE_QUERY_SHAPE, false);
		}

		shape->setQueryFilterData(hitFilter); 
		
		if (upRight)
		{ 
			PxTransform relativePose(PxQuat(PxHalfPi, PxVec3(0, 0, 1))); 
			shape->setLocalPose(relativePose); 
		} 
	} 

	return true;
}

// 初始化是否成功
bool PxActorWrap::InitOK()
{
	return mPxActor != NULL;
}

// 释放
void PxActorWrap::Release()
{
	if (mPxActor != NULL)
	{
		for (auto it = mInTraps.begin(); it != mInTraps.end(); it++)
		{
			auto pTrapWrap = (PxTrapWrap*)(it->TrapActor);

			pTrapWrap->mCatchActors.erase(this);
		}

		mInTraps.clear();

		PX_RELEASE(mPxActor);
	}
}

// 获取Actor绑定的Go对象
void* PxActorWrap::GetBindGoObj()
{
	CHECK_FALSE(InitOK(), , NULL);

	return mBindGoObj;
}

// 设置模式标记
bool PxActorWrap::SetActorModeFlags(ActorModeFlags actorModeFlags)
{	
	CHECK_FALSE(InitOK(), , false);

	// 调整是否不能被场景查询
	{
		bool newFlag = ((actorModeFlags & ActorMode::eNotBeQuery) != 0);
		bool oldFlag = ((mActorModeFlags & ActorMode::eNotBeQuery) != 0);

		if (newFlag != oldFlag)
		{
			bool canQueryShape = !(newFlag && ((mActorModeFlags & ActorMode::eNotBeTrap) != 0));

			{
				PxSceneWrapWriteLock lock(mInScene ? mSceneWrap : NULL);

				PxU32 nbShapes = mPxActor->getNbShapes();
				for (PxU32 i = 0; i < nbShapes; i++)
				{
					PxShape* shape;
					PxU32 n = mPxActor->getShapes(&shape, 1, i);
					if (n != 1)
					{
						continue;
					}

					shape->setFlag(PxShapeFlag::eSCENE_QUERY_SHAPE, canQueryShape);
				}
			}

			if (newFlag)
			{
				mActorModeFlags |= ActorMode::eNotBeQuery;
			}
			else
			{
				mActorModeFlags &= ~ActorMode::eNotBeQuery;
			}
		}
	}

	// 调整不能被区域触发器捕获
	{
		bool newFlag = ((actorModeFlags & ActorMode::eNotBeTrap) != 0);
		bool oldFlag = ((mActorModeFlags & ActorMode::eNotBeTrap) != 0);

		if (newFlag != oldFlag)
		{
			bool canQueryShape = !(newFlag && ((mActorModeFlags & ActorMode::eNotBeQuery) != 0));

			{
				PxSceneWrapWriteLock lock(mInScene ? mSceneWrap : NULL);

				PxU32 nbShapes = mPxActor->getNbShapes();
				for (PxU32 i = 0; i < nbShapes; i++)
				{
					PxShape* shape;
					PxU32 n = mPxActor->getShapes(&shape, 1, i);
					if (n != 1)
					{
						continue;
					}

					shape->setFlag(PxShapeFlag::eSCENE_QUERY_SHAPE, canQueryShape);
				}
			}

			if (newFlag)
			{
				mActorModeFlags |= ActorMode::eNotBeTrap;
			}
			else
			{
				mActorModeFlags &= ~ActorMode::eNotBeTrap;
			}

			ScanEnterTraps();
		}
	}

	return true;
}

// 查询模式标记
ActorModeFlags PxActorWrap::GetActorModeFlags()
{
	CHECK_FALSE(InitOK(), , 0);

	return mActorModeFlags;
}

// 设置过滤器
void PxActorWrap::SetHitFilter(PxFilterData hitFilter)
{
	CHECK_FALSE(InitOK(), , );

	PxSceneWrapWriteLock lock(mInScene ? mSceneWrap : NULL);

	PxU32 nbShapes = mPxActor->getNbShapes();
	for (PxU32 i = 0; i < nbShapes; i++)
	{
		PxShape* shape;
		PxU32 n = mPxActor->getShapes(&shape, 1, i);
		if (n != 1)
		{
			continue;
		}

		shape->setQueryFilterData(hitFilter);
	}
}

// 查询过滤器
PxFilterData PxActorWrap::GetHitFilter()
{
	CHECK_FALSE(InitOK(), , PxFilterData());


	PxSceneWrapReadLock lock(mInScene ? mSceneWrap : NULL);

	if (mPxActor->getNbShapes() <= 0)
	{
		return PxFilterData();
	}

	PxShape* shape;
	PxU32 n = mPxActor->getShapes(&shape, 1);
	if (n != 1)
	{
		return PxFilterData();
	}

	return shape->getQueryFilterData();
}

// 获取坐标与角度
PxTransform PxActorWrap::GetPose()
{
	CHECK_FALSE(InitOK(), , PxTransform());

	PxSceneWrapReadLock lock(mInScene ? mSceneWrap : NULL);

	return mPxActor->getGlobalPose();
}

// 设置坐标与角度	
bool PxActorWrap::SetPose(const PxTransform& pose)
{
	CHECK_FALSE(InitOK(), , false);

	{
		PxSceneWrapWriteLock lock(mInScene ? mSceneWrap : NULL);

		if (mPxActor->getConcreteType() != PxConcreteType::eRIGID_DYNAMIC)
		{
			return false;
		}

		mPxActor->setGlobalPose(pose, false);
	}

	ScanEnterTraps();

	return true;
}

// 测试单步移动
StepRes PxActorWrap::CheckStep(const physx::PxTransform& pose)
{
	CHECK_FALSE(InitOK(), , (StepRes{ NULL, ZeroTransForm(), Vector3{0,0,0}, false }));

	PxSceneWrapReadLock lock(mInScene ? mSceneWrap : NULL);

	// 可以被捕获
	if (mPxActor->getNbShapes() <= 0)
	{
		return StepRes{ NULL, ZeroTransForm(), Vector3{0,0,0}, false };
	}

	PxShape* shape;
	PxU32 n = mPxActor->getShapes(&shape, 1);
	if (n != 1)
	{
		return StepRes{ NULL, ZeroTransForm(), Vector3{0,0,0}, false };
	}

	auto geom = shape->getGeometry().any();
	PxSweepBuffer hitBuf;
	PxHitFlags hitFlags;
	hitFlags = PxHitFlag::ePOSITION | PxHitFlag::eNORMAL;
	PxQueryFilterData filter;	
	filter.flags |= PxQueryFlag::ePREFILTER;
	auto stepRes = StepRes{ NULL, ZeroTransForm(), Vector3{0,0,0}, false };
	bool rv = false;	
	auto curPose = PxShapeExt::getGlobalPose(*shape, *mPxActor);
	auto delta = pose.p - curPose.p;	
	auto unitDir = delta.getNormalized();

	switch (geom.getType())
	{
	case PxGeometryType::eBOX:
		rv = mSceneWrap->Sweep(shape->getGeometry().box(), curPose, unitDir, delta.magnitude(), hitBuf, hitFlags, filter, &mPxQFStep, NULL, 0.0f, false);
		break;
	case PxGeometryType::eCAPSULE:
		rv = mSceneWrap->Sweep(shape->getGeometry().capsule(), curPose, unitDir, delta.magnitude(), hitBuf, hitFlags, filter, &mPxQFStep, NULL, 0.0f, false);
		break;
	case PxGeometryType::eSPHERE:
		rv = mSceneWrap->Sweep(shape->getGeometry().sphere(), curPose, unitDir, delta.magnitude(), hitBuf, hitFlags, filter, &mPxQFStep, NULL, 0.0f, false);
		break;
	case PxGeometryType::eCONVEXMESH:
		rv = mSceneWrap->Sweep(shape->getGeometry().convexMesh(), curPose, unitDir, delta.magnitude(), hitBuf, hitFlags, filter, &mPxQFStep, NULL, 0.0f, false);
		break;
	default:
		return stepRes;
	}

	if (!rv || hitBuf.getNbAnyHits() <= 0)
	{
		stepRes.Ok = true;
		return stepRes;
	}

	auto hit = hitBuf.getAnyHit(0);

	stepRes.BlockActor = hit.actor->userData;
	stepRes.BlockPose.P.X = hit.position.x;
	stepRes.BlockPose.P.Y = hit.position.y;
	stepRes.BlockPose.P.Z = hit.position.z;
	stepRes.BlockPose.Q.X = curPose.q.x;
	stepRes.BlockPose.Q.Y = curPose.q.y;
	stepRes.BlockPose.Q.Z = curPose.q.z;
	stepRes.BlockPose.Q.W = curPose.q.w;
	stepRes.BlockNormal.X = hit.normal.x;
	stepRes.BlockNormal.Y = hit.normal.y;
	stepRes.BlockNormal.Z = hit.normal.z;

	return stepRes;
}

// 单步移动
StepRes PxActorWrap::Step(const physx::PxTransform& pose)
{	
	CHECK_FALSE(InitOK(), , (StepRes{ NULL, ZeroTransForm(), Vector3{0,0,0}, false }));

	auto stepRes = CheckStep(pose);
	if (!stepRes.Ok) 
	{
		return stepRes;
	}

	PxTransform newPose = pose;

	if 	(stepRes.BlockActor != NULL)
	{
		newPose.p.x = stepRes.BlockPose.P.X;
		newPose.p.y = stepRes.BlockPose.P.Y;
		newPose.p.z = stepRes.BlockPose.P.Z;
		newPose.q.x = stepRes.BlockPose.Q.X;
		newPose.q.y = stepRes.BlockPose.Q.Y;
		newPose.q.z = stepRes.BlockPose.Q.Z;
		newPose.q.w = stepRes.BlockPose.Q.W;
	}	

	{
		PxSceneWrapWriteLock lock(mInScene ? mSceneWrap : NULL);

		mPxActor->setGlobalPose(newPose);
	}

	ScanEnterTraps();

	return stepRes;
}

// 扫描进入的范围触发器
void PxActorWrap::ScanEnterTraps()
{
	CHECK_FALSE(InitOK(), , );

	if ((mActorModeFlags & ActorMode::eNotBeTrap) != 0)
	{
		// 不能被捕获
		for (auto it = mInTraps.rbegin(); it != mInTraps.rend(); it++)
		{
			switch (it->SelfStat)
			{
			case InTrapStat::eEnter:
				mInTraps.erase(it.base());
				break;
			case InTrapStat::eStay:
				it->SelfStat = InTrapStat::eLeave;
				break;			
			default:
				break;
			}
		}
	}
	else
	{
		PxSceneWrapReadLock lock(mInScene ? mSceneWrap : NULL);		

		// 可以被捕获
		if (mPxActor->getNbShapes() <= 0)
		{
			return;
		}

		PxShape* shape;
		PxU32 n = mPxActor->getShapes(&shape, 1);
		if (n != 1)
		{
			return;
		}

		auto geom = shape->getGeometry().any();
		auto pose = PxShapeExt::getGlobalPose(*shape, *mPxActor);
		PxOverlapBufferN<1> hitBuf;
		PxQueryFilterData filter;
		filter.flags |= PxQueryFlag::ePREFILTER | PxQueryFlag::ePOSTFILTER;
		mPxQFEnterTraps.mEnterTraps.clear();

		// 通过重叠查询触发器
		switch (geom.getType())
		{
		case PxGeometryType::eBOX:
			mSceneWrap->Overlap(shape->getGeometry().box(), pose, hitBuf, filter, &mPxQFEnterTraps, false);
			break;
		case PxGeometryType::eCAPSULE:
			mSceneWrap->Overlap(shape->getGeometry().capsule(), pose, hitBuf, filter, &mPxQFEnterTraps, false);
			break;
		case PxGeometryType::eSPHERE:
			mSceneWrap->Overlap(shape->getGeometry().sphere(), pose, hitBuf, filter, &mPxQFEnterTraps, false);
			break;
		case PxGeometryType::eCONVEXMESH:
			mSceneWrap->Overlap(shape->getGeometry().convexMesh(), pose, hitBuf, filter, &mPxQFEnterTraps, false);
			break;
		default:
			return;
		}

		// 检测离开的触发器
		for (auto it = mInTraps.rbegin(); it != mInTraps.rend(); it++)
		{
			// 已经离开触发器
			if (std::find(mPxQFEnterTraps.mEnterTraps.begin(), mPxQFEnterTraps.mEnterTraps.end(), it->TrapActor)
				== mPxQFEnterTraps.mEnterTraps.end())
			{
				switch (it->SelfStat)
				{
				case InTrapStat::eEnter:
					mInTraps.erase(it.base());
					break;
				case InTrapStat::eStay:
					it->SelfStat = InTrapStat::eLeave;
					break;
				default:
					break;
				}
			}
		}

		// 检测进入触发器
		for (auto it = mPxQFEnterTraps.mEnterTraps.begin(); it != mPxQFEnterTraps.mEnterTraps.end(); it++)
		{
			auto pTrapWrap = (PxTrapWrap*)(*it);

			// 检测是否已记录本Actor
			auto trapIt = pTrapWrap->mCatchActors.find(this);
			if (trapIt != pTrapWrap->mCatchActors.end())
			{
				continue;
			}

			// 进入触发器
			mInTraps.push_back(InTrap{ pTrapWrap, InTrapStat::eEnter });
			pTrapWrap->mCatchActors.insert(this);
		}		
	}
}

// 统计进入的范围触发器数量
int PxActorWrap::CountInTraps()
{
	CHECK_FALSE(InitOK(), , 0);
	return (int)mInTraps.size();
}

// 访问进入的范围触发器
InTrap* PxActorWrap::GetInTrap(int index)
{
	CHECK_FALSE(InitOK(), , NULL);

	if (index < 0 || index >= mInTraps.size())
	{
		return NULL;
	}

	return &(mInTraps[index]);
}

// 删除进入的范围触发器数据
void PxActorWrap::DeleteInTrap(int index)
{
	CHECK_FALSE(InitOK(), , );

	if (index < 0 || index >= mInTraps.size())
	{
		return;
	}

	auto pTrapWrap = (PxTrapWrap*)(mInTraps[index].TrapActor);

	mInTraps.erase(mInTraps.begin() + index);
	pTrapWrap->mCatchActors.erase(this);
}
