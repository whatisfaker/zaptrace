// Copyright (c) 2017 Uber Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package log

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

//Level 自定义log等级
// type Level int8

// const (
// 	//nolint
// 	DebugLevel Level = iota - 1
// 	InfoLevel
// 	WarnLevel
// 	ErrorLevel
// 	// DPanicLevel
// 	// PanicLevel
// 	FatalLevel = 5
// )

// Factory is the default logging wrapper that can create
// logger instances either for a given Context or context-less.
type Factory struct {
	ZapLogger *zap.Logger
	slogger   *spanLogger
	logger    *logger
	level     zap.AtomicLevel
	writer    io.Writer
}

func parseLevelString(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "panic":
		return zapcore.PanicLevel
	case "fatal":
		return zapcore.FatalLevel
	}
	return zapcore.InfoLevel
}

func zapLevelToString(lv zapcore.Level) string {
	switch lv {
	case zapcore.DebugLevel:
		return "debug"
	case zapcore.WarnLevel:
		return "warn"
	case zapcore.ErrorLevel:
		return "error"
	case zapcore.PanicLevel:
		return "panic"
	case zapcore.FatalLevel:
		return "fatal"
	case zapcore.InfoLevel:
		return "info"
	}
	return fmt.Sprintf("%d", lv)
}

func newZapCoreWithLevelEnable(level zapcore.Level, wr io.Writer) (zapcore.Core, zap.AtomicLevel) {
	logLevel := zap.NewAtomicLevel()
	logLevel.SetLevel(level)
	w := zapcore.AddSync(wr)
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(config),
		w,
		logLevel,
	)
	return core, logLevel
}

func NewStdLogger(level string) *Factory {
	return NewStdLoggerWithOptions(level, zap.AddCaller(), zap.AddCallerSkip(1))
}

func NewFileLogger(filePath string, level string) *Factory {
	return NewFileLoggerWithOptions(filePath, 1024, false, level, zap.AddCaller(), zap.AddCallerSkip(1))
}

// NewFileLogger create a file log to specified path, (level: debug, info, warn, error, panic, fatal)
func NewFileLoggerWithOptions(filePath string, maxFileSizeInMB int, compressed bool, level string, options ...zap.Option) *Factory {
	lv := parseLevelString(level)
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:  filePath,
		MaxSize:   maxFileSizeInMB, //MB
		LocalTime: true,
		Compress:  compressed,
	})
	core, logLevel := newZapCoreWithLevelEnable(lv, w)
	l := zap.New(core, options...)
	slogger := &spanLogger{l, nil}
	logger := &logger{l}
	ins := &Factory{l, slogger, logger, logLevel, w}
	return ins
}

// NewStdLogger create a standard out factory, (level: debug, info, warn, error, panic, fatal), callSkip just for zap.AddCallSkip option
func NewStdLoggerWithOptions(level string, options ...zap.Option) *Factory {
	wr := os.Stdout
	core, logLevel := newZapCoreWithLevelEnable(parseLevelString(level), wr)
	l := zap.New(core, options...)
	slogger := &spanLogger{l, nil}
	logger := &logger{l}
	ins := &Factory{l, slogger, logger, logLevel, wr}
	return ins
}

// Normal creates a context-unaware logger.
func (b *Factory) Normal() Logger {
	return b.logger
}

func (b *Factory) SetLevel(level string) {
	b.level.SetLevel(parseLevelString(level))
}

func (b *Factory) Level() string {
	return zapLevelToString(b.level.Level())
}

// Trace returns a context-aware Logger. If the context
// contains an OpenTracing span, all logging calls are also
// echo-ed into the span.
func (b *Factory) Trace(ctx context.Context) Logger {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		if span != b.slogger.span {
			b.slogger.span = span
		}
		return b.slogger
	}
	return b.logger
}

// With creates a child logger, and optionally adds some context fields to that logger.
func (b *Factory) With(fields ...zapcore.Field) *Factory {
	core, logLevel := newZapCoreWithLevelEnable(b.level.Level(), b.writer) //share same writer fake code
	NewOptionWithLevelAndFields := zap.WrapCore(func(c zapcore.Core) zapcore.Core {
		core = core.With(fields)
		return core
	})
	l := b.ZapLogger.WithOptions(NewOptionWithLevelAndFields)
	return &Factory{
		ZapLogger: l,
		slogger:   &spanLogger{l, nil},
		logger:    &logger{l},
		level:     logLevel,
		writer:    b.writer,
	}
}
