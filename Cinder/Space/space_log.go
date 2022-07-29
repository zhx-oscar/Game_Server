package Space

func (space *Space) Info(v ...interface{}) {
	logInfo(space.debugStr, v)
}

func (space *Space) Warn(v ...interface{}) {
	logWarn(space.debugStr, v)
}

func (space *Space) Error(v ...interface{}) {
	logError(space.debugStr, v)
}

func (space *Space) Debug(v ...interface{}) {
	logDebug(space.debugStr, v)
}

func (space *Space) Infof(format string, v ...interface{}) {
	logInfof(format, space.debugStr, v)
}

func (space *Space) Warnf(format string, v ...interface{}) {
	logWarnf(format, space.debugStr, v)
}

func (space *Space) Errorf(format string, v ...interface{}) {
	logErrorf(format, space.debugStr, v)
}

func (space *Space) Debugf(format string, v ...interface{}) {
	logDebugf(format, space.debugStr, v)
}
