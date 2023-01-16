package storage

import (
	"context"
	"errors"
	"time"
)

type Storage struct {
	Shortener
	Statistics
}

func NewStorage(shortener Shortener, stat Statistics) *Storage {
	return &Storage{
		Shortener:  shortener,
		Statistics: stat,
	}
}

type Shortener interface {
	SaveShortLink(ctx context.Context, linkID string, url string) error
	GetSource(ctx context.Context, linkID string) (string, error)
}

type Statistics interface {
	RegisterRedirect(ctx context.Context, linkID, sessionID string, at time.Time) error
}

var (
	ErrNoSourceUrl       = errors.New("given linkID is unknown")
	ErrPartialSave       = errors.New("short link saved not in all storages")
	ErrSaveFailed        = errors.New("short link save failed")
	ErrLinkIDAlreadyUsed = errors.New("given linkID is already used by another non-equal link")
)
