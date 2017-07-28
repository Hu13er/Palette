package logger

import (
	"bytes"
	"fmt"
	"io"
	"sync"
)

const (
	PanicLevel = 1 << iota
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel

	AllLevel = PanicLevel | ErrorLevel | WarnLevel | InfoLevel | DebugLevel
)

type LogWriter interface {
	Write(buffer []byte, level int) (n int, err error)
	SetHeader(header []byte)

	Copy() LogWriter
}

type LoggerIntegrate struct {
	loggers []LogWriter
}

func NewLoggerIntegrate(loggers ...LogWriter) *LoggerIntegrate {
	ary := make([]LogWriter, 0)
	ary = append(ary, loggers...)
	return &LoggerIntegrate{ary}
}

func (li *LoggerIntegrate) Append(logger ...LogWriter) {
	li.loggers = append(li.loggers, logger...)
}

func (li *LoggerIntegrate) Write(buffer []byte, level int) (n int, err error) {
	for _, writer := range li.loggers {
		n, err = writer.Write(buffer, level)
		if err != nil {
			return
		}
	}
	return
}

func (li *LoggerIntegrate) Copy() LogWriter {
	newAry := make([]LogWriter, len(li.loggers))
	for k, v := range li.loggers {
		newAry[k] = v.Copy()
	}
	return NewLoggerIntegrate(newAry...)
}

func (li *LoggerIntegrate) SetHeader(header []byte) {
	for _, writer := range li.loggers {
		writer.SetHeader(header)
	}
}

type Logger struct {
	writer      io.Writer
	showBitmask int
	prefix      map[int]string
	header      []byte
	mutex       sync.Mutex
}

func New(writer io.Writer) *Logger {
	return &Logger{writer: writer, showBitmask: AllLevel, prefix: map[int]string{}, header: nil, mutex: sync.Mutex{}}
}

func (l *Logger) SetPrefix(prefix map[int]string) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if prefix != nil {
		l.prefix = prefix
	}
}

func (l *Logger) SetLevelPrefix(lvl int, prefix string) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.prefix[lvl] = prefix
}

func (l *Logger) SetShowBitmask(mask int) {
	l.showBitmask = mask
}

func (l *Logger) SetMinShowLevel(lvl int) {
	l.SetShowBitmask((lvl << 1) - 1)
}

func (l *Logger) Write(buffer []byte, level int) (n int, err error) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if l.showBitmask&level == 0 {
		return
	}
	nbuffer := l.header
	nbuffer = append(nbuffer, []byte(l.prefix[level])...)
	nbuffer = append(nbuffer, buffer...)
	n, err = l.writer.Write(nbuffer)

	l.header = l.header[0:0]
	n = len(buffer)
	return
}

func (l *Logger) SetHeader(header []byte) {
	l.header = header
}

func (l *Logger) Copy() LogWriter {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	nl := *l
	nl.mutex.Unlock()

	nprefix := map[int]string{}
	for k, v := range l.prefix {
		nprefix[k] = v
	}

	nheader := make([]byte, len(l.header))
	for k, v := range l.header {
		nheader[k] = v
	}

	return &nl
}

type writeStream struct {
	logger LogWriter
	level  int
}

func (ws *writeStream) Write(buffer []byte) (int, error) {
	return ws.logger.Write(buffer, ws.level)
}

func (ws *writeStream) SetHeader(header []byte) {
	ws.logger.SetHeader(header)
}

func GetWriterStream(logger LogWriter, level int) io.Writer {
	return &writeStream{logger, level}
}

type LogFmt struct {
	Logger LogWriter
}

func (lf *LogFmt) Write(buffer []byte, level int) (n int, err error) {
	n, err = lf.Logger.Write(buffer, level)
	return
}

func (lf *LogFmt) SetHeader(header []byte) {
	lf.Logger.SetHeader(header)
}

func (lf *LogFmt) Copy() LogWriter {
	nlf := *lf
	nlf.Logger = nlf.Logger.Copy()
	return &nlf
}

func (lf *LogFmt) printf(level int, format string, args ...interface{}) (n int, err error) {
	buffer := make([]byte, 0)
	stream := bytes.NewBuffer(buffer)
	n, err = fmt.Fprintf(stream, format, args...)
	if err != nil {
		return
	}
	n, err = lf.Write(stream.Bytes(), level)
	return
}

func (lf *LogFmt) print(level int, args ...interface{}) (n int, err error) {
	buffer := make([]byte, 0)
	stream := bytes.NewBuffer(buffer)
	n, err = fmt.Fprint(stream, args...)
	if err != nil {
		return
	}
	n, err = lf.Write(stream.Bytes(), level)
	return
}

func (lf *LogFmt) println(level int, args ...interface{}) (n int, err error) {
	args = append(args, "\n")
	return lf.print(level, args...)
}

func (lf *LogFmt) Debugf(format string, args ...interface{}) (int, error) {
	return lf.printf(DebugLevel, format, args...)
}

func (lf *LogFmt) Debug(args ...interface{}) (int, error) {
	return lf.print(DebugLevel, args...)
}

func (lf *LogFmt) Debugln(args ...interface{}) (int, error) {
	return lf.println(DebugLevel, args...)
}

func (lf *LogFmt) Infof(format string, args ...interface{}) (int, error) {
	return lf.printf(InfoLevel, format, args...)
}

func (lf *LogFmt) Info(args ...interface{}) (int, error) {
	return lf.print(InfoLevel, args...)
}

func (lf *LogFmt) Infoln(args ...interface{}) (int, error) {
	return lf.println(InfoLevel, args...)
}

func (lf *LogFmt) Warnf(format string, args ...interface{}) (int, error) {
	return lf.printf(WarnLevel, format, args...)
}

func (lf *LogFmt) Warn(args ...interface{}) (int, error) {
	return lf.print(WarnLevel, args...)
}

func (lf *LogFmt) Warnln(args ...interface{}) (int, error) {
	return lf.println(WarnLevel, args...)
}

func (lf *LogFmt) Errorf(format string, args ...interface{}) (int, error) {
	return lf.printf(ErrorLevel, format, args...)
}

func (lf *LogFmt) Error(args ...interface{}) (int, error) {
	return lf.print(ErrorLevel, args...)
}

func (lf *LogFmt) Errorln(args ...interface{}) (int, error) {
	return lf.println(ErrorLevel, args...)
}

func (lf *LogFmt) Panicf(format string, args ...interface{}) (int, error) {
	return lf.printf(PanicLevel, format, args...)
}

func (lf *LogFmt) Panic(args ...interface{}) (int, error) {
	return lf.print(PanicLevel, args...)
}

func (lf *LogFmt) Panicln(args ...interface{}) (int, error) {
	return lf.println(PanicLevel, args...)
}

func (lf *LogFmt) WithHeader(args ...interface{}) *LogFmt {
	out := lf.Copy().(*LogFmt)
	stream := new(bytes.Buffer)
	fmt.Fprint(stream, args...)
	out.SetHeader(stream.Bytes())
	return out
}

func (lf *LogFmt) WithHeaderln(args ...interface{}) *LogFmt {
	args = append(args, "\n")
	return lf.WithHeader(args...)
}

func (lf *LogFmt) WithHeaderf(format string, args ...interface{}) *LogFmt {
	out := lf.Copy().(*LogFmt)
	stream := new(bytes.Buffer)
	fmt.Fprintf(stream, format, args...)
	out.SetHeader(stream.Bytes())
	return out
}
