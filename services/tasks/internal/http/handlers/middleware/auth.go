package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/student/tech-ip-sem2/services/tasks/internal/grpcclient"
	"github.com/student/tech-ip-sem2/shared/middleware"
)

func AuthMiddleware(authClient *grpcclient.AuthClient) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, `{"error":"missing authorization header"}`, http.StatusUnauthorized)
				return
			}
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				http.Error(w, `{"error":"invalid authorization format"}`, http.StatusUnauthorized)
				return
			}
			token := parts[1]

			requestID := middleware.GetRequestID(r.Context())

			valid, _, err := authClient.Verify(r.Context(), token, requestID)
			if err != nil {
				if errors.Is(err, context.DeadlineExceeded) {
					http.Error(w, `{"error":"auth service timeout"}`, http.StatusGatewayTimeout)
					return
				}
				http.Error(w, `{"error":"authorization service unavailable"}`, http.StatusServiceUnavailable)
				return
			}
			if !valid {
				http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
