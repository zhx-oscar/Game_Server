package ErrorCode

// 邮件相关 660~699

const (
	MailSuccess       = 0
	MailParamErr      = 2201 // 邮件参数错误
	MailUnknowErr     = 2202 // 邮件未知错误
	MailDBError       = 2203 // 邮件服错误
	MailBoxNoExist    = 2204 // 邮箱不存在
	MailNoExist       = 2205 // 邮件不存在
	MailAthNoExist    = 2206 // 邮件附件不存在
	MailTimeExpire    = 2207 // 邮件过期
	MailBagAddErr     = 2208 // 邮件附件加入背包失败
	MailTitleLenErr   = 2209 // 邮件标题太长
	MailContentLenErr = 2210 // 邮件内容太长
	MailAthLenErr     = 2211 // 附件超过限制，发送失败
	MailBoxEmpty      = 2212 // 邮箱是空的，没有邮件
	MailAllTimeExpire = 2213 // 所有邮件过期
)
