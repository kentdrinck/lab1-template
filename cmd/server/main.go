package main

import (
	"context"
	"flag"
	"log"
	"rsoi/internal"
	"rsoi/internal/config"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "config/k8s-dev.toml", "config")
	flag.Parse()

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Panic(err)
	}
	ctx := context.Background()
	app := internal.NewApp(cfg)
	if err := app.Run(ctx); err != nil {
		log.Panic(err)
	}
}
