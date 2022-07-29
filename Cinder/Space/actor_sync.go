package Space

import (
	"Cinder/Base/Message"
	"Cinder/Base/Prop"
)

func (actor *Actor) SyncProp(methodName string, args []byte, targets ...int) {

	as, _ := Message.UnPackArgs(args)
	actor.GetProp().GetCaller().CallMethod(methodName, as...)

	if !actor.IsReady() {
		return
	}

	msg := &Message.ActorPropNotify{
		SpaceID:    actor.GetSpace().GetID(),
		ActorID:    actor.GetID(),
		MethodName: methodName,
		Args:       args,
	}

	for _, target := range targets {

		switch target {
		case Prop.Target_Client:
			if ownerUser := actor.GetOwnerUser(); ownerUser != nil {
				_ = ownerUser.SendToClient(msg)
			} else {
				actor.Warn("actor no support target_client if no owner user")
			}
		case Prop.Target_All_Clients:
			actor.GetSpace().SendToAllClient(msg)
		case Prop.Target_Other_Clients:
			if ownerUserID := actor.GetOwnerUserID(); ownerUserID != "" {
				actor.GetSpace().SendToAllClientExceptOne(msg, ownerUserID)
			} else {
				actor.Warn("actor no support target_other_client is no owner user")
			}

		default:
			actor.Warn("no support target type", target)
		}
	}
}
