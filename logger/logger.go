// Package logger 提供了灵活的日志记录功能
// 支持多级日志、文件和控制台输出
package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// 定义ANSI颜色常量
const (
	colorReset   = "\033[0m"
	colorRed     = "\033[31m"
	colorGreen   = "\033[32m"
	colorYellow  = "\033[33m"
	colorBlue    = "\033[34m"
	colorMagenta = "\033[35m"
	colorCyan    = "\033[36m"
	colorWhite   = "\033[37m"
)

// LogLevel 定义了日志级别类型
type LogLevel int

// 日志级别常量
const (
	DEBUG  LogLevel = iota // 调试级别
	INFO                   // 信息级别
	WARN                   // 警告级别
	ERROR                  // 错误级别
	FATAL                  // 致命错误级别
	EASYGO                 // EasyGo 框架特定日志级别
)

// Logger 是日志记录器
// 支持多级日志、文件输出和并发安全
type Logger struct {
	*log.Logger
	level     LogLevel    // 日志级别
	logFile   *os.File    // 日志文件 (如果只输出到控制台或文件打开失败，则为 nil)
	mu        sync.Mutex  // 互斥锁，保证并发安全
	stdLogger *log.Logger // 标准日志记录器 (始终输出到 os.Stdout)
}

var (
	debugLogger *Logger
	infoLogger  *Logger
	warnLogger  *Logger
	errorLogger *Logger
)

// New 创建一个新的日志记录器
// level: 日志级别
// baseLogDir: 日志文件存储的根目录，例如 "logs"。如果为空，则只输出到控制台。
// logFileName: 日志文件的基础名称，例如 "app.log"。
func New(level LogLevel, baseLogDir, logFileName string) *Logger {
	l := &Logger{
		level:     level,
		stdLogger: log.New(os.Stdout, "", 0), // 初始化标准日志记录器
	}

	if baseLogDir != "" && logFileName != "" {
		if err := os.MkdirAll(baseLogDir, 0755); err != nil {
			l.stdLogger.Printf("无法创建日志目录 %s: %v", baseLogDir, err)
		} else {
			file, err := os.OpenFile(
				filepath.Join(baseLogDir, logFileName),
				os.O_CREATE|os.O_WRONLY|os.O_APPEND,
				0644,
			)
			if err != nil {
				l.stdLogger.Printf("无法打开日志文件 %s: %v", logFileName, err)
			} else {
				l.logFile = file
				l.Logger = log.New(file, "", log.LstdFlags)
			}
		}
	}

	return l
}

// log 内部日志记录方法
// level: 日志级别
// format: 格式化字符串
// v: 格式化参数
func (l *Logger) log(level LogLevel, format string, v ...interface{}) {
	// 检查日志级别
	if level < l.level {
		return
	}

	// 加锁保证并发安全
	l.mu.Lock()
	defer l.mu.Unlock()

	msg := fmt.Sprintf(format, v...)
	now := time.Now().Format(time.DateTime) // 获取当前时间并格式化

	var color string
	var levelStr string

	switch level {
	case DEBUG:
		color = colorBlue
		levelStr = "DEBUG"
	case INFO:
		color = colorGreen
		levelStr = "INFO"
	case WARN:
		color = colorYellow
		levelStr = "WARN"
	case ERROR:
		color = colorRed
		levelStr = "ERROR"
	case FATAL:
		color = colorRed
		levelStr = "FATAL"
	case EASYGO:
		color = colorMagenta
		levelStr = "EASYGO"
	default:
		color = colorReset
		levelStr = "UNKNOWN"
	}

	var logEntry string
	if level == EASYGO {
		logEntry = fmt.Sprintf("%s[EASYGO] %s %s%s", color, now, msg, colorReset)
	} else {
		logEntry = fmt.Sprintf("%s[EASYGO - %s] %s %s%s", color, levelStr, now, msg, colorReset)
	}

	// 输出到控制台
	l.stdLogger.Println(logEntry)

	// 输出到文件 (文件不写入颜色码)
	if l.logFile != nil {
		// 移除颜色码再写入文件，避免文件内容被颜色码污染
		fileLogEntry := fmt.Sprintf("[EASYGO - %s] %s %s", levelStr, now, msg)
		if level == EASYGO {
			fileLogEntry = fmt.Sprintf("[EASYGO] %s %s", now, msg)
		}
		fmt.Fprintln(l.logFile, fileLogEntry)
	}

	// 如果是致命错误，则退出程序
	if level == FATAL {
		os.Exit(1)
	}
}

// getLevelString 获取日志级别的字符串表示
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
	case EASYGO:
		return "EASYGO"
	default:
		return "UNKNOWN"
	}
}

// Debug 记录调试级别日志
func (l *Logger) Debug(format string, v ...interface{}) {
	l.log(DEBUG, format, v...)
}

// Info 记录信息级别日志
func (l *Logger) Info(format string, v ...interface{}) {
	l.log(INFO, format, v...)
}

// Warn 记录警告级别日志
func (l *Logger) Warn(format string, v ...interface{}) {
	l.log(WARN, format, v...)
}

// Error 记录错误级别日志
func (l *Logger) Error(format string, v ...interface{}) {
	l.log(ERROR, format, v...)
}

// Fatal 记录致命错误级别日志并退出程序
func (l *Logger) Fatal(format string, v ...interface{}) {
	l.log(FATAL, format, v...)
}

// EasyGo 记录EasyGo框架启动等特定日志
func (l *Logger) EasyGo(format string, v ...interface{}) {
	l.log(EASYGO, format, v...)
}

// Close 关闭日志文件
func (l *Logger) Close() {
	if l.logFile != nil {
		l.logFile.Close()
	}
}

// 默认日志记录器实例
// 注意：如果需要文件输出，请确保在 New() 中传递 baseLogDir 和 logFileName 参数
var defaultLogger = New(INFO, "logs", "app.log") // 修改默认日志器，使其只输出到控制台

// Init 初始化日志记录器
func Init() {
	debugLogger = New(DEBUG, "logs", "debug.log")
	infoLogger = New(INFO, "logs", "info.log")
	warnLogger = New(WARN, "logs", "warn.log")
	errorLogger = New(ERROR, "logs", "error.log")
}

// 包级别日志函数
func Error(format string, v ...interface{}) {
	if errorLogger != nil {
		errorLogger.Error(format, v...)
	}
}
func Info(format string, v ...interface{}) {
	if infoLogger != nil {
		infoLogger.Info(format, v...)
	}
}
