package app

import (
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zuyatna/shop-retail-employee-service/internal/config"
)

type App struct {
	Pool   *pgxpool.Pool
	Router http.Handler // nanti diisi di app/http.go
}

func New(cfg *config.Config) (*App, error) {
	pool, err := initPostgres(cfg)
	if err != nil {
		return nil, err
	}

	handler := NewHandler(pool, cfg)

	app := &App{
		Pool:   pool,
		Router: handler,
	}

	return app, nil
}

func (a *App) Close() {
	log.Println("closing database connection")
	a.Pool.Close()
}
