package adapterhttp

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/zuyatna/shop-retail-employee-service/internal/util/jwtutil"
)

type contextKey string

const UserClaimsKey contextKey = "user_claims"

func AuthMiddleware(signer *jwtutil.Signer) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				WriteErrorJSON(w, http.StatusUnauthorized, errors.New("missing authorization header"), "missing authorization header")
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				WriteErrorJSON(w, http.StatusUnauthorized, errors.New("invalid authorization header format"), "invalid authorization header format")
				return
			}

			tokenStr := parts[1]
			claims, err := signer.Parse(tokenStr)
			if err != nil {
				WriteErrorJSON(w, http.StatusUnauthorized, err, "invalid or expired token")
				return
			}

			ctx := context.WithValue(r.Context(), UserClaimsKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RoleMiddleware(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value(UserClaimsKey).(*jwtutil.Claims)
			if !ok || claims == nil {
				WriteErrorJSON(w, http.StatusUnauthorized, errors.New("unauthorized"), "user claims not found")
				return
			}

			isAllowed := false
			for _, role := range allowedRoles {
				if role == claims.Role {
					isAllowed = true
					break
				}
			}

			if !isAllowed {
				WriteErrorJSON(w, http.StatusForbidden, errors.New("forbidden"), "access denied: insufficient permissions")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
