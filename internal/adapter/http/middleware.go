package adapterhttp

import (
	"context"
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
				WriteErrorJSON(w, http.StatusUnauthorized, "missing authorization header")
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				WriteErrorJSON(w, http.StatusUnauthorized, "invalid authorization header format")
				return
			}

			tokenStr := parts[1]
			claims, err := signer.Parse(tokenStr)
			if err != nil {
				WriteErrorJSON(w, http.StatusUnauthorized, "invalid or expired token")
				return
			}

			ctx := context.WithValue(r.Context(), UserClaimsKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
