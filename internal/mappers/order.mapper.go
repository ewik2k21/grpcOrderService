package mappers

import (
	"github.com/ewik2k21/grpcOrderService/internal/models"
	order "github.com/ewik2k21/grpcOrderService/pkg/order_service_v1"
	"github.com/google/uuid"
)

func MapProtoToOrder(request *order.CreateOrderRequest) (*models.Order, error) {
	userId, err := uuid.Parse(request.GetUserId())
	if err != nil {
		return nil, err
	}

	marketId, err := uuid.Parse(request.GetMarketId())
	if err != nil {
		return nil, err
	}

	return &models.Order{
		UserId:    userId,
		MarketId:  marketId,
		OrderType: request.GetOrderType(),
		Price:     request.GetPrice(),
		Quantity:  request.GetPrice(),
	}, nil

}
