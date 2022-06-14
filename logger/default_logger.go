package logger

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

type Entry = logrus.Entry
type Ext1FieldLogger = logrus.Ext1FieldLogger
type FieldLogger = logrus.FieldLogger
type Fields = logrus.Fields
type Formatter = logrus.Formatter
type Hook = logrus.Hook
type Logger = OpenLogger
type Level = logrus.Level
type LevelHooks = logrus.LevelHooks
type MutexWrap = logrus.MutexWrap

const PanicLevel = logrus.PanicLevel
const FatalLevel = logrus.FatalLevel
const ErrorLevel = logrus.ErrorLevel
const WarnLevel = logrus.WarnLevel
const InfoLevel = logrus.InfoLevel
const DebugLevel = logrus.DebugLevel
const TraceLevel = logrus.TraceLevel

var AllLevels = logrus.AllLevels
var openStdLogger = openStdLoggerNew()

func newJSONFormatter() logrus.Formatter {
	formatter := new(JSONFormatter)
	formatter.TimestampFormat = "2006-01-02T15:04:05.000Z07:00"
	return formatter
}

func newTextFormatter() logrus.Formatter {
	formatter := new(logrus.TextFormatter)
	formatter.TimestampFormat = "2006-01-02T15:04:05.000Z07:00"
	return formatter
}

// New 生成带有指定格式的标准logger
func openStdLoggerNew() Logger {
	formatter := newJSONFormatter()

	nl := logrus.Logger{
		Out:          os.Stderr,
		Formatter:    formatter,
		Hooks:        make(LevelHooks),
		Level:        InfoLevel,
		ExitFunc:     os.Exit,
		ReportCaller: false,
	}
	return &CtxLogger{&nl}
}

func StandardLogger() Logger {
	return openStdLogger
}

func SetOutput(out io.Writer) {
	openStdLogger.SetOutput(out)
}

func GetOutput() (out io.Writer) {
	return openStdLogger.GetOutput()
}

func SetFormatter(formatter Formatter) {
	openStdLogger.SetFormatter(formatter)
}

func SetReportCaller(include bool) {
	openStdLogger.SetReportCaller(include)
}

func SetLevel(level logrus.Level) {
	openStdLogger.SetLevel(level)
}

func SetLevelWithShadow(level logrus.Level) {
	openStdLogger.SetLevel(level)
}

func AddHook(hook logrus.Hook) {
	openStdLogger.AddHook(hook)
}

func ParseLevel(level string) (Level, error) {
	return logrus.ParseLevel(level)
}

func NewLogrusEntry(l *logrus.Logger) *Entry {
	return logrus.NewEntry(l)
}

func WithError(err error) OpenLoggerEntry {
	return openStdLogger.WithError(err)
}

func WithField(key string, value interface{}) OpenLoggerEntry {
	return openStdLogger.WithField(key, value)
}

func WithFields(fields Fields) OpenLoggerEntry {
	return openStdLogger.WithFields(fields)
}

func WithTime(t time.Time) OpenLoggerEntry {
	return openStdLogger.WithTime(t)
}

func WithObject(obj interface{}) OpenLoggerEntry {
	return openStdLogger.WithObject(obj)
}

func Trace(ctx context.Context, args ...interface{}) {
	openStdLogger.Trace(ctx, args...)
}

func Debug(ctx context.Context, args ...interface{}) {
	openStdLogger.Debug(ctx, args...)
}

func Print(ctx context.Context, args ...interface{}) {
	openStdLogger.Print(ctx, args...)
}

func Info(ctx context.Context, args ...interface{}) {
	openStdLogger.Info(ctx, args...)
}

func Warn(ctx context.Context, args ...interface{}) {
	openStdLogger.Warn(ctx, args...)
}

func Warning(ctx context.Context, args ...interface{}) {
	openStdLogger.Warning(ctx, args...)
}

func Error(ctx context.Context, args ...interface{}) {
	openStdLogger.Error(ctx, args...)
}

func Panic(ctx context.Context, args ...interface{}) {
	openStdLogger.Panic(ctx, args...)
}

func Fatal(ctx context.Context, args ...interface{}) {
	openStdLogger.Fatal(ctx, args...)
}

func Tracef(ctx context.Context, format string, args ...interface{}) {
	openStdLogger.Tracef(ctx, format, args...)
}

func Debugf(ctx context.Context, format string, args ...interface{}) {
	openStdLogger.Debugf(ctx, format, args...)
}

func Printf(ctx context.Context, format string, args ...interface{}) {
	openStdLogger.Printf(ctx, format, args...)
}

func Infof(ctx context.Context, format string, args ...interface{}) {
	openStdLogger.Infof(ctx, format, args...)
}

func Warnf(ctx context.Context, format string, args ...interface{}) {
	openStdLogger.Warnf(ctx, format, args...)
}

func Warningf(ctx context.Context, format string, args ...interface{}) {
	openStdLogger.Warningf(ctx, format, args...)
}

func Errorf(ctx context.Context, format string, args ...interface{}) {
	openStdLogger.Errorf(ctx, format, args...)
}

func Panicf(ctx context.Context, format string, args ...interface{}) {
	openStdLogger.Panicf(ctx, format, args...)
}

func Fatalf(ctx context.Context, format string, args ...interface{}) {
	openStdLogger.Fatalf(ctx, format, args...)
}

func Traceln(ctx context.Context, args ...interface{}) {
	openStdLogger.Traceln(ctx, args...)
}

func Debugln(ctx context.Context, args ...interface{}) {
	openStdLogger.Debugln(ctx, args...)
}

func Println(ctx context.Context, args ...interface{}) {
	openStdLogger.Println(ctx, args...)
}

func Infoln(ctx context.Context, args ...interface{}) {
	openStdLogger.Infoln(ctx, args...)
}

func Warnln(ctx context.Context, args ...interface{}) {
	openStdLogger.Warnln(ctx, args...)
}

func Warningln(ctx context.Context, args ...interface{}) {
	openStdLogger.Warningln(ctx, args...)
}

func Errorln(ctx context.Context, args ...interface{}) {
	openStdLogger.Errorln(ctx, args...)
}

func Panicln(ctx context.Context, args ...interface{}) {
	openStdLogger.Panicln(ctx, args...)
}

func Fatalln(ctx context.Context, args ...interface{}) {
	openStdLogger.Fatalln(ctx, args...)
}
