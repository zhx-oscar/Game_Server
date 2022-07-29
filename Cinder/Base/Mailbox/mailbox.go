package Mailbox

type iMailBoxCtrl interface {
	getType() uint8
	getMailboxID() string
	onMailboxRet(retID string, err error, ret []interface{})
	marshal() ([]byte, error)
	unmarshal([]byte) error
}
