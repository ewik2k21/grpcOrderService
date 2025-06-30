package models

import (
	order "github.com/ewik2k21/grpcOrderService/pkg/order_service_v1"
	"github.com/google/uuid"
)

type Order struct {
	UserId    uuid.UUID
	MarketId  uuid.UUID
	OrderType order.OrderType
	Price     float64
	Quantity  float64
	Status    order.Status
}
