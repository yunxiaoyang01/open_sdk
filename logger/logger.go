package logger

import (
	"os"
	"reflect"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	shadowDirName = "shadow"
	LogDir        = "OPEN_LOGDIR"
	SvcName       = "OPEN_SVCNAME"
)

var AppName string

func NewLogger(opts ...Option) (l Logger, err error) {
	return NewLoggerWithOptions(newOptions(opts...))
}

func NewLoggerWithOptions(options Options) (l Logger, err error) {
	l = openStdLoggerNew()
	if err = initLoggerWithOptions(l, options); err != nil {
		return nil, errors.Wrap(err, "failed to initialize logger")
	}
	return l, nil
}

func initLoggerWithOptions(l Logger, options Options) (err error) {
	if options.Level != "" { // 如果配置里指定了日志等级，则解析并设置，否则默认等级是info。
		level, err := ParseLevel(options.Level)
		if err != nil {
			return errors.Wrapf(err, "failed to parse level(%s)", options.Level)
		}
		l.SetLevel(level)
	}
	if options.File != "" { // 如果配置里指定了日志文件，则解析并设置，否则默认写到stderr。
		err = handleFileOutput(l, options.File) // 设置output、压测标志
		if err != nil {
			return errors.Wrapf(err, "failed to set logger.Output and set flow_control")
		}
	}

	if options.ErrFile != "" { // 如果配置里指定了错误日志文件，则额外将等级为error(及以上)的日志复制一份写到该文件中。
		errWriter, err := os.OpenFile(options.ErrFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			return errors.Wrapf(err, "failed to open err file(%s)", options.ErrFile)
		}
		l.AddHook(NewErrWriterHook(errWriter))
	}

	return
}

func handleFileOutput(l Logger, fileName string) error {
	writer, err := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return errors.Wrapf(err, "failed to open file(%s)", fileName)
	}
	l.SetOutput(writer) // 设置正常日志输出
	return nil
}

// see from https://gitlab.xiaoduoai.com/marketing/base/blob/master/log/log.go
func parseFieldsFromObj(o interface{}) logrus.Fields {
	logFields := logrus.Fields{}

	val := reflect.ValueOf(o)
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return logFields
		}
		val = val.Elem()
	}
	for i := 0; i < val.NumField(); i++ {
		fValue := val.Field(i)
		fType := val.Type().Field(i)
		if !isZero(fValue) && fValue.IsValid() && fType.PkgPath == "" { // exported fields
			if fValue.Kind() == reflect.Struct ||
				(fValue.Kind() == reflect.Ptr &&
					fValue.Elem().Kind() == reflect.Struct) {
				fields := parseFieldsFromObj(fValue.Interface())
				if fType.Anonymous {
					for k, v := range fields {
						logFields[k] = v
					}
				} else {
					logFields[fType.Name] = fields
				}
			} else {
				logFields[fType.Name] = fValue.Interface()
			}
		}
	}
	return logFields
}

// see https://gitlab.xiaoduoai.com/marketing/base/blob/master/log/log.go
func isZero(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String:
		return len(v.String()) == 0
	case reflect.Ptr, reflect.Interface:
		return v.IsNil()
	case reflect.Slice:
		return v.Len() == 0
	case reflect.Map:
		return v.Len() == 0
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Struct: // 不去确认
		return false
	}
	return false
}
