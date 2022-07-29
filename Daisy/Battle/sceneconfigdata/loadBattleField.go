package sceneconfigdata

import (
	"Cinder/Base/linemath"
	"github.com/json-iterator/go"
	"io/ioutil"
	"sync"
)

type BornPoint struct {
	linemath.Vector2
	Angle float32
}

type BattleFieldCfg struct {
	AreaPoints   []linemath.Vector2
	EnemyPoints  []BornPoint
	PlayerPoints []BornPoint
}

var battleFieldCfgs sync.Map
var battleFieldLock sync.Mutex

func LoadBattleField(file string) (*BattleFieldCfg, error) {
	if v, ok := battleFieldCfgs.Load(file); ok {
		return v.(*BattleFieldCfg), nil
	}

	battleFieldLock.Lock()
	defer battleFieldLock.Unlock()

	if v, ok := battleFieldCfgs.Load(file); ok {
		return v.(*BattleFieldCfg), nil
	}

	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	cfg := &BattleFieldCfg{}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err = json.Unmarshal(data, cfg)
	if err != nil {
		return nil, err
	}

	battleFieldCfgs.Store(file, cfg)
	return cfg, nil
}
