package app

import (
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	adapterhttp "github.com/zuyatna/shop-retail-employee-service/internal/adapter/http"
	"github.com/zuyatna/shop-retail-employee-service/internal/adapter/repo"
	"github.com/zuyatna/shop-retail-employee-service/internal/adapter/storage"
	"github.com/zuyatna/shop-retail-employee-service/internal/config"
	"github.com/zuyatna/shop-retail-employee-service/internal/domain"
	"github.com/zuyatna/shop-retail-employee-service/internal/usecase"
	"github.com/zuyatna/shop-retail-employee-service/internal/util/clock"
	"github.com/zuyatna/shop-retail-employee-service/internal/util/idgen"
	"github.com/zuyatna/shop-retail-employee-service/internal/util/jwtutil"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewHandler(pool *pgxpool.Pool, mongoDB *mongo.Database, cfg *config.Config) http.Handler {
	idGenerator := idgen.NewUUIDv7Generator()

	realClock := clock.RealClock{}

	jwtSigner := &jwtutil.Signer{
		Secret: []byte(cfg.JWTSecret),
		Issuer: cfg.JWTIssuer,
		TTL:    time.Duration(cfg.JWTTTL) * time.Second,
	}

	employeeRepo := repo.NewPostgresEmployeeRepo(pool)
	attendanceRepo := repo.NewMongoAttendanceRepo(mongoDB)

	minioStorage, err := storage.NewMinioStorage(cfg)
	if err != nil {
		panic(err)
	}

	ctxTimeout := 5 * time.Second // Example timeout, can be from config

	employeeUsecase := usecase.NewEmployeeUsecase(employeeRepo, minioStorage, idGenerator, ctxTimeout)
	authUsecase := usecase.NewAuthUsecase(employeeRepo, jwtSigner, ctxTimeout)
	attendanceUsecase := usecase.NewAttendanceUsecase(attendanceRepo, employeeRepo, idGenerator, cfg, realClock, ctxTimeout)

	employeeHandler := adapterhttp.NewEmployeeHandler(employeeUsecase)
	authHandler := adapterhttp.NewAuthHandler(authUsecase)
	attendanceHandler := adapterhttp.NewAttendanceHandler(attendanceUsecase)

	authMiddleware := adapterhttp.AuthMiddleware(jwtSigner)

	requirePrivileged := adapterhttp.RoleMiddleware(string(domain.RoleAdmin), string(domain.RoleSupervisor))
	requireAllRoles := adapterhttp.RoleMiddleware(string(domain.RoleAdmin), string(domain.RoleSupervisor), string(domain.RoleStaff))

	mux := http.NewServeMux()

	mux.HandleFunc("POST /auth/login", authHandler.Login)

	mux.HandleFunc("GET /employees/me", authMiddleware(requireAllRoles(http.HandlerFunc(employeeHandler.GetMe))).ServeHTTP)
	mux.HandleFunc("POST /employees", authMiddleware(requirePrivileged(http.HandlerFunc(employeeHandler.Register))).ServeHTTP)
	mux.HandleFunc("GET /employees", authMiddleware(requirePrivileged(http.HandlerFunc(employeeHandler.GetAll))).ServeHTTP)
	mux.HandleFunc("GET /employees/{id}", authMiddleware(requirePrivileged(http.HandlerFunc(employeeHandler.GetByID))).ServeHTTP)
	mux.HandleFunc("PATCH /employees/{id}", authMiddleware(requirePrivileged(http.HandlerFunc(employeeHandler.Update))).ServeHTTP)
	mux.HandleFunc("POST /employees/{id}/photo", authMiddleware(requireAllRoles(http.HandlerFunc(employeeHandler.UploadPhoto))).ServeHTTP)
	mux.HandleFunc("DELETE /employees/{id}", authMiddleware(requirePrivileged(http.HandlerFunc(employeeHandler.Delete))).ServeHTTP)

	mux.HandleFunc("POST /attendances/checkin", authMiddleware(requireAllRoles(http.HandlerFunc(attendanceHandler.CheckIn))).ServeHTTP)
	mux.HandleFunc("POST /attendances/checkout", authMiddleware(requireAllRoles(http.HandlerFunc(attendanceHandler.CheckOut))).ServeHTTP)

	return mux
}
