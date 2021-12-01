package logging

import (
	log "github.com/sirupsen/logrus"
	"os"
)

func SetupLogger(isProduction bool) {
	if isProduction {
		log.SetFormatter(&log.JSONFormatter{})
		log.SetOutput(os.Stdout)
		return
	}

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
}
