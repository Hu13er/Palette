package logger

import (
	"io"
	"os"
)

var (
	log  *LogFmt
	slog *Logger
)

func Debug(args ...interface{}) (int, error) {
	return log.Debug(args...)
}

func Debugln(args ...interface{}) (int, error) {
	return log.Debugln(args...)
}

func Debugf(format string, args ...interface{}) (int, error) {
	return log.Debugf(format, args...)
}

func Info(args ...interface{}) (int, error) {
	return log.Info(args...)
}

func Infoln(args ...interface{}) (int, error) {
	return log.Infoln(args...)
}

func Infof(format string, args ...interface{}) (int, error) {
	return log.Infof(format, args...)
}

func Warn(args ...interface{}) (int, error) {
	return log.Warn(args...)
}

func Warnln(args ...interface{}) (int, error) {
	return log.Warnln(args...)
}

func Warnf(format string, args ...interface{}) (int, error) {
	return log.Warnf(format, args...)
}

func Error(args ...interface{}) (int, error) {
	return log.Error(args...)
}

func Errorln(args ...interface{}) (int, error) {
	return log.Errorln(args...)
}

func Errorf(format string, args ...interface{}) (int, error) {
	return log.Errorf(format, args...)
}

func Panic(args ...interface{}) (int, error) {
	return log.Panic(args...)
}

func Panicln(args ...interface{}) (int, error) {
	return log.Panicln(args...)
}

func Panicf(format string, args ...interface{}) (int, error) {
	return log.Panicf(format, args...)
}

func WithHeader(args ...interface{}) *LogFmt {
	return log.WithHeader(args...)
}

func WithHeaderln(args ...interface{}) *LogFmt {
	return log.WithHeaderln(args...)
}

func WithHeaderf(format string, args ...interface{}) *LogFmt {
	return log.WithHeaderf(format, args...)
}

func GetLogWriterStream(level int) io.Writer {
	return GetWriterStream(log, level)
}

func init() {
	prefix := map[int]string{
		DebugLevel: "(d) ",
		InfoLevel:  "(!) ",
		WarnLevel:  "[*] ",
		ErrorLevel: "[X] ",
		PanicLevel: "[#] ",
	}

	slog = New(os.Stderr)
	slog.SetPrefix(prefix)

	log = &LogFmt{slog}
}
