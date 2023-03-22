package logger

import (
	"go.uber.org/zap"
)

type LoggerWrap struct {
	config zap.Config
	logger *zap.Logger
}

func New(level string) *LoggerWrap {
	logWrap := LoggerWrap{}
	logWrap.config = zap.Config{
		Level:            zap.NewAtomicLevelAt(level),
		DisableCaller:    true,
		Development:      true,
		Encoding:         true,
		OutputPaths:      []string{"stdout", "file_log.log"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig:    zap.NewDevelopmentEncoderConfig(),
	}
	logWrap.logger = zap.Must(logWrap.config.Build()).Sugar()
	return &logWrap
}

func (l LoggerWrap) Info(msg string) {
	l.logger.Info(msg)
	//fmt.Println(msg)
}

func (l LoggerWrap) Warning(msg string) {
	l.logger.Warning(msg)
}

func (l LoggerWrap) Error(msg string) {
	l.logger.Error(msg)
}
