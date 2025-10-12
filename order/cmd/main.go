package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderV1API "github.com/pinai4/spaceship-factory/order/internal/api/order/v1"
	inventoryV1Client "github.com/pinai4/spaceship-factory/order/internal/client/grpc/inventory/v1"
	paymentV1Client "github.com/pinai4/spaceship-factory/order/internal/client/grpc/payment/v1"
	orderRepository "github.com/pinai4/spaceship-factory/order/internal/repository/order"
	orderService "github.com/pinai4/spaceship-factory/order/internal/service/order"
	orderV1 "github.com/pinai4/spaceship-factory/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/pinai4/spaceship-factory/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/pinai4/spaceship-factory/shared/pkg/proto/payment/v1"
)

const (
	httpPort = "8080"
	// HTTP server timeouts
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second

	inventoryServerAddress = "localhost:50051"
	paymentServerAddress   = "localhost:50052"
)

func main() {
	////////////////
	////////////////
	conn1, err := grpc.NewClient(
		inventoryServerAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("failed to connect: %v\n", err)
		return
	}
	defer func() {
		if cerr := conn1.Close(); cerr != nil {
			log.Printf("failed to close connect: %v", cerr)
		}
	}()

	// Create inventory gRPC client
	genInventoryClient := inventoryV1.NewInventoryServiceClient(conn1)

	////////////////
	conn2, err := grpc.NewClient(
		paymentServerAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("failed to connect: %v\n", err)
		return
	}
	defer func() {
		if cerr := conn2.Close(); cerr != nil {
			log.Printf("failed to close connect: %v", cerr)
		}
	}()

	// Create payment gRPC client
	genPaymentClient := paymentV1.NewPaymentServiceClient(conn2)
	//////////////////
	//////////////////

	repo := orderRepository.NewRepository()
	service := orderService.NewService(
		repo,
		paymentV1Client.NewClient(genPaymentClient),
		inventoryV1Client.NewClient(genInventoryClient),
	)
	api := orderV1API.NewAPI(service)

	// Create OpenAPI server
	orderServer, err := orderV1.NewServer(api)
	if err != nil {
		log.Printf("OpenAPI server creation error: %v", err)
		return
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))

	// Mount OpenAPI handler
	r.Mount("/", orderServer)

	server := &http.Server{
		Addr:              net.JoinHostPort("localhost", httpPort),
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout, // Protection against Slowloris attacks - a type of DDoS attack in which
		// the attacker deliberately sends HTTP headers slowly, keeping connections open and exhausting
		// the pool of available connections on the server. ReadHeaderTimeout forcibly closes the connection
		// if the client fails to send all headers within the allotted time.
	}

	// Run HTTP server in separate goroutine
	go func() {
		log.Printf("🚀 HTTP-server has been run on port %s\n", httpPort)
		err = server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("❌ HTTP-server running error: %v\n", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Printf("❌ Shutting down server error: %v\n", err)
	}

	log.Println("✅ Server has stopped")
}
