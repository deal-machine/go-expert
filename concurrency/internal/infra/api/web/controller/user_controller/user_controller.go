package user_controller

import (
	"concurrency/internal/usecase/user_usecase"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserController struct {
	FindUserByIdUseCase user_usecase.IFindUserByIdUseCase
}

func NewUserController(findUserById user_usecase.IFindUserByIdUseCase) *UserController {
	return &UserController{
		FindUserByIdUseCase: findUserById,
	}
}

func (u *UserController) FindById(c *gin.Context) {
	userId := c.Param("id")
	if err := uuid.Validate(userId); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	input := user_usecase.FindUserUserInput{ID: userId}
	userOutput, err := u.FindUserByIdUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusCreated, userOutput)
}
