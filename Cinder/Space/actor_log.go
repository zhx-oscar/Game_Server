package Space

func (actor *Actor) Info(v ...interface{}) {
	logInfo(actor.debugStr, v)
}

func (actor *Actor) Warn(v ...interface{}) {
	logWarn(actor.debugStr, v)
}

func (actor *Actor) Error(v ...interface{}) {
	logError(actor.debugStr, v)
}

func (actor *Actor) Debug(v ...interface{}) {
	logDebug(actor.debugStr, v)
}

func (actor *Actor) Infof(format string, v ...interface{}) {
	logInfof(format, actor.debugStr, v)
}

func (actor *Actor) Warnf(format string, v ...interface{}) {
	logWarnf(format, actor.debugStr, v)
}

func (actor *Actor) Errorf(format string, v ...interface{}) {
	logErrorf(format, actor.debugStr, v)
}

func (actor *Actor) Debugf(format string, v ...interface{}) {
	logDebugf(format, actor.debugStr, v)
}
