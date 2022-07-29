package dbdoc

type UsersGroupsDocKey struct {
	UserID  UserID  `bson:"userID"`
	GroupID GroupID `bson:"groupID"`
}

type UsersGroupsDoc struct {
	UserID     UserID     `bson:"userID"`
	GroupID    GroupID    `bson:"groupID"`
	SequenceID SequenceID `bson:"sequenceID"`
}
