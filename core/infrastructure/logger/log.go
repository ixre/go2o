package logger

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/ixre/go2o/core/infrastructure/fw"
	"github.com/sirupsen/logrus"
)

func init() {
	// 设置日志等级
	logrus.SetLevel(logrus.WarnLevel)
	logrus.SetFormatter(&logFormatter{"GO2O"})
}

type logFormatter struct {
	prefix string
}

func (m *logFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	timestamp := entry.Time.Format("2006-01-02 15:04:05")
	var msg string
	//entry.Logger.SetReportCaller(true)
	//HasCaller()为true才会有调用信息
	if entry.HasCaller() {
		fName := filepath.Base(entry.Caller.File)
		msg = fmt.Sprintf("%s [ %s][ %-7s] [%s:%d %s]: %s\n",
			timestamp,
			m.prefix,
			strings.ToUpper(entry.Level.String()),
			fName, entry.Caller.Line, entry.Caller.Function, entry.Message)
	} else {
		msg = fmt.Sprintf("%s [ %s][ %s]: %s\n", timestamp,
			m.prefix,
			strings.ToUpper(entry.Level.String()), entry.Message)
	}

	b.WriteString(msg)
	return b.Bytes(), nil
}

var log fw.ILogger = new(loggerImpl)

type loggerImpl struct {
}

// 输出调试日志
func (l *loggerImpl) Debug(format string, args ...interface{}) {
	logrus.Debugf(format, args...)
}

// 输出普通日志
func (l *loggerImpl) Info(format string, args ...interface{}) {
	logrus.Infof(format, args...)
}

// 输出错误日志
func (l *loggerImpl) Error(format string, args ...interface{}) {
	logrus.Errorf(format, args...)
}

// 输出警告日志
func (l *loggerImpl) Warn(format string, args ...interface{}) {
	logrus.Warnf(format, args...)

}

// 输出致命日志
func (l *loggerImpl) Fatal(format string, args ...interface{}) {
	logrus.Fatalf(format, args...)
}

// 输出调试日志
func Debug(format string, args ...interface{}) {
	log.Debug(format, args...)
}

// 输出普通日志
func Info(format string, args ...interface{}) {
	log.Info(format, args...)
}

// 输出错误日志
func Error(format string, args ...interface{}) {
	log.Error(format, args...)
}

// 输出警告日志
func Warn(format string, args ...interface{}) {
	log.Warn(format, args...)
}

// 输出致命日志
func Fatal(format string, args ...interface{}) {
	log.Fatal(format, args...)
}
