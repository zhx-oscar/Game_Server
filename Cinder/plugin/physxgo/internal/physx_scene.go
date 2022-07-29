package internal

/*
#cgo CFLAGS: -I../physxcwrap/install/incl
#include "PhysXCWrap.h"
#include "stdlib.h"
*/
import "C"
import (
	"Cinder/Base/linemath"
	"errors"
	"sync"
	"unsafe"
)

// sceneMap 所有场景
var sceneMap sync.Map

// PxScene 场景
type PxScene struct {
	mPxScene C.PxHandle
}

// init 初始化scene
func (scene *PxScene) init(multiThread bool) bool {
	scene.mPxScene = C.PxSdkCreatePxScene(pxSdk, C.bool(multiThread), unsafe.Pointer(scene))
	if scene.mPxScene == nil {
		return false
	}

	return true
}

// Release 销毁Scene
func (scene *PxScene) Release() {
	if scene.mPxScene == nil {
		return
	}

	C.ReleasePxScene(scene.mPxScene)
	scene.mPxScene = nil
	sceneMap.Delete(scene)
}

// Update Scene帧更新
func (scene *PxScene) Update(elapsedTime float32) {
	if scene.mPxScene == nil {
		return
	}

	C.PxSceneUpdate(scene.mPxScene, C.float(elapsedTime))
}

// AddBoxKinematic 放置Kinematic Box
func (scene *PxScene) AddBoxKinematic(pose TransForm, halfExtents linemath.Vector3, actorModeFlags ActorModeFlags, hitFilter HitFilter,
	actorEvents *ActorEvents) (IPxActor, error) {
	if scene.mPxScene == nil {
		return nil, errors.New("nil PxScene")
	}

	cPose := pose.toC()
	cHalfExtents := Vector3GoToC(halfExtents)
	cHitFilter := hitFilter.toC()

	pxActor := &PxActor{}

	cPxActor := C.PxSceneCreateBoxKinematic(scene.mPxScene, cPose, cHalfExtents, C.ActorModeFlags(actorModeFlags), cHitFilter, unsafe.Pointer(pxActor))
	if cPxActor == nil {
		return nil, errors.New("create PxActor Kinematic Box failed")
	}

	if !pxActor.init(cPxActor, actorModeFlags, actorEvents) {
		C.ReleasePxActor(cPxActor)
		return nil, errors.New("init GoActor Kinematic Box failed")
	}

	if !C.PxSceneAddActor(scene.mPxScene, cPxActor) {
		C.ReleasePxActor(cPxActor)
		return nil, errors.New("PxScene add PxActor Kinematic Box failed")
	}

	actorMap.Store(pxActor, 0)
	pxActor.emitEnterLeaveTrapEvent()
	return pxActor, nil
}

// AddBoxStatic 放置Static Box
func (scene *PxScene) AddBoxStatic(pose TransForm, halfExtents linemath.Vector3, actorModeFlags ActorModeFlags, hitFilter HitFilter,
	actorEvents *ActorEvents) (IPxActor, error) {
	if scene.mPxScene == nil {
		return nil, errors.New("nil PxScene")
	}

	cPose := pose.toC()
	cHalfExtents := Vector3GoToC(halfExtents)
	cHitFilter := hitFilter.toC()

	pxActor := &PxActor{}

	cPxActor := C.PxSceneCreateBoxStatic(scene.mPxScene, cPose, cHalfExtents, C.ActorModeFlags(actorModeFlags), cHitFilter, unsafe.Pointer(pxActor))
	if cPxActor == nil {
		return nil, errors.New("create PxActor Static Box failed")
	}

	if !pxActor.init(cPxActor, actorModeFlags, actorEvents) {
		C.ReleasePxActor(cPxActor)
		return nil, errors.New("init GoActor Static Box failed")
	}

	if !C.PxSceneAddActor(scene.mPxScene, cPxActor) {
		C.ReleasePxActor(cPxActor)
		return nil, errors.New("PxScene add PxActor Static Box failed")
	}

	actorMap.Store(pxActor, 0)
	pxActor.emitEnterLeaveTrapEvent()
	return pxActor, nil
}

// AddSphereKinematic 放置Kinematic Sphere
func (scene *PxScene) AddSphereKinematic(pose TransForm, radius float32, actorModeFlags ActorModeFlags, hitFilter HitFilter,
	actorEvents *ActorEvents) (IPxActor, error) {
	if scene.mPxScene == nil {
		return nil, errors.New("nil PxScene")
	}

	cPose := pose.toC()
	cHitFilter := hitFilter.toC()

	pxActor := &PxActor{}

	cPxActor := C.PxSceneCreateSphereKinematic(scene.mPxScene, cPose, C.float(radius), C.ActorModeFlags(actorModeFlags), cHitFilter, unsafe.Pointer(pxActor))
	if cPxActor == nil {
		return nil, errors.New("create PxActor Kinematic Sphere failed")
	}

	if !pxActor.init(cPxActor, actorModeFlags, actorEvents) {
		C.ReleasePxActor(cPxActor)
		return nil, errors.New("init GoActor Kinematic Sphere failed")
	}

	if !C.PxSceneAddActor(scene.mPxScene, cPxActor) {
		C.ReleasePxActor(cPxActor)
		return nil, errors.New("PxScene add PxActor Kinematic Sphere failed")
	}

	actorMap.Store(pxActor, 0)
	pxActor.emitEnterLeaveTrapEvent()
	return pxActor, nil
}

// AddSphereStatic 放置Static Sphere
func (scene *PxScene) AddSphereStatic(pose TransForm, radius float32, actorModeFlags ActorModeFlags, hitFilter HitFilter,
	actorEvents *ActorEvents) (IPxActor, error) {
	if scene.mPxScene == nil {
		return nil, errors.New("nil PxScene")
	}

	cPose := pose.toC()
	cHitFilter := hitFilter.toC()

	pxActor := &PxActor{}

	cPxActor := C.PxSceneCreateSphereStatic(scene.mPxScene, cPose, C.float(radius), C.ActorModeFlags(actorModeFlags), cHitFilter, unsafe.Pointer(pxActor))
	if cPxActor == nil {
		return nil, errors.New("create PxActor Static Sphere failed")
	}

	if !pxActor.init(cPxActor, actorModeFlags, actorEvents) {
		C.ReleasePxActor(cPxActor)
		return nil, errors.New("init GoActor Static Sphere failed")
	}

	if !C.PxSceneAddActor(scene.mPxScene, cPxActor) {
		C.ReleasePxActor(cPxActor)
		return nil, errors.New("PxScene add PxActor Static Sphere failed")
	}

	actorMap.Store(pxActor, 0)
	pxActor.emitEnterLeaveTrapEvent()
	return pxActor, nil
}

// AddCapsuleKinematic 放置Kinematic Capsule
func (scene *PxScene) AddCapsuleKinematic(pose TransForm, radius float32, halfHeight float32, actorModeFlags ActorModeFlags, hitFilter HitFilter,
	actorEvents *ActorEvents) (IPxActor, error) {
	if scene.mPxScene == nil {
		return nil, errors.New("nil PxScene")
	}

	cPose := pose.toC()
	cHitFilter := hitFilter.toC()

	pxActor := &PxActor{}

	cPxActor := C.PxSceneCreateCapsuleKinematic(scene.mPxScene, cPose, C.float(radius), C.float(halfHeight), C.ActorModeFlags(actorModeFlags), cHitFilter, unsafe.Pointer(pxActor))
	if cPxActor == nil {
		return nil, errors.New("create PxActor Kinematic Capsule failed")
	}

	if !pxActor.init(cPxActor, actorModeFlags, actorEvents) {
		C.ReleasePxActor(cPxActor)
		return nil, errors.New("init GoActor Kinematic Capsule failed")
	}

	if !C.PxSceneAddActor(scene.mPxScene, cPxActor) {
		C.ReleasePxActor(cPxActor)
		return nil, errors.New("PxScene add PxActor Kinematic Capsule failed")
	}

	actorMap.Store(pxActor, 0)
	pxActor.emitEnterLeaveTrapEvent()
	return pxActor, nil
}

// AddCapsuleStatic 放置Static Capsule
func (scene *PxScene) AddCapsuleStatic(pose TransForm, radius float32, halfHeight float32, actorModeFlags ActorModeFlags, hitFilter HitFilter,
	actorEvents *ActorEvents) (IPxActor, error) {
	if scene.mPxScene == nil {
		return nil, errors.New("nil PxScene")
	}

	cPose := pose.toC()
	cHitFilter := hitFilter.toC()

	pxActor := &PxActor{}

	cPxActor := C.PxSceneCreateCapsuleStatic(scene.mPxScene, cPose, C.float(radius), C.float(halfHeight), C.ActorModeFlags(actorModeFlags), cHitFilter, unsafe.Pointer(pxActor))
	if cPxActor == nil {
		return nil, errors.New("create PxActor Static Capsule failed")
	}

	if !pxActor.init(cPxActor, actorModeFlags, actorEvents) {
		C.ReleasePxActor(cPxActor)
		return nil, errors.New("init GoActor Static Capsule failed")
	}

	if !C.PxSceneAddActor(scene.mPxScene, cPxActor) {
		C.ReleasePxActor(cPxActor)
		return nil, errors.New("PxScene add PxActor Static Capsule failed")
	}

	actorMap.Store(pxActor, 0)
	pxActor.emitEnterLeaveTrapEvent()
	return pxActor, nil
}

// RaycastOne 射线检测单个
func (scene *PxScene) RaycastOne(origin linemath.Vector3, unitDir linemath.Vector3, distance float32,
	hitModeFlags HitModeFlags, hitFilter HitFilter) (*Hit, error) {
	if scene.mPxScene == nil {
		return nil, errors.New("nil PxScene")
	}

	cOrigin := Vector3GoToC(origin)
	cUnitDir := Vector3GoToC(unitDir)
	cHitFilter := hitFilter.toC()

	var cHit C.Hit

	rv := C.PxSceneRaycast(scene.mPxScene, cOrigin, cUnitDir, C.float(distance), C.HitModeFlags(hitModeFlags), cHitFilter, &cHit, 1)
	if rv < 0 {
		return nil, errors.New("raycast one failed")
	} else if rv == 0 {
		return nil, nil
	}

	var hit Hit
	hit.fromC(&cHit)

	return &hit, nil
}

// RaycastMany 射线检测多个
func (scene *PxScene) RaycastMany(origin linemath.Vector3, unitDir linemath.Vector3, distance float32,
	hitModeFlags HitModeFlags, hitFilter HitFilter, maxHit int32) ([]Hit, error) {
	if scene.mPxScene == nil {
		return nil, errors.New("nil PxScene")
	}

	if maxHit <= 0 {
		return nil, errors.New("raycast many maxHit invalid")
	}

	cOrigin := Vector3GoToC(origin)
	cUnitDir := Vector3GoToC(unitDir)
	cHitFilter := hitFilter.toC()

	cHits := (*C.Hit)(C.malloc(C.sizeof_Hit * C.size_t(maxHit)))
	defer C.free(unsafe.Pointer(cHits))

	rv := C.PxSceneRaycast(scene.mPxScene, cOrigin, cUnitDir, C.float(distance), C.HitModeFlags(hitModeFlags), cHitFilter, cHits, C.size_t(maxHit))
	if rv < 0 {
		return nil, errors.New("raycast many failed")
	} else if rv == 0 {
		return nil, nil
	}

	hits := make([]Hit, rv)
	for i := 0; i < int(rv); i++ {
		hits[i].fromC((*C.Hit)(unsafe.Pointer(uintptr(unsafe.Pointer(cHits)) + uintptr(C.sizeof_Hit*C.int(i)))))
	}

	return hits, nil
}

// SweepOne 滑动检测单个
func (scene *PxScene) SweepOne(geom Geometry, pose TransForm, unitDir linemath.Vector3, distance float32, inflation float32,
	hitModeFlags HitModeFlags, hitFilter HitFilter) (*Hit, error) {
	if scene.mPxScene == nil {
		return nil, errors.New("nil PxScene")
	}

	cGeom := geom.toC()
	cPose := pose.toC()
	cUnitDir := Vector3GoToC(unitDir)
	cHitFilter := hitFilter.toC()

	var cHit C.Hit

	rv := C.PxSceneSweep(scene.mPxScene, cGeom, cPose, cUnitDir, C.float(distance), C.float(inflation),
		C.HitModeFlags(hitModeFlags), cHitFilter, &cHit, 1)
	if rv < 0 {
		return nil, errors.New("sweep one failed")
	} else if rv == 0 {
		return nil, nil
	}

	var hit Hit
	hit.fromC(&cHit)

	return &hit, nil
}

// SweepMany 滑动检测多个
func (scene *PxScene) SweepMany(geom Geometry, pose TransForm, unitDir linemath.Vector3, distance float32, inflation float32,
	hitModeFlags HitModeFlags, hitFilter HitFilter, maxHit int32) ([]Hit, error) {
	if scene.mPxScene == nil {
		return nil, errors.New("nil PxScene")
	}

	if maxHit <= 0 {
		return nil, errors.New("sweep many maxHit invalid")
	}

	cGeom := geom.toC()
	cPose := pose.toC()
	cUnitDir := Vector3GoToC(unitDir)
	cHitFilter := hitFilter.toC()

	cHits := (*C.Hit)(C.malloc(C.sizeof_Hit * C.size_t(maxHit)))
	defer C.free(unsafe.Pointer(cHits))

	rv := C.PxSceneSweep(scene.mPxScene, cGeom, cPose, cUnitDir, C.float(distance), C.float(inflation), C.HitModeFlags(hitModeFlags),
		cHitFilter, cHits, C.size_t(maxHit))
	if rv < 0 {
		return nil, errors.New("sweep many failed")
	} else if rv == 0 {
		return nil, nil
	}

	hits := make([]Hit, rv)
	for i := 0; i < int(rv); i++ {
		hits[i].fromC((*C.Hit)(unsafe.Pointer(uintptr(unsafe.Pointer(cHits)) + uintptr(C.sizeof_Hit*C.int(i)))))
	}

	return hits, nil
}

// OverlapOne 重叠检测单个
func (scene *PxScene) OverlapOne(geom Geometry, pose TransForm, hitModeFlags HitModeFlags, hitFilter HitFilter) (*Hit, error) {
	if scene.mPxScene == nil {
		return nil, errors.New("nil PxScene")
	}

	cGeom := geom.toC()
	cPose := pose.toC()
	cHitFilter := hitFilter.toC()

	var cHit C.Hit

	rv := C.PxSceneOverlap(scene.mPxScene, cGeom, cPose, C.HitModeFlags(hitModeFlags), cHitFilter, &cHit, 1)
	if rv < 0 {
		return nil, errors.New("overlap one failed")
	} else if rv == 0 {
		return nil, nil
	}

	var hit Hit
	hit.fromC(&cHit)

	return &hit, nil
}

// OverlapMany 重叠检测多个
func (scene *PxScene) OverlapMany(geom Geometry, pose TransForm, hitModeFlags HitModeFlags, hitFilter HitFilter, maxHit int32) ([]Hit, error) {
	if scene.mPxScene == nil {
		return nil, errors.New("nil PxScene")
	}

	if maxHit <= 0 {
		return nil, errors.New("overlap many maxHit invalid")
	}

	cGeom := geom.toC()
	cPose := pose.toC()
	cHitFilter := hitFilter.toC()

	cHits := (*C.Hit)(C.malloc(C.sizeof_Hit * C.size_t(maxHit)))
	defer C.free(unsafe.Pointer(cHits))

	rv := C.PxSceneOverlap(scene.mPxScene, cGeom, cPose, C.HitModeFlags(hitModeFlags), cHitFilter, cHits, C.size_t(maxHit))
	if rv < 0 {
		return nil, errors.New("overlap many failed")
	} else if rv == 0 {
		return nil, nil
	}

	hits := make([]Hit, rv)
	for i := 0; i < int(rv); i++ {
		hits[i].fromC((*C.Hit)(unsafe.Pointer(uintptr(unsafe.Pointer(cHits)) + uintptr(C.sizeof_Hit*C.int(i)))))
	}

	return hits, nil
}
