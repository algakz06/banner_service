package http

import (
	"net/http"

	"github.com/algakz/banner_service/models"
	"github.com/algakz/banner_service/pkg/auth"
	"github.com/algakz/banner_service/pkg/banner"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	useCase banner.UseCase
}

func NewHandler(useCase banner.UseCase) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

func (h *Handler) UserGet(ctx *gin.Context) {
	logrus.Debugf("hello World! from UserGet")
}

func (h *Handler) Get(ctx *gin.Context) {
}

type CreateBanner struct {
	TagIds    []int                  `json:"tag_ids"`
	FeatureId int                    `json:"feature_id"`
	Content   map[string]interface{} `json:"content"`
	IsActive  bool                   `json:"is_active"`
}

type createBannerResponse struct {
	BannerId int `json:"banner_id"`
}

func (h *Handler) Create(ctx *gin.Context) {
	inp := new(CreateBanner)
	if err := ctx.BindJSON(inp); err != nil {
		logrus.Errorf("error occured while processing json to CreateBanner struct: %s", err.Error())
		ctx.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}
	user := ctx.MustGet(auth.CtxUserKey).(*models.User)
	banner := &models.Banner{
		TagIds:    inp.TagIds,
		FeatureId: inp.FeatureId,
		Content:   inp.Content,
		IsActive:  inp.IsActive,
	}
	banner_id, err := h.useCase.CreateBanner(ctx, banner, user)
	if err != nil {
		logrus.Errorf("error returned from useCase.CreateBanner: %s", err.Error())
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}
  ctx.JSON(http.StatusCreated, &createBannerResponse{
    BannerId: banner_id,
  })
}

func (h *Handler) Delete(ctx *gin.Context) {
}

func (h *Handler) Update(ctx *gin.Context) {
}
