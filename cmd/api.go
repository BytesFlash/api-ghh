package cmd

import (
	"github.com/Unleash/unleash-client-go/v3"
	"github.com/imedcl/manager-api/pkg/config"
	"github.com/imedcl/manager-api/pkg/routes"
)

func Start(client *unleash.Client) {
	cfg := config.New()
	routes.Create(cfg, client)
}
