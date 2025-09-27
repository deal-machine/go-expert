package integration

import (
	"concurrency/internal/entity"
	"concurrency/internal/infra/database"
	"context"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func startupMongoContainer(ctx context.Context) (testcontainers.Container, error) {
	containerRequest := testcontainers.ContainerRequest{
		Image:           "mongo",
		ExposedPorts:    []string{"27017/tcp"},
		WaitingFor:      wait.ForListeningPort("27017/tcp").WithStartupTimeout(60 * time.Second),
		AlwaysPullImage: false,
	}
	genericContainer := testcontainers.GenericContainerRequest{
		ContainerRequest: containerRequest,
		Started:          true,
	}
	return testcontainers.GenericContainer(ctx, genericContainer)
}

func TestCreateAndExpireAuction(t *testing.T) {
	ctx := context.Background()

	mongoC, err := startupMongoContainer(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer mongoC.Terminate(ctx)

	host, _ := mongoC.Host(ctx)
	port, _ := mongoC.MappedPort(ctx, "27017")
	uri := "mongodb://" + host + ":" + port.Port()

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		t.Fatal(err)
	}
	defer client.Disconnect(ctx)

	collection := client.Database("auctions_test").Collection("auctions")
	repo := &database.AuctionRepository{
		Collection:        collection,
		ExpirationAuction: 2 * time.Second,
	}

	auction, err := entity.NewAuction(
		"productName", "category_name", "description", entity.New,
	)
	if err != nil {
		t.Fatal(err)
	}

	if err := repo.Create(ctx, *auction); err != nil {
		t.Fatal(err)
	}

	// valida inserção
	var a entity.Auction
	err = collection.FindOne(ctx, bson.M{"_id": auction.ID}).Decode(&a)
	if err != nil {
		t.Fatal("Auction not found in DB:", err)
	}

	// espera a expiração da goroutine
	time.Sleep(3 * time.Second)

	err = collection.FindOne(ctx, bson.M{"_id": auction.ID}).Decode(&a)
	if err != nil {
		t.Fatal(err)
	}

	if a.Status != entity.Completed {
		t.Fatal("Expected COMPLETED, got", a.Status)
	}
}
