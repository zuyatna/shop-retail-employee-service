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
)

func NewHandler(pool *pgxpool.Pool, cfg *config.Config) http.Handler {
	idGenerator := idgen.NewUUIDv7Generator()

	employeeRepo := repo.NewPostgresEmployeeRepo(pool)

	ctxTimeout := 5 * time.Second // Example timeout, can be from config
	employeeUsecase := usecase.NewEmployeeUsecase(employeeRepo, idGenerator, ctxTimeout)

	employeeHandler := adapterhttp.NewEmployeeHandler(employeeUsecase)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /employees", employeeHandler.Register)
	mux.HandleFunc("GET /employees/{id}", employeeHandler.GetByID)
	mux.HandleFunc("PATCH /employees/{id}", employeeHandler.Update)
	mux.HandleFunc("DELETE /employees/{id}", employeeHandler.Delete)

	// Middlewares can be added here like logging, recovery, etc.

	return mux
}
