package internal

import (
	"Cinder/Base/Core"
	"Cinder/Mail/mailapi/types"
	"fmt"
	"strings"
	"testing"
	"time"

	assert "github.com/arl/assertgo"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/suite"
)

type _Suite struct {
	suite.Suite

	svc      *_MailService
	coreInst Core.ICore

	toUserID string
}

func newSuite(t *testing.T) *_Suite {
	return &_Suite{
		svc:      NewMailService(),
		coreInst: Core.New(),

		toUserID: "test_to_user_id",
	}
}

func (s *_Suite) SetupSuite() {
	info := Core.NewDefaultInfo()
	info.ServiceType = "test"
	svcID := fmt.Sprintf("%s_%s", info.ServiceType, uuid.NewV4())
	svcID = strings.Replace(svcID, "-", "", -1)
	assert.True(len(svcID) < 64) // nsq topic requires
	info.ServiceID = svcID
	if err := s.coreInst.Init(info); err != nil {
		panic(err)
	}
}

func (s *_Suite) TearDownSuite() {
	s.coreInst.Destroy()
}

func (s *_Suite) TestLoginLogout() {
	var err error
	err = s.svc.Login("userID")
	s.NoError(err)
	err = s.svc.Logout("userID")
	s.NoError(err)
}

func (s *_Suite) TestSend() {
	s.send()
}

func (s *_Suite) send() {
	now := time.Now()
	exp := now.Add(time.Second * 100)
	mail := types.Mail{
		ID:       "ID",
		From:     "From",
		FromNick: "FromNick",
		To:       s.toUserID,
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
	err := s.svc.Send(&mail)
	s.NoError(err)
}

func (s *_Suite) TestBroadcast() {
	s.broadcast()
}

func (s *_Suite) broadcast() {
	now := time.Now()
	exp := now.Add(time.Second * 100)
	mail := types.Mail{
		ID:       "ID",
		From:     "From",
		FromNick: "FromNick",

		Title: "Title",
		Body:  "Body",
		Attachments: []*types.Attachment{
			&types.Attachment{
				ItemID: 1111,
			},
			&types.Attachment{
				ItemID: 2222,
			},
		},

		SendTime:   now,
		ExpireTime: exp,
	}
	err := s.svc.Broadcast(&mail)
	s.NoError(err)
}

func (s *_Suite) TestListMail() {
	var err error
	const kUser = "userID"
	err = s.svc.Login(kUser) // 保证在DB中已创建
	s.NoError(err)
	defer s.svc.Logout(kUser)

	mails, errMails := s.svc.ListMail(kUser)
	_ = mails
	s.NoError(errMails)
}

func (s *_Suite) TestDelete() {
	err := s.svc.Delete(s.toUserID, "mailID")
	s.Error(err) // parse mail id
}

func (s *_Suite) TestMarkAsRead() {
	err := s.svc.MarkAsRead(s.toUserID, "mailID")
	s.Error(err) // parse mail id
}

func (s *_Suite) TestMarkAsUnread() {
	err := s.svc.MarkAsUnread(s.toUserID, "mailID")
	s.Error(err) // parse mail id
}

func (s *_Suite) TestMarkAttachmentsAsReceived() {
	err := s.svc.MarkAttachmentsAsReceived(s.toUserID, "mailID")
	s.Error(err) // parse mail id
}

func (s *_Suite) TestMarkAttachmentsAsUnreceived() {
	err := s.svc.MarkAttachmentsAsUnreceived(s.toUserID, "mailID")
	s.Error(err) // parse mail id
}

// 让 go test 执行测试
func TestMailService(t *testing.T) {
	suite.Run(t, newSuite(t))
}
