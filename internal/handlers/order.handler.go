package handlers

import (
	"context"
	"github.com/ewik2k21/grpcOrderService/internal/services"
	order "github.com/ewik2k21/grpcOrderService/pkg/order_service_v1"
	"log/slog"
)

type OrderHandler struct {
	order.UnimplementedOrderServiceServer
	service services.OrderService
	logger  *slog.Logger
}

func NewOrderHandler(
	logger *slog.Logger,
	service *services.OrderService,
) *OrderHandler {
	return &OrderHandler{
		logger:  logger,
		service: *service,
	}
}

func (h *OrderHandler) CreateOrder(ctx context.Context, request *order.CreateOrderRequest) (*order.CreateOrderResponse, error) {

	userRole := request.GetUserRole()

	orderId, status, err := h.service.CreateOrder(ctx, userRole, request)
	if err != nil {
		return nil, err
	}

	return &order.CreateOrderResponse{
		OrderId: orderId,
		Status:  *status,
	}, nil

}

func (h *OrderHandler) GetOrderStatus(ctx context.Context, req *order.GetOrderStatusRequest) (*order.GetOrderStatusResponse, error) {

	userId := req.GetUserId()
	orderId := req.GetOrderId()

	status, err := h.service.GetOrderStatus(userId, orderId)
	if err != nil {
		return nil, err
	}

	return &order.GetOrderStatusResponse{
		Status: *status,
	}, nil
}
