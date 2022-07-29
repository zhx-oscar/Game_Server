package svcid

import (
	"fmt"
)

// FormatSvcID 用服务类型，区号，服务器号格式化一个服务ID.
// 服务ID用作 NSQ 的主题。
// 匹配服需要跨区发消息，所以需要规定服务ID的格式。
func FormatSvcID(serviceType string, areaID string, serverID string) string {
	return fmt.Sprintf("%s_%s_%s", serviceType, areaID, serverID)
}
