package DistributeLock

import (
	"Cinder/Cache"
	"errors"
	"time"
)

type ILocker interface {
	Lock() error
	Unlock() error
}

var ErrTimeout = errors.New("get lock timeout")

func New(key string, opts ...Option) ILocker {
	if key == "" {
		return nil
	}

	rl := &_RedisLocker{
		key: "lock:" + key,
		opts: Options{
			Expire:   DefaultExpire,
			Timeout:  DefaultTimeout,
			Interval: DefaultInterval,
		},
		isLocking: false,
	}

	for _, o := range opts {
		o(&rl.opts)
	}

	return rl
}

type _RedisLocker struct {
	key       string
	expireAt  time.Time
	opts      Options
	isLocking bool
}

func (rl *_RedisLocker) Lock() error {

	if rl.isLocking {
		panic("relock")
	}

	timeoutAt := time.Now().Add(rl.opts.Timeout)
	for {
		now := time.Now()
		if now.After(timeoutAt) {
			return ErrTimeout
		}

		ret := Cache.RedisDB.SetNX(rl.key, "lock", rl.opts.Expire)
		if ret.Val() {
			rl.expireAt = now.Add(rl.opts.Expire)
			rl.isLocking = true
			return nil
		}

		time.Sleep(rl.opts.Interval)
	}

	return nil
}

func (rl *_RedisLocker) Unlock() error {

	if !rl.isLocking {
		return nil
	}

	if time.Now().After(rl.expireAt) {
		rl.isLocking = false
		return nil
	}

	rl.isLocking = false
	r := Cache.RedisDB.Del(rl.key)
	return r.Err()
}
