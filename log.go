// Copyright (C) 2023-2024 by Ubaldo Porcheddu <ubaldo@eja.it>

package main

import (
	"fmt"
	"log"
	"regexp"
)

func Log(level int, args ...interface{}) {
	msg := ""
	switch level {
	case 1:
		msg = "[E]"
	case 2:
		msg = "[W]"
	case 3:
		msg = "[I]"
	case 4:
		msg = "[D]"
	case 5:
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

func Error(args ...interface{}) {
	Log(1, args...)
}

func Warn(args ...interface{}) {
	Log(2, args...)
}

func Info(args ...interface{}) {
	Log(3, args...)
}

func Debug(args ...interface{}) {
	Log(4, args...)
}

func Trace(args ...interface{}) {
	Log(5, args...)
}
