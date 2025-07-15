package services

import (
	"context"
	"fmt"
	"github.com/ewik2k21/grpcOrderService/internal/mappers"
	"github.com/ewik2k21/grpcOrderService/internal/repositories"
	order "github.com/ewik2k21/grpcOrderService/pkg/order_service_v1"
	pkg "github.com/ewik2k21/grpcSpotInstrumentService/pkg/spot_instrument_v1"
	"github.com/google/uuid"
	"log/slog"
)

type OrderService struct {
	repo   repositories.OrderRepository
	client pkg.SpotInstrumentServiceClient
	logger *slog.Logger
}

//redis todo

func NewOrderService(
	repo repositories.OrderRepository,
	client pkg.SpotInstrumentServiceClient,
	logger *slog.Logger,
) *OrderService {
	return &OrderService{
		repo:   repo,
		client: client,
		logger: logger,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, userRole pkg.UserRole, request *order.CreateOrderRequest) (string, *order.Status, error) {

	resp, err := s.client.ViewMarkets(
		ctx,
		&pkg.ViewMarketsRequest{
			UserRole: userRole,
		})
	if err != nil {
		s.logger.Error("error request view markets from clients", slog.String("error", err.Error()))
		return "", nil, err
	}

	markets, err := mappers.MapProtoToMarkets(resp)
	if err != nil {
		s.logger.Error("failed mapping proto to markets", slog.String("error", err.Error()))
		return "", nil, err
	}

	marketId := request.GetMarketId()
	var ok bool = false

	mapOrder, err := mappers.MapProtoToOrder(request)
	if err != nil {
		s.logger.Error("failed mapping proto to order", slog.String("error", err.Error()))
		return "", nil, err
	}

	var orderId *uuid.UUID
	var status *order.Status
	for _, market := range markets {
		if market.ID.String() == marketId {
			ok = true
			orderId, status, err = s.repo.CreateOrder(mapOrder)
		}
	}

	if !ok {
		return "", nil, fmt.Errorf("needed market not found")
	}

	return orderId.String(), status, nil

}

func (s *OrderService) GetOrderStatus(userIdString, orderIdString string) (*order.Status, error) {
	userId, err := uuid.Parse(userIdString)
	if err != nil {
		s.logger.Error("failed parse userId", slog.String("error", err.Error()))
		return nil, err
	}
	orderId, err := uuid.Parse(orderIdString)
	if err != nil {
		s.logger.Error("failed parse orderId", slog.String("error", err.Error()))
		return nil, err
	}

	status, err := s.repo.GetOrderStatus(userId, orderId)
	if err != nil {
		s.logger.Error("error get order status from repo", slog.String("error", err.Error()))
		return nil, err
	}

	return status, nil

}
