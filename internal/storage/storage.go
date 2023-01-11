package storage

import (
	"context"
	"errors"
)

type Storage struct {
	Shortener
}

func NewStorage(shortener Shortener) *Storage {
	return &Storage{Shortener: shortener}
}

type Shortener interface {
	SaveShortLink(ctx context.Context, linkID string, url string) error
	GetSource(ctx context.Context, linkID string) (string, error)
}

var (
	ErrNoSourceUrl       = errors.New("given linkID is unknown")
	ErrPartialSave       = errors.New("short link saved not in all storages")
	ErrSaveFailed        = errors.New("short link save failed")
	ErrLinkIDAlreadyUsed = errors.New("given linkID is already used by another non-equal link")
)
