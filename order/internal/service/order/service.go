package order

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"

	"github.com/xgmsx/rsf/order/internal/model"
	"github.com/xgmsx/rsf/order/internal/repository"
	"github.com/xgmsx/rsf/order/internal/service"
	inventoryV1 "github.com/xgmsx/rsf/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/xgmsx/rsf/shared/pkg/proto/payment/v1"
)

var _ service.OrderService = (*orderService)(nil)

type orderService struct {
	repo            repository.OrderRepository
	inventoryClient inventoryV1.InventoryServiceClient
	paymentClient   paymentV1.PaymentServiceClient
}

func NewOrderService(
	repo repository.OrderRepository,
	inventoryClient inventoryV1.InventoryServiceClient,
	paymentClient paymentV1.PaymentServiceClient,
) *orderService {
	return &orderService{repo: repo, inventoryClient: inventoryClient, paymentClient: paymentClient}
}

func (s *orderService) CreateOrder(ctx context.Context, input model.CreateOrderInput) (model.CreateOrderOutput, error) {
	partUUIDs := make([]string, len(input.PartUuids))
	for i, partUUID := range input.PartUuids {
		partUUIDs[i] = partUUID
	}

	listPartsResponse, err := s.inventoryClient.ListParts(ctx, &inventoryV1.ListPartsRequest{
		Filter: &inventoryV1.PartsFilter{
			Uuids: partUUIDs,
		},
	})
	if err != nil {
		log.Printf("error while fetching inventory: %v\n", err)
		return model.CreateOrderOutput{}, model.ErrFailedToFetchInventory
	}

	parts := listPartsResponse.GetParts()
	if len(parts) != len(input.PartUuids) {

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
		err = fmt.Errorf("the following partUuid(s) do not exist: %v: %w", missing, model.ErrPartDoesNotExist)
		return model.CreateOrderOutput{}, err
	}

	var totalPrice float64
	for _, part := range parts {
		totalPrice += part.GetPrice()
	}

	order := model.Order{
		UserUUID:   input.UserUUID,
		OrderUUID:  uuid.New(),
		PartUUIDs:  input.PartUuids,
		Status:     model.OrderStatusPENDINGPAYMENT,
		TotalPrice: totalPrice,
	}
	err = s.repo.Update(ctx, order)
	if err != nil {
		return model.CreateOrderOutput{}, err
	}

	output := model.CreateOrderOutput{
		OrderUUID:  order.OrderUUID,
		TotalPrice: order.TotalPrice,
	}
	return output, nil
}

func (s *orderService) GetOrder(ctx context.Context, orderUUID string) (model.Order, error) {
	return s.repo.Get(ctx, orderUUID)
}

func (s *orderService) CancelOrder(ctx context.Context, orderUUID string) (model.Order, error) {
	order, err := s.repo.Get(ctx, orderUUID)
	if err != nil {
		return model.Order{}, err
	}

	if order.Status == model.OrderStatusPAID {
		return model.Order{}, model.ErrOrderAlreadyPaid
	}

	order.Status = model.OrderStatusCANCELLED

	err = s.repo.Update(ctx, order)
	if err != nil {
		return model.Order{}, err
	}

	return order, nil
}

func (s *orderService) PayOrder(ctx context.Context, input model.PayOrderInput) (model.PayOrderOutput, error) {
	order, err := s.repo.Get(ctx, input.OrderUUID.String())
	if err != nil {
		return model.PayOrderOutput{}, err
	}

	paymentMethodsMap := map[model.PaymentMethod]paymentV1.PaymentMethod{
		model.PaymentMethodCARD:          paymentV1.PaymentMethod_PAYMENT_METHOD_CARD,
		model.PaymentMethodSBP:           paymentV1.PaymentMethod_PAYMENT_METHOD_SBP,
		model.PaymentMethodCREDITCARD:    paymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD,
		model.PaymentMethodINVESTORMONEY: paymentV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY,
	}
	paymentMethod, ok := paymentMethodsMap[input.PaymentMethod]
	if !ok {
		return model.PayOrderOutput{}, model.ErrPaymentMethodIsNotSupported
	}

	response, err := s.paymentClient.PayOrder(ctx, &paymentV1.PayOrderRequest{
		UserUuid:      order.UserUUID.String(),
		OrderUuid:     order.OrderUUID.String(),
		PaymentMethod: paymentMethod,
	})
	if err != nil {
		log.Println("failed to process payment:", err)
		return model.PayOrderOutput{}, err
	}

	transactionUUID, err := uuid.Parse(response.TransactionUuid)
	if err != nil {
		log.Println("failed to parse transaction UUID:", err)
		return model.PayOrderOutput{}, err
	}

	order.Status = model.OrderStatusPAID
	order.PaymentMethod = &input.PaymentMethod
	order.TransactionUUID = &transactionUUID

	err = s.repo.Update(ctx, order)
	if err != nil {
		return model.PayOrderOutput{}, err
	}

	return model.PayOrderOutput{
		TransactionUUID: transactionUUID.String(),
	}, nil
}
