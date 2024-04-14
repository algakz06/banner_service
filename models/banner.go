package models

type Banner struct {
	BannerId  string                 `json:"banner_id"`
	IsActive  bool                   `json:"is_active"`
	FeatureId int                    `json:"feature_id"`
	TagIds    []int                  `json:"tag_ids"`
	Content   map[string]interface{} `json:"content"`
	CreatedAt string                 `json:"created_at"`
	UpdatedAt string                 `json:"updated_at"`
}
