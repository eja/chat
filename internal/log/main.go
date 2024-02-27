// Copyright (C) 2023-2024 by Ubaldo Porcheddu <ubaldo@eja.it>

package log

import (
	"fmt"
	"log"
	"regexp"

	"github.com/eja/chat/internal/sys"
)

const (
	logLevelError = 1
	logLevelWarn  = 2
	logLevelInfo  = 3
	logLevelDebug = 4
	logLevelTrace = 5
)

func logMessage(level int, args ...interface{}) {
	msg := ""
	switch level {
	case logLevelError:
		msg = "[E]"
	case logLevelWarn:
		msg = "[W]"
	case logLevelInfo:
		msg = "[I]"
	case logLevelDebug:
		msg = "[D]"
	case logLevelTrace:
		msg = "[T]"
	}

	for _, arg := range args {
		if str, ok := arg.(string); ok {
			arg = regexp.MustCompile(`[\n\t\s]+`).ReplaceAllString(str, " ")
		}
		msg += fmt.Sprintf(" %v", arg)
	}
	if level <= sys.Options.LogLevel {
		log.Println(msg)
	}
}

func Error(args ...interface{}) {
	logMessage(logLevelError, args...)
}

func Warn(args ...interface{}) {
	logMessage(logLevelWarn, args...)
}

func Info(args ...interface{}) {
	logMessage(logLevelInfo, args...)
}

func Debug(args ...interface{}) {
	logMessage(logLevelDebug, args...)
}

func Trace(args ...interface{}) {
	logMessage(logLevelTrace, args...)
}
