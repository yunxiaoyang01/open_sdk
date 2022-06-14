package logger

import (
	"context"
	"testing"
)

func Test_InitLogger(t *testing.T) {
	options := Options{
		Level:   "debug",
		File:    "/var/log/open/open_logger.log",
		ErrFile: "/var/log/open/open_logger.err.log",
	}
	err := ResetStandardWithOptions(options)
	if err != nil {
		t.Errorf("failed to reset standard logger: %v", err)
	}
	ctx := context.Background()
	Debug(ctx, "debug logger print")
	Infof(ctx, "info logger print")
	Warnf(ctx, "warn logger print")
	Error(ctx, "error logger print")
}
