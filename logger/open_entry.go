package logger

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
)

// 实现接口 OpenLoggerEntry
type OpenLogShadowEntry struct {
	*logrus.Entry
	nl *logrus.Logger
}

// logger 以及 entry的公用接口，方法一致，业务使用体验一致
type OpenLoggerEntry interface {
	WithField(key string, value interface{}) OpenLoggerEntry
	WithFields(fields Fields) OpenLoggerEntry
	WithError(err error) OpenLoggerEntry
	WithTime(t time.Time) OpenLoggerEntry
	WithObject(obj interface{}) OpenLoggerEntry
	Tracef(ctx context.Context, format string, args ...interface{})
	Debugf(ctx context.Context, format string, args ...interface{})
	Infof(ctx context.Context, format string, args ...interface{})
	Printf(ctx context.Context, format string, args ...interface{})
	Warnf(ctx context.Context, format string, args ...interface{})
	Warningf(ctx context.Context, format string, args ...interface{})
	Errorf(ctx context.Context, format string, args ...interface{})
	Fatalf(ctx context.Context, format string, args ...interface{})
	Panicf(ctx context.Context, format string, args ...interface{})
	Logf(ctx context.Context, level Level, format string, args ...interface{})
	Log(ctx context.Context, level Level, args ...interface{})
	Trace(ctx context.Context, args ...interface{})
	Debug(ctx context.Context, args ...interface{})
	Info(ctx context.Context, args ...interface{})
	Print(ctx context.Context, args ...interface{})
	Warn(ctx context.Context, args ...interface{})
	Warning(ctx context.Context, args ...interface{})
	Error(ctx context.Context, args ...interface{})
	Fatal(ctx context.Context, args ...interface{})
	Panic(ctx context.Context, args ...interface{})
	Logln(ctx context.Context, level Level, args ...interface{})
	Traceln(ctx context.Context, args ...interface{})
	Debugln(ctx context.Context, args ...interface{})
	Infoln(ctx context.Context, args ...interface{})
	Println(ctx context.Context, args ...interface{})
	Warnln(ctx context.Context, args ...interface{})
	Warningln(ctx context.Context, args ...interface{})
	Errorln(ctx context.Context, args ...interface{})
	Fatalln(ctx context.Context, args ...interface{})
	Panicln(ctx context.Context, args ...interface{})
}

func (en OpenLogShadowEntry) WithField(key string, value interface{}) OpenLoggerEntry {
	return &OpenLogShadowEntry{en.Entry.WithField(key, value), en.nl}
}

func (en OpenLogShadowEntry) WithFields(fields Fields) OpenLoggerEntry {
	return &OpenLogShadowEntry{en.Entry.WithFields(fields), en.nl}
}

func (en OpenLogShadowEntry) WithError(err error) OpenLoggerEntry {
	return &OpenLogShadowEntry{en.Entry.WithError(err), en.nl}
}

func (en OpenLogShadowEntry) WithTime(t time.Time) OpenLoggerEntry {
	return &OpenLogShadowEntry{en.Entry.WithTime(t), en.nl}
}

func (en OpenLogShadowEntry) WithObject(obj interface{}) OpenLoggerEntry {
	fields := parseFieldsFromObj(obj)
	return &OpenLogShadowEntry{en.Entry.WithFields(fields), en.nl}
}

func (en OpenLogShadowEntry) setLogger(ctx context.Context) {
	en.Entry.Logger = en.nl
}

func logLevelTrans(ctx context.Context, originLevel logrus.Level) logrus.Level {
	if originLevel > logrus.InfoLevel {
		return logrus.InfoLevel
	}

	return originLevel
}

func (en OpenLogShadowEntry) Tracef(ctx context.Context, format string, args ...interface{}) {
	en.setLogger(ctx)
	en.Entry.WithContext(ctx).Logf(logLevelTrans(ctx, TraceLevel), format, args...)
}

func (en OpenLogShadowEntry) Debugf(ctx context.Context, format string, args ...interface{}) {
	en.setLogger(ctx)
	en.Entry.WithContext(ctx).Logf(logLevelTrans(ctx, DebugLevel), format, args...)
}

func (en OpenLogShadowEntry) Infof(ctx context.Context, format string, args ...interface{}) {
	en.setLogger(ctx)
	en.Entry.WithContext(ctx).Logf(InfoLevel, format, args...)
}

func (en OpenLogShadowEntry) Printf(ctx context.Context, format string, args ...interface{}) {
	en.setLogger(ctx)
	en.Entry.WithContext(ctx).Printf(format, args...)
}

func (en OpenLogShadowEntry) Warnf(ctx context.Context, format string, args ...interface{}) {
	en.setLogger(ctx)
	en.Entry.WithContext(ctx).Logf(WarnLevel, format, args...)
}

func (en OpenLogShadowEntry) Warningf(ctx context.Context, format string, args ...interface{}) {
	en.setLogger(ctx)
	en.Entry.WithContext(ctx).Warnf(format, args...)
}

func (en OpenLogShadowEntry) Errorf(ctx context.Context, format string, args ...interface{}) {
	en.setLogger(ctx)
	en.Entry.WithContext(ctx).Logf(ErrorLevel, format, args...)
}

func (en OpenLogShadowEntry) Fatalf(ctx context.Context, format string, args ...interface{}) {
	en.setLogger(ctx)
	en.Entry.WithContext(ctx).Fatalf(format, args...)
}

func (en OpenLogShadowEntry) Panicf(ctx context.Context, format string, args ...interface{}) {
	en.setLogger(ctx)
	en.Entry.WithContext(ctx).Logf(PanicLevel, format, args...)
}

func (en OpenLogShadowEntry) Logf(ctx context.Context, level Level, format string, args ...interface{}) {
	en.setLogger(ctx)
	en.Entry.WithContext(ctx).Logf(logLevelTrans(ctx, level), format, args...)
}

func (en OpenLogShadowEntry) Log(ctx context.Context, level Level, args ...interface{}) {
	en.setLogger(ctx)
	en.Entry.WithContext(ctx).Log(logLevelTrans(ctx, level), args...)
}

func (en OpenLogShadowEntry) Trace(ctx context.Context, args ...interface{}) {
	en.setLogger(ctx)
	en.Entry.WithContext(ctx).Log(logLevelTrans(ctx, TraceLevel), args...)
}

func (en OpenLogShadowEntry) Debug(ctx context.Context, args ...interface{}) {
	en.setLogger(ctx)
	en.Entry.WithContext(ctx).Log(logLevelTrans(ctx, DebugLevel), args...)
}

func (en OpenLogShadowEntry) Info(ctx context.Context, args ...interface{}) {
	en.setLogger(ctx)
	en.Entry.WithContext(ctx).Log(InfoLevel, args...)
}

func (en OpenLogShadowEntry) Print(ctx context.Context, args ...interface{}) {
	en.setLogger(ctx)
	en.Entry.WithContext(ctx).Print(args...)
}

func (en OpenLogShadowEntry) Warn(ctx context.Context, args ...interface{}) {
	en.setLogger(ctx)
	en.Entry.WithContext(ctx).Log(WarnLevel, args...)
}

func (en OpenLogShadowEntry) Warning(ctx context.Context, args ...interface{}) {
	en.setLogger(ctx)
	en.Entry.WithContext(ctx).Warn(args...)
}

func (en OpenLogShadowEntry) Error(ctx context.Context, args ...interface{}) {
	en.setLogger(ctx)
	en.Entry.WithContext(ctx).Log(ErrorLevel, args...)
}

func (en OpenLogShadowEntry) Fatal(ctx context.Context, args ...interface{}) {
	en.setLogger(ctx)
	en.Entry.WithContext(ctx).Fatal(args...)
}

func (en OpenLogShadowEntry) Panic(ctx context.Context, args ...interface{}) {
	en.setLogger(ctx)
	en.Entry.WithContext(ctx).Panic(args...)
}

func (en OpenLogShadowEntry) Logln(ctx context.Context, level Level, args ...interface{}) {
	en.setLogger(ctx)
	en.Entry.WithContext(ctx).Logln(logLevelTrans(ctx, level), args...)
}

func (en OpenLogShadowEntry) Traceln(ctx context.Context, args ...interface{}) {
	en.setLogger(ctx)
	en.Entry.WithContext(ctx).Logln(logLevelTrans(ctx, TraceLevel), args...)
}

func (en OpenLogShadowEntry) Debugln(ctx context.Context, args ...interface{}) {
	en.setLogger(ctx)
	en.Entry.WithContext(ctx).Logln(logLevelTrans(ctx, DebugLevel), args...)
}

func (en OpenLogShadowEntry) Infoln(ctx context.Context, args ...interface{}) {
	en.setLogger(ctx)
	en.Entry.WithContext(ctx).Logln(InfoLevel, args...)
}

func (en OpenLogShadowEntry) Println(ctx context.Context, args ...interface{}) {
	en.setLogger(ctx)
	en.Entry.WithContext(ctx).Println(args...)
}

func (en OpenLogShadowEntry) Warnln(ctx context.Context, args ...interface{}) {
	en.setLogger(ctx)
	en.Entry.WithContext(ctx).Logln(WarnLevel, args...)
}

func (en OpenLogShadowEntry) Warningln(ctx context.Context, args ...interface{}) {
	en.setLogger(ctx)
	en.Entry.WithContext(ctx).Logln(WarnLevel, args...)
}

func (en OpenLogShadowEntry) Errorln(ctx context.Context, args ...interface{}) {
	en.setLogger(ctx)
	en.Entry.WithContext(ctx).Logln(ErrorLevel, args...)
}

func (en OpenLogShadowEntry) Fatalln(ctx context.Context, args ...interface{}) {
	en.setLogger(ctx)
	en.Entry.WithContext(ctx).Fatalln(args...)
}

func (en OpenLogShadowEntry) Panicln(ctx context.Context, args ...interface{}) {
	en.setLogger(ctx)
	en.Entry.WithContext(ctx).Logln(PanicLevel, args...)
}
