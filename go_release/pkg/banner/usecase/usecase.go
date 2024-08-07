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
	tag_ids []int,
	feature_id int,
) (models.Banner, error) {
	banner, err := b.bannerRepo.GetUserBanner(ctx, tag_ids, feature_id)
	return banner, err
}

func (b *BannerUseCase) GetBanners(
	ctx context.Context,
	tag_id int,
	feature_id int,
	limit int,
	offset int,
) ([]models.Banner, error) {
	banners, err := b.bannerRepo.GetBanners(ctx, tag_id, feature_id, limit, offset)
	return banners, err
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
	err := b.bannerRepo.UpdateBanner(ctx, banner)
	return err
}

func (b *BannerUseCase) DeleteBanner(ctx context.Context, banner_id int) error {
	err := b.bannerRepo.DeleteBanner(ctx, banner_id)
	if err != nil {
		return err
	}
	return nil
}
