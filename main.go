package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/omegaatt36/gotasker/api"
	"github.com/omegaatt36/gotasker/logging"
	"github.com/omegaatt36/gotasker/persistance/database"
	"github.com/omegaatt36/gotasker/util"
)

var (
	appPort  *string
	logLevel *string
	appENV   *string

	redisHost     *string
	redisPort     *string
	redisPassword *string
)

func parseConfig() {
	_appPort := util.GetENV("APP_PORT", "8070")
	_appENV := util.GetENV("APP_ENV", "dev")
	_logLevel := util.GetENV("LOG_LEVEL", "debug")
	_redisHost := util.GetENV("REDIS_HOST", "localhost")
	_redisPort := util.GetENV("REDIS_PORT", "6379")
	_redisPassword := util.GetENV("REDIS_PASSWORD", "")

	appPort = flag.String("app-port", _appPort, "server port\ndefault to 8070 or the value of the APP_PORT env var, if it is set")
	appENV = flag.String("app-env", _logLevel, "app env\nmust be one of [dev, prod]\ndefault to dev or the value of the APP_ENV env var, if it is set")
	logLevel = flag.String("log-level", _appENV, "log level\nmust be one of [debug, info, warn, error, fatal]\ndefault to debug or the value of the LOG_LEVEL env var, if it is set")
	redisHost = flag.String("redis-host", _redisHost, "redis host\ndefault to localhost or the value of the REDIS_HOST env var, if it is set")
	redisPort = flag.String("redis-port", _redisPort, "redis port\ndefault to 6379 or the value of the REDIS_PORT env var, if it is set")
	redisPassword = flag.String("redis-password", _redisPassword, "redis port\ndefault to 6379 or the value of the REDIS_PASSWORD env var, if it is set")

	flag.Parse()
}

func main() {
	ctx := context.Background()

	parseConfig()

	logging.Init(*appENV == "prod", *logLevel)

	database.Initialize(ctx, fmt.Sprintf("%s:%s", *redisHost, *redisPort), *redisPassword)

	stopped := api.NewServer().Start(ctx, *appPort)
	<-stopped

	logging.Info("api stopped")
}
