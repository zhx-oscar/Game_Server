package friend

import (
	"Cinder/Chat/chatapi"
	"Cinder/Chat/rpcproc/logic/rpc"
	"Cinder/Chat/rpcproc/logic/types"
	"Cinder/Chat/rpcproc/logic/user/internal/friend/callback"
	"Cinder/Chat/rpcproc/logic/user/internal/friend/dbutil"
	"Cinder/Chat/rpcproc/logic/user/internal/oflinfos"
	"errors"
	"fmt"
	"sync"

	assert "github.com/arl/assertgo"
	log "github.com/cihub/seelog"
)

type UserID = types.UserID

type _SrvIDGetter interface {
	GetSrvID() string
}

type FriendMgr struct {
	mtx sync.Mutex

	userID      UserID // 主人ID
	srvIDGetter _SrvIDGetter

	// 好友列表
	friends map[UserID]bool
}

func NewFriendMgr(userID UserID, srvIDGetter _SrvIDGetter) *FriendMgr {
	return &FriendMgr{
		userID:      userID,
		srvIDGetter: srvIDGetter,

		friends: make(map[UserID]bool),
	}
}

// SetFriends 设置好友列表
func (f *FriendMgr) SetFriends(ids []UserID) {
	f.mtx.Lock()
	defer f.mtx.Unlock()

	f.friends = make(map[UserID]bool)
	for _, id := range ids {
		f.friends[id] = true
	}
	delete(f.friends, f.userID) // 排除自身
}

/* 加好友流程如下：
1. A -> B 请求加好友
	1. A.SendRequest()
	2. B.RecvRequest(), B不在线则上线时触发
2. B -> A 响应加好友
	1. B.SendResponse()
	2. A.RecvResponse(), A不在线则上线时触发
如果不在线，请求或响应先存DB, 待上线时再处理。
*/

// SendRequest 发送申请加好友请求
func (f *FriendMgr) SendRequest(friendID UserID, reqInfo []byte) error {
	f.mtx.Lock()
	defer f.mtx.Unlock()

	if friendID == f.userID {
		return nil // 不允许，会mtx重入死锁
	}
	if hasReachedMaxFriendCount(len(f.friends)) {
		return chatapi.ErrSelfReachedMaxFriendCount
	}
	if err := checkPeerMaxAddFriendReq(friendID); err != nil {
		return err
	}

	// 不管是否在线都存DB，待应答后再删
	db := dbutil.FriendRequestsUtil(friendID)
	if err := db.Add(f.userID, reqInfo); err != nil {
		return fmt.Errorf("db add request: %w", err)
	}

	userFriendMgr := userMgr.GetUserFriendMgr(friendID)
	if userFriendMgr != nil {
		// 在线就直接处理
		userFriendMgr.RecvRequest(f.userID, reqInfo)
	}
	return nil
}

// RecvRequest 接收加好友请求
func (f *FriendMgr) RecvRequest(fromID UserID, reqInfo []byte) {
	f.mtx.Lock()
	defer f.mtx.Unlock()
	srvID := f.srvIDGetter.GetSrvID()
	req := &dbutil.DocFriendRequest{
		UserID:  f.userID,
		FromID:  fromID,
		ReqInfo: reqInfo,
	}
	go callback.RpcOneAddFriendReq(srvID, req)
}

// SendResponse 发送加好友响应
func (f *FriendMgr) SendResponse(fromID UserID, ok bool) (e error) {
	f.mtx.Lock()
	defer f.mtx.Unlock()

	// log.Debugf("(*FriendMgr).SendResponse() userID=%v fromID=%v ok=%v", f.userID, fromID, ok)
	if fromID == f.userID {
		return nil // 不允许，会mtx重入死锁
	}

	// 判断是否存在对应的请求
	if err := f.checkRequestExists(fromID); err != nil {
		return fmt.Errorf("check request exist: %w", err)
	}

	// 成功后删除请求记录
	defer func() {
		if e == nil {
			dbutil.FriendRequestsUtil(f.userID).DeleteRequest(fromID)
		}
	}()

	if ok {
		if err := checkMaxFriendCount(f.userID, len(f.friends), fromID); err != nil {
			return fmt.Errorf("check max friend count: %w", err)
		}

		err := f.addFriendOnOK(fromID) // 我同意fromID的加好友请求，此时即确立好友关系
		if err != nil {
			return fmt.Errorf("addFriend: %w", err)
		}
	}

	userFriendMgr := userMgr.GetUserFriendMgr(fromID)
	if userFriendMgr != nil {
		// 在线就直接处理
		userFriendMgr.RecvResponse(f.userID, ok)
		return nil
	}

	// 不在线就存DB
	db := dbutil.FriendResponsesUtil(fromID)
	if err := db.Add(f.userID, ok); err != nil {
		// TODO: 丢失响应，影响不大，因为好友已经添加
		return fmt.Errorf("db add response: %w", err)
	}
	return nil
}

// RecvResponse 接收加好友响应
func (f *FriendMgr) RecvResponse(fromID UserID, ok bool) {
	f.mtx.Lock()
	defer f.mtx.Unlock()
	// log.Debugf("(*FriendMgr).RecvResponse(fromID=%v, ok=%v)", fromID, ok)
	srvID := f.srvIDGetter.GetSrvID()
	rsp := &dbutil.DocFriendResponse{
		UserID:      f.userID,
		ResponderID: fromID,
		OK:          ok,
	}
	go callback.RpcOneAddFriendRet(srvID, rsp)
}

// HandleReqResp 从DB加载加好友请求和应答，处理。
func (f *FriendMgr) HandleReqResp() error {
	f.mtx.Lock()
	defer f.mtx.Unlock()

	// 请求记录将在收到应答后删除
	if err := f.handleRequests(); err != nil {
		return fmt.Errorf("pop and handle requests: %w", err)
	}
	// 应答记录处理后立即删除
	if err := f.popAndHandleResponses(); err != nil {
		return fmt.Errorf("pop and handle responses: %w", err)
	}
	return nil
}

// handleRequests 从DB加载加好友请求, 处理。
// 请求记录将在收到应答后删除.
func (f *FriendMgr) handleRequests() error {
	db := dbutil.FriendRequestsUtil(f.userID)
	requests, err := db.Query()
	if err != nil {
		return fmt.Errorf("db load: %w", err)
	}
	assert.True(requests != nil)

	// 处理请求
	srvID := f.srvIDGetter.GetSrvID()
	go callback.RpcManyAddFriendReqs(srvID, requests)
	return nil
}

// popAndHandleResponses 从DB加载并删除加好友应答，然后处理。
func (f *FriendMgr) popAndHandleResponses() error {
	db := dbutil.FriendResponsesUtil(f.userID)
	responses, err := db.Pop()
	if err != nil {
		return fmt.Errorf("db load: %w", err)
	}
	assert.True(responses != nil)

	// 处理应答
	srvID := f.srvIDGetter.GetSrvID()
	go callback.RpcManyAddFriendRets(srvID, responses)
	return nil
}

// DeleteFriendActive 删好友, 主动删对方。
func (f *FriendMgr) DeleteFriendActive(friendID UserID) error {
	f.mtx.Lock()
	defer f.mtx.Unlock()

	if _, ok := f.friends[friendID]; !ok {
		return nil // 无此好友，按成功处理
	}

	if err := dbutil.UsersUtil(f.userID).DeleteFriend(friendID); err != nil {
		return fmt.Errorf("db user util DeleteFriend: %w", err)
	}
	delete(f.friends, friendID)

	// 不管 friendID 是否在线，更改其DB中的好友列表
	if err := dbutil.UsersUtil(friendID).DeleteFriend(f.userID); err != nil {
		log.Errorf("db users util DeleteFriend error: %v", err) // TODO: 会产生数据不一致
	}
	if friendMgr := userMgr.GetUserFriendMgr(friendID); friendMgr != nil {
		friendMgr.DeleteFriendPassive(f.userID) // 如果对方在线，则同步更新其内存中的好友列表，并通知
	}
	return nil
}

// DeleteFriendPassive 被动删好友。即被人删。
// 仅更新内存好友列表，不会写DB, 因为DB数据由调用者更改，不管我是否在线。
func (f *FriendMgr) DeleteFriendPassive(friendID UserID) {
	f.mtx.Lock()
	defer f.mtx.Unlock()
	delete(f.friends, friendID)

	srvID := f.srvIDGetter.GetSrvID()
	targetID := string(f.userID)
	// 在协程中通知，释放锁
	go func() {
		ret := rpc.Rpc(srvID, "RPC_FriendDeleted", targetID, string(friendID))
		if ret.Err != nil {
			log.Errorf("RPC_FriendDeleted: %w", ret.Err)
		}
	}()
}

// GetFriendList 获取好友列表
func (f *FriendMgr) GetFriendList() ([]*chatapi.FriendInfo, error) {
	f.mtx.Lock()
	defer f.mtx.Unlock()

	ids := f.getFriendIDs()
	return oflinfos.GetFriendInfos(ids)
}

// GetFriendIDs 获取好友ID列表
func (f *FriendMgr) GetFriendIDs() []UserID {
	f.mtx.Lock()
	defer f.mtx.Unlock()

	return f.getFriendIDs()
}

func (f *FriendMgr) GetFriendCount() int {
	f.mtx.Lock()
	defer f.mtx.Unlock()

	return len(f.friends)
}

// getFriendIDs 获取好友ID列表
func (f *FriendMgr) getFriendIDs() []UserID {
	result := make([]UserID, 0, len(f.friends))
	for id, _ := range f.friends {
		result = append(result, id)
	}
	return result
}

// addFriend 添加好友，双方写DB并更新内存。
// 我同意fromID的加好友请求后触发。
func (f *FriendMgr) addFriendOnOK(fromID UserID) error {
	assert.True(fromID != f.userID)
	if _, ok := f.friends[fromID]; ok {
		return nil // 已为好友
	}

	if err := dbutil.UsersUtil(f.userID).AddFriend(fromID); err != nil {
		return fmt.Errorf("db users util add friend: %w", err)
	}
	f.friends[fromID] = true

	// 不管 fromID 是否在线，总是更新 DB
	if err := dbutil.UsersUtil(fromID).AddFriend(f.userID); err != nil {
		log.Errorf("db user util add friend error: %v", err) // TODO: 数据不一致
	}

	// 如果 fromID 在线，更新其内存
	if fromFriendMgr := userMgr.GetUserFriendMgr(fromID); fromFriendMgr != nil {
		fromFriendMgr.AddFriendWithoutDB(f.userID)
	}
	return nil
}

// AddFriendWithoutDB 添加好友，不写DB.
// 用于最终确认好友后更新内存。
func (f *FriendMgr) AddFriendWithoutDB(responderID UserID) {
	f.mtx.Lock()
	defer f.mtx.Unlock()
	assert.True(f.userID != responderID)
	f.friends[responderID] = true
}

// checkRequestExists 检查加好友请求是否存在，存在则返回 nil, 否则返回错误
func (f *FriendMgr) checkRequestExists(fromID UserID) error {
	if has, err := dbutil.FriendRequestsUtil(f.userID).Has(fromID); err != nil {
		return err
	} else if has {
		return nil
	}
	return errors.New("no request")
}
