package sceneconfigdata

import (
	"Cinder/Base/linemath"
	log "github.com/cihub/seelog"
	"github.com/json-iterator/go"
	"io/ioutil"
	"sync"
)

const (
	TriggerType_None = iota
	TriggerType_Jump
	TriggerType_Speed
	TriggerType_TeamFork
	TriggerType_PersonalFork
	TriggerType_Battle
	TriggerType_Chest        = 9
	TriggerType_SpawnMonster = 13
)

type TriggerCfg struct {
	ID              string
	TriggerPxFilter uint32
	HalfExtends     linemath.Vector3
	TriggerType     uint32
	TriggerParam    TriggerParam
	Location        linemath.Vector3
	Rotation        linemath.Vector3
}

type TriggerParam struct {
	Positive    interface{}
	Opposite    interface{}
	JumpGravity float32
	JumpSpeed   float32
}

type TriggerParamJump struct {
	IsIn   bool
	OutId  string
	OutLoc linemath.Vector3
}

type TriggerParamSpeed struct {
	Speed float32
}

type TriggerParamTeamFork struct {
	Separate []uint32
}

type TriggerParamPersonalFork struct {
	Separate []uint32
}

type TriggerParamBattle struct {
	BattleID uint32
}

type TriggerParamChest struct {
	Probability uint32           `json:"possibility"`
	ChestLoc    linemath.Vector3 `json:"position"`
	ChestRot    linemath.Vector3 `json:"rotation"`
	DropID      uint32           `json:"dropID"`
}

type TriggerParamSpawnMonster struct {
	BattleID uint32
}

var triggerCfg sync.Map
var triggerCfgLock sync.RWMutex

func LoadTrigger(file string) ([]*TriggerCfg, error) {
	if v, ok := triggerCfg.Load(file); ok == true {
		return v.([]*TriggerCfg), nil
	}

	triggerCfgLock.Lock()
	defer triggerCfgLock.Unlock()

	if v, ok := triggerCfg.Load(file); ok == true {
		return v.([]*TriggerCfg), nil
	}

	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	cfg := make([]*TriggerCfg, 0)
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	for _, value := range cfg {
		if value.TriggerType == TriggerType_Jump {
			value.TriggerParam.Positive = &TriggerParamJump{}
			value.TriggerParam.Opposite = &TriggerParamJump{}
		} else if value.TriggerType == TriggerType_Speed {
			value.TriggerParam.Positive = &TriggerParamSpeed{}
			value.TriggerParam.Opposite = &TriggerParamSpeed{}
		} else if value.TriggerType == TriggerType_TeamFork {
			value.TriggerParam.Positive = &TriggerParamTeamFork{}
			value.TriggerParam.Opposite = &TriggerParamTeamFork{}
		} else if value.TriggerType == TriggerType_PersonalFork {
			value.TriggerParam.Positive = &TriggerParamPersonalFork{}
			value.TriggerParam.Opposite = &TriggerParamPersonalFork{}
		} else if value.TriggerType == TriggerType_Battle {
			value.TriggerParam.Positive = &TriggerParamBattle{}
			value.TriggerParam.Opposite = &TriggerParamBattle{}
		} else if value.TriggerType == TriggerType_Chest {
			value.TriggerParam.Positive = &TriggerParamChest{}
			value.TriggerParam.Opposite = &TriggerParamChest{}
		} else if value.TriggerType == TriggerType_SpawnMonster {
			value.TriggerParam.Positive = &TriggerParamSpawnMonster{}
			value.TriggerParam.Opposite = &TriggerParamSpawnMonster{}
		}
	}

	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	//convert to z-up
	for i := 0; i < len(cfg); i++ {
		cfg[i].Location = ConvertVector3(cfg[i].Location)
		cfg[i].Rotation = ConvertVector3(cfg[i].Rotation)
		cfg[i].HalfExtends = ConvertVector3(cfg[i].HalfExtends)
	}

	findTrigger := func(id string) *TriggerCfg {
		for _, value := range cfg {
			if value.ID == id {
				return value
			}
		}
		return nil
	}

	//填充OutLoc
	for _, value := range cfg {
		if value.TriggerType == TriggerType_Jump {
			positive := value.TriggerParam.Positive.(*TriggerParamJump)
			t := findTrigger(positive.OutId)
			if t != nil {
				positive.OutLoc = t.Location
			}

			if positive.IsIn && t == nil {
				log.Errorf("jump trigger %s 配置有误", value.ID)
			}

			opposite := value.TriggerParam.Opposite.(*TriggerParamJump)
			t = findTrigger(opposite.OutId)
			if t != nil {
				opposite.OutLoc = t.Location
			}

			if opposite.IsIn && t == nil {
				log.Errorf("jump trigger %s 配置有误", value.ID)
			}
		}
	}

	triggerCfg.Store(file, cfg)
	return cfg, nil
}
