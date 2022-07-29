package Security

import (
	"fmt"
	"math/big"
	"testing"
	"time"
)

func Test_key(t *testing.T) {
	dhBase = big.NewInt(3)
	dhPrime, _ = big.NewInt(0).SetString("0x7FFFFFC3", 0)

	a := big.NewInt(123456789)

	println(big.NewInt(0).Exp(dhBase, a, dhPrime).String())

}

func Test_dec(t *testing.T) {

	enc := NewCrypt([]byte("xxyxxyxxxy"))
	//	dec := NewCrypt([]byte("123"))

	src := []byte("wxjwxjwxjwxj,hahahahaha")
	dest1 := make([]byte, 1000)
	dest2 := make([]byte, 1000)

	n, _ := enc.Enc(src, dest1)
	println(string(dest1[:n]))

	n, _ = enc.Enc(dest1[:n], dest2)
	println(string(dest2[:n]))

}

func Test_time(t *testing.T) {

	println(time.Now().UnixNano())

	println(time.Now().UnixNano() / (1e9))

	println(time.Now().Unix())

	tt, err := time.Parse(time.RFC3339, "2016-01-02T15:04:05Z")

	if err != nil {
		println(err.Error())
	}

	println(tt.UnixNano())

	d := time.Now().Sub(tt)

	println(d.Seconds())
	fmt.Printf("%.5f\n", d.Seconds())

}
