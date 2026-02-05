package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type LocalStorage struct {
	basePath string
}

func NewLocalStorage(basePath string) (*LocalStorage, error) {
	if err := os.MkdirAll(basePath, 0o750); err != nil {
		return nil, fmt.Errorf("create storage directory: %w", err)
	}
	return &LocalStorage{basePath: basePath}, nil
}

type objectMeta struct {
	ContentType string `json:"content_type"`
}

// safePath resolves a storage key to an absolute path and verifies it remains
// within basePath. Returns an error if the key would escape the storage root.
func (s *LocalStorage) safePath(key string) (string, error) {
	joined := filepath.Join(s.basePath, key)
	absPath, err := filepath.Abs(joined)
	if err != nil {
		return "", fmt.Errorf("resolve path: %w", err)
	}
	absBase, err := filepath.Abs(s.basePath)
	if err != nil {
		return "", fmt.Errorf("resolve base: %w", err)
	}
	if !strings.HasPrefix(absPath, absBase+string(filepath.Separator)) {
		return "", fmt.Errorf("invalid key: path traversal detected")
	}
	return absPath, nil
}

func (s *LocalStorage) Put(_ context.Context, key string, reader io.Reader, contentType string) error {
	path, err := s.safePath(key)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
		return fmt.Errorf("create directory: %w", err)
	}

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer f.Close()

	if _, err := io.Copy(f, reader); err != nil {
		os.Remove(path)
		return fmt.Errorf("write file: %w", err)
	}

	meta := objectMeta{ContentType: contentType}
	metaBytes, _ := json.Marshal(meta)
	if err := os.WriteFile(path+".meta", metaBytes, 0o640); err != nil {
		os.Remove(path)
		return fmt.Errorf("write metadata: %w", err)
	}

	return nil
}

func (s *LocalStorage) Get(_ context.Context, key string) (io.ReadCloser, string, error) {
	path, err := s.safePath(key)
	if err != nil {
		return nil, "", err
	}

	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, "", ErrNotFound
		}
		return nil, "", fmt.Errorf("open file: %w", err)
	}

	var meta objectMeta
	metaBytes, err := os.ReadFile(path + ".meta")
	if err == nil {
		json.Unmarshal(metaBytes, &meta)
	}

	return f, meta.ContentType, nil
}

func (s *LocalStorage) Delete(_ context.Context, key string) error {
	path, err := s.safePath(key)
	if err != nil {
		return err
	}
	os.Remove(path + ".meta")

	if err := os.Remove(path); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("delete file: %w", err)
	}
	return nil
}
