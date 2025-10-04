package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	orderV1 "github.com/pinai4/microservices-course-project/shared/pkg/openapi/order/v1"
)

const (
	httpPort = "8080"
	// HTTP server timeouts
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
)

var OrderAlreadyExistsError = errors.New("order already exists")

type OrderStorage struct {
	mu    sync.RWMutex
	order map[string]*orderV1.Order
}

func NewOrderStorage() *OrderStorage {
	return &OrderStorage{
		order: make(map[string]*orderV1.Order),
	}
}

// CreateOrder update order by ID
func (s *OrderStorage) CreateOrder(order *orderV1.Order) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.order[order.OrderUUID.String()]; ok {
		return OrderAlreadyExistsError
	}

	s.order[order.OrderUUID.String()] = order

	return nil
}

// GetOrder return order by ID
func (s *OrderStorage) GetOrder(id uuid.UUID) *orderV1.Order {
	s.mu.RLock()
	defer s.mu.RUnlock()

	order, ok := s.order[id.String()]
	if !ok {
		return nil
	}

	return order
}

// UpdateOrder update order by ID
func (s *OrderStorage) UpdateOrder(id uuid.UUID, order *orderV1.Order) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.order[id.String()] = order
}

// OrderHandler implement interface orderV1.Handler for handling Order API http requests
type OrderHandler struct {
	storage *OrderStorage
	//orderV1.UnimplementedHandler
}

func NewOrderHandler(storage *OrderStorage) *OrderHandler {
	return &OrderHandler{
		storage: storage,
	}
}

func (h *OrderHandler) CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.CreateOrderRes, error) {
	orderID := uuid.New()
	var totalPrice float32 = 9.99

	order := &orderV1.Order{
		OrderUUID:       orderID,
		UserUUID:        req.UserUUID,
		PartUuids:       req.PartUuids,
		TotalPrice:      totalPrice,
		TransactionUUID: orderV1.OptUUID{},
		PaymentMethod:   orderV1.OptOrderPaymentMethod{},
		Status:          orderV1.OrderStatusPENDINGPAYMENT,
	}

	if err := h.storage.CreateOrder(order); err != nil {
		return &orderV1.InternalServerError{
			Code:    500,
			Message: "Internal Server Error",
		}, nil
	}

	return &orderV1.CreateOrderResponse{
		OrderUUID:  orderID,
		TotalPrice: totalPrice,
	}, nil
}

func (h *OrderHandler) GetOrder(ctx context.Context, params orderV1.GetOrderParams) (orderV1.GetOrderRes, error) {
	order := h.storage.GetOrder(params.OrderUUID)
	if order == nil {
		return &orderV1.NotFoundError{
			Code:    404,
			Message: "Order with ID '" + params.OrderUUID.String() + "' not found",
		}, nil
	}

	return order, nil
}

func (h *OrderHandler) ProcessOrderPayment(
	ctx context.Context,
	req *orderV1.ProcessOrderPaymentRequest,
	params orderV1.ProcessOrderPaymentParams,
) (orderV1.ProcessOrderPaymentRes, error) {
	order := h.storage.GetOrder(params.OrderUUID)
	if order == nil {
		return &orderV1.NotFoundError{
			Code:    404,
			Message: "Order with ID '" + params.OrderUUID.String() + "' not found",
		}, nil
	}

	transactionID := uuid.New()
	order.TransactionUUID.SetTo(transactionID)
	order.Status = orderV1.OrderStatusPAID
	order.PaymentMethod.SetTo(orderV1.OrderPaymentMethod(req.GetPaymentMethod()))
	//h.storage.UpdateOrder(order.OrderUUID, order)

	return &orderV1.ProcessOrderPaymentResponse{
		TransactionUUID: transactionID,
	}, nil
}

func (h *OrderHandler) CancelOrder(ctx context.Context, params orderV1.CancelOrderParams) (orderV1.CancelOrderRes, error) {
	order := h.storage.GetOrder(params.OrderUUID)
	if order == nil {
		return &orderV1.NotFoundError{
			Code:    404,
			Message: "Order with ID '" + params.OrderUUID.String() + "' not found",
		}, nil
	}

	if order.Status == orderV1.OrderStatusPAID {
		return &orderV1.ConflictError{
			Code:    409,
			Message: "Order already has been paid and can't be canceled",
		}, nil
	}

	order.Status = orderV1.OrderStatusCANCELLED
	//h.storage.UpdateOrder(order.OrderUUID, order)

	return &orderV1.CancelOrderNoContent{}, nil
}

func (h *OrderHandler) NewError(ctx context.Context, err error) *orderV1.GenericErrorStatusCode {
	return &orderV1.GenericErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: orderV1.GenericError{
			Code:    orderV1.NewOptInt(http.StatusInternalServerError),
			Message: orderV1.NewOptString(err.Error()),
		},
	}
}

func main() {
	storage := NewOrderStorage()

	orderHandler := NewOrderHandler(storage)

	// Create OpenAPI server
	orderServer, err := orderV1.NewServer(orderHandler)
	if err != nil {
		log.Fatalf("OpenAPI server creation error: %v", err)
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
