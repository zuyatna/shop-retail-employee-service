package router

import (
	"net/http"

	httpAdapter "github.com/zuyatna/shop-retail-employee-service/internal/adapter/http"
)

func EmployeeRoutes(authHandler *httpAdapter.AuthHandler, empHandler *httpAdapter.EmployeeHandler, authMiddleware *httpAdapter.AuthMiddleware) *http.ServeMux {
	mux := http.NewServeMux()

	// public routes
	mux.HandleFunc("POST /login", authHandler.Login)

	// protected routes
	mux.Handle("POST /employee", authMiddleware.WithAuth(http.HandlerFunc(empHandler.Create)))
	mux.Handle("GET /employees", authMiddleware.WithAuth(http.HandlerFunc(empHandler.List)))
	mux.Handle("GET /employee/", authMiddleware.WithAuth(http.HandlerFunc(empHandler.Get)))
	mux.Handle("GET /employee/me", authMiddleware.WithAuth(http.HandlerFunc(empHandler.GetMe)))
	mux.Handle("PUT /employee/", authMiddleware.WithAuth(http.HandlerFunc(empHandler.Update)))
	mux.Handle("DELETE /employee/", authMiddleware.WithAuth(http.HandlerFunc(empHandler.Delete)))

	// employee photo routes
	mux.Handle("GET /employee/photo/", authMiddleware.WithAuth(http.HandlerFunc(empHandler.GetPhoto)))
	mux.Handle("PUT /employee/photo/", authMiddleware.WithAuth(http.HandlerFunc(empHandler.PutPhotoMultipart)))
	mux.Handle("DELETE /employee/photo/", authMiddleware.WithAuth(http.HandlerFunc(empHandler.DeletePhoto)))

	return mux
}
