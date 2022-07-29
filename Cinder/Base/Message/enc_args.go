package Message

import (
	"Cinder/Base/ProtoDef"
	"Cinder/Base/linemath"
	"bytes"
	"encoding/binary"
	"errors"
	"math"
	"reflect"

	"github.com/golang/protobuf/proto"
)

const (
	Arg_Int8        = 2
	Arg_UInt8       = 3
	Arg_Int16       = 4
	Arg_UInt16      = 5
	Arg_Int32       = 6
	Arg_UInt32      = 7
	Arg_Int64       = 8
	Arg_UInt64      = 9
	Arg_String      = 10
	Arg_ByteArray   = 11
	Arg_Proto       = 12
	Arg_Bool        = 14
	Arg_Float32     = 15
	Arg_Float64     = 16
	Arg_Vector2     = 17
	Arg_Vector3     = 18
	Arg_Quaternion  = 19
	Arg_UInt32Array = 20
)

type _ArgsEncoder struct {
	buf   *bytes.Buffer
	order binary.ByteOrder
}

func newArgsEnc() *_ArgsEncoder {
	return &_ArgsEncoder{
		buf:   bytes.NewBuffer(nil),
		order: binary.LittleEndian,
	}
}

func (enc *_ArgsEncoder) Encode(buf []byte, args []interface{}) ([]byte, error) {
	enc.buf = bytes.NewBuffer(buf)
	err := enc.encode(args)
	if err != nil {
		return nil, err
	}

	return enc.buf.Bytes(), nil
}

func (enc *_ArgsEncoder) encode(args []interface{}) error {

	for _, arg := range args {

		switch arg.(type) {
		case bool:
			enc.writeInt8(Arg_Bool)
			if arg.(bool) {
				enc.writeInt8(1)
			} else {
				enc.writeInt8(0)
			}

		case int8:
			enc.writeInt8(Arg_Int8)
			enc.writeInt8(arg.(int8))

		case uint8:
			enc.writeInt8(Arg_UInt8)
			enc.writeUInt8(arg.(uint8))

		case int16:
			enc.writeInt8(Arg_Int16)
			enc.writeInt16(arg.(int16))

		case uint16:
			enc.writeInt8(Arg_UInt16)
			enc.writeUInt16(arg.(uint16))

		case int32:
			enc.writeInt8(Arg_Int32)
			enc.writeInt32(arg.(int32))

		case uint32:
			enc.writeInt8(Arg_UInt32)
			enc.writeUInt32(arg.(uint32))

		case int64:
			enc.writeInt8(Arg_Int64)
			enc.writeInt64(arg.(int64))

		case uint64:
			enc.writeInt8(Arg_UInt64)
			enc.writeUInt64(arg.(uint64))

		case float32:
			enc.writeInt8(Arg_Float32)
			enc.writeFloat32(arg.(float32))

		case float64:
			enc.writeInt8(Arg_Float64)
			enc.writeFloat64(arg.(float64))

		case string:
			enc.writeInt8(Arg_String)
			enc.writeString(arg.(string))

		case []byte:
			enc.writeInt8(Arg_ByteArray)
			enc.writeBytes(arg.([]byte))
		case *linemath.Vector2:
			o := arg.(*linemath.Vector2)

			enc.writeInt8(Arg_Vector2)
			enc.writeFloat32(o.X)
			enc.writeFloat32(o.Y)
		case *linemath.Vector3:
			o := arg.(*linemath.Vector3)

			enc.writeInt8(Arg_Vector3)
			enc.writeFloat32(o.X)
			enc.writeFloat32(o.Y)
			enc.writeFloat32(o.Z)
		case *linemath.Quaternion:
			o := arg.(*linemath.Quaternion)

			enc.writeInt8(Arg_Quaternion)
			enc.writeFloat32(o.X)
			enc.writeFloat32(o.Y)
			enc.writeFloat32(o.Z)
			enc.writeFloat32(o.W)
		case []uint32:
			enc.writeInt8(Arg_UInt32Array)
			enc.writeUInt16(uint16(len(arg.([]uint32))))
			for i := 0; i < len(arg.([]uint32)); i++ {
				enc.writeUInt32(arg.([]uint32)[i])
			}
		default:

			pm, ok := arg.(proto.Message)
			if !ok {
				return errors.New("no support")
			}

			pid, err := ProtoDef.GetIDByName(reflect.TypeOf(arg).Elem().Name())
			if err != nil {
				return errors.New("no protodef ")
			}

			enc.writeInt8(Arg_Proto)
			enc.writeInt16(int16(pid))

			if !reflect.ValueOf(pm).IsNil() {
				buf, err := proto.Marshal(pm)
				if err != nil {
					return err
				}
				enc.writeBytes(buf)
			} else {
				enc.writeUInt32(0)
			}
		}
	}

	return nil
}

func (enc *_ArgsEncoder) writeInt8(x int8) {
	_ = enc.buf.WriteByte(byte(x))
}

func (enc *_ArgsEncoder) writeUInt8(b uint8) {
	_ = enc.buf.WriteByte(b)
}

func (enc *_ArgsEncoder) writeInt16(b int16) {
	buf := make([]byte, 2)
	enc.order.PutUint16(buf, uint16(b))
	enc.buf.Write(buf)
}

func (enc *_ArgsEncoder) writeInt32(b int32) {
	buf := make([]byte, 4)
	enc.order.PutUint32(buf, uint32(b))
	enc.buf.Write(buf)
}

func (enc *_ArgsEncoder) writeInt64(b int64) {
	buf := make([]byte, 8)
	enc.order.PutUint64(buf, uint64(b))
	enc.buf.Write(buf)
}

func (enc *_ArgsEncoder) writeUInt16(b uint16) {
	buf := make([]byte, 2)
	enc.order.PutUint16(buf, b)
	enc.buf.Write(buf)
}

func (enc *_ArgsEncoder) writeUInt32(b uint32) {
	buf := make([]byte, 4)
	enc.order.PutUint32(buf, b)
	enc.buf.Write(buf)
}

func (enc *_ArgsEncoder) writeUInt64(b uint64) {
	buf := make([]byte, 8)
	enc.order.PutUint64(buf, b)
	enc.buf.Write(buf)
}

func (enc *_ArgsEncoder) writeBytes(bs []byte) {
	enc.writeUInt32(uint32(len(bs)))
	enc.buf.Write(bs)
}

func (enc *_ArgsEncoder) writeString(s string) {
	enc.writeUInt32(uint32(len(s)))
	enc.buf.Write([]byte(s))
}

func (enc *_ArgsEncoder) writeFloat32(f float32) {
	enc.writeUInt32(math.Float32bits(f))
}

func (enc *_ArgsEncoder) writeFloat64(f float64) {
	enc.writeUInt64(math.Float64bits(f))
}

func (enc *_ArgsEncoder) writeUInt32Slice(bs []uint32) {
	for i := 0; i < len(bs); i++ {
		enc.writeUInt32(bs[i])
	}
}
