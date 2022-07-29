package main

import (
	cConst "Cinder/Base/Const"
	"Cinder/Mail/mailapi"
	mailtypes "Cinder/Mail/mailapi/types"
	"Daisy/Const"
	"Daisy/Data"
	"Daisy/ErrorCode"
	"Daisy/ItemProto"
	"Daisy/Proto"
	"encoding/json"
	"sort"
	"strconv"
	"strings"
	"time"
)

var mailConfigMaxNum = 100
var _MailReadTime = "ReadTime"
var _MailBroadCastExpireTime = "BroadCastExpireTime"

func (r *_Role) gm_fiendDataChange(args string) int32 {
	user := r.GetOwnerUser()
	if user == nil {
		r.Error("[findNewTitle] role's user is nil")
	} else {
		user.Rpc(cConst.Game, "RPC_UpdateChatUserActivateData")
	}
	return 1
}

// p1: 标题, p2: 内容， p4: 道具id p5: 道具类型 p6: 道具数量
// gm_mail 发送广播邮件
func (r *_Role) gm_mail(args string) int32 {
	list := strings.Split(args, ",")

	title, content, to, toName := "", "", "", ""
	var items []*Proto.MailAttachment
	_len := len(list)
	if _len >= 1 {
		title = list[0]
	}
	if _len >= 2 {
		content = list[1]
	}
	if _len >= 3 {
		_len -= 2
		_len /= 3
		for i := 0; i < _len; i++ {
			id, err := strconv.Atoi(list[i*3+2])
			if err != nil {
				return ErrorCode.Success
			}
			ty, err := strconv.Atoi(list[i*3+3])
			if err != nil {
				return ErrorCode.Success
			}
			num, err := strconv.Atoi(list[i*3+4])
			if err != nil {
				return ErrorCode.Success
			}
			items = append(items, &Proto.MailAttachment{ItemID: uint32(id), ItemType: uint32(ty), ItemNum: uint32(num)})
		}
	}
	return r.SendMail(title, content, to, toName, true, items)

}

// p1: 标题, p2: 内容， p3：目标用户的roleID  (role.GetOwnerID()) p4:道具id p5: 道具类型 p6 :道具数量
// gm_mailto 发送邮件给固定的玩家
func (r *_Role) gm_mailto(args string) int32 {
	list := strings.Split(args, ",")
	title, content, to, toName := "", "", "", ""
	var items []*Proto.MailAttachment
	_len := len(list)
	if _len < 3 {
		return -1
	}

	title = list[0]
	content = list[1]
	to = list[2]

	if _len >= 4 {
		_len -= 3
		_len /= 3
		for i := 0; i < _len; i++ {
			id, err := strconv.Atoi(list[i*3+3])
			if err != nil {
				return ErrorCode.Success
			}
			ty, err := strconv.Atoi(list[i*3+4])
			if err != nil {
				return ErrorCode.Success
			}
			num, err := strconv.Atoi(list[i*3+5])
			if err != nil {
				return ErrorCode.Success
			}
			items = append(items, &Proto.MailAttachment{ItemID: uint32(id), ItemType: uint32(ty), ItemNum: uint32(num)})
		}
	}

	r.SendMail(title, content, to, toName, false, items)

	return 0
}

// mailOnline 登陆user(在上线时使用)
func (r *_Role) mailOnline() {
	r.Debug("[role] mailOnline")
	r.ms = mailapi.GetMailService()
	if r.ms == nil {
		r.Error("[mailOnline] r.mailService is nil")
		return
	}
	err := r.ms.Login(r.GetOwnerUserID())
	if err != nil {
		r.Errorf("[mailOnline] mail login error :%s", err)
		return
	}

	mails, errList := r.ms.ListMail(r.GetOwnerUserID())
	if errList != nil {
		r.Errorf("[mailOnline] list mail error: %s", errList)
		return
	}

	// 排序
	sortMails := r.mailSortFromMailApi(mails)

	// 超过邮件数量上限检查(邮件服不会发送超过有效期的邮件)
	mailConfigMaxNum = int(Data.GetMailGeneralConfig(Data.MailConfigMaxNum))

	if len(sortMails) > mailConfigMaxNum {
		sortMails = sortMails[len(sortMails)-mailConfigMaxNum:]
	}

	// 填充邮箱
	for i := 0; i < len(sortMails); i++ {
		r.prop.Data.MailBox.Mails = append(r.prop.Data.MailBox.Mails, sortMails[i])
	}
}

// mailSort 邮件排序  排序规则（发送时间大>发送时间小>已读时间小>已读时间大）
func (r *_Role) mailSortFromMailApi(mails []*mailapi.Mail) []*Proto.Mail {
	now := time.Now().Unix()
	sortMails := make([]*Proto.Mail, 0)
	for i := 0; i < len(mails); i++ {
		mail := r.convertToProtoMail(mails[i])

		// 如果过期了，不添加到内存中
		if mail.ExpireTime < now {
			continue
		}
		sortMails = append(sortMails, mail)
	}
	if len(sortMails) <= 1 {
		return sortMails
	}
	return r.mailSort(sortMails)
}

type SortMail []*Proto.Mail

func (m SortMail) Len() int {
	return len(m)
}

func (m SortMail) Less(i, j int) bool {
	// m[i] 未读
	if !mailIsReadState(m[i]) {
		if !mailIsReadState(m[j]) {
			if m[i].ReadTime > m[j].ReadTime {
				return true
			} else {
				return false
			}
		} else {
			return true
		}
	} else {
		// m[i] 已读
		if !mailIsReadState(m[j]) {
			// m[j] 未读
			return false
		} else {
			// m[j] 已读
			if m[i].ReadTime < m[j].ReadTime {
				return true
			} else {
				return false
			}
		}
	}
}

func (m SortMail) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

func (r *_Role) mailSort(mails []*Proto.Mail) []*Proto.Mail {
	sortMail := SortMail{}
	for i := 0; i < len(mails); i++ {
		sortMail = append(sortMail, mails[i])
	}
	sort.Sort(sortMail)
	return sortMail
}

// 判断一封邮件是否是已读状态
func mailIsReadState(mail *Proto.Mail) bool {
	if len(mail.Attachments) > 0 {
		if mail.IsReceived {
			return true
		} else {
			return false
		}
	} else {
		if mail.IsRead {
			return true
		} else {
			return false
		}
	}
}

// mailOffline 退出user(在下线时使用)
func (r *_Role) mailOffline() {
	r.Debug("[role] mailOffline ")
	if r.ms == nil {
		r.Error("[mailOffline] r.mailservice is nil")
		return
	}
	err := r.ms.Logout(r.GetOwnerUserID())
	if err != nil {
		r.Errorf("[mailOffline] mail logout error: %s", err)
		return
	}
}

func (r *_Role) SendMail(title, content, to, toName string, isBroadcast bool, items []*Proto.MailAttachment) int32 {
	_mail := &mailapi.Mail{}

	_mail.IsBroadcast = isBroadcast
	_mail.To = to
	_mail.ToNick = toName
	_mail.From = r.GetOwnerUserID()
	_mail.FromNick = r.prop.Data.Base.GetName()
	_mail.Title = title
	_mail.Body = content

	if r.IsOverAttachmentLimit(items) {
		r.Error("[SendMail] 附件物品超过8种，不允许发送邮件")
		return ErrorCode.MailAthLenErr
	}
	if items != nil {
		for _, v := range items {
			data, _ := v.Marshal()
			_mail.Attachments = append(_mail.Attachments, &mailtypes.Attachment{Data: data})
		}
	}

	// 设置过期时间
	var expireSecs uint32
	_mail.SendTime = time.Now()
	if len(_mail.Attachments) > 0 {
		expireSecs = Data.GetMailGeneralConfig(Data.MailConfigAttachExpireTime)
	} else {
		expireSecs = Data.GetMailGeneralConfig(Data.MailConfigExpireTime)
	}
	_mail.ExpireTime = _mail.SendTime.Add(time.Second * time.Duration(expireSecs))

	// 设置已读时间（邮件排序使用）
	profile_ := map[string]string{
		_MailReadTime:            strconv.Itoa(int(_mail.SendTime.Unix())),
		_MailBroadCastExpireTime: strconv.Itoa(int(_mail.ExpireTime.Unix())),
	}
	profile, err := json.Marshal(profile_)
	if err != nil {
		r.Error("[SendMail] mail Marshal ExtData failed", err)
	}
	_mail.ExtData = profile

	if isBroadcast && mailapi.GetMailService().Broadcast(_mail) == nil {
		return ErrorCode.MailSuccess
	} else if mailapi.GetMailService().Send(_mail) == nil {
		return ErrorCode.MailSuccess
	}

	return ErrorCode.MailUnknowErr
}

// SyncInsertMailFromSrv 插入邮件并同步到客户端
func (r *_Role) SyncInsertMailFromSrv(mail *mailtypes.Mail) {
	r.Info("[SyncInsertMailFromSrv] begin ", r.GetOwnerUserID())
	protoMail := r.convertToProtoMail(mail)
	if protoMail == nil {
		return
	}

	// 移除过期邮件
	r.removeExpireMail()
	// 排序
	r.prop.Data.MailBox.Mails = r.mailSort(r.prop.Data.MailBox.Mails)
	// 限制个数（如果超过100用新邮件代替老邮件） 移除最后一封
	mails := r.getRoleMails()
	if len(mails) > mailConfigMaxNum-1 {
		if r.ms == nil {
			r.Error("[mailOnline] r.mailService is nil")
			return
		}
		// 删除老邮件
		index := len(mails) - 1
		r.ms.Delete(r.GetOwnerUserID(), mails[index].MailID)
		r.prop.SyncMailBoxRemoveMail(mails[index].MailID)
	}
	myMail := r.convertToProtoMail(mail)

	r.prop.SyncInsertMail(myMail)

	notify := r.getRedPointKey(Const.RedPointType_notifyMail, "")
	_, ok := r.prop.Data.RedPointsData[notify]

	// 没有红点，增加红点
	if !ok {
		redPointData := &Proto.RedPointInfo{
			Value:      1,
			CreateTime: time.Now().Unix(),
		}
		r.prop.SyncAddRedPoint(notify, redPointData)
	}
}

//--------------------------------------------------
// rpc 函数调用系列

// oneKeyGetAttachments 一键领取邮件附件
func (r *_Role) oneKeyGetAttachmentsOfMail() (int32, *Proto.MailAwardItems) {
	mails := r.getRoleMails()
	if mails == nil {
		return ErrorCode.MailBoxEmpty, nil
	}
	unreceive := make([]string, 0)
	for _, v := range mails {
		// 有附件并且附件未接收
		if !v.IsReceived && len(v.Attachments) > 0 {
			unreceive = append(unreceive, v.MailID)
		}
	}

	ec, award, errors := r.getBatchAttachments(unreceive)
	result := &Proto.MailAwardItems{
		AwardItems: award,
		BagErrors:  errors,
	}
	r.notifyMailRedPoint()

	return ec, result
}

// getAttachments 领取一封邮件的附件
func (r *_Role) getAttachmentsOfMail(mailId string) (int32, *Proto.MailAwardItems) {
	ec, result, errors := r.getOneAttachments(mailId)
	ret := &Proto.MailAwardItems{
		AwardItems: result,
		BagErrors:  errors,
	}
	// 返回码不为空的时候，可能结构里会返回邮件错误码，所以还是要先给结构赋值
	if ec != ErrorCode.MailSuccess {
		return ec, ret
	}
	return ErrorCode.MailSuccess, ret
}

// markMailAsRead 将邮件设置为已读
func (r *_Role) markMailAsRead(mailId string) int32 {
	mail := r.getMailFromMailBox(mailId)
	if mail == nil {
		return ErrorCode.Failure
	}
	now := time.Now()
	if mail.ExpireTime < now.Unix() {
		r.prop.SyncMailBoxRemoveMail(mail.MailID)
		return ErrorCode.MailTimeExpire
	}
	if mail.IsRead {
		return ErrorCode.MailSuccess
	}
	if r.ms == nil {
		r.Error("[markMailAsRead] r.mailService is nil")
		return ErrorCode.MailUnknowErr
	}

	err := r.ms.MarkAsRead(r.GetOwnerUserID(), mailId)
	if err != nil {
		r.Errorf("[MarkMailAsRead] mark mail as read error: %s", err)
		return ErrorCode.MailDBError
	}

	r.prop.SyncUpdateMailRead(mailId, true)

	// 没有附件的已读邮件设置过期时间 && 设置已读时间
	if len(mail.Attachments) == 0 {
		expireTime, ec := r.setMailExpireTime(mailId, Data.MailConfigHasReadExpireTime, now, mail.IsBroadcast)
		if ec != ErrorCode.Success {
			return ec
		}
		ec = r.setMailHasReadTimeAndBroadcastMailExpireTime(mailId, now, expireTime)
		if ec != ErrorCode.Success {
			return ec
		}
		r.prop.SyncUpdateMailExpireTimeAndReadTime(mailId, expireTime.Unix(), now.Unix())
	}
	r.notifyMailRedPoint()
	return ErrorCode.MailSuccess
}

func (r *_Role) setMailExpireTime(mailId string, timeCode uint32, now time.Time, isBroadcast bool) (time.Time, int32) {
	expireSecs := Data.GetMailGeneralConfig(timeCode)
	expireTime := now.Add(time.Second * time.Duration(expireSecs))
	if !isBroadcast {
		if r.ms == nil {
			r.Error("[setMailExpireTime] r.mailService is nil")
			return expireTime, ErrorCode.MailUnknowErr
		}
		// 不是广播邮件
		err := r.ms.SetExpireTime(r.GetOwnerUserID(), mailId, expireTime)
		if err != nil {
			r.Error("[setMailExpireTime] mail service set expire time failed ", err)
			return expireTime, ErrorCode.MailDBError
		}
	}
	return expireTime, ErrorCode.Success
}

// setMailHasReadTime 设置邮件已读时间（是自定义的已读）
func (r *_Role) setMailHasReadTimeAndBroadcastMailExpireTime(mailId string, now time.Time, expireTime time.Time) int32 {
	// todo 会不会有什么隐患。从int64直接到int。会被截断吗？
	// 设置已读时间（邮件排序使用）
	profile_ := map[string]string{
		_MailReadTime:            strconv.Itoa(int(now.Unix())),
		_MailBroadCastExpireTime: strconv.Itoa(int(expireTime.Unix())),
	}
	profile, err := json.Marshal(profile_)
	if err != nil {
		r.Error("[setMailHasReadTime] mail Marshal ExtData failed", err)
		return ErrorCode.MailUnknowErr
	}
	if r.ms == nil {
		r.Error("[setMailHasReadTimeAndBroadcastMailExpireTime] r.mailService is nil")
		return ErrorCode.MailUnknowErr
	}

	err = r.ms.SetExtData(r.GetOwnerUserID(), mailId, profile)
	if err != nil {
		r.Error("[setMailHasReadTime] mail service set extData failed ", err)
		return ErrorCode.MailDBError
	}
	return ErrorCode.Success
}

// deleteMailsHasRead 删除已读邮件
func (r *_Role) deleteMailsHasRead() int32 {
	mails := r.getRoleMails()
	if mails == nil {
		return ErrorCode.MailBoxEmpty
	}
	hasRead := make([]string, 0)
	for _, v := range mails {
		if v.IsRead {
			// 读过了但是有附件没有领取,也不能删除
			if len(v.Attachments) > 0 && !v.IsReceived {
				continue
			}
			hasRead = append(hasRead, v.MailID)
		}
	}
	if len(hasRead) == 0 {
		return ErrorCode.MailSuccess
	}
	if r.ms == nil {
		r.Error("[deleteMailsHasRead] r.mailService is nil")
		return ErrorCode.MailUnknowErr
	}

	err := r.ms.BatchDelete(r.GetOwnerUserID(), hasRead)
	if err != nil {
		r.Errorf("[deleteMailsHasRead] batch delete error: %s", err)
		return ErrorCode.MailDBError
	}
	// 内存删除邮件
	// 给客户端发消息
	for _, v := range hasRead {
		r.prop.SyncMailBoxRemoveMail(v)
	}

	return ErrorCode.Success
}

//-------------------------------------------------------------
// 工具类函数

// removeExpireMail 删除过期邮件
func (r *_Role) removeExpireMail() {
	now := time.Now().Unix()
	for _, mail := range r.prop.Data.MailBox.Mails {
		if mail.ExpireTime < now {
			r.prop.SyncMailBoxRemoveMail(mail.MailID)
		}
	}
}

// convertToProtoMail 将mailapi里的邮件结构转化为自定义的Mail结构
func (r *_Role) convertToProtoMail(mail *mailtypes.Mail) *Proto.Mail {
	if mail == nil {
		return nil
	}
	attachments := make([]*Proto.MailAttachment, 0, len(mail.Attachments))
	for _, ele := range mail.Attachments {
		att := &Proto.MailAttachment{}
		att.Unmarshal(ele.Data)

		attachments = append(attachments, att)
	}

	fromProfile := make(map[string]string)
	err := json.Unmarshal(mail.ExtData, &fromProfile)
	if err != nil {
		r.Error("[convertToProtoMail] mail.ExtData Unmatshal failed")
		// todo 这里失败需要返回吗
	}
	// 之所以这里不放 time.Time类型是因为从[]byte 里我不知道怎么转换成time.Time。可是我知道转换成int 。所以这里存 int了。
	readTime, err := strconv.Atoi(fromProfile[_MailReadTime])
	if err != nil {
		r.Error("[convertToProtoMail] read readTime from mail.ExtData failed")
		// todo 这里失败需要返回吗？ 意味这没有已读时间，不赋值就是0？
	}
	broadCastMailExpireTime, err := strconv.Atoi(fromProfile[_MailBroadCastExpireTime])
	if err != nil {
		r.Error("[convertToProtoMail] read broadcast mail expire Time from mail.ExtData failed")
		// todo 这里失败需要返回吗？ 意味这没有已读时间，不赋值就是0？
	}

	myMail := &Proto.Mail{
		SenderID:    mail.From,
		SenderName:  mail.FromNick,
		ReceiveID:   mail.To,
		IsBroadcast: mail.IsBroadcast,
		IsRead:      mail.IsRead,
		IsReceived:  mail.IsReceived,
		Title:       mail.Title,
		Content:     mail.Body,
		Attachments: attachments,
		SendTime:    mail.SendTime.Unix(),
		MailID:      mail.ID,
		ReadTime:    int64(readTime),
	}
	if myMail.IsBroadcast {
		myMail.ExpireTime = int64(broadCastMailExpireTime)
	} else {
		myMail.ExpireTime = mail.ExpireTime.Unix()
	}
	return myMail
}

// getBatchAttachments 得到批量获取邮件附件
func (r *_Role) getBatchAttachments(mailIDs []string) (int32, []*Proto.OfflineAwardItem, []int32) {
	award := make([]*Proto.OfflineAwardItem, 0)
	for _, mailID := range mailIDs {
		ec, result, errors := r.getOneAttachments(mailID)
		if ec == ErrorCode.MailTimeExpire {
			continue
		}
		if ec != ErrorCode.MailSuccess {
			return ec, nil, errors
		}

		for _, v := range result {
			award = append(award, v)
		}
	}
	if len(award) == 0 {
		return ErrorCode.MailAllTimeExpire, nil, nil
	}
	return ErrorCode.MailSuccess, award, nil
}

// getOneAttachments 得到一封邮件的附件
func (r *_Role) getOneAttachments(mailID string) (int32, []*Proto.OfflineAwardItem, []int32) {
	attachments, ec1 := r.getAttachmentsContext(mailID)
	if attachments == nil {
		return ec1, nil, nil
	}

	// 转化为掉落
	award, items := r.transAttachmentsToDropMaterial(attachments)
	// 判断背包是否能装下所有道具
	bagError, ok := r.CanAddItemList(award)

	// 判断背包能否装下所有固定道具
	myItems := make([]ItemProto.IItem, 0)
	for _, v := range items {
		tmp := ItemProto.CreateIItemByData(v)
		myItems = append(myItems, tmp)
	}
	ok2 := true
	bagError2 := make([]int32, 0)
	key := make(map[int32]int32, 0)
	if len(myItems) != 0 {
		bagError2, ok2 = r.CanAddItemIList(myItems)
	}

	var allBagError []int32
	for _, v := range bagError {
		allBagError = append(allBagError, v)
		key[v] = 1
	}
	for _, v := range bagError2 {
		_, ook := key[v]
		if ook {
			continue
		}
		allBagError = append(allBagError, v)
	}

	if !ok || !ok2 {
		mailError := r.transBagErrorToMailError(allBagError)
		return ErrorCode.MailBagAddErr, nil, mailError
	}

	// 将附件置为已领取
	ec := r.markAttachmentReived(mailID)
	if ec != ErrorCode.MailSuccess {
		return ec, nil, nil
	}

	// 添加物品进入背包
	ret := r.transAwardMaterialToItem(award)
	// 添加固定道具进入背包
	ret2 := r.AddItems(myItems)

	final := make([]*Proto.OfflineAwardItem, 0)
	for _, v := range ret {
		final = append(final, v)
	}
	for _, v := range ret2 {
		final = append(final, &Proto.OfflineAwardItem{ID: v.GetID(), Type: Proto.ItemEnum_Type(v.GetType()), Num: v.GetNum()})
	}

	return ErrorCode.MailSuccess, final, nil

}

// getAttachments 获取附件内容
func (r *_Role) getAttachmentsContext(mailID string) ([]*Proto.MailAttachment, int32) {
	mail := r.getMailFromMailBox(mailID)
	if mail == nil {
		return nil, ErrorCode.MailSuccess
	}
	// 邮件过期
	now := time.Now().Unix()
	if mail.ExpireTime < now {
		r.prop.SyncMailBoxRemoveMail(mail.MailID)
		return nil, ErrorCode.MailTimeExpire
	}

	// 无附件
	if len(mail.Attachments) == 0 {
		return nil, ErrorCode.MailSuccess
	}
	// 已收取
	if mail.IsReceived {
		return nil, ErrorCode.MailSuccess
	}

	return mail.Attachments, ErrorCode.MailSuccess
}

// transAttachmentsToDropMaterial 将附件类型转化为掉落物品类型
func (r *_Role) transAttachmentsToDropMaterial(attachments []*Proto.MailAttachment) ([]*Proto.DropMaterial, []*Proto.Item) {
	var all []*Proto.DropMaterial
	var items []*Proto.Item
	for _, att := range attachments {
		all = append(all, &Proto.DropMaterial{MaterialId: att.ItemID, MaterialType: att.ItemType, MaterialNum: att.ItemNum})
		if att.Data != nil {
			items = append(items, att.Data)
		}
	}
	return all, items
}

// transBagErrorToMailError 将背包错误码转化为邮件错误码
func (r *_Role) transBagErrorToMailError(bagError []int32) []int32 {
	mailErrors := make([]int32, 0)
	for _, v := range bagError {
		switch v {
		case int32(Proto.ContainerEnum_EquipBag):
			mailErrors = append(mailErrors, ErrorCode.EquipBagNotEnoughSpace)
		case int32(Proto.ContainerEnum_SkillBag):
			mailErrors = append(mailErrors, ErrorCode.SkillBagNotEnoughSpace)
		}
	}
	return mailErrors
}

// transAwardMaterialToItem 将奖励材料转化为背包道具（并加入背包）
func (r *_Role) transAwardMaterialToItem(material []*Proto.DropMaterial) []*Proto.OfflineAwardItem {
	itemData := r.AddItemList(material)
	if len(itemData) == 0 {
		return nil
	}
	result := make([]*Proto.OfflineAwardItem, 0)
	for _, val := range itemData {
		if val.GetNum == 0 {
			continue
		}
		awardItem := &Proto.OfflineAwardItem{
			ID:       val.ItemData.Base.ID,
			Type:     val.ItemData.Base.Type,
			Num:      val.GetNum,
			ConfigID: val.ItemData.Base.ConfigID,
		}
		result = append(result, awardItem)
	}
	return result
}

// todo 代码需复检

// MarkAttachmentReived 设置附件为已接收已读取
func (r *_Role) markAttachmentReived(mailId string) int32 {
	mail := r.getMailFromMailBox(mailId)
	if mail == nil {
		return ErrorCode.Failure
	}

	// 附件不存在或者已领取或者已读
	if len(mail.Attachments) == 0 || mail.IsReceived {
		return ErrorCode.MailAthNoExist
	}
	if r.ms == nil {
		r.Error("[markAttachmentReived] r.mailService is nil")
		return ErrorCode.MailUnknowErr
	}

	if !mail.IsRead {
		readErr := r.ms.MarkAsRead(r.GetOwnerUserID(), mailId)
		if readErr != nil {
			r.Errorf("[MarkAttachmentReived] failed to mark mail as read: %s", readErr)
			return ErrorCode.MailDBError
		}
		r.prop.SyncUpdateMailRead(mailId, true)
	}

	received := r.ms.MarkAttachmentsAsReceived(r.GetOwnerUserID(), mailId)
	if received != nil {
		r.Errorf("[MarkAttachmentReived] failed to mark mail attachments as received: %s", received)
		return ErrorCode.MailDBError
	}

	r.prop.SyncUpdateMailReceived(mailId, true)
	now := time.Now()

	expireTime, ec := r.setMailExpireTime(mailId, Data.MailConfigHasReadExpireTime, now, mail.IsBroadcast)
	if ec != ErrorCode.Success {
		return ec
	}
	ec = r.setMailHasReadTimeAndBroadcastMailExpireTime(mailId, now, expireTime)
	if ec != ErrorCode.Success {
		return ec
	}
	r.prop.SyncUpdateMailExpireTimeAndReadTime(mailId, expireTime.Unix(), now.Unix())
	return ErrorCode.MailSuccess
}

func (r *_Role) getMailBox() *Proto.MailBox {
	return r.prop.Data.MailBox
}

func (r *_Role) getRoleMails() []*Proto.Mail {
	if r.prop.Data.MailBox == nil {
		return nil
	}
	return r.prop.Data.MailBox.Mails
}

func (r *_Role) getMailFromMailBox(mailId string) *Proto.Mail {
	box := r.getMailBox()
	if box == nil {
		return nil
	}
	for _, v := range box.Mails {
		if v.MailID == mailId {
			return v
		}
	}
	r.Errorf("[getMailFromMailBox] 邮箱里找不到邮件 %d", mailId)
	return nil
}

// notifyMailRedPoint 增加或者移除邮件红点提示
func (r *_Role) notifyMailRedPoint() {
	notify := r.getRedPointKey(Const.RedPointType_notifyMail, "")
	_, ok := r.prop.Data.RedPointsData[notify]
	canNotify := r.batchCheckMailUnRead()

	// 没有红点又有未读邮件，增加红点
	if !ok && canNotify {
		redPointData := &Proto.RedPointInfo{
			Value:      1,
			CreateTime: time.Now().Unix(),
		}
		r.prop.SyncAddRedPoint(notify, redPointData)
	}

	// 有红点并且没有未读邮件，移除红点
	if ok && !canNotify {
		r.prop.SyncRemoveRedPoint(notify)
	}
}

func (r *_Role) batchCheckMailUnRead() bool {
	if r.prop.Data.MailBox == nil {
		return false
	}

	for _, v := range r.prop.Data.MailBox.Mails {
		if !v.IsRead && !v.IsReceived {
			return true
		}
	}
	return false
}

func (r *_Role) IsOverAttachmentLimit(items []*Proto.MailAttachment) bool {
	var i uint32
	i = 0
	for _, v := range items {
		if v.Data != nil {
			i++
		}
		i++
	}
	limit := Data.GetMailGeneralConfig(Data.MainConfigAttachLen)
	if i > limit {
		return true
	}

	return false
}
