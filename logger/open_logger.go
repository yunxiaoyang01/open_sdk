package logger

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

// OpenLogger的接口
type OpenLogger interface {
	SetOutput(out io.Writer)
	GetOutput() (out io.Writer)
	SetFormatter(formatter Formatter)
	SetReportCaller(include bool)
	SetLevel(level logrus.Level)
	AddHook(hook logrus.Hook)
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

type CtxLogger struct {
	n *logrus.Logger // normal logger
}

func NewCtxLogger() OpenLogger {
	formatter := new(JSONFormatter)
	formatter.TimestampFormat = "2006-01-02T15:04:05.000Z07:00"

	n := logrus.Logger{
		Out:          os.Stderr,
		Formatter:    formatter,
		Hooks:        make(LevelHooks),
		Level:        InfoLevel,
		ExitFunc:     os.Exit,
		ReportCaller: false,
	}
	return &CtxLogger{&n}
}

func (cl *CtxLogger) newOpenLogShadowEntry() OpenLoggerEntry {
	return &OpenLogShadowEntry{logrus.NewEntry(cl.n), cl.n}
}

func (cl *CtxLogger) SetOutput(out io.Writer) {
	cl.n.SetOutput(out)
}

func (cl *CtxLogger) GetOutput() (out io.Writer) {
	return cl.n.Out
}

func (cl *CtxLogger) SetFormatter(formatter Formatter) {
	cl.n.SetFormatter(formatter)
}

func (cl *CtxLogger) SetReportCaller(include bool) {
	cl.n.SetReportCaller(include)
}

func (cl *CtxLogger) SetLevel(level logrus.Level) {
	cl.n.SetLevel(level)
}

func (cl *CtxLogger) AddHook(hook logrus.Hook) {
	cl.n.AddHook(hook)
}

func (cl *CtxLogger) WithField(key string, value interface{}) OpenLoggerEntry {
	// 借用logrus.Logger本身Entry的管理机制来创建Entry,下同
	return &OpenLogShadowEntry{cl.n.WithField(key, value), cl.n}
}

func (cl *CtxLogger) WithFields(fields Fields) OpenLoggerEntry {
	return &OpenLogShadowEntry{cl.n.WithFields(fields), cl.n}
}

func (cl *CtxLogger) WithError(err error) OpenLoggerEntry {
	return &OpenLogShadowEntry{cl.n.WithError(err), cl.n}
}

func (cl *CtxLogger) WithTime(t time.Time) OpenLoggerEntry {
	return &OpenLogShadowEntry{cl.n.WithTime(t), cl.n}
}

func (cl *CtxLogger) WithObject(obj interface{}) OpenLoggerEntry {
	fields := parseFieldsFromObj(obj)
	return &OpenLogShadowEntry{cl.n.WithFields(fields), cl.n}
}

func (cl *CtxLogger) Tracef(ctx context.Context, format string, args ...interface{}) {
	cl.newOpenLogShadowEntry().Tracef(ctx, format, args...)
}

func (cl *CtxLogger) Debugf(ctx context.Context, format string, args ...interface{}) {
	cl.newOpenLogShadowEntry().Debugf(ctx, format, args...)
}

func (cl *CtxLogger) Infof(ctx context.Context, format string, args ...interface{}) {
	cl.newOpenLogShadowEntry().Infof(ctx, format, args...)
}

func (cl *CtxLogger) Printf(ctx context.Context, format string, args ...interface{}) {
	cl.newOpenLogShadowEntry().Printf(ctx, format, args...)
}

func (cl *CtxLogger) Warnf(ctx context.Context, format string, args ...interface{}) {
	cl.newOpenLogShadowEntry().Warnf(ctx, format, args...)
}

func (cl *CtxLogger) Warningf(ctx context.Context, format string, args ...interface{}) {
	cl.newOpenLogShadowEntry().Warningf(ctx, format, args...)
}

func (cl *CtxLogger) Errorf(ctx context.Context, format string, args ...interface{}) {
	cl.newOpenLogShadowEntry().Errorf(ctx, format, args...)
}

func (cl *CtxLogger) Fatalf(ctx context.Context, format string, args ...interface{}) {
	cl.newOpenLogShadowEntry().Fatalf(ctx, format, args...)
}

func (cl *CtxLogger) Panicf(ctx context.Context, format string, args ...interface{}) {
	cl.newOpenLogShadowEntry().Panicf(ctx, format, args...)
}

func (cl *CtxLogger) Logf(ctx context.Context, level Level, format string, args ...interface{}) {
	cl.newOpenLogShadowEntry().Logf(ctx, level, format, args...)
}

func (cl *CtxLogger) Log(ctx context.Context, level Level, args ...interface{}) {
	cl.newOpenLogShadowEntry().Log(ctx, level, args...)
}

func (cl *CtxLogger) Trace(ctx context.Context, args ...interface{}) {
	cl.newOpenLogShadowEntry().Trace(ctx, args...)
}

func (cl *CtxLogger) Debug(ctx context.Context, args ...interface{}) {
	cl.newOpenLogShadowEntry().Debug(ctx, args...)
}

func (cl *CtxLogger) Info(ctx context.Context, args ...interface{}) {
	cl.newOpenLogShadowEntry().Info(ctx, args...)
}

func (cl *CtxLogger) Print(ctx context.Context, args ...interface{}) {
	cl.newOpenLogShadowEntry().Print(ctx, args...)
}

func (cl *CtxLogger) Warn(ctx context.Context, args ...interface{}) {
	cl.newOpenLogShadowEntry().Warn(ctx, args...)
}

func (cl *CtxLogger) Warning(ctx context.Context, args ...interface{}) {
	cl.newOpenLogShadowEntry().Warning(ctx, args...)
}

func (cl *CtxLogger) Error(ctx context.Context, args ...interface{}) {
	cl.newOpenLogShadowEntry().Error(ctx, args...)
}

func (cl *CtxLogger) Fatal(ctx context.Context, args ...interface{}) {
	cl.newOpenLogShadowEntry().Fatal(ctx, args...)
}

func (cl *CtxLogger) Panic(ctx context.Context, args ...interface{}) {
	cl.newOpenLogShadowEntry().Panic(ctx, args...)
}

func (cl *CtxLogger) Logln(ctx context.Context, level Level, args ...interface{}) {
	cl.newOpenLogShadowEntry().Logln(ctx, level, args...)
}

func (cl *CtxLogger) Traceln(ctx context.Context, args ...interface{}) {
	cl.newOpenLogShadowEntry().Traceln(ctx, args...)
}

func (cl *CtxLogger) Debugln(ctx context.Context, args ...interface{}) {
	cl.newOpenLogShadowEntry().Debugln(ctx, args...)
}

func (cl *CtxLogger) Infoln(ctx context.Context, args ...interface{}) {
	cl.newOpenLogShadowEntry().Infoln(ctx, args...)
}

func (cl *CtxLogger) Println(ctx context.Context, args ...interface{}) {
	cl.newOpenLogShadowEntry().Println(ctx, args...)
}

func (cl *CtxLogger) Warnln(ctx context.Context, args ...interface{}) {
	cl.newOpenLogShadowEntry().Warnln(ctx, args...)
}

func (cl *CtxLogger) Warningln(ctx context.Context, args ...interface{}) {
	cl.newOpenLogShadowEntry().Warningln(ctx, args...)
}

func (cl *CtxLogger) Errorln(ctx context.Context, args ...interface{}) {
	cl.newOpenLogShadowEntry().Errorln(ctx, args...)
}

func (cl *CtxLogger) Fatalln(ctx context.Context, args ...interface{}) {
	cl.newOpenLogShadowEntry().Fatalln(ctx, args...)
}

func (cl *CtxLogger) Panicln(ctx context.Context, args ...interface{}) {
	cl.newOpenLogShadowEntry().Panicln(ctx, args...)
}
