package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/xgmsx/rsf/shared/pkg/swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	orderV1 "github.com/xgmsx/rsf/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/xgmsx/rsf/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/xgmsx/rsf/shared/pkg/proto/payment/v1"
)

const (
	httpPort          = "8080"
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
)

const (
	ORDER_STATUS_UNKNOWN = iota
	ORDER_STATUS_PENDING_PAYMENT
	ORDER_STATUS_PAID
	ORDER_STATUS_CANCELED
)

type Order struct {
	OrderUuid       uuid.UUID
	UserUuid        uuid.UUID
	PartsUuids      []uuid.UUID
	TotalPrice      float64
	TransactionUuid *uuid.UUID
	PaymentMethod   *string
	Status          uint8
}

type OrderStorage struct {
	mu     sync.RWMutex
	orders map[uuid.UUID]*orderV1.Order
}

func NewOrderStorage() *OrderStorage {
	return &OrderStorage{
		orders: make(map[uuid.UUID]*orderV1.Order),
	}
}

func (s *OrderStorage) UpdateOrder(order *orderV1.Order) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.orders[order.OrderUUID] = order

	return nil
}

type OrderHandler struct {
	inventoryClient inventoryV1.InventoryServiceClient
	paymentClient   paymentV1.PaymentServiceClient
	storage         *OrderStorage
}

func NewOrderHandler(
	inventoryClient inventoryV1.InventoryServiceClient,
	paymentClient paymentV1.PaymentServiceClient,
	storage *OrderStorage,
) *OrderHandler {
	return &OrderHandler{
		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,
		storage:         storage,
	}
}

func (h *OrderHandler) CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.CreateOrderRes, error) {
	orderUuid := uuid.New()

	partUUIDs := make([]string, len(req.PartUuids))
	for i, partUUID := range req.PartUuids {
		partUUIDs[i] = partUUID
	}

	listPartsResponse, err := h.inventoryClient.ListParts(ctx, &inventoryV1.ListPartsRequest{
		Filter: &inventoryV1.PartsFilter{
			Uuids: partUUIDs,
		},
	})
	if err != nil {
		return &orderV1.InternalServerError{
			Code:    500,
			Message: "failed to fetch inventory service",
		}, err
	}

	parts := listPartsResponse.GetParts()

	if len(parts) != len(req.PartUuids) {

		returned := make(map[string]struct{}, len(parts))
		for _, part := range parts {
			returned[part.GetUuid()] = struct{}{}
		}

		var missing []string
		for _, reqUUID := range partUUIDs {
			if _, ok := returned[reqUUID]; !ok {
				missing = append(missing, reqUUID)
			}
		}

		msg := fmt.Sprintf("The following partUuid(s) do not exist: %v", missing)
		return &orderV1.BadRequestError{
			Code:    400,
			Message: msg,
		}, nil
	}

	var totalPrice float64
	for _, part := range parts {
		totalPrice += part.GetPrice()
	}

	order := &orderV1.Order{
		UserUUID:   req.UserUUID,
		OrderUUID:  orderUuid,
		PartUuids:  req.PartUuids,
		Status:     orderV1.OrderStatusPENDINGPAYMENT,
		TotalPrice: totalPrice,
	}
	err = h.storage.UpdateOrder(order)
	if err != nil {
		return nil, h.NewError(ctx, err)
	}
	response := &orderV1.CreateOrderResponse{
		OrderUUID:  order.OrderUUID,
		TotalPrice: float32(order.TotalPrice),
	}
	return response, nil
}

func (h *OrderHandler) PayOrder(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.PayOrderParams) (orderV1.PayOrderRes, error) {
	h.storage.mu.Lock()
	defer h.storage.mu.Unlock()

	order, ok := h.storage.orders[params.OrderUUID]

	if !ok {
		return &orderV1.NotFoundError{
			Code:    404,
			Message: "Order with UUID: '" + params.OrderUUID.String() + "' not found",
		}, nil
	}

	paymentMethodsMap := map[orderV1.PayOrderRequestPaymentMethod]paymentV1.PaymentMethod{
		orderV1.PayOrderRequestPaymentMethodCARD:          paymentV1.PaymentMethod_PAYMENT_METHOD_CARD,
		orderV1.PayOrderRequestPaymentMethodSBP:           paymentV1.PaymentMethod_PAYMENT_METHOD_SBP,
		orderV1.PayOrderRequestPaymentMethodCREDITCARD:    paymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD,
		orderV1.PayOrderRequestPaymentMethodINVESTORMONEY: paymentV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY,
	}

	paymentMethod, ok := paymentMethodsMap[req.PaymentMethod]
	if !ok {
		return &orderV1.BadRequestError{
			Code:    400,
			Message: "Payment method '" + string(req.PaymentMethod) + "' is not supported",
		}, nil
	}

	response, err := h.paymentClient.PayOrder(ctx, &paymentV1.PayOrderRequest{
		UserUuid:      order.UserUUID.String(),
		OrderUuid:     order.OrderUUID.String(),
		PaymentMethod: paymentMethod,
	})
	if err != nil {
		return &orderV1.InternalServerError{
			Code:    500,
			Message: "failed to process payment",
		}, err
	}

	transactionUUID, err := uuid.Parse(response.TransactionUuid)
	if err != nil {
		return &orderV1.InternalServerError{
			Code:    500,
			Message: "failed to process payment",
		}, err
	}

	order.Status = orderV1.OrderStatusPAID
	order.PaymentMethod = orderV1.NewOptNilOrderPaymentMethod(
		orderV1.OrderPaymentMethod(req.PaymentMethod),
	)
	order.TransactionUUID = orderV1.NewOptNilUUID(transactionUUID)

	return &orderV1.PayOrderResponse{
		TransactionUUID: transactionUUID,
	}, nil
}

func (h *OrderHandler) GetOrder(ctx context.Context, params orderV1.GetOrderParams) (orderV1.GetOrderRes, error) {
	h.storage.mu.RLock()
	defer h.storage.mu.RUnlock()

	order, ok := h.storage.orders[params.OrderUUID]
	if !ok {
		return &orderV1.NotFoundError{
			Code:    404,
			Message: "Order with UUID: '" + params.OrderUUID.String() + "' not found",
		}, nil
	}

	resp := &orderV1.Order{
		OrderUUID:       order.OrderUUID,
		UserUUID:        order.UserUUID,
		PartUuids:       order.PartUuids,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: order.TransactionUUID,
		PaymentMethod:   order.PaymentMethod,
		Status:          order.Status,
	}

	return resp, nil
}

func (h *OrderHandler) CancelOrder(ctx context.Context, params orderV1.CancelOrderParams) (orderV1.CancelOrderRes, error) {
	h.storage.mu.Lock()
	defer h.storage.mu.Unlock()

	order, ok := h.storage.orders[params.OrderUUID]
	if !ok {
		return &orderV1.NotFoundError{
			Code:    404,
			Message: "Order with UUID: '" + params.OrderUUID.String() + "' not found",
		}, nil
	}

	switch order.Status {
	case orderV1.OrderStatusPENDINGPAYMENT:
		// Меняем статус на CANCELLED
		order.Status = orderV1.OrderStatusCANCELLED

		// Возвращаем 204 No Content (nil, nil)
		return &orderV1.CancelOrderNoContent{}, nil

	case orderV1.OrderStatusPAID:
		// Заказ уже оплачен, отменить нельзя
		return &orderV1.ConflictError{
			Code:    409,
			Message: "Order already paid and cannot be cancelled",
		}, nil

	default:
		// Для других статусов можно вернуть 404 или 409, в зависимости от логики
		return &orderV1.BadRequestError{
			Code:    400,
			Message: "Order cannot be cancelled in its current status",
		}, nil
	}
}

// NewError создает новую ошибку в формате GenericError
func (h *OrderHandler) NewError(_ context.Context, err error) *orderV1.GenericErrorStatusCode {
	return &orderV1.GenericErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: orderV1.GenericError{
			Code:    orderV1.NewOptInt(http.StatusInternalServerError),
			Message: orderV1.NewOptString(err.Error()),
		},
	}
}

func main() {
	// Создаем хранилище для данных о погоде
	storage := NewOrderStorage()

	paymentConn, err := grpc.NewClient(
		"localhost:50052",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to connect to Payment Service: %v\n", err)
	}

	inventoryConn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to connect to Payment Service: %v\n", err)
	}

	// Инициализируем grpc-клиенты к другим сервисам
	paymentConn.Connect()
	paymentClient := paymentV1.NewPaymentServiceClient(paymentConn)
	inventoryConn.Connect()
	inventoryClient := inventoryV1.NewInventoryServiceClient(inventoryConn)

	// Инициализируем обработчик сервиса
	orderHandler := NewOrderHandler(inventoryClient, paymentClient, storage)

	// Инициализируем роутер Chi
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))

	// Монтируем обработчик OpenAPI (/api/v1/orders, /api/v1/orders/{order_uuid} и т.д.)
	orderApiV1, err := orderV1.NewServer(orderHandler)
	if err != nil {
		log.Fatalf("ошибка создания сервера OpenAPI: %v", err)
	}
	r.Mount("/api/", orderApiV1)

	// Монтируем обработчик Swagger UI
	r.Mount("/swagger", swagger.NewSwaggerHandler(
		"/swagger/", "order_v1.swagger.json", "api"))
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger/", http.StatusMovedPermanently)
	})

	// Запускаем HTTP-сервер
	server := &http.Server{
		Addr:              net.JoinHostPort("0.0.0.0", httpPort),
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	go func() {
		log.Printf("🚀 HTTP-сервер запущен на порту %s\n", httpPort)
		err = server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("❌ Ошибка запуска сервера: %v\n", err)
		}
	}()

	// Создаем graceful shutdown контекст для ожидания сигнала завершения
	notify := make(chan os.Signal, 1)
	signal.Notify(notify, syscall.SIGINT, syscall.SIGTERM)
	<-notify

	tCtx, tCancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer tCancel()

	log.Println("🛑 Завершение работы сервера...")
	err = server.Shutdown(tCtx)
	if err != nil {
		log.Printf("❌ Ошибка при остановке сервера: %v\n", err)
	}
	log.Println("✅ Сервер остановлен")
}
