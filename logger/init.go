package logger

import (
	"github.com/bool64/ctxd"
	"github.com/bool64/zapctxd"
)

// InitLogger initializes logger.
func InitLogger(cfg Config, devMode bool) *zapctxd.Logger {
	return zapctxd.New(zapctxd.Config{
		Level:   cfg.Level,
		DevMode: devMode,
		FieldNames: ctxd.FieldNames{
			Timestamp: "timestamp",
			Message:   "message",
		},
		StripTime: cfg.LockTime,
		Output:    cfg.Output,
	})
}
