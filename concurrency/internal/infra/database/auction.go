package database

import (
	"concurrency/internal/entity"
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type AuctionEntityMongo struct {
	ID          string                  `bson:"_id"`
	ProductName string                  `bson:"product_name"`
	Category    string                  `bson:"category"`
	Description string                  `bson:"description"`
	Condition   entity.ProductCondition `bson:"condition"`
	Status      entity.AuctionStatus    `bson:"status"`
	Timestamp   int64                   `bson:"timestamp"`
}

type AuctionRepository struct {
	Collection        *mongo.Collection
	ExpirationAuction time.Duration
}

func NewAuctionRepository(db *mongo.Database) *AuctionRepository {
	expiration := getExpiration()
	return &AuctionRepository{
		Collection:        db.Collection("auctions"),
		ExpirationAuction: expiration,
	}
}

func getExpiration() time.Duration {
	expiration := os.Getenv("EXPIRATION_AUCTION")
	duration, err := time.ParseDuration(expiration)
	if err != nil {
		return 10 * time.Minute
	}
	return duration
}

func (ar *AuctionRepository) Create(ctx context.Context, auction entity.Auction) error {
	_, err := ar.Collection.InsertOne(ctx, AuctionEntityMongo{
		ID:          auction.ID,
		ProductName: auction.ProductName,
		Category:    auction.Category,
		Description: auction.Description,
		Condition:   auction.Condition,
		Status:      auction.Status,
		Timestamp:   auction.Timestamp.Unix(),
	})
	if err != nil {
		log.Println("error on InsertOne auction", err)
		return err
	}

	expCtx, cancel := context.WithCancel(context.Background())
	go func() {
		defer cancel()
		select {
		case <-time.After(ar.ExpirationAuction):
			filter := bson.M{"_id": auction.ID}
			update := bson.M{"$set": bson.M{"status": entity.Completed}}
			if _, err := ar.Collection.UpdateOne(expCtx, filter, update); err != nil {
				log.Println("Error on UpdateOne auction", err)
				return
			}
			log.Printf("Auction %v expired and was marked as COMPLETED\n", auction.ID)
		case <-expCtx.Done():
			log.Println("Context Timeout")
		}
	}()

	return nil
}

func (ar *AuctionRepository) FindById(ctx context.Context, id string) (*entity.Auction, error) {
	var auctionEntity AuctionEntityMongo
	err := ar.Collection.FindOne(ctx, bson.M{"_id": id}).Decode(&auctionEntity)
	if err != nil {
		return nil, err
	}
	return &entity.Auction{
		ID:          auctionEntity.ID,
		ProductName: auctionEntity.ProductName,
		Category:    auctionEntity.Category,
		Description: auctionEntity.Description,
		Condition:   auctionEntity.Condition,
		Status:      auctionEntity.Status,
		Timestamp:   time.Unix(auctionEntity.Timestamp, 0),
	}, nil
}

func (ar *AuctionRepository) FindBy(ctx context.Context, status entity.AuctionStatus, category, productName string) ([]entity.Auction, error) {
	filter := bson.M{}
	if status != 0 {
		filter["status"] = status
	}
	if category != "" {
		filter["category"] = category
	}
	if productName != "" {
		filter["product_name"] = bson.Regex{
			Pattern: productName,
			Options: "i",
		}
	}

	cur, err := ar.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var auctionsEntity []AuctionEntityMongo
	if err = cur.All(ctx, &auctionsEntity); err != nil {
		return nil, err
	}

	var auctions []entity.Auction
	for _, ae := range auctionsEntity {
		auctions = append(auctions, entity.Auction{
			ID:          ae.ID,
			ProductName: ae.ProductName,
			Category:    ae.Category,
			Description: ae.Description,
			Condition:   ae.Condition,
			Status:      ae.Status,
			Timestamp:   time.Unix(ae.Timestamp, 0),
		})
	}

	return auctions, nil
}
