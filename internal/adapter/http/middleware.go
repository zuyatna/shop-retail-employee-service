package http

import (
	"context"
	"net/http"
	"strings"

	"github.com/zuyatna/shop-retail-employee-service/internal/domain"
	jwtutils "github.com/zuyatna/shop-retail-employee-service/internal/utils/jwt"
)

type ctxKey string

const (
	ctxUserID ctxKey = "uid"
	ctxEmail  ctxKey = "email"
	ctxRole   ctxKey = "role"
)

type AuthMiddleware struct {
	parser interface {
		Parse(tokenStr string) (*jwtutils.Claims, error)
	}
}

func NewAuthMiddleware(parser interface {
	Parse(tokenStr string) (*jwtutils.Claims, error)
}) *AuthMiddleware {
	return &AuthMiddleware{parser: parser}
}

func (a *AuthMiddleware) WithAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
			writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "missing or invalid authorization header"})
			return
		}
		tokenStr := strings.TrimSpace(authHeader[len("Bearer "):])
		claims, err := a.parser.Parse(tokenStr)
		if err != nil {
			writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid token: " + err.Error()})
			return
		}
		ctx := context.WithValue(r.Context(), ctxUserID, claims.UserID)
		ctx = context.WithValue(ctx, ctxEmail, claims.Email)
		ctx = context.WithValue(ctx, ctxRole, claims.Role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getCallerRoleFromContext(r *http.Request) domain.Role {
	if v := r.Context().Value(ctxRole); v != nil {
		if roleStr, ok := v.(string); ok {
			return domain.Role(roleStr)
		}
		if role, ok := v.(domain.Role); ok {
			return role
		}
	}
	return domain.RoleStaff
}

func getCalledIDFromContext(r *http.Request) string {
	if v := r.Context().Value(ctxUserID); v != nil {
		if id, ok := v.(string); ok {
			return id
		}
	}
	return ""
}
