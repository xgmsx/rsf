package order

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"

	"github.com/xgmsx/rsf/order/internal/client"
	"github.com/xgmsx/rsf/order/internal/model"
	"github.com/xgmsx/rsf/order/internal/repository"
	def "github.com/xgmsx/rsf/order/internal/service"
)

var _ def.OrderService = (*orderService)(nil)

type orderService struct {
	repo            repository.OrderRepository
	inventoryClient client.InventoryClient
	paymentClient   client.PaymentClient
}

func NewOrderService(
	repo repository.OrderRepository,
	inventoryClient client.InventoryClient,
	paymentClient client.PaymentClient,
) *orderService {
	return &orderService{repo: repo, inventoryClient: inventoryClient, paymentClient: paymentClient}
}

func (s *orderService) CreateOrder(ctx context.Context, input model.CreateOrderInput) (model.CreateOrderOutput, error) {
	parts, err := s.inventoryClient.GetParts(ctx, input.PartUUIDs)
	if err != nil {
		log.Printf("error while fetching inventory: %v\n", err)
		return model.CreateOrderOutput{}, model.ErrFailedToFetchInventory
	}

	if len(parts) != len(input.PartUUIDs) {
		returned := make(map[string]struct{}, len(parts))
		for _, part := range parts {
			returned[part.GetUuid()] = struct{}{}
		}

		var missing []string
		for _, reqUUID := range input.PartUUIDs {
			if _, ok := returned[reqUUID.String()]; !ok {
				missing = append(missing, reqUUID.String())
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
		PartUUIDs:  input.PartUUIDs,
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

	txUUID, err := s.paymentClient.PayOrder(ctx, order.UserUUID, order.OrderUUID, input.PaymentMethod)
	if err != nil {
		log.Println("failed to process payment:", err)
		return model.PayOrderOutput{}, err
	}

	order.Status = model.OrderStatusPAID
	order.PaymentMethod = &input.PaymentMethod
	order.TransactionUUID = txUUID

	err = s.repo.Update(ctx, order)
	if err != nil {
		return model.PayOrderOutput{}, err
	}

	return model.PayOrderOutput{
		TransactionUUID: *txUUID,
	}, nil
}
