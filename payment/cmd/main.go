package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	paymentV1 "github.com/pinai4/spaceship-factory/shared/pkg/proto/payment/v1"
)

const grpcPort = 50052

type paymentService struct {
	paymentV1.UnimplementedPaymentServiceServer
}

func (s *paymentService) PayOrder(ctx context.Context, request *paymentV1.PayOrderRequest) (*paymentV1.PayOrderResponse, error) {
	tranID := uuid.New()
	return &paymentV1.PayOrderResponse{TransactionUuid: tranID.String()}, nil
}

func PrinterInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		resp, err := handler(ctx, req)
		if err == nil {
			if r, ok := resp.(*paymentV1.PayOrderResponse); ok && r != nil {
				log.Printf("Payment was successful, %s\n", r.GetTransactionUuid())
				// log.Printf("Request: %#v\n", req)
			}
		}

		return resp, err
	}
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
		return
	}

	// Create GRPC server
	s := grpc.NewServer(grpc.ChainUnaryInterceptor(
		PrinterInterceptor(),
	))

	// Register our service
	service := &paymentService{}

	paymentV1.RegisterPaymentServiceServer(s, service)

	// Enable GRPC reflection to simplify debugging
	reflection.Register(s)

	go func() {
		log.Printf("🚀 gRPC server listening on %d\n", grpcPort)
		if err := s.Serve(lis); err != nil {
			log.Printf("failed to serve: %v\n", err)
			return
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down gRPC server...")
	s.GracefulStop()
	log.Println("✅ Server stopped")
}
