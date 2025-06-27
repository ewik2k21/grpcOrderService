package repositories

import (
	"errors"
	"fmt"
	"github.com/ewik2k21/grpcOrderService/internal/models"
	order "github.com/ewik2k21/grpcOrderService/pkg/order_service_v1"
	"github.com/google/uuid"
	"log/slog"
)

type IOrderRepository interface {
	CreateOrder(order *models.Order) (uuid.UUID, *order.Status, error)
	GetOrderStatus(userId, orderId uuid.UUID) (*order.Status, error)
}

type OrderRepository struct {
	orders map[string]models.Order
	logger *slog.Logger
}

func NewOrderRepository(logger *slog.Logger) *OrderRepository {
	return &OrderRepository{
		orders: make(map[string]models.Order),
		logger: logger,
	}
}

func (r *OrderRepository) CreateOrder(newOrder *models.Order) (*uuid.UUID, *order.Status, error) {
	orderId := uuid.New()
	newOrder.Status = order.Status_CREATED
	if _, ok := r.orders[orderId.String()]; ok {
		err := fmt.Errorf("order already created")
		r.logger.Error("order already created", slog.String("error", err.Error()))
		return nil, nil, err
	}

	r.orders[orderId.String()] = *newOrder
	r.logger.Info("order successfully created")

	return &orderId, &newOrder.Status, nil
}

func (r *OrderRepository) GetOrderStatus(userId, orderId uuid.UUID) (*order.Status, error) {
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
