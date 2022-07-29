// package oflinfos 管理离线用户FriendInfo, 并为在线和离线用户提供统一的获取FriendInfo的接口。
// Mgr caches FriendInfo of offline users.
// 管理器是离线用户 FriendInfo 缓存。
// 限制个数。
// 获取时可按需从DB加载。
// 上线时删除。离线时添加。
// GetFriendInfo() 优先查询在线用户，然后查缓存，然后读DB来获取FriendInfo.
package oflinfos
