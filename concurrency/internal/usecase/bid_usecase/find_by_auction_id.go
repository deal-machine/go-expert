package bid_usecase

import (
	"concurrency/internal/entity"
	"context"
	"time"
)

type FindByAuctionIdInput struct {
	AuctionId string `json:"auction_id"`
}
type FindByAuctionIdOutput struct {
	ID        string    `json:"id"`
	UserId    string    `json:"user_id"`
	AuctionId string    `json:"auction_id"`
	Amount    float64   `json:"amount"`
	Timestamp time.Time `json:"timestamp" time_format:"2006-01-02 15:04:05"`
}

type IFindByAuctionIdUseCase interface {
	Execute(ctx context.Context, input FindByAuctionIdInput) ([]FindByAuctionIdOutput, error)
}

type FindByAuctionId struct {
	BidRepository entity.IBidRepository
}

func (b *FindByAuctionId) Execute(ctx context.Context, input FindByAuctionIdInput) ([]FindByAuctionIdOutput, error) {
	bids, err := b.BidRepository.FindByAuctionId(ctx, input.AuctionId)
	if err != nil {
		return nil, err
	}
	var bidOutput []FindByAuctionIdOutput
	for _, bid := range bids {
		bidOutput = append(bidOutput, FindByAuctionIdOutput{
			ID:        bid.ID,
			UserId:    bid.UserId,
			AuctionId: bid.AuctionId,
			Amount:    bid.Amount,
			Timestamp: bid.Timestamp,
		})
	}
	return bidOutput, nil
}
