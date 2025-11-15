package internal

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"rsoi/internal/api"
	"rsoi/internal/config"
	"rsoi/internal/repo"
	"rsoi/internal/service"

	_ "github.com/lib/pq" // PostgreSQL драйвер
)

type App struct {
	cfg *config.AppConfig
	api *api.RestApi
}

func NewApp(cfg *config.AppConfig) *App {
	return &App{
		cfg: cfg,
	}
}

func (a *App) Run(ctx context.Context) error {
	dsn := fmt.Sprintf("user=%s dbname=%s password=%s host=%s port=%s sslmode=disable", a.cfg.DB.User, a.cfg.DB.DBName, a.cfg.DB.Password, a.cfg.DB.Host, a.cfg.DB.Port)
	log.Println("dsn:", dsn)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	r := repo.NewPostgresRepo(db)
	s := service.NewService(r)
	a.api = api.NewApi(s)
	a.api.Init()
	return a.api.Run(a.cfg.Addr)
}
