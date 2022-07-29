package maildoc

import (
	"Cinder/Mail/mailapi/types"
	"Cinder/Mail/rpcproc/handler/delmark/mailid"

	"go.mongodb.org/mongo-driver/bson/primitive"

	log "github.com/cihub/seelog"
)

// convert between types.Mail and MailDoc

func MailToDoc(mail types.Mail) Mail {
	atts := make([]*Attachment, 0, len(mail.Attachments))
	for _, att := range mail.Attachments {
		atts = append(atts, &Attachment{
			ItemID: att.ItemID,
			Count:  att.Count,
			Data:   att.Data,
		})
	}

	result := Mail{
		From:     mail.From,
		FromNick: mail.FromNick,
		To:       mail.To,
		ToNick:   mail.ToNick,
		Title:    mail.Title,
		Body:     mail.Body,

		State: &MailState{
			IsRead:                mail.IsRead,
			IsAttachmentsReceived: mail.IsReceived,
			ExtData:               mail.ExtData,
		},

		Attachments: atts,

		SendTime:   mail.SendTime,
		ExpireTime: mail.ExpireTime,
	}
	return result
}

func DocUserToMail(doc Mail, oid primitive.ObjectID) types.Mail {
	const isBroadcast = false
	return docToMail(doc, oid, isBroadcast)
}

func DocBcToMail(doc Mail, oid primitive.ObjectID) types.Mail {
	const isBroadcast = true
	return docToMail(doc, oid, isBroadcast)
}

func docToMail(doc Mail, oid primitive.ObjectID, isBroadcast bool) types.Mail {
	// 需要设置 mail id
	mailIDStr, err := mailid.GetMailIDStr(isBroadcast, oid)
	if err != nil {
		log.Error("failed to marshal mail id: %v", err)
		// 但是依旧继续，只是ID是无效的
	}

	atts := make([]*types.Attachment, 0, len(doc.Attachments))
	for _, att := range doc.Attachments {
		atts = append(atts, &types.Attachment{
			ItemID: att.ItemID,
			Count:  att.Count,
			Data:   att.Data,
		})
	}

	result := types.Mail{
		IsBroadcast: isBroadcast,

		ID:       mailIDStr,
		From:     doc.From,
		FromNick: doc.FromNick,
		To:       doc.To,
		ToNick:   doc.ToNick,

		Title:  doc.Title,
		Body:   doc.Body,
		IsRead: doc.State.IsRead,

		Attachments: atts,

		IsReceived: doc.State.IsAttachmentsReceived,

		SendTime:   doc.SendTime,
		ExpireTime: doc.ExpireTime,

		ExtData: doc.State.ExtData,
	}
	return result
}
