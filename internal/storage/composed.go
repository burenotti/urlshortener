package storage

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
)

type ComposedShortener struct {
	mainStorage Shortener
	cacheStore  Shortener
}

func NewComposedShortener(main Shortener, cache Shortener) *ComposedShortener {
	return &ComposedShortener{mainStorage: main, cacheStore: cache}
}

func (c ComposedShortener) SaveShortLink(ctx context.Context, linkID string, url string) error {
	logger := logrus.WithField("link_id", linkID)
	logger.Info("trying to save short link to main storage")
	err := c.mainStorage.SaveShortLink(ctx, linkID, url)
	if err != nil {
		logger.WithField("error", err.Error()).Error("saving to main storage failed")
		return ErrSaveFailed
	}
	err = c.cacheStore.SaveShortLink(ctx, linkID, url)
	if err != nil {
		logger.WithField("error", err.Error()).Error("saving to cache storage failed")
		return ErrPartialSave
	}
	return nil
}

func (c ComposedShortener) GetSource(ctx context.Context, linkID string) (string, error) {
	logger := logrus.WithField("link_id", linkID)
	storages := []Shortener{c.cacheStore, c.mainStorage}
	for i, store := range storages {
		srcUrl, err := store.GetSource(ctx, linkID)
		l := logger.WithField("storage", i)
		if errors.Is(err, ErrNoSourceUrl) {
			l.Info("link_id hasn't been found in store")
		} else if err != nil {

			l.WithField("error", err.Error()).Error("an unexpected error occurred on request to store")
		} else {
			l.WithField("url", srcUrl).Info("has found url in store")
			return srcUrl, nil
		}
	}
	return "", ErrNoSourceUrl
}
