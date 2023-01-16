package service

import (
	"context"
	"github.com/burenotti/urlshortener/internal/storage"
	"time"
)

type Service struct {
	Shortener
	Statistics
}

func NewService(storage *storage.Storage) *Service {
	return &Service{
		Shortener:  NewShortenerService(storage.Shortener, storage.Statistics),
		Statistics: NewStatisticsService(storage.Statistics),
	}
}

type Shortener interface {
	CreateShortLink(ctx context.Context, url string) (string, error)
	GetSourceForRedirect(ctx context.Context, sessionID string, linkID string) (string, error)
	GetSource(ctx context.Context, linkID string) (string, error)
}

type Statistics interface {
	RegisterRedirect(ctx context.Context, linkID, sessionID string, at time.Time) error
}
