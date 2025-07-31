package main

import (
	"context"
	"errors"
	"fmt"
	inventoryV1 "github.com/xgmsx/rsf/shared/pkg/proto/inventory/v1"
	"github.com/xgmsx/rsf/shared/pkg/swagger"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	"github.com/xgmsx/rsf/shared/pkg/interceptor"
	paymentV1 "github.com/xgmsx/rsf/shared/pkg/proto/payment/v1"
)

const (
	grpcPort = 50051
	httpPort = 8080
)

type PaymentService struct {
	paymentV1.UnimplementedPaymentServiceServer
}

func NewPaymentService() *PaymentService {
	return &PaymentService{}
}

func (s *PaymentService) PayOrder(_ context.Context, request *paymentV1.PayOrderRequest) (*paymentV1.PayOrderResponse, error) {
	if request.PaymentMethod == paymentV1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED {
		return nil, status.Error(codes.InvalidArgument, "invalid payment method provided")
	}
	transactionUUID := uuid.New()
	log.Printf("processing request=%v with transactionUUID=%v", request, transactionUUID)
	return &paymentV1.PayOrderResponse{
		TransactionUuid: transactionUUID.String(),
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
		return
	}
	defer func() {
		if cerr := lis.Close(); cerr != nil {
			log.Printf("failed to close listener: %v\n", cerr)
		}
	}()

	// Инициализируем сервисный слой
	service := NewPaymentService()

	// Инициализируем gRPC сервер
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpc.UnaryServerInterceptor(interceptor.LoggerInterceptor()),
		),
	)
	paymentV1.RegisterPaymentServiceServer(server, service)
	reflection.Register(server)

	// Запускаем gRPC сервер
	go func() {
		log.Printf("🚀 gRPC server listening on %d\n", grpcPort)
		err = server.Serve(lis)
		if err != nil {
			log.Printf("failed to serve: %v\n", err)
			return
		}
	}()

	// Инициализируем HTTP сервер с gRPC Gateway и Swagger UI
	go func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// Создаем мультиплексор для HTTP запросов в gRPC-gateway
		mux := runtime.NewServeMux()
		err = inventoryV1.RegisterInventoryServiceHandlerFromEndpoint(
			ctx,
			mux,
			fmt.Sprintf("localhost:%d", grpcPort),
			[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
		)
		if err != nil {
			log.Printf("Failed to register gateway: %v\n", err)
			return
		}

		// Создаем мультиплексор для Swagger UI
		httpMux := http.NewServeMux()
		httpMux.Handle("/api/", mux)

		httpMux.Handle("/swagger/", swagger.NewSwaggerHandler(
			"/swagger/", "payment.swagger.json", "api"))

		httpMux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/" {
				http.Redirect(w, r, "/swagger/", http.StatusMovedPermanently)
				return
			}
			http.NotFound(w, r)
		}))

		// Создаем HTTP gateway сервер
		var gwServer = &http.Server{
			Addr:              fmt.Sprintf(":%d", httpPort),
			Handler:           httpMux,
			ReadHeaderTimeout: 10 * time.Second,
		}

		// Запускаем HTTP сервер
		log.Printf("🌐 HTTP server with gRPC-Gateway and Swagger UI listening on %d\n", httpPort)
		err = gwServer.ListenAndServe()
		if err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Printf("Failed to serve HTTP: %v\n", err)
			return
		}
	}()

	// Graceful shutdown
	notify := make(chan os.Signal, 1)
	signal.Notify(notify, syscall.SIGINT, syscall.SIGTERM)
	<-notify

	log.Println("🛑 Завершение работы сервера...")
	server.GracefulStop()
	log.Println("✅ Сервер остановлен")
}
