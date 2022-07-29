package mqnsq

import (
	log "github.com/cihub/seelog"
)

type mqLogger struct {
	logger log.LoggerInterface
}

func (l *mqLogger) Output(calldepth int, s string) error {
	l.logger.Error(s)
	return nil
}
