package dbdoc

type UsersOfflineMessagesDoc struct {
	UserID     UserID `bson:"userID"`
	FromID     string `bson:"fromID`
	FromNick   string `bson:"fromNick"`
	FromData   []byte `bson:"fromData"`
	SendTime   int64  `bson:"sendTime"`
	MsgContent []byte `bson:"msgContent"`
}
