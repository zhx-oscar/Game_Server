package Space

import (
	"Cinder/Base/Message"
	BaseUser "Cinder/Base/User"
)

func (user *User) clientNotifyUserEnter() {
	user.notifySelfEnterSpace()
	user.notifySelfBatchAOI()
	user.notifyOthersEnterAOI()
}

func (user *User) notifySelfEnterSpace() {
	msg := &Message.EnterSpace{}
	msg.SpaceID = user.GetSpace().GetID()

	if user.GetSpace().GetOwnerUser() != nil {
		msg.OwnerID = user.GetSpace().GetOwnerUser().GetID()
	} else {
		msg.OwnerID = ""
	}

	msg.SpaceInfo, _ = user.GetSpace().GetProp().MarshalPart()
	_ = user.SendToClient(msg)
	user.Infof("notifySelfEnterSpace spaceID:%s, ownerID:%s", msg.SpaceID, msg.OwnerID)
}

func (user *User) notifySelfBatchAOI() {
	bmsg := &Message.BatchEnterAOI{
		Info: make([]Message.EnterAOIInfo, 0, 10),
	}

	user.GetSpace().TraversalActor(func(actor IActor) {
		info := Message.EnterAOIInfo{}
		info.IsUser = false
		info.ID = actor.GetID()
		info.Type = actor.GetType()

		info.OwnerID = actor.GetOwnerUserID()

		if actor.GetProp() != nil {
			info.PropType = actor.GetPropType()
			if user.GetID() == actor.GetOwnerUserID() {
				info.Properties, _ = actor.GetProp().Marshal()
			} else {
				info.Properties, _ = actor.GetProp().MarshalPart()
			}

		}

		if info.Properties == nil {
			info.Properties = []byte{}
		}

		bmsg.Info = append(bmsg.Info, info)
	})

	user.GetSpace().TraversalUser(func(user BaseUser.IUser) bool {

		info := Message.EnterAOIInfo{}
		info.IsUser = true
		info.ID = user.GetID()
		info.PropType = user.GetPropType()
		if user.GetProp() != nil {
			info.Properties, _ = user.GetProp().MarshalPart()
		}

		bmsg.Info = append(bmsg.Info, info)

		return true
	})

	if len(bmsg.Info) > 0 {
		if err := user.SendToClient(bmsg); err != nil {
			user.Error("batch enter aoi err:", err)
		}
	}
	user.Debug("batch enter aoi")
	for i := 0; i < len(bmsg.Info); i++ {
		info := bmsg.Info[i]
		user.Debugf("enter aoi ID:%s, ownerID:%s, isUser:%v", info.ID, info.OwnerID, info.IsUser)
	}
}

func (user *User) notifyOthersEnterAOI() {

	msg := &Message.EnterAOI{}
	msg.Info = Message.EnterAOIInfo{}
	msg.Info.IsUser = true
	msg.Info.ID = user.GetID()
	if user.GetProp() != nil {
		msg.Info.PropType = user.GetPropType()
		msg.Info.Properties, _ = user.GetProp().MarshalPart()
	}

	user.GetSpace().SendToAllClientExceptOne(msg, user.GetID())
}

func (user *User) clientNotifyUserLeave() {
	user.notifyOthersLeaveAOI()
	user.notifySelfClearAOI()
	user.notifySelfLeaveSpace()
}

func (user *User) notifySelfLeaveSpace() {
	bmsg := &Message.LeaveSpace{}
	bmsg.SpaceID = user.GetSpace().GetID()
	_ = user.SendToClient(bmsg)
}

func (user *User) notifySelfClearAOI() {
	msg := &Message.ClearAOI{SpaceID: user.GetSpace().GetID()}
	_ = user.SendToClient(msg)
}

func (user *User) notifyOthersLeaveAOI() {
	msg := &Message.LeaveAOI{
		IsUser: true,
		ID:     user.GetID(),
	}

	user.GetSpace().SendToAllClientExceptOne(msg, user.GetID())
}
