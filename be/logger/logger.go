package logger

import (
	"be/constant"
	"github.com/sirupsen/logrus"
	"github.com/yukitsune/lokirus"
	"errors"
)

var Logger = logrus.New()

func InitLogger() error {
	if constant.APPCONFIG == nil {
		return errors.New("APPCONFIG is nil; cannot initialize logger")
	}

	opts := lokirus.NewLokiHookOptions().
		WithLevelMap(lokirus.LevelMap{logrus.PanicLevel: "critical"}).
		WithFormatter(&logrus.JSONFormatter{
			PrettyPrint: true,
		}).
		WithStaticLabels(lokirus.Labels{
			"app":         constant.APPCONFIG.AppName,
			"environment": constant.APPCONFIG.Environment,
		})
	hook := lokirus.NewLokiHookWithOpts(
		constant.APPCONFIG.LokiURL,
		opts,
		logrus.InfoLevel,
		logrus.WarnLevel,
		logrus.ErrorLevel,
		logrus.FatalLevel)
	Logger.AddHook(hook)
	return nil
}
