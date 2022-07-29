package main

import (
	"Cinder/Base/linemath"
	"Cinder/plugin/physxgo"
	"Daisy/Battle/sceneconfigdata"
	"fmt"
	log "github.com/cihub/seelog"
	"math"
	"sync"
)

var scenes sync.Map
var loadLock sync.Mutex

func init() {
	err := physxgo.InitPxSdk(false, "127.0.0.1", 5425)
	if err != nil {
		log.Error(err)
		return
	}
}

// 转换物理层角度
func ToPxAngle(v float32) float32 {
	angle := math.Mod(float64(v), 360)
	if angle > 180 {
		angle = angle - 360
	} else if angle < -180 {
		angle = angle + 360
	}
	angle = angle / 360 * linemath.PI2

	return float32(angle)
}

func LoadPxScene(mapID string) (physxgo.IPxScene, error) {
	scene, ok := scenes.Load(mapID)
	if ok {
		return scene.(physxgo.IPxScene), nil
	}

	loadLock.Lock()
	defer loadLock.Unlock()

	scene, ok = scenes.Load(mapID)
	if ok {
		return scene.(physxgo.IPxScene), nil
	}

	scene, err := physxgo.CreatePxScene(true)
	if err != nil {
		return nil, err
	}

	triggers, err := sceneconfigdata.LoadTrigger(fmt.Sprintf("../res/MapData/%s/trigger.json", mapID))
	if err != nil {
		return nil, err
	}

	for _, trigger := range triggers {
		rotation := linemath.Vector3{X: ToPxAngle(trigger.Rotation.X), Y: ToPxAngle(trigger.Rotation.Y), Z: ToPxAngle(trigger.Rotation.Z)}
		actor, err := scene.(physxgo.IPxScene).AddBoxStatic(physxgo.TransForm{P: trigger.Location, Q: linemath.EulerToQuaternion(rotation)},
			trigger.HalfExtends,
			physxgo.ActorMode_eNone,
			physxgo.NoHitFilter,
			nil)

		if err != nil {
			return nil, err
		}

		actor.SetUserData(trigger)
	}

	scenes.Store(mapID, scene)

	return scene.(physxgo.IPxScene), nil
}
