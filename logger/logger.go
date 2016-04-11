package logger

import (
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"
)

// SimpleFormatter is a formatter to describe simple logs
type SimpleFormatter struct {
	EnableColors     bool
	LogDepth         int
	EnableStackTrace []logrus.Level
}

var l *logrus.Logger

func init() {
	l = logrus.New()
	l.Level = logrus.DebugLevel
	l.Out = os.Stdout
}

//
// Convenience methods
//

// Debug is a logger for debug
func Debug(args ...interface{}) {
	l.Debug(args...)
}

// Debugf is a logger for debug with format
func Debugf(format string, args ...interface{}) {
	l.Debug(fmt.Sprintf(format, args...))
}

// Info is a logger for info
func Info(args ...interface{}) {
	l.Info(args...)
}

// Infof is a logger for info with format
func Infof(format string, args ...interface{}) {
	l.Info(fmt.Sprintf(format, args...))
}

// Warn is a logger for warn
func Warn(args ...interface{}) {
	l.Warn(args...)
}

// Warnf is a logger for warn with format
func Warnf(format string, args ...interface{}) {
	l.Warn(fmt.Sprintf(format, args...))
}

// Error is a logger for error
func Error(args ...interface{}) {
	l.Error(args...)

}

// Errorf is a logger for error with format
func Errorf(format string, args ...interface{}) {
	l.Error(fmt.Sprintf(format, args...))

}
