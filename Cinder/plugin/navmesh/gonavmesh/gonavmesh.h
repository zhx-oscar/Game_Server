#ifndef GONAVMESH_H
#define GONAVMESH_H

#ifndef GO_NAVIGATION
#ifdef _WIN32
#	define GO_NAVIGATION __declspec(dllexport)
#else
#	define GO_NAVIGATION
#endif
#endif

#ifdef __cplusplus
extern "C" {
#endif

static const int MAX_POLYS;

typedef struct {
    float x;
    float y;
    float z;
} NavVector;

typedef struct GoNavMeshT{} *GoNavMesh;
typedef struct GoNavMeshQueryT{} *GoNavMeshQuery;

GO_NAVIGATION GoNavMesh LoadMesh(const char* path);
GO_NAVIGATION void DeleteMesh(GoNavMesh mesh);

GO_NAVIGATION GoNavMeshQuery CreateQuery(GoNavMesh mesh, int maxNodes);
GO_NAVIGATION void DeleteQuery(GoNavMeshQuery query);

GO_NAVIGATION int FindPath(GoNavMeshQuery goQuery, NavVector start, NavVector end, NavVector* path);

#ifdef __cplusplus
}
#endif

#endif // GONAVMESH_H