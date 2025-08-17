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

	orderApiV1 "github.com/xgmsx/rsf/order/internal/api/v1/order"
	orderRepo "github.com/xgmsx/rsf/order/internal/repository/order"
	orderService "github.com/xgmsx/rsf/order/internal/service/order"
	orderV1 "github.com/xgmsx/rsf/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/xgmsx/rsf/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/xgmsx/rsf/shared/pkg/proto/payment/v1"
	"github.com/xgmsx/rsf/shared/pkg/swagger"
)

const (
	httpPort          = "8080"
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
)

func main() {
	paymentConn, err := grpc.NewClient(
		"api-payment:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to connect to Payment service: %v\n", err)
	}

	inventoryConn, err := grpc.NewClient(
		"api-inventory:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to connect to Inventory service: %v\n", err)
	}

	// Инициализируем grpc-клиенты к другим сервисам
	paymentConn.Connect()
	paymentClient := paymentV1.NewPaymentServiceClient(paymentConn)
	inventoryConn.Connect()
	inventoryClient := inventoryV1.NewInventoryServiceClient(inventoryConn)

	// Инициализируем слои приложения
	repository := orderRepo.NewOrderRepository()
	service := orderService.NewOrderService(repository, inventoryClient, paymentClient)
	api := orderApiV1.NewOrderAPI(service)

	// Инициализируем HTTP сервер
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))

	// Монтируем order API
	orderApiRouter, err := orderV1.NewServer(api)
	if err != nil {
		log.Fatalf("ошибка инициализации openapi.order.v1: %v", err)
	}
	r.Mount("/api/", orderApiRouter)

	// Монтируем Swagger UI
	r.Mount("/swagger", swagger.NewSwaggerHandler(
		"/swagger/", "order_v1.swagger.json", "api"))
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger/", http.StatusMovedPermanently)
	})

	// Запускаем HTTP сервер
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
