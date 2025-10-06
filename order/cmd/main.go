package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"slices"
	"sync"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	orderV1 "github.com/pinai4/microservices-course-project/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/pinai4/microservices-course-project/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/pinai4/microservices-course-project/shared/pkg/proto/payment/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	httpPort = "8080"
	// HTTP server timeouts
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second

	inventoryServerAddress = "localhost:50051"
	paymentServerAddress   = "localhost:50052"
)

var OrderAlreadyExistsError = errors.New("order already exists")

type OrderStorage struct {
	mu     sync.RWMutex
	orders map[string]*orderV1.Order
}

func NewOrderStorage() *OrderStorage {
	return &OrderStorage{
		orders: make(map[string]*orderV1.Order),
	}
}

// CreateOrder update order by ID
func (s *OrderStorage) CreateOrder(order *orderV1.Order) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.orders[order.OrderUUID.String()]; ok {
		return OrderAlreadyExistsError
	}

	s.orders[order.OrderUUID.String()] = order

	return nil
}

// GetOrder return order by ID
func (s *OrderStorage) GetOrder(id uuid.UUID) *orderV1.Order {
	s.mu.RLock()
	defer s.mu.RUnlock()

	order, ok := s.orders[id.String()]
	if !ok {
		return nil
	}

	return order
}

// UpdateOrder update order by ID
func (s *OrderStorage) UpdateOrder(id uuid.UUID, order *orderV1.Order) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.orders[id.String()] = order
}

// OrderHandler implement interface orderV1.Handler for handling Order API http requests
type OrderHandler struct {
	storage         *OrderStorage
	inventoryClient inventoryV1.InventoryServiceClient
	paymentClient   paymentV1.PaymentServiceClient
	//orderV1.UnimplementedHandler
}

func NewOrderHandler(
	storage *OrderStorage,
	inventoryClient inventoryV1.InventoryServiceClient,
	paymentClient paymentV1.PaymentServiceClient,
) *OrderHandler {
	return &OrderHandler{
		storage:         storage,
		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,
	}
}

func (h *OrderHandler) CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.CreateOrderRes, error) {
	if len(req.PartUuids) == 0 {
		return &orderV1.BadRequestError{
			Code:    400,
			Message: "parts list is empty",
		}, nil
	}
	// check the availability of the parts in stock
	checkPartIDs := make([]string, len(req.PartUuids))
	for i, p := range req.PartUuids {
		checkPartIDs[i] = p.String()
	}

	resp, err := h.inventoryClient.ListParts(context.TODO(), &inventoryV1.ListPartsRequest{
		Filter: &inventoryV1.PartsFilter{
			Uuids: checkPartIDs,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("inventory service error: %v", err)
	}
	if len(resp.GetParts()) != len(checkPartIDs) {
		return nil, fmt.Errorf("ordered parts are absent from stock")
	}
	for _, p := range resp.GetParts() {
		if !slices.Contains(checkPartIDs, p.GetUuid()) {
			return nil, fmt.Errorf("ordered parts are absent from stock")
		}
	}
	/////////
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

	parsePaymentMethod := func(s string) paymentV1.PaymentMethod {
		if val, ok := paymentV1.PaymentMethod_value[s]; ok {
			return paymentV1.PaymentMethod(val)
		}
		return paymentV1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED
	}

	resp, err := h.paymentClient.PayOrder(context.TODO(), &paymentV1.PayOrderRequest{
		OrderUuid:     params.OrderUUID.String(),
		UserUuid:      order.UserUUID.String(),
		PaymentMethod: parsePaymentMethod(fmt.Sprintf("PAYMENT_METHOD_%s", req.GetPaymentMethod())),
	})
	if err != nil {
		return nil, fmt.Errorf("payment service error: %v", err)
	}

	tranID, err := uuid.Parse(resp.GetTransactionUuid())
	if err != nil {
		return nil, fmt.Errorf("transaction uuid parse error: %v", err)
	}

	order.TransactionUUID.SetTo(tranID)
	order.Status = orderV1.OrderStatusPAID
	order.PaymentMethod.SetTo(orderV1.OrderPaymentMethod(req.GetPaymentMethod()))
	//h.storage.UpdateOrder(order.OrderUUID, order)

	return &orderV1.ProcessOrderPaymentResponse{
		TransactionUUID: tranID,
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
	inventoryClient := inventoryV1.NewInventoryServiceClient(conn1)
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
	paymentClient := paymentV1.NewPaymentServiceClient(conn2)
	//////////////////
	//////////////////

	orderHandler := NewOrderHandler(storage, inventoryClient, paymentClient)

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
