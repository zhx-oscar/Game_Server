package sceneconfigdata

import (
	"Cinder/Base/linemath"
	"github.com/json-iterator/go"
	"io/ioutil"
	"sync"
)

type PathCfg struct {
	LeaderSpeed   float32
	FirstStartIdx []uint32
	RestartIdx    []uint32
	Points        []linemath.Vector3
	Path          []PointRelation
}

type PointRelation struct {
	Name     string
	Idx      uint32
	Previous []uint32
	Nexts    []uint32
}

var pathCfg sync.Map
var pathCfgLock sync.RWMutex

func LoadPath(file string) (*PathCfg, error) {
	if v, ok := pathCfg.Load(file); ok {
		return v.(*PathCfg), nil
	}

	pathCfgLock.Lock()
	defer pathCfgLock.Unlock()

	if v, ok := pathCfg.Load(file); ok {
		return v.(*PathCfg), nil
	}

	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	cfg := &PathCfg{}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err = json.Unmarshal(data, cfg)
	if err != nil {
		return nil, err
	}

	//convert to z-up
	for i := 0; i < len(cfg.Points); i++ {
		cfg.Points[i] = ConvertVector3(cfg.Points[i])
	}

	pathCfg.Store(file, cfg)
	return cfg, nil
}
