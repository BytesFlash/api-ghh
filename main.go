package main

import (
	"fmt"

	"github.com/Unleash/unleash-client-go/v3"
	api "github.com/imedcl/manager-api/cmd"
	"github.com/imedcl/manager-api/pkg/config"
)

type metricsInterface struct {
}

var (
	client *unleash.Client
	err    error
)

// @title Autentia Manager Documentation
// @version 1
// @Description api autentia manager

// @securityDefinitions.apikey barerToken
// @in header
// @name Authorization

func main() {
	cfg := config.New()
	client, err = unleash.NewClient(
		unleash.WithListener(&unleash.DebugListener{}),
		unleash.WithInstanceId("ELX9hN749PGHkuqrSfCt"),
		unleash.WithAppName(cfg.Environment()), // Set to the running environment of your application
		unleash.WithListener(&metricsInterface{}),
	)
	if err != nil {
		fmt.Println("ERROR", err)
	}
	api.Start(client)
}
