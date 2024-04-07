package main

import (
	"context"
	"flag"

	"github.com/omegaatt36/gotasker/api"
	"github.com/omegaatt36/gotasker/logging"
	"github.com/omegaatt36/gotasker/util"
)

var (
	appPort  *string
	logLevel *string
	appENV   *string
)

func parseConfig() {
	_appPort := util.GetENV("APP_PORT", "8070")
	_appENV := util.GetENV("APP_ENV", "dev")
	_logLevel := util.GetENV("LOG_LEVEL", "debug")

	appPort = flag.String("app-port", _appPort, "server port\ndefault to 8070 or the value of the APP_PORT env var, if it is set")
	appENV = flag.String("app-env", _logLevel, "server port\nmust be one of [dev, prod]\ndefault to dev or the value of the APP_ENV env var, if it is set")
	logLevel = flag.String("log-level", _appENV, "server port\nmust be one of [debug, info, warn, error, fatal]\ndefault to debug or the value of the LOG_LEVEL env var, if it is set")

	flag.Parse()
}

func main() {
	parseConfig()

	logging.Init(*appENV == "prod", *logLevel)
	stopped := api.NewServer().Start(context.Background(), *appPort)
	<-stopped

	logging.Info("api stopped")
}
