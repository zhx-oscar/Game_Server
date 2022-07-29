package Message

import (
	"Cinder/Base/Security"
	"fmt"
	"reflect"
	"testing"
)

type _TestStruct struct {
	U8   uint8
	U16  uint16
	I32  int32
	F64  float64
	Str  string
	Data []byte
	U16s []uint16
}

func (ts *_TestStruct) GetID() uint16 {
	return 100
}

func BenchmarkPack(b *testing.B) {
	def.addDef(&_TestStruct{})

	ts := &_TestStruct{}
	ts.Str = "Hello"
	ts.Data = []byte("Hello")
	ts.U16s = []uint16{1, 2, 3}
	//crypt := Security.NewCrypt([]byte("HelloWorld"))

	b.ResetTimer()
	b.ReportAllocs()
	var err error
	for i := 0; i < b.N; i++ {
		if _, err = pack(ts, nil, 12, nil); err != nil {
			b.Fatal(err)
		}
	}
}

func TestEncode(t *testing.T) {
	def.addDef(&_TestStruct{})

	ts := &_TestStruct{}
	ts.Str = "Hello"
	ts.Data = []byte("Hello")
	ts.U16s = []uint16{1, 2, 3}
	n, err := CalcSize(reflect.Indirect(reflect.ValueOf(ts)))
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("MsgSize:", n)

	crypt := Security.NewCrypt([]byte("HelloWorld"))

	data, err := pack(ts, crypt, 12, nil)
	fmt.Println("Data:", data)

	tss, no, err := unpack(data, crypt)
	fmt.Println(tss, no)

	if !reflect.DeepEqual(ts, tss) {
		t.Fatal("mismatch")
	}
}
