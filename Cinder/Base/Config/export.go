package Config

const (
	KeyAdd    = 0
	KeyDelete = 1
)

type IConfig interface {
	SetValueAndKeepAlive(key string, value string) error
	SetValue(key string, value string) error
	SetValueAndOvertime(key string, value string, ttl int64) error
	GetValue(key string) (string, error)
	GetValuesByPrefix(prefix string) ([]string, []string, error)
	WatchKey(key string, cb func(int, string, string)) (int, error)
	WatchKeys(keyPrefix string, cb func(int, string, string)) (int, error)
	CancelWatch(watchHandle int) error
}

var Inst IConfig
