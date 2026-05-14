package storage

import (
	"context"
	"errors"
)

var (
	ErrNotFound   = errors.New("Not found")
	ErrCodeExists = errors.New("Code exists")
)

type Storage interface {
	Save(ctx context.Context, code, longUrl string) error
	Get(ctx context.Context, code string) (string, error)
}

// type UrlItem struct {
// 	id         int64
// 	short_code string
// 	long_url   string
// 	created_at time.Time
// }
