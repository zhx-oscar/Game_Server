package conf

import "Daisy/Data"

func init() {
	Data.AttachHotUpdate(func(isHotUpdate bool) {
		HotUpdateConfs(false, true)
	})
}
