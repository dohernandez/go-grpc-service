package logger

import (
	"io"

	"go.uber.org/zap/zapcore"
)

type Config struct {
	Level      zapcore.Level `envconfig:"LOG_LEVEL" default:"error"`
	FieldNames string        `envconfig:"LOG_FILENAMES" default:"true"`
	Output     io.Writer
	// LockTime disables time variance in logger.
	LockTime bool

	// CallerSkip configures how deeply func calls should be skipped, default 1.
	CallerSkip int
}
