package Core

import (
	"fmt"
	log "github.com/cihub/seelog"
)

func (srv *_Core) Debug(v ...interface{}) {
	param := []interface{}{srv.getPrefix()}
	param = append(param, v...)
	log.Debug(param...)
}

func (srv *_Core) Info(v ...interface{}) {
	param := []interface{}{srv.getPrefix()}
	param = append(param, v...)
	log.Info(param...)
}

func (srv *_Core) Warning(v ...interface{}) {
	param := []interface{}{srv.getPrefix()}
	param = append(param, v...)
	log.Warn(param...)
}

func (srv *_Core) Error(v ...interface{}) {
	param := []interface{}{srv.getPrefix()}
	param = append(param, v...)
	log.Error(param...)
}

func (srv *_Core) getPrefix() string {
	return fmt.Sprintf("[SRV|%s|%s]", srv.GetServiceType(), srv.GetServiceID())
}
