package mailapi

import (
	"Cinder/Mail/mailapi/types"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestUnmarshalMail(t *testing.T) {
	now := time.Now()
	exp := now.AddDate(0, 0, 7)
	mail := types.Mail{
		ID:     "ID",
		IsRead: true,
		Attachments: []*types.Attachment{
			&types.Attachment{
				ItemID: 123,
				Count:  234,
				Data:   []byte("data"),
			},
			&types.Attachment{
				ItemID: 2222,
			},
		},

		SendTime:   now,
		ExpireTime: exp,
		ExtData:    []byte("ExtData"),
	}

	assert := require.New(t)
	buf, err1 := json.Marshal(mail)
	assert.NoError(err1)

	m2, err := UnmarshalMail(buf)
	assert.NoError(err)
	assert.NotNil(m2)
	assert.Equal("ID", m2.ID)
	assert.Empty(m2.From)
	assert.True(m2.IsRead)
	assert.False(m2.IsReceived)
	assert.Equal(2, len(m2.Attachments))
	a1 := m2.Attachments[0]
	assert.EqualValues(123, a1.ItemID)
	assert.EqualValues(234, a1.Count)
	assert.Equal("data", string(a1.Data))
	a2 := m2.Attachments[1]
	assert.EqualValues(2222, a2.ItemID)
	assert.Zero(a2.Count)
	assert.Zero(a2.Data)
	assert.True(now.Equal(m2.SendTime))
	assert.True(exp.Equal(m2.ExpireTime))
	assert.Equal("ExtData", string(m2.ExtData))
}
