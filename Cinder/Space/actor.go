package Space

import (
	"Cinder/Base/Core"
	"Cinder/Base/Message"
	"Cinder/Base/Prop"
	"Cinder/Base/event"
	"fmt"
)

type _IActor interface {
	IActor
	InitBase(data *_InitInfo)
	SetReady()
	DestroyBase()
}

type Actor struct {
	id          string
	typ         string
	userData    interface{}
	mgr         IActorMgr
	ownerUserID string
	realPtr     interface{}
	debugStr    string
	isReady     bool

	Prop.IPropOwner
	space ISpace
	event.ILocalEventDispatcher
}

func (actor *Actor) InitBase(data *_InitInfo) {
	actor.id = data.ID
	actor.typ = data.Type
	actor.debugStr = fmt.Sprintf("[A:%s:%s]", actor.typ, actor.id)
	actor.mgr = data.Owner.(IActorMgr)
	actor.realPtr = data.RealPtr
	actor.space = data.OwnerRealPtr.(ISpace)
	actor.userData = data.UserData
	actor.ownerUserID = data.OwnerUserID

	actor.IPropOwner = Core.Inst.NewPropOwner(data.RealPtr)
	actor.IPropOwner.InitPropOwner(data.PropData)
	if p := actor.GetProp(); p != nil {
		p.GetCaller().SetParentCaller(actor.GetSpace().(_ISpace).GetSafeCall())
	}

	actor.ILocalEventDispatcher = event.GetLocalEventDispatcher()
}

func (actor *Actor) DestroyBase() {
	actor.IPropOwner.DestroyPropOwner()
}

func (actor *Actor) GetID() string {
	return actor.id
}

func (actor *Actor) GetType() string {
	return actor.typ
}

func (actor *Actor) GetUserData() interface{} {
	return actor.userData
}

func (actor *Actor) GetRealPtr() interface{} {
	return actor.realPtr
}

func (actor *Actor) GetSpace() ISpace {
	return actor.space
}

func (actor *Actor) GetOwnerUser() IUser {
	u, _ := actor.GetSpace().GetUser(actor.ownerUserID)
	return u
}

func (actor *Actor) DestroySelf() {
	_ = actor.mgr.RemoveActor(actor.GetID())
}

func (actor *Actor) LoopBase() {

}

func (actor *Actor) SetOwnerUserID(id string) {
	if ownerUser, _ := actor.GetSpace().GetUser(id); ownerUser != nil {
		actor.ownerUserID = id

		var propData []byte
		var propType string

		if actor.GetProp() != nil {
			propData, _ = actor.GetProp().Marshal()
			propType = actor.GetPropType()
		} else {
			propData = []byte{}
		}

		ownerUser.SendToClient(&Message.ActorRefreshOwnerUser{
			ActorID:     actor.GetID(),
			OwnerUserID: id,
			PropType:    propType,
			PropData:    propData,
		})

		actor.Info("SetOwnerUserID ID:", id)
	}
}

func (actor *Actor) GetOwnerUserID() string {
	return actor.ownerUserID
}

func (actor *Actor) IsReady() bool { return actor.isReady }
func (actor *Actor) SetReady()     { actor.isReady = true }
