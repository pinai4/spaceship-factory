package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"

	authV1 "github.com/pinai4/spaceship-factory/shared/pkg/proto/auth/v1"
	commonV1 "github.com/pinai4/spaceship-factory/shared/pkg/proto/common/v1"
)

const (
	grpcAddress = "0.0.0.0:50051"
	userID      = "00000000-0000-0000-0000-111111111111"
	sessionID   = "00000000-0000-0000-0000-222222222222"
)

type mockAuthServer struct {
	authV1.UnimplementedAuthServiceServer
}

func (s *mockAuthServer) Whoami(_ context.Context, _ *authV1.WhoamiRequest) (*authV1.WhoamiResponse, error) {
	now := time.Now()
	return &authV1.WhoamiResponse{
		Session: &commonV1.Session{
			Uuid:      sessionID,
			CreatedAt: timestamppb.New(now),
			UpdatedAt: timestamppb.New(now.Add(24 * time.Hour)),
		},
		User: &commonV1.User{
			Uuid: userID,
			Info: &commonV1.UserInfo{
				Login: "test_login",
				Email: "test_email",
			},
			CreatedAt: timestamppb.New(now),
		},
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", grpcAddress) //nolint:gosec
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
		return
	}
	grpcServer := grpc.NewServer()
	authV1.RegisterAuthServiceServer(grpcServer, &mockAuthServer{})

	reflection.Register(grpcServer)

	go func() {
		log.Printf("🚀 gRPC server listening on %s\n", grpcAddress)
		if err := grpcServer.Serve(lis); err != nil {
			log.Printf("failed to serve: %v\n", err)
			return
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down gRPC server...")
	grpcServer.GracefulStop()
	log.Println("✅ Server stopped")
}
