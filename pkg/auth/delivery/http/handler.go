package http

import (
	"net/http"

	"github.com/algakz/banner_service/pkg/auth"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	useCase auth.UseCase
}

func NewHandler(useCase auth.UseCase) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

type signInput struct {
  Username string `json:"username" binding:"required"`
  Password string `json:"password" binding:"required"`
}

func (h *Handler) SignUp(c *gin.Context) {
	inp := new(signInput)

	if err := c.BindJSON(inp); err != nil {
		c.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}

	if err := h.useCase.SignUp(c.Request.Context(), inp.Username, inp.Password); err != nil {
    logrus.Fatalf("error while processing userCase.SignUp: %s", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

type signInResponse struct {
	Token string `json:"token"`
}

func (h *Handler) SignIn(c *gin.Context) {
	inp := new(signInput)

	if err := c.BindJSON(inp); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	token, err := h.useCase.SignIn(c.Request.Context(), inp.Username, inp.Password)
	if err != nil {
		if err == auth.ErrUserNotFound {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, signInResponse{Token: token})
}
