package main

import (
	"concurrency/configs/database/mongodb"
	"concurrency/internal/infra/api/web/controller/auction_controller"
	"concurrency/internal/infra/api/web/controller/bid_controller"
	"concurrency/internal/infra/api/web/controller/user_controller"
	"concurrency/internal/infra/database"
	"concurrency/internal/usecase/auction_usecase"
	"concurrency/internal/usecase/bid_usecase"
	"concurrency/internal/usecase/user_usecase"
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func main() {
	if err := godotenv.Load("cmd/auctions/.env"); err != nil {
		log.Fatalln("Error on loading environment variables")
		return
	}
	ctx := context.Background()

	db, err := mongodb.NewMongoDBConnection(ctx)
	if err != nil {
		log.Fatalln("Error on connecting on database", err)
		return
	}

	router := gin.Default()

	userController, auctionController, bidController := initDependencies(db)

	router.GET("/auctions", auctionController.Find)
	router.GET("/auctions/:id", auctionController.FindById)
	router.GET("/auctions/:id/winning", auctionController.FindWinningBidByAuctionId)
	router.POST("/auctions", auctionController.Create)

	router.POST("/bids", bidController.Create)
	router.GET("/bids/auction/:id", bidController.FindByAuctionId)

	router.GET("/users/:id", userController.FindById)

	router.Run(":" + os.Getenv("APP_PORT"))
}

func initDependencies(db *mongo.Database) (user *user_controller.UserController, auction *auction_controller.AuctionController, bid *bid_controller.BidController) {

	userRepository := database.NewUserRepository(db)
	auctionRepository := database.NewAuctionRepository(db)
	bidRepository := database.NewBidRepository(db, auctionRepository)

	findUserByIdUseCase := user_usecase.FindUserById{UserRepository: userRepository}

	userController := user_controller.NewUserController(&findUserByIdUseCase)

	createAuctionUseCase := auction_usecase.CreateAuction{AuctionRepository: auctionRepository}
	findAuctionByIdUseCase := auction_usecase.FindAuctionById{AuctionRepository: auctionRepository}
	findWinningByAuctionIdUseCase := auction_usecase.FindWinningBidByAuctionId{AuctionRepository: auctionRepository}
	findAuctionsUseCase := auction_usecase.FindAuctions{AuctionRepository: auctionRepository}

	auctionController := auction_controller.NewAuctionController(&createAuctionUseCase, &findAuctionByIdUseCase, &findWinningByAuctionIdUseCase, &findAuctionsUseCase)

	createBidUseCase := bid_usecase.CreateBid{BidRepository: bidRepository}
	findBidByAuctionId := bid_usecase.FindByAuctionId{BidRepository: bidRepository}
	bidController := bid_controller.NewBidController(&createBidUseCase, &findBidByAuctionId)

	return userController, auctionController, bidController
}
