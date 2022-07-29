package sceneconfigdata

import "testing"

func Test_LoadTrigger(t *testing.T) {
	cfg, err := LoadTrigger("../../../res/MapData/1/trigger.json")
	if err != nil {
		t.Error(err)
		return
	}
	//t.Logf("triggerCfg:%+v", cfg)
	if len(cfg) > 0 && cfg[0].TriggerType == TriggerType_Jump {
		t.Logf("%+v", cfg[0].TriggerParam.Positive.(*TriggerParamJump).OutId)
	}
}
