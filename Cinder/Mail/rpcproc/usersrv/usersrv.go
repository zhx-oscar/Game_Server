package usersrv

import (
	"Cinder/Base/Const"
	"Cinder/Base/Core"
	"Cinder/Mail/rpcproc/rpc"
	"Cinder/Mail/rpcproc/usersrv/db"
	"fmt"
	"sort"
	"sync"

	log "github.com/cihub/seelog"
)

// _UserSrvIDs 管理用户服类型
type _UserSrvIDs struct {
	mtx sync.Mutex
	ids map[string]bool
}

var userSrvIDs = &_UserSrvIDs{
	ids: make(map[string]bool),
}

// LoadFromDB 从DB加载 userSrvID 列表
func LoadFromDB() error {
	log.Infof("load user server IDs from DB")
	ids, err := db.GetUserSrvIDs().Load()
	if err != nil {
		return nil
	}
	userSrvIDs.Set(ids)
	return nil
}

// Insert 插入用户服ID, 写DB, 并广播
func InsertAndBroadcast(srvID string) error {
	if !userSrvIDs.Insert(srvID) {
		return nil
	}

	// DB 插入 srvID
	if err := db.GetUserSrvIDs().Insert(srvID); err != nil {
		userSrvIDs.Delete(srvID) // 回滚
		return err
	}

	// 插入成功，同步到所有邮件服
	ids, err := Core.Inst.GetSrvIDSByType(Const.Mail)
	if err != nil {
		log.Errorf("GetSrvIDSByType(ConstMail) error: %v", err) // 不应该出错
		return nil
	}

	rpc.RpcByIDs(ids, "RPC_SyncUserSrvID", srvID)
	return nil
}

// Sync 同步用户服ID
func Sync(srvID string) {
	_ = userSrvIDs.Insert(srvID)
}

func GetUserSrvIDs() []string {
	userSrvIDs.Check()
	return userSrvIDs.Get()
}

// Insert 插入用户服ID, 返回是否插入成功
func (l *_UserSrvIDs) Insert(srvID string) bool {
	l.mtx.Lock()
	defer l.mtx.Unlock()

	if _, ok := l.ids[srvID]; ok {
		return false
	}
	l.ids[srvID] = true
	l.logServerIDs()
	return true
}

func (l *_UserSrvIDs) logServerIDs() {
	log.Infof("server IDs: %s", l.getServerIDsStr())
}

func (l *_UserSrvIDs) getServerIDsStr() string {
	a := make([]string, 0, len(l.ids))
	for id, _ := range l.ids {
		a = append(a, id)
	}
	sort.Strings(a)
	return fmt.Sprintf("%v", a)
}

// Remove 删除用户服ID
func (l *_UserSrvIDs) Delete(srvID string) {
	l.mtx.Lock()
	defer l.mtx.Unlock()

	delete(l.ids, srvID)
	log.Infof("delete user server ID: %v, now server IDs: %s", srvID, l.getServerIDsStr())
}

func (l *_UserSrvIDs) Set(ids map[string]bool) {
	l.mtx.Lock()
	defer l.mtx.Unlock()

	l.ids = ids
	l.logServerIDs()
}

func (l *_UserSrvIDs) Get() []string {
	l.mtx.Lock()
	defer l.mtx.Unlock()

	result := make([]string, 0, len(l.ids))
	for id, _ := range l.ids {
		result = append(result, id)
	}
	return result
}

func (l *_UserSrvIDs) Check() {
	l.mtx.Lock()
	defer l.mtx.Unlock()

	failedIDs := make([]string, 0, len(l.ids))
	for id, _ := range l.ids {
		if _, err := Core.Inst.GetSrvTypeByID(id); err != nil {
			failedIDs = append(failedIDs, id)
		}
	}
	if len(failedIDs) == 0 {
		return
	}

	for _, id := range failedIDs {
		delete(l.ids, id)
	}

	log.Infof("remove user server IDs: %v, now server IDs: %s", failedIDs, l.getServerIDsStr())
	db.GetUserSrvIDs().RemoveIDs(failedIDs)
}
