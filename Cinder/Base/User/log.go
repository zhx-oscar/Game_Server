package User

import (
	"Cinder/Base/Log"
	log "github.com/cihub/seelog"
)

var ulogger log.LoggerInterface

func init() {
	var err error

	ulogger, err = log.CloneLogger(Log.GetLogger())
	if err != nil {
		panic(err)
	}

	ulogger.SetAdditionalStackDepth(1)
}

func (u *User) fmtLogMessage(v []interface{}) []interface{} {
	var fmtList []interface{}
	fmtList = append(fmtList, u.debugStr)
	for _, arg := range v {
		fmtList = append(fmtList, " ", arg)
	}
	return fmtList
}

func (u *User) Info(v ...interface{}) {
	fmtList := u.fmtLogMessage(v)
	ulogger.Info(fmtList...)
}

func (u *User) Warn(v ...interface{}) {
	fmtList := u.fmtLogMessage(v)
	ulogger.Warn(fmtList...)
}

func (u *User) Error(v ...interface{}) {
	fmtList := u.fmtLogMessage(v)
	ulogger.Error(fmtList...)
}

func (u *User) Debug(v ...interface{}) {
	fmtList := u.fmtLogMessage(v)
	ulogger.Debug(fmtList...)
}

func (u *User) Infof(format string, v ...interface{}) {
	ff := "%s" + format
	params := []interface{}{u.debugStr}
	params = append(params, v...)
	ulogger.Infof(ff, params...)
}

func (u *User) Warnf(format string, v ...interface{}) {
	ff := "%s" + format
	params := []interface{}{u.debugStr}
	params = append(params, v...)
	ulogger.Warnf(ff, params...)
}

func (u *User) Errorf(format string, v ...interface{}) {
	ff := "%s" + format
	params := []interface{}{u.debugStr}
	params = append(params, v...)
	ulogger.Errorf(ff, params...)
}

func (u *User) Debugf(format string, v ...interface{}) {
	ff := "%s" + format
	params := []interface{}{u.debugStr}
	params = append(params, v...)
	ulogger.Debugf(ff, params...)
}
