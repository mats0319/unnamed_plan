package mlog

import (
	"log/slog"
	"testing"
	"time"
)

func TestCustomLog(t *testing.T) {
	InitializeTest()
	defer Close()

	logger := slog.Default().WithGroup("group1").WithGroup("group2")

	Log(logger, slog.LevelDebug, "log msg",
		slog.Int("key1", 10), slog.String("key2", "value2"))
}

func TestLogLevel(t *testing.T) {
	InitializeTest()
	defer Close()

	Debug("debug level log")
	Info("info level log")
	Warn("warn level log")
	Error("error level log")
}

func TestLogSplitFile(t *testing.T) {
	InitializeTest()
	defer Close()

	lastSize := handler.Size
	currentSize := handler.Size
	for lastSize <= currentSize { // exit loop when emit log split
		lastSize = currentSize

		slog.SetDefault(slog.New(handler.WithGroup("groupName")))
		Log(nil, 100, "this is a long long test log message",
			slog.String("key1", "value1"),
			slog.String("key2", "value2"),
		)

		currentSize = handler.Size
	}

	time.Sleep(time.Second * 3) // 阻塞，避免异步压缩goroutine因主程序退出而停止执行
}
