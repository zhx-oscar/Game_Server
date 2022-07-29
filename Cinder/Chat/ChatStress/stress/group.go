package stress

import (
	"Cinder/Chat/chatapi"
)

func init() {
	register(addIntoGroupAndKick)
	register(getGroupMembers)
	register(sendGroupMessage0)
	register(sendGroupMessage1)
	register(sendGroupMessage2)
	register(sendGroupMessage3)
	register(sendGroupMessage10)
	register(sendGroupMessage100)
	register(getHistoryMessages)
}

func addIntoGroupAndKick(iGo GoroutineIndex, i _RunIndex) {
	grp := firstCreateGroup(iGo, i)
	errAdd := grp.AddIntoGroup("member")
	panicIfError(errAdd)
	errKick := grp.KickFromGroup("member")
	panicIfError(errKick)
}

func getGroupMembers(iGo GoroutineIndex, i _RunIndex) {
	grp := firstCreateGroup(iGo, i)
	_, err := grp.GetGroupMembers()
	panicIfError(err)
}

func sendGroupMessage0(iGo GoroutineIndex, i _RunIndex) {
	sendGroupMessageN(iGo, i, 0)
}

func sendGroupMessage1(iGo GoroutineIndex, i _RunIndex) {
	sendGroupMessageN(iGo, i, 1)
}

func sendGroupMessage2(iGo GoroutineIndex, i _RunIndex) {
	sendGroupMessageN(iGo, i, 2)
}

func sendGroupMessage3(iGo GoroutineIndex, i _RunIndex) {
	sendGroupMessageN(iGo, i, 3)
}

func sendGroupMessage10(iGo GoroutineIndex, i _RunIndex) {
	sendGroupMessageN(iGo, i, 10)
}

func sendGroupMessage100(iGo GoroutineIndex, i _RunIndex) {
	sendGroupMessageN(iGo, i, 100)
}

func sendGroupMessageN(iGo GoroutineIndex, i _RunIndex, nGroupSize int) {
	grp := firstCreateGroup(iGo, i)
	if i == 0 && nGroupSize <= 10 { // n 10 以上需要 setup() 中初始化
		for j := 0; j < nGroupSize; j++ {
			roleID := getRoleID(iGo, _RunIndex(j))
			err := grp.AddIntoGroup(roleID)
			panicIfError(err)
			_, errLogin := chatapi.Login(roleID, "nick", []byte("data"))
			panicIfError(errLogin)
		}
	}
	roleID0 := getRoleID(iGo, 0)
	err := grp.SendGroupMessage(roleID0, []byte("msgContent"))
	panicIfError(err)
}

func getHistoryMessages(iGo GoroutineIndex, i _RunIndex) {
	grp := firstCreateGroup(iGo, i)
	if grp.GetHistoryMessages(10) == nil {
		flushAndPanic("failed to get group history messages")
	}
}
