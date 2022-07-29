#include "physx_scene.h"
#include "physx_sdk.h"
#include "common.h"
#include <thread>

using namespace physx;

PxSceneWrap::PxSceneWrap()
: mPxCpuDispatcher(NULL), mPxScene(NULL), mPxMaterial(NULL), mPxSdkWrap(NULL), mMultiThread(false), mBindGoObj(NULL)
{
}

PxSceneWrap::~PxSceneWrap()
{
	Release();
}

// 查询包装类型
PxWrap::PxType PxSceneWrap::GetPxType()
{
	return PxWrap::PxScene;
}

// 初始化
bool PxSceneWrap::Init(PxSdkWrap* pSdkWrap, bool multiThread, void* bindGoObj)
{
	CHECK_TRUE(InitOK(), , true);
	CHECK_NULL(pSdkWrap, , false);
	CHECK_FALSE(pSdkWrap->InitOK(), , false);

	PxSceneDesc sceneDesc(pSdkWrap->mPxPhysics->getTolerancesScale());

	mMultiThread = multiThread;
	if (mMultiThread)
	{
		sceneDesc.flags |= PxSceneFlag::eREQUIRE_RW_LOCK;
	}

	mPxCpuDispatcher = PxDefaultCpuDispatcherCreate(std::thread::hardware_concurrency());
	CHECK_NULL(mPxCpuDispatcher, , false);

	sceneDesc.cpuDispatcher = mPxCpuDispatcher;
	sceneDesc.filterShader = PxDefaultSimulationFilterShader;

	mPxScene = pSdkWrap->mPxPhysics->createScene(sceneDesc);
	CHECK_NULL(mPxScene, { Release(); }, false);

	mPxMaterial = pSdkWrap->mPxPhysics->createMaterial(0.0f, 0.0f, 0.0f);
	CHECK_NULL(mPxMaterial, { Release(); }, false);	

	mPxScene->userData = this;
	mPxSdkWrap = pSdkWrap;
	mBindGoObj = bindGoObj;

	PxPvdSceneClient* pvdClient = mPxScene->getScenePvdClient();
	if (pvdClient)
	{
		pvdClient->setScenePvdFlag(PxPvdSceneFlag::eTRANSMIT_CONSTRAINTS, true);
		pvdClient->setScenePvdFlag(PxPvdSceneFlag::eTRANSMIT_CONTACTS, true);
		pvdClient->setScenePvdFlag(PxPvdSceneFlag::eTRANSMIT_SCENEQUERIES, true);
	}

	return true;
}

// 初始化是否成功
bool PxSceneWrap::InitOK()
{
	return mPxScene != NULL;
}

// 释放
void PxSceneWrap::Release()
{
	PX_RELEASE(mPxMaterial);
	PX_RELEASE(mPxScene);
	PX_RELEASE(mPxCpuDispatcher);
	mPxSdkWrap = NULL;
}		   

// 帧更新
void PxSceneWrap::Update(PxReal elapsedTime)
{
	CHECK_FALSE(InitOK(), , );

	if (mPxScene != NULL && elapsedTime > 0)
	{
		PxSceneWrapWriteLock lock(this);
		mPxScene->simulate(elapsedTime);
		mPxScene->fetchResults(true);
	}
}