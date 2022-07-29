#pragma once
#pragma warning(disable:26812) 

class PxWrap {
public:
	// 包装类型定义
	enum PxType {
		PxSdk,    // 物理Sdk
		PxScene,  // 物理Scene
		PxActor   // 物理Actor
	};
	
	virtual ~PxWrap() = 0;
	
	// 查询包装类型
	virtual PxType GetPxType() = 0;	
};

inline PxWrap::~PxWrap() {}

