package DB

import (
	"Cinder/Base/Util"
	"testing"
)

func BenchmarkPropUtil(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {

			util, err := NewPropUtil(Util.GetGUID(), "TestProp")
			if err != nil {
				b.Fatal(err)
				return
			}

			if _, err = util.GetBsonData(); err != nil {
				b.Fatal(err)
				return
			}
		}
	})
}

func TestNewPropUtil(t *testing.T) {

	util, err := NewPropUtil("5ec1e4e17d85345288be99cf", "RoleProp")
	if err != nil {
		t.Fatal(err)
		return
	}

	data, err := util.GetBsonData()
	if err != nil {
		t.Fatal(err)
		return
	}

	t.Log(data)
}
