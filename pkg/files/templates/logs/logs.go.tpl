/*Generated code do not modify it*/
package logs

import (
    log "github.com/sirupsen/logrus"
)

func Init(logLevel string) {
    level, err := log.ParseLevel(logLevel)
    if err != nil {
        log.Warnf("Unknown log level %s defaulting to warning level", logLevel)
        level = log.WarnLevel
    }

    log.SetLevel(level)
    log.SetFormatter(&log.JSONFormatter{})
    log.SetReportCaller(true)
}
