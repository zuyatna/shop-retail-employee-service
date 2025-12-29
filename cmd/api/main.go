package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	httpAdapter "github.com/zuyatna/shop-retail-employee-service/internal/adapter/http"
	"github.com/zuyatna/shop-retail-employee-service/internal/adapter/repo"
	"github.com/zuyatna/shop-retail-employee-service/internal/config"
	"github.com/zuyatna/shop-retail-employee-service/internal/router"
	"github.com/zuyatna/shop-retail-employee-service/internal/usecase"
	"github.com/zuyatna/shop-retail-employee-service/internal/util/idgen"
	"github.com/zuyatna/shop-retail-employee-service/internal/util/jwtutil"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found, fallback to system env")
	}

	ctx := context.Background()

	// database config
	cfg, _ := pgxpool.ParseConfig(config.Load().DatabaseURL())
	cfg.MaxConns = 50 // adjust based on application's needs
	cfg.MinConns = 5
	cfg.MaxConnLifetime = 30 * time.Minute
	cfg.HealthCheckPeriod = 1 * time.Minute

	// database connection
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer pool.Close()

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := pool.Ping(pingCtx); err != nil {
		log.Fatalf("Database ping failed: %v\n", err)
	}

	// repo & usecase
	pgRepo := repo.NewPostgresEmployeeRepo(pool, 5*time.Second)
	uuidGen := idgen.NewUUIDv7Generator()
	empUsecase := usecase.NewEmployeeUsecase(pgRepo, uuidGen)
	signer := &jwtutil.Signer{Secret: []byte(config.Load().JWTSecret), Issuer: config.Load().JWTIssuer, TTL: time.Duration(config.Load().JWTTTL) * time.Second}
	authUsecase := usecase.NewAuthUsecase(pgRepo, signer)

	// handler
	empHandler := httpAdapter.NewEmployeeHandler(empUsecase)
	authHandler := httpAdapter.NewAuthHandler(authUsecase)

	// middleware
	authMiddleware := httpAdapter.NewAuthMiddleware(signer)

	mux := router.EmployeeRoutes(ctx, authHandler, empHandler, authMiddleware)

	server := &http.Server{
		Addr:              config.Load().HTTPAddr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	go func() {
		log.Printf("Starting server on %s\n", config.Load().HTTPAddr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on %s: %v\n", config.Load().HTTPAddr, err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Shutdown error: %v", err)
	}
	log.Println("Server gracefully stopped")
}
