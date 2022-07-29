package Prop

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

func fmtLogMessage(prefix string, v []interface{}) []interface{} {
	var fmtList []interface{}
	fmtList = append(fmtList, prefix)
	for _, arg := range v {
		fmtList = append(fmtList, " ", arg)
	}
	return fmtList
}

func infof(format string, prefix string, v ...interface{}) {
	ff := "%s " + format
	params := []interface{}{prefix}
	params = append(params, v...)
	ulogger.Infof(ff, params...)
}

func warnf(format string, prefix string, v ...interface{}) {
	ff := "%s " + format
	params := []interface{}{prefix}
	params = append(params, v...)
	ulogger.Warnf(ff, params...)
}

func errorf(format string, prefix string, v ...interface{}) {
	ff := "%s " + format
	params := []interface{}{prefix}
	params = append(params, v...)
	ulogger.Errorf(ff, params...)
}

func debugf(format string, prefix string, v ...interface{}) {
	ff := "%s " + format
	params := []interface{}{prefix}
	params = append(params, v...)
	ulogger.Debugf(ff, params...)
}

///////////////////////////////////////////////////////////////////////////////

func (obj *Object) Info(v ...interface{}) {
	fmtList := fmtLogMessage(obj.debugStr, v)
	ulogger.Info(fmtList...)
}

func (obj *Object) Warn(v ...interface{}) {
	fmtList := fmtLogMessage(obj.debugStr, v)
	ulogger.Warn(fmtList...)
}

func (obj *Object) Error(v ...interface{}) {
	fmtList := fmtLogMessage(obj.debugStr, v)
	ulogger.Error(fmtList...)
}

func (obj *Object) Debug(v ...interface{}) {
	fmtList := fmtLogMessage(obj.debugStr, v)
	ulogger.Debug(fmtList...)
}

func (obj *Object) Infof(format string, v ...interface{}) {
	infof(format, obj.debugStr, v...)
}

func (obj *Object) Warnf(format string, v ...interface{}) {
	warnf(format, obj.debugStr, v...)
}

func (obj *Object) Errorf(format string, v ...interface{}) {
	errorf(format, obj.debugStr, v...)
}

func (obj *Object) Debugf(format string, v ...interface{}) {
	debugf(format, obj.debugStr, v...)
}

///////////////////////////////////////////////////////////////////////////////

func (o *_Owner) Info(v ...interface{}) {
	fmtList := fmtLogMessage(o.debugStr, v)
	ulogger.Info(fmtList...)
}

func (o *_Owner) Warn(v ...interface{}) {
	fmtList := fmtLogMessage(o.debugStr, v)
	ulogger.Warn(fmtList...)
}

func (o *_Owner) Error(v ...interface{}) {
	fmtList := fmtLogMessage(o.debugStr, v)
	ulogger.Error(fmtList...)
}

func (o *_Owner) Debug(v ...interface{}) {
	fmtList := fmtLogMessage(o.debugStr, v)
	ulogger.Debug(fmtList...)
}

func (o *_Owner) Infof(format string, v ...interface{}) {
	infof(format, o.debugStr, v...)
}

func (o *_Owner) Warnf(format string, v ...interface{}) {
	warnf(format, o.debugStr, v...)
}

func (o *_Owner) Errorf(format string, v ...interface{}) {
	errorf(format, o.debugStr, v...)
}

func (o *_Owner) Debugf(format string, v ...interface{}) {
	debugf(format, o.debugStr, v...)
}
