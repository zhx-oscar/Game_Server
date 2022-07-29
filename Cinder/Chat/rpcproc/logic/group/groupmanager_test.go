package group

import (
	"Cinder/Chat/rpcproc/logic/user"
	"Cinder/Chat/rpcproc/logic/usermgr"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	// 包变量 group._GroupMgr 应该在 init() 之前初始化
	user.SetGroupMgr(GetGroupMgr())
}

func TestSaveAllToDBOnExit(t *testing.T) {
	const kGroup = "group"
	assert := assert.New(t)
	GetGroupMgr().InsertMembersToGroup(kGroup, []UserID{"abc"})
	GetGroupMgr().SaveAllToDBOnExit()
	err := GetGroupMgr().DeleteGroup(kGroup)
	assert.Nil(err)
}

func TestGetGroupOfflineMessageMap(t *testing.T) {
	assert := assert.New(t)
	const kG1 = "test_group1"
	const kG2 = "test_group2"
	const kU1 = "test_user1"
	const kU2 = "test_user2"
	usermgr.GetUserMgr().RemoveUser(kU1)
	usermgr.GetUserMgr().RemoveUser(kU2)

	mgr := GetGroupMgr()

	var err error
	err = mgr.DeleteGroup(kG1)
	assert.Nil(err)
	err = mgr.DeleteGroup(kG2)
	assert.Nil(err)

	g1, errG1 := mgr.GetOrLoadGroup(kG1)
	assert.Nil(errG1)
	assert.NotNil(g1)
	err = g1.AddMembers([]UserID{kU1, kU2})
	assert.Nil(err)

	g2, errG2 := mgr.GetOrLoadGroup(kG2)
	assert.Nil(errG2)
	assert.NotNil(g2)
	err = g2.AddMembers([]UserID{kU1, kU2})
	assert.Nil(err)

	g1.SendGroupMessage(kU1, []byte("1111"))
	g1.SendGroupMessage(kU1, []byte("2222"))
	g1.SendGroupMessage(kU1, []byte("3333"))
	g2.SendGroupMessage(kU1, []byte("21111"))
	g2.SendGroupMessage(kU1, []byte("22222"))
	g2.SendGroupMessage(kU1, []byte("23333"))

	g2s := make(map[GroupID]SequenceID)
	g2s[kG1] = 0
	g2s[kG2] = 1
	msgMap := mgr.GetGroupOfflineMessageMap(g2s)
	for groupID, msgs := range msgMap {
		fmt.Printf("group: %v\n", groupID)
		for _, msg := range msgs {
			fmt.Printf("msg: %v\n", msg)
		}
	}

	err = mgr.DeleteGroup(kG1)
	assert.Nil(err)
	err = mgr.DeleteGroup(kG2)
	assert.Nil(err)
}
