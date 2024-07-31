package models

import "time"

type Banner struct {
	BannerId  int                    `json:"banner_id"  db:"id"`
	IsActive  bool                   `json:"is_active"  db:"is_active"`
	FeatureId int                    `json:"feature_id" db:"feature_id"`
	TagIds    []int                  `json:"tag_ids"    db:"tag_ids"`
	Content   map[string]interface{} `json:"content"    db:"content"`
	CreatedAt time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt time.Time              `json:"updated_at" db:"updated_at"`
}
