package service

import (
	"context"
	"fmt"
	"github.com/burenotti/urlshortener/internal/storage"
	"github.com/sirupsen/logrus"
	"hash/fnv"
	"time"
)

type ShortenerService struct {
	shortenerStore storage.Shortener
	statsStore     storage.Statistics
}

func (s ShortenerService) CreateShortLink(ctx context.Context, url string) (string, error) {
	linkID := s.generateShortLink(url)
	err := s.shortenerStore.SaveShortLink(ctx, linkID, url)
	return linkID, err
}

func (s ShortenerService) GetSource(ctx context.Context, linkID string) (string, error) {
	source, err := s.shortenerStore.GetSource(ctx, linkID)
	return source, err
}

func (s ShortenerService) GetSourceForRedirect(ctx context.Context, sessionID string, linkID string) (string, error) {
	source, err := s.GetSource(ctx, linkID)
	if err == nil {
		err = s.statsStore.RegisterRedirect(ctx, sessionID, linkID, time.Now())
		if err != nil {
			logrus.WithField("error", err).Error("registering redirect failed")
		}
		return source, nil

	}
	return source, err
}

// Naively generate short link using url and timestamp
func (s ShortenerService) generateShortLink(url string) string {
	data := fmt.Sprintf("%s@%d", url, time.Now().Unix())
	hasher := fnv.New32a()
	_, _ = hasher.Write([]byte(data))
	return fmt.Sprintf("%x", hasher.Sum32())
}

func NewShortenerService(storage storage.Shortener, stats storage.Statistics) Shortener {
	return ShortenerService{shortenerStore: storage, statsStore: stats}
}
