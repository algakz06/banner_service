package banner

import "errors"

var (
	ErrBannerAlreadyExists = errors.New(
		"banner with these feature_id and tag_ids already exists. did you mean update?",
	)
  ErrNoBannerFound = errors.New(
    "no banner found with got id",
  )
)
