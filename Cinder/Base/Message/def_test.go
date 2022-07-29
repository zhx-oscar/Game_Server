package Message

import "testing"

func Test_Def(t *testing.T) {

	def := newDef()

	def.addDef(&ClientValidateReq{})
	def.addDef(&ClientValidateRet{})

	m, _ := def.fetchMessage(1)

	if m == nil {
		println("wrong")
	} else {
		println("right")
	}
}
