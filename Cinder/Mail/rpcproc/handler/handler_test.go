package handler

import (
	"Cinder/Base/Const"
	"Cinder/Mail/mailapi/types"
	"Cinder/Mail/rpcproc/handler/bcmail"
	"Cinder/Mail/rpcproc/handler/delmark"
	"Cinder/Mail/rpcproc/handler/list"
	"Cinder/Mail/rpcproc/handler/loginout"
	"Cinder/Mail/rpcproc/handler/usermail"
	"Cinder/Mail/rpcproc/mockcore"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandlers(t *testing.T) {
	mockCore := &mockcore.MockCore{}
	mockCore.SetupMock()
	defer mockCore.TearDownMock()
	a := mock.Anything
	mockCore.On("GetSrvIDSByType", Const.Mail).Return([]string{"1", "2"}, nil)
	mockCore.On("GetSrvTypeByID", a).Return("", nil)
	mockCore.On("RpcByID", a, a, a).Return(mockCore.Ch())

	assert := require.New(t)
	const kUser = "user_test_list"
	loginout.Login(kUser, "peerSrvID") // 创建后才会有邮件
	defer loginout.Logout(kUser)

	var err error
	var mails []types.Mail
	usermail.Send(getMail(kUser))
	bcmail.Broadcast(getMail(kUser))
	mails, err = list.ListMails(kUser)
	assert.True(len(mails) >= 2)
	assert.NoError(err)
	// fmt.Printf("mails: %v", mails)

	for _, mail := range mails {
		err = delmark.MarkAsRead(kUser, mail.ID)
		assert.NoError(err)
		err = delmark.MarkAsUnread(kUser, mail.ID)
		assert.NoError(err)
		err = delmark.MarkAttachmentsAsReceived(kUser, mail.ID)
		assert.NoError(err)
		err = delmark.MarkAttachmentsAsUnreceived(kUser, mail.ID)
		assert.NoError(err)
		err = delmark.SetExpireTime(kUser, mail.ID, time.Now().Add(time.Second*10))
		if mail.IsBroadcast {
			assert.Error(err)
		} else {
			assert.NoError(err)
		}
		err = delmark.DeleteMail(kUser, mail.ID)
		assert.NoError(err)
	}

	usermail.Send(getMail(kUser))
	bcmail.Broadcast(getMail(kUser))
	mails, err = list.ListMails(kUser)
	assert.Equal(2, len(mails))
	assert.NoError(err)
	for _, mail := range mails {
		err = delmark.DeleteMail(kUser, mail.ID)
		assert.NoError(err)
	}

	usermail.Send(getMail(kUser))
	bcmail.Broadcast(getMail(kUser))
	mails, err = list.ListMails(kUser)
	assert.Equal(2, len(mails))
	assert.NoError(err)
	for _, mail := range mails {
		err = delmark.MarkAsUnread(kUser, mail.ID)
		assert.NoError(err)
		err = delmark.DeleteMail(kUser, mail.ID)
		assert.NoError(err)
	}

	usermail.Send(getMail(kUser))
	bcmail.Broadcast(getMail(kUser))
	mails, err = list.ListMails(kUser)
	assert.Equal(2, len(mails))
	assert.NoError(err)
	for _, mail := range mails {
		err = delmark.MarkAttachmentsAsReceived(kUser, mail.ID)
		assert.NoError(err)
	}

	mails, err = list.ListMails(kUser)
	assert.Equal(2, len(mails))
	assert.NoError(err)
	for _, mail := range mails {
		assert.Equal(true, mail.IsReceived)
		err = delmark.DeleteMail(kUser, mail.ID)
		assert.NoError(err)
	}

	usermail.Send(getMail(kUser))
	bcmail.Broadcast(getMail(kUser))
	mails, err = list.ListMails(kUser)
	assert.Equal(2, len(mails))
	assert.NoError(err)
	for _, mail := range mails {
		err = delmark.MarkAttachmentsAsUnreceived(kUser, mail.ID)
		assert.NoError(err)
		err = delmark.DeleteMail(kUser, mail.ID)
		assert.NoError(err)
	}

	// test set ext data
	usermail.Send(getMail(kUser))
	bcmail.Broadcast(getMail(kUser))
	mails, err = list.ListMails(kUser)
	assert.Equal(2, len(mails))
	assert.NoError(err)
	for _, mail := range mails {
		err = delmark.SetExtData(kUser, mail.ID, []byte("abc"))
		assert.NoError(err)
	}

	mails, err = list.ListMails(kUser)
	assert.Equal(2, len(mails))
	assert.NoError(err)
	for _, mail := range mails {
		assert.Equal([]byte("abc"), mail.ExtData)
		err = delmark.DeleteMail(kUser, mail.ID)
		assert.NoError(err)
	}
}

func getMail(to string) types.Mail {
	now := time.Now()
	exp := now.Add(time.Second * 100)
	mail := types.Mail{
		ID:       "ID",
		From:     "From",
		FromNick: "FromNick",
		To:       to,
		ToNick:   "ToNick",

		Title:  "Title",
		Body:   "Body",
		IsRead: false,
		Attachments: []*types.Attachment{
			&types.Attachment{
				ItemID: 1234,
				Count:  23,
				Data:   []byte("data"),
			},
			&types.Attachment{
				ItemID: 333,
			},
		},
		IsReceived: false,

		SendTime:   now,
		ExpireTime: exp,

		ExtData: []byte("ExtData"),
	}
	return mail
}
