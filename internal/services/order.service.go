package services

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ewik2k21/grpcOrderService/internal/mappers"
	"github.com/ewik2k21/grpcOrderService/internal/repositories"
	order "github.com/ewik2k21/grpcOrderService/pkg/order_service_v1"
	pkg "github.com/ewik2k21/grpcSpotInstrumentService/pkg/spot_instrument_v1"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"log/slog"
	"time"
)

type OrderService struct {
	repo        *repositories.OrderRepository
	client      pkg.SpotInstrumentServiceClient
	logger      *slog.Logger
	redisClient *redis.Client
	cacheTTL    time.Duration
}

//redis todo

func NewOrderService(
	repo *repositories.OrderRepository,
	client pkg.SpotInstrumentServiceClient,
	logger *slog.Logger,
	redisClient *redis.Client,
	cacheTTL time.Duration,
) *OrderService {
	return &OrderService{
		repo:        repo,
		client:      client,
		logger:      logger,
		redisClient: redisClient,
		cacheTTL:    cacheTTL,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, userRole pkg.UserRole, request *order.CreateOrderRequest) (string, *order.Status, error) {
	cacheKey := fmt.Sprintf("markets:%v", userRole.String())

	cachedData, err := s.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var cachedResp pkg.ViewMarketsResponse
		if err := json.Unmarshal([]byte(cachedData), &cachedResp); err == nil {
			return CheckMarkets(&cachedResp, request, s)
		}
	}

	resp, err := s.client.ViewMarkets(
		ctx,
		&pkg.ViewMarketsRequest{
			UserRole: userRole,
		})
	if err != nil {
		s.logger.Error("error request view markets from clients", slog.String("error", err.Error()))
		return "", nil, err
	}

	dataBytes, err := json.Marshal(resp)
	if err == nil {
		err = s.redisClient.SetEx(ctx, cacheKey, dataBytes, s.cacheTTL).Err()
		if err != nil {
			s.logger.Error("failed to cache data: %v", err)
		} else {
			s.logger.Info("data cached")
		}
	}

	return CheckMarkets(resp, request, s)

}

func CheckMarkets(resp *pkg.ViewMarketsResponse, request *order.CreateOrderRequest, s *OrderService) (string, *order.Status, error) {
	markets, err := mappers.MapProtoToMarkets(resp)
	if err != nil {
		s.logger.Error("failed mapping proto to markets", slog.String("error", err.Error()))
		return "", nil, err
	}

	marketId := request.GetMarketId()
	var ok = false

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

func (s *OrderService) StreamOrderUpdates(ctx context.Context) (*order.OrderStatusUpdateResponse, error) {
	s.repo.UpdateAllOrdersToLast()

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(5 * time.Second):
			orders := s.repo.GetOrders()
			for id, ord := range orders {
				if s.repo.HasChanged(id) {
					update := &order.OrderStatusUpdateResponse{
						Status: ord.Status,
					}
					s.repo.UpdateLastState(id)
					return update, nil
				}
			}
		}
	}
}

func (s *OrderService) UpdateOrderStatus(id string, status *order.Status) (*order.Status, error) {
	s.repo.UpdateOrderStatus(id, *status)
	return status, nil
}
