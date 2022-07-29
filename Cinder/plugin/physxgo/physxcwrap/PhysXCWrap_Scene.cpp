#include "PhysXCWrap.h"
#include "physx_scene.h"
#include "physx_sdk.h"
#include "common.h"

// 创建Scene
PHYSXCWRAP_API PxHandle PxSdkCreatePxScene(PxHandle sdk, bool multiThread, void* bindGoObj)
{
	auto pSdkWrap = ConvertPxWrap<PxSdkWrap>(sdk);
	CHECK_NULL(pSdkWrap, , NULL);

	return pSdkWrap->CreateScene(multiThread, bindGoObj);
}

// 销毁Scene
PHYSXCWRAP_API bool ReleasePxScene(PxHandle scene)
{
	auto pSceneWrap = ConvertPxWrap<PxSceneWrap>(scene);
	CHECK_NULL(pSceneWrap, , false);

	delete pSceneWrap;

	return true;
}

// Scene帧更新
PHYSXCWRAP_API void PxSceneUpdate(PxHandle scene, float elapsedTime)
{
	auto pSceneWrap = ConvertPxWrap<PxSceneWrap>(scene);
	CHECK_NULL(pSceneWrap, , );

	pSceneWrap->Update(elapsedTime);
}