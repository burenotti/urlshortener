package service

import (
	"context"
	"fmt"
	"github.com/burenotti/urlshortener/internal/storage"
	"hash/fnv"
	"time"
)

type ShortenerService struct {
	storage storage.Shortener
}

func (s ShortenerService) CreateShortLink(ctx context.Context, url string) (string, error) {
	linkID := s.generateShortLink(url)
	err := s.storage.SaveShortLink(ctx, linkID, url)
	return linkID, err
}

func (s ShortenerService) GetSource(ctx context.Context, linkID string) (string, error) {
	return s.storage.GetSource(ctx, linkID)
}

// Naively generate short link using url and timestamp
func (s ShortenerService) generateShortLink(url string) string {
	data := fmt.Sprintf("%s@%d", url, time.Now().Unix())
	hasher := fnv.New32a()
	_, _ = hasher.Write([]byte(data))
	return fmt.Sprintf("%x", hasher.Sum32())
}

func NewShortenerService(storage storage.Shortener) Shortener {
	return ShortenerService{storage: storage}
}
