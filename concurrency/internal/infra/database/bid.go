package database

import (
	"concurrency/internal/entity"
	"context"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type BidEntityMongo struct {
	ID        string  `bson:"_id"`
	UserId    string  `bson:"user_id"`
	AuctionId string  `bson:"auction_id"`
	Amount    float64 `bson:"amount"`
	Timestamp int64   `bson:"timestamp"`
}

type BidRepository struct {
	Collection        *mongo.Collection
	AuctionRepository *AuctionRepository
}

func NewBidRepository(db *mongo.Database, ar *AuctionRepository) *BidRepository {
	return &BidRepository{
		Collection:        db.Collection("bids"),
		AuctionRepository: ar,
	}
}

func (br *BidRepository) Create(ctx context.Context, bids []entity.Bid) error {
	var wg sync.WaitGroup
	var err error

	for _, bid := range bids {
		wg.Add(1)

		go func(b entity.Bid) {
			defer wg.Done()

			auction, err := br.AuctionRepository.FindById(ctx, b.AuctionId)
			if err != nil {
				log.Println("Error on get auction by id into BidRepository", err)
				return
			}
			if auction.Status != entity.Active {
				log.Println("AuctionStatus is not active", auction.Status)
				return
			}
			bidEntity := &BidEntityMongo{
				ID:        b.ID,
				UserId:    b.UserId,
				AuctionId: auction.ID,
				Amount:    b.Amount,
				Timestamp: b.Timestamp.Unix(),
			}

			if _, err = br.Collection.InsertOne(ctx, &bidEntity); err != nil {
				log.Println("Error on insert Bid", err)
				return
			}
		}(bid)
	}

	wg.Wait()
	return err
}

func (br *BidRepository) FindByAuctionId(ctx context.Context, auctionId string) ([]entity.Bid, error) {
	filter := bson.M{"auction_id": auctionId}

	var bidsEntity []BidEntityMongo
	cur, err := br.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	if err = cur.All(ctx, &bidsEntity); err != nil {
		return nil, err
	}
	var bids []entity.Bid

	for _, b := range bidsEntity {
		bids = append(bids, entity.Bid{
			ID:        b.ID,
			UserId:    b.UserId,
			AuctionId: b.AuctionId,
			Amount:    b.Amount,
			Timestamp: time.Unix(b.Timestamp, 0),
		})
	}

	return bids, nil
}

func (br *BidRepository) FindWinningByAuctionId(ctx context.Context, auctionId string) (*entity.Bid, error) {
	filter := bson.M{"auction_id": auctionId}
	opts := options.FindOne().SetSort(bson.D{{Key: "amount", Value: -1}})

	var bidEntity BidEntityMongo
	if err := br.Collection.FindOne(ctx, filter, opts).Decode(&bidEntity); err != nil {
		return nil, err
	}

	return &entity.Bid{
		ID:        bidEntity.ID,
		UserId:    bidEntity.UserId,
		AuctionId: bidEntity.AuctionId,
		Amount:    bidEntity.Amount,
		Timestamp: time.Unix(bidEntity.Timestamp, 0),
	}, nil
}
