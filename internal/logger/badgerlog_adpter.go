package logger

import (
	"fmt"

	"github.com/dgraph-io/badger/v3"
	"github.com/sirupsen/logrus"
)

var _ badger.Logger = (*BadgerLogAdapter)(nil)

type BadgerLogAdapter struct {
	log logrus.FieldLogger
}

func NewBadgerLogAdapter() *BadgerLogAdapter {
	return &BadgerLogAdapter{
		log: GetLogrusInstance(),
	}
}

func (b *BadgerLogAdapter) Errorf(msg string, args ...interface{}) {
	b.log.Error(fmt.Sprintf(msg, args...))
}

func (b *BadgerLogAdapter) Warningf(msg string, args ...interface{}) {
	b.log.Warn(fmt.Sprintf(msg, args...))
}

func (b *BadgerLogAdapter) Infof(msg string, args ...interface{}) {
	b.log.Info(fmt.Sprintf(msg, args...))
}

func (b *BadgerLogAdapter) Debugf(msg string, args ...interface{}) {
	b.log.Debug(fmt.Sprintf(msg, args...))
}
