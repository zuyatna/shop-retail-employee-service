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
	"github.com/zuyatna/shop-retail-employee-service/internal/usecase"
	"github.com/zuyatna/shop-retail-employee-service/internal/utils/idgen"
	jwtutils "github.com/zuyatna/shop-retail-employee-service/internal/utils/jwt"
)

func main() {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Println("No .env file found, fallback to system env")
	}

	ctx := context.Background()
	cfg := config.Load()

	pool, err := pgxpool.New(ctx, cfg.DatabaseURL())
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
	signer := &jwtutils.Signer{Secret: []byte(cfg.JWTSecret), Issuer: cfg.JWTIssuer, TTL: time.Duration(cfg.JWTTTL) * time.Second}
	authUsecase := usecase.NewAuthUsecase(pgRepo, signer)

	// handler
	empHandler := httpAdapter.NewEmployeeHandler(empUsecase)
	authHandler := httpAdapter.NewAuthHandler(authUsecase)

	// middleware
	authMiddleware := httpAdapter.NewAuthMiddleware(signer)

	mux := http.NewServeMux()

	// public routes
	mux.HandleFunc("POST /login", authHandler.Login)

	// protected routes
	mux.Handle("POST /employee", authMiddleware.WithAuth(http.HandlerFunc(empHandler.Create)))
	mux.Handle("GET /employees", authMiddleware.WithAuth(http.HandlerFunc(empHandler.List)))
	mux.Handle("GET /employee/", authMiddleware.WithAuth(http.HandlerFunc(empHandler.Get)))
	mux.Handle("PUT /employee/", authMiddleware.WithAuth(http.HandlerFunc(empHandler.Update)))
	mux.Handle("DELETE /employee/", authMiddleware.WithAuth(http.HandlerFunc(empHandler.Delete)))

	server := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("Starting server on %s\n", cfg.HTTPAddr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on %s: %v\n", cfg.HTTPAddr, err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("shutdown error: %v", err)
	}
	log.Println("Server gracefully stopped")
}
