package Logger

import (
	"fmt"
	"time"
)

const (
	LogDebug = iota
	LogVerbose
	LogInfo
	LogNotice
	LogWarn
	LogError
	LogMerge
)

var errText = [7]string{
	"DEBUG",
	"VERBOSE",
	"INFO",
	"NOTICE",
	"WARN",
	"ERROR",
	"MERGE",
}

var loggerLevel = LogVerbose

func SetLogLevel(level ...int) {
	_logL := LogVerbose
	if len(level) > 0 {
		_logL = level[0]
	}
	if _logL < 0 || _logL > 6 {
		return
	}
	loggerLevel = _logL
}

func Log(level int, format string, a ...interface{}) {
	if level < loggerLevel {
		return
	}
	logMessage := fmt.Sprintf(format, a...)
	println(fmt.Sprintf("[%s, %s]%s", errText[level], time.Now().Format("2006-01-02 15:04:05.000000000"), logMessage))
}

func LogD(format string, a ...interface{}) {
	Log(LogDebug, format, a...)
}

func LogV(format string, a ...interface{}) {
	Log(LogVerbose, format, a...)
}

func LogI(format string, a ...interface{}) {
	Log(LogInfo, format, a...)
}

func LogN(format string, a ...interface{}) {
	Log(LogNotice, format, a...)
}

func LogW(format string, a ...interface{}) {
	Log(LogWarn, format, a...)
}

func LogE(format string, a ...interface{}) {
	Log(LogError, format, a...)
}

func LogM(format string, a ...interface{}) {
	Log(LogMerge, format, a...)
}
