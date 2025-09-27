package bid_controller

import (
	"concurrency/internal/usecase/bid_usecase"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BidController struct {
	CreateBidUseCase   bid_usecase.ICreateBidUseCase
	FindBidByAuctionId bid_usecase.IFindByAuctionIdUseCase
}

func NewBidController(
	createBidUseCase bid_usecase.ICreateBidUseCase,
	findBidByAuctionId bid_usecase.IFindByAuctionIdUseCase,
) *BidController {
	return &BidController{
		CreateBidUseCase:   createBidUseCase,
		FindBidByAuctionId: findBidByAuctionId,
	}
}

func (b *BidController) Create(c *gin.Context) {
	var input bid_usecase.CreateBidInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	if err := b.CreateBidUseCase.Execute(c.Request.Context(), input); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.Status(http.StatusCreated)
}

func (b *BidController) FindByAuctionId(c *gin.Context) {
	auctionId := c.Param("id")
	if err := uuid.Validate(auctionId); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	input := bid_usecase.FindByAuctionIdInput{AuctionId: auctionId}
	output, err := b.FindBidByAuctionId.Execute(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, output)
}
