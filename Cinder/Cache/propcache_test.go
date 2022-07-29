package Cache

import (
	"reflect"
	"testing"
)

func TestSetGetPropCache(t *testing.T) {
	propType := "TestProp"
	propID := "helloworld"
	propCache := []byte("Men always remember love because of romance only")

	if err := SetPropCache(propType, propID, propCache); err != nil {
		t.Fatal(err)
		return
	}

	if data, err := GetPropCache(propType, propID); err != nil {
		t.Fatal(err)
		return
	} else {
		if !reflect.DeepEqual(data, propCache) {
			t.Fatal("Mismatch", data, propCache)
			return
		}
	}
}

func TestSetGetPropCacheList(t *testing.T) {
	propTypes := []string{"TestProp1", "TestProp2"}
	propIDs := []string{"1", "2"}
	propCaches := [][]byte{[]byte("hello"), []byte("world")}

	if err := SetPropCacheList(propTypes, propIDs, propCaches); err != nil {
		t.Fatal(err)
		return
	}

	if data, err := GetPropCacheList(propTypes, propIDs); err != nil {
		t.Fatal(err)
		return
	} else {
		if !reflect.DeepEqual(data, propCaches) {
			t.Fatal("Mismatch", data, propCaches)
			return
		}
	}
}
