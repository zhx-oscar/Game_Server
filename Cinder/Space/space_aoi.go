package Space

import (
	"Cinder/Base/Message"
)

func (space *Space) clientNotifyActorEnter(actor IActor) {

	msg := &Message.EnterAOI{}
	msg.Info.IsUser = false
	msg.Info.ID = actor.GetID()
	msg.Info.Type = actor.GetType()
	msg.Info.OwnerID = actor.GetOwnerUserID()

	if actor.GetProp() != nil {
		msg.Info.PropType = actor.GetPropType()
		msg.Info.Properties, _ = actor.GetProp().MarshalPart()
	}

	if msg.Info.Properties == nil {
		msg.Info.Properties = []byte{}
	}

	if actor.GetOwnerUserID() == "" {
		space.SendToAllClient(msg)
	} else {
		space.SendToAllClientExceptOne(msg, actor.GetOwnerUserID())

		if actor.GetProp() != nil {
			msg.Info.Properties, _ = actor.GetProp().Marshal()
		}

		if actor.GetOwnerUser() != nil {
			_ = actor.GetOwnerUser().SendToClient(msg)
		}
	}

}

func (space *Space) clientNotifyActorLeave(actor IActor) {
	msg := &Message.LeaveAOI{
		IsUser: false,
		ID:     actor.GetID(),
	}

	space.SendToAllClient(msg)
}
