package app

import (
	"context"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zuyatna/shop-retail-employee-service/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
)

type App struct {
	Pool    *pgxpool.Pool
	MongoDB *mongo.Database
	Router  http.Handler
}

func New(cfg *config.Config) (*App, error) {
	pool, err := initPostgres(cfg)
	if err != nil {
		return nil, err
	}

	mongoDB, err := initMongo(cfg)
	if err != nil {
		return nil, err
	}

	handler := NewHandler(pool, mongoDB, cfg)

	app := &App{
		Pool:    pool,
		MongoDB: mongoDB,
		Router:  handler,
	}

	return app, nil
}

func (a *App) Close() {
	log.Println("closing database connection")
	a.Pool.Close()

	if a.MongoDB != nil {
		if err := a.MongoDB.Client().Disconnect(context.Background()); err != nil {
			log.Printf("error disconnecting MongoDB client: %v", err)
		}
	}
}
