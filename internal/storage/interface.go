package storage

import "context"

type Storage interface {
	GenerateUploadURL(ctx context.Context, key string, contentType string) (string, error)
	ObjectExists(ctx context.Context, key string) (bool, int64, error)
	DeleteObject(ctx context.Context, key string) error
}
