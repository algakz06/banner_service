package http

import (
	"net/http"
	"strings"

	"github.com/algakz/banner_service/pkg/auth"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type AuthMiddleware struct {
	usecase auth.UseCase
}

func NewAuthMiddleware(usecase auth.UseCase) gin.HandlerFunc {
	return (&AuthMiddleware{
		usecase: usecase,
	}).Handle
}

func (m *AuthMiddleware) Handle(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if headerParts[0] != "Bearer" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	user, err := m.usecase.ParseToken(c.Request.Context(), headerParts[1])
	if err != nil {
		status := http.StatusInternalServerError
		if err == auth.ErrInvalidAccessToken {
			status = http.StatusUnauthorized
		}
    logrus.Errorf("error occured while parsing token: %s", err.Error())

		c.AbortWithStatus(status)
		return
	}

	c.Set(auth.CtxUserKey, user)
}