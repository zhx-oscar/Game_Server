package Security

import (
	"math"
	"math/big"
	"math/rand"
	"sync"
	"time"
)

var (
	rnd        = rand.New(rand.NewSource(time.Now().UnixNano()))
	dhBase     = big.NewInt(3)
	dhPrime, _ = big.NewInt(0).SetString("0x7FFFFFC3", 0)
	maxNum     = big.NewInt(math.MaxInt64)

	randLock = sync.Mutex{}
)

func Pair() (privateKey *big.Int, publicKey *big.Int) {
	randLock.Lock()
	s := big.NewInt(0).Rand(rnd, maxNum)
	randLock.Unlock()
	//s := big.NewInt(rnd.Int63())
	p := big.NewInt(0).Exp(dhBase, s, dhPrime)
	return s, p
}

func Key(privateKey *big.Int, otherKey *big.Int) *big.Int {
	key := big.NewInt(0).Exp(otherKey, privateKey, dhPrime)
	return key
}
