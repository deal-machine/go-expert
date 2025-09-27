package auction_usecase

import (
	"concurrency/internal/entity"
	"context"
	"time"
)

type FindAuctionsOutput struct {
	ID          string                  `json:"id"`
	ProductName string                  `json:"product_name"`
	Category    string                  `json:"category"`
	Description string                  `json:"description"`
	Condition   entity.ProductCondition `json:"condition"`
	Status      entity.AuctionStatus    `json:"status"`
	Timestamp   time.Time               `json:"timestamp" time_format:"2006-01-02 15:04:05"`
}

type FindAuctionsInput struct {
	Status      entity.AuctionStatus `json:"status"`
	Category    string               `json:"category"`
	ProductName string               `json:"product_name"`
}

type IFindAuctionsUseCase interface {
	Execute(ctx context.Context, input FindAuctionsInput) ([]FindAuctionsOutput, error)
}

type FindAuctions struct {
	AuctionRepository entity.IAuctionRepository
}

func (a *FindAuctions) Execute(ctx context.Context, input FindAuctionsInput) ([]FindAuctionsOutput, error) {
	auctions, err := a.AuctionRepository.FindBy(ctx, entity.AuctionStatus(input.Status), input.Category, input.ProductName)
	if err != nil {
		return nil, err
	}
	var auctionsOutput []FindAuctionsOutput
	for _, auction := range auctions {
		auctionsOutput = append(auctionsOutput, FindAuctionsOutput{
			ID:          auction.ID,
			ProductName: auction.ProductName,
			Category:    auction.Category,
			Description: auction.Description,
			Condition:   entity.ProductCondition(auction.Condition),
			Status:      entity.AuctionStatus(auction.Status),
			Timestamp:   auction.Timestamp,
		})
	}
	return auctionsOutput, nil
}
