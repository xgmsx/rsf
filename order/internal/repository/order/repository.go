package order

import (
	"context"
	"sync"

	"github.com/xgmsx/rsf/order/internal/model"
	def "github.com/xgmsx/rsf/order/internal/repository"
)

var _ def.OrderRepository = (*orderRepository)(nil)

type orderRepository struct {
	mu     sync.RWMutex
	orders map[string]*model.Order
}

func NewOrderRepository() *orderRepository {
	return &orderRepository{
		orders: make(map[string]*model.Order),
	}
}

func (r *orderRepository) Get(ctx context.Context, orderUUID string) (model.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	order, ok := r.orders[orderUUID]
	if !ok {
		return model.Order{}, model.ErrOrderNotFound
	}

	return *order, nil
}

func (r *orderRepository) Update(ctx context.Context, order model.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.orders[order.OrderUUID.String()] = &order
	return nil
}
