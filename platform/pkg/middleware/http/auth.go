package http

import (
	"context"
	"net/http"

	grpcAuth "github.com/pinai4/spaceship-factory/platform/pkg/middleware/grpc"
	authV1 "github.com/pinai4/spaceship-factory/shared/pkg/proto/auth/v1"
	commonV1 "github.com/pinai4/spaceship-factory/shared/pkg/proto/common/v1"
)

const SessionUUIDHeader = "X-Session-Uuid"

// IAMClient is an alias for the generated gRPC client
type IAMClient = authV1.AuthServiceClient

// AuthMiddleware is the middleware for HTTP request authentication
type AuthMiddleware struct {
	iamClient IAMClient
}

// NewAuthMiddleware creates a new authentication middleware
func NewAuthMiddleware(iamClient IAMClient) *AuthMiddleware {
	return &AuthMiddleware{
		iamClient: iamClient,
	}
}

// Handle processes an HTTP request with authentication
func (m *AuthMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract session UUID from header
		sessionUUID := r.Header.Get(SessionUUIDHeader)
		if sessionUUID == "" {
			writeErrorResponse(w, http.StatusUnauthorized, "MISSING_SESSION", "Authentication required")
			return
		}

		// Validate the session through the IAM service
		whoamiRes, err := m.iamClient.Whoami(r.Context(), &authV1.WhoamiRequest{
			SessionUuid: sessionUUID,
		})
		if err != nil {
			writeErrorResponse(w, http.StatusUnauthorized, "INVALID_SESSION", "Authentication failed")
			return
		}

		// Add the user and session UUID to context using functions from grpc middleware
		ctx := r.Context()
		ctx = grpcAuth.AddSessionUUIDToContext(ctx, sessionUUID)
		// Also add the user to context
		ctx = context.WithValue(ctx, grpcAuth.GetUserContextKey(), whoamiRes.User)

		// Pass control to the next handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserFromContext extracts the user from the context
func GetUserFromContext(ctx context.Context) (*commonV1.User, bool) {
	return grpcAuth.GetUserFromContext(ctx)
}

// GetUserIDFromContext extracts the user id from the context
func GetUserIDFromContext(ctx context.Context) string {
	var id string
	u, ok := grpcAuth.GetUserFromContext(ctx)
	if ok {
		id = u.GetUuid()
	}

	return id
}

// GetSessionUUIDFromContext extracts the session UUID from the context
func GetSessionUUIDFromContext(ctx context.Context) (string, bool) {
	return grpcAuth.GetSessionUUIDFromContext(ctx)
}
