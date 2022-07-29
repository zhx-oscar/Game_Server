package stress

import (
	"Cinder/Chat/chatapi"

	"github.com/spf13/viper"
)

func Setup(iGo GoroutineIndex) {
	test := viper.GetString("test")
	switch test {
	case "sendGroupMessage100":
		setupSendGroupMessage100(iGo)
	}
}

func setupSendGroupMessage100(iGo GoroutineIndex) {
	// fmt.Printf("setupSendGroupMessage100\n")
	grp := firstCreateGroup(iGo, 0)
	for j := 0; j < 100; j++ {
		roleID := getRoleID(iGo, _RunIndex(j))
		err := grp.AddIntoGroup(roleID)
		panicIfError(err)
		_, errLogin := chatapi.Login(roleID, "nick", []byte("data"))
		panicIfError(errLogin)
	}
}
