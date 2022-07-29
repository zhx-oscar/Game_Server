package main

import (
	"Cinder/Base/linemath"
	"Cinder/plugin/physxgo"
	"Cinder/plugin/physxgo/internal"
	"fmt"
	"time"
)

func main() {
	if err := physxgo.InitPxSdk(true, "127.0.0.1", 5425); err != nil {
		panic(err)
	}

	pxScene, err := physxgo.CreatePxScene(true)
	if err != nil || pxScene == nil {
		panic(err)
	}

	trap, err := pxScene.AddBoxKinematic(physxgo.TransForm{P: linemath.Vector3{
		X: 0,
		Y: 0,
		Z: 0,
	}, Q: linemath.Quaternion{
		X: 0,
		Y: 0,
		Z: 0,
		W: 1,
	}}, linemath.Vector3{
		X: 10,
		Y: 10,
		Z: 10,
	}, physxgo.ActorMode_eTrap, physxgo.NoHitFilter, &physxgo.ActorEvents{
		OnActorEnterCallback: func(self, actor internal.IPxActor) {
			fmt.Println("actor in")
		},
		OnActorLeaveCallback: func(self, actor internal.IPxActor) {
			fmt.Println("actor out")
		},
	})
	if err != nil {
		panic(err)
	}
	defer trap.Release()

	actor, err := pxScene.AddCapsuleKinematic(physxgo.TransForm{P: linemath.Vector3{
		X: 0,
		Y: 0,
		Z: 0,
	}, Q: linemath.Quaternion{
		X: 0,
		Y: 0,
		Z: 0,
		W: 1,
	}}, 10, 10, physxgo.ActorMode_eNone, physxgo.NoHitFilter, &physxgo.ActorEvents{
		OnEnterTrapCallback: func(self, trap internal.IPxActor) {
			fmt.Println("in trap")
		},
		OnLeaveTrapCallback: func(self, trap internal.IPxActor) {
			fmt.Println("out trap")
		},
	})
	if err != nil {
		panic(err)
	}
	defer actor.Release()

	go func() {
		for {
			pxScene.Update(0.1)
			time.Sleep(100 * time.Millisecond)
		}
	}()

	for {
		var a string
		fmt.Scanln(&a)

		switch a {
		case "w":
			actor.Step(physxgo.TransForm{
				P: func() linemath.Vector3 {
					t := actor.GetPose().P
					t.Z += 5
					return t
				}(),
				Q: actor.GetPose().Q,
			})
		case "a":
			actor.Step(physxgo.TransForm{
				P: func() linemath.Vector3 {
					t := actor.GetPose().P
					t.X -= 5
					return t
				}(),
				Q: actor.GetPose().Q,
			})
		case "s":
			actor.Step(physxgo.TransForm{
				P: func() linemath.Vector3 {
					t := actor.GetPose().P
					t.Z -= 5
					return t
				}(),
				Q: actor.GetPose().Q,
			})
		case "d":
			actor.Step(physxgo.TransForm{
				P: func() linemath.Vector3 {
					t := actor.GetPose().P
					t.X += 5
					return t
				}(),
				Q: actor.GetPose().Q,
			})
		}
	}

	physxgo.ShutPxSdk()
}
