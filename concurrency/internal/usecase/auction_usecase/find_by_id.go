package auction_usecase

import (
	"concurrency/internal/entity"
	"context"
	"time"
)

type FindAuctionByIdInput struct {
	ID string `json:"id"`
}
type FindAuctionByIdOutput struct {
	ID          string                  `json:"id"`
	ProductName string                  `json:"product_name"`
	Category    string                  `json:"category"`
	Description string                  `json:"description"`
	Condition   entity.ProductCondition `json:"condition"`
	Status      entity.AuctionStatus    `json:"status"`
	Timestamp   time.Time               `json:"timestamp" time_format:"2006-01-02 15:04:05"`
}

type IFindAuctionByIdUseCase interface {
	Execute(ctx context.Context, input FindAuctionByIdInput) (*FindAuctionByIdOutput, error)
}

type FindAuctionById struct {
	AuctionRepository entity.IAuctionRepository
}

func (a *FindAuctionById) Execute(ctx context.Context, input FindAuctionByIdInput) (*FindAuctionByIdOutput, error) {
	auction, err := a.AuctionRepository.FindById(ctx, input.ID)
	if err != nil {
		return nil, err
	}
	return &FindAuctionByIdOutput{
		ID:          auction.ID,
		ProductName: auction.ProductName,
		Category:    auction.Category,
		Description: auction.Description,
		Condition:   entity.ProductCondition(auction.Condition),
		Status:      entity.AuctionStatus(auction.Status),
		Timestamp:   auction.Timestamp,
	}, nil
}
