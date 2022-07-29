package Space

import (
	"Cinder/Base/Log"
	log "github.com/cihub/seelog"
)

var slogger log.LoggerInterface

func init() {
	var err error
	slogger, err = log.CloneLogger(Log.GetLogger())
	if err != nil {
		panic(err)
	}

	slogger.SetAdditionalStackDepth(2)
}

func fmtLogMessage(prefix string, v []interface{}) []interface{} {
	var fmtList []interface{}
	fmtList = append(fmtList, prefix)
	for _, arg := range v {
		fmtList = append(fmtList, " ", arg)
	}
	return fmtList
}

func logInfo(prefix string, v []interface{}) {
	fmtList := fmtLogMessage(prefix, v)
	slogger.Info(fmtList...)
}

func logWarn(prefix string, v []interface{}) {
	fmtList := fmtLogMessage(prefix, v)
	slogger.Warn(fmtList...)
}

func logError(prefix string, v []interface{}) {
	fmtList := fmtLogMessage(prefix, v)
	slogger.Error(fmtList...)
}

func logDebug(prefix string, v []interface{}) {
	fmtList := fmtLogMessage(prefix, v)
	slogger.Debug(fmtList...)
}

func logInfof(format string, prefix string, v []interface{}) {
	ff := "%s" + format
	params := []interface{}{prefix}
	params = append(params, v...)
	slogger.Infof(ff, params...)
}

func logWarnf(format string, prefix string, v []interface{}) {
	ff := "%s" + format
	params := []interface{}{prefix}
	params = append(params, v...)
	slogger.Warnf(ff, params...)
}

func logErrorf(format string, prefix string, v []interface{}) {
	ff := "%s" + format
	params := []interface{}{prefix}
	params = append(params, v...)
	slogger.Errorf(ff, params...)
}

func logDebugf(format string, prefix string, v []interface{}) {
	ff := "%s" + format
	params := []interface{}{prefix}
	params = append(params, v...)
	slogger.Debugf(ff, params...)
}
