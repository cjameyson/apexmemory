package storage

import (
	"context"
	"errors"
	"io"
)

var ErrNotFound = errors.New("object not found")

// Storage is the interface for object storage backends.
type Storage interface {
	Put(ctx context.Context, key string, reader io.Reader, contentType string) error
	Get(ctx context.Context, key string) (io.ReadCloser, string, error)
	Delete(ctx context.Context, key string) error
}
