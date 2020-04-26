package tracing

import (
	"fmt"

	"github.com/whatisfaker/zaptrace/log"
)

type jaegerLoggerImp struct {
	logger log.Logger
}

func (c jaegerLoggerImp) Error(msg string) {
	c.logger.Error(msg)
}

func (c jaegerLoggerImp) Infof(msg string, args ...interface{}) {
	c.logger.Info(fmt.Sprintf(msg, args...))
}
