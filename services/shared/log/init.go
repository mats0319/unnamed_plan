package mlog

import (
	"github.com/mats9693/unnamed_plan/services/shared/config"
	"github.com/mats9693/unnamed_plan/services/shared/const"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
	"time"
)

var zLog *zap.Logger

func Logger() *zap.Logger {
	return zLog
}

func Init() error {
	if zLog != nil { // have initialized
		log.Println("already initialized")
		return nil
	}

	ws, err := logWriteSyncer()
	if err != nil {
		return err
	}

	coreSlice := make([]zapcore.Core, 0, 2)
	coreSlice = append(coreSlice, zapcore.NewCore(logEncoder(), ws, logLevel())) // log file
	if mconfig.GetConfigLevel() == mconst.ConfigLevel_Dev || mconfig.GetConfigLevel() == mconst.ConfigLevel_Default {
		coreSlice = append(coreSlice, zapcore.NewCore(logEncoder(), os.Stdout, logLevel())) // console
	}

	core := zapcore.NewTee(coreSlice...)
	zLog = zap.New(core, zap.AddCaller())

	zLog.Info("> Config init finish.")
	zLog.Info("> Log init finish.")

	return nil
}

func logEncoder() zapcore.Encoder {
	return zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		// Keys can be anything except the empty string.
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "name",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     timeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	})
}

func logWriteSyncer() (zapcore.WriteSyncer, error) {
	file, err := os.OpenFile(mconst.LogFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Println("create log file failed, error:", err)
		return nil, err
	}

	return zapcore.AddSync(file), nil
}

func logLevel() (level zapcore.Level) {
	if mconfig.GetConfigLevel() == mconst.ConfigLevel_Dev {
		level = zapcore.DebugLevel
	} else {
		level = zapcore.InfoLevel
	}

	return
}

func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + t.Format("2006-01-02 15:04:05") + "]")
}
