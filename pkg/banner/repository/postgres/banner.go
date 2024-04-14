package postgres

import (
	"context"
	"fmt"

	"github.com/algakz/banner_service/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type BannerPostgres struct {
	dbpool *pgxpool.Pool
}

func NewBannerPostgres(dbpool *pgxpool.Pool) *BannerPostgres {
	return &BannerPostgres{dbpool: dbpool}
}

func (b *BannerPostgres) GetUserBanner(
	ctx context.Context,
	tags_id []int,
	feature_id int,
) ([]*models.Banner, error) {
	return nil, nil
}

func (b *BannerPostgres) GetBanners(
	ctx context.Context,
	tag_id int,
	feature_id int,
	limit int,
	offset int,
) ([]*models.Banner, error) {
	return nil, nil
}

func (b *BannerPostgres) CreateBanner(
	ctx context.Context,
	banner *models.Banner,
	user *models.User,
) (int, error) {
	tx, err := b.dbpool.Begin(ctx)
	if err != nil {
		logrus.Errorf("error occured while starting tx: %s", err.Error())
		return 0, err
	}

	var banner_id int

	createBannerQuery := "INSERT INTO banner (feature_id, content, is_active) VALUES ($1, $2, $3) RETURNING id"
	err = tx.QueryRow(ctx, createBannerQuery, banner.FeatureId, banner.Content, banner.IsActive).
		Scan(&banner_id)
	if err != nil {
		logrus.Errorf("error occured while inserting to banner table: %s", err.Error())
		tx.Rollback(ctx)
		return 0, err
	}

	addTagsQuery := fmt.Sprintf(
		"INSERT INTO banner_tag (banner_id, tag_id) VALUES (%s, $1)",
		fmt.Sprint(banner_id),
	)
	for _, tag_id := range banner.TagIds {
		commandTag, err := tx.Exec(ctx, addTagsQuery, tag_id)
		if err != nil {
			logrus.Errorf(
				"error occured while inserting to banner_tag table with tag_id=%s: %s",
				fmt.Sprint(tag_id),
				err.Error(),
			)
			tx.Rollback(ctx)
			return 0, err
		}
		if commandTag.RowsAffected() != 1 {
			logrus.Errorf("expected one row to be affected, got %d", commandTag.RowsAffected())
			tx.Rollback(ctx)
			return 0, fmt.Errorf("expected one row to be affected, got %d", commandTag.RowsAffected())
		}
	}
  err = tx.Commit(ctx)
  if err != nil {
    logrus.Errorf("error occured while commit tx: %s", err.Error())
    tx.Rollback(ctx)
    return 0, err
  }

	return banner_id, nil
}

func (b *BannerPostgres) UpdateBanner(ctx context.Context, banner *models.Banner) error {
	return nil
}

func (b *BannerPostgres) DeleteBanner(ctx context.Context, banner_id int) error {
  query := "DELETE FROM banner WHERE id = $1"
  commandTag, err := b.dbpool.Exec(ctx, query, banner_id)

  if err != nil {
    logrus.Errorf("error occured while deleting banner by id=%s, error: %s", fmt.Sprint(banner_id), err.Error())
    return err
  }

  if commandTag.RowsAffected() != 1 {
    err = fmt.Errorf("expected one row to be affected, got %d", commandTag.RowsAffected())
    logrus.Error(err)
    return err
  }

  return nil
}
