package logger

import (
	"fmt"
	"github.com/sirupsen/logrus"
)

func Debug(args ...any) {
	log.Debugln(args)
}

func Info(args ...any) {
	log.Infoln(args)
}

func Warn(args ...any) {
	log.Warnln(args)
}

func DebugF(format string, args ...any) {
	log.Debugf(format, args)
}

func InfoF(format string, args ...any) {
	log.Infof(format, args)
}

func WarnF(format string, args ...any) {
	log.Warnf(format, args)
}

func DebugByKv(msg string, args ...any) {
	fields := make(logrus.Fields)
	handleFields(fields, args)
	log.WithFields(fields).Debug(msg)
}

func InfoByKv(msg string, args ...any) {
	fields := make(logrus.Fields)
	handleFields(fields, args)
	log.WithFields(fields).Info(msg)
}

func WarnByKv(msg string, args ...any) {
	fields := make(logrus.Fields)
	handleFields(fields, args)
	log.WithFields(fields).Warn(msg)
}

func handleFields(fields logrus.Fields, args []any) {
	for i := 0; i < len(args); i += 2 {
		if i+1 < len(args) {
			fields[fmt.Sprintf("%v", args[i])] = args[i+1]
		} else {
			fields[fmt.Sprintf("%v", args[i])] = ""
		}
	}
}

func GinInfoByKv(msg string, args ...any) {
	fields := make(logrus.Fields)
	handleFields(fields, args)
	ginLog.WithFields(fields).Info(msg)
}
