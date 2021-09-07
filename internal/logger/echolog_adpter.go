package logger

import (
	"encoding/json"
	"io"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/sirupsen/logrus"
)

var _ echo.Logger = (*EchoLogAdapter)(nil)

type EchoLogAdapter struct {
	log logrus.FieldLogger
}

func NewEchoLogAdapter() *EchoLogAdapter {
	return &EchoLogAdapter{
		log: GetLogrusInstance(),
	}
}

func (e *EchoLogAdapter) Output() io.Writer {
	var w io.Writer

	logger, ok := e.log.(*logrus.Logger)
	if ok {
		w = logger.Out
	}

	if w == nil {
		w = os.Stderr
	}

	return w
}

func (e *EchoLogAdapter) SetOutput(w io.Writer) { /* nope */ }

func (e *EchoLogAdapter) Level() log.Lvl {
	w := logrus.InfoLevel

	logger, ok := e.log.(*logrus.Logger)
	if ok {
		w = logger.GetLevel()
	}

	return toEchoLevel(w)
}

func (e *EchoLogAdapter) SetLevel(v log.Lvl) { /* nope */ }

func (e *EchoLogAdapter) SetHeader(h string) { /* nope */ }

func (e *EchoLogAdapter) Formatter() logrus.Formatter {
	var w logrus.Formatter

	logger, ok := e.log.(*logrus.Logger)
	if ok {
		w = logger.Formatter
	}

	return w
}

func (e *EchoLogAdapter) SetFormatter(formatter logrus.Formatter) { /* nope */ }

func (e *EchoLogAdapter) Prefix() string {
	return ""
}

func (e *EchoLogAdapter) SetPrefix(p string) {

}

func (e *EchoLogAdapter) Print(i ...interface{}) {
	e.log.Print(i...)
}

func (e *EchoLogAdapter) Printf(format string, args ...interface{}) {
	e.log.Printf(format, args...)
}

func (e *EchoLogAdapter) Printj(j log.JSON) {
	b, err := json.Marshal(j)
	if err != nil {
		panic(err)
	}

	e.log.Println(string(b))
}

func (e *EchoLogAdapter) Debug(i ...interface{}) {
	e.log.Debug(i...)
}

func (e *EchoLogAdapter) Debugf(format string, args ...interface{}) {
	e.log.Debugf(format, args...)
}

func (e *EchoLogAdapter) Debugj(j log.JSON) {
	b, err := json.Marshal(j)
	if err != nil {
		panic(err)
	}

	e.log.Debugln(string(b))
}

func (e *EchoLogAdapter) Info(i ...interface{}) {
	e.log.Info(i...)
}

func (e *EchoLogAdapter) Infof(format string, args ...interface{}) {
	e.log.Infof(format, args...)
}

func (e *EchoLogAdapter) Infoj(j log.JSON) {
	b, err := json.Marshal(j)
	if err != nil {
		panic(err)
	}

	e.log.Infoln(string(b))
}

func (e *EchoLogAdapter) Warn(i ...interface{}) {
	e.log.Warn(i...)
}

func (e *EchoLogAdapter) Warnf(format string, args ...interface{}) {
	e.log.Warnf(format, args...)
}

func (e *EchoLogAdapter) Warnj(j log.JSON) {
	b, err := json.Marshal(j)
	if err != nil {
		panic(err)
	}

	e.log.Warnln(string(b))
}

func (e *EchoLogAdapter) Error(i ...interface{}) {
	e.log.Error(i...)
}

func (e *EchoLogAdapter) Errorf(format string, args ...interface{}) {
	e.log.Errorf(format, args...)
}

func (e *EchoLogAdapter) Errorj(j log.JSON) {
	b, err := json.Marshal(j)
	if err != nil {
		panic(err)
	}

	e.log.Errorln(string(b))
}

func (e *EchoLogAdapter) Fatal(i ...interface{}) {
	e.log.Fatal(i...)
}

func (e *EchoLogAdapter) Fatalf(format string, args ...interface{}) {
	e.log.Fatalf(format, args...)
}

func (e *EchoLogAdapter) Fatalj(j log.JSON) {
	b, err := json.Marshal(j)
	if err != nil {
		panic(err)
	}

	e.log.Fatalln(string(b))
}

func (e *EchoLogAdapter) Panic(i ...interface{}) {
	e.log.Panic(i...)
}

func (e *EchoLogAdapter) Panicf(format string, args ...interface{}) {
	e.log.Panicf(format, args...)
}

func (e *EchoLogAdapter) Panicj(j log.JSON) {
	b, err := json.Marshal(j)
	if err != nil {
		panic(err)
	}

	e.log.Panicln(string(b))
}

func toEchoLevel(level logrus.Level) log.Lvl {
	switch level {
	case logrus.DebugLevel:
		return log.DEBUG
	case logrus.InfoLevel:
		return log.INFO
	case logrus.WarnLevel:
		return log.WARN
	case logrus.ErrorLevel:
		return log.ERROR
	}

	return log.OFF
}
