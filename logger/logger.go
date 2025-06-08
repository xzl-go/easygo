package logger

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"sync"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

type Logger struct {
	level     LogLevel
	logFile   *os.File
	mu        sync.Mutex
	stdLogger *log.Logger
}

func New(level LogLevel, logPath string) *Logger {
	var logFile *os.File
	var err error

	if logPath != "" {
		logFile, err = os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("无法打开日志文件：%v", err)
		}
	}

	return &Logger{
		level:     level,
		logFile:   logFile,
		stdLogger: log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile),
	}
}

func (l *Logger) log(level LogLevel, format string, v ...interface{}) {
	if level < l.level {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	_, file, line, _ := runtime.Caller(2)
	msg := fmt.Sprintf(format, v...)
	logEntry := fmt.Sprintf("[%s] %s:%d %s", getLevelString(level), file, line, msg)

	l.stdLogger.Println(logEntry)

	if l.logFile != nil {
		fmt.Fprintln(l.logFile, logEntry)
	}

	if level == FATAL {
		os.Exit(1)
	}
}

func getLevelString(level LogLevel) string {
	switch level {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

func (l *Logger) Debug(format string, v ...interface{}) {
	l.log(DEBUG, format, v...)
}

func (l *Logger) Info(format string, v ...interface{}) {
	l.log(INFO, format, v...)
}

func (l *Logger) Warn(format string, v ...interface{}) {
	l.log(WARN, format, v...)
}

func (l *Logger) Error(format string, v ...interface{}) {
	l.log(ERROR, format, v...)
}

func (l *Logger) Fatal(format string, v ...interface{}) {
	l.log(FATAL, format, v...)
}

func (l *Logger) Close() {
	if l.logFile != nil {
		l.logFile.Close()
	}
}

var defaultLogger = New(INFO, "")
