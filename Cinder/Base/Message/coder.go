package Message

import (
	"Cinder/Base/Security"
	"Cinder/Base/Util"
	"errors"

	log "github.com/cihub/seelog"
	"github.com/golang/snappy"
)

var (
	ErrMessageInvalid  = errors.New("message invalid")
	ErrMessageTooLarge = errors.New("message too large")
	ErrMessageLength   = errors.New("message length error")
	ErrCryptNotExisted = errors.New("crypt not existed")
)

const CompressSize = 20 * 1024

func isCompression(flag byte) bool {
	return flag&(1) != 0
}

func isEncryption(flag byte) bool {
	return flag&(1<<1) != 0
}

func combineFlag(compressionFlag bool, encryptionFlag bool) byte {
	var flag byte = 0

	if compressionFlag {
		flag |= 1
	}

	if encryptionFlag {
		flag |= 1 << 1
	}

	return flag
}

func compress(data []byte, tempBuf []byte) int {
	ret := snappy.Encode(tempBuf, data)
	return len(ret)
}

func decompress(data []byte, tempBuf []byte) (int, error) {
	ret, err := snappy.Decode(tempBuf, data)
	if err != nil {
		return 0, err
	}
	return len(ret), nil
}

func pack(msg IMessage, crypt Security.ICrypt, msgNo uint32, tempBuf []byte) ([]byte, error) {
	if msg == nil {
		return nil, ErrMessageInvalid
	}

	msgLen, err := CalcMessageSize(msg)
	if err != nil {
		return nil, err
	}

	var buf []byte
	if tempBuf == nil {
		// 最终消息缓存区, 可能存在消息很大且不能被压缩的情况, 所以需要比原始消息更大的缓存
		buf = make([]byte, GetMaxMessageEncodedSizeWithLen(msgLen))
	} else {
		buf = tempBuf
	}

	// 如果消息需要压缩, 则一开始先把消息序列化到一块临时内存中
	var encBuf []byte
	var compressionFlag bool
	if NeedCompress(msgLen) {
		encBuf = Util.Get(msgLen)
		defer Util.Put(encBuf)
		compressionFlag = true
	} else {
		encBuf = buf[HeadLen:]
	}

	enc := newEnc()
	l, err := enc.Encode(encBuf[0:0], msg)
	if err != nil {
		return nil, err
	}

	// 压缩
	if compressionFlag {
		cl := compress(encBuf[:l], buf[HeadLen:])
		buf = buf[:HeadLen+cl]
	} else {
		buf = buf[:HeadLen+l]
	}

	// 加密
	encryptionFlag := false
	if crypt != nil {
		encryptionFlag = true
		_, err = crypt.Enc(buf[HeadLen:], buf[HeadLen:])
		if err != nil {
			return nil, err
		}
	}

	flag := combineFlag(compressionFlag, encryptionFlag)
	l = len(buf[HeadLen:])
	if l > MaxMessageLen {
		return nil, ErrMessageTooLarge
	}

	buf[0] = byte(l)
	buf[1] = byte(l >> 8)
	buf[2] = byte(l >> 16)
	buf[3] = flag

	id := msg.GetID()
	buf[4] = byte(id)
	buf[5] = byte(id >> 8)

	if msgNo != 0 {
		buf[6] = byte(msgNo)
		buf[7] = byte(msgNo >> 8)
		buf[8] = byte(msgNo >> 16)
		buf[9] = byte(msgNo >> 24)
	}

	return buf, nil
}

func unpack(buf []byte, crypt Security.ICrypt) (IMessage, uint32, error) {

	l := int(buf[0]) | int(buf[1])<<8 | int(buf[2])<<16
	flag := buf[3]
	id := uint16(buf[4]) | uint16(buf[5])<<8
	msgNo := uint32(buf[6]) | uint32(buf[7])<<8 | uint32(buf[8])<<16 | uint32(buf[9])<<24

	if l+HeadLen != len(buf) {
		return nil, 0, ErrMessageLength
	}

	////////////////////////////////

	m, err := def.fetchMessage(id)
	if err != nil {
		return nil, 0, err
	}

	// 解密
	if isEncryption(flag) {
		if crypt == nil {
			return nil, 0, ErrCryptNotExisted
		}

		if _, err = crypt.Dec(buf[HeadLen:], buf[HeadLen:]); err != nil {
			return nil, 0, err
		}
	}

	// 压缩
	tbuf := buf[HeadLen:]
	if isCompression(flag) {
		var decodeLen int
		decodeLen, err = snappy.DecodedLen(tbuf)
		if err != nil {
			return nil, 0, err
		}

		dbuf := Util.Get(decodeLen)
		defer Util.Put(dbuf)

		var dl int
		dl, err = decompress(tbuf, dbuf)
		if err != nil {
			return nil, 0, err
		}
		tbuf = dbuf[:dl]
	}

	dec := newDec()
	err = dec.Decode(tbuf, m)
	if err != nil {
		return nil, 0, err
	}

	return m, msgNo, err
}

func packArgs(args []interface{}) []byte {

	var tempBuf = make([]byte, 4*1024)

	enc := newArgsEnc()
	mb, err := enc.Encode(tempBuf[0:0], args)
	if err != nil {
		log.Error("pack args err ", err, " Args: ", args)
		return nil
	}

	ret := make([]byte, len(mb))
	copy(ret, mb)

	return ret
}

func unPackArgs(buf []byte) ([]interface{}, error) {

	dec := newArgsDec()
	v, err := dec.Decode(buf)
	if err != nil {
		return nil, err
	}
	return v, nil
}
