package handlers

import (
	"context"
	"github.com/ewik2k21/grpcOrderService/internal/services"
	order "github.com/ewik2k21/grpcOrderService/pkg/order_service_v1"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
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
	ctx, span := otel.Tracer("OrderService").Start(ctx, "CreateOrder")
	defer span.End()

	span.SetAttributes(attribute.String("user.role", request.GetUserRole().String()))
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

func (h *OrderHandler) GetOrderStatus(
	ctx context.Context,
	req *order.GetOrderStatusRequest,
) (*order.GetOrderStatusResponse, error) {

	ctx, span := otel.Tracer("OrderService").Start(ctx, "GetOrderStatus")
	defer span.End()

	span.SetAttributes(attribute.String("user.role", req.GetUserId()))

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

func (h *OrderHandler) StreamOrderUpdates(
	req *order.StreamOrderUpdatesRequest,
	stream order.OrderService_StreamOrderUpdatesServer,
) error {
	ctx := stream.Context()
	ctx, span := otel.Tracer("OrderService").Start(ctx, "StreamOrderUpdates")
	defer span.End()

	span.SetAttributes(attribute.String("user.role", req.GetUserRole().String()))

	update, err := h.service.StreamOrderUpdates(ctx)
	if err != nil {
		return err
	}
	if err = stream.Send(update); err != nil {
		return err
	}
	return nil
}

func (h *OrderHandler) UpdateOrderStatus(ctx context.Context, req *order.UpdateOrderStatusRequest) (*order.UpdateOrderStatusResponse, error) {
	ctx, span := otel.Tracer("OrderService").Start(ctx, "UpdateOrderStatus")
	defer span.End()

	span.SetAttributes(attribute.String("user.role", req.GetStatus().String()))

	id, newStatus := req.GetOrderId(), req.GetStatus()

	status, err := h.service.UpdateOrderStatus(id, &newStatus)
	if err != nil {
		return nil, err
	}
	return &order.UpdateOrderStatusResponse{
		Status: *status,
	}, nil
}
