package sceneconfigdata

import "testing"

func Test_LoadPath(t *testing.T) {
	cfg, err := LoadPath("../../../../res/MapData/1/path.json")
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("pathCfg:%+v", cfg)
}
