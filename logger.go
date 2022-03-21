package glogger

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type LogLevel struct {
	name          string
	level         int
	prefix        string
	formatterFlag int
}

func (l *LogLevel) isLowerOrEqualThan(log LogLevel) bool {
	return l.level <= log.level
}

func (l *LogLevel) getPrefix() string {
	if len(l.prefix) > 0 {
		return l.name + ": " + l.prefix + " "
	}
	return l.name + ": "
}

var (
	logLevelDebug LogLevel = LogLevel{name: "DEBUG", level: 1, formatterFlag: log.Ldate | log.Ltime | log.Lshortfile}
	logLevelInfo  LogLevel = LogLevel{name: "INFO", level: 2, formatterFlag: log.Ldate | log.Ltime | log.Lshortfile}
	logLevelWarn  LogLevel = LogLevel{name: "WARN", level: 3, formatterFlag: log.Ldate | log.Ltime | log.Lshortfile}
	logLevelError LogLevel = LogLevel{name: "ERROR", level: 4, formatterFlag: log.Ldate | log.Ltime | log.Llongfile}
)

var logLevel LogLevel = checkLogLevel()

func UpdateLogLevel() {
	logLevel = checkLogLevel()
}

func GetLogLevel() LogLevel {
	return logLevel
}

func loadPrefix(level string) string {
	prefixKey := "G_LOG_PREFIX"
	if len(level) > 0 {
		prefixLevel := os.Getenv(prefixKey + "_" + level)
		if len(prefixLevel) > 0 {
			return prefixLevel
		}
	}
	return os.Getenv(prefixKey)
}

func formatter(level string) int {
	formatKey := "G_LOG_FORMAT"
	if len(level) > 0 {
		formatKey = formatKey + "_" + level
	}
	format := os.Getenv(formatKey)

	if len(format) <= 0 {
		return 0
	}

	flags := strings.Split(format, "|")

	flag := 0
	for _, value := range flags {
		if value == "Ldate" { // the date in the local time zone: 2009/01/23
			flag += log.Ldate
		} else if value == "Ltime" { // the time in the local time zone: 01:23:23
			flag += log.Ltime
		} else if value == "Lmicroseconds" { // microsecond resolution: 01:23:23.123123.  assumes Ltime.
			flag += log.Lmicroseconds
		} else if value == "Llongfile" { // full file name and line number: /a/b/c/d.go:23
			flag += log.Llongfile
		} else if value == "Lshortfile" { // final file name element and line number: d.go:23. overrides Llongfile
			flag += log.Lshortfile
		} else if value == "LUTC" { // if Ldate or Ltime is set, use UTC rather than the local time zone
			flag += log.LUTC
		} else if value == "Lmsgprefix" { // move the "prefix" from the beginning of the line to before the message
			flag += log.Lmsgprefix
		} else if value == "LstdFlags" { // initial values for the standard logger
			flag += log.LstdFlags
		}
	}

	return flag
}

func checkLogLevel() LogLevel {
	godotenv.Load()
	level := os.Getenv("G_LOG_LEVEL")

	internalLog := *log.New(os.Stdout, "", log.LstdFlags)
	internalLog.Printf("Log Level: %v", level)

	if len(level) <= 0 {
		return logLevelDebug
	}

	logLevelDebug.prefix = loadPrefix(logLevelDebug.name)
	logLevelInfo.prefix = loadPrefix(logLevelInfo.name)
	logLevelWarn.prefix = loadPrefix(logLevelWarn.name)
	logLevelError.prefix = loadPrefix(logLevelError.name)

	if flag := formatter(logLevelDebug.name); flag > 0 {
		logLevelDebug.formatterFlag = flag
	}
	if flag := formatter(logLevelInfo.name); flag > 0 {
		logLevelInfo.formatterFlag = flag
	}
	if flag := formatter(logLevelWarn.name); flag > 0 {
		logLevelWarn.formatterFlag = flag
	}
	if flag := formatter(logLevelError.name); flag > 0 {
		logLevelError.formatterFlag = flag
	}

	if logLevelDebug.name == level {
		return logLevelDebug
	}
	if logLevelInfo.name == level {
		return logLevelInfo
	}
	if logLevelWarn.name == level {
		return logLevelWarn
	}
	if logLevelError.name == level {
		return logLevelError
	}

	return logLevelDebug
}

func Debug(format string, v ...interface{}) {
	if logLevel.isLowerOrEqualThan(logLevelDebug) {
		logDebug := *log.New(os.Stdout, logLevelDebug.getPrefix(), logLevelDebug.formatterFlag)
		logDebug.Output(2, fmt.Sprintf(format, v...))
	}
}

func Info(format string, v ...interface{}) {
	if logLevel.isLowerOrEqualThan(logLevelInfo) {
		logInfo := *log.New(os.Stdout, logLevelInfo.getPrefix(), logLevelInfo.formatterFlag)
		logInfo.Output(2, fmt.Sprintf(format, v...))
	}
}

func Warn(format string, v ...interface{}) {
	if logLevel.isLowerOrEqualThan(logLevelWarn) {
		logWarning := *log.New(os.Stderr, logLevelWarn.getPrefix(), logLevelWarn.formatterFlag)
		logWarning.Output(2, fmt.Sprintf(format, v...))
	}
}

func Error(format string, v ...interface{}) {
	if logLevel.isLowerOrEqualThan(logLevelError) {
		logError := *log.New(os.Stderr, logLevelError.getPrefix(), logLevelError.formatterFlag)
		logError.Output(2, fmt.Sprintf(format, v...))
	}
}

func Fatal(format string, v ...interface{}) {
	logError := *log.New(os.Stderr, logLevelError.getPrefix(), logLevelError.formatterFlag)
	logError.Output(2, fmt.Sprintf(format, v...))
	os.Exit(1)
}
