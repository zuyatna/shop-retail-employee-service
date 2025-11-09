package router

import (
	"context"
	"net/http"

	httpAdapter "github.com/zuyatna/shop-retail-employee-service/internal/adapter/http"
)

func EmployeeRoutes(ctx context.Context, authHandler *httpAdapter.AuthHandler, empHandler *httpAdapter.EmployeeHandler, authMiddleware *httpAdapter.AuthMiddleware) *http.ServeMux {
	mux := http.NewServeMux()

	// public routes
	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		authHandler.Login(ctx, w, r)
	})

	// ======================================================================
	// protected routes
	// ======================================================================
	// POST /employee
	mux.HandleFunc("/employee", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		authMiddleware.WithAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			empHandler.Create(ctx, w, r)
		})).ServeHTTP(w, r)
	})
	// GET /employees
	mux.HandleFunc("/employees", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		authMiddleware.WithAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			empHandler.List(ctx, w, r)
		})).ServeHTTP(w, r)
	})
	// GET, PUT, DELETE /employee/{id}
	mux.HandleFunc("/employee/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			authMiddleware.WithAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				empHandler.Get(ctx, w, r)
			})).ServeHTTP(w, r)
		case http.MethodPut:
			authMiddleware.WithAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				empHandler.Update(ctx, w, r)
			})).ServeHTTP(w, r)
		case http.MethodDelete:
			authMiddleware.WithAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				empHandler.Delete(ctx, w, r)
			})).ServeHTTP(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})
	// GET /employee/me
	mux.HandleFunc("/employee/me", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		authMiddleware.WithAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			empHandler.GetMe(ctx, w, r)
		})).ServeHTTP(w, r)
	})

	// ======================================================================
	// employee photo routes
	// ======================================================================
	// GET, PUT, DELETE /employee/photo/{id}
	mux.HandleFunc("/employee/photo/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			authMiddleware.WithAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				empHandler.GetPhoto(ctx, w, r)
			})).ServeHTTP(w, r)
		case http.MethodPut:
			authMiddleware.WithAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				empHandler.PutPhotoMultipart(ctx, w, r)
			})).ServeHTTP(w, r)
		case http.MethodDelete:
			authMiddleware.WithAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				empHandler.DeletePhoto(ctx, w, r)
			})).ServeHTTP(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	return mux
}
