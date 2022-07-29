package Config

import (
	_ "Cinder/Base/ServerConfig"
	"context"
	"errors"
	"github.com/spf13/viper"
	"go.etcd.io/etcd/clientv3"
	"sync"
	"time"

	log "github.com/cihub/seelog"
)

func init() {
	Inst = newConfig(viper.GetString("ETCD.Addr"))
}

func newConfig(addr string) IConfig {

	c, _ := clientv3.New(
		clientv3.Config{
			Endpoints:   []string{addr},
			DialTimeout: 5 * time.Second,
		})

	return &_EtcdClient{
		cli:         c,
		watchHandle: 1,
		cancelMap:   make(map[int]context.CancelFunc),
	}
}

type _EtcdClient struct {
	cli *clientv3.Client

	watchHandle int
	cancelMap   map[int]context.CancelFunc
	mtx         sync.Mutex
}

func (c *_EtcdClient) SetValue(key string, value string) error {
	_, err := c.cli.Put(context.Background(), key, value)
	return err
}

func (c *_EtcdClient) SetValueAndOvertime(key string, value string, ttl int64) error {
	r, err := c.cli.Grant(context.Background(), ttl)
	if err != nil {
		return err
	}

	_, err = c.cli.Put(context.Background(), key, value, clientv3.WithLease(r.ID))
	return err
}

func (c *_EtcdClient) SetValueAndKeepAlive(key string, value string) error {

	r, err := c.cli.Grant(context.Background(), 60)
	if err != nil {
		return err
	}

	_, err = c.cli.Put(context.Background(), key, value, clientv3.WithLease(r.ID))
	if err != nil {
		return err
	}

	reps, err := c.cli.KeepAlive(context.Background(), r.ID)
	if err != nil {
		return err
	}

	go func() {
		for {
			if _, ok := <-reps; !ok {
				log.Error("SetValueAndKeepAlive ETCD failed")
				return
			}
		}
	}()

	return nil
}

func (c *_EtcdClient) GetValue(key string) (string, error) {

	r, err := c.cli.Get(context.Background(), key)

	if err != nil {
		return "", err
	}

	if len(r.Kvs) != 1 {
		return "", errors.New("GetValue no key value pair return" + key)
	}

	return string(r.Kvs[0].Value), nil
}

func (c *_EtcdClient) GetValuesByPrefix(prefix string) ([]string, []string, error) {

	r, err := c.cli.Get(context.Background(), prefix, clientv3.WithPrefix())

	if err != nil {
		return nil, nil, err
	}

	keys := make([]string, 0, 5)
	values := make([]string, 0, 5)

	for _, ev := range r.Kvs {
		keys = append(keys, string(ev.Key))
		values = append(values, string(ev.Value))
	}

	return keys, values, nil
}

func (c *_EtcdClient) WatchKey(key string, cb func(int, string, string)) (int, error) {

	ctx, cancelFunc := context.WithCancel(context.Background())

	r := c.cli.Watch(ctx, key)

	var watchHandle int
	c.mtx.Lock()

	watchHandle = c.watchHandle
	c.cancelMap[c.watchHandle] = cancelFunc
	c.watchHandle++

	c.mtx.Unlock()

	go func() {
		for w := range r {
			for _, ev := range w.Events {
				if ev.Type == 0 {
					cb(KeyAdd, string(ev.Kv.Key), string(ev.Kv.Value))
				} else if ev.Type == 1 {
					cb(KeyDelete, string(ev.Kv.Key), string(ev.Kv.Value))
				}
			}
		}
	}()

	return watchHandle, nil
}

func (c *_EtcdClient) WatchKeys(keyPrefix string, cb func(int, string, string)) (int, error) {
	ctx, cancelFunc := context.WithCancel(context.Background())

	r := c.cli.Watch(ctx, keyPrefix, clientv3.WithPrefix())

	var watchHandle int
	c.mtx.Lock()

	watchHandle = c.watchHandle
	c.cancelMap[c.watchHandle] = cancelFunc
	c.watchHandle++

	c.mtx.Unlock()

	go func() {
		for w := range r {
			for _, ev := range w.Events {
				if ev.Type == 0 {
					cb(KeyAdd, string(ev.Kv.Key), string(ev.Kv.Value))
				} else if ev.Type == 1 {
					cb(KeyDelete, string(ev.Kv.Key), string(ev.Kv.Value))
				}
			}
		}
	}()

	return watchHandle, nil
}

func (c *_EtcdClient) CancelWatch(watchHandle int) error {

	c.mtx.Lock()
	defer c.mtx.Unlock()

	cancelFunc, ok := c.cancelMap[watchHandle]
	if !ok {
		return errors.New("invalid watchHandle")
	}

	delete(c.cancelMap, watchHandle)
	cancelFunc()

	return nil
}
