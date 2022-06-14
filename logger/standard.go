package logger

import (
	"github.com/pkg/errors"
)

type Standard Options

var _ = ResetStandard()

func ResetStandard(opts ...Option) (err error) {
	return ResetStandardWithOptions(newOptions(opts...))
}

func ResetStandardWithOptions(options Options) (err error) {
	l := StandardLogger()
	if err = initLoggerWithOptions(l, options); err != nil {
		return errors.Wrap(err, "failed to initialize logger")
	}

	return nil
}
