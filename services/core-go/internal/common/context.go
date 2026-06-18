package common

import (
	"context"
	"net/http"
	permDomain "github.com/user/seta-dam-backend/services/core-go/internal/permission/domain"
)

type contextKey string

const UserContextKey contextKey = "user_context"

func ExtractUserContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get("X-User-ID")
		role := r.Header.Get("X-User-Role")
		username := r.Header.Get("X-User-Username")

		if userID != "" && role != "" {
			uCtx := &permDomain.UserContext{
				UserID:   userID,
				Role:     role,
				Username: username,
			}
			ctx := context.WithValue(r.Context(), UserContextKey, uCtx)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		next.ServeHTTP(w, r)
	})
}

func GetUserContext(ctx context.Context) *permDomain.UserContext {
	val := ctx.Value(UserContextKey)
	if val == nil {
		return nil
	}
	return val.(*permDomain.UserContext)
}
