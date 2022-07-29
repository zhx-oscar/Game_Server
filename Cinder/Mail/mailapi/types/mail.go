package types

import (
	"time"
)

type Mail struct {
	IsBroadcast bool // 是否全服广播邮件

	ID       string // 唯一标识，发送时为空，接收时才会有
	From     string
	FromNick string
	To       string
	ToNick   string

	Title       string
	Body        string
	IsRead      bool
	Attachments []*Attachment
	IsReceived  bool

	SendTime   time.Time
	ExpireTime time.Time

	// 额外数据, 可用 SetExtData() 更改
	// 可保存任意数据, 例如阅读进度，标签.
	ExtData []byte
}

type Attachment struct {
	ItemID uint32
	Count  uint32
	Data   []byte
}
