package Security

import (
	"crypto/rc4"
	"errors"
)

type ICrypt interface {
	Enc(data []byte, tempBuf []byte) (int, error)
	Dec(data []byte, tempBuf []byte) (int, error)
}

func NewCrypt(key []byte) ICrypt {
	r := &_Crypt{
		key: key,
	}
	return r
}

type _Crypt struct {
	key []byte
}

func (c *_Crypt) Enc(data []byte, tempBuf []byte) (int, error) {

	if c.key == nil {
		return 0, errors.New("no crypt")
	}

	if cap(tempBuf) < len(data) {
		return 0, errors.New("no enough space")
	}

	rc, _ := rc4.NewCipher(c.key)
	rc.XORKeyStream(tempBuf, data)
	return len(data), nil
}

func (c *_Crypt) Dec(data []byte, tempBuf []byte) (int, error) {
	return c.Enc(data, tempBuf)
}
