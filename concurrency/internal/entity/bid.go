package entity

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Bid struct {
	ID        string
	UserId    string
	AuctionId string
	Amount    float64
	Timestamp time.Time
}

func NewBid(userId, auctionId string, amount float64) (*Bid, error) {
	bid := &Bid{
		ID:        uuid.NewString(),
		UserId:    userId,
		AuctionId: auctionId,
		Amount:    amount,
		Timestamp: time.Now(),
	}
	if err := bid.Validate(); err != nil {
		return nil, err
	}
	return bid, nil
}

func (b *Bid) Validate() error {
	var error error = nil
	if b.Amount <= 0 {
		error = errors.New("amount must be positive, greater than zero")
	}
	if err := uuid.Validate(b.ID); err != nil {
		error = errors.New("id must be a valid uuid")
	}
	if err := uuid.Validate(b.AuctionId); err != nil {
		error = errors.New("auctionId must be a valid uuid")
	}
	if err := uuid.Validate(b.UserId); err != nil {
		error = errors.New("userId must be a valid uuid")
	}
	return error
}

type IBidRepository interface {
	Create(ctx context.Context, bids []Bid) error
	FindByAuctionId(ctx context.Context, auctionId string) ([]Bid, error)
	FindWinningByAuctionId(ctx context.Context, auctionId string) (*Bid, error)
}
