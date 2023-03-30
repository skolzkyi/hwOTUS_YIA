package logger

import (
	"go.uber.org/zap"
)

type LoggerWrap struct {
	config zap.Config
	logger *zap.SugaredLogger
}

func New(level string) (*LoggerWrap, error) {
	logWrap := LoggerWrap{}
	zlevel, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return nil, err
	}
	logWrap.config = zap.Config{
		Level:            zlevel,
		DisableCaller:    true,
		Development:      true,
		Encoding:         "console",
		OutputPaths:      []string{"stdout", "file_log.log"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig:    zap.NewDevelopmentEncoderConfig(),
	}
	logWrap.logger = zap.Must(logWrap.config.Build()).Sugar()
	return &logWrap, nil
}

func (l LoggerWrap) Info(msg string) {
	l.logger.Info(msg)
}

func (l LoggerWrap) Warning(msg string) {
	l.logger.Warn(msg)
}

func (l LoggerWrap) Error(msg string) {
	l.logger.Error(msg)
}

func (l LoggerWrap) Fatal(msg string) {
	l.logger.Fatal(msg)
}
