package grpc

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	authV1 "github.com/pinai4/spaceship-factory/shared/pkg/proto/auth/v1"
	commonV1 "github.com/pinai4/spaceship-factory/shared/pkg/proto/common/v1"
)

const (
	// SessionUUIDMetadataKey is the key for passing the session UUID in gRPC metadata
	SessionUUIDMetadataKey = "session-uuid"
)

type contextKey string

const (
	// userContextKey is the key for storing the user in the context
	userContextKey contextKey = "user"
	// sessionUUIDContextKey is the key for storing the session UUID in the context
	sessionUUIDContextKey contextKey = "session-uuid"
)

// IAMClient is an alias for the generated gRPC client
type IAMClient = authV1.AuthServiceClient

// AuthInterceptor is the interceptor for gRPC authentication
type AuthInterceptor struct {
	iamClient IAMClient
}

// NewAuthInterceptor creates a new authentication interceptor
func NewAuthInterceptor(iamClient IAMClient) *AuthInterceptor {
	return &AuthInterceptor{
		iamClient: iamClient,
	}
}

// Unary returns a unary server interceptor for authentication
func (i *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		authCtx, err := i.authenticate(ctx)
		if err != nil {
			return nil, err
		}

		return handler(authCtx, req)
	}
}

// authenticate performs authentication and adds the user to the context
func (i *AuthInterceptor) authenticate(ctx context.Context) (context.Context, error) {
	// Extract metadata from context
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "missing metadata")
	}

	// Get session UUID from metadata
	sessionUUIDs := md.Get(SessionUUIDMetadataKey)
	if len(sessionUUIDs) == 0 {
		return nil, status.Error(codes.Unauthenticated, "missing session-uuid in metadata")
	}

	sessionUUID := sessionUUIDs[0]
	if sessionUUID == "" {
		return nil, status.Error(codes.Unauthenticated, "empty session-uuid")
	}

	// Validate session via IAM service
	whoamiRes, err := i.iamClient.Whoami(ctx, &authV1.WhoamiRequest{
		SessionUuid: sessionUUID,
	})
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, fmt.Sprintf("invalid session: %v", err))
	}

	// Add the user and session UUID to the context
	authCtx := context.WithValue(ctx, userContextKey, whoamiRes.User)
	authCtx = context.WithValue(authCtx, sessionUUIDContextKey, sessionUUID)
	return authCtx, nil
}

// GetUserFromContext retrieves the user from the context
func GetUserFromContext(ctx context.Context) (*commonV1.User, bool) {
	user, ok := ctx.Value(userContextKey).(*commonV1.User)
	return user, ok
}

// GetUserContextKey returns the context key for the user
func GetUserContextKey() contextKey {
	return userContextKey
}

// GetSessionUUIDFromContext retrieves the session UUID from the context
func GetSessionUUIDFromContext(ctx context.Context) (string, bool) {
	sessionUUID, ok := ctx.Value(sessionUUIDContextKey).(string)
	return sessionUUID, ok
}

// AddSessionUUIDToContext adds the session UUID to the context
func AddSessionUUIDToContext(ctx context.Context, sessionUUID string) context.Context {
	return context.WithValue(ctx, sessionUUIDContextKey, sessionUUID)
}

// ForwardSessionUUIDToGRPC adds the session UUID from the context to outgoing gRPC metadata
func ForwardSessionUUIDToGRPC(ctx context.Context) context.Context {
	sessionUUID, ok := GetSessionUUIDFromContext(ctx)
	if !ok || sessionUUID == "" {
		return ctx
	}

	return metadata.AppendToOutgoingContext(ctx, SessionUUIDMetadataKey, sessionUUID)
}
