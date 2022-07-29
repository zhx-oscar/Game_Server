package Message

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"reflect"
)

type _Encoder struct {
	buf   *bytes.Buffer
	order binary.ByteOrder
}

func newEnc() *_Encoder {
	return &_Encoder{
		buf:   bytes.NewBuffer(nil),
		order: binary.LittleEndian,
	}
}

func (enc *_Encoder) Encode(buf []byte, data interface{}) (int, error) {
	enc.buf = bytes.NewBuffer(buf)
	err := enc.encode(reflect.Indirect(reflect.ValueOf(data)))
	if err != nil {
		return 0, err
	}

	return enc.buf.Len(), nil
}

func (enc *_Encoder) encode(v reflect.Value) error {

	switch v.Kind() {
	case reflect.Array:
		fallthrough
	case reflect.Slice:
		if v.Type().Elem().Kind() == reflect.Uint8 {
			enc.writeBytes(v.Bytes())
		} else {
			l := v.Len()
			enc.writeInt16(int16(l))
			for i := 0; i < l; i++ {
				err := enc.encode(v.Index(i))
				if err != nil {
					return err
				}
			}
		}
	case reflect.Struct:
		l := v.NumField()
		for i := 0; i < l; i++ {
			err := enc.encode(v.Field(i))
			if err != nil {
				return err
			}
		}
	case reflect.Bool:
		if v.Bool() {
			enc.writeInt8(1)
		} else {
			enc.writeInt8(0)
		}
	case reflect.Int8:
		enc.writeInt8(int8(v.Int()))
	case reflect.Int16:
		enc.writeInt16(int16(v.Int()))
	case reflect.Int32:
		enc.writeInt32(int32(v.Int()))
	case reflect.Int64:
		enc.writeInt64(v.Int())

	case reflect.Uint8:
		enc.writeUInt8(uint8(v.Uint()))
	case reflect.Uint16:
		enc.writeUInt16(uint16(v.Uint()))
	case reflect.Uint32:
		enc.writeUInt32(uint32(v.Uint()))
	case reflect.Uint64:
		enc.writeUInt64(v.Uint())

	case reflect.Float32:
		enc.writeFloat32(float32(v.Float()))
	case reflect.Float64:
		enc.writeFloat64(v.Float())
	case reflect.String:
		enc.writeString(v.String())
	default:
		return errors.New(fmt.Sprintf("%s , %d", "not support data kind ", v.Kind()))
	}

	return nil
}

func (enc *_Encoder) writeInt8(x int8) {
	_ = enc.buf.WriteByte(byte(x))
}

func (enc *_Encoder) writeUInt8(b uint8) {
	_ = enc.buf.WriteByte(b)
}

func (enc *_Encoder) writeInt16(b int16) {
	buf := make([]byte, 2)
	enc.order.PutUint16(buf, uint16(b))
	enc.buf.Write(buf)
}

func (enc *_Encoder) writeInt32(b int32) {
	buf := make([]byte, 4)
	enc.order.PutUint32(buf, uint32(b))
	enc.buf.Write(buf)
}

func (enc *_Encoder) writeInt64(b int64) {
	buf := make([]byte, 8)
	enc.order.PutUint64(buf, uint64(b))
	enc.buf.Write(buf)
}

func (enc *_Encoder) writeUInt16(b uint16) {
	buf := make([]byte, 2)
	enc.order.PutUint16(buf, b)
	enc.buf.Write(buf)
}

func (enc *_Encoder) writeUInt32(b uint32) {
	buf := make([]byte, 4)
	enc.order.PutUint32(buf, b)
	enc.buf.Write(buf)
}

func (enc *_Encoder) writeUInt64(b uint64) {
	buf := make([]byte, 8)
	enc.order.PutUint64(buf, b)
	enc.buf.Write(buf)
}

func (enc *_Encoder) writeBytes(bs []byte) {
	enc.writeUInt32(uint32(len(bs)))
	enc.buf.Write(bs)
}

func (enc *_Encoder) writeString(s string) {
	enc.writeUInt32(uint32(len(s)))
	enc.buf.Write([]byte(s))
}

func (enc *_Encoder) writeFloat32(f float32) {
	enc.writeUInt32(math.Float32bits(f))
}

func (enc *_Encoder) writeFloat64(f float64) {
	enc.writeUInt64(math.Float64bits(f))
}
