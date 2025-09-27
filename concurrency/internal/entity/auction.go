package entity

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Auction struct {
	ID          string
	ProductName string
	Category    string
	Description string
	Condition   ProductCondition
	Status      AuctionStatus
	Timestamp   time.Time
}

func NewAuction(productName, category, description string, condition ProductCondition) (*Auction, error) {
	auction := &Auction{
		ID:          uuid.NewString(),
		ProductName: productName,
		Category:    category,
		Description: description,
		Condition:   condition,
		Status:      Active,
		Timestamp:   time.Now(),
	}
	if err := auction.Validate(); err != nil {
		return nil, err
	}
	return auction, nil
}

func (a *Auction) Validate() error {
	var err error = nil
	if a.Condition != New && a.Condition != Used && a.Condition != Refurbished {
		err = errors.New("invalid condition, must to be New (0)")
	}
	if len(a.Category) < 1 || len(a.Description) < 1 || len(a.ProductName) < 1 {
		err = errors.New("invalid category, description or productName must to be greater than 1 character")
	}
	if a.Status != Active {
		err = errors.New("invalid status, must to be Active (0) or Completed (1)")
	}
	return err
}

type ProductCondition int
type AuctionStatus int

const (
	Active AuctionStatus = iota
	Completed
)

const (
	New ProductCondition = iota
	Used
	Refurbished
)

type IAuctionRepository interface {
	Create(ctx context.Context, auction Auction) error
	FindById(ctx context.Context, id string) (*Auction, error)
	FindBy(ctx context.Context, status AuctionStatus, category, productName string) ([]Auction, error)
}
