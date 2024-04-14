package usecase

import (
	"context"

	"github.com/algakz/banner_service/models"
	"github.com/algakz/banner_service/pkg/banner"
)

type BannerUseCase struct {
	bannerRepo banner.Repository
}

func NewBannerUseCase(bannerRepo banner.Repository) *BannerUseCase {
	return &BannerUseCase{
		bannerRepo: bannerRepo,
	}
}

func (b *BannerUseCase) GetUserBanner(
	ctx context.Context,
	tags_id []int,
	feature_id int,
) ([]*models.Banner, error) {
	return nil, nil
}

func (b *BannerUseCase) GetBanners(
	ctx context.Context,
	tag_id int,
	feature_id int,
	limit int,
	offset int,
) ([]*models.Banner, error) {
	return nil, nil
}

func (b *BannerUseCase) CreateBanner(
	ctx context.Context,
	banner *models.Banner,
	user *models.User,
) (int, error) {
	banner_id, err := b.bannerRepo.CreateBanner(ctx, banner, user)
	if err != nil {
		return 0, err
	}
	return banner_id, nil
}

func (b *BannerUseCase) UpdateBanner(ctx context.Context, banner *models.Banner) error {
	return nil
}

func (b *BannerUseCase) DeleteBanner(ctx context.Context, banner_id int) error {
	err := b.bannerRepo.DeleteBanner(ctx, banner_id)
	if err != nil {
		return err
	}
	return nil
}
