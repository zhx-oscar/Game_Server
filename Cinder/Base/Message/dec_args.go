package Message

import (
	"Cinder/Base/ProtoDef"
	"Cinder/Base/linemath"
	"bytes"
	"encoding/binary"
	"errors"
	"math"

	"github.com/golang/protobuf/proto"
)

type _ArgsDecoder struct {
	buf   *bytes.Buffer
	order binary.ByteOrder
}

func newArgsDec() *_ArgsDecoder {
	return &_ArgsDecoder{
		buf:   nil,
		order: binary.LittleEndian,
	}
}

func (dec *_ArgsDecoder) Decode(data []byte) ([]interface{}, error) {
	dec.buf = bytes.NewBuffer(data)
	return dec.decode()
}

func (dec *_ArgsDecoder) decode() ([]interface{}, error) {

	ret := make([]interface{}, 0, 5)

	for {

		var typ int8
		var err error
		var v interface{}

		if typ, err = dec.readInt8(); err != nil {
			break
		}

		switch typ {
		case Arg_Bool:
			v, err = dec.readInt8()
			if v.(int8) == 1 {
				v = true
			} else {
				v = false
			}
		case Arg_Int8:
			v, err = dec.readInt8()
		case Arg_UInt8:
			v, err = dec.readUInt8()
		case Arg_Int16:
			v, err = dec.readInt16()
		case Arg_UInt16:
			v, err = dec.readUInt16()
		case Arg_Int32:
			v, err = dec.readInt32()
		case Arg_UInt32:
			v, err = dec.readUInt32()
		case Arg_Int64:
			v, err = dec.readInt64()
		case Arg_UInt64:
			v, err = dec.readUInt64()
		case Arg_Float32:
			v, err = dec.readFloat32()
		case Arg_Float64:
			v, err = dec.readFloat64()
		case Arg_String:
			v, err = dec.readString()
		case Arg_ByteArray:
			v, err = dec.readBytes()
		case Arg_UInt32Array:
			v, err = dec.readUInt32Array()
		case Arg_Vector2:
			o := linemath.NewVector2()
			o.X, err = dec.readFloat32()
			o.Y, err = dec.readFloat32()
			v = o
		case Arg_Vector3:
			o := linemath.NewVector3()
			o.X, err = dec.readFloat32()
			o.Y, err = dec.readFloat32()
			o.Z, err = dec.readFloat32()
			v = o
		case Arg_Quaternion:
			o := linemath.NewQuaternion()
			o.X, err = dec.readFloat32()
			o.Y, err = dec.readFloat32()
			o.Z, err = dec.readFloat32()
			o.W, err = dec.readFloat32()
			v = o
		case Arg_Proto:
			id, err := dec.readInt16()
			if err != nil {
				return nil, errors.New("wrong format")
			}

			v, err = ProtoDef.GetProtoMessageByID(int(id))
			if err != nil {
				return nil, errors.New("no proto id")
			}

			buf, err := dec.readBytes()
			if err != nil {
				return nil, errors.New("get proto wrong")
			}

			err = proto.Unmarshal(buf, v.(proto.Message))

		default:
			return nil, errors.New("no support type")
		}

		if err != nil {
			return nil, err
		}

		ret = append(ret, v)
	}

	return ret, nil
}

func (dec *_ArgsDecoder) readInt8() (int8, error) {
	b, err := dec.buf.ReadByte()
	return int8(b), err
}

func (dec *_ArgsDecoder) readUInt8() (uint8, error) {
	return dec.buf.ReadByte()
}

func (dec *_ArgsDecoder) readInt16() (int16, error) {
	buf := make([]byte, 2)

	n, err := dec.buf.Read(buf)
	if err != nil {
		return 0, err
	}

	if n != 2 {
		return 0, errors.New("read buf error")
	}

	return int16(dec.order.Uint16(buf)), nil
}

func (dec *_ArgsDecoder) readInt32() (int32, error) {
	buf := make([]byte, 4)

	n, err := dec.buf.Read(buf)
	if err != nil {
		return 0, err
	}

	if n != 4 {
		return 0, errors.New("read buf error")
	}

	return int32(dec.order.Uint32(buf)), nil
}

func (dec *_ArgsDecoder) readInt64() (int64, error) {
	buf := make([]byte, 8)

	n, err := dec.buf.Read(buf)
	if err != nil {
		return 0, err
	}

	if n != 8 {
		return 0, errors.New("read buf error")
	}

	return int64(dec.order.Uint64(buf)), nil
}

func (dec *_ArgsDecoder) readUInt16() (uint16, error) {
	buf := make([]byte, 2)

	n, err := dec.buf.Read(buf)
	if err != nil {
		return 0, err
	}

	if n != 2 {
		return 0, errors.New("read buf error")
	}

	return dec.order.Uint16(buf), nil
}

func (dec *_ArgsDecoder) readUInt32() (uint32, error) {
	buf := make([]byte, 4)

	n, err := dec.buf.Read(buf)
	if err != nil {
		return 0, err
	}

	if n != 4 {
		return 0, errors.New("read buf error")
	}

	return dec.order.Uint32(buf), nil
}

func (dec *_ArgsDecoder) readUInt64() (uint64, error) {
	buf := make([]byte, 8)

	n, err := dec.buf.Read(buf)
	if err != nil {
		return 0, err
	}

	if n != 8 {
		return 0, errors.New("read buf error")
	}

	return dec.order.Uint64(buf), nil
}

func (dec *_ArgsDecoder) readBytes() ([]byte, error) {

	n, err := dec.readUInt32()
	if err != nil {
		return nil, err
	}

	if n < 0 {
		return nil, errors.New("format error")
	}

	buf := make([]byte, n)
	rn, err := dec.buf.Read(buf)
	if err != nil || rn != int(n) {
		return nil, errors.New("read buf error")
	}

	return buf, nil

}

func (dec *_ArgsDecoder) readString() (string, error) {
	bs, err := dec.readBytes()
	if err != nil {
		return "", err
	}

	return string(bs), err
}

func (dec *_ArgsDecoder) readFloat32() (float32, error) {

	d, err := dec.readUInt32()
	if err != nil {
		return 0, err
	}

	return math.Float32frombits(d), nil
}

func (dec *_ArgsDecoder) readFloat64() (float64, error) {

	d, err := dec.readUInt64()
	if err != nil {
		return 0, err
	}

	return math.Float64frombits(d), nil
}

func (dec *_ArgsDecoder) readUInt32Array() ([]uint32, error) {
	n, err := dec.readUInt16()
	if err != nil {
		return nil, err
	}

	if n < 0 {
		return nil, errors.New("format error")
	}
	res := make([]uint32, n)
	for i := 0; i < int(n); i++ {
		num, err := dec.readUInt32()
		if err != nil {
			return nil, err
		}
		res[i] = num
	}
	return res, nil
}
