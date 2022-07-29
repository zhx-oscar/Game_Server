package Message

func init() {
	def.addDef(&ClientValidateReq{})
	def.addDef(&ClientValidateRet{})
	def.addDef(&ClientRpcReq{})
	def.addDef(&ClientRpcRet{})
	def.addDef(&ForwardUserMessage{})
	def.addDef(&ForwardUsersMessage{})
	def.addDef(&ForwardAllUsersMessage{})
	def.addDef(&EnterSpace{})
	def.addDef(&LeaveSpace{})
	def.addDef(&SpaceBroadcastToClient{})
	def.addDef(&BatchEnterAOI{})
	def.addDef(&ClearAOI{})
	def.addDef(&EnterAOI{})
	def.addDef(&LeaveAOI{})
	def.addDef(&HeartBeat{})

	def.addDef(&UserBroadcastCreate{})
	def.addDef(&UserBroadcastDestroy{})
	def.addDef(&UserRpcReq{})
	def.addDef(&UserRpcRet{})

	def.addDef(&UsersRpcReq{})
	def.addDef(&AllUsersRpcReq{})

	def.addDef(&UserLoginReq{})
	def.addDef(&UserLoginRet{})
	def.addDef(&UserLogoutReq{})
	def.addDef(&RpcReq{})
	def.addDef(&RpcRet{})
	def.addDef(&MQHello{})

	def.addDef(&UserPropNotify{})
	def.addDef(&SpacePropNotify{})
	def.addDef(&ActorPropNotify{})
	def.addDef(&PropObjectPropNotify{})

	def.addDef(&SpaceOwnerChange{})
	def.addDef(&ClientUserRpc{})

	def.addDef(&PropDataReq{})
	def.addDef(&PropDataRet{})
	def.addDef(&PropDataFlushReq{})
	def.addDef(&PropDataFlushRet{})
	def.addDef(&PropNotify{})

	def.addDef(&PropObjectOpenReq{})
	def.addDef(&PropObjectOpenRet{})
	def.addDef(&PropObjectCloseReq{})
	def.addDef(&PropObjectCloseRet{})

	def.addDef(&ActorRefreshOwnerUser{})

	def.addDef(&UserDestroyReq{})
	def.addDef(&UserDestroyRet{})

	def.addDef(&MailboxReq{})
	def.addDef(&MailboxRet{})

	def.addDef(&PropCacheFlushReq{})
	def.addDef(&PropCacheFlushRet{})
}
