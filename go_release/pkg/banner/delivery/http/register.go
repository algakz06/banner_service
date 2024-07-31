package http

import (
	"github.com/algakz/banner_service/pkg/banner"
	"github.com/gin-gonic/gin"
)

func RegisterHTTPEndpoints(router *gin.RouterGroup, uc banner.UseCase) {
	h := NewHandler(uc)

	bannerEndpoints := router.Group("")
	{
		bannerEndpoints.GET("/user_banner", h.UserGet)
		banner := router.Group("banner")
		{
			banner.GET("", h.Get)
			banner.POST("", h.Create)
      banner.PATCH("/:id", h.Update)
      banner.DELETE("/:id", h.Delete)
		}
	}
}
