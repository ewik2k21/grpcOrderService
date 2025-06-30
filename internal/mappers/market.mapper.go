package mappers

import (
	"github.com/ewik2k21/grpcOrderService/internal/models"
	spot "github.com/ewik2k21/grpcSpotInstrumentService/pkg/spot_instrument_v1"
	"github.com/google/uuid"
)

func MapProtoToMarkets(resp *spot.ViewMarketsResponse) ([]*models.Market, error) {
	markets := resp.GetMarkets()
	res := make([]*models.Market, 0, len(markets))
	for _, market := range markets {
		marketId, err := uuid.Parse(market.Id)
		if err != nil {
			return nil, err
		}
		res = append(res, &models.Market{
			ID:        marketId,
			Name:      market.Name,
			Enabled:   true,
			DeletedAt: nil,
		})
	}
	return res, nil
}
