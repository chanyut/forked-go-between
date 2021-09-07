package logger

import (
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var instance logrus.FieldLogger
var once sync.Once

func GetLogrusInstance() logrus.FieldLogger {
	once.Do(func() {
		log := logrus.New()

		if level, err := logrus.ParseLevel(viper.GetString("log-level")); err != nil {
			log.SetLevel(level)
		}

		if viper.GetBool("log-json") {
			log.SetFormatter(&logrus.JSONFormatter{})
		}

		instance = log
	})

	return instance
}
