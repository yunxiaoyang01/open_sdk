package logger

import (
	"io"

	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

type OpenLfsHook struct {
	*lfshook.LfsHook
	isCopyTestLog bool
}

// 将等级为error(及以上)的日志复制一份写到errWriter。
func NewErrWriterHook(errWriter io.Writer) *OpenLfsHook {
	lfsh := NewOpenLfsHook(
		lfshook.WriterMap{
			ErrorLevel: errWriter,
			FatalLevel: errWriter,
			PanicLevel: errWriter,
		}, newJSONFormatter())
	lfsh.SetIsCopyTestLog(false) // 压测流量不复制
	return lfsh
}

func NewOpenLfsHook(output interface{}, formatter logrus.Formatter) *OpenLfsHook {
	return &OpenLfsHook{lfshook.NewHook(output, formatter), false}
}

// 覆盖LfsHook的同名方法，控制压测日志的输出
func (hook *OpenLfsHook) Fire(entry *logrus.Entry) error {
	return hook.LfsHook.Fire(entry)
}

func (hook *OpenLfsHook) SetIsCopyTestLog(isCopyTestLog bool) {
	hook.isCopyTestLog = isCopyTestLog
}
