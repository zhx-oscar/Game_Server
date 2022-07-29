#include "physx_sdk.h"
#include "physx_scene.h"
#include "common.h"

using namespace physx;

PxSdkWrap::PxSdkWrap()
: mPxFoundation(NULL), mPxPvd(NULL), mPxPhysics(NULL), mPxCooking(NULL)
{
}

PxSdkWrap::~PxSdkWrap()
{
	Release();
}

// 查询包装类型
PxWrap::PxType PxSdkWrap::GetPxType()
{
	return PxWrap::PxSdk;
}

// 初始化
bool PxSdkWrap::Init(bool usePvd, const char* host, int port, unsigned int timeoutInMs)
{
	CHECK_TRUE(InitOK(), , true);	

	mPxFoundation = PxCreateFoundation(PX_PHYSICS_VERSION, mPxAllocator, mPxErrorCallback);
	CHECK_NULL(mPxFoundation, { Release(); }, false);

	if (usePvd)
	{
		mPxPvd = PxCreatePvd(*mPxFoundation);
		CHECK_NULL(mPxPvd, { Release(); }, false);

		auto transport = PxDefaultPvdSocketTransportCreate(host, port, timeoutInMs);
		CHECK_NULL(transport, { Release(); }, false);

		CHECK_FALSE(mPxPvd->connect(*transport, PxPvdInstrumentationFlag::eALL), { Release(); }, false);
	}

	mPxPhysics = PxCreatePhysics(PX_PHYSICS_VERSION, *mPxFoundation, PxTolerancesScale(), false, mPxPvd);
	CHECK_NULL(mPxPhysics, { Release(); }, false);

	mPxCooking = PxCreateCooking(PX_PHYSICS_VERSION, *mPxFoundation, PxCookingParams(PxTolerancesScale()));
	CHECK_NULL(mPxCooking, { Release(); }, false);	

	return true;
}

// 是否初始化成功
bool PxSdkWrap::InitOK()
{
	return mPxPhysics != NULL;
}

// 释放
void PxSdkWrap::Release()
{
	PX_RELEASE(mPxPhysics);
	if (mPxPvd != NULL)
	{
		auto transport = mPxPvd->getTransport();
		PX_RELEASE(mPxPvd);
		PX_RELEASE(transport);
	}
	PX_RELEASE(mPxCooking);
	PX_RELEASE(mPxFoundation);
}

// 创建物理Scene
PxSceneWrap* PxSdkWrap::CreateScene(bool multiThread, void* bindGoObj)
{
	auto pSceneWrap = new PxSceneWrap();
	CHECK_NULL(pSceneWrap, , NULL);

	CHECK_FALSE(pSceneWrap->Init(this, multiThread, bindGoObj), { delete pSceneWrap; }, NULL);

	return pSceneWrap;
}