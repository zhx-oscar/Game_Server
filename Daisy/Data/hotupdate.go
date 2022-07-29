package Data

// attachFuncList 附加热更新函数列表
var attachFuncList []func(isHotUpdate bool)

// AttachHotUpdate 附加热更新
func AttachHotUpdate(fun func(isHotUpdate bool)) {
	attachFuncList = append(attachFuncList, fun)
	fun(false)
}

// execHotUpdate 执行热更新
func execHotUpdate(isHotUpdate bool) {
	for _, fun := range attachFuncList {
		fun(isHotUpdate)
	}
}
