package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"slices"
	"sync"
	"syscall"
	"time"

	"github.com/go-faster/errors"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	inventoryV1 "github.com/pinai4/microservices-course-project/shared/pkg/proto/inventory/v1"
)

const grpcPort = 50051

var ErrPartAlreadyExists = errors.New("part already exists")

type PartStorage struct {
	mu    sync.RWMutex
	parts map[string]*inventoryV1.Part
}

func NewPartStorage() *PartStorage {
	return &PartStorage{
		parts: make(map[string]*inventoryV1.Part),
	}
}

// AddPart adds part
func (s *PartStorage) AddPart(part *inventoryV1.Part) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.parts[part.Uuid]; ok {
		return ErrPartAlreadyExists
	}

	s.parts[part.Uuid] = part

	return nil
}

// GetPart returns part by ID
func (s *PartStorage) GetPart(id string) *inventoryV1.Part {
	s.mu.RLock()
	defer s.mu.RUnlock()

	part, ok := s.parts[id]
	if !ok {
		return nil
	}

	return part
}

// GetPartsList returns parts list
func (s *PartStorage) GetPartsList(filter *inventoryV1.PartsFilter) []*inventoryV1.Part {
	s.mu.RLock()
	defer s.mu.RUnlock()

	hasTagListsIntersectionFunc := func(a, b []string) bool {
		set := make(map[string]struct{}, len(a))
		for _, v := range a {
			set[v] = struct{}{}
		}

		for _, v := range b {
			if _, ok := set[v]; ok {
				return true
			}
		}
		return false
	}

	parts := make([]*inventoryV1.Part, 0, len(s.parts))
	for _, p := range s.parts {
		if len(filter.GetUuids()) > 0 && !slices.Contains(filter.GetUuids(), p.GetUuid()) {
			continue
		}
		if len(filter.GetNames()) > 0 && !slices.Contains(filter.GetNames(), p.GetName()) {
			continue
		}
		if len(filter.GetCategories()) > 0 && !slices.Contains(filter.GetCategories(), p.GetCategory()) {
			continue
		}
		if len(filter.GetManufacturerCountries()) > 0 && (p.GetManufacturer() == nil || !slices.Contains(filter.GetManufacturerCountries(), p.GetManufacturer().GetCountry())) {
			continue
		}
		if len(filter.GetTags()) > 0 && (len(p.GetTags()) == 0 || !hasTagListsIntersectionFunc(filter.GetTags(), p.GetTags())) {
			continue
		}
		parts = append(parts, p)
	}

	return parts
}

type inventoryService struct {
	storage *PartStorage
	inventoryV1.UnimplementedInventoryServiceServer
}

func newInventoryService(storage *PartStorage) *inventoryService {
	return &inventoryService{
		storage: storage,
	}
}

func (s *inventoryService) GetPart(ctx context.Context, request *inventoryV1.GetPartRequest) (*inventoryV1.GetPartResponse, error) {
	part := s.storage.GetPart(request.GetUuid())
	if part == nil {
		return nil, status.Errorf(
			codes.NotFound,
			"Part with ID '%s' not found",
			request.GetUuid(),
		)
	}

	return &inventoryV1.GetPartResponse{Part: part}, nil
}

func (s *inventoryService) ListParts(ctx context.Context, request *inventoryV1.ListPartsRequest) (*inventoryV1.ListPartsResponse, error) {
	parts := s.storage.GetPartsList(request.GetFilter())

	return &inventoryV1.ListPartsResponse{Parts: parts}, nil
}

func seed(storage *PartStorage) error {
	now := timestamppb.New(time.Now())

	part1 := &inventoryV1.Part{
		Uuid:          uuid.Nil.String(),
		Name:          "Turbo Engine X200",
		Description:   "High-performance turbo engine suitable for small aircraft.",
		Price:         125000.50,
		StockQuantity: 8,
		Category:      inventoryV1.Category_CATEGORY_ENGINE,
		Dimensions: &inventoryV1.Dimensions{
			Length: 120.5,
			Width:  80.2,
			Height: 95.3,
			Weight: 450.0,
		},
		Manufacturer: &inventoryV1.Manufacturer{
			Name:    "AeroTech Industries",
			Country: "USA",
			Website: "https://aerotech.example.com",
		},
		Tags: []string{"engine", "turbo", "aircraft"},
		Metadata: map[string]*inventoryV1.Value{
			"power_kw":      {Kind: &inventoryV1.Value_DoubleValue{DoubleValue: 980.5}},
			"certified":     {Kind: &inventoryV1.Value_BoolValue{BoolValue: true}},
			"serial_number": {Kind: &inventoryV1.Value_StringValue{StringValue: "SN-ENGX200-001"}},
		},
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := storage.AddPart(part1); err != nil {
		return err
	}

	part2 := &inventoryV1.Part{
		Uuid:          uuid.NewString(),
		Name:          "Titanium Wing Panel",
		Description:   "Lightweight titanium alloy wing panel with anti-corrosion coating.",
		Price:         32000.0,
		StockQuantity: 25,
		Category:      inventoryV1.Category_CATEGORY_WING,
		Dimensions: &inventoryV1.Dimensions{
			Length: 250.0,
			Width:  60.0,
			Height: 5.0,
			Weight: 120.0,
		},
		Manufacturer: &inventoryV1.Manufacturer{
			Name:    "SkyMetal Works",
			Country: "Germany",
			Website: "https://skymetal.example.com",
		},
		Tags: []string{"wing", "titanium", "aircraft"},
		Metadata: map[string]*inventoryV1.Value{
			"material":     {Kind: &inventoryV1.Value_StringValue{StringValue: "Titanium Alloy"}},
			"batch_number": {Kind: &inventoryV1.Value_Int64Value{Int64Value: 20241001}},
			"is_tested":    {Kind: &inventoryV1.Value_BoolValue{BoolValue: true}},
		},
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := storage.AddPart(part2); err != nil {
		return err
	}

	return nil
}

func main() {
	// Init storage
	storage := NewPartStorage()
	if err := seed(storage); err != nil {
		log.Printf("failed to seed: %v\n", err)
		return
	}

	// Register our service
	service := newInventoryService(storage)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
		return
	}

	// Create GRPC server
	s := grpc.NewServer()

	inventoryV1.RegisterInventoryServiceServer(s, service)

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
