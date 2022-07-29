package internal

/*
#cgo CFLAGS: -I../physxcwrap/install/incl
#cgo windows LDFLAGS: -L../physxcwrap/install/lib/release -lPhysXCWrap
#cgo linux LDFLAGS: -L../physxcwrap/install/lib/release -L../PhysX-4.1/install/linux/PhysX/bin/linux.clang/release -Wl,--start-group -lPhysXCWrap -lPhysXCharacterKinematic_static_64 -lPhysXCommon_static_64 -lPhysXCooking_static_64 -lPhysXExtensions_static_64 -lPhysXFoundation_static_64 -lPhysXPvdSDK_static_64 -lPhysX_static_64 -lPhysXVehicle_static_64 -lm -lstdc++ -ldl -Wl,--end-group
#include "PhysXCWrap.h"
#include "stdlib.h"
*/
import "C"
import (
	"errors"
	"unsafe"
)

var pxSdk C.PxHandle = nil

// InitPxSdk 初始化PxSdk
func InitPxSdk(usePvd bool, pvdHost string, port int32) error {
	if usePvd {
		cPvdHost := C.CString(pvdHost)
		defer C.free(unsafe.Pointer(cPvdHost))

		cPvdPort := C.int32_t(port)

		pxSdk = C.CreatePxSdkConnectPvd(cPvdHost, cPvdPort, 10000)

	} else {
		pxSdk = C.CreatePxSdk()
	}

	if pxSdk == nil {
		return errors.New("init PxSdk failed")
	}

	return nil
}

// ShutPxSdk 销毁PxSdk
func ShutPxSdk() {
	if pxSdk != nil {
		C.ReleasePxSdk(pxSdk)
		pxSdk = nil
	}
}

// CreatePxScene 创建PxScene
func CreatePxScene(multiThread bool) (IPxScene, error) {
	pxScene := &PxScene{}
	if !pxScene.init(multiThread) {
		return nil, errors.New("init PxScene failed")
	}
	sceneMap.Store(pxScene, 0)
	return pxScene, nil
}
