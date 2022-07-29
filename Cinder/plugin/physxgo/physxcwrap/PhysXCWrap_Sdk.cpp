#include "PhysXCWrap.h"
#include "physx_sdk.h"
#include "common.h"

using namespace physx;

// 创建Sdk
PHYSXCWRAP_API PxHandle CreatePxSdk()
{
	auto pSdkWrap = new PxSdkWrap();
	CHECK_NULL(pSdkWrap, , NULL);

	CHECK_FALSE(pSdkWrap->Init(false, "", 0, 0), { delete pSdkWrap; }, NULL);

	return pSdkWrap;
}

// 创建Sdk并连接pvd工具
PHYSXCWRAP_API PxHandle CreatePxSdkConnectPvd(const char* pvdHost, int32_t pvdPort, uint32_t timeoutInMs)
{
	auto pSdkWrap = new PxSdkWrap();
	CHECK_NULL(pSdkWrap, , NULL);

	CHECK_FALSE(pSdkWrap->Init(true, pvdHost, pvdPort, timeoutInMs), { delete pSdkWrap; }, NULL);

	return pSdkWrap;
}

// 销毁Sdk
PHYSXCWRAP_API bool ReleasePxSdk(PxHandle sdk)
{
	auto pSdkWrap = ConvertPxWrap<PxSdkWrap>(sdk);
	CHECK_NULL(pSdkWrap, , false);

	delete pSdkWrap;

	return true;
}

TransForm ZeroTransForm()
{
	TransForm rv;
	rv.P = Vector3{ 0, 0, 0 };
	rv.Q = Quat{ 0, 0, 0, 0 };
	rv.Q.W = 1;
	return rv;
}