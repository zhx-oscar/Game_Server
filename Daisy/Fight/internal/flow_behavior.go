package internal

// _BehaviorFlow 行为树流程
type _BehaviorFlow struct {
	scene *Scene
}

// init 初始化
func (flow *_BehaviorFlow) init(scene *Scene) {
	flow.scene = scene
}

// update 帧更新
func (flow *_BehaviorFlow) update() {
	for _, pawn := range flow.scene.pawnList {
		if !pawn.IsAlive() {
			continue
		}

		pawn.behaviorUpdate()
	}
}

// AllAIPause 所有AI暂停
func (flow *_BehaviorFlow) AllAIPause(pause bool) {
	for _, pawn := range flow.scene.pawnList {
		pawn.AIPause(pause)
	}
}

// AllAIPauseSkipOne 除指定pawn以外所有AI暂停
func (flow *_BehaviorFlow) AllAIPauseSkipOne(skipPawn *Pawn, pause bool) {
	for _, pawn := range flow.scene.pawnList {
		if pawn.Equal(skipPawn) {
			continue
		}
		pawn.AIPause(pause)
	}
}
