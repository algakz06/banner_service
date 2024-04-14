package banner

import (
	"context"

	"github.com/algakz/banner_service/models"
)

type UseCase interface {
  	GetUserBanner(
		ctx context.Context,
		tag_ids []int,
		feature_id int,
	) (models.Banner, error)
	GetBanners(
		ctx context.Context,
    tag_id int,
    feature_id int,
    limit int,
    offset int,
	) ([]models.Banner, error)
	CreateBanner(
		ctx context.Context,
		banner *models.Banner,
		user *models.User,
	) (int, error)
	DeleteBanner(ctx context.Context, banner_id int) error
	UpdateBanner(ctx context.Context, banner *models.Banner) error
}
