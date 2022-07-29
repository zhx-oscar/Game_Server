package tosrvs

import (
	"Cinder/Base/Message"
	"Cinder/Matcher/matchapi/mtypes"
	"encoding/json"

	assert "github.com/arl/assertgo"
	log "github.com/cihub/seelog"
)

// jsonMarshalMsgs 打包成多个 NotifyMsgsToOneSrv
// 结果 json 串不能大于 Message.MaxMessageLen - 1024，不然就分成多段。
func jsonMarshalMsgs(msgs []mtypes.NotifyMsgToOneSrv) [][]byte {
	for i := 1; i <= len(msgs); i++ {
		if bufs, ok := jsonMarshalMsgsN(msgs, i); ok {
			return bufs
		}
	}
	log.Errorf("message is too large") // 单个数据都太大，无法打包
	return nil
}

// jsonMarshalMsgs 打包成n个 NotifyMsgsToOneSrv
// 如果有一个数据长度太大，则返回 (nil, false)
func jsonMarshalMsgsN(msgs []mtypes.NotifyMsgToOneSrv, n int) ([][]byte, bool) {
	assert.True(n > 0)
	totalLen := len(msgs)
	assert.True(n <= totalLen)

	const kMaxBufLen = Message.MaxMessageLen - 1024
	result := make([][]byte, 0, n)
	segmentLen := len(msgs) / n         // 每段长度
	iEnd := totalLen - segmentLen*(n-1) // 当前段结束，第1段为最大段
	for iBegin := 0; iEnd <= totalLen; {
		buf := jsonMarshalMsgsToOne(msgs[iBegin:iEnd])
		if len(buf) > kMaxBufLen {
			return nil, false // 打包后太长了
		}
		result = append(result, buf)

		// 下一段
		iBegin = iEnd
		iEnd += segmentLen
	}
	return result, true
}

// jsonMarshalMsgsToOne 打包成一个 NotifyMsgsToOneSrv
// 出错则返回 nil
func jsonMarshalMsgsToOne(msgs []mtypes.NotifyMsgToOneSrv) []byte {
	msg := mtypes.NotifyMsgsToOneSrv{
		Msgs: msgs,
	}
	buf, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("failed to json marshal: %s", err)
	}
	return buf
}
