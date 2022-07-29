package Message

import (
	"Cinder/Base/Security"
	"errors"
	"fmt"
	"github.com/golang/snappy"
	"reflect"
)

const (
	MaxMessageLen = 0xffffff
	HeadLen       = 10
)

type IMessage interface {
	GetID() uint16
}

func Pack(msg IMessage) ([]byte, error) {
	return pack(msg, nil, 0, nil)
}

func PackWithBuf(msg IMessage, buf []byte) ([]byte, error) {
	return pack(msg, nil, 0, buf)
}

func PackWithSKey(msg IMessage, crypt Security.ICrypt, msgNo uint32) ([]byte, error) {
	return pack(msg, crypt, msgNo, nil)
}

func Unpack(buf []byte) (IMessage, error) {
	msg, _, err := unpack(buf, nil)
	return msg, err
}

func UnpackWithSKey(buf []byte, crypt Security.ICrypt) (IMessage, uint32, error) {
	return unpack(buf, crypt)
}

func PackArgs(args ...interface{}) []byte {
	return packArgs(args)
}

func UnPackArgs(buf []byte) ([]interface{}, error) {
	return unPackArgs(buf)
}

func CalcMessageSize(msg IMessage) (int, error) {
	return CalcSize(reflect.Indirect(reflect.ValueOf(msg)))
}

func CalcSize(v reflect.Value) (int, error) {
	switch v.Kind() {
	case reflect.Array:
		fallthrough
	case reflect.Slice:
		if v.Type().Elem().Kind() == reflect.Uint8 {
			return 4 + len(v.Bytes()), nil
		} else {
			l := v.Len()
			sum := 2
			for i := 0; i < l; i++ {
				n, err := CalcSize(v.Index(i))
				if err != nil {
					return 0, err
				}
				sum += n
			}
			return sum, nil
		}

	case reflect.Struct:
		l := v.NumField()
		sum := 0
		for i := 0; i < l; i++ {
			n, err := CalcSize(v.Field(i))
			if err != nil {
				return 0, err
			}
			sum += n
		}
		return sum, nil

	case reflect.Int8, reflect.Uint8, reflect.Bool:
		return 1, nil
	case reflect.Int16, reflect.Uint16:
		return 2, nil
	case reflect.Int32, reflect.Uint32, reflect.Float32:
		return 4, nil
	case reflect.Int64, reflect.Uint64, reflect.Float64:
		return 8, nil
	case reflect.String:
		return 4 + len(v.String()), nil

	default:
		return 0, errors.New(fmt.Sprintf("%s , %d", "not support data kind ", v.Kind()))
	}
}

func GetMaxMessageEncodedSize(message IMessage) (int, error) {
	msgLen, err := CalcMessageSize(message)
	if err != nil {
		return 0, err
	}

	return GetMaxMessageEncodedSizeWithLen(msgLen), nil
}

func GetMaxMessageEncodedSizeWithLen(msgLen int) int {
	if NeedCompress(msgLen) {
		return snappy.MaxEncodedLen(msgLen) + HeadLen
	} else {
		return msgLen + HeadLen
	}
}

func NeedCompress(size int) bool {
	return size > CompressSize
}
