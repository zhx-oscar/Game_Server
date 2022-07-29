package navmesh

import (
	"Cinder/Base/linemath"
	"testing"
)

func TestQuery_FindPath(t *testing.T) {
	query := CreateQuery("all_tiles_navmesh.bin", 2048)
	if query == nil {
		t.Fatal("query failed")
		return
	}

	start := linemath.Vector3{
		X: 2.256470,
		Y: 9.998184,
		Z: -5.886890,
	}

	end := linemath.Vector3{
		X: 40.594093,
		Y: 9.998184,
		Z: 4.742132,
	}

	path, finded, err := query.FindPath(start, end)
	if err != nil {
		t.Fatal(err)
	}

	if finded {
		t.Log(path)
	}
}
