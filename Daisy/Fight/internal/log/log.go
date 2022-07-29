package log

import (
	"fmt"
	"github.com/cihub/seelog"
	"github.com/spf13/viper"
	"sync/atomic"
)

// LogLevel 日志等级
const (
	Level_Off     int32 = iota // 关闭
	Level_Console              // 打印到控制台
	Level_File                 // 打印到日志文件
)

var logger seelog.LoggerInterface
var logDir string
var logLevel int32

func init() {
	dir := "../seelog/"
	levelStr := "off"
	if v := viper.GetString("FightLog.LogDir"); v != "" {
		dir = v
	}
	if v := viper.GetString("FightLog.LogLevel"); v != "" {
		levelStr = v
	}
	load(dir, levelStr)
}

func load(dir string, levelStr string) {
	defer seelog.Flush()
	seelog.Info("FightLogDir:", dir, ", FightLogLevel:", levelStr)

	switch levelStr {
	case "console":
		atomic.StoreInt32(&logLevel, Level_Console)
	case "file", "on":
		atomic.StoreInt32(&logLevel, Level_File)
	case "off":
		atomic.StoreInt32(&logLevel, Level_Off)
	default:
		seelog.Info("FightLogLevel unsupported, use off! [console/file|on/off]")
		atomic.StoreInt32(&logLevel, Level_Off)
	}

	logDir = dir

	if logger != nil {
		logger.Close()
	}
	logger = newLogger(atomic.LoadInt32(&logLevel), dir)
}

func newLogger(level int32, dir string) seelog.LoggerInterface {
	path := ""

	switch level {
	case Level_Off:
		fallthrough
	case Level_Console:
		path = `<seelog minlevel="debug" maxlevel="error">
				<outputs formatid="main">
					<filter levels="debug,info,error">
						<console />
					</filter>	
				</outputs>
				<formats>
					<format id="main" format="%Date(2006-01-02 15:04:05.999) %LEV [%File:%Line] %Msg%n"/>  
				</formats>
			</seelog>`
	case Level_File:
		fileName := dir + "Fight"
		path = `<seelog minlevel="debug" maxlevel="error">
				<outputs formatid="main">
					<filter levels="debug,info,error">
						<buffered size="10000" flushperiod="1000">
							<rollingfile type="date" filename="` + fileName + `.seelog" datepattern="2006.01.02.15" fullname="true" maxrolls="168"/>  
						</buffered>
						<filter levels="error">
							<console />
						</filter>
					</filter>
				</outputs>
				<formats>
					<format id="main" format="%Date(2006-01-02 15:04:05.999) %LEV [%File:%Line] %Msg%n"/>  
				</formats>
			</seelog>`
	}

	tLogger, err := seelog.LoggerFromConfigAsString(path)
	if err != nil {
		panic(err)
	}
	tLogger.SetAdditionalStackDepth(1)

	return tLogger
}

func SetLogLevel(level int32) {
	if level != Level_File && level != Level_Console {
		level = Level_Off
	}
	atomic.StoreInt32(&logLevel, level)

	if logger != nil {
		logger.Close()
	}
	logger = newLogger(atomic.LoadInt32(&logLevel), logDir)
}

func GetLogLevel() int32 {
	return atomic.LoadInt32(&logLevel)
}

func SetLogPath(dir string) {
	logDir = dir

	if logger != nil {
		logger.Close()
	}
	logger = newLogger(atomic.LoadInt32(&logLevel), logDir)
}

func Debug(v ...interface{}) {
	if atomic.LoadInt32(&logLevel) != Level_Off {
		logger.Debug(v...)
	}
}

func Debugf(format string, v ...interface{}) {
	if atomic.LoadInt32(&logLevel) != Level_Off {
		logger.Debugf(format, v...)
	}
}

func Debugfun(fun func() string) {
	if fun == nil {
		return
	}

	level := atomic.LoadInt32(&logLevel)
	if level == Level_Off {
		return
	}

	if atomic.LoadInt32(&logLevel) != Level_Off {
		logger.Debug(fun())
	}
}

func Info(v ...interface{}) {
	logger.Info(v...)
}

func Infof(format string, v ...interface{}) {
	logger.Infof(format, v...)
}

func Infofun(fun func() string) {
	if fun == nil {
		return
	}

	logger.Info(fun())
}

func Error(v ...interface{}) {
	logger.Error(v...)
}

func Errorf(format string, v ...interface{}) {
	logger.Errorf(format, v...)
}

func Errorfun(fun func() string) {
	if fun == nil {
		return
	}

	logger.Error(fun())
}

func Flush() {
	if logger != nil {
		logger.Flush()
	}
}

type Cache string

func (cache *Cache) Print(v ...interface{}) {
	level := atomic.LoadInt32(&logLevel)
	if level != Level_Off {
		*cache += Cache(fmt.Sprint(v...))
	}
}

func (cache *Cache) Printf(format string, v ...interface{}) {
	level := atomic.LoadInt32(&logLevel)
	if level != Level_Off {
		*cache += Cache(fmt.Sprintf(format, v...))
	}
}

func (cache *Cache) PrintFun(fun func() string) {
	if fun == nil {
		return
	}

	level := atomic.LoadInt32(&logLevel)
	if level != Level_Off {
		*cache += Cache(fun())
	}
}

func (cache *Cache) DebugFlush() {
	level := atomic.LoadInt32(&logLevel)
	if level != Level_Off {
		logger.Debug(*cache)
	}

	*cache = ""
}

func (cache *Cache) InfoFlush() {
	level := atomic.LoadInt32(&logLevel)
	if level != Level_Off {
		logger.Info(*cache)
	}

	*cache = ""
}

func (cache *Cache) ErrorFlush() {
	level := atomic.LoadInt32(&logLevel)
	if level != Level_Off {
		logger.Error(*cache)
	}

	*cache = ""
}
