#pragma once

#include "PxPhysicsAPI.h"			 
#include "physx_wrap.h"

class PxSdkWrap;
class PxSceneWrap;

// Sdk包装
class PxSdkWrap	: public PxWrap
{
public:
	PxSdkWrap();
	virtual ~PxSdkWrap();

	// 查询包装类型
	PxWrap::PxType GetPxType();

	// 初始化
	bool Init(bool usePvd, const char* host, int port, unsigned int timeoutInMs);

	// 初始化是否成功
	bool InitOK();

	// 释放
	void Release();	

	// 创建物理Scene
	PxSceneWrap* CreateScene(bool multiThread, void* bindGoObj);

public:
	physx::PxDefaultAllocator mPxAllocator;
	physx::PxDefaultErrorCallback mPxErrorCallback;
	physx::PxFoundation* mPxFoundation;
	physx::PxPhysics* mPxPhysics;
	physx::PxCooking* mPxCooking;
	physx::PxPvd* mPxPvd;		
};