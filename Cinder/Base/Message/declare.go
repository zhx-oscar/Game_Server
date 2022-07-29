package Message

const (
	ID_Client_Validate_Req = 1
	ID_Client_Validate_Ret = 2

	ID_Forward_User_Message = 7
	ID_Enter_Space          = 8
	ID_Leave_Space          = 9

	ID_Client_Rpc_Req = 10
	ID_Client_Rpc_Ret = 11

	ID_Space_Broadcast_To_Client = 13

	ID_Enter_AOI = 14
	ID_Leave_AOI = 15

	ID_Batch_EnterAOI = 16
	ID_Clear_AOI      = 17

	ID_Heart_Beat = 20

	ID_User_Broadcast_Create  = 26
	ID_User_Broadcast_Destroy = 27

	ID_User_Rpc_Req = 30
	ID_User_Rpc_Ret = 31

	ID_User_Login_Req = 33
	ID_User_Login_Ret = 34

	ID_User_Logout_Req = 35

	ID_Rpc_Req = 37
	ID_Rpc_Ret = 38

	ID_MQ_Hello = 39

	ID_User_Prop_Notify  = 42
	ID_Space_Prop_Notify = 43
	ID_Actor_Prop_Notify = 44

	ID_Space_Owner_Change = 47

	ID_Client_User_Rpc = 50

	ID_Prop_Data_Req = 52
	ID_Prop_Data_Ret = 53

	ID_Prop_Notify = 55

	ID_PropObject_Open_Req  = 56
	ID_PropObject_Open_Ret  = 57
	ID_PropObject_Close_Req = 58
	ID_PropObject_Close_Ret = 59

	ID_Prop_Object_Prop_Notify   = 60
	ID_Forward_Users_Message     = 61
	ID_Forward_All_Users_Message = 62

	ID_Users_Rpc_Req     = 63
	ID_All_Users_Rpc_Req = 64

	ID_Actor_Refresh_Owner_User = 65

	ID_Prop_Data_Flush_Req = 66

	ID_User_Destroy_Req = 67
	ID_User_Destroy_Ret = 68

	ID_Prop_Data_Flush_Ret = 69

	ID_Mailbox_Req = 70
	ID_Mailbox_Ret = 71

	ID_Prop_Cache_Flush_Req = 72
	ID_Prop_Cache_Flush_Ret = 73

	ID_Rpc_Forbidden_Ret = 74
)

type ClientValidateReq struct {
	Version int32
	ID      string
	Token   string
	MsgSNo  uint32
}

func (m *ClientValidateReq) GetID() uint16 {
	return ID_Client_Validate_Req
}

type ClientValidateRet struct {
	OK           byte
	ERR          string
	UserPropType string
	ProtoDef     []byte
	UserData     []byte
	ServerTime   int64
}

func (m *ClientValidateRet) GetID() uint16 {
	return ID_Client_Validate_Ret
}

type ClientRpcReq struct {
	UserID     string
	SrvType    string
	MethodName string
	Args       []byte
	CBIndex    uint32
	Source     byte
}

func (m *ClientRpcReq) GetID() uint16 {
	return ID_Client_Rpc_Req
}

type ClientRpcRet struct {
	UserID  string
	CBIndex uint32
	Ret     []byte
	Source  byte
}

func (m *ClientRpcRet) GetID() uint16 {
	return ID_Client_Rpc_Ret
}

type RpcForbiddenRet struct {
	UserID  string
	CBIndex uint32
	Source  byte
}

func (m *RpcForbiddenRet) GetID() uint16 {
	return ID_Rpc_Forbidden_Ret
}

type ForwardUserMessage struct {
	TargetSrv string
	UserID    string
	MsgData   []byte
}

func (m *ForwardUserMessage) GetID() uint16 {
	return ID_Forward_User_Message
}

type ForwardUsersMessage struct {
	TargetSrv string
	UserIDS   []string
	MsgData   []byte
}

func (m *ForwardUsersMessage) GetID() uint16 {
	return ID_Forward_Users_Message
}

type ForwardAllUsersMessage struct {
	TargetSrv string
	MsgData   []byte
}

func (m *ForwardAllUsersMessage) GetID() uint16 {
	return ID_Forward_All_Users_Message
}

type EnterSpace struct {
	SpaceID   string
	OwnerID   string
	SpaceInfo []byte
}

func (m *EnterSpace) GetID() uint16 {
	return ID_Enter_Space
}

type LeaveSpace struct {
	SpaceID string
}

func (m *LeaveSpace) GetID() uint16 {
	return ID_Leave_Space
}

type SpaceBroadcastToClient struct {
	UserList     []string
	ExceptUserID string
	//MsgID        uint16
	//MsgFlag      byte
	MsgData []byte
}

func (m *SpaceBroadcastToClient) GetID() uint16 {
	return ID_Space_Broadcast_To_Client
}

type EnterAOIInfo struct {
	IsUser     bool
	ID         string
	Type       string
	OwnerID    string
	PropType   string
	Properties []byte
}

type BatchEnterAOI struct {
	Info []EnterAOIInfo
}

func (m *BatchEnterAOI) GetID() uint16 {
	return ID_Batch_EnterAOI
}

type ClearAOI struct {
	SpaceID string
}

func (m *ClearAOI) GetID() uint16 {
	return ID_Clear_AOI
}

type EnterAOI struct {
	Info EnterAOIInfo
}

func (m *EnterAOI) GetID() uint16 {
	return ID_Enter_AOI
}

type LeaveAOI struct {
	IsUser bool
	ID     string
}

func (m *LeaveAOI) GetID() uint16 {
	return ID_Leave_AOI
}

type HeartBeat struct {
	ClientSendTime int64
	ServerTime     int64
}

func (m *HeartBeat) GetID() uint16 {
	return ID_Heart_Beat
}

type UserBroadcastCreate struct {
	UserID  string
	SrvID   string
	SrvType string
}

func (m *UserBroadcastCreate) GetID() uint16 {
	return ID_User_Broadcast_Create
}

type UserBroadcastDestroy struct {
	UserID  string
	SrvID   string
	SrvType string
}

func (m *UserBroadcastDestroy) GetID() uint16 {
	return ID_User_Broadcast_Destroy
}

type UserRpcReq struct {
	UserID     string
	MethodName string
	Args       []byte
	RetID      string
}

func (m *UserRpcReq) GetID() uint16 {
	return ID_User_Rpc_Req
}

type UserRpcRet struct {
	UserID string
	Ret    []byte
	Err    string
	RetID  string
}

func (m *UserRpcRet) GetID() uint16 {
	return ID_User_Rpc_Ret
}

type UsersRpcReq struct {
	UserIDS    []string
	MethodName string
	Args       []byte
}

func (m *UsersRpcReq) GetID() uint16 {
	return ID_Users_Rpc_Req
}

type AllUsersRpcReq struct {
	MethodName string
	Args       []byte
}

func (m *AllUsersRpcReq) GetID() uint16 {
	return ID_All_Users_Rpc_Req
}

type UserLoginReq struct {
	UserID string
}

func (m *UserLoginReq) GetID() uint16 {
	return ID_User_Login_Req
}

type UserLoginRet struct {
	UserID   string
	PropType string
	PropDef  []byte
	UserData []byte
}

func (m *UserLoginRet) GetID() uint16 {
	return ID_User_Login_Ret
}

type UserLogoutReq struct {
	UserID string
}

func (m *UserLogoutReq) GetID() uint16 {
	return ID_User_Logout_Req
}

type RpcReq struct {
	MethodName string
	Args       []byte
	RetID      string
}

func (m *RpcReq) GetID() uint16 {
	return ID_Rpc_Req
}

type RpcRet struct {
	Ret   []byte
	Err   string
	RetID string
}

func (m *RpcRet) GetID() uint16 {
	return ID_Rpc_Ret
}

type MQHello struct {
	Greeting string
}

func (m *MQHello) GetID() uint16 {
	return ID_MQ_Hello
}

type UserPropNotify struct {
	Target     int8
	UserID     string
	MethodName string
	Args       []byte
}

func (m *UserPropNotify) GetID() uint16 {
	return ID_User_Prop_Notify
}

type SpacePropNotify struct {
	SpaceID    string
	MethodName string
	Args       []byte
}

func (m *SpacePropNotify) GetID() uint16 {
	return ID_Space_Prop_Notify
}

type ActorPropNotify struct {
	SpaceID    string
	ActorID    string
	MethodName string
	Args       []byte
}

func (m *ActorPropNotify) GetID() uint16 {
	return ID_Actor_Prop_Notify
}

type PropObjectPropNotify struct {
	ObjID      string
	MethodName string
	Args       []byte
}

func (m *PropObjectPropNotify) GetID() uint16 {
	return ID_Prop_Object_Prop_Notify
}

type SpaceOwnerChange struct {
	UserID string
}

func (m *SpaceOwnerChange) GetID() uint16 {
	return ID_Space_Owner_Change
}

type ClientActorProp struct {
	ObserverType string
	Data         []byte
	SrvMethod    string
}

type ClientActorProps struct {
	ActorID string
	Props   []ClientActorProp
}

type ClientUserRpc struct {
	Target     byte
	UserID     string
	MethodName string
	Args       []byte
}

func (m *ClientUserRpc) GetID() uint16 {
	return ID_Client_User_Rpc
}

type PropDataReq struct {
	TempID      string
	ID          string
	PropType    string
	ServiceType string
}

func (m *PropDataReq) GetID() uint16 {
	return ID_Prop_Data_Req
}

type PropDataRet struct {
	TempID string
	ID     string
	Data   []byte
	Err    string
}

func (m *PropDataRet) GetID() uint16 {
	return ID_Prop_Data_Ret
}

type PropDataFlushReq struct {
	TempID      string
	ID          string
	Type        string
	ServiceType string
}

func (m *PropDataFlushReq) GetID() uint16 {
	return ID_Prop_Data_Flush_Req
}

type PropDataFlushRet struct {
	TempID string
	Err    string
}

func (m *PropDataFlushRet) GetID() uint16 {
	return ID_Prop_Data_Flush_Ret
}

type PropNotify struct {
	ID          string
	Type        string
	MethodName  string
	Args        []byte
	ServiceType string
}

func (m *PropNotify) GetID() uint16 {
	return ID_Prop_Notify
}

type PropObjectOpenReq struct {
	ID     string
	UserID string
	Typ    string
}

func (m *PropObjectOpenReq) GetID() uint16 {
	return ID_PropObject_Open_Req
}

type PropObjectOpenRet struct {
	ID       string
	UserID   string
	SrvID    string
	PropType string
	PropData []byte
}

func (m *PropObjectOpenRet) GetID() uint16 {
	return ID_PropObject_Open_Ret
}

type PropObjectCloseReq struct {
	ID     string
	UserID string
}

func (m *PropObjectCloseReq) GetID() uint16 {
	return ID_PropObject_Close_Req
}

type PropObjectCloseRet struct {
	ID     string
	UserID string
}

func (m *PropObjectCloseRet) GetID() uint16 {
	return ID_PropObject_Close_Ret
}

type ActorRefreshOwnerUser struct {
	ActorID     string
	OwnerUserID string
	PropType    string
	PropData    []byte
}

func (m *ActorRefreshOwnerUser) GetID() uint16 {
	return ID_Actor_Refresh_Owner_User
}

type UserDestroyReq struct {
	UserID string
}

func (m *UserDestroyReq) GetID() uint16 {
	return ID_User_Destroy_Req
}

type UserDestroyRet struct {
	UserID string
}

func (m *UserDestroyRet) GetID() uint16 {
	return ID_User_Destroy_Ret
}

type MailboxReq struct {
	MailBoxType uint8
	MailBoxID   string
	SpaceID     string
	TargetID    string
	RetID       string
	MethodName  string
	Args        []byte
}

func (m *MailboxReq) GetID() uint16 {
	return ID_Mailbox_Req
}

type MailboxRet struct {
	MailboxID string
	Ret       []byte
	Err       string
	RetID     string
}

func (m *MailboxRet) GetID() uint16 {
	return ID_Mailbox_Ret
}

type PropCacheFlushReq struct {
	TempID      string
	ID          string
	Type        string
	ServiceType string
}

func (m *PropCacheFlushReq) GetID() uint16 {
	return ID_Prop_Cache_Flush_Req
}

type PropCacheFlushRet struct {
	TempID string
	Err    string
}

func (m *PropCacheFlushRet) GetID() uint16 {
	return ID_Prop_Cache_Flush_Ret
}
