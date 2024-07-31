package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/algakz/banner_service/models"
	"github.com/algakz/banner_service/pkg/auth"
	bn "github.com/algakz/banner_service/pkg/banner"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	useCase bn.UseCase
}

func NewHandler(useCase bn.UseCase) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

func (h *Handler) UserGet(ctx *gin.Context) {
	user := ctx.MustGet(auth.CtxUserKey).(*models.User)
	q_tag_ids := ""
	q_feature_id := ""

	q_tag_ids, _ = ctx.GetQuery("tag_ids")
	q_feature_id, _ = ctx.GetQuery("feature_id")

	if q_tag_ids == "" || q_feature_id == "" {
		err := fmt.Errorf("tag_id and feature_id not found")
		logrus.Error(err)
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	feature_id, err := strconv.Atoi(q_feature_id)
	if err != nil {
		logrus.Error(err)
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	var tag_ids []int
	err = json.Unmarshal([]byte(q_tag_ids), &tag_ids)
	if err != nil {
		logrus.Error(err)
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	banner, err := h.useCase.GetUserBanner(ctx, tag_ids, feature_id)
	if user.Role == "user" {
		if !banner.IsActive {
			ctx.AbortWithStatus(http.StatusNotFound)
			return
		}
	}
	ctx.JSON(http.StatusOK, banner.Content)
}

func (h *Handler) Get(ctx *gin.Context) {
	user := ctx.MustGet(auth.CtxUserKey).(*models.User)
	if user.Role != "admin" {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}
	q_tag_id := ""
	q_feature_id := ""
	q_limit := ""
	q_offset := ""

	q_tag_id, _ = ctx.GetQuery("tag_id")
	q_feature_id, _ = ctx.GetQuery("feature_id")

	if q_tag_id == "" && q_feature_id == "" {
		err := fmt.Errorf("tag_id and feature_id not found. i should put at least one of them")
		logrus.Error(err)
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	tag_id := 0
	var err error
	if q_tag_id != "" {
		tag_id, err = strconv.Atoi(q_tag_id)
		if err != nil {
			logrus.Error(err)
			ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}
	}

	feature_id := 0
	if q_feature_id != "" {
		feature_id, err = strconv.Atoi(q_feature_id)
		if err != nil {
			logrus.Error(err)
			ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}
	}

	q_limit, _ = ctx.GetQuery("limit")
	q_offset, _ = ctx.GetQuery("offset")

	var limit, offset int
	if q_limit == "" {
		limit = 50
	} else {
		limit, err = strconv.Atoi(q_limit)
		if err != nil {
			logrus.Error(err)
			ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}
	}

	if q_offset == "" {
		offset = 0
	} else {
		offset, err = strconv.Atoi(q_offset)
		if err != nil {
			logrus.Error(err)
			ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}
	}

	banner_list, err := h.useCase.GetBanners(ctx, tag_id, feature_id, limit, offset)
	if err != nil {
		logrus.Error(err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, banner_list)
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
		if err == bn.ErrBannerAlreadyExists {
			ctx.AbortWithError(http.StatusConflict, err)
		}
	}
	if err != nil {
		logrus.Errorf("error returned from useCase.CreateBanner: %s", err.Error())
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ctx.JSON(http.StatusCreated, &createBannerResponse{
		BannerId: banner_id,
	})
}

func (h *Handler) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	user := ctx.MustGet(auth.CtxUserKey).(*models.User)
	if user.Role != "admin" {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}
	banner_id, err := strconv.Atoi(id)
	if err != nil {
		logrus.Error(err)
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	err = h.useCase.DeleteBanner(ctx, banner_id)
	if err != nil {
		if strings.HasPrefix(err.Error(), "error") {
			logrus.Error(err)
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		} else {
			logrus.Error(err)
			ctx.AbortWithStatus(http.StatusNotFound)
			return
		}
	}
	ctx.AbortWithStatus(http.StatusNoContent)
}

type UpdateBanner struct {
	BannerId  int                    `json:"banner_id"`
	TagIds    []int                  `json:"tag_ids"`
	FeatureId int                    `json:"feature_id"`
	Content   map[string]interface{} `json:"content"`
	IsActive  bool                   `json:"is_active"`
}

func (h *Handler) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	banner_id, err := strconv.Atoi(id)
	if err != nil {
		logrus.Errorf("error occured converting param id to int: %s", err.Error())
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	inp := new(UpdateBanner)
	if err := ctx.BindJSON(inp); err != nil {
		logrus.Errorf("error occured while processing json to UpdateBanner struct: %s", err.Error())
		ctx.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}
	user := ctx.MustGet(auth.CtxUserKey).(*models.User)
	if user.Role != "admin" {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}
	banner := &models.Banner{
		BannerId:  banner_id,
		TagIds:    inp.TagIds,
		FeatureId: inp.FeatureId,
		Content:   inp.Content,
		IsActive:  inp.IsActive,
	}
	err = h.useCase.UpdateBanner(ctx, banner)
	if err == bn.ErrNoBannerFound {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ctx.AbortWithStatus(http.StatusOK)
	return
}
