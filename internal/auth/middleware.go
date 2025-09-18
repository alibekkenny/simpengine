package auth

import (
	"context"
	"net/http"
)

type contextKey string

const userContextKey contextKey = "user_id"
const roleContextKey contextKey = "role"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("jwt")
		if err != nil {
			http.Error(w, "missing auth cookie", http.StatusUnauthorized)
			return
		}

		claims, err := ParseJWT(cookie.Value)
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		userID, ok := claims["user_id"].(float64)
		if !ok {
			http.Error(w, "invalid subject", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userContextKey, int64(userID))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserIDFromContext(ctx context.Context) (int64, bool) {
	id, ok := ctx.Value(userContextKey).(int64)
	return int64(id), ok
}

func GetRoleFromContext(ctx context.Context) (string, bool) {
	role, ok := ctx.Value(roleContextKey).(string)
	return role, ok
}
