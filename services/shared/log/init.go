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

var zlog *zap.Logger

func Logger() *zap.Logger {
	return zlog
}

func Init() {
	if zlog != nil { // have initialized
		return
	}

	core := zapcore.NewCore(logEncoder(), logWriteSyncer(), logLevel())
	zlog = zap.New(core, zap.AddCaller())

	zlog.Info("> Config init finish.")
	zlog.Info("> Log init finish.")
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

func logWriteSyncer() zapcore.WriteSyncer {
	file, err := os.Create(mconst.LogFileName)
	if err != nil {
		log.Println("create log file failed, error:", err)
		os.Exit(-1)
	}

	return zapcore.AddSync(file)
}

func logLevel() (level zapcore.Level) {
	if mconfig.GetConfigLevel() == mconst.ConfigDevLevel {
		level = zapcore.DebugLevel
	} else {
		level = zapcore.InfoLevel
	}

	return
}

func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + t.Format("2006-01-02 15:04:05") + "]")
}
