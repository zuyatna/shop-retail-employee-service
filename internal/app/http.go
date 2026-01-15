package app

import (
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	adapterhttp "github.com/zuyatna/shop-retail-employee-service/internal/adapter/http"
	"github.com/zuyatna/shop-retail-employee-service/internal/adapter/repo"
	"github.com/zuyatna/shop-retail-employee-service/internal/config"
	"github.com/zuyatna/shop-retail-employee-service/internal/usecase"
	"github.com/zuyatna/shop-retail-employee-service/internal/util/idgen"
	"github.com/zuyatna/shop-retail-employee-service/internal/util/jwtutil"
)

func NewHandler(pool *pgxpool.Pool, cfg *config.Config) http.Handler {
	idGenerator := idgen.NewUUIDv7Generator()

	jwtSigner := &jwtutil.Signer{
		Secret: []byte(cfg.JWTSecret),
		Issuer: cfg.JWTIssuer,
		TTL:    time.Duration(cfg.JWTTTL),
	}

	employeeRepo := repo.NewPostgresEmployeeRepo(pool)

	ctxTimeout := 5 * time.Second // Example timeout, can be from config
	employeeUsecase := usecase.NewEmployeeUsecase(employeeRepo, idGenerator, ctxTimeout)
	authUsecase := usecase.NewAuthUsecase(employeeRepo, jwtSigner, ctxTimeout)

	employeeHandler := adapterhttp.NewEmployeeHandler(employeeUsecase)
	authHandler := adapterhttp.NewAuthHandler(authUsecase)

	authMiddleware := adapterhttp.AuthMiddleware(jwtSigner)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /auth/login", authHandler.Login)

	mux.HandleFunc("POST /employees", authMiddleware(http.HandlerFunc(employeeHandler.Register)).ServeHTTP)
	mux.HandleFunc("GET /employees/{id}", authMiddleware(http.HandlerFunc(employeeHandler.GetByID)).ServeHTTP)
	mux.HandleFunc("PATCH /employees/{id}", authMiddleware(http.HandlerFunc(employeeHandler.Update)).ServeHTTP)
	mux.HandleFunc("DELETE /employees/{id}", authMiddleware(http.HandlerFunc(employeeHandler.Delete)).ServeHTTP)

	return mux
}
