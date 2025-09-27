package bid_usecase

import (
	"concurrency/internal/entity"
	"context"
	"log"
	"os"
	"strconv"
	"time"
)

type CreateBidInput struct {
	UserId    string  `json:"user_id"`
	AuctionId string  `json:"auction_id"`
	Amount    float64 `json:"amount"`
}
type CreateBidOutput struct {
	ID        string    `json:"id"`
	UserId    string    `json:"user_id"`
	AuctionId string    `json:"auction_id"`
	Amount    float64   `json:"amount"`
	Timestamp time.Time `json:"timestamp" time_format:"2006-01-02 15:04:05"`
}

type ICreateBidUseCase interface {
	Execute(ctx context.Context, input CreateBidInput) error
}

type CreateBid struct {
	BidRepository entity.IBidRepository

	timer               *time.Timer
	maxBatchSize        int
	batchInsertInterval time.Duration
	bidChannel          chan entity.Bid
}

func NewCreateBidUseCase(bidRepository entity.IBidRepository) ICreateBidUseCase {
	interval := getBatchInsertInterval()
	batchSize := getMaxBatchSize()

	bidUseCase := &CreateBid{
		BidRepository:       bidRepository,
		timer:               time.NewTimer(interval),
		maxBatchSize:        batchSize,
		batchInsertInterval: interval,
		bidChannel:          make(chan entity.Bid, batchSize),
	}

	bidUseCase.triggerCreateRoutine(context.Background())

	return bidUseCase
}

func getBatchInsertInterval() time.Duration {
	batchInsertInterval := os.Getenv("BATCH_INSERT_INTERVAL")
	duration, err := time.ParseDuration(batchInsertInterval)
	if err != nil {
		return 3 * time.Minute
	}
	return duration
}
func getMaxBatchSize() int {
	value, err := strconv.Atoi(os.Getenv("MAX_BATCH_SIZE"))
	if err != nil {
		return 10
	}
	return value
}

var bidBatch []entity.Bid

func (b *CreateBid) triggerCreateRoutine(ctx context.Context) {
	go func() {
		defer close(b.bidChannel)

		for {
			select {
			case bid, ok := <-b.bidChannel:
				if !ok {
					if len(bidBatch) > 0 {
						if err := b.BidRepository.Create(ctx, bidBatch); err != nil {
							log.Println("error on persist bidBatch", err)
						}
					}
					return
				}
				bidBatch = append(bidBatch, bid)
				if len(bidBatch) >= b.maxBatchSize {
					if err := b.BidRepository.Create(ctx, bidBatch); err != nil {
						log.Println("error on persist bidBatch", err)
					}
					bidBatch = nil
					b.timer.Reset(b.batchInsertInterval)
				}
			case <-b.timer.C:
				if err := b.BidRepository.Create(ctx, bidBatch); err != nil {
					log.Println("error on persist bidBatch", err)
				}
				bidBatch = nil
				b.timer.Reset(b.batchInsertInterval)
			}

		}
	}()
}

func (b *CreateBid) Execute(ctx context.Context, input CreateBidInput) error {
	bid, err := entity.NewBid(input.UserId, input.AuctionId, input.Amount)
	if err != nil {
		return err
	}
	b.bidChannel <- *bid

	return nil
}
