package log

import (
	"fmt"
	"os"
	"runtime"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

type Map = map[string]interface{}

var log *_Logger

func init() {
	// symPath, err := os.Executable()
	// if err != nil {
	// 	panic(err)
	// }
	// path, err := filepath.EvalSymlinks(symPath)
	// if err != nil {
	// 	panic(err)
	// }
	// _, execName := filepath.Split(path)

	// config := rollingwriter.Config{
	// 	TimeTagFormat:          time.DateTime,
	// 	LogPath:                "./logs",
	// 	FileName:               execName,
	// 	MaxRemain:              5,
	// 	RollingPolicy:          rollingwriter.TimeRolling,
	// 	RollingTimePattern:     "0 0 0 * * *",
	// 	RollingVolumeSize:      "64M",
	// 	WriterMode:             "lock",
	// 	BufferWriterThershould: (1 << 20) * 8,
	// 	Compress:               false,
	// }
	// writer, err := rollingwriter.NewWriterFromConfig(&config)
	// if err != nil {
	// 	panic(err)
	// }

	writer, err := rotatelogs.New(
		"./logs/app_%Y%m%d.log",
		// rotatelogs.WithLinkName("./logs/app.log"),
		rotatelogs.WithMaxAge(24*time.Hour*730),
		rotatelogs.WithRotationTime(time.Hour*24),
	)
	if err != nil {
		panic(err)
	}

	tmp := logrus.New()
	tmp.SetLevel(logrus.DebugLevel)
	tmp.SetOutput(writer)
	tmp.SetReportCaller(true)
	tmp.ExitFunc = os.Exit
	tmp.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.DateTime,
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "time",
			logrus.FieldKeyLevel: "level",
			logrus.FieldKeyMsg:   "desc",
			logrus.FieldKeyFile:  "file",
			logrus.FieldKeyFunc:  "func",
		},
		DisableTimestamp:  false,
		DisableHTMLEscape: false,
		CallerPrettyfier: func(_ *runtime.Frame) (string, string) {
			pc, file, line, ok := runtime.Caller(10)
			if !ok {
				return "", ""
			}

			return runtime.FuncForPC(pc).Name(), fmt.Sprintf("%s:%d", file, line)
		},
		PrettyPrint: true,
	})

	log = &_Logger{
		log: tmp,
	}

	// Logger = tmp
}

type _Logger struct {
	log *logrus.Logger
}

func (l *_Logger) Debugf(format string, args ...any) {
	l.log.Debugf(format, args...)
}

func (l *_Logger) Infof(format string, args ...any) {
	l.log.Infof(format, args...)
}

func (l *_Logger) Warnf(format string, args ...any) {
	l.log.Warnf(format, args...)
}

func (l *_Logger) Errorf(format string, args ...any) {
	l.log.Errorf(format, args...)
}

func (l *_Logger) Panicf(format string, args ...any) {
	l.log.Panicf(format, args...)
}

func (l *_Logger) DebugWf(format string, fields Map, args ...any) {
	l.log.WithFields(fields).Debugf(format, args...)
}

func (l *_Logger) InfoWf(format string, fields Map, args ...any) {
	l.log.WithFields(fields).Infof(format, args...)
}

func (l *_Logger) WarnWf(format string, fields Map, args ...any) {
	l.log.WithFields(fields).Warnf(format, args...)
}

func (l *_Logger) ErrorWf(format string, fields Map, args ...any) {
	l.log.WithFields(fields).Errorf(format, args...)
}

func (l *_Logger) PanicWf(format string, fields Map, args ...any) {
	l.log.WithFields(fields).Panicf(format, args...)
}
