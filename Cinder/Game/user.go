package Game

import (
	"Cinder/Base/Const"
	"Cinder/Base/Core"
	"Cinder/Base/Message"
	"Cinder/Base/ProtoDef"
	BaseUser "Cinder/Base/User"
	"sync"
)

type User struct {
	BaseUser.User
	matchSrvID string

	propObjectMap sync.Map
}

func (user *User) LateInitBase() {
	var propData []byte
	if user != nil {
		propData, _ = user.GetProp().Marshal()
	}

	user.SendToPeerServer(Const.Agent, &Message.UserLoginRet{
		UserID:   user.GetID(),
		PropType: user.GetPropType(),
		PropDef:  ProtoDef.GetProtoDefData(),
		UserData: propData,
	})
}

func (user *User) DestroyBase() {
	user.clearPropObjectMap()
}

func (user *User) clearPropObjectMap() {
	user.propObjectMap.Range(func(id, srvID interface{}) bool {
		if err := Core.Inst.Send(srvID.(string), &Message.PropObjectCloseReq{
			ID:     id.(string),
			UserID: user.GetID(),
		}); err != nil {
			user.Error("clearPropObjectMap Send Msg err", err)
		}

		return true
	})
}
