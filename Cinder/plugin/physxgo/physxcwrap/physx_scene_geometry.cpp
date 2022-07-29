#include "physx_scene.h"
#include "physx_sdk.h"
#include "common.h"

using namespace physx;

bool PxSceneWrap::BuildHeightFieldGeometry(PxHeightFieldGeometry& geom, const std::vector<int16_t>& heightmap, unsigned columns, unsigned rows, const PxVec3& scale)
{
	CHECK_FALSE(InitOK(), , false);

	unsigned hfNumVerts = columns * rows;
	PxHeightFieldSample* samples = (PxHeightFieldSample*)malloc(sizeof(PxHeightFieldSample) * hfNumVerts);
	CHECK_NULL(samples, , false);
	memset(samples, 0, hfNumVerts * sizeof(PxHeightFieldSample));

	for (unsigned row = 0; row < rows; row++)
	{
		for (unsigned col = 0; col < columns; col++)
		{
			int index = col + row * columns;
			samples[index].height = heightmap[index];
		}
	}

	PxHeightFieldDesc hfDesc;
	hfDesc.format = PxHeightFieldFormat::eS16_TM;
	hfDesc.nbColumns = columns;
	hfDesc.nbRows = rows;
	hfDesc.samples.data = samples;
	hfDesc.samples.stride = sizeof(PxHeightFieldSample);

	PxHeightField* heightField = mPxSdkWrap->mPxCooking->createHeightField(hfDesc, mPxScene->getPhysics().getPhysicsInsertionCallback());
	CHECK_NULL(heightField, free(samples), false);	

	geom.heightField = heightField;
	geom.columnScale = scale.x;
	geom.heightScale = scale.y;
	geom.rowScale = scale.z;

	free(samples);
	return true;
}

bool PxSceneWrap::BuildMeshGeometry(PxTriangleMeshGeometry& geom, const PxVec3& scale, const std::vector<float>& vb, const std::vector<uint16_t>& ib)
{
	CHECK_FALSE(InitOK(), , false);

	PxTriangleMeshDesc meshDesc;
	meshDesc.points.count = PxU32(vb.size() / 3);
	meshDesc.triangles.count = PxU32(ib.size() / 3);
	meshDesc.points.stride = sizeof(float) * 3;
	meshDesc.triangles.stride = sizeof(uint16_t) * 3;
	meshDesc.points.data = vb.data();
	meshDesc.triangles.data = ib.data();
	meshDesc.flags |= PxMeshFlag::e16_BIT_INDICES | PxMeshFlag::eFLIPNORMALS;

	PxDefaultMemoryOutputStream streamout;	
	if (!mPxSdkWrap->mPxCooking->cookTriangleMesh(meshDesc, streamout))
	{
		return false;
	}

	PxDefaultMemoryInputData streamin(streamout.getData(), streamout.getSize());
	PxTriangleMesh* triangleMesh = mPxScene->getPhysics().createTriangleMesh(streamin);
	CHECK_NULL(triangleMesh, , false);

	PxMeshScale meshScale = PxMeshScale(PxVec3{ scale.x ,scale.y ,scale.z }, PxQuat(PxIdentity));
	geom.triangleMesh = triangleMesh;
	geom.scale = meshScale;

	return true;
}
