package mlog

import (
	"log"
	"log/slog"
	"strings"

	mconfig "github.com/mats0319/unnamed_plan/server/internal/config"
)

var handler *Handler

// InitializeTest 正式的初始化函数需要先从配置文件读取配置，而在一些不需要配置的场景中就需要一个测试用初始化函数了
// 例如测试备份/恢复功能、测试日志结构、分文件等
func InitializeTest() {
	var err error
	handler, err = newHandler("log.log", 1, slog.LevelDebug)
	if err != nil {
		log.Fatalln("open log file failed, error:", err)
	}

	slog.SetDefault(slog.New(handler))
}

func Initialize() {
	var err error
	handler, err = newHandler("log.log", 1, getLogLevel())
	if err != nil {
		log.Fatalln("open log file failed, error:", err)
	}

	slog.SetDefault(slog.New(handler))

	Info("> Config init.") // 不是这里才初始化的，但是只有这里（日志）初始化之后才能使用自定义结构打印这句话
	Info("> Log init.")
}

func Close() {
	handler.close()
}

func Debug(msg string) {
	slog.Debug(msg)
}

func Info(msg string) {
	slog.Info(msg)
}

func Warn(msg string) {
	slog.Warn(msg)
}

func Error(msg string) {
	slog.Error(msg)
}

// Log 适用于需要通过'WithAttrs'/'WithGroup'定制输出内容、生成新的logger实例的场景
// 例如：http请求、db操作......
//
// 其实可以把slog default logger / slog level复制到这里，在使用过程中就不需要使用slog
// 已调整为与默认打印函数使用相同的调用层级数（日志中的代码位置一项，使用'mlog.Info()'/'mlog.Log()'均显示为该函数位置）
func Log(logger *slog.Logger, level slog.Level, msg string, fields ...any) {
	if logger == nil {
		logger = slog.Default()
	}

	switch level {
	case slog.LevelInfo:
		logger.Info(msg, fields...)
	case slog.LevelWarn:
		logger.Warn(msg, fields...)
	case slog.LevelError:
		logger.Error(msg, fields...)
	default:
		logger.Debug(msg, fields...)
	}
}

func getLogLevel() slog.Level {
	levelStr := mconfig.GetLevel()

	var level slog.Level
	switch strings.ToLower(levelStr) {
	case "error":
		level = slog.LevelError
	case "warn":
		level = slog.LevelWarn
	case "info":
		level = slog.LevelInfo
	default: // 'debug' and other unknown levels
		level = slog.LevelDebug
	}

	return level
}
