package navmesh

/*
#cgo windows LDFLAGS: -L. -lgonavmesh
#cgo linux LDFLAGS: -L. -lgonavmesh -lm -lstdc++
#include "gonavmesh/gonavmesh.h"
#include <stdlib.h>
*/
import "C"
import (
	"Cinder/Base/linemath"
	"errors"
	"sync"
	"unsafe"
)

var staticMeshes sync.Map

var (
	ErrFindNearestStartPoly = errors.New("find nearest start poly failed")
	ErrFindNearestEndPoly   = errors.New("find nearest end poly failed")
	ErrFindPath             = errors.New("find path failed")
	ErrUnknown              = errors.New("unknonw error")
)

type Query struct {
	cgoQuery C.GoNavMeshQuery
}

// Load 加载静态navmesh数据
func Load(path string) C.GoNavMesh {
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))

	v, _ := staticMeshes.LoadOrStore(path, C.LoadMesh(cPath))

	if v != nil {
		return v.(C.GoNavMesh)
	}

	return nil
}

// Delete 卸载静态navmesh数据
func Delete(path string) {
	v, ok := staticMeshes.Load(path)
	if !ok {
		return
	}

	C.DeleteMesh(v.(C.GoNavMesh))
	staticMeshes.Delete(path)
}

// CreateQuery 创建寻路接口, 非线程安全, 不用时需要销毁防止内存泄漏
func CreateQuery(path string, maxNodes int) *Query {
	mesh := Load(path)
	query := C.CreateQuery(mesh, C.int(maxNodes))
	if query == nil {
		return nil
	}

	return &Query{cgoQuery: query}
}

// DestroyQuery 销毁寻路接口
func DestroyQuery(query *Query) {
	C.DeleteQuery(query.cgoQuery)
}

func (q *Query) FindPath(start, end linemath.Vector3) ([]linemath.Vector3, bool, error) {
	cStart := gv2cv(start)
	cEnd := gv2cv(end)
	path := make([]C.NavVector, 256)
	size := C.FindPath(q.cgoQuery, cStart, cEnd, (*C.NavVector)((unsafe.Pointer)(&path[0])))

	if size <= 0 {
		switch size {
		case 0:
			return nil, false, nil
		case -1:
			return nil, false, ErrFindNearestStartPoly
		case -2:
			return nil, false, ErrFindNearestEndPoly
		case -3:
			return nil, false, ErrFindPath
		default:
			return nil, false, ErrUnknown
		}
	}

	result := make([]linemath.Vector3, 0, size)
	for i := 0; i < int(size); i++ {
		result = append(result, cv2gv(path[i]))
	}

	return result, true, nil
}

func cv2gv(v C.NavVector) linemath.Vector3 {
	return linemath.Vector3{
		X: float32(v.x),
		Y: float32(v.y),
		Z: float32(v.z),
	}
}

func gv2cv(v linemath.Vector3) C.NavVector {
	pos := C.NavVector{}
	pos.x = C.float(v.X)
	pos.y = C.float(v.Y)
	pos.z = C.float(v.Z)
	return pos
}
