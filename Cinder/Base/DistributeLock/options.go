package DistributeLock

import "time"

var (
	DefaultExpire   = 5 * time.Second
	DefaultTimeout  = 50 * time.Second
	DefaultInterval = 10 * time.Millisecond
)

type Options struct {
	Expire   time.Duration
	Timeout  time.Duration
	Interval time.Duration
}

type Option func(options *Options)

func Expire(expire time.Duration) Option {
	return func(o *Options) {
		o.Expire = expire
	}
}

func Timeout(timeout time.Duration) Option {
	return func(o *Options) {
		o.Timeout = timeout
	}
}

func Interval(interval time.Duration) Option {
	return func(o *Options) {
		o.Interval = interval
	}
}
