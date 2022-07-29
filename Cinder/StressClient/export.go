package StressClient

/*
type IUser interface {
	GetID() string
	GetProp() Prop.IProp
}

type IActor interface {
	GetID() string
	GetType() string
	GetProp() Prop.IProp
}

type ISpace interface {
	AddDelegateToEnterSpace(del func())
	AddDelegateToLeaveSpace()
	AddDelegateToAddUser(del func(user IUser))
	AddDelegateToRemoveUser(del func(user IUser))
	AddDelegateToAddActor(del func(actor IActor))
	AddDelegateToRemoveActor(del func(actor IActor))

	GetID() string
	GetProp() Prop.IProp
}
*/

type IClient interface {
	GetID() string
	GetUserName() string

	AddDelegateToNetMessage(del func(string, ...interface{}))

	Rpc(methodName string, args ...interface{})
	SpaceRpc(methodName string, args ...interface{})
	//ConnectToPropObject(id string)
	//DisconnectToPropObject(id string)

	//GetSpace() ISpace
	//GetLocalUser() IUser
}

type IInit interface {
	Init()
}

type IDestroy interface {
	Destroy()
}

type ILoop interface {
	Loop()
}

type IClientPool interface {
	Start()
	Close()
}
