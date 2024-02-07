// Copyright (C) 2023-2024 by Ubaldo Porcheddu <ubaldo@eja.it>

package main

import (
	"fmt"
	"log"
	"regexp"
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
	if level <= Options.LogLevel {
		log.Println(msg)
	}
}

func logError(args ...interface{}) {
	logMessage(logLevelError, args...)
}

func logWarn(args ...interface{}) {
	logMessage(logLevelWarn, args...)
}

func logInfo(args ...interface{}) {
	logMessage(logLevelInfo, args...)
}

func logDebug(args ...interface{}) {
	logMessage(logLevelDebug, args...)
}

func logTrace(args ...interface{}) {
	logMessage(logLevelTrace, args...)
}
