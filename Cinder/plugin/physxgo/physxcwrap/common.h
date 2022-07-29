#pragma once

#include "physx_wrap.h"

#define CHECK_NULL(p, code, rt) \
if (NULL == (p)) \
{ \
	code; \
	return rt; \
}

#define CHECK_TRUE(cond, code, rt) \
if ((cond)) \
{ \
	code; \
	return rt; \
}

#define CHECK_FALSE(cond, code, rt) \
if (!(cond)) \
{ \
	code; \
	return rt; \
}

#define PX_RELEASE(p)	if(p)	{ p->release(); p = NULL;	}

template<class T>
T* ConvertPxWrap(void* pxHandle)
{
	if (NULL == pxHandle)
	{
		return NULL;
	}

	static T model;

	auto pWrap = (PxWrap*)pxHandle;
	if (pWrap->GetPxType() != model.GetPxType())
	{
		return NULL;
	}

	return (T*)pWrap;
}