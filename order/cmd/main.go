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
	"sync"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	order_v1 "github.com/ploskirev/go-rocket/shared/pkg/openapi/order/v1"
	inventory_v1 "github.com/ploskirev/go-rocket/shared/pkg/proto/inventory/v1"
	payment_v1 "github.com/ploskirev/go-rocket/shared/pkg/proto/payment/v1"
)

const (
	inventoryAddress = "localhost:50051"
	paymentAddress   = "localhost:50052"
	httpPort         = "8080"
	// Таймауты для HTTP-сервера
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
)

type OrderStatus string

const (
	PENDING_PAYMENT OrderStatus = "PENDING_PAYMENT"
	PAID            OrderStatus = "PAID"
	CANCELLED       OrderStatus = "CANCELLED"
)

type Order struct {
	OrderUUID       string
	UserUUID        string
	PartUUIDs       []string
	TotalPrice      float64
	TransactionUUID *string
	PaymentMethod   *string
	Status          OrderStatus
}

type Storage map[string]Order

type OrderHandler struct {
	ic      inventory_v1.InventoryServiceClient
	pc      payment_v1.PaymentServiceClient
	mu      sync.RWMutex
	storage *Storage
}

func NewOrderHandler(ic inventory_v1.InventoryServiceClient, pc payment_v1.PaymentServiceClient) *OrderHandler {
	return &OrderHandler{
		ic:      ic,
		pc:      pc,
		storage: &Storage{},
	}
}

func (h *OrderHandler) CreateOrder(ctx context.Context, req *order_v1.CreateOrderRequest) (order_v1.CreateOrderRes, error) {
	fmt.Println("USER UUID: ", req.UserUUID, " , PART UUIDS: ", req.PartUuids)

	res, err := h.ic.ListPart(ctx, &inventory_v1.ListPartsRequest{})
	if err != nil {
		log.Printf("Error get list parts: %s", err)
		return &order_v1.InternalServerError{
			Code:    500,
			Message: fmt.Sprintf("Error get list parts: %s", err),
		}, nil
	}

	totalPrice := float64(0)

	partsMap := map[string]*inventory_v1.Part{}
	for _, p := range res.Parts {
		partsMap[p.Uuid] = p
	}

	for _, pu := range req.PartUuids {
		if partInfo, ok := partsMap[pu]; !ok {
			log.Printf("Error parts not found")
			return &order_v1.NotFoundError{
				Code:    404,
				Message: "Error parts not found",
			}, nil
		} else {
			totalPrice += partInfo.Price
		}
	}

	orderUUID, err := uuid.NewV6()
	if err != nil {
		log.Printf("Error create uuid: %s", err)
		return &order_v1.InternalServerError{
			Code:    500,
			Message: fmt.Sprintf("Error get create uuid: %s", err),
		}, nil
	}
	orderUUIDString := orderUUID.String()

	h.mu.Lock()
	defer h.mu.Unlock()

	st := *h.storage

	st[orderUUIDString] = Order{
		OrderUUID:  orderUUIDString,
		UserUUID:   req.UserUUID,
		PartUUIDs:  req.PartUuids,
		TotalPrice: totalPrice,
		Status:     PENDING_PAYMENT,
	}

	order := &order_v1.OrderDto{
		OrderUUID:  orderUUIDString,
		TotalPrice: float32(totalPrice),
	}

	return order, nil
}

func (h *OrderHandler) PayOrder(ctx context.Context, req *order_v1.PayOrderRequest, params order_v1.PayOrderParams) (order_v1.PayOrderRes, error) {
	fmt.Println("ORDER UUID: ", params.OrderUUID, " , PAYMENT METHOD: ", req.PaymentMethod)

	h.mu.Lock()
	defer h.mu.Unlock()

	st := *h.storage

	order, ok := st[params.OrderUUID]
	if !ok {
		log.Printf("Error order not found")
		return &order_v1.NotFoundError{
			Code:    404,
			Message: "Error order not found",
		}, nil
	}

	paymentMethodMap := map[order_v1.PaymentMethod]payment_v1.PaymentMethod{
		order_v1.PaymentMethodPAYMENTMETHODUNKNOWN:       payment_v1.PaymentMethod_UNKNOWN,
		order_v1.PaymentMethodPAYMENTMETHODCARD:          payment_v1.PaymentMethod_CARD,
		order_v1.PaymentMethodPAYMENTMETHODSBP:           payment_v1.PaymentMethod_SBP,
		order_v1.PaymentMethodPAYMENTMETHODCREDITCARD:    payment_v1.PaymentMethod_CREDIT_CARD,
		order_v1.PaymentMethodPAYMENTMETHODINVESTORMONEY: payment_v1.PaymentMethod_INVESTOR_MONEY,
	}

	payRes, err := h.pc.PayOrder(ctx, &payment_v1.PayOrderRequest{
		OrderUuid:     params.OrderUUID,
		UserUuid:      order.UserUUID,
		PaymentMethod: paymentMethodMap[req.PaymentMethod],
	})
	if err != nil {
		log.Printf("Pay service error: %s", err)
		return &order_v1.InternalServerError{
			Code:    500,
			Message: fmt.Sprintf("Pay service: %s", err),
		}, nil
	}

	order.Status = PAID
	order.TransactionUUID = &payRes.TransactionUuid
	order.PaymentMethod = (*string)(&req.PaymentMethod)

	res := &order_v1.PayOrderResponse{
		TransactionUUID: payRes.TransactionUuid,
	}

	return res, nil
}

func (h *OrderHandler) GetOrder(_ context.Context, params order_v1.GetOrderParams) (order_v1.GetOrderRes, error) {
	fmt.Println("ORDER UUID: ", params.OrderUUID)

	h.mu.RLock()
	defer h.mu.RUnlock()

	st := *h.storage

	for _, o := range st {
		fmt.Println("ORDER UUID2: ", o.OrderUUID)
	}

	order, ok := st[params.OrderUUID]
	if !ok {
		log.Printf("Error order not found")
		return &order_v1.NotFoundError{
			Code:    404,
			Message: "Error order not found",
		}, nil
	}

	transactionUUID := ""
	if order.TransactionUUID != nil {
		transactionUUID = *order.TransactionUUID
	}

	paymentMethod := ""
	if order.PaymentMethod != nil {
		paymentMethod = *order.PaymentMethod
	}

	res := &order_v1.GetOrderResponse{
		TransactionUUID: transactionUUID,
		OrderUUID:       order.OrderUUID,
		UserUUID:        order.UserUUID,
		PartUuids:       order.PartUUIDs,
		TotalPrice:      float32(order.TotalPrice),
		Status:          order_v1.OrderStatus(order.Status),
		PaymentMethod:   order_v1.PaymentMethod(paymentMethod),
	}

	return res, nil
}

func (h *OrderHandler) CancelOrder(ctx context.Context, params order_v1.CancelOrderParams) (order_v1.CancelOrderRes, error) {
	fmt.Println("ORDER UUID: ", params.OrderUUID)

	h.mu.Lock()
	defer h.mu.Unlock()

	st := *h.storage

	order, ok := st[params.OrderUUID]
	if !ok {
		log.Printf("Error order not found")
		return &order_v1.NotFoundError{
			Code:    404,
			Message: "Error order not found",
		}, nil
	}

	if order.Status == PAID {
		log.Printf("Wrong status")
		return &order_v1.ConflictError{
			Code:    409,
			Message: fmt.Sprintf("Conflict! Wrong status: %s", order.Status),
		}, nil
	}

	if order.Status == PENDING_PAYMENT {
		order.Status = CANCELLED
	}

	res := &order_v1.CancelOrderNoContent{}

	return res, nil
}

func (h *OrderHandler) NewError(_ context.Context, err error) *order_v1.GenericErrorStatusCode {
	return &order_v1.GenericErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: order_v1.GenericError{
			Code:    order_v1.NewOptInt(http.StatusInternalServerError),
			Message: order_v1.NewOptString(err.Error()),
		},
	}
}

func main() {
	inventoryConn, err := grpc.NewClient(
		inventoryAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("failed to connect inventory service: %v\n", err)
		return
	}
	defer func() {
		if cerr := inventoryConn.Close(); cerr != nil {
			log.Printf("failed to close connect inventory service: %v", cerr)
		}
	}()
	paymentConn, err := grpc.NewClient(
		paymentAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("failed to connect payment service: %v\n", err)
		return
	}
	defer func() {
		if cerr := paymentConn.Close(); cerr != nil {
			log.Printf("failed to close connect payment service: %v", cerr)
		}
	}()

	ic := inventory_v1.NewInventoryServiceClient(inventoryConn)
	pc := payment_v1.NewPaymentServiceClient(paymentConn)
	orderHandler := NewOrderHandler(ic, pc)

	orderServer, err := order_v1.NewServer(orderHandler)
	if err != nil {
		log.Fatalf("ошибка создания Orders сервера: %v", err)
	}

	// Инициализируем роутер Chi
	r := chi.NewRouter()

	// Добавляем middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))
	// r.Use(customMiddleware.RequestLogger)

	// Монтируем обработчики OpenAPI
	r.Mount("/", orderServer)

	// Запускаем HTTP-сервер
	server := &http.Server{
		Addr:              net.JoinHostPort("localhost", httpPort),
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout, // Защита от Slowloris атак - тип DDoS-атаки, при которой
		// атакующий умышленно медленно отправляет HTTP-заголовки, удерживая соединения открытыми и истощая
		// пул доступных соединений на сервере. ReadHeaderTimeout принудительно закрывает соединение,
		// если клиент не успел отправить все заголовки за отведенное время.
	}

	// Запускаем сервер в отдельной горутине
	go func() {
		log.Printf("🚀 HTTP-сервер запущен на порту %s\n", httpPort)
		err = server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("❌ Ошибка запуска сервера: %v\n", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Завершение работы сервера...")

	// Создаем контекст с таймаутом для остановки сервера
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Printf("❌ Ошибка при остановке сервера: %v\n", err)
	}

	log.Println("✅ Сервер остановлен")
}
