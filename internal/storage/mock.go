package storage

import (
	"context"
	"github.com/stretchr/testify/mock"
)

type ShortenerStorageMock struct {
	mock.Mock
}

func (m *ShortenerStorageMock) SaveShortLink(ctx context.Context, linkID string, url string) error {
	args := m.Called(linkID, url)
	return args.Error(0)
}

func (m *ShortenerStorageMock) GetSource(ctx context.Context, linkID string) (string, error) {
	args := m.Called(linkID)
	return args.String(0), args.Error(1)
}
