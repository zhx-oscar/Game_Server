package Message

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"reflect"
)

type _Decoder struct {
	buf   *bytes.Buffer
	order binary.ByteOrder
}

func newDec() *_Decoder {
	return &_Decoder{
		buf:   nil,
		order: binary.LittleEndian,
	}
}

func (dec *_Decoder) Decode(data []byte, obj interface{}) error {
	dec.buf = bytes.NewBuffer(data)

	err := dec.decode(reflect.Indirect(reflect.ValueOf(obj)))
	if err != nil {
		return err
	}

	return nil
}

func (dec *_Decoder) decode(v reflect.Value) error {

	switch v.Kind() {
	case reflect.Array:
		fallthrough
	case reflect.Slice:

		if v.Type().Elem().Kind() == reflect.Uint8 {
			d, err := dec.readBytes()
			if err != nil {
				return err
			}
			v.SetBytes(d)
		} else {

			l, err := dec.readInt16()
			if err != nil {
				return err
			}

			t := v.Type()
			tt := t.Elem()

			for i := 0; i < int(l); i++ {

				ev := reflect.New(tt)
				err = dec.decode(ev.Elem())
				if err != nil {
					return err
				}

				v.Set(reflect.Append(v, ev.Elem()))
			}
		}

	case reflect.Struct:
		l := v.NumField()
		for i := 0; i < l; i++ {
			err := dec.decode(v.Field(i))
			if err != nil {
				return err
			}
		}
	case reflect.Bool:
		d, err := dec.readInt8()
		if err != nil {
			return err
		}

		if d == 1 {
			v.SetBool(true)
		} else if d == 0 {
			v.SetBool(false)
		} else {
			return errors.New("bool error")
		}

	case reflect.Int8:
		d, err := dec.readInt8()
		if err != nil {
			return err
		}
		v.SetInt(int64(d))
	case reflect.Int16:
		d, err := dec.readInt16()
		if err != nil {
			return err
		}
		v.SetInt(int64(d))
	case reflect.Int32:
		d, err := dec.readInt32()
		if err != nil {
			return err
		}
		v.SetInt(int64(d))
	case reflect.Int64:
		d, err := dec.readInt64()
		if err != nil {
			return err
		}
		v.SetInt(d)

	case reflect.Uint8:
		d, err := dec.readUInt8()
		if err != nil {
			return err
		}
		v.SetUint(uint64(d))
	case reflect.Uint16:
		d, err := dec.readUInt16()
		if err != nil {
			return err
		}
		v.SetUint(uint64(d))
	case reflect.Uint32:
		d, err := dec.readUInt32()
		if err != nil {
			return err
		}
		v.SetUint(uint64(d))
	case reflect.Uint64:
		d, err := dec.readUInt64()
		if err != nil {
			return err
		}
		v.SetUint(d)

	case reflect.Float32:
		d, err := dec.readFloat32()
		if err != nil {
			return err
		}
		v.SetFloat(float64(d))
	case reflect.Float64:
		d, err := dec.readFloat64()
		if err != nil {
			return err
		}
		v.SetFloat(d)
	case reflect.String:
		s, err := dec.readString()
		if err != nil {
			return err
		}
		v.SetString(s)

	default:
		return errors.New(fmt.Sprintf("%s , %d", "not support data kind ", v.Kind()))
	}

	return nil
}

func (dec *_Decoder) readInt8() (int8, error) {
	b, err := dec.buf.ReadByte()
	return int8(b), err
}

func (dec *_Decoder) readUInt8() (uint8, error) {
	return dec.buf.ReadByte()
}

func (dec *_Decoder) readInt16() (int16, error) {
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

func (dec *_Decoder) readInt32() (int32, error) {
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

func (dec *_Decoder) readInt64() (int64, error) {
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

func (dec *_Decoder) readUInt16() (uint16, error) {
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

func (dec *_Decoder) readUInt32() (uint32, error) {
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

func (dec *_Decoder) readUInt64() (uint64, error) {
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

func (dec *_Decoder) readBytes() ([]byte, error) {

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

func (dec *_Decoder) readString() (string, error) {
	bs, err := dec.readBytes()
	if err != nil {
		return "", err
	}

	return string(bs), err
}

func (dec *_Decoder) readFloat32() (float32, error) {

	d, err := dec.readUInt32()
	if err != nil {
		return 0, err
	}

	return math.Float32frombits(d), nil
}

func (dec *_Decoder) readFloat64() (float64, error) {

	d, err := dec.readUInt64()
	if err != nil {
		return 0, err
	}

	return math.Float64frombits(d), nil
}
