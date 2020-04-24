package main

// "github.com/robfig/cron"
import (
	"xe_currencymanager/config"
	"xe_currencymanager/service"

	logger "github.com/sirupsen/logrus"
)

// load envi

func main() {
	logger.SetFormatter(&logger.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "02-01-2006 15:04:05",
	})

	config.InitViper()
	service.InitService()
}
