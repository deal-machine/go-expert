package auction_usecase

import (
	"concurrency/internal/entity"
	"context"
	"log"
	"time"
)

type CreateAuctionInput struct {
	ProductName string                  `json:"product_name"`
	Category    string                  `json:"category"`
	Description string                  `json:"description"`
	Condition   entity.ProductCondition `json:"condition"`
}
type CreateAuctionOutput struct {
	ID          string                  `json:"id"`
	ProductName string                  `json:"product_name"`
	Category    string                  `json:"category"`
	Description string                  `json:"description"`
	Condition   entity.ProductCondition `json:"condition"`
	Status      entity.AuctionStatus    `json:"status"`
	Timestamp   time.Time               `json:"timestamp" time_format:"2006-01-02 15:04:05"`
}

type ICreateAuctionUseCase interface {
	Execute(ctx context.Context, input CreateAuctionInput) (*CreateAuctionOutput, error)
}

type CreateAuction struct {
	AuctionRepository entity.IAuctionRepository
}

func (a *CreateAuction) Execute(ctx context.Context, input CreateAuctionInput) (*CreateAuctionOutput, error) {
	auction, err := entity.NewAuction(
		input.ProductName,
		input.Category,
		input.Description,
		entity.ProductCondition(input.Condition),
	)
	if err != nil {
		log.Println("Error on create auction entity", err)
		return nil, err
	}
	if err := a.AuctionRepository.Create(ctx, *auction); err != nil {
		log.Println("Error on create auction on database", err)
		return nil, err
	}
	return &CreateAuctionOutput{
		ID:          auction.ID,
		ProductName: auction.ProductName,
		Category:    auction.Category,
		Description: auction.Description,
		Condition:   entity.ProductCondition(auction.Condition),
		Status:      entity.AuctionStatus(auction.Status),
		Timestamp:   auction.Timestamp,
	}, nil
}
