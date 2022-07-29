package internal

// _MovementFlow 移动控制器流程
type _MovementFlow struct {
	scene *Scene
}

// init 初始化
func (flow *_MovementFlow) init(scene *Scene) {
	flow.scene = scene
}

// update 帧更新
func (flow *_MovementFlow) update() {
	for _, pawn := range flow.scene.pawnList {
		if !pawn.IsAlive() {
			continue
		}

		pawn.controllerUpdate()
	}
}
