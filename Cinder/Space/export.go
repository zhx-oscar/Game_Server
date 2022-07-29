package Space

import (
	"Cinder/Base/Core"
	"Cinder/Base/Message"
	"Cinder/Base/Prop"
	BaseUser "Cinder/Base/User"
	"Cinder/Base/Util"
	"time"
)

var Inst ISpaceServer

type IActor interface {
	Prop.IPropOwner
	GetID() string
	GetType() string
	GetUserData() interface{}
	GetRealPtr() interface{}
	GetSpace() ISpace
	SetOwnerUserID(id string)
	GetOwnerUserID() string
	GetOwnerUser() IUser

	DestroySelf()
}

type IUser interface {
	BaseUser.IUser
	GetSpace() ISpace

	SendToAllClient(msg Message.IMessage)
	SendToAllClientExceptMe(msg Message.IMessage)
}

type ISpaceTime interface {
	GetTime() time.Time
	GetDeltaTime() time.Duration
}

type IActorMgr interface {
	RegisterActor(actorType string, protoType IActor)

	AddActor(actorType string, actorID string, ownerUserID string, propData []byte, userData interface{}) (string, error)
	RemoveActor(actorID string) error

	GetActor(id string) (IActor, error)
	UpdateActors()

	TraversalActor(cb func(actor IActor))

	DestroyAllActor()
}

type ISpace interface {
	Util.ISafeCall
	ISpaceTime
	IActorMgr
	Prop.IPropOwner
	GetID() string
	GetUserData() interface{}
	GetOwnerUser() IUser

	GetUser(userID string) (IUser, error)
	TraversalUser(cb func(user BaseUser.IUser) bool)

	SendToAllClient(msg Message.IMessage)
	SendToAllClientExceptOne(msg Message.IMessage, exceptUserID string)

	DestroySelf()
}

type ISpaceServer interface {
	Core.ICore

	CreateSpace(id string, spacePropData []byte, userData interface{}) string
	DestroySpace(id string) error

	EnterSpace(userID string, spaceID string) error
	LeaveSpace(userID string) error

	GetSpace(id string) (ISpace, error)
	TraversalSpace(func(space ISpace))

	GetUser(id string) (BaseUser.IUser, error)
}

func Init(areaID string, serverID string, spacePT ISpace, userPT IUser, rpcProc interface{}) error {
	return _Init(areaID, serverID, spacePT, userPT, rpcProc)
}

func Destroy() {
	_Destroy()
}
