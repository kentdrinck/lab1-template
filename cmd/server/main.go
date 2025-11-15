package main

import (
	"context"
	"log"
	"rsoi/internal"
	"rsoi/internal/config"
)

func main() {
	cfg, err := config.LoadConfig("config/development.toml")
	if err != nil {
		log.Panic(err)
	}
	ctx := context.Background()
	app := internal.NewApp(cfg)
	if err := app.Run(ctx); err != nil {
		log.Panic(err)
	}
}
