package storage

import (
	"bytes"
	"context"
	"io"
	"testing"
)

func TestLocalStorage_PutGetDelete(t *testing.T) {
	dir := t.TempDir()
	s, err := NewLocalStorage(dir)
	if err != nil {
		t.Fatalf("NewLocalStorage: %v", err)
	}

	ctx := context.Background()
	key := "assets/test-user/test-file"
	data := []byte("hello world")
	contentType := "text/plain"

	if err := s.Put(ctx, key, bytes.NewReader(data), contentType); err != nil {
		t.Fatalf("Put: %v", err)
	}

	reader, ct, err := s.Get(ctx, key)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	defer reader.Close()

	got, _ := io.ReadAll(reader)
	if !bytes.Equal(got, data) {
		t.Errorf("data mismatch: got %q, want %q", got, data)
	}
	if ct != contentType {
		t.Errorf("content type: got %q, want %q", ct, contentType)
	}

	if err := s.Delete(ctx, key); err != nil {
		t.Fatalf("Delete: %v", err)
	}

	_, _, err = s.Get(ctx, key)
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound after delete, got %v", err)
	}
}

func TestLocalStorage_GetNotFound(t *testing.T) {
	dir := t.TempDir()
	s, err := NewLocalStorage(dir)
	if err != nil {
		t.Fatalf("NewLocalStorage: %v", err)
	}

	_, _, err = s.Get(context.Background(), "nonexistent")
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestLocalStorage_DeleteIdempotent(t *testing.T) {
	dir := t.TempDir()
	s, err := NewLocalStorage(dir)
	if err != nil {
		t.Fatalf("NewLocalStorage: %v", err)
	}

	if err := s.Delete(context.Background(), "nonexistent"); err != nil {
		t.Errorf("expected nil error for idempotent delete, got %v", err)
	}
}
