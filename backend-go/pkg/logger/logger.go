package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

// Init 初始化日志
func Init(level string) {
	Log = logrus.New()

	// 设置日志格式
	Log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// 设置日志级别
	switch level {
	case "debug":
		Log.SetLevel(logrus.DebugLevel)
	case "info":
		Log.SetLevel(logrus.InfoLevel)
	case "warn":
		Log.SetLevel(logrus.WarnLevel)
	case "error":
		Log.SetLevel(logrus.ErrorLevel)
	default:
		Log.SetLevel(logrus.InfoLevel)
	}

	// 设置输出
	Log.SetOutput(os.Stdout)
}

// GetLogger 获取日志实例
func GetLogger() *logrus.Logger {
	return Log
}
