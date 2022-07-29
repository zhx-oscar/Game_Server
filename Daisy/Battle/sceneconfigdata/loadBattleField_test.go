package sceneconfigdata

import "testing"

func Test_LoadBattleField(t *testing.T) {
	cfg, err := LoadBattleField("../../../res/MapData/1/Battle_1.json")
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("battleFieldCfg:%v", *cfg)
}
