package log

import (
	"testing"
	"time"

	"go.uber.org/zap"
)

func TestLogSetLevelConcur(t *testing.T) {
	//特殊使用，在自己独立的协程里 创建一个子logger来独立控制level
	logger := NewStdLogger("debug")
	l2 := logger.With(zap.String("c", "1"))
	l2.SetLevel("error")
	go func() {
		for i := 10; i > 0; i-- {
			l2.Normal().Info("t1 info")
			time.Sleep(time.Second)
		}
	}()
	go func() {
		l3 := l2.With()
		l3.SetLevel("info")
		for i := 10; i > 0; i-- {
			l3.Normal().Info("t2 info")
			time.Sleep(time.Second)
		}
	}()
	<-time.After(11 * time.Second)
}

func TestLogSetLevel(t *testing.T) {
	//初始log debug level
	logger := NewStdLogger("debug")
	logger.Normal().Debug("L1 Debug")
	logger.Normal().Info("L1 Info")
	logger.Normal().Warn("L1 Warn")
	logger.Normal().Error("L1 Error")

	//继承log
	logger2 := logger.With(zap.String("child", "child log"))
	logger2.Normal().Debug("L2 Debug")
	logger2.Normal().Info("L2 Info")
	logger2.Normal().Warn("L2 Warn")
	logger2.Normal().Error("L2 Error")

	//动态设置level
	logger.SetLevel("warn")
	logger.Normal().Debug("L1 Debug Not show")
	logger.Normal().Info("L1 Info Not show")
	logger.Normal().Warn("L1 Warn")
	logger.Normal().Error("L1 Error")
	logger2.Normal().Debug("L2 Debug")
	logger2.Normal().Info("L2 Info")
	logger2.Normal().Warn("L2 Warn")
	logger2.Normal().Error("L2 Error")

	//动态设置logger2
	logger2.SetLevel("error")
	logger.Normal().Debug("L1 Debug Not show")
	logger.Normal().Info("L1 Info Not show")
	logger.Normal().Warn("L1 Warn")
	logger.Normal().Error("L1 Error")
	logger2.Normal().Debug("L2 Debug Not show")
	logger2.Normal().Info("L2 Info Not show")
	logger2.Normal().Warn("L2 Warn Not show")
	logger2.Normal().Error("L2 Error")
}
