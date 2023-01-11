package service

import (
	"context"
	"github.com/burenotti/urlshortener/internal/storage"
)

type Service struct {
	Shortener
}

func NewService(storage storage.Shortener) *Service {
	return &Service{
		Shortener: NewShortenerService(storage),
	}
}

type Shortener interface {
	CreateShortLink(ctx context.Context, url string) (string, error)
	GetSource(ctx context.Context, linkID string) (string, error)
}
