#include <string.h>
#include <stdio.h>
#include "DetourNavMesh.h"
#include "DetourNavMeshQuery.h"
#include "DetourCommon.h"

#ifndef GO_NAVIGATION
#ifdef _WIN32
#define GO_NAVIGATION __declspec(dllexport)
#else
#define GO_NAVIGATION
#endif
#endif

extern "C" {

    static const int NAVMESHSET_MAGIC = 'M'<<24 | 'S'<<16 | 'E'<<8 | 'T'; //'MSET';
    static const int NAVMESHSET_VERSION = 1;
    static const int MAX_POLYS = 256;

    struct NavMeshSetHeader
    {
    	int magic;
    	int version;
    	int numTiles;
    	dtNavMeshParams params;
    };

    struct NavMeshTileHeader
    {
    	dtTileRef tileRef;
    	int dataSize;
    };

    typedef struct {
        float x;
        float y;
        float z;
    } NavVector;

    typedef struct{} *GoNavMesh;
    typedef struct{} *GoNavMeshQuery;

    GO_NAVIGATION GoNavMesh LoadMesh(const char* path)
    {
        FILE* fp = fopen(path, "rb");
        if (!fp) return 0;

        // Read header.
        NavMeshSetHeader header;
        size_t readLen = fread(&header, sizeof(NavMeshSetHeader), 1, fp);
        if (readLen != 1)
        {
            fclose(fp);
            return 0;
        }
        if (header.magic != NAVMESHSET_MAGIC)
        {
            fclose(fp);
            return 0;
        }
        if (header.version != NAVMESHSET_VERSION)
        {
            fclose(fp);
            return 0;
        }

        dtNavMesh* mesh = dtAllocNavMesh();
        if (!mesh)
        {
            fclose(fp);
            return 0;
        }
        dtStatus status = mesh->init(&header.params);
        if (dtStatusFailed(status))
        {
            fclose(fp);
            return 0;
        }

        // Read tiles.
        for (int i = 0; i < header.numTiles; ++i)
        {
            NavMeshTileHeader tileHeader;
            readLen = fread(&tileHeader, sizeof(tileHeader), 1, fp);
            if (readLen != 1)
            {
                fclose(fp);
                return 0;
            }

            if (!tileHeader.tileRef || !tileHeader.dataSize)
                break;

            unsigned char* data = (unsigned char*)dtAlloc(tileHeader.dataSize, DT_ALLOC_PERM);
            if (!data) break;
            memset(data, 0, tileHeader.dataSize);
            readLen = fread(data, tileHeader.dataSize, 1, fp);
            if (readLen != 1)
            {
                dtFree(data);
                fclose(fp);
                return 0;
            }

            mesh->addTile(data, tileHeader.dataSize, DT_TILE_FREE_DATA, tileHeader.tileRef, 0);
        }

        fclose(fp);

        return reinterpret_cast<GoNavMesh>(mesh);
    }

    GO_NAVIGATION void DeleteMesh(GoNavMesh mesh)
    {
        dtNavMesh* dtMesh = reinterpret_cast<dtNavMesh*>(mesh);
        dtFreeNavMesh(dtMesh);
    }

    GO_NAVIGATION GoNavMeshQuery CreateQuery(GoNavMesh mesh, int maxNodes)
    {
        dtNavMeshQuery* query = dtAllocNavMeshQuery();
        dtStatus status = query->init(reinterpret_cast<dtNavMesh*>(mesh), maxNodes);
        if (dtStatusFailed(status))
        {
            printf("buildTiledNavigation: Could not init Detour navmesh query");
        	return 0;
        }

        return reinterpret_cast<GoNavMeshQuery>(query);
    }

    GO_NAVIGATION void DeleteQuery(GoNavMeshQuery query)
    {
        dtNavMeshQuery* dtQuery = reinterpret_cast<dtNavMeshQuery*>(query);
        dtFreeNavMeshQuery(dtQuery);
    }

    GO_NAVIGATION int FindPath(GoNavMeshQuery goQuery, NavVector start, NavVector end, NavVector* path)
    {
        dtNavMeshQuery* query = reinterpret_cast<dtNavMeshQuery*>(goQuery);

        dtQueryFilter filter;
        filter.setIncludeFlags(0xffff ^ 0x10);
        filter.setExcludeFlags(0);

        float polyPickExt[3];
        polyPickExt[0] = 2;
        polyPickExt[1] = 4;
        polyPickExt[2] = 2;

        float m_spos[3];
        float m_epos[3];
        m_spos[0] = start.x;
        m_spos[1] = start.y;
        m_spos[2] = start.z;
        m_epos[0] = end.x;
        m_epos[1] = end.y;
        m_epos[2] = end.z;

        dtPolyRef startRef;
        dtPolyRef endRef;
        dtStatus status;

        status = query->findNearestPoly(m_spos, polyPickExt, &filter, &startRef, 0);
        if (dtStatusFailed(status))
        {
            return -1;
        }

        status = query->findNearestPoly(m_epos, polyPickExt, &filter, &endRef, 0);
        if (dtStatusFailed(status))
        {
            return -2;
        }

        dtPolyRef polys[MAX_POLYS];
        int npolys;
        float straightPath[MAX_POLYS*3];
        unsigned char straightPathFlags[MAX_POLYS];
        dtPolyRef straightPathPolys[MAX_POLYS];
        int nstraightPath;
        status = query->findPath(startRef, endRef, m_spos, m_epos, &filter, polys, &npolys, MAX_POLYS);
        if (dtStatusFailed(status))
        {
            return -3;
        }

        if (npolys)
        {
            // In case of partial path, make sure the end point is clamped to the last polygon.
			float epos[3];
			dtVcopy(epos, m_epos);
			if (polys[npolys-1] != endRef)
				query->closestPointOnPoly(polys[npolys-1], m_epos, epos, 0);

			query->findStraightPath(m_spos, epos, polys, npolys, straightPath, straightPathFlags,
				straightPathPolys, &nstraightPath, MAX_POLYS, 0);

            for (int i = 0; i < nstraightPath; i++)
            {
                path[i].x = straightPath[i*3];
                path[i].y = straightPath[i*3 + 1];
                path[i].z = straightPath[i*3 + 2];
            }

            return nstraightPath;
        }

        return 0;
    }
}
