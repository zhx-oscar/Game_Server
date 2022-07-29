package Space

func (space *Space) onAddActor(actor _IActor) {
	space.Info("Actor enter space", actor.GetType(), actor.GetID())
	space.clientNotifyActorEnter(actor)
}

func (space *Space) onRemoveActor(actor _IActor) {
	space.Info("Actor leave space ", actor.GetType(), actor.GetID())
	space.clientNotifyActorLeave(actor)
}
