// callback 包发送 RPC_AddFriendReq 和 RPC_AddFriendRet 回调
// 为避免死锁，回调时应该先释放锁，如新建协程执行。
package callback
