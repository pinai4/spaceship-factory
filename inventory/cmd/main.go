package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	inventoryV1API "github.com/pinai4/spaceship-factory/inventory/internal/api/inventory/v1"
	"github.com/pinai4/spaceship-factory/inventory/internal/model"
	"github.com/pinai4/spaceship-factory/inventory/internal/repository"
	partRepository "github.com/pinai4/spaceship-factory/inventory/internal/repository/part"
	partService "github.com/pinai4/spaceship-factory/inventory/internal/service/part"
	inventoryV1 "github.com/pinai4/spaceship-factory/shared/pkg/proto/inventory/v1"
)

const grpcPort = 50051

func seed(repo repository.PartRepository) error {
	now := time.Now()

	part1 := model.Part{
		UUID:          uuid.Nil.String(),
		Name:          "Turbo Engine X200",
		Description:   "High-performance turbo engine suitable for small aircraft.",
		Price:         125000.50,
		StockQuantity: 8,
		Category:      model.CategoryEngine,
		Dimensions: model.Dimensions{
			Length: 120.5,
			Width:  80.2,
			Height: 95.3,
			Weight: 450.0,
		},
		Manufacturer: model.Manufacturer{
			Name:    "AeroTech Industries",
			Country: "USA",
			Website: "https://aerotech.example.com",
		},
		Tags: []string{"engine", "turbo", "aircraft"},
		Metadata: map[string]any{
			"power_kw":      980.5,
			"certified":     true,
			"serial_number": "SN-ENGX200-001",
		},
		CreatedAt: now,
	}
	if err := repo.Add(context.Background(), part1); err != nil {
		return err
	}

	part2 := model.Part{
		UUID:          uuid.NewString(),
		Name:          "Titanium Wing Panel",
		Description:   "Lightweight titanium alloy wing panel with anti-corrosion coating.",
		Price:         32000.0,
		StockQuantity: 25,
		Category:      model.CategoryWing,
		Dimensions: model.Dimensions{
			Length: 250.0,
			Width:  60.0,
			Height: 5.0,
			Weight: 120.0,
		},
		Manufacturer: model.Manufacturer{
			Name:    "SkyMetal Works",
			Country: "Germany",
			Website: "https://skymetal.example.com",
		},
		Tags: []string{"wing", "titanium", "aircraft"},
		Metadata: map[string]any{
			"material":     "Titanium Alloy",
			"batch_number": 20241001,
			"is_tested":    true,
		},
		CreatedAt: now,
	}
	if err := repo.Add(context.Background(), part2); err != nil {
		return err
	}

	return nil
}

func main() {
	// Init repo
	repo := partRepository.NewRepository()
	if err := seed(repo); err != nil {
		log.Printf("failed to seed: %v\n", err)
		return
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
		return
	}

	// Create GRPC server
	s := grpc.NewServer()

	// Register our service
	service := partService.NewService(repo)
	api := inventoryV1API.NewAPI(service)

	inventoryV1.RegisterInventoryServiceServer(s, api)

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
