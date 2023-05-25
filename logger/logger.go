package logger

import (
	"bufio"
	"fmt"
	nested "github.com/antonfisher/nested-logrus-formatter"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
	"time"
)

var log = logrus.New()
var ginLog = logrus.New()

// var logName string = "roman"
var path string = "./log/roman.log"

func GinLogInit() {
	logInit(ginLog, "api")
}

func Init() {
	logInit(log, "server")
	log.AddHook(newFileHook())
}

func logInit(logs *logrus.Logger, ModuleName string) {
	src, err := os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		panic(err.Error())
	}
	writer := bufio.NewWriter(src)
	logs.SetOutput(writer)

	logs.SetFormatter(&nested.Formatter{
		TimestampFormat: "2006-01-02 15:04:05.000",
		HideKeys:        false,
		FieldsOrder:     []string{"FilePath"},
		CallerFirst:     false,
		NoFieldsColors:  true,
	})
	//Log.SetFormatter(&NewFormatter{})

	logs.SetLevel(logrus.DebugLevel)
	//log.AddHook(newFileHook())

	newLfsHook(logs, ModuleName, path, time.Hour*24*30)
	return
}

func newLfsHook(logs *logrus.Logger, ModuleName, path string, maxAge time.Duration) {
	lfHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: newWriter(ModuleName, path, "debug", maxAge),
		logrus.InfoLevel:  newWriter(ModuleName, path, "info", maxAge),
		logrus.WarnLevel:  newWriter(ModuleName, path, "warn", maxAge),
		logrus.ErrorLevel: newWriter(ModuleName, path, "error", maxAge),
		logrus.FatalLevel: newWriter(ModuleName, path, "fatal", maxAge),
		logrus.PanicLevel: newWriter(ModuleName, path, "panic", maxAge),
		//}, &NewFormatter{})
	}, &nested.Formatter{
		TimestampFormat: "2006-01-02 15:04:05.000",
		HideKeys:        false,
		FieldsOrder:     []string{"FilePath"},
		NoFieldsColors:  true,
	})

	logs.AddHook(lfHook)
}

func newWriter(moduleName, path, level string, exTime time.Duration) *rotatelogs.RotateLogs {
	fullPath := path + "-" + moduleName + "." + level
	writer, err := rotatelogs.New(
		fullPath+".%Y%m%d",
		//rotatelogs.WithLinkName(fullPath),
		//rotatelogs.ForceNewFile(),
		rotatelogs.WithMaxAge(-1),
		//rotatelogs.WithRotationSize(),
		rotatelogs.WithRotationCount(90),
	)
	if err != nil {
		log.Errorf("config local file system for logger error: %v", err)
		panic(err)
	}
	return writer
}

type NewFormatter struct {
}

func (n *NewFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	for k, v := range entry.Data {
		fmt.Printf("key: %v, val: %v\n", k, v)
	}
	//fmt.Println("data:", entry.Data)
	msg := fmt.Sprintf("[%s] [%s] %s\n", time.Now().Local().Format("2006-01-02 15:04:05.000"), strings.ToUpper(entry.Level.String()), entry.Message)
	return []byte(msg), nil
}
