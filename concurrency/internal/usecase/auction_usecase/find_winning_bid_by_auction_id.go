package auction_usecase

import (
	"concurrency/internal/entity"
	"context"
	"time"
)

type FindWinningBidByAuctionIdInput struct {
	AuctionId string `json:"auction_id"`
}
type AuctionOutput struct {
	ID          string                  `json:"id"`
	ProductName string                  `json:"product_name"`
	Category    string                  `json:"category"`
	Description string                  `json:"description"`
	Condition   entity.ProductCondition `json:"condition"`
	Status      entity.AuctionStatus    `json:"status"`
	Timestamp   time.Time               `json:"timestamp" time_format:"2006-01-02 15:04:05"`
}
type BidOutput struct {
	ID        string    `json:"id"`
	UserId    string    `json:"user_id"`
	AuctionId string    `json:"auction_id"`
	Amount    float64   `json:"amount"`
	Timestamp time.Time `json:"timestamp" time_format:"2006-01-02 15:04:05"`
}
type FindWinningBidByAuctionIdOutput struct {
	Auction AuctionOutput `json:"auction"`
	Bid     *BidOutput    `json:"bid,omitempty"`
}

type IFindWinningByAuctionIdUseCase interface {
	Execute(ctx context.Context, input FindWinningBidByAuctionIdInput) (*FindWinningBidByAuctionIdOutput, error)
}

type FindWinningBidByAuctionId struct {
	AuctionRepository entity.IAuctionRepository
	BidRepository     entity.IBidRepository
}

func (b *FindWinningBidByAuctionId) Execute(ctx context.Context, input FindWinningBidByAuctionIdInput) (*FindWinningBidByAuctionIdOutput, error) {
	auction, err := b.AuctionRepository.FindById(ctx, input.AuctionId)
	if err != nil {
		return nil, err
	}
	auctionOutput := AuctionOutput{
		ID:          auction.ID,
		ProductName: auction.ProductName,
		Category:    auction.Category,
		Description: auction.Description,
		Condition:   entity.ProductCondition(auction.Condition),
		Status:      entity.AuctionStatus(auction.Status),
		Timestamp:   auction.Timestamp,
	}

	bid, err := b.BidRepository.FindWinningByAuctionId(ctx, input.AuctionId)
	if err != nil {
		return &FindWinningBidByAuctionIdOutput{
			Auction: auctionOutput,
			Bid:     nil,
		}, nil
	}
	bidOutput := &BidOutput{
		ID:        bid.ID,
		UserId:    bid.UserId,
		AuctionId: bid.AuctionId,
		Amount:    bid.Amount,
		Timestamp: bid.Timestamp,
	}
	return &FindWinningBidByAuctionIdOutput{
		Auction: auctionOutput,
		Bid:     bidOutput,
	}, nil
}
