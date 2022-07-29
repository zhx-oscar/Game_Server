package stress

import (
	"Cinder/Base/Util"
	"Cinder/Chat/chatapi"
	"strconv"
	"sync"

	log "github.com/cihub/seelog"
)

var (
	rolePrefix   = "role_" + Util.GetGUID() + "_"
	groupPrefix  = "group_" + Util.GetGUID() + "_"
	goroutineMap sync.Map // [int]interface{}
)

func getSuffix(iGo GoroutineIndex, i _RunIndex) string {
	return strconv.Itoa(int(iGo)) + "_" + strconv.Itoa(int(i))
}

func getRoleID(iGo GoroutineIndex, i _RunIndex) string {
	return rolePrefix + getSuffix(iGo, i)
}

func getGroupID(iGo GoroutineIndex, i _RunIndex) string {
	return groupPrefix + getSuffix(iGo, i)
}

func panicIfError(err error) {
	if err != nil {
		flushAndPanic(err.Error())
	}
}

func flushAndPanic(msg string) {
	log.Flush()
	panic(msg)
}

func firstLoginUser(iGo GoroutineIndex, i _RunIndex) chatapi.IUser {
	if i != 0 {
		return getUser(iGo)
	}
	user, err := chatapi.Login(getRoleID(iGo, i), "nick", []byte("data"))
	panicIfError(err)
	goroutineMap.Store(iGo, user)
	return user
}

func getUser(iGo GoroutineIndex) chatapi.IUser {
	i, ok := goroutineMap.Load(iGo)
	if !ok {
		flushAndPanic("can not get user")
	}
	return i.(chatapi.IUser)
}

func firstCreateGroup(iGo GoroutineIndex, i _RunIndex) chatapi.IGroup {
	if i != 0 {
		return getGroup(iGo)
	}
	grp, err := chatapi.GetOrCreateGroup(getGroupID(iGo, i))
	panicIfError(err)
	goroutineMap.Store(iGo, grp)
	return grp
}

func getGroup(iGo GoroutineIndex) chatapi.IGroup {
	i, ok := goroutineMap.Load(iGo)
	if !ok {
		flushAndPanic("can not get group")
	}
	return i.(chatapi.IGroup)
}
