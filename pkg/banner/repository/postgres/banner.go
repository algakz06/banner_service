package postgres

import (
	"context"
	"fmt"
	"sort"

	"github.com/algakz/banner_service/models"
	bn "github.com/algakz/banner_service/pkg/banner"
	"github.com/jackc/pgx/v5"
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
) ([]models.Banner, error) {
	if tag_id == 0 {
		banners, err := GetBannersByFeatureId(ctx, b.dbpool, feature_id, limit, offset)
		return banners, err
	}
	if feature_id == 0 {
		banners, err := GetBannersByTagId(ctx, b.dbpool, tag_id, limit, offset)
		return banners, err
	}
	banners, err := GetBannersByTagIdAndFeatureId(ctx, b.dbpool, tag_id, feature_id, limit, offset)
	return banners, err
}

func (b *BannerPostgres) CreateBanner(
	ctx context.Context,
	banner *models.Banner,
	user *models.User,
) (int, error) {
  sort.Ints(banner.TagIds)
  query := `
  SELECT 1
  FROM banner b
  JOIN banner_tag bt ON bt.banner_id = b.id
  WHERE b.feature_id = $1
  GROUP BY b.id
  HAVING ARRAY_AGG(bt.tag_id order by tag_id) = $2
  LIMIT 1
  `
  commandTag, _ := b.dbpool.Exec(ctx, query, banner.FeatureId, banner.TagIds)

  if commandTag.RowsAffected() == 1 {
    return 0, bn.ErrBannerAlreadyExists
  }
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
			return 0, fmt.Errorf(
				"expected one row to be affected, got %d",
				commandTag.RowsAffected(),
			)
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
		logrus.Errorf(
			"error occured while deleting banner by id=%s, error: %s",
			fmt.Sprint(banner_id),
			err.Error(),
		)
		return err
	}

	if commandTag.RowsAffected() != 1 {
		err = fmt.Errorf("expected one row to be affected, got %d", commandTag.RowsAffected())
		logrus.Error(err)
		return err
	}

	return nil
}

func GetBannerIdsByTagId(
	ctx context.Context,
	dbpool *pgxpool.Pool,
	tag_id int,
	limit int,
	offset int,
) ([]int, error) {
	getBannerIdsQuery := "SELECT banner_id FROM banner_tag WHERE tag_id=$1"
	rows, err := dbpool.Query(ctx, getBannerIdsQuery, tag_id)
	if err != nil {
		return nil, fmt.Errorf("error querying banner_ids: %w", err)
	}
	defer rows.Close()

	var banner_ids []int
	for rows.Next() {
		var banner_id int
		if err := rows.Scan(&banner_id); err != nil {
			return nil, fmt.Errorf("error scanning banner_id: %w", err)
		}
		banner_ids = append(banner_ids, banner_id)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error reading banner_ids: %w", err)
	}

	return banner_ids, nil
}

func GetBannersByTagId(
	ctx context.Context,
	dbpool *pgxpool.Pool,
	tag_id int,
	limit int,
	offset int,
) ([]models.Banner, error) {
	banner_ids, err := GetBannerIdsByTagId(ctx, dbpool, tag_id, limit, offset)
	if err != nil {
		return nil, err
	}
	query := `SELECT id, is_active, feature_id, content, created_at, updated_at, ARRAY_AGG(bt.tag_id) as tag_ids 
  FROM banner b 
  JOIN banner_tag bt ON bt.banner_id = b.id
  GROUP BY b.id
  HAVING b.id = ANY($1) AND b.version_number = 1`

	rows, err := dbpool.Query(ctx, query, banner_ids)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	banners, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Banner])
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return banners, err
}

func GetBannersByFeatureId(
	ctx context.Context,
	dbpool *pgxpool.Pool,
	feature_id int,
	limit int,
	offset int,
) ([]models.Banner, error) {
	query := `SELECT id, is_active, feature_id, content, created_at, updated_at, ARRAY_AGG(bt.tag_id) as tag_ids
  FROM banner b
  JOIN banner_tag bt ON bt.banner_id = b.id
  GROUP BY b.id
  HAVING b.feature_id = $1 AND b.version_number = 1
`
	rows, err := dbpool.Query(ctx, query, feature_id)
	if err != nil {
		logrus.Errorf(
			"error occured while processing query: %s with 1=%s",
			query,
			fmt.Sprint(feature_id),
		)
		return nil, err
	}
  logrus.Infof("collecting rows started")
	banners, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Banner])
  logrus.Infof("collecting rows ended")
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
  logrus.Infof("banners was collected successfully")
	return banners, err
}

func GetBannersByTagIdAndFeatureId(
	ctx context.Context,
	dbpool *pgxpool.Pool,
	tag_id int,
	feature_id int,
	limit int,
	offset int,
) ([]models.Banner, error) {
	banner_ids, err := GetBannerIdsByTagId(ctx, dbpool, tag_id, limit, offset)
	if err != nil {
		return nil, err
	}
	query := `SELECT id, is_active, feature_id, content, created_at, updated_at, ARRAY_AGG(bt.tag_id) as tag_ids 
  FROM banner b 
  JOIN banner_tag bt ON bt.banner_id = b.id
  GROUP BY b.id
  HAVING b.id = ANY($1) AND b.feature_id = $2 AND b.version_number = 1`

	rows, err := dbpool.Query(ctx, query, banner_ids, feature_id)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	banners, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Banner])
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return banners, err
}
