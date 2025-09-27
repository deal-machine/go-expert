package auction_controller

import (
	"concurrency/internal/entity"
	"concurrency/internal/usecase/auction_usecase"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuctionController struct {
	CreateAuctionUseCase          auction_usecase.ICreateAuctionUseCase
	FindAuctionByIdUseCase        auction_usecase.IFindAuctionByIdUseCase
	FindWinningByAuctionIdUseCase auction_usecase.IFindWinningByAuctionIdUseCase
	FindAuctionsUseCase           auction_usecase.IFindAuctionsUseCase
}

func NewAuctionController(
	createAuctionUseCase auction_usecase.ICreateAuctionUseCase,
	findAuctionByIdUseCase auction_usecase.IFindAuctionByIdUseCase,
	findWinningByAuctionIdUseCase auction_usecase.IFindWinningByAuctionIdUseCase,
	findAuctionsUseCase auction_usecase.IFindAuctionsUseCase,
) *AuctionController {
	return &AuctionController{
		CreateAuctionUseCase:          createAuctionUseCase,
		FindAuctionByIdUseCase:        findAuctionByIdUseCase,
		FindWinningByAuctionIdUseCase: findWinningByAuctionIdUseCase,
		FindAuctionsUseCase:           findAuctionsUseCase,
	}
}

func (a *AuctionController) Create(c *gin.Context) {
	var input auction_usecase.CreateAuctionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Println("here", err)
		c.JSON(http.StatusBadRequest, err)
		return
	}
	output, err := a.CreateAuctionUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, output)
}

func (a *AuctionController) FindById(c *gin.Context) {
	auctionId := c.Param("id")
	if err := uuid.Validate(auctionId); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	input := auction_usecase.FindAuctionByIdInput{ID: auctionId}
	output, err := a.FindAuctionByIdUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, output)
}

func (a *AuctionController) FindWinningBidByAuctionId(c *gin.Context) {
	auctionId := c.Param("id")
	if err := uuid.Validate(auctionId); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	input := auction_usecase.FindWinningBidByAuctionIdInput{AuctionId: auctionId}
	output, err := a.FindWinningByAuctionIdUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, output)
}

func (a *AuctionController) Find(c *gin.Context) {
	category := c.Query("category")
	productName := c.Query("productName")

	status, err := strconv.Atoi(c.Query("status"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	input := auction_usecase.FindAuctionsInput{Status: entity.AuctionStatus(status), Category: category, ProductName: productName}

	output, err := a.FindAuctionsUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, output)
}
