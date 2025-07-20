package repositories

import (
	"errors"
	"fmt"
	"github.com/ewik2k21/grpcOrderService/internal/models"
	order "github.com/ewik2k21/grpcOrderService/pkg/order_service_v1"
	"github.com/google/uuid"
	"log/slog"
	"sync"
)

type IOrderRepository interface {
	CreateOrder(order *models.Order) (uuid.UUID, *order.Status, error)
	GetOrderStatus(userId, orderId uuid.UUID) (*order.Status, error)
	GetOrders() map[string]*models.Order
	UpdateOrderStatus(orderID string, status order.Status)
	HasChanged(orderID string) bool
	UpdateLastState(orderID string)
	UpdateAllOrdersToLast()
}

type OrderRepository struct {
	orders     map[string]*models.Order
	lastOrders map[string]*models.Order
	logger     *slog.Logger
	mu         sync.RWMutex
}

func NewOrderRepository(logger *slog.Logger) *OrderRepository {
	return &OrderRepository{
		orders:     make(map[string]*models.Order),
		lastOrders: make(map[string]*models.Order),
		logger:     logger,
	}
}

func (r *OrderRepository) CreateOrder(newOrder *models.Order) (*uuid.UUID, *order.Status, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	orderId := uuid.New()
	newOrder.Status = order.Status_CREATED
	if _, ok := r.orders[orderId.String()]; ok {
		err := fmt.Errorf("order already created")
		r.logger.Error("order already created", slog.String("error", err.Error()))
		return nil, nil, err
	}

	r.orders[orderId.String()] = newOrder
	r.logger.Info("order successfully created")

	return &orderId, &newOrder.Status, nil
}

func (r *OrderRepository) GetOrderStatus(userId, orderId uuid.UUID) (*order.Status, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	neededOrder, ok := r.orders[orderId.String()]
	if !ok {
		err := errors.New("failed get order by orderId")
		r.logger.Error("failed get order by orderId", slog.String("error", err.Error()))
		return nil, err
	}

	if neededOrder.UserId != userId {
		err := errors.New("wrong user id")
		r.logger.Error("wrong user id in order", slog.String("error", err.Error()))
		return nil, err
	}

	return &neededOrder.Status, nil
}

func (r *OrderRepository) GetOrders() map[string]*models.Order {
	r.mu.RLock()
	defer r.mu.RUnlock()
	orders := make(map[string]*models.Order, len(r.orders))
	for k, v := range r.orders {
		orders[k] = v
	}
	return orders
}

func (r *OrderRepository) UpdateOrderStatus(orderID string, status order.Status) {
	r.mu.Lock()
	defer r.mu.Unlock()
	needOrder := r.orders[orderID]
	needOrder.Status = status
}

func (r *OrderRepository) HasChanged(orderID string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	curr := r.orders[orderID]
	last := r.lastOrders[orderID]
	if last != nil {
		return curr != nil
	}
	return curr.Status != last.Status
}

func (r *OrderRepository) UpdateLastState(orderID string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.lastOrders[orderID] = r.orders[orderID]
}

func (r *OrderRepository) UpdateAllOrdersToLast() {
	r.mu.Lock()
	defer r.mu.Unlock()
	for id := range r.orders {
		r.lastOrders[id] = r.orders[id]
	}
}
