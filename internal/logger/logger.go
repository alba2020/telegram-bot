package logger

import (
	"log"

	"go.uber.org/zap"
)

type Mode int

const (
	Production Mode = iota
	Development
)

var logger *zap.SugaredLogger

func init() {
	Init(Development)
}

func Init(mode Mode) {
	var newLogger *zap.Logger
	var err error
	if mode == Development {
		newLogger, err = zap.NewDevelopment()
	} else if mode == Production {
		cfg := zap.NewProductionConfig()
		cfg.DisableCaller = true
		cfg.DisableStacktrace = true
		cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
		newLogger, err = cfg.Build()
	} else {
		log.Fatal("Unknown logger mode")
	}
	if err != nil {
		log.Fatal("Cannot init zap", err)
	}

	opt := zap.AddCallerSkip(1)
	newLogger = newLogger.WithOptions(opt)
	logger = newLogger.Sugar()
}

func Debug(args ...any) {
	logger.Debug(args)
}

func Debugf(template string, args ...any) {
	logger.Debugf(template, args)
}

func Info(args ...any) {
	logger.Info(args)
}

func Infof(template string, args ...any) {
	logger.Infof(template, args)
}

func Warn(args ...any) {
	logger.Warn(args)
}

func Error(args ...any) {
	logger.Warn(args)
}

func Fatal(args ...any) {
	logger.Fatal(args)
}
