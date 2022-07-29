package Space

import (
	"Cinder/Base/Message"
	"Cinder/Base/Prop"
)

func (space *Space) SyncProp(methodName string, args []byte, targets ...int) {

	as, _ := Message.UnPackArgs(args)
	_, _ = space.GetProp().GetCaller().CallMethod(methodName, as...)

	msg := &Message.SpacePropNotify{
		SpaceID:    space.GetID(),
		MethodName: methodName,
		Args:       args,
	}

	for _, target := range targets {
		switch target {
		case Prop.Target_All_Clients:
			space.SendToAllClient(msg)
		default:
			space.Warn("space no support target type", target)
		}

	}
}
